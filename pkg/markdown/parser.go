package markdown

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Parse parses a markdown string into a Document AST.
func Parse(input string) *Document {
	doc := &Document{
		FrontMatter: make(map[string]string),
	}

	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	idx := 0

	// 1. Check for YAML frontmatter at start
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		idx = 1
		for idx < len(lines) {
			line := lines[idx]
			if strings.TrimSpace(line) == "---" {
				idx++
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Strip quotes if present
				val = strings.Trim(val, `"'`)
				if key != "" {
					doc.FrontMatter[key] = val
				}
			}
			idx++
		}
	}

	// 2. Parse blocks for the rest of document
	doc.Children = parseBlocks(lines, idx, 0)
	return doc
}

// parseBlocks parses a slice of lines starting from startIdx until the end.
func parseBlocks(lines []string, startIdx int, baseLineNum int) []Node {
	var nodes []Node
	i := startIdx

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Blank lines
		if trimmed == "" {
			i++
			continue
		}

		// Fenced code block
		if isCodeFenceStart(trimmed) {
			fence := trimmed[:3]
			lang := strings.TrimSpace(trimmed[3:])
			var codeLines []string
			i++
			for i < len(lines) {
				currTrimmed := strings.TrimSpace(lines[i])
				if strings.HasPrefix(currTrimmed, fence) {
					i++
					break
				}
				codeLines = append(codeLines, lines[i])
				i++
			}
			nodes = append(nodes, &CodeBlock{
				Language: lang,
				Content:  strings.Join(codeLines, "\n"),
			})
			continue
		}

		// Thematic break (HR): 3 or more -, *, _
		if isThematicBreak(trimmed) {
			nodes = append(nodes, &ThematicBreak{})
			i++
			continue
		}

		// ATX Heading: # to ######
		if level, headingText, ok := parseHeading(line); ok {
			nodes = append(nodes, &Heading{
				Level:   level,
				RawText: headingText,
				Inlines: ParseInlines(headingText),
				ID:      slugify(headingText),
				LineNum: baseLineNum + i + 1,
			})
			i++
			continue
		}

		// Table: must have | and the next line must be table delimiter
		if isTableCandidate(line) && i+1 < len(lines) && isTableDelimiter(lines[i+1]) {
			table, nextIdx := parseTable(lines, i)
			nodes = append(nodes, table)
			i = nextIdx
			continue
		}

		// Blockquote: starts with >
		if strings.HasPrefix(trimmed, ">") {
			var bqLines []string
			for i < len(lines) {
				currTrimmed := strings.TrimSpace(lines[i])
				if !strings.HasPrefix(currTrimmed, ">") {
					// Blank line or break terminates blockquote
					break
				}
				// Remove leading '>' and optional single space
				rest := strings.TrimPrefix(currTrimmed, ">")
				if strings.HasPrefix(rest, " ") {
					rest = rest[1:]
				}
				bqLines = append(bqLines, rest)
				i++
			}
			bqNodes := parseBlocks(bqLines, 0, baseLineNum+i)
			nodes = append(nodes, &Blockquote{
				Children: bqNodes,
			})
			continue
		}

		// List (unordered or ordered)
		if isListItem(line) {
			list, nextIdx := parseList(lines, i)
			nodes = append(nodes, list)
			i = nextIdx
			continue
		}

		// Paragraph: collect lines until blank line or block start
		var pLines []string
		for i < len(lines) {
			curr := lines[i]
			currTrimmed := strings.TrimSpace(curr)
			if currTrimmed == "" {
				i++
				break
			}
			// If next line starts a heading, code fence, hr, or list, stop paragraph
			if isCodeFenceStart(currTrimmed) || isThematicBreak(currTrimmed) || isListItem(curr) || strings.HasPrefix(currTrimmed, ">") {
				break
			}
			if _, _, ok := parseHeading(curr); ok {
				break
			}
			if isTableCandidate(curr) && i+1 < len(lines) && isTableDelimiter(lines[i+1]) {
				break
			}

			pLines = append(pLines, currTrimmed)
			i++
		}

		rawText := strings.Join(pLines, " ")
		if strings.TrimSpace(rawText) != "" {
			nodes = append(nodes, &Paragraph{
				RawText: rawText,
				Inlines: ParseInlines(rawText),
			})
		}
	}

	return nodes
}

