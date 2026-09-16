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
)

// StartServer launches an HTTP server displaying the live-reloaded markdown document.
func StartServer(filePath string, port int) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	fileName := filepath.Base(absPath)

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

		// Inject live reload script before </body>
		reloadScript := `
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
  </script>
</body>`
		finalHTML := strings.Replace(htmlContent, "</body>", reloadScript, 1)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(finalHTML))
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
	fmt.Printf("💡 Edit and save your file; the browser will refresh automatically.\n")
	fmt.Printf("Press Ctrl+C to stop.\n\n")

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.ListenAndServe()
}
