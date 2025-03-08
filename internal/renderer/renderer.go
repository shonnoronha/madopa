package renderer

import (
	"github.com/shonnnoronha/madopa/internal/parser"
)

type Renderer interface {
	Render(doc *parser.Document) (string, error)
}

