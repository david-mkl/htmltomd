package markdown

import (
	"fmt"
	"strings"
)

// Header represents a markdown header
type Header struct {
	headerType headerType
	Content    string
}

// List represents either an ordered or unorder markdown list
type List struct {
	ordinal string
	Items   []string
}

// Paragraph represents a block of string content
type Paragraph struct {
	Content string
}

// Codeblock represents preformatted text such as code.
// "lang" specifies the programming language of the code.
type CodeBlock struct {
	Lang string
	Code string
}

// HorizontalRule represents a markdown horizontal rule seporator
type HorizontalRule struct {
}

// Table represents a markdown table
type Table struct {
	Headers []string
	Rows    [][]string
}

// NewUnorderedList creates a new List with the unordered ordinal.
func NewUnorderedList(items []string) List {
	return List{ordinal: unorderedChar, Items: items}
}

// NewOrderedList creates a new List with the unordered ordinal.
func NewOrderedList(items []string) List {
	return List{ordinal: orderedChar, Items: items}
}

// String renders the header with the number of #'s according to the type
func (h Header) String() string {
	return string(h.headerType) + " " + h.Content
}

// String renders the items in the list and prefixes each item with the specificed ordinal
func (l List) String() string {
	listItems := make([]string, len(l.Items))
	for idx, itemContent := range l.Items {
		listItems[idx] = l.ordinal + " " + strings.Trim(itemContent, " ")
	}

	return strings.Join(listItems, "\n")
}

// String renders the paragraph content into a block of text
func (p Paragraph) String() string {
	return p.Content
}

// String wraps the codeblock in ``` with the specified language
func (cb CodeBlock) String() string {
	return strings.Join([]string{"```" + cb.Lang, cb.Code, "```"}, "\n")
}

// String renders the horizontal rule as "---"
func (hr HorizontalRule) String() string {
	return "---"
}

// String renders the table to to a markdown table.
// if "headers" are specified, then a row of dividers will be
// rendered after the headers.
func (t Table) String() string {
	var mdTable []string

	if len(t.Headers) > 0 {
		dividers := make([]string, len(t.Headers))
		for i := range t.Headers {
			dividers[i] = "---"
		}
		mdTable = append(mdTable, fmt.Sprintf("| %s |", strings.Join(t.Headers, " | ")))
		mdTable = append(mdTable, fmt.Sprintf("| %s |", strings.Join(dividers, " | ")))
	}

	for _, row := range t.Rows {
		mdTable = append(mdTable, fmt.Sprintf("| %s |", strings.Join(row, " | ")))
	}

	return strings.Join(mdTable, "\n")
}
