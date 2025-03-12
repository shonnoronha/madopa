package renderer

import (
	"bytes"
	"fmt"
	"html"

	"github.com/shonnnoronha/madopa/internal/parser"
)

type HTMLRenderer struct {
	buffer *bytes.Buffer
	opts   *Options
}

type Options struct {
	EscapeHTML    bool
	HardLineBreak bool
	SoftLineBreak bool
	AutoLink      bool
	Strikethrough bool
	Table         bool
	TaskList      bool
}

func NewHTMLRenderer(opts *Options) *HTMLRenderer {
	if opts == nil {
		opts = &Options{
			EscapeHTML: true,
		}
	}
	return &HTMLRenderer{
		buffer: &bytes.Buffer{},
		opts:   opts,
	}
}

func (r *HTMLRenderer) Render(doc *parser.Document) (string, error) {
	r.buffer.Reset()

	for _, block := range doc.Blocks {
		if err := r.renderBlock(block); err != nil {
			return "", err
		}
	}

	return r.buffer.String(), nil
}

func (r *HTMLRenderer) renderBlock(block parser.Block) error {
	switch b := block.(type) {
	case *parser.Heading:
		level := b.Level
		r.buffer.WriteString(fmt.Sprintf("<h%d>", level))
		if err := r.renderInlines(b.Text); err != nil {
			return err
		}
		r.buffer.WriteString(fmt.Sprintf("</h%d>\n", level))

	case *parser.Paragraph:
		r.buffer.WriteString("<p>")
		if err := r.renderInlines(b.Text); err != nil {
			return err
		}
		r.buffer.WriteString("</p>\n")

	case *parser.CodeBlock:
		if b.Lang != "" {
			r.buffer.WriteString(fmt.Sprintf("<pre><code class=\"%s\">", b.Lang))
		} else {
			r.buffer.WriteString("<pre><code>")
		}
		content := b.Code
		if r.opts.EscapeHTML {
			content = html.EscapeString(content)
		}
		r.buffer.WriteString(content)
		r.buffer.WriteString("</code></pre>\n")

	case *parser.Table:
		r.buffer.WriteString("<table>\n")

		r.buffer.WriteString("<thead>\n")
		r.buffer.WriteString("<tr>\n")

		alignDir := ""

		for i, cell := range b.Headers {
			if i < len(b.Alignments) {
				switch b.Alignments[i] {
				case parser.AlignLeft:
					alignDir = " align=\"left\""
				case parser.AlignCenter:
					alignDir = " align=\"center\""
				case parser.AlignRight:
					alignDir = " align=\"right\""
				}
			}
			r.buffer.WriteString(fmt.Sprintf("<th%s>", alignDir))
			if err := r.renderInlines(cell.Content); err != nil {
				return err
			}
			r.buffer.WriteString("</th>\n")
		}
		r.buffer.WriteString("</tr>\n")
		r.buffer.WriteString("</thead>\n")

		if len(b.Rows) > 0 {
			r.buffer.WriteString("<tbody>\n")
			for _, row := range b.Rows {
				r.buffer.WriteString("<tr>\n")
				for i, cell := range row {
					if i < len(b.Alignments) {
						switch b.Alignments[i] {
						case parser.AlignLeft:
							alignDir = " align=\"left\""
						case parser.AlignCenter:
							alignDir = " align=\"center\""
						case parser.AlignRight:
							alignDir = " align=\"right\""
						}
					}
					r.buffer.WriteString(fmt.Sprintf("<td%s>", alignDir))
					if err := r.renderInlines(cell.Content); err != nil {
						return err
					}
					r.buffer.WriteString("</td>\n")
				}

				r.buffer.WriteString("</tr>\n")
			}
			r.buffer.WriteString("</tbody>\n")
		}

		r.buffer.WriteString("</table>\n")

	default:
		r.buffer.WriteString(fmt.Sprintf("<!-- Unsupported block type: %T -->\n", b))
	}

	return nil
}

func (r *HTMLRenderer) renderInlines(inlines []parser.Inline) error {
	for _, inline := range inlines {
		switch i := inline.(type) {
		case *parser.Text:
			content := i.Content
			if r.opts.EscapeHTML {
				content = html.EscapeString(content)
			}
			r.buffer.WriteString(content)

		case *parser.BoldItalic:
			r.buffer.WriteString("<strong><em>")
			if err := r.renderInlines(i.Content); err != nil {
				return err
			}
			r.buffer.WriteString("</em></strong>")

		case *parser.Bold:
			r.buffer.WriteString("<strong>")
			if err := r.renderInlines(i.Content); err != nil {
				return err
			}
			r.buffer.WriteString("</strong>")

		case *parser.Italic:
			r.buffer.WriteString("<em>")
			if err := r.renderInlines(i.Content); err != nil {
				return err
			}
			r.buffer.WriteString("</em>")

		case *parser.Link:
			r.buffer.WriteString("<a href=\"")
			r.buffer.WriteString(i.URL)
			r.buffer.WriteString("\">")
			if err := r.renderInlines(i.Text); err != nil {
				return err
			}
			r.buffer.WriteString("</a>")

		default:
			r.buffer.WriteString(fmt.Sprintf("<!-- Unsupported inline type: %T -->", i))
		}
	}
	return nil
}