// Helpers for block parsing

func isCodeFenceStart(s string) bool {
	return strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~")
}

func isThematicBreak(s string) bool {
	clean := strings.ReplaceAll(s, " ", "")
	if len(clean) < 3 {
		return false
	}
	ch := clean[0]
	if ch != '-' && ch != '*' && ch != '_' {
		return false
	}
	for i := 1; i < len(clean); i++ {
		if clean[i] != ch {
			return false
		}
	}
	return true
}

func parseHeading(line string) (int, string, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimmed, "#") {
		return 0, "", false
	}
	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level > 6 {
		return 0, "", false
	}
	if level < len(trimmed) && trimmed[level] != ' ' && trimmed[level] != '\t' {
		return 0, "", false
	}
	text := strings.TrimSpace(trimmed[level:])
	// Strip optional trailing hashes: "### Title ###" -> "Title"
	text = strings.TrimRight(text, "#")
	text = strings.TrimSpace(text)
	return level, text, true
}

func isTableCandidate(line string) bool {
	return strings.Contains(line, "|")
}

func isTableDelimiter(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "|") {
		return false
	}
	trimmed = strings.Trim(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		cell := strings.TrimSpace(p)
		if len(cell) == 0 {
			return false
		}
		// Cell must consist of '-', optional ':' at ends
		dashCount := 0
		for idx, r := range cell {
			if r == ':' {
				if idx != 0 && idx != len(cell)-1 {
					return false
				}
			} else if r == '-' {
				dashCount++
			} else {
				return false
			}
		}
		if dashCount == 0 {
			return false
		}
	}
	return true
}

func parseTable(lines []string, startIdx int) (*Table, int) {
	headerLine := lines[startIdx]
	delimLine := lines[startIdx+1]

	headers := splitTableCells(headerLine)
	delims := splitTableCells(delimLine)

	alignments := make([]Alignment, len(headers))
	for i, d := range delims {
		if i >= len(alignments) {
			break
		}
		cell := strings.TrimSpace(d)
		leftCol := strings.HasPrefix(cell, ":")
		rightCol := strings.HasSuffix(cell, ":")
		if leftCol && rightCol {
			alignments[i] = AlignCenter
		} else if rightCol {
			alignments[i] = AlignRight
		} else if leftCol {
			alignments[i] = AlignLeft
		} else {
			alignments[i] = AlignNone
		}
	}

	var headerCells []TableCell
	for _, h := range headers {
		txt := strings.TrimSpace(h)
		headerCells = append(headerCells, TableCell{
			RawText: txt,
			Inlines: ParseInlines(txt),
		})
	}

	var rows [][]TableCell
	i := startIdx + 2
	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == "" || !strings.Contains(line, "|") {
			break
		}
		rawCells := splitTableCells(line)
		var row []TableCell
		for idx, c := range rawCells {
			if idx >= len(headerCells) {
				break
			}
			txt := strings.TrimSpace(c)
			row = append(row, TableCell{
				RawText: txt,
				Inlines: ParseInlines(txt),
			})
		}
		// Pad with empty cells if row has fewer cells than headers
		for len(row) < len(headerCells) {
			row = append(row, TableCell{RawText: "", Inlines: nil})
		}
		rows = append(rows, row)
		i++
	}

	return &Table{
		Headers:    headerCells,
		Alignments: alignments,
		Rows:       rows,
	}, i
}

