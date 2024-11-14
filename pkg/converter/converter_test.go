package converter

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/david-mk-lawrence/htmltomd/pkg/markdown"
)

type TestSelectionConverter struct {
	Transformer *Transformer
}

func (c *TestSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
	return doc.Find("body").First()
}

func (c *TestSelectionConverter) FindTitle(doc *goquery.Document) string {
	return c.Transformer.CleanText(doc.Find("h1").Text())
}

func (c *TestSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
	return s.Find("p")
}

func (c *TestSelectionConverter) HandleMatchedSelection(i int, s *goquery.Selection, md *markdown.Doc, toMD SelectionToMD) {
	md.AddParagraph(c.Transformer.CleanText(s.Text()))
}

func TestDocumentToMarkdown(t *testing.T) {
	doc := newTestDoc(`
<html>
	<body>
		<h1>Těst Title</h1>
		<div>
			<p>Test Paragraph 1</p>
			<p>Test Pâragraph 2</p>
		</div>
	</body>
</html>
`)

	tc := NewTextCleaner(&TextCleanerConf{
		AsciiOnly: true,
	})
	s := &TestSelectionConverter{
		Transformer: NewTransformer(&TransformerConf{
			TextCleaner: tc,
		}),
	}
	d := &DocumentConverterConf{
		TextCleaner: tc,
	}
	c := NewDocumentConverter(s, d)
	mdDoc := c.DocumentToMarkdown(doc)

	result := mdDoc.String()
	expected := "# Tst Title\n\nTest Paragraph 1\n\nTest Pragraph 2"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestDocumentToMarkdownWithUnicode(t *testing.T) {
	doc := newTestDoc(`
<html>
	<body>
		<h1>Tèsț Tïtlē</h1>
		<div>
			<p>Tëst Paragrãph 1</p>
			<p>Těst Parägřaph 2</p>
		</div>
	</body>
</html>
`)

	tc := NewTextCleaner(&TextCleanerConf{
		AsciiOnly: false,
	})
	s := &TestSelectionConverter{
		Transformer: NewTransformer(&TransformerConf{
			TextCleaner: tc,
		}),
	}
	d := &DocumentConverterConf{
		TextCleaner: tc,
	}
	c := NewDocumentConverter(s, d)
	mdDoc := c.DocumentToMarkdown(doc)

	result := mdDoc.String()
	expected := "# Tèsț Tïtlē\n\nTëst Paragrãph 1\n\nTěst Parägřaph 2"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
