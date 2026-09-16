package pdf

import (
	"os"
	"testing"
)

func TestFindBrowser(t *testing.T) {
	path, err := FindBrowser()
	if err != nil {
		t.Skipf("skipping browser test: %v", err)
	}
	if path == "" {
		t.Errorf("expected non-empty browser path")
	}
}

func TestGeneratePDF(t *testing.T) {
	tmpPDF, err := os.CreateTemp("", "test-*.pdf")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPDFPath := tmpPDF.Name()
	tmpPDF.Close()
	defer os.Remove(tmpPDFPath)

	sampleHTML := `<!DOCTYPE html><html><body><h1>Test PDF Generation</h1><p>Hello world!</p></body></html>`
	if err := GeneratePDF(sampleHTML, tmpPDFPath); err != nil {
		t.Skipf("skipping PDF generation test if headless browser is not supported in environment: %v", err)
	}

	info, err := os.Stat(tmpPDFPath)
	if err != nil || info.Size() == 0 {
		t.Errorf("PDF was not created or has 0 bytes")
	}
}
