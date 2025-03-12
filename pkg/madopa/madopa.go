package madopa

import (
	"github.com/shonnnoronha/madopa/internal/parser"
	"github.com/shonnnoronha/madopa/internal/renderer"
)

func Convert(markdown string, renderer renderer.Renderer) (string, error) {
	doc, err := parser.Parse(markdown)
	if err != nil {
		return "", err
	}

	html, err := renderer.Render(doc)
	if err != nil {
		return "", err
	}

	return html, nil
}

func ConvertWithOptions(markdown string, opts *renderer.Options) (string, error) {
	doc, err := parser.Parse(markdown)
	if err != nil {
		return "", err
	}

	htmlRenderer := renderer.NewHTMLRenderer(opts)

	html, err := htmlRenderer.Render(doc)
	if err != nil {
		return "", err
	}

	return html, nil
}
