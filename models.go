package htmlparser

type HtmlParser interface {
	ParseHtml(html string) ([]*HtmlTag, error)
}

type HtmlSerializer interface {
	RenderHtml(tags []*HtmlTag) string
}

type HtmlTag struct {
	Name          string
	InnerHtml     string
	InnerContent  string
	IsSelfClosing bool
	Attributes    map[string]HtmlAttribute
	Parent        *HtmlTag
	Children      []*HtmlTag
	Pos           Position
	htmlStart     int
}

type HtmlAttribute struct {
	Name         string
	Value        string
	IsValueExist bool
}

type Position struct {
	Line   int
	Column int
}
