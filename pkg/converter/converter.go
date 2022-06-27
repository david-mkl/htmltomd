// Converts HTML content into markdown content
package converter

import (
	"github.com/david-mk-lawrence/htmltomd/pkg/markdown"

	"github.com/PuerkitoBio/goquery"
)

// DefaultSearchPattern defines a default pattern to search for elements that will
// contain content for the markdown document
const DefaultSearchPattern = "p,span,hr,h1,h2,h3,h4,h5,h6,ul,ol,div,table"

// FindDocumentSelection is a callable that finds DOM elements in the given the Document
type FindDocumentSelection func(*goquery.Document) *goquery.Selection

// FindText is a callable that finds text in the given the HTMLDoc
type FindText func(*goquery.Document) string

// FindSelection is a callable that finds DOM elements in the given selection
type FindSelection func(*goquery.Selection) *goquery.Selection

// SelectionToMD is a callable that converts a selection to a markdown document
type SelectionToMD func(*goquery.Selection, markdown.DocConfig) *markdown.Doc

// HandleSelection is a callable that is given a selection, a markdown document to
// add to, and a callable to convert child elements to markdown documents
type HandleSelection func(int, *goquery.Selection, *markdown.Doc, SelectionToMD)

// DocumentConverter is a struct that can convert an HTML document into a markdown document
type DocumentConverter struct {
	SelectionConv SelectionConverter
}

// SelectionConverter is an interface that converts a style of HTML document to markdown.
// The interface allows for customization to handle a specific and known HTML structure.
type SelectionConverter interface {
	FindRootElement(*goquery.Document) *goquery.Selection
	FindTitle(*goquery.Document) string
	FindContentElements(*goquery.Selection) *goquery.Selection
	HandleMatchedSelection(int, *goquery.Selection, *markdown.Doc, SelectionToMD)
}

// SelectionConverterConfig contains parameters that a SelectionConvert will can use to be more customizable
type SelectionConverterConfig struct {
	Transformer            *Transformer
	RootElementFinder      FindDocumentSelection
	TitleFinder            FindText
	ContentSelector        FindSelection
	ContentSelectorHandler HandleSelection
}

// DocumentToMarkdown converts the HTML doc to markdown
func (c *DocumentConverter) DocumentToMarkdown(doc *goquery.Document) *markdown.Doc {
	root := c.SelectionConv.FindRootElement(doc)
	title := CleanText(c.SelectionConv.FindTitle(doc))
	mdDoc := c.SelectionToMarkdown(root, markdown.DocConfig{Title: &title})

	return mdDoc
}

// SelectionToMarkdown creates a new markdown document, and searches for content to add to the markdown doc.
// It hands off handling of matched selections to the SelectionConverter since it depends heavily
// on the HTML structure of the original document.
func (c *DocumentConverter) SelectionToMarkdown(elm *goquery.Selection, docConf markdown.DocConfig) *markdown.Doc {
	mdDoc := markdown.NewDoc(docConf)

	c.SelectionConv.FindContentElements(elm).Each(func(i int, elm *goquery.Selection) {
		c.SelectionConv.HandleMatchedSelection(i, elm, mdDoc, c.SelectionToMarkdown)
	})

	return mdDoc
}
