package markdown

// NodeType identifies the kind of block node.
type NodeType string

const (
	NodeDocument      NodeType = "Document"
	NodeFrontMatter   NodeType = "FrontMatter"
	NodeHeading       NodeType = "Heading"
	NodeParagraph     NodeType = "Paragraph"
	NodeBlockquote    NodeType = "Blockquote"
	NodeCodeBlock     NodeType = "CodeBlock"
	NodeList          NodeType = "List"
	NodeListItem      NodeType = "ListItem"
	NodeTable         NodeType = "Table"
	NodeThematicBreak NodeType = "ThematicBreak"
)

// Node represents any markdown block node.
type Node interface {
	NodeType() NodeType
}

// InlineType identifies the kind of inline node.
type InlineType string

const (
	InlineText           InlineType = "Text"
	InlineEmphasis       InlineType = "Emphasis"
	InlineStrong         InlineType = "Strong"
	InlineStrongEmphasis InlineType = "StrongEmphasis"
	InlineStrikethrough  InlineType = "Strikethrough"
	InlineCodeSpan       InlineType = "CodeSpan"
	InlineLink           InlineType = "Link"
	InlineImage          InlineType = "Image"
)

// InlineNode represents an inline element within text.
type InlineNode interface {
	InlineType() InlineType
}

// TextInline represents plain text.
type TextInline struct {
	Content string
}

func (t TextInline) InlineType() InlineType { return InlineText }

// EmphasisInline represents italic text (*text* or _text_).
type EmphasisInline struct {
	Children []InlineNode
}

func (e EmphasisInline) InlineType() InlineType { return InlineEmphasis }

// StrongInline represents bold text (**text** or __text__).
type StrongInline struct {
	Children []InlineNode
}

func (s StrongInline) InlineType() InlineType { return InlineStrong }

// StrongEmphasisInline represents bold and italic text (***text*** or ___text___).
type StrongEmphasisInline struct {
	Children []InlineNode
}

func (se StrongEmphasisInline) InlineType() InlineType { return InlineStrongEmphasis }

// StrikethroughInline represents deleted text (~~text~~).
type StrikethroughInline struct {
	Children []InlineNode
}

func (s StrikethroughInline) InlineType() InlineType { return InlineStrikethrough }

// CodeSpanInline represents inline code (`code`).
type CodeSpanInline struct {
	Code string
}

func (c CodeSpanInline) InlineType() InlineType { return InlineCodeSpan }

// LinkInline represents a link ([text](url "title")).
type LinkInline struct {
	Text  string
	URL   string
	Title string
}

func (l LinkInline) InlineType() InlineType { return InlineLink }

// ImageInline represents an embedded image (![alt](url "title")).
type ImageInline struct {
	Alt   string
	URL   string
	Title string
}

func (i ImageInline) InlineType() InlineType { return InlineImage }

// Document is the root node of the parsed AST.
type Document struct {
	FrontMatter map[string]string
	Children    []Node
}

func (d *Document) NodeType() NodeType { return NodeDocument }

// Heading represents an ATX heading (# to ######).
type Heading struct {
	Level   int
	RawText string
	Inlines []InlineNode
	ID      string
	LineNum int
}

func (h *Heading) NodeType() NodeType { return NodeHeading }

// Paragraph represents a block of regular text.
type Paragraph struct {
	RawText string
	Inlines []InlineNode
}

func (p *Paragraph) NodeType() NodeType { return NodeParagraph }

// Blockquote represents a blockquote (> text).
type Blockquote struct {
	Children []Node
}

func (b *Blockquote) NodeType() NodeType { return NodeBlockquote }

// CodeBlock represents a fenced code block with optional syntax language.
type CodeBlock struct {
	Language string
	Content  string
}

func (c *CodeBlock) NodeType() NodeType { return NodeCodeBlock }

// List represents an ordered or unordered list.
type List struct {
	Ordered bool
	Start   int
	Items   []*ListItem
}

func (l *List) NodeType() NodeType { return NodeList }

// ListItem represents an item inside a List.
type ListItem struct {
	IsTask   bool
	Checked  bool
	Inlines  []InlineNode
	Children []Node // For nested lists or sub-blocks
}

func (li *ListItem) NodeType() NodeType { return NodeListItem }

// Alignment defines the column alignment in a table.
type Alignment string

const (
	AlignNone   Alignment = "none"
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

// TableCell represents a single cell in a table row.
type TableCell struct {
	RawText string
	Inlines []InlineNode
}

// Table represents a GFM pipe table.
type Table struct {
	Headers    []TableCell
	Alignments []Alignment
	Rows       [][]TableCell
}

func (t *Table) NodeType() NodeType { return NodeTable }

// ThematicBreak represents a horizontal rule (---, ***, ___).
type ThematicBreak struct{}

func (tb *ThematicBreak) NodeType() NodeType { return NodeThematicBreak }
