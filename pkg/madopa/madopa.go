package madopa

import (
	"github.com/shonnnoronha/madopa/internal/parser"
	"github.com/shonnnoronha/madopa/internal/renderer"
)

type Renderer struct {
	options renderer.Options
}

func (r *Renderer) SetEscapeHTML(escapeHTML bool) {
	r.options.EscapeHTML = escapeHTML
}

func (r *Renderer) SetHardLineBreak(hardLineBreak bool) {
	r.options.HardLineBreak = hardLineBreak
}

func (r *Renderer) SetSoftLineBreak(softLineBreak bool) {
	r.options.SoftLineBreak = softLineBreak
}

func (r *Renderer) SetIncludeCss(includeCss bool) {
	r.options.IncludeCSS = includeCss
}

func (r *Renderer) SetCssFilePath(cssFilePath string) {
	r.options.CssFilePath = cssFilePath
}

func (r *Renderer) SetSyntaxHighlight(highlight bool) {
	r.options.IncludeSyntaxHighlight = highlight
}

func (r *Renderer) NewHTMLRenderer() renderer.Renderer {
	return renderer.NewHTMLRenderer(&r.options)
}

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
