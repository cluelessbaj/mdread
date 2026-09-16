package term

import (
	"fmt"
	"strings"

	"mdreader/pkg/markdown"
)

// TerminalOptions configures terminal rendering behavior.
type TerminalOptions struct {
	Width    int
	UseColor bool
}

// RenderTerminal converts a Document AST into an ANSI-formatted terminal string.
func RenderTerminal(doc *markdown.Document, opts TerminalOptions) string {
	if opts.Width <= 0 {
		opts.Width, _ = GetTerminalSize()
	}
	if opts.Width > 100 {
		opts.Width = 100
	}
	if opts.Width < 40 {
		opts.Width = 40
	}

	var sb strings.Builder

	// 1. Render FrontMatter if present
	if len(doc.FrontMatter) > 0 {
		renderFrontMatter(&sb, doc.FrontMatter, opts)
	}

	// 2. Render Children
	for i, child := range doc.Children {
		renderBlock(&sb, child, opts, 0)
		if i < len(doc.Children)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func renderFrontMatter(sb *strings.Builder, fm map[string]string, opts TerminalOptions) {
	title := " METADATA "
	boxWidth := opts.Width - 4
	if boxWidth < 30 {
		boxWidth = 30
	}

	topBar := "┌─" + Style(title, FgBrightCyan, Bold) + strings.Repeat("─", boxWidth-VisibleLen(title)-2) + "┐"
	sb.WriteString(Style(topBar, FgBrightBlack) + "\n")

	for k, v := range fm {
		line := fmt.Sprintf("│ %s: %s", Style(k, FgCyan, Bold), v)
		padding := boxWidth - VisibleLen(line) + 1
		if padding < 1 {
			padding = 1
		}
		sb.WriteString(line + strings.Repeat(" ", padding) + Style("│", FgBrightBlack) + "\n")
	}

	bottomBar := "└" + strings.Repeat("─", boxWidth) + "┘"
	sb.WriteString(Style(bottomBar, FgBrightBlack) + "\n\n")
}

func renderBlock(sb *strings.Builder, node markdown.Node, opts TerminalOptions, depth int) {
	indent := strings.Repeat("  ", depth)

	switch n := node.(type) {
	case *markdown.Heading:
		renderHeading(sb, n, opts, indent)

	case *markdown.Paragraph:
		renderParagraph(sb, n, opts, indent)

	case *markdown.Blockquote:
		renderBlockquote(sb, n, opts, depth)

	case *markdown.CodeBlock:
		renderCodeBlock(sb, n, opts, indent)

	case *markdown.List:
		renderList(sb, n, opts, depth)

	case *markdown.Table:
		renderTable(sb, n, opts, indent)

	case *markdown.ThematicBreak:
		lineWidth := opts.Width - VisibleLen(indent)
		if lineWidth < 10 {
			lineWidth = 10
		}
		sb.WriteString(indent + Style(strings.Repeat("─", lineWidth), FgBrightBlack) + "\n")
	}
}

func renderHeading(sb *strings.Builder, h *markdown.Heading, opts TerminalOptions, indent string) {
	text := renderInlines(h.Inlines)
	lineWidth := opts.Width - VisibleLen(indent)

	switch h.Level {
	case 1:
		banner := "█ " + text
		sb.WriteString("\n" + indent + Style(banner, FgBrightCyan, Bold) + "\n")
		sb.WriteString(indent + Style(strings.Repeat("═", lineWidth), FgBrightCyan) + "\n")
	case 2:
		banner := "■ " + text
		sb.WriteString("\n" + indent + Style(banner, FgBrightBlue, Bold) + "\n")
		sb.WriteString(indent + Style(strings.Repeat("─", lineWidth), FgBrightBlue) + "\n")
	case 3:
		sb.WriteString(indent + Style("◆ "+text, FgBrightYellow, Bold) + "\n")
	case 4:
		sb.WriteString(indent + Style("● "+text, FgBrightMagenta, Bold) + "\n")
	case 5:
		sb.WriteString(indent + Style("○ "+text, FgBrightWhite, Bold) + "\n")
	default:
		sb.WriteString(indent + Style("▸ "+text, FgBrightBlack, Italic) + "\n")
	}
}

func renderParagraph(sb *strings.Builder, p *markdown.Paragraph, opts TerminalOptions, indent string) {
	text := renderInlines(p.Inlines)
	maxTextWidth := opts.Width - VisibleLen(indent) - 2
	if maxTextWidth < 20 {
		maxTextWidth = 20
	}

	lines := WrapText(text, maxTextWidth)
	for _, line := range lines {
		sb.WriteString(indent + line + "\n")
	}
}

func renderBlockquote(sb *strings.Builder, b *markdown.Blockquote, opts TerminalOptions, depth int) {
	var innerSb strings.Builder
	for _, child := range b.Children {
		renderBlock(&innerSb, child, opts, 0)
	}

	bar := Style("│ ", FgBrightBlue, Bold)
	lines := strings.Split(strings.TrimRight(innerSb.String(), "\n"), "\n")
	for _, l := range lines {
		sb.WriteString(strings.Repeat("  ", depth) + bar + Style(l, Italic) + "\n")
	}
}

func renderCodeBlock(sb *strings.Builder, cb *markdown.CodeBlock, opts TerminalOptions, indent string) {
	langBadge := ""
	if cb.Language != "" {
		langBadge = " [" + cb.Language + "] "
	}

	boxWidth := opts.Width - VisibleLen(indent) - 2
	if boxWidth < 30 {
		boxWidth = 30
	}

	badgeLen := VisibleLen(langBadge)
	dashes := boxWidth - badgeLen - 3
	if dashes < 2 {
		dashes = 2
	}

	topBar := "┌─" + Style(langBadge, FgBrightYellow, Bold) + strings.Repeat("─", dashes) + "┐"
	sb.WriteString(indent + Style(topBar, FgBrightBlack) + "\n")

	highlightedLines := HighlightCode(cb.Content, cb.Language)
	for _, line := range highlightedLines {
		sb.WriteString(indent + Style("│ ", FgBrightBlack) + line + "\n")
	}

	bottomBar := "└" + strings.Repeat("─", boxWidth) + "┘"
	sb.WriteString(indent + Style(bottomBar, FgBrightBlack) + "\n")
}

func renderList(sb *strings.Builder, l *markdown.List, opts TerminalOptions, depth int) {
	indent := strings.Repeat("  ", depth)

	for i, item := range l.Items {
		var prefix string
		if item.IsTask {
			if item.Checked {
				prefix = Style("[✓] ", FgBrightGreen, Bold)
			} else {
				prefix = Style("[ ] ", FgBrightBlack)
			}
		} else if l.Ordered {
			prefix = Style(fmt.Sprintf("%d. ", l.Start+i), FgBrightCyan)
		} else {
			bullets := []string{"• ", "○ ", "▪ "}
			b := bullets[depth%len(bullets)]
			prefix = Style(b, FgBrightYellow)
		}

		itemText := renderInlines(item.Inlines)
		prefixLen := VisibleLen(prefix)
		maxTextWidth := opts.Width - VisibleLen(indent) - prefixLen - 2
		if maxTextWidth < 20 {
			maxTextWidth = 20
		}

		wrapped := WrapText(itemText, maxTextWidth)
		if len(wrapped) > 0 {
			sb.WriteString(indent + prefix + wrapped[0] + "\n")
			subIndent := indent + strings.Repeat(" ", prefixLen)
			for _, wl := range wrapped[1:] {
				sb.WriteString(subIndent + wl + "\n")
			}
		} else {
			sb.WriteString(indent + prefix + "\n")
		}

		// Render child blocks (sublists, paragraphs, etc.)
		for _, child := range item.Children {
			renderBlock(sb, child, opts, depth+1)
		}
	}
}

func renderTable(sb *strings.Builder, t *markdown.Table, opts TerminalOptions, indent string) {
	colCount := len(t.Headers)
	if colCount == 0 {
		return
	}

	colWidths := make([]int, colCount)
	for i, h := range t.Headers {
		w := VisibleLen(renderInlines(h.Inlines))
		if w > colWidths[i] {
			colWidths[i] = w
		}
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < colCount {
				w := VisibleLen(renderInlines(cell.Inlines))
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	// Add 2 padding spaces to each column
	for i := range colWidths {
		colWidths[i] += 2
		if colWidths[i] < 4 {
			colWidths[i] = 4
		}
	}

	// Draw top border: ┌──────┬──────┐
	var top strings.Builder
	top.WriteString(indent + Style("┌", FgBrightBlack))
	for i, w := range colWidths {
		top.WriteString(Style(strings.Repeat("─", w), FgBrightBlack))
		if i < colCount-1 {
			top.WriteString(Style("┬", FgBrightBlack))
		}
	}
	top.WriteString(Style("┐", FgBrightBlack))
	sb.WriteString(top.String() + "\n")

	// Header row: │ Header 1 │ Header 2 │
	var hdr strings.Builder
	hdr.WriteString(indent + Style("│", FgBrightBlack))
	for i, h := range t.Headers {
		cellStr := renderInlines(h.Inlines)
		padded := formatCell(cellStr, colWidths[i], t.Alignments[i])
		hdr.WriteString(Style(padded, Bold, FgBrightWhite))
		hdr.WriteString(Style("│", FgBrightBlack))
	}
	sb.WriteString(hdr.String() + "\n")

	// Separator border: ├──────┼──────┤
	var sep strings.Builder
	sep.WriteString(indent + Style("├", FgBrightBlack))
	for i, w := range colWidths {
		sep.WriteString(Style(strings.Repeat("─", w), FgBrightBlack))
		if i < colCount-1 {
			sep.WriteString(Style("┼", FgBrightBlack))
		}
	}
	sep.WriteString(Style("┤", FgBrightBlack))
	sb.WriteString(sep.String() + "\n")

	// Data rows
	for _, row := range t.Rows {
		var rowSb strings.Builder
		rowSb.WriteString(indent + Style("│", FgBrightBlack))
		for i := 0; i < colCount; i++ {
			cellStr := ""
			if i < len(row) {
				cellStr = renderInlines(row[i].Inlines)
			}
			align := markdown.AlignNone
			if i < len(t.Alignments) {
				align = t.Alignments[i]
			}
			padded := formatCell(cellStr, colWidths[i], align)
			rowSb.WriteString(padded)
			rowSb.WriteString(Style("│", FgBrightBlack))
		}
		sb.WriteString(rowSb.String() + "\n")
	}

	// Bottom border: └──────┴──────┘
	var bot strings.Builder
	bot.WriteString(indent + Style("└", FgBrightBlack))
	for i, w := range colWidths {
		bot.WriteString(Style(strings.Repeat("─", w), FgBrightBlack))
		if i < colCount-1 {
			bot.WriteString(Style("┴", FgBrightBlack))
		}
	}
	bot.WriteString(Style("┘", FgBrightBlack))
	sb.WriteString(bot.String() + "\n")
}

func formatCell(content string, width int, align markdown.Alignment) string {
	content = " " + content + " "
	switch align {
	case markdown.AlignRight:
		return padLeft(content, width)
	case markdown.AlignCenter:
		return PadCenter(content, width)
	default:
		return PadRight(content, width)
	}
}

func padLeft(s string, width int) string {
	curr := VisibleLen(s)
	if curr >= width {
		return s
	}
	return strings.Repeat(" ", width-curr) + s
}

func renderInlines(inlines []markdown.InlineNode) string {
	var sb strings.Builder
	for _, in := range inlines {
		switch n := in.(type) {
		case markdown.TextInline:
			sb.WriteString(n.Content)
		case markdown.CodeSpanInline:
			sb.WriteString(Style(" "+n.Code+" ", BgDarkGray, FgBrightYellow))
		case markdown.EmphasisInline:
			sb.WriteString(Style(renderInlines(n.Children), Italic))
		case markdown.StrongInline:
			sb.WriteString(Style(renderInlines(n.Children), Bold))
		case markdown.StrongEmphasisInline:
			sb.WriteString(Style(renderInlines(n.Children), Bold, Italic))
		case markdown.StrikethroughInline:
			sb.WriteString(Style(renderInlines(n.Children), Strike, Dim))
		case markdown.LinkInline:
			linkText := Style(n.Text, FgBrightCyan, Underline)
			if n.URL != "" && n.URL != n.Text {
				sb.WriteString(linkText + " " + Style("("+n.URL+")", FgBrightBlack))
			} else {
				sb.WriteString(linkText)
			}
		case markdown.ImageInline:
			sb.WriteString(Style("[🖼 "+n.Alt+"]", FgBrightMagenta) + " " + Style("("+n.URL+")", FgBrightBlack))
		}
	}
	return sb.String()
}
