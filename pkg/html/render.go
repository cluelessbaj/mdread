package html

import (
	"fmt"
	"html"
	"strings"

	"mdreader/pkg/markdown"
)

// RenderHTML converts a Document AST into a standalone, styled HTML document.
func RenderHTML(doc *markdown.Document, title string) string {
	if title == "" && doc.FrontMatter["title"] != "" {
		title = doc.FrontMatter["title"]
	}
	if title == "" {
		title = "Markdown Document"
	}

	var body strings.Builder

	// Render FrontMatter if present
	if len(doc.FrontMatter) > 0 {
		body.WriteString("<div class=\"metadata-box\">\n")
		body.WriteString("<h3>Metadata</h3>\n<dl>\n")
		for k, v := range doc.FrontMatter {
			body.WriteString(fmt.Sprintf("  <dt>%s</dt><dd>%s</dd>\n", html.EscapeString(k), html.EscapeString(v)))
		}
		body.WriteString("</dl>\n</div>\n")
	}

	for _, child := range doc.Children {
		renderHTMLBlock(&body, child)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s</title>
  <style>
    :root {
      --bg: #0d1117;
      --fg: #c9d1d9;
      --card-bg: #161b22;
      --border: #30363d;
      --accent: #58a6ff;
      --code-bg: #1f242c;
      --blockquote: #8b949e;
      --table-alt: #161b22;
    }
    @media (prefers-color-scheme: light) {
      :root {
        --bg: #ffffff;
        --fg: #24292f;
        --card-bg: #f6f8fa;
        --border: #d0d7de;
        --accent: #0969da;
        --code-bg: #f6f8fa;
        --blockquote: #57606a;
        --table-alt: #f6f8fa;
      }
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
      line-height: 1.6;
      color: var(--fg);
      background-color: var(--bg);
      margin: 0;
      padding: 2rem 1rem;
    }
    .container {
      max-width: 860px;
      margin: 0 auto;
    }
    h1, h2, h3, h4, h5, h6 {
      color: var(--fg);
      margin-top: 1.5em;
      margin-bottom: 0.5em;
      font-weight: 600;
    }
    h1 { border-bottom: 1px solid var(--border); padding-bottom: 0.3em; color: var(--accent); }
    h2 { border-bottom: 1px solid var(--border); padding-bottom: 0.2em; }
    p { margin: 1em 0; }
    a { color: var(--accent); text-decoration: none; }
    a:hover { text-decoration: underline; }
    code {
      font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
      font-size: 85%%;
      background-color: var(--code-bg);
      padding: 0.2em 0.4em;
      border-radius: 4px;
      border: 1px solid var(--border);
    }
    pre {
      background-color: var(--code-bg);
      border: 1px solid var(--border);
      border-radius: 6px;
      padding: 16px;
      overflow-x: auto;
    }
    pre code {
      background: none;
      padding: 0;
      border: none;
    }
    blockquote {
      margin: 1em 0;
      padding: 0 1em;
      color: var(--blockquote);
      border-left: 0.25em solid var(--border);
    }
    table {
      border-collapse: collapse;
      width: 100%%;
      margin: 1.5em 0;
    }
    th, td {
      border: 1px solid var(--border);
      padding: 8px 12px;
    }
    tr:nth-child(even) {
      background-color: var(--table-alt);
    }
    th {
      background-color: var(--card-bg);
      font-weight: 600;
    }
    hr {
      border: 0;
      height: 1px;
      background: var(--border);
      margin: 2em 0;
    }
    .metadata-box {
      background: var(--card-bg);
      border: 1px solid var(--border);
      border-radius: 6px;
      padding: 1rem;
      margin-bottom: 2rem;
    }
    .metadata-box h3 { margin-top: 0; }
    dl { display: grid; grid-template-columns: 120px 1fr; row-gap: 0.3em; margin: 0; }
    dt { font-weight: bold; color: var(--accent); }
    .task-item { list-style: none; margin-left: -1.2em; }

    /* Print & PDF Stylesheet */
    @page {
      size: auto;
      margin: 1.6cm 1.4cm;
    }
    @media print {
      body {
        background: #ffffff !important;
        color: #111827 !important;
        font-size: 11pt !important;
        line-height: 1.5 !important;
        padding: 0 !important;
      }
      .container {
        max-width: 100%% !important;
        margin: 0 !important;
        padding: 0 !important;
      }
      .no-print {
        display: none !important;
      }
      h1, h2, h3, h4, h5, h6 {
        color: #111827 !important;
        page-break-after: avoid !important;
        break-after: avoid !important;
      }
      h1 { color: #0969da !important; }
      pre, blockquote, table, .metadata-box {
        page-break-inside: avoid !important;
        break-inside: avoid !important;
      }
      pre {
        background-color: #f8fafc !important;
        border: 1px solid #e2e8f0 !important;
        color: #0f172a !important;
      }
      code {
        background-color: #f1f5f9 !important;
        border: 1px solid #cbd5e1 !important;
        color: #0f172a !important;
      }
      table, th, td {
        border-color: #cbd5e1 !important;
      }
      th {
        background-color: #f1f5f9 !important;
        color: #0f172a !important;
      }
      tr:nth-child(even) {
        background-color: #f8fafc !important;
      }
      a {
        color: #0969da !important;
        text-decoration: underline !important;
      }
      .metadata-box {
        background-color: #f8fafc !important;
        border: 1px solid #cbd5e1 !important;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    %s
  </div>
</body>
</html>`, html.EscapeString(title), body.String())
}

func renderHTMLBlock(sb *strings.Builder, node markdown.Node) {
	switch n := node.(type) {
	case *markdown.Heading:
		sb.WriteString(fmt.Sprintf("<h%d id=\"%s\">%s</h%d>\n", n.Level, n.ID, renderHTMLInlines(n.Inlines), n.Level))
	case *markdown.Paragraph:
		sb.WriteString(fmt.Sprintf("<p>%s</p>\n", renderHTMLInlines(n.Inlines)))
	case *markdown.Blockquote:
		sb.WriteString("<blockquote>\n")
		for _, child := range n.Children {
			renderHTMLBlock(sb, child)
		}
		sb.WriteString("</blockquote>\n")
	case *markdown.CodeBlock:
		langClass := ""
		if n.Language != "" {
			langClass = fmt.Sprintf(" class=\"language-%s\"", html.EscapeString(n.Language))
		}
		sb.WriteString(fmt.Sprintf("<pre><code%s>%s</code></pre>\n", langClass, html.EscapeString(n.Content)))
	case *markdown.List:
		tag := "ul"
		if n.Ordered {
			tag = "ol"
		}
		sb.WriteString(fmt.Sprintf("<%s>\n", tag))
		for _, item := range n.Items {
			classAttr := ""
			prefix := ""
			if item.IsTask {
				classAttr = " class=\"task-item\""
				if item.Checked {
					prefix = "<input type=\"checkbox\" checked disabled> "
				} else {
					prefix = "<input type=\"checkbox\" disabled> "
				}
			}
			sb.WriteString(fmt.Sprintf("  <li%s>%s%s", classAttr, prefix, renderHTMLInlines(item.Inlines)))
			for _, child := range item.Children {
				renderHTMLBlock(sb, child)
			}
			sb.WriteString("</li>\n")
		}
		sb.WriteString(fmt.Sprintf("</%s>\n", tag))
	case *markdown.Table:
		sb.WriteString("<table>\n<thead>\n  <tr>\n")
		for i, h := range n.Headers {
			align := ""
			if i < len(n.Alignments) && n.Alignments[i] != markdown.AlignNone {
				align = fmt.Sprintf(" align=\"%s\"", n.Alignments[i])
			}
			sb.WriteString(fmt.Sprintf("    <th%s>%s</th>\n", align, renderHTMLInlines(h.Inlines)))
		}
		sb.WriteString("  </tr>\n</thead>\n<tbody>\n")
		for _, row := range n.Rows {
			sb.WriteString("  <tr>\n")
			for i, cell := range row {
				align := ""
				if i < len(n.Alignments) && n.Alignments[i] != markdown.AlignNone {
					align = fmt.Sprintf(" align=\"%s\"", n.Alignments[i])
				}
				sb.WriteString(fmt.Sprintf("    <td%s>%s</td>\n", align, renderHTMLInlines(cell.Inlines)))
			}
			sb.WriteString("  </tr>\n")
		}
		sb.WriteString("</tbody>\n</table>\n")
	case *markdown.ThematicBreak:
		sb.WriteString("<hr>\n")
	}
}

func renderHTMLInlines(inlines []markdown.InlineNode) string {
	var sb strings.Builder
	for _, in := range inlines {
		switch n := in.(type) {
		case markdown.TextInline:
			sb.WriteString(html.EscapeString(n.Content))
		case markdown.CodeSpanInline:
			sb.WriteString(fmt.Sprintf("<code>%s</code>", html.EscapeString(n.Code)))
		case markdown.EmphasisInline:
			sb.WriteString(fmt.Sprintf("<em>%s</em>", renderHTMLInlines(n.Children)))
		case markdown.StrongInline:
			sb.WriteString(fmt.Sprintf("<strong>%s</strong>", renderHTMLInlines(n.Children)))
		case markdown.StrongEmphasisInline:
			sb.WriteString(fmt.Sprintf("<strong><em>%s</em></strong>", renderHTMLInlines(n.Children)))
		case markdown.StrikethroughInline:
			sb.WriteString(fmt.Sprintf("<del>%s</del>", renderHTMLInlines(n.Children)))
		case markdown.LinkInline:
			titleAttr := ""
			if n.Title != "" {
				titleAttr = fmt.Sprintf(" title=\"%s\"", html.EscapeString(n.Title))
			}
			sb.WriteString(fmt.Sprintf("<a href=\"%s\"%s>%s</a>", html.EscapeString(n.URL), titleAttr, html.EscapeString(n.Text)))
		case markdown.ImageInline:
			titleAttr := ""
			if n.Title != "" {
				titleAttr = fmt.Sprintf(" title=\"%s\"", html.EscapeString(n.Title))
			}
			sb.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"%s\"%s>", html.EscapeString(n.URL), html.EscapeString(n.Alt), titleAttr))
		}
	}
	return sb.String()
}
