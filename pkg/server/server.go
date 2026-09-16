package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mdreader/pkg/html"
	"mdreader/pkg/markdown"
	"mdreader/pkg/pdf"
)

// StartServer launches an HTTP server displaying the live-reloaded markdown document.
func StartServer(filePath string, port int) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	fileName := filepath.Base(absPath)
	docBaseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading %s: %v", fileName, err), http.StatusInternalServerError)
			return
		}

		doc := markdown.Parse(string(data))
		htmlContent := html.RenderHTML(doc, fileName)

		// Floating top toolbar for easy PDF printing and copying
		toolbar := `
  <div class="no-print" style="position: sticky; top: 0; z-index: 1000; background: rgba(22, 27, 34, 0.92); backdrop-filter: blur(8px); border-bottom: 1px solid #30363d; padding: 10px 18px; margin: -2rem -1rem 1.5rem -1rem; display: flex; align-items: center; justify-content: space-between; font-size: 13px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);">
    <div style="font-weight: 600; display: flex; align-items: center; gap: 8px;">
      <span style="font-size: 16px;">📖</span> <span>` + fileName + `</span>
      <span style="background: #238636; color: white; padding: 2px 6px; border-radius: 10px; font-size: 11px;">Live</span>
    </div>
    <div style="display: flex; gap: 8px;">
      <button onclick="window.print()" style="background: #21262d; border: 1px solid #30363d; color: #c9d1d9; padding: 5px 12px; border-radius: 6px; cursor: pointer; display: flex; align-items: center; gap: 6px; font-size: 13px;">
        🖨️ Print / Save PDF
      </button>
      <a href="/__export-pdf" download style="background: #238636; border: 1px solid rgba(240,246,252,0.1); color: #ffffff; padding: 5px 12px; border-radius: 6px; text-decoration: none; cursor: pointer; display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 500;">
        ⬇️ Download PDF
      </a>
      <button onclick="copyDocument()" id="copyBtn" style="background: #21262d; border: 1px solid #30363d; color: #c9d1d9; padding: 5px 12px; border-radius: 6px; cursor: pointer; display: flex; align-items: center; gap: 6px; font-size: 13px;">
        📋 Copy Text
      </button>
    </div>
  </div>
`

		// Inject live reload and copy scripts before </body>
		scripts := `
  <script>
    let lastMod = 0;
    async function checkReload() {
      try {
        let res = await fetch('/__modtime');
        let mod = await res.text();
        if (lastMod && mod !== lastMod) {
          window.location.reload();
        }
        lastMod = mod;
      } catch (e) {}
    }
    setInterval(checkReload, 800);
    checkReload();

    async function copyDocument() {
      const btn = document.getElementById('copyBtn');
      try {
        const text = document.querySelector('.container').innerText;
        await navigator.clipboard.writeText(text);
        btn.innerText = '✅ Copied!';
        setTimeout(() => { btn.innerText = '📋 Copy Text'; }, 2000);
      } catch (e) {
        btn.innerText = '❌ Failed';
        setTimeout(() => { btn.innerText = '📋 Copy Text'; }, 2000);
      }
    }
  </script>
</body>`

		pageWithToolbar := strings.Replace(htmlContent, "<div class=\"container\">", "<div class=\"container\">\n"+toolbar, 1)
		finalHTML := strings.Replace(pageWithToolbar, "</body>", scripts, 1)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(finalHTML))
	})

	http.HandleFunc("/__export-pdf", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(absPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
			return
		}

		doc := markdown.Parse(string(data))
		htmlContent := html.RenderHTML(doc, fileName)

		tmpPDF, err := os.CreateTemp("", "mdread-*.pdf")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
			return
		}
		tmpPDFPath := tmpPDF.Name()
		tmpPDF.Close()
		defer os.Remove(tmpPDFPath)

		if err := pdf.GeneratePDF(htmlContent, tmpPDFPath); err != nil {
			http.Error(w, fmt.Sprintf("Error generating PDF: %v", err), http.StatusInternalServerError)
			return
		}

		pdfBytes, err := os.ReadFile(tmpPDFPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading generated PDF: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", docBaseName))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
		w.Write(pdfBytes)
	})

	http.HandleFunc("/__modtime", func(w http.ResponseWriter, r *http.Request) {
		info, err := os.Stat(absPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%d", info.ModTime().UnixNano())
	})

	fmt.Printf("\n🚀 Live Markdown Preview Server running at: http://%s/\n", addr)
	fmt.Printf("📄 Watching: %s\n", absPath)
	fmt.Printf("💡 Use the toolbar to 'Save as PDF' or 'Download PDF' anytime!\n")
	fmt.Printf("Press Ctrl+C to stop.\n\n")

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return srv.ListenAndServe()
}
