package parser

import (
	"fmt"
	"regexp"
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

		block, err := p.parseBlock(p.line)
		// fmt.Printf("Block type: %T\n", block)
		if err != nil {
			return nil, err
		}

		doc.Blocks = append(doc.Blocks, block)
	}

	return doc, nil
}

func (p *parser) parseBlock(line string) (Block, error) {
	trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)

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

	if strings.HasPrefix(trimmedLine, ">") {
		return p.parseBlockquote()
	}

	_, _, ok, _ := p.parseListItem(p.line)
	if ok {
		return p.parseList(p.line)
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

		// Parse image
		if strings.HasPrefix(text[i:], "!") {
			if currentText.Len() > 0 {
				inlines = append(inlines, &Text{Content: currentText.String()})
				currentText.Reset()
			}

			if strings.HasPrefix(text[i+1:], "[") {
				end := strings.Index(text[i+2:], "]")
				if end != -1 && strings.HasPrefix(text[i+2+end+1:], "(") {
					altText := text[i+2 : i+2+end]

					urlStart := i + 2 + end + 2
					urlAndTitle := text[urlStart:]
					urlAndTitleEnd := strings.Index(urlAndTitle, ")")

					if urlAndTitleEnd != -1 {
						urlAndTitle = urlAndTitle[:urlAndTitleEnd]

						titleMatch := regexp.MustCompile(`^(.*?)\s+(.*?)$`).FindStringSubmatch(urlAndTitle)

						var linkURL, title string

						if titleMatch != nil {
							linkURL = titleMatch[1]
							title = titleMatch[2]

							if strings.HasPrefix(title, "\"") && strings.HasSuffix(title, "\"") {
								title = title[1 : len(title)-1]
							}
						} else {
							linkURL = urlAndTitle
						}

						inlines = append(inlines,
							&Image{
								Alt:   altText,
								Src:   linkURL,
								Title: title,
							},
						)

						i += (2 + end + 2 + urlAndTitleEnd + 1)
						continue
					}

				}
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

func (p *parser) parseList(line string) (*List, error) {
	var items []*ListItem
	var listType ListType

	trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
	if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") {
		listType = UnorderedList
	} else {
		listType = OrderedList
	}

	for p.pos <= len(p.input) {
		inlineElements, level, isListItem, itemType := p.parseListItem(p.line)

		if !isListItem {
			break
		}

		newItem := &ListItem{
			Level:   level,
			Content: inlineElements,
		}

		if len(items) == 0 || level == 0 {
			items = append(items, newItem)
		} else {
			parent := FindListItemParent(items, level)
			if parent != nil {
				if parent.Children == nil {
					parent.Children = &List{
						Items: []*ListItem{newItem},
						Type:  itemType,
					}
				} else {
					parent.Children.Items = append(parent.Children.Items, newItem)
				}
			} else {
				items = append(items, newItem)
			}
		}

		if p.pos >= len(p.input) {
			break
		}
		p.readLine()

		if strings.TrimSpace(p.line) == "" {
			break
		}

		_, _, nextIsListItem, _ := p.parseListItem(p.line)
		if !nextIsListItem {
			break
		}
	}
	return &List{
		Items: items,
		Type:  listType,
	}, nil
}

func (p *parser) parseListItem(line string) ([]Inline, int, bool, ListType) {
	trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
	indentation := len(line) - len(trimmedLine)
	level := indentation / 2

	if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") {
		text := strings.TrimSpace(trimmedLine[2:])
		return p.parseInline(text), level, true, UnorderedList
	}

	if match := regexp.MustCompile(`^\d+\.\s+`).FindString(trimmedLine); match != "" {
		text := strings.TrimSpace(trimmedLine[len(match):])
		return p.parseInline(text), level, true, OrderedList
	}

	return nil, 0, false, 0
}

func FindListItemParent(items []*ListItem, level int) *ListItem {
	if len(items) == 0 {
		return nil
	}

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		if item.Level == level-1 {
			return item
		}

		if item.Children != nil && item.Level < level {
			return FindListItemParent(item.Children.Items, level)
		}
	}

	return nil
}

func FindBlockQuoteItemParent(items []*BlockquoteItem, level int) *BlockquoteItem {
	if len(items) == 0 {
		return nil
	}

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		if item.Level == level-1 {
			return item
		}

		if item.Children != nil && item.Level < level {
			return FindBlockQuoteItemParent(item.Children.Items, level)
		}
	}
	return nil
}

func (p *parser) parseBlockquote() (*Blockquote, error) {
	blockquote := &Blockquote{}

	blockquote.Items = append(blockquote.Items, &BlockquoteItem{
		Content: p.parseInline(p.line[1:]),
		Level:   1,
	})

	for p.pos < len(p.input) {
		p.readLine()
		trimmedLine := strings.TrimSpace(p.line)

		if strings.HasPrefix(trimmedLine, ">") {
			level := strings.Count(trimmedLine, ">")
			trimmedLine = strings.TrimSpace(p.line[2*level-1:])
			if trimmedLine != "" {
				newItem := &BlockquoteItem{
					Content: p.parseInline(trimmedLine),
					Level:   level,
				}
				if level == 1 {
					blockquote.Items = append(blockquote.Items, newItem)
				} else {
					parent := FindBlockQuoteItemParent(blockquote.Items, level)
					if parent != nil {
						if parent.Children == nil {
							parent.Children = &Blockquote{
								Items: []*BlockquoteItem{newItem},
							}
						} else {
							parent.Children.Items = append(parent.Children.Items, newItem)
						}
					}
				}
				fmt.Println(level, trimmedLine)
			}

		} else {
			break
		}
	}
	return blockquote, nil
}
