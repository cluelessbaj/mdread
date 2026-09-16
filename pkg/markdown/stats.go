package markdown

import (
	"strings"
	"time"
)

// DocumentStats holds metrics about a markdown document.
type DocumentStats struct {
	WordCount      int
	CharCount      int
	LineCount      int
	ReadingTime    time.Duration
	HeadingCount   int
	CodeBlockCount int
	TableCount     int
	ListCount      int
	ListItemCount  int
}

// CalculateStats computes textual and structural metrics for a document.
func CalculateStats(doc *Document, rawInput string) DocumentStats {
	stats := DocumentStats{}

	// Textual metrics
	stats.CharCount = len(rawInput)
	lines := strings.Split(rawInput, "\n")
	stats.LineCount = len(lines)

	words := strings.Fields(rawInput)
	stats.WordCount = len(words)

	// Average reading speed: ~200 words per minute
	minutes := float64(stats.WordCount) / 200.0
	stats.ReadingTime = time.Duration(minutes * float64(time.Minute))
	if stats.WordCount > 0 && stats.ReadingTime < time.Second {
		stats.ReadingTime = time.Second
	}

	// Structural metrics
	var walkNodes func(nodes []Node)
	walkNodes = func(nodes []Node) {
		for _, n := range nodes {
			switch node := n.(type) {
			case *Heading:
				stats.HeadingCount++
			case *CodeBlock:
				stats.CodeBlockCount++
			case *Table:
				stats.TableCount++
			case *List:
				stats.ListCount++
				stats.ListItemCount += len(node.Items)
				for _, item := range node.Items {
					walkNodes(item.Children)
				}
			case *Blockquote:
				walkNodes(node.Children)
			}
		}
	}
	walkNodes(doc.Children)

	return stats
}
