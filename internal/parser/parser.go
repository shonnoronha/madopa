package parser

import (
	"fmt"
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

		doc.Blocks = append(doc.Blocks, block)
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

	if strings.Contains(trimmedLine, "|") {
		currentPos := p.pos
		currentLine := p.line

		if p.pos < len(p.input) {
			originalPos := p.pos
			p.readLine()

			if strings.Contains(p.line, "|") && strings.Contains(p.line, "-") {
				p.pos = currentPos
				p.line = currentLine
				return p.parseTable()
			}

			// Restore the pos to the original position
			p.pos = originalPos
			p.readLine()
		}
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
			} else {
				currentText.WriteString("***")
				i += 3
				continue
			}
		}

		// Parse bold text
		if strings.HasPrefix(text[i:], "**") || strings.HasPrefix(text[i:], "__") {
			marker := text[i : i+2]

			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+2:], marker)
			if end != -1 {
				boldText := text[i+2 : i+2+end]
				inlines = append(inlines, &Bold{Content: p.parseInline(boldText)})
				i += (2 + end + 2)
				continue
			} else {
				currentText.WriteString(marker)
				i += 2
				continue
			}
		}

		// Parse italic text
		if strings.HasPrefix(text[i:], "_") || strings.HasPrefix(text[i:], "*") {
			marker := text[i : i+1]

			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+1:], marker)
			if end != -1 {
				italicText := text[i+1 : i+1+end]
				inlines = append(inlines, &Italic{Content: p.parseInline(italicText)})
				i += (1 + end + 1)
				continue
			} else {
				currentText.WriteString(marker)
				i++
				continue
			}
		}

		// Parse link
		if strings.HasPrefix(text[i:], "[") {
			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+1:], "]")
			// fmt.Println(text[i+1:i+1+end], string(text[i]), string(text[i+end+1]))

			if end != -1 && text[i+2+end] == '(' {
				// fmt.Println("valid link")
				linkText := text[i+1 : i+1+end]
				linkEnd := strings.Index(text[i+1+end+1:], ")")
				if linkEnd != -1 {
					linkURL := text[i+end+3 : i+end+linkEnd+2]
					inlines = append(inlines, &Link{
						Text: p.parseInline(linkText),
						URL:  linkURL})
					i += (1 + end + 1 + linkEnd + 1)
					continue
				} else {
					currentText.WriteString("[")
					i++
					continue
				}

			} else {
				currentText.WriteString("[")
				i++
				continue
			}
		}

		// Parse inline code
		if strings.HasPrefix(text[i:], "`") {
			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			end := strings.Index(text[i+1:], "`")
			if end != -1 {
				codeText := text[i+1 : i+1+end]
				inlines = append(inlines, &CodeInline{Content: codeText})
				i += (1 + end + 1)
				continue
			} else {
				currentText.WriteString("`")
				i++
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

func (p *parser) parseTable() (*Table, error) {
	headerLine := p.line
	p.readLine()

	delimiterLine := p.line

	alignments := p.parseTableAlignments(delimiterLine)
	headerCells := p.parseTableRow(headerLine)

	if len(headerCells) != len(alignments) {
		return nil, fmt.Errorf("number of header cells (%d) doesn't match number of columns in delimiter row (%d)",
			len(headerCells), len(alignments))
	}

	var rows [][]TableCell

	for p.pos < len(p.input) {
		p.readLine()

		// check for the end of the table
		if p.line == "" || !strings.Contains(p.line, "|") {
			break
		}

		row := p.parseTableRow(p.line)
		if len(row) < len(headerCells) {
			padded := make([]TableCell, len(headerCells))
			copy(padded, row)
			row = padded
		} else if len(row) > len(headerCells) {
			row = row[:len(headerCells)]
		}

		rows = append(rows, row)
	}

	return &Table{
		Headers:    headerCells,
		Rows:       rows,
		Alignments: alignments,
	}, nil
}

func (p *parser) parseTableRow(line string) []TableCell {
	parts := strings.Split(line, "|")

	if parts[0] == "" {
		parts = parts[1:]
	}
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	cells := make([]TableCell, len(parts))
	for i, part := range parts {
		content := strings.TrimSpace(part)
		cells[i] = TableCell{
			Content: p.parseInline(content),
		}
	}
	return cells
}

func (p *parser) parseTableAlignments(delimiterLine string) []Alignment {
	parts := strings.Split(delimiterLine, "|")

	if parts[0] == "" {
		parts = parts[1:]
	}
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	alignments := make([]Alignment, len(parts))
	for i, part := range parts {
		part := strings.TrimSpace(part)
		if part == "" {
			continue
		}
		hasRight := strings.HasPrefix(part, ":")
		hasLeft := strings.HasSuffix(part, ":")
		if hasRight && hasLeft {
			alignments[i] = AlignCenter
		} else if hasLeft {
			alignments[i] = AlignLeft
		} else if hasRight {
			alignments[i] = AlignRight
		} else {
			alignments[i] = AlignDefault
		}
	}
	return alignments
}

func (p *parser) parseLink() (*Link, error) {
	return nil, nil
}
