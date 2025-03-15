package renderer

import (
	"github.com/shonnnoronha/madopa/internal/parser"
)

type Options struct {
	EscapeHTML             bool
	HardLineBreak          bool
	SoftLineBreak          bool
	IncludeCSS             bool
	CssFilePath            string
	IncludeSyntaxHighlight bool
}

type Renderer interface {
	Render(doc *parser.Document) (string, error)
	renderBlock(block parser.Block) error
	renderInlines(inlines []parser.Inline) error
}
