// Package markdown provides structs to build and render markdown documents
package markdown

import (
	"fmt"
	"strings"

	"github.com/david-mk-lawrence/htmltomd/pkg/util"
)

type headerType string

const (
	h1 headerType = "#"
	h2 headerType = "##"
	h3 headerType = "###"
	h4 headerType = "####"
	h5 headerType = "#####"
	h6 headerType = "######"

	unorderedChar = "*"
	orderedChar   = "1."

	defaultSeparator = "\n\n"
)

var (
	headerMap = map[string]headerType{
		"h1": h1,
		"h2": h2,
		"h3": h3,
		"h4": h4,
		"h5": h5,
		"h6": h6,
	}

	// headers are reduced one level to make
	// room for the page title
	reducedHeaderMap = map[string]headerType{
		"h1": h2,
		"h2": h3,
		"h3": h4,
		"h4": h5,
		"h5": h6,
		"h6": h6,
	}
)

// Doc represents a markdown document.
// It holds a list of fmt.Stringer which represent
// each block of text that will be rendered
type Doc struct {
	content       []fmt.Stringer
	title         *string
	separator     *string
	reduceHeaders *bool
}

// DocConfig contains parameters that are used to intialize a new Doc.
type DocConfig struct {
	Title         *string
	ReduceHeaders *bool
	Separator     *string
}

// NewDoc intializes a new Doc.
func NewDoc(conf DocConfig) *Doc {
	doc := &Doc{title: conf.Title}

	if conf.ReduceHeaders == nil {
		doc.reduceHeaders = util.Bool(true)
	} else {
		doc.reduceHeaders = conf.ReduceHeaders
	}
	if conf.Separator == nil {
		doc.separator = util.String(defaultSeparator)
	} else {
		doc.separator = conf.Separator
	}

	return doc
}

// GetConfig retrieves the documentation config.
func (d *Doc) GetConfig() DocConfig {
	return DocConfig{
		Title:         d.title,
		Separator:     d.separator,
		ReduceHeaders: d.reduceHeaders,
	}
}

// GetRenderConfig retrieves the config without document content config
// like "title". This is convenient when needing to preserve config for
// rendering only when creating a child document from a parent. Otherwise,
// if all the config is copied to the child, then "title" would be rendered
// twice.
func (d *Doc) GetRenderConfig() DocConfig {
	return DocConfig{
		Separator:     d.separator,
		ReduceHeaders: d.reduceHeaders,
	}
}

// AddContent adds a block to the document.
func (d *Doc) AddContent(content fmt.Stringer) {
	d.content = append(d.content, content)
}

// AddDoc adds another markdown document as a block to this document.
func (d *Doc) AddDoc(subdoc *Doc) {
	d.AddContent(subdoc)
}

// AddHeader adds a section header to the document.
// If the content is empty, it will not be added to the document.
// If the doc has reduceHeaders set, then the headerTag
// will be reduced one level.
func (d *Doc) AddHeader(headerTag string, content string) {
	if len(content) > 0 {
		var hType headerType
		if d.reduceHeaders != nil && *d.reduceHeaders {
			hType = reducedHeaderMap[headerTag]
		} else {
			hType = headerMap[headerTag]
		}

		d.AddContent(Header{
			headerType: hType,
			Content:    content,
		})
	}
}

// AddUnorderedList adds a List to the document.
// It sets the unordered prefix ordinal.
func (d *Doc) AddUnorderedList(items []string) {
	d.AddContent(NewUnorderedList(items))
}

// AddUnorderedList adds a List to the document.
// It sets the ordered prefix ordinal.
func (d *Doc) AddOrderedList(items []string) {
	d.AddContent(NewOrderedList(items))
}

// AddParagraph adds a block of text to the document.
// If the content being added is empty, it will not be
// added to the document.
func (d *Doc) AddParagraph(content string) {
	if len(content) > 0 {
		d.AddContent(Paragraph{Content: content})
	}
}

// AddCodeBlock adds a block of code to the document.
func (d *Doc) AddCodeBlock(lang string, code string) {
	d.AddContent(CodeBlock{
		Lang: lang,
		Code: code,
	})
}

// AddHorizontalRule adds a horizontal rule to the document
func (d *Doc) AddHorizontalRule() {
	d.AddContent(HorizontalRule{})
}

// AddTable adds a table to the document.
func (d *Doc) AddTable(headers []string, rows [][]string) {
	d.AddContent(Table{Headers: headers, Rows: rows})
}

// Content renders just the content of the document.
// It does not include the title.
func (d Doc) Content() string {
	var contentLines []string
	for _, content := range d.content {
		rendered := content.String()
		if len(rendered) > 0 {
			contentLines = append(contentLines, rendered)
		}
	}

	return strings.Join(contentLines, d.getSeparator())
}

// Title renders just the title of the document
func (d *Doc) Title() string {
	title := ""
	if d.title != nil {
		title = string(h1) + " " + *d.title
	}
	return title
}

// String renders the title with the content.
func (d *Doc) String() string {
	title := d.Title()
	if title != "" {
		title += d.getSeparator()
	}

	return title + d.Content()
}

func (d *Doc) getSeparator() string {
	sep := defaultSeparator
	if d.separator != nil {
		sep = *d.separator
	}
	return sep
}
