package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"mdreader/pkg/html"
	"mdreader/pkg/markdown"
	"mdreader/pkg/server"
	"mdreader/pkg/term"
)

const version = "1.0.0"

func printUsage() {
	fmt.Printf(`%s - Terminal Markdown Reader & Parser (v%s)

Usage:
  mdread [flags] <file.md>
  cat <file.md> | mdread [flags]

Flags:
  -t, --toc          Display the Table of Contents / Outline
  -s, --stats        Display document statistics (word count, reading time, elements)
      --html         Output standalone HTML document
      --serve [port] Start a local live-reload preview server (default port 8080)
      --no-pager     Disable the interactive terminal pager (less -R)
  -w, --width <num>  Set output column width (defaults to terminal width, max 100)
  -v, --version      Print version information
  -h, --help         Show this help message

Examples:
  mdread README.md                 # Read markdown in terminal with auto-pager
  mdread --toc README.md           # Inspect document outline
  mdread --stats README.md         # View word count and reading time
  mdread --serve 8080 README.md    # View in browser with live reload
  mdread --html doc.md > doc.html  # Export to standalone styled HTML
  curl -s https://... | mdread     # Read markdown directly from stdin
`, term.Style("mdread", term.FgBrightCyan, term.Bold), version)
}

func main() {
	var (
		showHelp    bool
		showVersion bool
		showTOC     bool
		showStats   bool
		exportHTML  bool
		noPager     bool
		servePort   int
		width       int
	)

	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version (shorthand)")
	flag.BoolVar(&showTOC, "toc", false, "Show Table of Contents")
	flag.BoolVar(&showTOC, "t", false, "Show Table of Contents (shorthand)")
	flag.BoolVar(&showStats, "stats", false, "Show document statistics")
	flag.BoolVar(&showStats, "s", false, "Show document statistics (shorthand)")
	flag.BoolVar(&exportHTML, "html", false, "Output standalone HTML")
	flag.BoolVar(&noPager, "no-pager", false, "Disable pager")
	flag.IntVar(&servePort, "serve", 0, "Start local HTTP preview server on specified port (e.g. 8080)")
	flag.IntVar(&width, "width", 0, "Set output width")
	flag.IntVar(&width, "w", 0, "Set output width (shorthand)")

	flag.Usage = printUsage
	flag.Parse()

	if showHelp {
		printUsage()
		return
	}

	if showVersion {
		fmt.Printf("mdread version %s\n", version)
		return
	}

	args := flag.Args()
	var input string
	var filename string

	// Check if reading from stdin or file
	if len(args) == 0 || args[0] == "-" {
		// Read from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 && len(args) == 0 {
			printUsage()
			os.Exit(1)
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
		filename = "stdin"
	} else {
		filename = args[0]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
			os.Exit(1)
		}
		input = string(data)
	}

	// Web preview server mode
	if servePort > 0 {
		if filename == "stdin" {
			fmt.Fprintln(os.Stderr, "Error: --serve requires a real file on disk (not stdin).")
			os.Exit(1)
		}
		if err := server.StartServer(filename, servePort); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Parse markdown document into AST
	doc := markdown.Parse(input)

	// TOC mode
	if showTOC {
		renderTOC(doc)
		return
	}

	// Stats mode
	if showStats {
		renderStats(doc, input)
		return
	}

	// HTML Export mode
	if exportHTML {
		htmlDoc := html.RenderHTML(doc, filename)
		fmt.Println(htmlDoc)
		return
	}

	// Terminal Reader mode (Default)
	rendered := term.RenderTerminal(doc, term.TerminalOptions{
		Width:    width,
		UseColor: true,
	})

	if err := term.DisplayWithPager(rendered, noPager); err != nil {
		fmt.Fprintf(os.Stderr, "Display error: %v\n", err)
	}
}

func renderTOC(doc *markdown.Document) {
	toc := markdown.ExtractTOC(doc)
	if len(toc) == 0 {
		fmt.Println(term.Style("No headings found in document.", term.Dim))
		return
	}

	boxWidth := 60
	topBar := "┌─ " + term.Style("TABLE OF CONTENTS", term.FgBrightCyan, term.Bold) + " " + strings.Repeat("─", boxWidth-22) + "┐"
	fmt.Println(term.Style(topBar, term.FgBrightBlack))

	for _, entry := range toc {
		indent := strings.Repeat("  ", entry.Level-1)
		bullet := "•"
		switch entry.Level {
		case 1:
			bullet = term.Style("■", term.FgBrightCyan, term.Bold)
		case 2:
			bullet = term.Style("□", term.FgBrightBlue)
		case 3:
			bullet = term.Style("◆", term.FgBrightYellow)
		default:
			bullet = term.Style("▸", term.FgBrightBlack)
		}

		lineInfo := ""
		if entry.LineNum > 0 {
			lineInfo = term.Style(fmt.Sprintf(" (L%d)", entry.LineNum), term.FgBrightBlack)
		}

		title := entry.Title
		fmt.Printf("│ %s%s %s%s\n", indent, bullet, title, lineInfo)
	}

	bottomBar := "└" + strings.Repeat("─", boxWidth) + "┘"
	fmt.Println(term.Style(bottomBar, term.FgBrightBlack))
}

func renderStats(doc *markdown.Document, rawInput string) {
	stats := markdown.CalculateStats(doc, rawInput)

	boxWidth := 46
	topBar := "┌─ " + term.Style("DOCUMENT STATISTICS", term.FgBrightCyan, term.Bold) + " " + strings.Repeat("─", boxWidth-24) + "┐"
	fmt.Println(term.Style(topBar, term.FgBrightBlack))

	printStatLine("Words", fmt.Sprintf("%d", stats.WordCount))
	printStatLine("Characters", fmt.Sprintf("%d", stats.CharCount))
	printStatLine("Lines", fmt.Sprintf("%d", stats.LineCount))

	readTimeStr := fmt.Sprintf("%d sec", int(stats.ReadingTime.Seconds()))
	if stats.ReadingTime >= 60 {
		mins := int(stats.ReadingTime.Minutes())
		secs := int(stats.ReadingTime.Seconds()) % 60
		readTimeStr = fmt.Sprintf("%d min %d sec", mins, secs)
	}
	printStatLine("Est. Reading Time", readTimeStr)
	printStatLine("Headings", fmt.Sprintf("%d", stats.HeadingCount))
	printStatLine("Code Blocks", fmt.Sprintf("%d", stats.CodeBlockCount))
	printStatLine("Tables", fmt.Sprintf("%d", stats.TableCount))
	printStatLine("Lists", fmt.Sprintf("%d (%d items)", stats.ListCount, stats.ListItemCount))

	bottomBar := "└" + strings.Repeat("─", boxWidth) + "┘"
	fmt.Println(term.Style(bottomBar, term.FgBrightBlack))
}

func printStatLine(label, value string) {
	styledLabel := term.Style(fmt.Sprintf("%-18s", label+":"), term.FgCyan)
	styledVal := term.Style(fmt.Sprintf("%18s", value), term.Bold, term.FgBrightWhite)
	fmt.Printf("│  %s %s  │\n", styledLabel, styledVal)
}
