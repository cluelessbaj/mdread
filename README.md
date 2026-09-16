# mdread 📖

A fast, lightweight, zero-dependency Markdown parser, terminal reader, and PDF exporter written in Go.

## Features

- **Rich Terminal Reader**: Styled headers, syntax-highlighted code blocks, blockquotes, task checkboxes `[✓]`, and ASCII/Unicode tables.
- **Smart Pager**: Automatically detects terminal height and uses `less -R` when reading long documents.
- **PDF Export (`--pdf`)**: Convert and copy markdown document previews directly into beautifully styled PDF files.
- **Live Browser Preview (`--serve [port]`)**: Starts a local HTTP server with automatic file reloading, a "Save as PDF" print button, and a one-click "Download PDF" button.
- **Table of Contents (`--toc`)**: Outlines heading hierarchy with line numbers.
- **Document Statistics (`--stats`)**: Computes word count, character count, estimated reading time, and structural element counts.
- **HTML Export (`--html`)**: Exports clean standalone HTML with responsive GitHub-like styling and `@media print` formatting.
- **Stdin Support**: Read directly from pipes (e.g. `cat file.md | mdread` or `curl -s URL | mdread`).
- **Zero External Dependencies**: Built using standard Go.

## Installation

```bash
cd ~/Projects/mdreader
make install
```
This builds and installs `mdread` into `~/go/bin/mdread` (which is already in your `$PATH`).

## Usage Examples

### 1. PDF Export
```bash
# Export markdown directly to a PDF
mdread --pdf README.md
# Output: README.pdf

# Export to a custom PDF file name
mdread --pdf report.pdf document.md
```

### 2. Live Web Preview with PDF Actions
```bash
# Start local preview server on port 8080
mdread --serve 8080 README.md
```
The browser view features a sticky toolbar with:
- **🖨️ Print / Save PDF**: Opens browser print dialog with page-break-optimized styling.
- **⬇️ Download PDF**: Generates and downloads the PDF directly from the server.
- **📋 Copy Text**: Copies formatted plain text of the document to your clipboard.

### 3. Terminal Reading
```bash
# Read a markdown file with automatic pager support
mdread README.md

# Read without pager
mdread --no-pager README.md

# View Table of Contents / Outline
mdread --toc README.md

# View document metrics (word count, reading time)
mdread --stats README.md

# Export to standalone styled HTML
mdread --html README.md > document.html

# Pipe input from stdin
curl -s https://raw.githubusercontent.com/golang/go/master/README.md | mdread
```

## Programmatic API

You can also import `mdreader/pkg/markdown` as a library in your own Go projects:

```go
package main

import (
    "fmt"
    "mdreader/pkg/markdown"
    "mdreader/pkg/term"
)

func main() {
    doc := markdown.Parse("# Hello **World**\nThis is a test.")
    output := term.RenderTerminal(doc, term.TerminalOptions{Width: 80, UseColor: true})
    fmt.Print(output)
}
```

## License

MIT License - Copyright (c) 2026 cluelessbaj
