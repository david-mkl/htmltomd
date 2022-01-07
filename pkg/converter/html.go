package converter

import (
	"github.com/david-mk-lawrence/html-to-md/pkg/markdown"

	"github.com/PuerkitoBio/goquery"
)

// HTMLSelectionConverter converts generic HTML pages to markdown
type HTMLSelectionConverter struct {
	Transformer            *Transformer
	RootElementFinder      FindDocumentSelection
	TitleFinder            FindText
	ContentSelector        FindSelection
	ContentSelectorHandler HandleSelection
}

// NewHTMLSelectionConverter intializes a HTMLSelectionConverter with default function calls.
func NewHTMLSelectionConverter(conf SelectionConverterConfig) *HTMLSelectionConverter {
	c := &HTMLSelectionConverter{}

	if conf.Transformer != nil {
		c.Transformer = conf.Transformer
	} else {
		c.Transformer = &Transformer{}
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
func (c *HTMLSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
	return c.RootElementFinder(doc)
}

// FindTitle finds the title of the document.
func (c *HTMLSelectionConverter) FindTitle(doc *goquery.Document) string {
	return c.TitleFinder(doc)
}

// FindContentElements finds the selections that that should be iterated over for content
func (c *HTMLSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
	return c.ContentSelector(s)
}

// HandleMatchedSelection handles matched selections from FindContentElements.
func (c *HTMLSelectionConverter) HandleMatchedSelection(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
	c.ContentSelectorHandler(i, elm, mdDoc, toMD)
}

func (c *HTMLSelectionConverter) defaultRootElementFinder(doc *goquery.Document) *goquery.Selection {
	return doc.Find("body").First()
}

func (c *HTMLSelectionConverter) defaultTitleFinder(doc *goquery.Document) string {
	return CleanText(doc.Find("head").First().ChildrenFiltered("title").First().Text())
}

func (c *HTMLSelectionConverter) defaultContentSelector(s *goquery.Selection) *goquery.Selection {
	return s.ChildrenFiltered(DefaultSearchPattern)
}

func (c *HTMLSelectionConverter) defaultContentSelectorHandler(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
	c.Transformer.RemoveScripts(elm)
	c.Transformer.ReplaceAll(elm)

	tag := elm.Nodes[0].Data
	switch tag {
	case "p", "span":
		mdDoc.AddParagraph(CleanText(elm.Text()))
	case "hr":
		mdDoc.AddHorizontalRule()
	case "h1", "h2", "h3", "h4", "h5", "h6":
		mdDoc.AddHeader(tag, CleanText(elm.Text()))
	case "ul", "ol":
		mdDoc.AddContent(c.Transformer.ToList(elm))
	case "table":
		mdDoc.AddContent(c.Transformer.ToTable(elm))
	case "div":
		// Recurse through the div
		mdDoc.AddDoc(toMD(elm, mdDoc.GetRenderConfig()))
	}
}
