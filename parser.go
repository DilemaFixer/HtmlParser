package htmlparser

import (
	"fmt"
	"sort"
	"strings"
)

type DefaultHtmlParser struct {
	customHandlers map[string]func(*Scanner) (*HtmlTag, error)
}

func NewHtmlParser() HtmlParser {
	return &DefaultHtmlParser{
		customHandlers: make(map[string]func(*Scanner) (*HtmlTag, error)),
	}
}

func NewHtmlSerializer() HtmlSerializer {
	return &DefaultHtmlParser{
		customHandlers: make(map[string]func(*Scanner) (*HtmlTag, error)),
	}
}

func (parser *DefaultHtmlParser) AddCustomAttributeHandler(tagName string, handler func(*Scanner) (*HtmlTag, error)) {
	if handler == nil || strings.TrimSpace(tagName) == "" {
		return
	}
	parser.customHandlers[tagName] = handler
}

func (parser *DefaultHtmlParser) ParseHtml(html string) ([]*HtmlTag, error) {
	scanner := NewScanner(html)
	closeStack := NewStack[string]()
	var roots []*HtmlTag
	var current *HtmlTag = nil

	for !scanner.EOF() {
		scanner.SkipWhitespace()
		if scanner.EOF() {
			break
		}

		if scanner.Current() == '<' {
			if scanner.Position()+4 <= scanner.Len() && scanner.Slice(scanner.Position(), scanner.Position()+4) == "<!--" {
				scanner.ConsumeUntilString("-->")
				if !scanner.MatchString("-->") {
					return nil, fmt.Errorf("Html parsing error: unclosed comment at %s", scanner.Location())
				}
				continue
			}

			if scanner.Position()+9 <= scanner.Len() && scanner.Slice(scanner.Position(), scanner.Position()+9) == "<!DOCTYPE" {
				scanner.ConsumeUntil(func(r rune) bool { return r == '>' })
				if !scanner.Match('>') {
					return nil, fmt.Errorf("Html parsing error: unclosed DOCTYPE at %s", scanner.Location())
				}
				continue
			}

			if scanner.PeekNext() == '/' {
				if closeStack.IsEmpty() {
					return nil, fmt.Errorf("Html parsing error: superfluous closing tag at %s", scanner.Location())
				}

				htmlEndPos := scanner.Position()
				closingTag, err := parseClosingTag(scanner)
				if err != nil {
					return nil, err
				}

				expected, ok := closeStack.Pop()
				if !ok || expected != closingTag {
					return nil, fmt.Errorf("Html parsing error: invalid closing tag '%s', expected '%s' at %s", closingTag, expected, scanner.Location())
				}

				if current != nil && current.htmlStart > 0 && current.htmlStart < htmlEndPos {
					current.InnerHtml = strings.TrimSpace(scanner.Slice(current.htmlStart, htmlEndPos))
				}

				current = current.Parent
				continue
			}

			tag, err := parsingOpenTag(parser, scanner)
			if err != nil {
				return nil, err
			}

			if current == nil {
				roots = append(roots, tag)
			} else {
				tag.Parent = current
				current.Children = append(current.Children, tag)
			}

			if !tag.IsSelfClosing {
				current = tag
				closeStack.Push(tag.Name)
			}

		} else {
			contentStart := scanner.Position()
			scanner.ConsumeUntil(func(r rune) bool { return r == '<' })
			content := strings.TrimSpace(scanner.SliceFrom(contentStart))

			if content != "" && current != nil {
				current.InnerContent += content
			}
		}
	}

	if !closeStack.IsEmpty() {
		unclosed, _ := closeStack.Pop()
		return nil, fmt.Errorf("Html parsing error: unclosed tag '%s'", unclosed)
	}

	return roots, nil
}

func parseClosingTag(scanner *Scanner) (string, error) {
	if !(scanner.Match('<') && scanner.Match('/')) {
		return "", fmt.Errorf("expected '</' at %s", scanner.Location())
	}

	scanner.SkipWhitespace()
	tagName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '>' && r != ' ' && r != '\t' && r != '\n' && r != '\r'
	})

	if tagName == "" {
		return "", fmt.Errorf("empty tag name at %s", scanner.Location())
	}

	scanner.SkipWhitespace()
	if !scanner.Match('>') {
		return "", fmt.Errorf("expected '>' at %s", scanner.Location())
	}

	return tagName, nil
}

