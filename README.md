# mdread 📖

A fast, lightweight, zero-dependency Markdown parser and terminal reader written in Go.

## Features

- **Rich Terminal Reader**: Styled headers, syntax-highlighted code blocks, blockquotes, task checkboxes `[✓]`, and ASCII/Unicode tables.
- **Smart Pager**: Automatically detects terminal height and uses `less -R` when reading long documents.
- **Table of Contents (`--toc`)**: Outlines heading hierarchy with line numbers.
- **Document Statistics (`--stats`)**: Computes word count, character count, estimated reading time, and structural element counts.
- **HTML Export (`--html`)**: Exports clean standalone HTML with responsive GitHub-like styling.
- **Live Browser Preview (`--serve [port]`)**: Starts a local HTTP server with automatic reloading when the markdown file is edited on disk.
- **Stdin Support**: Read directly from pipes (e.g. `cat file.md | mdread` or `curl -s URL | mdread`).
- **Zero External Dependencies**: Built 100% using standard Go.

## Installation

```bash
cd ~/Projects/mdreader
make install
```
This builds and copies `mdread` into `~/go/bin/mdread` (which is already in your `$PATH`).

## Usage Examples

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

# Preview in web browser with live reload on file edit
mdread --serve 8080 README.md

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