func splitTableCells(line string) []string {
	trimmed := strings.TrimSpace(line)
	// Remove outer pipes if present
	if strings.HasPrefix(trimmed, "|") {
		trimmed = trimmed[1:]
	}
	if strings.HasSuffix(trimmed, "|") {
		trimmed = trimmed[:len(trimmed)-1]
	}

	var cells []string
	var current strings.Builder
	escaped := false

	for _, r := range trimmed {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '|' {
			cells = append(cells, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}
	cells = append(cells, current.String())
	return cells
}

var (
	unorderedListRegex = regexp.MustCompile(`^(\s*)([\*\-\+])\s+(.*)$`)
	orderedListRegex   = regexp.MustCompile(`^(\s*)(\d+)[\.\)]\s+(.*)$`)
	taskListRegex      = regexp.MustCompile(`^\[([ xX])\]\s*(.*)$`)
)

func isListItem(line string) bool {
	return unorderedListRegex.MatchString(line) || orderedListRegex.MatchString(line)
}

func parseList(lines []string, startIdx int) (*List, int) {
	firstLine := lines[startIdx]
	ordered := orderedListRegex.MatchString(firstLine)
	startNum := 1
	if ordered {
		m := orderedListRegex.FindStringSubmatch(firstLine)
		if len(m) >= 3 {
			startNum, _ = strconv.Atoi(m[2])
		}
	}

	list := &List{
		Ordered: ordered,
		Start:   startNum,
	}

	i := startIdx
	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			// Look ahead: if next non-empty line is an indented line or list item, continue; else break
			j := i + 1
			for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
				j++
			}
			if j < len(lines) && (isListItem(lines[j]) || strings.HasPrefix(lines[j], "  ") || strings.HasPrefix(lines[j], "\t")) {
				i = j
				line = lines[i]
			} else {
				break
			}
		}

		var itemText string
		var itemIndent string
		matched := false

		if ordered {
			if m := orderedListRegex.FindStringSubmatch(line); len(m) >= 4 {
				itemIndent = m[1]
				itemText = m[3]
				matched = true
			}
		} else {
			if m := unorderedListRegex.FindStringSubmatch(line); len(m) >= 4 {
				itemIndent = m[1]
				itemText = m[3]
				matched = true
			}
		}

		if !matched {
			// Not a matching list item; could be end of list
			break
		}

		item := &ListItem{}

		// Check for task list checkbox: [ ] or [x]
		if tm := taskListRegex.FindStringSubmatch(itemText); len(tm) >= 3 {
			item.IsTask = true
			item.Checked = strings.ToLower(tm[1]) == "x"
			itemText = tm[2]
		}

		item.Inlines = ParseInlines(itemText)

		// Collect indented continuation lines or sublists
		i++
		var subLines []string
		for i < len(lines) {
			subLine := lines[i]
			subTrimmed := strings.TrimSpace(subLine)
			if subTrimmed == "" {
				// Empty line inside list item
				if i+1 < len(lines) && (strings.HasPrefix(lines[i+1], itemIndent+"  ") || strings.HasPrefix(lines[i+1], "\t")) {
					subLines = append(subLines, "")
					i++
					continue
				}
				break
			}

			// If it's a new top-level item with same or less indent, break
			if isListItem(subLine) {
				currIndent := getIndent(subLine)
				if len(currIndent) <= len(itemIndent) {
					break
				}
			} else if len(getIndent(subLine)) <= len(itemIndent) {
				break
			}

			// Strip list item's indentation prefix for nested block parsing
			stripped := subLine
			if strings.HasPrefix(stripped, itemIndent+"    ") {
				stripped = stripped[len(itemIndent)+4:]
			} else if strings.HasPrefix(stripped, itemIndent+"  ") {
				stripped = stripped[len(itemIndent)+2:]
			} else if strings.HasPrefix(stripped, "\t") {
				stripped = stripped[1:]
			}
			subLines = append(subLines, stripped)
			i++
		}

		if len(subLines) > 0 {
			item.Children = parseBlocks(subLines, 0, 0)
		}

		list.Items = append(list.Items, item)
	}

	return list, i
}

func getIndent(s string) string {
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return s[:i]
		}
	}
	return s
}