func parsingOpenTag(parser *DefaultHtmlParser, scanner *Scanner) (*HtmlTag, error) {
	startLine, startColumn := scanner.Line(), scanner.Column()
	if !scanner.Match('<') {
		return nil, fmt.Errorf("expected '<' at %s", scanner.Location())
	}

	scanner.SkipWhitespace()
	tagName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '>' && r != '/' && !isWhitespace(r)
	})
	if tagName == "" {
		return nil, fmt.Errorf("empty tag name at %s", scanner.Location())
	}

	if handler, ok := parser.customHandlers[tagName]; ok {
		scanner.SetLocation(startLine, startColumn)
		return handler(scanner)
	}

	if tagName == "style" || tagName == "script" {
		for scanner.Current() != '>' && !scanner.EOF() {
			scanner.Take()
		}
		if !scanner.Match('>') {
			return nil, fmt.Errorf("expected '>' at %s", scanner.Location())
		}

		contentStart := scanner.Position()
		endTag := fmt.Sprintf("</%s>", tagName)
		scanner.ConsumeUntilString(endTag)
		content := scanner.Slice(contentStart, scanner.Position())
		scanner.MatchString(endTag)

		return &HtmlTag{
			Name:          tagName,
			Attributes:    make(map[string]HtmlAttribute),
			Children:      nil,
			InnerHtml:     content,
			IsSelfClosing: true,
			Pos:           Position{Line: startLine, Column: startColumn},
		}, nil
	}

	tag := &HtmlTag{
		Name:       tagName,
		Attributes: make(map[string]HtmlAttribute),
		Children:   make([]*HtmlTag, 0),
		Pos:        Position{Line: startLine, Column: startColumn},
	}

	// --- атрибуты ---
	scanner.SkipWhitespace()
	for !scanner.EOF() && scanner.Current() != '>' && scanner.Current() != '/' {
		attr, err := parseAttribute(scanner)
		if err != nil {
			return nil, err
		}
		tag.Attributes[attr.Name] = attr
		scanner.SkipWhitespace()
	}

	if scanner.Current() == '/' {
		scanner.Take()
		tag.IsSelfClosing = true
	}

	if !scanner.Match('>') {
		return nil, fmt.Errorf("expected '>' at %s", scanner.Location())
	}

	tag.htmlStart = scanner.Position()
	return tag, nil
}

func parseAttribute(scanner *Scanner) (HtmlAttribute, error) {
	scanner.SkipWhitespace()
	attrName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '=' && r != '>' && r != '/' && !isWhitespace(r)
	})

	if attrName == "" {
		return HtmlAttribute{}, fmt.Errorf("empty attribute name at %s", scanner.Location())
	}

	attr := HtmlAttribute{Name: attrName, IsValueExist: false}
	scanner.SkipWhitespace()

	if scanner.Current() == '=' {
		scanner.Take()
		scanner.SkipWhitespace()
		attr.IsValueExist = true

		if scanner.Current() == '"' || scanner.Current() == '\'' {
			quote := scanner.Current()
			scanner.Take()
			valueStart := scanner.Position()
			scanner.ConsumeUntil(func(r rune) bool { return r == quote })
			attr.Value = strings.TrimSpace(scanner.SliceFrom(valueStart))
			if !scanner.Match(quote) {
				return HtmlAttribute{}, fmt.Errorf("unclosed attribute value at %s", scanner.Location())
			}
		} else {
			attr.Value = strings.TrimSpace(scanner.ConsumeWhile(func(r rune) bool {
				return !isWhitespace(r) && r != '>' && r != '/'
			}))
		}
	}

	return attr, nil
}

func PrintHtmlTree(tag *HtmlTag) {
	printHtmlTreeRecursive(tag, 0)
}

func printHtmlTreeRecursive(tag *HtmlTag, depth int) {
	if tag == nil {
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s<%s", indent, tag.Name)

	// стабильный порядок атрибутов
	keys := make([]string, 0, len(tag.Attributes))
	for k := range tag.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := tag.Attributes[k]
		if attr.IsValueExist {
			fmt.Printf(" %s=\"%s\"", attr.Name, attr.Value)
		} else {
			fmt.Printf(" %s", attr.Name)
		}
	}

	if tag.IsSelfClosing {
		fmt.Printf("/>\n")
		return
	}

	fmt.Printf(">")

	if tag.InnerContent != "" {
		fmt.Printf("%s", tag.InnerContent)
	}

	if len(tag.Children) > 0 {
		fmt.Printf("\n")
		for _, child := range tag.Children {
			printHtmlTreeRecursive(child, depth+1)
		}
		fmt.Printf("%s", indent)
	}

	fmt.Printf("</%s>\n", tag.Name)
}

func (parser *DefaultHtmlParser) RenderHtml(tags []*HtmlTag) string {
	var sb strings.Builder
	for _, tag := range tags {
		renderTag(&sb, tag)
	}
	return sb.String()
}

func renderTag(sb *strings.Builder, tag *HtmlTag) {
	if tag == nil {
		return
	}

	sb.WriteString("<")
	sb.WriteString(tag.Name)

	keys := make([]string, 0, len(tag.Attributes))
	for k := range tag.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := tag.Attributes[k]
		if attr.IsValueExist {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name, attr.Value))
		} else {
			sb.WriteString(" " + attr.Name)
		}
	}

	if tag.IsSelfClosing {
		sb.WriteString("/>")
		return
	}

	sb.WriteString(">")
	if tag.InnerContent != "" {
		sb.WriteString(tag.InnerContent)
	}
	for _, child := range tag.Children {
		renderTag(sb, child)
	}
	sb.WriteString("</")
	sb.WriteString(tag.Name)
	sb.WriteString(">")
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
