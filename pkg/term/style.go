package term

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ANSI escape sequences
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Underline   = "\033[4m"
	Blink       = "\033[5m"
	Inverse     = "\033[7m"
	Hidden      = "\033[8m"
	Strike      = "\033[9m"

	FgBlack     = "\033[30m"
	FgRed       = "\033[31m"
	FgGreen     = "\033[32m"
	FgYellow    = "\033[33m"
	FgBlue      = "\033[34m"
	FgMagenta   = "\033[35m"
	FgCyan      = "\033[36m"
	FgWhite     = "\033[37m"

	FgBrightBlack   = "\033[90m"
	FgBrightRed     = "\033[91m"
	FgBrightGreen   = "\033[92m"
	FgBrightYellow  = "\033[93m"
	FgBrightBlue    = "\033[94m"
	FgBrightMagenta = "\033[95m"
	FgBrightCyan    = "\033[96m"
	FgBrightWhite   = "\033[97m"

	BgBlack     = "\033[40m"
	BgDarkGray  = "\033[100m"
	BgBlue      = "\033[44m"
	BgCyan      = "\033[46m"
)

// Style wraps text in ANSI escape codes and resets at the end.
func Style(s string, codes ...string) string {
	if len(codes) == 0 || s == "" {
		return s
	}
	return strings.Join(codes, "") + s + Reset
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// StripANSI removes all ANSI escape sequences to compute printable character length.
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// VisibleLen returns the visual rune width of a string disregarding ANSI escape codes.
func VisibleLen(s string) int {
	clean := StripANSI(s)
	return utf8.RuneCountInString(clean)
}

// GetTerminalSize returns columns and rows of the current terminal.
func GetTerminalSize() (width int, height int) {
	width, height = 80, 24
	out, err := exec.Command("stty", "size").Output()
	if err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			if h, err := strconv.Atoi(parts[0]); err == nil && h > 0 {
				height = h
			}
			if w, err := strconv.Atoi(parts[1]); err == nil && w > 0 {
				width = w
			}
		}
	}
	// Fallback to COLUMNS and LINES env vars if set
	if envW := os.Getenv("COLUMNS"); envW != "" {
		if w, err := strconv.Atoi(envW); err == nil && w > 0 {
			width = w
		}
	}
	if envH := os.Getenv("LINES"); envH != "" {
		if h, err := strconv.Atoi(envH); err == nil && h > 0 {
			height = h
		}
	}
	return width, height
}

// PadRight pads a string with spaces up to target visual width.
func PadRight(s string, width int) string {
	curr := VisibleLen(s)
	if curr >= width {
		return s
	}
	return s + strings.Repeat(" ", width-curr)
}

// PadCenter centers a string with spaces up to target visual width.
func PadCenter(s string, width int) string {
	curr := VisibleLen(s)
	if curr >= width {
		return s
	}
	left := (width - curr) / 2
	right := width - curr - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// WrapText wraps text according to printable width, breaking at word boundaries.
func WrapText(text string, width int) []string {
	if width <= 10 {
		width = 80
	}
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	var current strings.Builder
	currLen := 0

	for _, w := range words {
		wLen := VisibleLen(w)
		if currLen == 0 {
			current.WriteString(w)
			currLen = wLen
		} else if currLen+1+wLen <= width {
			current.WriteByte(' ')
			current.WriteString(w)
			currLen += 1 + wLen
		} else {
			lines = append(lines, current.String())
			current.Reset()
			current.WriteString(w)
			currLen = wLen
		}
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return lines
}

// HighlightCode applies simple syntax coloring to code lines.
func HighlightCode(code string, lang string) []string {
	lines := strings.Split(code, "\n")
	var result []string

	keywords := map[string]bool{
		"func": true, "return": true, "package": true, "import": true, "var": true,
		"const": true, "type": true, "struct": true, "interface": true, "if": true,
		"else": true, "for": true, "range": true, "switch": true, "case": true,
		"default": true, "go": true, "select": true, "defer": true, "nil": true,
		"true": true, "false": true, "def": true, "class": true, "class_name": true,
		"function": true, "let": true, "const_val": true, "export": true,
	}

	types := map[string]bool{
		"string": true, "int": true, "int64": true, "float64": true, "bool": true,
		"byte": true, "rune": true, "error": true, "any": true,
	}

	for lineIdx, line := range lines {
		trimmed := strings.TrimSpace(line)
		lineNumPrefix := Style(fmt.Sprintf("%3d │ ", lineIdx+1), FgBrightBlack)

		// Comment highlighting
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			result = append(result, lineNumPrefix+Style(line, FgBrightBlack, Italic))
			continue
		}

		// Simple word-based highlighting
		words := strings.Split(line, " ")
		var styledWords []string
		for _, w := range words {
			cleanWord := strings.Trim(w, `(),.{}:;[]*&`)
			if keywords[cleanWord] {
				styledWords = append(styledWords, strings.Replace(w, cleanWord, Style(cleanWord, FgBrightMagenta, Bold), 1))
			} else if types[cleanWord] {
				styledWords = append(styledWords, strings.Replace(w, cleanWord, Style(cleanWord, FgBrightYellow), 1))
			} else if strings.HasPrefix(w, `"`) || strings.HasPrefix(w, `'`) || strings.HasPrefix(w, "`") {
				styledWords = append(styledWords, Style(w, FgBrightGreen))
			} else {
				styledWords = append(styledWords, w)
			}
		}
		result = append(result, lineNumPrefix+strings.Join(styledWords, " "))
	}
	return result
}
