package renderer

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"os"

	"github.com/shonnnoronha/madopa/internal/parser"
)

const (
	defaultCssFilePath    = "./internal/renderer/styles/dark_blog.css"
	defaultScriptFilePath = "./internal/renderer/scripts/highlight.html"
)

type HTMLRenderer struct {
	buffer *bytes.Buffer
	opts   *Options
}

func NewHTMLRenderer(opts *Options) *HTMLRenderer {
	return &HTMLRenderer{
		buffer: &bytes.Buffer{},
		opts:   opts,
	}
}

func (r *HTMLRenderer) Render(doc *parser.Document) (string, error) {
	r.buffer.Reset()

	if r.opts.IncludeCSS {
		cssFilePath := r.opts.CssFilePath
		if cssFilePath == "" {
			cssFilePath = defaultCssFilePath
		}

		r.buffer.WriteString("<!DOCTYPE html>\n")
		r.buffer.WriteString("<html>\n<head>\n")
		r.buffer.WriteString("<meta charset=\"UTF-8\">\n")
		r.buffer.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
		r.buffer.WriteString("<title>Markdown Blog</title>\n")
		r.buffer.WriteString("<style>\n")

		cssContent, err := os.ReadFile(cssFilePath)
		if err != nil {
			log.Println("Error reading CSS file", err)
			return "", err
		}
		r.buffer.Write(cssContent)
		r.buffer.WriteString("\n</style>\n</head>\n<body>\n")

		scriptContent, err := os.ReadFile(defaultScriptFilePath)
		if err != nil {
			log.Println("Error reading Script file", err)
			return "", err
		}
		r.buffer.Write(scriptContent)

		r.buffer.WriteString("<div class=\"container\">\n")
		r.buffer.WriteString("<article class=\"post\">\n")
	}

	for _, block := range doc.Blocks {
		if err := r.renderBlock(block); err != nil {
			return "", err
		}
	}

	if r.opts.IncludeCSS {
		r.buffer.WriteString("\n</article>\n</div>\n</body>\n</html>")
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

	case *parser.List:
		if b.Type == parser.OrderedList {
			r.buffer.WriteString("<ol>\n")
		} else {
			r.buffer.WriteString("<ul>\n")
		}

		if err := r.renderListItems(b.Items); err != nil {
			return err
		}

		if b.Type == parser.OrderedList {
			r.buffer.WriteString("</ol>\n")
		} else {
			r.buffer.WriteString("</ul>\n")
		}

	case *parser.Blockquote:
		r.buffer.WriteString("<blockquote>")

		if err := r.renderBlockQuoteItems(b.Items); err != nil {
			return err
		}

		r.buffer.WriteString("</blockquote>")

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

		case *parser.CodeInline:
			r.buffer.WriteString("<code>")
			content := i.Content
			if r.opts.EscapeHTML {
				content = html.EscapeString(content)
			}
			r.buffer.WriteString(content)
			r.buffer.WriteString("</code>")

		case *parser.Image:
			r.buffer.WriteString("<img src=\"")
			if r.opts.EscapeHTML {
				r.buffer.WriteString(html.EscapeString(i.Src))
			} else {
				r.buffer.WriteString(i.Src)
			}
			r.buffer.WriteString("\" alt=\"")
			if r.opts.EscapeHTML {
				r.buffer.WriteString(html.EscapeString(i.Alt))
			} else {
				r.buffer.WriteString(i.Alt)
			}
			r.buffer.WriteString("\"")
			if i.Title != "" {
				r.buffer.WriteString("\" title=\"")
				if r.opts.EscapeHTML {
					r.buffer.WriteString(html.EscapeString(i.Title))
				} else {
					r.buffer.WriteString(i.Title)
				}
				r.buffer.WriteString("\"")
			}
			r.buffer.WriteString(">")

		default:
			r.buffer.WriteString(fmt.Sprintf("<!-- Unsupported inline type: %T -->", i))
		}
	}
	return nil
}

func (r *HTMLRenderer) renderListItems(items []*parser.ListItem) error {
	for _, item := range items {
		r.buffer.WriteString("<li>")
		if err := r.renderInlines(item.Content); err != nil {
			return err
		}
		if item.Children != nil {
			if item.Children.Type == parser.OrderedList {
				r.buffer.WriteString("\n<ol>\n")
			} else {
				r.buffer.WriteString("\n<ul>\n")
			}

			if err := r.renderListItems(item.Children.Items); err != nil {
				return err
			}

			if item.Children.Type == parser.OrderedList {
				r.buffer.WriteString("</ol>\n")
			} else {
				r.buffer.WriteString("</ul>\n")
			}
		}
		r.buffer.WriteString("</li>\n")
	}
	return nil
}

func (r *HTMLRenderer) renderBlockQuoteItems(items []*parser.BlockquoteItem) error {
	for _, item := range items {
		err := r.renderInlines(item.Content)
		r.buffer.WriteString("<br>")
		if err != nil {
			return err
		}

		if item.Children != nil {
			r.buffer.WriteString("<blockquote>")
			if err := r.renderBlockQuoteItems(item.Children.Items); err != nil {
				return err
			}
			r.buffer.WriteString("</blockquote>")
		}
	}
	return nil
}
