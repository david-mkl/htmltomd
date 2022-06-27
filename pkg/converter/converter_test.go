package converter

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/david-mk-lawrence/htmltomd/pkg/markdown"
)

type TestSelectionConverter struct {
}

func (c *TestSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
	return doc.Find("body").First()
}

func (c *TestSelectionConverter) FindTitle(doc *goquery.Document) string {
	return CleanText(doc.Find("h1").Text())
}

func (c *TestSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
	return s.Find("p")
}

func (c *TestSelectionConverter) HandleMatchedSelection(i int, s *goquery.Selection, md *markdown.Doc, toMD SelectionToMD) {
	md.AddParagraph(CleanText(s.Text()))
}

func TestDocumentToMarkdown(t *testing.T) {
	doc := newTestDoc(`
<html>
	<body>
		<h1>Test Title</h1>
		<div>
			<p>Test Paragraph 1</p>
			<p>Test Paragraph 2</p>
		</div>
	</body>
</html>
`)

	s := TestSelectionConverter{}
	c := DocumentConverter{SelectionConv: &s}
	mdDoc := c.DocumentToMarkdown(doc)

	result := mdDoc.String()
	expected := "# Test Title\n\nTest Paragraph 1\n\nTest Paragraph 2"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
