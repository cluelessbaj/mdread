package markdown

import (
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	input := `---
title: My Documentation
author: Alice
date: 2026-09-16
---

# Main Title

This is content.`

	doc := Parse(input)
	if doc.FrontMatter["title"] != "My Documentation" {
		t.Fatalf("expected title 'My Documentation', got %q", doc.FrontMatter["title"])
	}
	if doc.FrontMatter["author"] != "Alice" {
		t.Fatalf("expected author 'Alice', got %q", doc.FrontMatter["author"])
	}
	if len(doc.Children) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(doc.Children))
	}
	h, ok := doc.Children[0].(*Heading)
	if !ok || h.Level != 1 || h.RawText != "Main Title" {
		t.Fatalf("expected H1 'Main Title', got %+v", doc.Children[0])
	}
}

func TestParseHeadingsAndTOC(t *testing.T) {
	input := `# Header 1
Some intro.
## Subheader 1.1
Content.
### Deep Header 1.1.1
More content.
`
	doc := Parse(input)
	toc := ExtractTOC(doc)

	if len(toc) != 3 {
		t.Fatalf("expected 3 TOC entries, got %d", len(toc))
	}
	if toc[0].Title != "Header 1" || toc[0].Level != 1 {
		t.Errorf("unexpected toc[0]: %+v", toc[0])
	}
	if toc[1].Title != "Subheader 1.1" || toc[1].Level != 2 {
		t.Errorf("unexpected toc[1]: %+v", toc[1])
	}
	if toc[2].Title != "Deep Header 1.1.1" || toc[2].Level != 3 {
		t.Errorf("unexpected toc[2]: %+v", toc[2])
	}
}

func TestParseCodeBlock(t *testing.T) {
	input := "```go\nfunc main() {\n    println(\"hello\")\n}\n```"
	doc := Parse(input)

	if len(doc.Children) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Children))
	}
	cb, ok := doc.Children[0].(*CodeBlock)
	if !ok {
		t.Fatalf("expected *CodeBlock, got %T", doc.Children[0])
	}
	if cb.Language != "go" {
		t.Errorf("expected language 'go', got %q", cb.Language)
	}
	expectedCode := "func main() {\n    println(\"hello\")\n}"
	if cb.Content != expectedCode {
		t.Errorf("expected code:\n%s\ngot:\n%s", expectedCode, cb.Content)
	}
}

func TestParseTable(t *testing.T) {
	input := `| Header 1 | Header 2 | Header 3 |
| :--- | :---: | ---: |
| Left | Center | Right |
| Val 1 | Val 2 | Val 3 |
`
	doc := Parse(input)
	if len(doc.Children) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Children))
	}
	tbl, ok := doc.Children[0].(*Table)
	if !ok {
		t.Fatalf("expected *Table, got %T", doc.Children[0])
	}
	if len(tbl.Headers) != 3 {
		t.Fatalf("expected 3 headers, got %d", len(tbl.Headers))
	}
	if tbl.Alignments[0] != AlignLeft || tbl.Alignments[1] != AlignCenter || tbl.Alignments[2] != AlignRight {
		t.Errorf("unexpected alignments: %+v", tbl.Alignments)
	}
	if len(tbl.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows))
	}
	if tbl.Rows[0][0].RawText != "Left" || tbl.Rows[0][1].RawText != "Center" || tbl.Rows[0][2].RawText != "Right" {
		t.Errorf("unexpected row 0 values: %+v", tbl.Rows[0])
	}
}

func TestParseListsAndTasks(t *testing.T) {
	input := `- [ ] Task 1
- [x] Task 2 completed
- Regular item 3
`
	doc := Parse(input)
	if len(doc.Children) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Children))
	}
	list, ok := doc.Children[0].(*List)
	if !ok {
		t.Fatalf("expected *List, got %T", doc.Children[0])
	}
	if len(list.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list.Items))
	}
	if !list.Items[0].IsTask || list.Items[0].Checked {
		t.Errorf("item 0 should be unchecked task, got %+v", list.Items[0])
	}
	if !list.Items[1].IsTask || !list.Items[1].Checked {
		t.Errorf("item 1 should be checked task, got %+v", list.Items[1])
	}
	if list.Items[2].IsTask {
		t.Errorf("item 2 should not be task, got %+v", list.Items[2])
	}
}

func TestParseInlines(t *testing.T) {
	text := "Hello **bold** and *italic* and `code` and [link](https://golang.org) and ~~strike~~"
	inlines := ParseInlines(text)

	var types []InlineType
	for _, in := range inlines {
		types = append(types, in.InlineType())
	}

	expectedCount := 10
	if len(inlines) != expectedCount {
		t.Fatalf("expected %d inlines, got %d (%v)", expectedCount, len(inlines), types)
	}

	b, ok := inlines[1].(StrongInline)
	if !ok || InlinesToPlainText(b.Children) != "bold" {
		t.Errorf("expected bold 'bold', got %+v", inlines[1])
	}

	it, ok := inlines[3].(EmphasisInline)
	if !ok || InlinesToPlainText(it.Children) != "italic" {
		t.Errorf("expected italic 'italic', got %+v", inlines[3])
	}

	cs, ok := inlines[5].(CodeSpanInline)
	if !ok || cs.Code != "code" {
		t.Errorf("expected code 'code', got %+v", inlines[5])
	}

	lnk, ok := inlines[7].(LinkInline)
	if !ok || lnk.Text != "link" || lnk.URL != "https://golang.org" {
		t.Errorf("expected link, got %+v", inlines[7])
	}

	st, ok := inlines[9].(StrikethroughInline)
	if !ok || InlinesToPlainText(st.Children) != "strike" {
		t.Errorf("expected strikethrough 'strike', got %+v", inlines[9])
	}
}

func TestCalculateStats(t *testing.T) {
	input := `# Title

Paragraph with ten words in this sentence to test counting nicely.

- Item 1
- Item 2

` + "```go\ncode\n```\n"

	doc := Parse(input)
	stats := CalculateStats(doc, input)

	if stats.HeadingCount != 1 {
		t.Errorf("expected 1 heading, got %d", stats.HeadingCount)
	}
	if stats.CodeBlockCount != 1 {
		t.Errorf("expected 1 code block, got %d", stats.CodeBlockCount)
	}
	if stats.ListCount != 1 || stats.ListItemCount != 2 {
		t.Errorf("expected 1 list with 2 items, got %d lists, %d items", stats.ListCount, stats.ListItemCount)
	}
	if stats.WordCount == 0 {
		t.Errorf("expected word count > 0, got %d", stats.WordCount)
	}
}
