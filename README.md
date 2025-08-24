![Library not ready](https://img.shields.io/badge/Library-not%20ready%20for%20contributions-red)

# HtmlParser

A lightweight and flexible HTML parser written in Go that builds an Abstract Syntax Tree (AST) from HTML documents.

## Features

- ✅ Parse HTML documents into AST
- ✅ Handle self-closing tags
- ✅ Support for `<style>` and `<script>` raw content tags
- ✅ HTML comments parsing
- ✅ DOCTYPE declarations
- ✅ Custom tag handlers
- ✅ Tree manipulation (clone, remove)
- ✅ HTML serialization back to string
- ✅ Attribute type conversion utilities
- ✅ Position tracking for debugging

## Installation

```bash
go get github.com/DilemaFixer/HtmlParser
```

## Quick Start

```go
package main

import (
    "fmt"
    parser "github.com/DilemaFixer/HtmlParser"
)

func main() {
    html := `<div class="container">
        <h1>Hello World</h1>
        <p>This is a paragraph with <strong>bold</strong> text.</p>
    </div>`

    htmlParser := parser.NewHtmlParser()
    ast, err := htmlParser.ParseHtml(html)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the AST tree
    for _, tag := range ast {
        parser.PrintHtmlTree(tag)
    }
}
```

## Data Structures

### HtmlTag
```go
type HtmlTag struct {
    Name          string                    // Tag name (e.g., "div", "p")
    InnerHtml     string                    // Raw HTML content between tags
    InnerContent  string                    // Text content between tags
    IsSelfClosing bool                      // Whether the tag is self-closing
    Attributes    map[string]HtmlAttribute  // Tag attributes
    Parent        *HtmlTag                  // Parent tag reference
    Children      []*HtmlTag                // Child tags
    Pos           Position                  // Position in source document
}
```

### HtmlAttribute
```go
type HtmlAttribute struct {
    Name         string  // Attribute name
    Value        string  // Attribute value
    IsValueExist bool    // Whether the attribute has a value
}
```

### Position
```go
type Position struct {
    Line   int  // Line number in source
    Column int  // Column number in source
}
```

## Core Functionality

### Parsing HTML

```go
htmlParser := parser.NewHtmlParser()
ast, err := htmlParser.ParseHtml(htmlString)
if err != nil {
    // Handle parsing error
}

// ast is []*HtmlTag - slice of root elements
```

### Working with Attributes

```go
// Check if attribute exists
if tag.HasAttribute("class") {
    // Get attribute
    attr := tag.GetAttribute("class")
    if attr != nil {
        className := attr.AsString()
    }
}

// Set attribute
tag.SetAttribute("id", "my-id")

// Remove attribute
tag.RemoveAttribute("class")

// Type conversions
widthAttr := tag.GetAttribute("width")
if widthAttr != nil {
    width, err := widthAttr.AsInt()
    if err == nil {
        // use width as integer
    }
}
```

### Supported Type Conversions
- `AsString()` - string
- `AsBool()` - bool with error
- `AsInt()`, `AsInt8()`, `AsInt16()`, `AsInt32()`, `AsInt64()` - integers with error
- `AsUint()`, `AsUint8()`, `AsUint16()`, `AsUint32()`, `AsUint64()` - unsigned integers with error
- `AsFloat32()`, `AsFloat64()` - floats with error

### Tree Manipulation

```go
// Clone tag with all children up to specified depth
clonedTag, err := tag.CloneDown(2) // Clone 2 levels deep

// Clone tag and its parents up to specified depth
rootClone, err := tag.CloneUp(3, false) // Clone up 3 parent levels

// Remove child from parent
parent.RemoveChild(childTag)
```

### HTML Serialization

```go
// Serialize AST back to HTML string
serializer := parser.NewHtmlSerializer()
htmlString := serializer.RenderHtml(ast)
fmt.Println(htmlString)
```

### Custom Tag Handlers

```go
htmlParser := parser.NewHtmlParser()

// Add custom handler for specific tags
htmlParser.AddCustomAttributeHandler("my-tag", func(scanner *Scanner) (*HtmlTag, error) {
    // Custom parsing logic for <my-tag>
    // Return custom HtmlTag
})
```

## Advanced Examples

### Complete Parsing Example

```go
package main

import (
    "fmt"
    parser "github.com/DilemaFixer/HtmlParser"
)

const HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>Sample Page</title>
    <style>
        .highlight { background: yellow; }
    </style>
</head>
<body>
    <!-- This is a comment -->
    <div class="container" id="main">
        <h1 class="highlight">Welcome</h1>
        <p>Paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
        <img src="image.jpg" alt="Description" width="100" height="50" />
        <ul>
            <li data-id="1">Item 1</li>
            <li data-id="2">Item 2</li>
        </ul>
    </div>
    <script>
        console.log("Hello from script");
    </script>
</body>
</html>`

func main() {
    htmlParser := parser.NewHtmlParser()
    ast, err := htmlParser.ParseHtml(HTML)
    if err != nil {
        fmt.Println("Parsing error:", err)
        return
    }

    // Find specific elements
    findDivs(ast)

    // Modify and serialize
    modifyAndSerialize(ast)
}

func findDivs(tags []*HtmlTag) {
    for _, tag := range tags {
        if tag.Name == "div" {
            fmt.Printf("Found div with class: %s\n", 
                tag.GetAttribute("class").AsString())
        }
        // Recursively search children
        findDivs(tag.Children)
    }
}

func modifyAndSerialize(ast []*HtmlTag) {
    // Find and modify the first h1
    for _, tag := range ast {
        h1 := findFirstH1(tag)
        if h1 != nil {
            h1.SetAttribute("style", "color: red;")
            h1.InnerContent = "Modified Title"
            break
        }
    }

    // Serialize back to HTML
    serializer := parser.NewHtmlSerializer()
    html := serializer.RenderHtml(ast)
    fmt.Println("Modified HTML:")
    fmt.Println(html)
}

func findFirstH1(tag *HtmlTag) *HtmlTag {
    if tag.Name == "h1" {
        return tag
    }
    for _, child := range tag.Children {
        if result := findFirstH1(child); result != nil {
            return result
        }
    }
    return nil
}
```

### Working with Forms

```go
// Find all form inputs
func extractFormData(tag *HtmlTag) map[string]string {
    formData := make(map[string]string)
    
    if tag.Name == "input" {
        nameAttr := tag.GetAttribute("name")
        valueAttr := tag.GetAttribute("value")
        
        if nameAttr != nil && valueAttr != nil {
            formData[nameAttr.AsString()] = valueAttr.AsString()
        }
    }
    
    // Recursively process children
    for _, child := range tag.Children {
        childData := extractFormData(child)
        for k, v := range childData {
            formData[k] = v
        }
    }
    
    return formData
}
```

## Error Handling

The parser provides detailed error messages with line and column information:

```go
ast, err := htmlParser.ParseHtml(invalidHtml)
if err != nil {
    fmt.Printf("Parsing failed: %s\n", err.Error())
    // Output: Html parsing error: unclosed tag 'div' at 5:12
}
```

## Supported HTML Features

- ✅ Standard HTML tags
- ✅ Self-closing tags (`<img />`, `<br />`)
- ✅ Raw content tags (`<style>`, `<script>`)
- ✅ HTML comments (`<!-- comment -->`)
- ✅ DOCTYPE declarations
- ✅ Attributes with and without values
- ✅ Quoted and unquoted attribute values
- ✅ Nested tag structures
- ✅ Text content preservation

## License

This project is part of the HtmlPuzzles project.

## Contributing

Feel free to submit issues and pull requests to improve the parser functionality.
