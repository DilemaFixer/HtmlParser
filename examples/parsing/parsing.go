package main

import (
	"fmt"

	parser "github.com/DilemaFixer/HtmlParser"
)

const HTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Parser: mini-test</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <style>
    /* style for testing <style> */
    .hi { font-weight: bold }
  </style>
</head>
<body data-app="demo" data-version="1.0">
  <!-- Comment for parser test -->
  <h1 class="hi">Hello &amp; entity check: &lt;tag&gt;</h1>

  <p id="p1">Paragraph with <strong>bold</strong> and <em>italic</em>, also with an <a href="#p1">anchor</a>.</p>

  <img src="about:blank" alt="empty image" width="1" height="1" />

  <ul>
    <li data-i="1">first</li>
    <li data-i="2">second</li>
    <li data-i="3">third</li>
  </ul>

  <table>
    <thead><tr><th>Key</th><th>Value</th></tr></thead>
    <tbody>
      <tr><td>α</td><td>observation</td></tr>
      <tr><td>β</td><td>test</td></tr>
    </tbody>
  </table>

  <form action="/echo" method="get">
    <input name="q" value="text" />
    <input type="checkbox" checked />
    <button type="submit" disabled>Submit</button>
  </form>

  <br />

  <script>
    // minimal script for tag <script>
    document.body.setAttribute("data-seen", "true");
  </script>
</body>
</html>`

func main() {
	htmlParser := parser.NewHtmlParser()
	ast, err := htmlParser.ParseHtml(HTML)
	if err != nil {
		fmt.Println("Html parsing error:", err)
		return
	}

	for _, tag := range ast {
		parser.PrintHtmlTree(tag)
	}
}
