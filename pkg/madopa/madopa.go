package madopa

import (
	"github.com/shonnnoronha/madopa/internal/parser"
	"github.com/shonnnoronha/madopa/internal/renderer"
)

func ConvertToHTML(markdown string) (string, error) {
	doc, err := parser.Parse(markdown)
	if err != nil {
		return "", err
	}

	htmlRenderer := renderer.NewHTMLRenderer(nil)

	html, err := htmlRenderer.Render(doc)
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
