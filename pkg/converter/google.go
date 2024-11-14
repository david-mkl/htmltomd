package converter

import (
	"strings"

	"github.com/david-mk-lawrence/htmltomd/pkg/markdown"

	"github.com/PuerkitoBio/goquery"
)

// GoogleSelectionConverter converts the Google Doc HTML page to markdown
type GoogleSelectionConverter struct {
	Transformer            *Transformer
	RootElementFinder      FindDocumentSelection
	TitleFinder            FindText
	ContentSelector        FindSelection
	ContentSelectorHandler HandleSelection
}

// NewGoogleSelectionConverter intializes a GoogleSelectionConverter with default function calls.
func NewGoogleSelectionConverter(conf SelectionConverterConfig) *GoogleSelectionConverter {
	c := &GoogleSelectionConverter{}

	if conf.Transformer != nil {
		c.Transformer = conf.Transformer
		if c.Transformer.textCleaner == nil {
			c.Transformer.textCleaner = NewTextCleaner(nil)
		}
	} else {
		c.Transformer = NewTransformer(nil)
	}

	if conf.RootElementFinder != nil {
		c.RootElementFinder = conf.RootElementFinder
	} else {
		c.RootElementFinder = c.defaultRootElementFinder
	}

	if conf.TitleFinder != nil {
		c.TitleFinder = conf.TitleFinder
	} else {
		c.TitleFinder = c.defaultTitleFinder
	}

	if conf.ContentSelector != nil {
		c.ContentSelector = conf.ContentSelector
	} else {
		c.ContentSelector = c.defaultContentSelector
	}

	if conf.ContentSelectorHandler != nil {
		c.ContentSelectorHandler = conf.ContentSelectorHandler
	} else {
		c.ContentSelectorHandler = c.defaultContentSelectorHandler
	}

	return c
}

// FindRootElement finds the root element.
func (c *GoogleSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
	return c.RootElementFinder(doc)
}

// FindTitle finds the title of the document.
func (c *GoogleSelectionConverter) FindTitle(doc *goquery.Document) string {
	return c.TitleFinder(doc)
}

// FindContentElements finds the selections that that should be iterated over for content
func (c *GoogleSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
	return c.ContentSelector(s)
}

// HandleMatchedSelection handles matched selections from FindContentElements.
func (c *GoogleSelectionConverter) HandleMatchedSelection(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
	c.ContentSelectorHandler(i, elm, mdDoc, toMD)
}

func (c *GoogleSelectionConverter) defaultRootElementFinder(doc *goquery.Document) *goquery.Selection {
	return doc.Find("body").First()
}

func (c *GoogleSelectionConverter) defaultTitleFinder(doc *goquery.Document) string {
	return c.Transformer.CleanText(doc.Find("head").First().ChildrenFiltered("title").First().Text())
}

func (c *GoogleSelectionConverter) defaultContentSelector(s *goquery.Selection) *goquery.Selection {
	return s.ChildrenFiltered(DefaultSearchPattern)
}

func (c *GoogleSelectionConverter) defaultContentSelectorHandler(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
	c.Transformer.RemoveScripts(elm)
	c.Transformer.ReplaceAll(elm)

	tag := elm.Nodes[0].Data
	switch tag {
	case "p", "span":
		mdDoc.AddParagraph(c.Transformer.CleanText(elm.Text()))
	case "hr":
		// Ignore the page breaks that indicate a new page in the google doc
		if !c.isPageBreak(elm) {
			mdDoc.AddHorizontalRule()
		}
	case "h1", "h2", "h3", "h4", "h5", "h6":
		mdDoc.AddHeader(tag, c.Transformer.CleanText(elm.Text()))
	case "ul", "ol":
		mdDoc.AddContent(c.Transformer.ToList(elm))
	case "table":
		mdDoc.AddContent(c.Transformer.ToTable(elm))
	case "div":
		// Recurse through the div
		mdDoc.AddDoc(toMD(elm, mdDoc.GetRenderConfig()))
	}
}

func (c *GoogleSelectionConverter) isPageBreak(elm *goquery.Selection) bool {
	if style, exists := elm.Attr("style"); exists {
		if strings.Contains(style, "page-break") {
			return true
		}
	}
	return false
}
