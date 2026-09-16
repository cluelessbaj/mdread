---
title: Markdown Reader & Parser Demo
author: Antigravity
version: 1.0.0
category: Tooling
---

# Welcome to mdread

`mdread` is a fast, lightweight, and beautiful Markdown parser and reader built natively in **Go**.

---

## Features Overview

This tool allows you to comfortably read and inspect markdown files from your terminal or in your web browser.

### Key Highlights
- **Zero Dependencies**: Powered exclusively by the Go standard library.
- **Rich Terminal Styling**: Colors, underlines, unicode box drawings, and syntax highlighting.
- **Smart Pager**: Automatically uses `less -R` for long documents.
- **Table of Contents**: Inspect document outlines at a glance.
- **Live Browser Preview**: Automatic reloading while you write.

---

## Typography & Formatting

You can write **bold text**, *italic text*, ***bold and italic***, or ~~strikethrough~~.
You can also use inline `code spans` like `fmt.Println("Hello!")`, or link to external references like [Go Language](https://golang.org).

### Blockquotes

> "Simplicity is prerequisite for reliability."
> — Edsger W. Dijkstra

Nested quotes are also supported:
> First level of quotation.
> > Second nested quotation level.

---

## Code Blocks

Fenced code blocks are neatly encased in stylized terminal boxes with language identifiers and line numbers.

```go
package main

import "fmt"

func main() {
    message := "Hello from mdread!"
    fmt.Println(message)
}
```

```bash
# Build and install mdread
go build -o bin/mdread ./cmd/mdread
./bin/mdread README.md
```

---

## Lists & Tasks

### Task Progress
- [x] Create AST and line-by-line block parser
- [x] Implement ANSI terminal renderer
- [x] Add auto-pager integration (`less -R`)
- [x] Implement Outline / Table of Contents mode (`--toc`)
- [x] Add Document Statistics metrics (`--stats`)
- [ ] Add PDF export backend

### Nested Numbered List
1. Getting Started
   1. Install Go 1.22+
   2. Clone or navigate to the repository
2. Running the Tool
   - Use `mdread <file.md>` for terminal reading
   - Use `mdread --serve 8080 <file.md>` for web preview

---

## Formatted Tables

Pipe tables with column alignment (left, center, right):

| Command | Option | Description | Status |
| :--- | :---: | :--- | ---: |
| `mdread file.md` | Default | Terminal reader with auto pager | Stable |
| `mdread --toc` | `-t` | Heading outline tree | Active |
| `mdread --stats` | `-s` | Word count & reading time | Active |
| `mdread --serve` | `[port]` | Browser preview with auto reload | Active |
| `mdread --html` | None | Export standalone HTML | Ready |

---

*Thank you for using mdread!*
