package parser

type Inline interface {
	isInline()
}

type Text struct {
	Content string
}

type Bold struct {
	Content []Inline
}

type Italic struct {
	Content []Inline
}

type BoldItalic struct {
	Content []Inline
}

type Link struct {
	Text []Inline
	URL  string
}

type CodeInline struct {
	Content string
}

type Image struct {
	Alt   string
	Src   string
	Title string
}

func (t Text) isInline()       {}
func (b Bold) isInline()       {}
func (i Italic) isInline()     {}
func (l Link) isInline()       {}
func (b BoldItalic) isInline() {}
func (c CodeInline) isInline() {}
func (i Image) isInline()      {}
