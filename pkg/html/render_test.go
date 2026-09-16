package html

import (
	"strings"
	"testing"

	"mdreader/pkg/markdown"
)

func TestRenderHTML(t *testing.T) {
	input := `# Title

Paragraph with [link](https://example.com) and ` + "`code`" + `.

` + "```bash\necho 123\n```" + `

| Name | Role |
| :--- | :--- |
| Bob  | Dev  |
`
	doc := markdown.Parse(input)
	output := RenderHTML(doc, "Test Doc")

	if !strings.Contains(output, "<h1 id=\"title\">Title</h1>") {
		t.Errorf("HTML missing h1: %s", output)
	}
	if !strings.Contains(output, "<a href=\"https://example.com\">link</a>") {
		t.Errorf("HTML missing link: %s", output)
	}
	if !strings.Contains(output, "<code>code</code>") {
		t.Errorf("HTML missing code span: %s", output)
	}
	if !strings.Contains(output, "<code class=\"language-bash\">echo 123</code>") {
		t.Errorf("HTML missing code block: %s", output)
	}
	if !strings.Contains(output, "<table>") || !strings.Contains(output, "Bob</td>") {
		t.Errorf("HTML missing table elements: %s", output)
	}
}
