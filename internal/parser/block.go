package parser

type Block interface {
	isBlock()
}

type Heading struct {
	Level int
	Text  []Inline
}

type Paragraph struct {
	Text []Inline
}

type List struct {
	Items []*ListItem
	Type  ListType
}

type ListItem struct {
	Level    int
	Content  []Inline
	Children *List
}

type CodeBlock struct {
	Lang string
	Code string
}

type Table struct {
	Headers    []TableCell
	Rows       [][]TableCell
	Alignments []Alignment
}

type TableCell struct {
	Content []Inline
}

type Blockquote struct {
	Items []*BlockquoteItem
}

type BlockquoteItem struct {
	Level    int
	Content  []Inline
	Children *Blockquote
}

type Alignment int

const (
	AlignDefault Alignment = iota
	AlignLeft
	AlignCenter
	AlignRight
)

type ListType int

const (
	UnorderedList ListType = iota
	OrderedList
)

func (h Heading) isBlock()    {}
func (p Paragraph) isBlock()  {}
func (l List) isBlock()       {}
func (li ListItem) isBlock()  {}
func (c CodeBlock) isBlock()  {}
func (t Table) isBlock()      {}
func (b Blockquote) isBlock() {}
