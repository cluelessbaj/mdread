package markdown

// TOCEntry represents a single item in the table of contents.
type TOCEntry struct {
	Level   int
	Title   string
	ID      string
	LineNum int
}

// ExtractTOC traverses the document and returns all headings in order.
func ExtractTOC(doc *Document) []TOCEntry {
	var entries []TOCEntry
	var collectHeadings func(nodes []Node)
	collectHeadings = func(nodes []Node) {
		for _, node := range nodes {
			switch n := node.(type) {
			case *Heading:
				entries = append(entries, TOCEntry{
					Level:   n.Level,
					Title:   InlinesToPlainText(n.Inlines),
					ID:      n.ID,
					LineNum: n.LineNum,
				})
			case *Blockquote:
				collectHeadings(n.Children)
			case *List:
				for _, item := range n.Items {
					collectHeadings(item.Children)
				}
			}
		}
	}
	collectHeadings(doc.Children)
	return entries
}
