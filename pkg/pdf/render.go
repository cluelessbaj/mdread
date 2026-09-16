package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var candidateBrowsers = []string{
	"chromium",
	"google-chrome",
	"google-chrome-stable",
	"brave-browser",
	"brave",
	"microsoft-edge",
	"microsoft-edge-stable",
	"chromium-browser",
}

// FindBrowser locates a chromium-based browser capable of headless PDF printing.
func FindBrowser() (string, error) {
	for _, name := range candidateBrowsers {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no supported browser found (looked for: chromium, google-chrome, brave, microsoft-edge)")
}

// GeneratePDF converts styled HTML content into a PDF file at outputPath.
func GeneratePDF(htmlContent string, outputPath string) error {
	browserPath, err := FindBrowser()
	if err != nil {
		return fmt.Errorf("cannot generate PDF: %w\nInstall chromium or google-chrome, or use 'mdread --serve' and print to PDF from your browser", err)
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	// Create temporary HTML file
	tmpFile, err := os.CreateTemp("", "mdread-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary HTML file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temporary HTML: %w", err)
	}
	tmpFile.Close()

	// Run headless browser to print PDF
	cmd := exec.Command(
		browserPath,
		"--headless",
		"--disable-gpu",
		"--no-pdf-header-footer",
		fmt.Sprintf("--print-to-pdf=%s", absOutput),
		tmpPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("browser failed to generate PDF: %w (output: %s)", err, string(out))
	}

	// Verify output file exists and is not empty
	stat, err := os.Stat(absOutput)
	if err != nil || stat.Size() == 0 {
		return fmt.Errorf("PDF was not created or is empty: %w", err)
	}

	return nil
}
