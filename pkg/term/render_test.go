package term

import (
	"strings"
	"testing"

	"mdreader/pkg/markdown"
)

func TestRenderTerminalBasic(t *testing.T) {
	input := `# Test Document

This is a paragraph with **bold** and ` + "`code`" + `.

> A quote to remember

` + "```go\nfmt.Println(\"test\")\n```" + `

- [x] Done task
- [ ] Todo task

| Col A | Col B |
| :--- | ---: |
| 1 | 2 |
`
	doc := markdown.Parse(input)
	output := RenderTerminal(doc, TerminalOptions{Width: 80, UseColor: false})

	if !strings.Contains(output, "Test Document") {
		t.Errorf("rendered output missing header: %s", output)
	}
	if !strings.Contains(output, "bold") {
		t.Errorf("rendered output missing bold: %s", output)
	}
	if !strings.Contains(output, "Col A") {
		t.Errorf("rendered output missing table: %s", output)
	}
	if !strings.Contains(output, "[✓]") {
		t.Errorf("rendered output missing checked task: %s", output)
	}
	if !strings.Contains(output, "[ ]") {
		t.Errorf("rendered output missing unchecked task: %s", output)
	}
}

func TestVisibleLen(t *testing.T) {
	styled := Style("Hello", Bold, FgRed)
	if VisibleLen(styled) != 5 {
		t.Errorf("expected visible length 5, got %d", VisibleLen(styled))
	}
}