func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
				b.WriteByte('-')
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// -------------------------------------------------------------
// Inline Parsing
// -------------------------------------------------------------

// ParseInlines parses inline markdown formatting within a string.
func ParseInlines(text string) []InlineNode {
	var nodes []InlineNode
	var buf strings.Builder

	flushBuf := func() {
		if buf.Len() > 0 {
			nodes = append(nodes, TextInline{Content: buf.String()})
			buf.Reset()
		}
	}

	i := 0
	n := len(text)

	for i < n {
		// 1. Inline code: `code`
		if text[i] == '`' {
			closeIdx := strings.IndexByte(text[i+1:], '`')
			if closeIdx != -1 {
				flushBuf()
				code := text[i+1 : i+1+closeIdx]
				nodes = append(nodes, CodeSpanInline{Code: code})
				i = i + 1 + closeIdx + 1
				continue
			}
		}

		// 2. Image: ![alt](url)
		if i+1 < n && text[i] == '!' && text[i+1] == '[' {
			if img, endIdx, ok := parseImageInline(text, i); ok {
				flushBuf()
				nodes = append(nodes, img)
				i = endIdx
				continue
			}
		}

		// 3. Link: [text](url)
		if text[i] == '[' {
			if link, endIdx, ok := parseLinkInline(text, i); ok {
				flushBuf()
				nodes = append(nodes, link)
				i = endIdx
				continue
			}
		}

		// 4. Autolink: <http://...> or <https://...>
		if text[i] == '<' {
			if autoLink, endIdx, ok := parseAutoLink(text, i); ok {
				flushBuf()
				nodes = append(nodes, autoLink)
				i = endIdx
				continue
			}
		}

		// 5. Strikethrough: ~~text~~
		if i+1 < n && text[i] == '~' && text[i+1] == '~' {
			closeIdx := strings.Index(text[i+2:], "~~")
			if closeIdx != -1 {
				flushBuf()
				inner := text[i+2 : i+2+closeIdx]
				nodes = append(nodes, StrikethroughInline{Children: ParseInlines(inner)})
				i = i + 2 + closeIdx + 2
				continue
			}
		}

		// 6. Strong + Emphasis: ***text*** or ___text___
		if i+2 < n && (strings.HasPrefix(text[i:], "***") || strings.HasPrefix(text[i:], "___")) {
			delim := text[i : i+3]
			closeIdx := strings.Index(text[i+3:], delim)
			if closeIdx != -1 {
				flushBuf()
				inner := text[i+3 : i+3+closeIdx]
				nodes = append(nodes, StrongEmphasisInline{Children: ParseInlines(inner)})
				i = i + 3 + closeIdx + 3
				continue
			}
		}

		// 7. Strong: **text** or __text__
		if i+1 < n && (strings.HasPrefix(text[i:], "**") || strings.HasPrefix(text[i:], "__")) {
			delim := text[i : i+2]
			closeIdx := strings.Index(text[i+2:], delim)
			if closeIdx != -1 {
				flushBuf()
				inner := text[i+2 : i+2+closeIdx]
				nodes = append(nodes, StrongInline{Children: ParseInlines(inner)})
				i = i + 2 + closeIdx + 2
				continue
			}
		}

		// 8. Emphasis: *text* or _text_
		if text[i] == '*' || text[i] == '_' {
			delim := text[i]
			// Ensure it's not intra-word for underscore (e.g. some_var_name)
			if delim == '_' && i > 0 && isWordChar(rune(text[i-1])) {
				buf.WriteByte(text[i])
				i++
				continue
			}
			closeIdx := findMatchingDelimiter(text, i+1, delim)
			if closeIdx != -1 {
				flushBuf()
				inner := text[i+1 : closeIdx]
				nodes = append(nodes, EmphasisInline{Children: ParseInlines(inner)})
				i = closeIdx + 1
				continue
			}
		}

		// Escaped character: \* or \_ etc.
		if text[i] == '\\' && i+1 < n && isMarkdownSpecialChar(text[i+1]) {
			buf.WriteByte(text[i+1])
			i += 2
			continue
		}

		buf.WriteByte(text[i])
		i++
	}

	flushBuf()
	return nodes
}

