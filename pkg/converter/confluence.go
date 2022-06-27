package converter

import (
	"fmt"
	"strings"

	"github.com/david-mk-lawrence/htmltomd/pkg/markdown"
	"github.com/david-mk-lawrence/htmltomd/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

var (
	confluencePanelNoteClass    = "panel" // Confluence handle's "note" panels differently for some reason
	confluencePanelContentClass = "confluence-information-macro-body"
	confluencePanelInfoClass    = "confluence-information-macro-information"
	confluencePanelWarningClass = "confluence-information-macro-note"
	confluencePanelTipClass     = "confluence-information-macro-tip"
	confluencePanelErrorClass   = "confluence-information-macro-warning"
)

// ConfluenceSelectionConverter converts the Confluence HTML page to markdown.
// Tags controls which HTML tags will be searched when looking for content.
// If not set, then the defaultTags will be used.
type ConfluenceSelectionConverter struct {
	Transformer            *Transformer
	RootElementFinder      FindDocumentSelection
	TitleFinder            FindText
	ContentSelector        FindSelection
	ContentSelectorHandler HandleSelection
}

// NewConfluenceSelectionConverter intializes a ConfluenceSelectionConverter with default function calls.
func NewConfluenceSelectionConverter(conf SelectionConverterConfig) *ConfluenceSelectionConverter {
	c := &ConfluenceSelectionConverter{}

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
func (c *ConfluenceSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
	return c.RootElementFinder(doc)
}

// FindTitle finds the title of the document.
func (c *ConfluenceSelectionConverter) FindTitle(doc *goquery.Document) string {
	return c.TitleFinder(doc)
}

// FindContentElements finds the selections that that should be iterated over for content
func (c *ConfluenceSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
	return c.ContentSelector(s)
}

// HandleMatchedSelection handles matched selections from FindContentElements.
func (c *ConfluenceSelectionConverter) HandleMatchedSelection(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
	c.ContentSelectorHandler(i, elm, mdDoc, toMD)
}

func (c *ConfluenceSelectionConverter) defaultRootElementFinder(doc *goquery.Document) *goquery.Selection {
	return doc.Find("#main-content").First()
}

func (c *ConfluenceSelectionConverter) defaultTitleFinder(doc *goquery.Document) string {
	return CleanText(doc.Find("#title-text").First().Text())
}

func (c *ConfluenceSelectionConverter) defaultContentSelector(s *goquery.Selection) *goquery.Selection {
	return s.ChildrenFiltered(DefaultSearchPattern)
}

func (c *ConfluenceSelectionConverter) defaultContentSelectorHandler(i int, elm *goquery.Selection, mdDoc *markdown.Doc, toMD SelectionToMD) {
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
		if c.isPanel(elm) {
			mdDoc.AddDoc(c.toPanel(elm, mdDoc.GetRenderConfig(), toMD))
		} else if c.isCodeBlock(elm) {
			mdDoc.AddContent(c.toCodeBlock(elm))
		} else {
			// Recurse through the div
			mdDoc.AddDoc(toMD(elm, mdDoc.GetRenderConfig()))
		}
	}
}

func (c *ConfluenceSelectionConverter) toPanel(elm *goquery.Selection, docConf markdown.DocConfig, toMD SelectionToMD) *markdown.Doc {
	// Recursively convert the content in the panel since it may contain lists, code blocks, etc
	// which will have been missed by the root since they aren't direct children
	doc := toMD(elm.Find("."+confluencePanelContentClass).First(), docConf)

	// Only markdown for hugo will be converted to panels
	if c.Transformer.Format != "hugo" {
		return doc
	}

	noticeType := "note"
	if elm.HasClass(confluencePanelNoteClass) {
		noticeType = "note"
	} else if elm.HasClass(confluencePanelInfoClass) {
		noticeType = "info"
	} else if elm.HasClass(confluencePanelWarningClass) {
		noticeType = "warning"
	} else if elm.HasClass(confluencePanelTipClass) {
		noticeType = "tip"
	} else if elm.HasClass(confluencePanelErrorClass) {
		noticeType = "error"
	}

	// Wrap the content with another document with a single space as separator
	// This will place the shortcode wrappers directory before and after the content
	wrapper := markdown.NewDoc(markdown.DocConfig{Separator: util.String("\n")})
	wrapper.AddParagraph(fmt.Sprintf("{{%% notice %s %%}}", noticeType))
	wrapper.AddDoc(doc)
	wrapper.AddParagraph("{{% /notice %}}")

	return wrapper
}

func (c *ConfluenceSelectionConverter) toCodeBlock(elm *goquery.Selection) markdown.CodeBlock {
	preBlock := elm.Find("pre").First()

	lang := "txt"
	if dataParams, exist := preBlock.Attr("data-syntaxhighlighter-params"); exist {
		// example: "brush: php; gutter: false; theme: Confluence"
		for _, param := range strings.Split(dataParams, ";") {
			keyVal := strings.Split(param, ":")
			if strings.Trim(keyVal[0], " ") == "brush" {
				lang = strings.Trim(keyVal[1], " ")
			}
		}
	}

	// Confluence sets the default language as PHP on "pre" blocks, so we can't reliably know whether the
	// user actually specified PHP. Better to default to txt, since that will be more common than PHP
	if lang == "php" {
		lang = "txt"
	}

	return markdown.CodeBlock{Lang: lang, Code: preBlock.Text()}
}

func (c *ConfluenceSelectionConverter) isCodeBlock(elm *goquery.Selection) bool {
	return elm.HasClass("code")
}

func (c *ConfluenceSelectionConverter) isPanel(elm *goquery.Selection) bool {
	// This is checking that the *only* class is panel, which is different than HasClass
	// Confluence treats "note" panels differently than other types, and has only the "panel" class
	if class, _ := elm.Attr("class"); class == confluencePanelNoteClass {
		return true
	}

	return elm.HasClass("confluence-information-macro")
}
