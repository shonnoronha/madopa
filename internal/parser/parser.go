package parser

import (
	"strings"
	"unicode"
)

type Document struct {
	Blocks []Block
}

type parser struct {
	input string
	pos   int
	line  string
}

func (p *parser) parse() (*Document, error) {
	doc := &Document{
		Blocks: make([]Block, 0),
	}

	for p.pos < len(p.input) {
		p.readLine()

		if strings.TrimSpace(p.line) == "" {
			continue
		}

		block, err := p.parseBlock()
		// fmt.Printf("Block type: %T\n", block)
		if err != nil {
			return nil, err
		}

		if block != nil {
			doc.Blocks = append(doc.Blocks, block)
		}
	}

	return doc, nil
}

func (p *parser) parseBlock() (Block, error) {
	trimmedLine := strings.TrimLeftFunc(p.line, unicode.IsSpace)

	if strings.HasPrefix(trimmedLine, "#") {
		return p.parseHeading()
	}

	if strings.HasPrefix(trimmedLine, "```") {
		return p.parseCodeBlock()
	}

	return p.parseParagraph()
}

func (p *parser) readLine() {
	end := strings.IndexByte(p.input[p.pos:], '\n')
	if end == -1 {
		end = len(p.input) - p.pos
	}
	p.line = p.input[p.pos : p.pos+end]
	p.pos += end + 1
}

func (p *parser) parseParagraph() (*Paragraph, error) {
	return &Paragraph{
		Text: p.parseInline(p.line),
	}, nil
}

func (p *parser) parseHeading() (*Heading, error) {
	level := 0
	for _, char := range p.line {
		if char == '#' {
			level++
		} else {
			break
		}
	}
	if level > 6 {
		level = 6
	}

	text := strings.TrimSpace(strings.TrimPrefix(p.line, strings.Repeat("#", level)))

	return &Heading{
		Level: level,
		Text:  p.parseInline(text),
	}, nil
}

func (p *parser) parseCodeBlock() (*CodeBlock, error) {
	language := strings.TrimSpace(strings.TrimPrefix(p.line, "```"))
	codeContent := ""

	for p.pos < len(p.input) {
		p.readLine()
		if strings.TrimSpace(p.line) == "```" {
			break
		}
		codeContent += p.line + "\n"
	}

	return &CodeBlock{
		Lang: language,
		Code: strings.TrimSpace(codeContent),
	}, nil
}

func (p *parser) parseInline(text string) []Inline {
	var inlines []Inline
	var currentText strings.Builder
	var i int

	for i < len(text) {
		// Nested bold and italic text
		if strings.HasPrefix(text[i:], "***") {
			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+3:], "***")
			if end != -1 {
				boldItalicText := text[i+3 : i+3+end]
				inlines = append(inlines, &BoldItalic{Content: p.parseInline(boldItalicText)})
				i += (3 + end + 3)
				continue
			}
		}

		// Parse bold text
		if strings.HasPrefix(text[i:], "**") {

			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+2:], "**")
			if end != -1 {
				boldText := text[i+2 : i+2+end]
				inlines = append(inlines, &Bold{Content: p.parseInline(boldText)})
				i += (2 + end + 2)
				continue
			}
		}

		// Parse italic text
		if strings.HasPrefix(text[i:], "_") {
			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+1:], "_")
			if end != -1 {
				italicText := text[i+1 : i+1+end]
				inlines = append(inlines, &Italic{Content: p.parseInline(italicText)})
				i += (1 + end + 1)
				continue
			}
		}

		currentText.WriteByte(text[i])
		i++
	}

	if currentText.Len() > 0 {
		inlines = append(inlines, &Text{Content: currentText.String()})
	}

	return inlines
}

func Parse(markdown string) (*Document, error) {
	normalizedMarkdown := strings.ReplaceAll(markdown, "\r\n", "\n")
	if !strings.HasSuffix(normalizedMarkdown, "\n") {
		normalizedMarkdown += "\n"
	}
	p := &parser{
		input: normalizedMarkdown,
		pos:   0,
	}
	return p.parse()
}