func isMarkdownSpecialChar(b byte) bool {
	return b == '*' || b == '_' || b == '`' || b == '[' || b == ']' || b == '(' || b == ')' || b == '#' || b == '+' || b == '-' || b == '!' || b == '~' || b == '\\'
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func findMatchingDelimiter(text string, start int, delim byte) int {
	for i := start; i < len(text); i++ {
		if text[i] == '\\' {
			i++
			continue
		}
		if text[i] == delim {
			if delim == '_' && i+1 < len(text) && isWordChar(rune(text[i+1])) {
				continue
			}
			return i
		}
	}
	return -1
}

func parseLinkInline(text string, start int) (LinkInline, int, bool) {
	// [text](url "title")
	bracketClose := findBalancedClose(text, start, '[', ']')
	if bracketClose == -1 || bracketClose+1 >= len(text) || text[bracketClose+1] != '(' {
		return LinkInline{}, 0, false
	}
	parenClose := findBalancedClose(text, bracketClose+1, '(', ')')
	if parenClose == -1 {
		return LinkInline{}, 0, false
	}

	linkText := text[start+1 : bracketClose]
	linkDest := strings.TrimSpace(text[bracketClose+2 : parenClose])

	url := linkDest
	title := ""
	if quoteIdx := strings.IndexAny(linkDest, `"'`); quoteIdx != -1 {
		url = strings.TrimSpace(linkDest[:quoteIdx])
		title = strings.Trim(linkDest[quoteIdx:], `"' `)
	}

	return LinkInline{
		Text:  linkText,
		URL:   url,
		Title: title,
	}, parenClose + 1, true
}

func parseImageInline(text string, start int) (ImageInline, int, bool) {
	// ![alt](url "title")
	link, endIdx, ok := parseLinkInline(text[1:], start)
	if !ok {
		return ImageInline{}, 0, false
	}
	return ImageInline{
		Alt:   link.Text,
		URL:   link.URL,
		Title: link.Title,
	}, endIdx + 1, true
}

func parseAutoLink(text string, start int) (LinkInline, int, bool) {
	closeIdx := strings.IndexByte(text[start+1:], '>')
	if closeIdx == -1 {
		return LinkInline{}, 0, false
	}
	inner := text[start+1 : start+1+closeIdx]
	if strings.HasPrefix(inner, "http://") || strings.HasPrefix(inner, "https://") || strings.Contains(inner, "@") {
		url := inner
		if strings.Contains(inner, "@") && !strings.HasPrefix(inner, "mailto:") {
			url = "mailto:" + inner
		}
		return LinkInline{
			Text: inner,
			URL:  url,
		}, start + 1 + closeIdx + 1, true
	}
	return LinkInline{}, 0, false
}

func findBalancedClose(text string, start int, openCh, closeCh byte) int {
	depth := 0
	for i := start; i < len(text); i++ {
		if text[i] == '\\' {
			i++
			continue
		}
		if text[i] == openCh {
			depth++
		} else if text[i] == closeCh {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// InlinesToPlainText extracts plain text from a list of inline nodes.
func InlinesToPlainText(inlines []InlineNode) string {
	var b strings.Builder
	for _, in := range inlines {
		switch n := in.(type) {
		case TextInline:
			b.WriteString(n.Content)
		case CodeSpanInline:
			b.WriteString(n.Code)
		case EmphasisInline:
			b.WriteString(InlinesToPlainText(n.Children))
		case StrongInline:
			b.WriteString(InlinesToPlainText(n.Children))
		case StrongEmphasisInline:
			b.WriteString(InlinesToPlainText(n.Children))
		case StrikethroughInline:
			b.WriteString(InlinesToPlainText(n.Children))
		case LinkInline:
			b.WriteString(n.Text)
		case ImageInline:
			b.WriteString(n.Alt)
		}
	}
	return b.String()
}
