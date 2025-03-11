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
	Items []ListItem
}

type ListItem struct {
	Level int
	Text  []Inline
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

type Alignment int

const (
	AlignDefault Alignment = iota
	AlignLeft
	AlignCenter
	AlignRight
)

func (h Heading) isBlock()   {}
func (p Paragraph) isBlock() {}
func (l List) isBlock()      {}
func (li ListItem) isBlock() {}
func (c CodeBlock) isBlock() {}
func (t Table) isBlock()     {}
