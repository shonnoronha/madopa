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

type Link struct {
    Text string
    URL string
}

func (t Text) isInline() {}
func (b Bold) isInline() {}
func (i Italic) isInline() {}
func (l Link) isInline() {}

