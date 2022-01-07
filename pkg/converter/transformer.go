package converter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/david-mk-lawrence/html-to-md/pkg/markdown"

	"github.com/PuerkitoBio/goquery"
)

var asciiFilter = regexp.MustCompile("[[:^ascii:]]")

// SelectionCallback is a function that handles a goquery.Selection
type SelectionCallback = func(i int, s *goquery.Selection)

// Transformer converts HTML DOM elements into markdown elements
type Transformer struct {
	Format string
}

// RemoveScripts removes any script, style, or link tags from the DOM element.
func (t *Transformer) RemoveScripts(elm *goquery.Selection) {
	elm.Find("*").RemoveFiltered("style,script,link")
}

// Transform finds all elements matching the pattern and calls
// each given callback on each child element.
func (t *Transformer) Transform(pattern string, elm *goquery.Selection, callbacks ...SelectionCallback) {
	elm.Find(pattern).Each(func(i int, s *goquery.Selection) {
		t.Transforms(i, s, callbacks...)
	})
}

// Transforms calls each callback on the given DOM element.
func (t *Transformer) Transforms(i int, s *goquery.Selection, callbacks ...SelectionCallback) {
	for _, cb := range callbacks {
		cb(i, s)
	}
}

// ToList transforms the "ul" or "ol" dom element to a markdown List.
func (t *Transformer) ToList(list *goquery.Selection) markdown.List {
	var items []string
	tag := list.Nodes[0].Data
	list.ChildrenFiltered("li").Each(func(i int, li *goquery.Selection) {
		items = append(items, CleanText(li.Text()))
	})

	if tag == "ol" {
		return markdown.NewOrderedList(items)
	}
	return markdown.NewUnorderedList(items)
}

// ToTable transforms the "table" dom element to a markdown Table.
func (t *Transformer) ToTable(table *goquery.Selection) markdown.Table {
	headerElms := getTableHeaders(table)
	rowElms := getTableRows(table)

	headers := make([]string, len(headerElms.Nodes))
	headerElms.Each(func(i int, th *goquery.Selection) {
		headers[i] = CleanText(th.Text())
	})

	var rows [][]string

	rowElms.Each(func(i int, tr *goquery.Selection) {
		cellElms := tr.Find("td")
		cells := make([]string, len(cellElms.Nodes))
		cellElms.Each(func(j int, td *goquery.Selection) {
			cells[j] = CleanText(td.Text())
		})
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	})

	return markdown.Table{Headers: headers, Rows: rows}
}

func getTableHeaders(table *goquery.Selection) (headerElms *goquery.Selection) {
	thead := table.Find("thead")
	if len(thead.Nodes) > 0 {
		var headerRoot *goquery.Selection

		firstRow := thead.Find("tr").First()
		if len(firstRow.Nodes) > 0 {
			headerRoot = firstRow
		} else {
			headerRoot = thead
		}

		headerElms = headerRoot.Find("th")
		if len(headerElms.Nodes) == 0 {
			headerElms = headerRoot.Find("td")
		}
	} else {
		headerElms = table.Find("th")
	}

	return
}

func getTableRows(table *goquery.Selection) (rowElms *goquery.Selection) {
	tbody := table.Find("tbody")
	if len(tbody.Nodes) > 0 {
		rowElms = tbody.Find("tr")
	} else {
		rowElms = table.Find("tr")
	}
	return
}

// ReplaceAll runs all the default replacement functions
func (t *Transformer) ReplaceAll(elm *goquery.Selection) {
	t.ReplaceBolds(elm)
	t.ReplaceItalics(elm)
	t.ReplaceAnchors(elm)
	t.ReplaceInlineCodes(elm)
	t.ReplaceImages(elm)
}

// ReplaceAnchors finds all child "a" tags and replaces them in place with markdown links.
func (t *Transformer) ReplaceAnchors(elm *goquery.Selection) {
	t.Transform("a", elm, t.ReplaceAnchor)
}

// ReplaceAnchor replaces the DOM element in place with a markdown link.
func (t *Transformer) ReplaceAnchor(i int, s *goquery.Selection) {
	if href, exists := s.Attr("href"); exists {
		text := CleanText(s.Text())
		s.ReplaceWithHtml(fmt.Sprintf("[%s](%s)", text, href))
	}
}

// ReplaceImages finds all child "img" tags and replaces them in place with markdown image links.
func (t *Transformer) ReplaceImages(elm *goquery.Selection) {
	t.Transform("img", elm, t.ReplaceImage)
}

// ReplaceImage replaces the DOM element in place with a markdown image link.
// If the Transformer is rendering for Hugo, then will replace with a Hugo figure shortcode.
func (t *Transformer) ReplaceImage(i int, s *goquery.Selection) {
	if src, exists := s.Attr("src"); exists {
		alt, _ := s.Attr("alt")
		if t.Format == "hugo" {
			s.ReplaceWithHtml(fmt.Sprintf("{{< figure src=\"./%s\" alt=\"%s\" >}}", src, alt))
		} else {
			s.ReplaceWithHtml(fmt.Sprintf("![%s](%s)", alt, src))
		}
	}
}

// ReplaceInlineCodes finds all child "code" tags and replaces them in place with text content wrapped in "`".
func (t *Transformer) ReplaceInlineCodes(elm *goquery.Selection) {
	t.Transform("code", elm, t.ReplaceInlineCode)
}

// ReplaceInlineCode replaces the DOM element in place with text content wrapped in "`".
func (t *Transformer) ReplaceInlineCode(i int, s *goquery.Selection) {
	html, _ := s.Html()
	s.ReplaceWithHtml(fmt.Sprintf("`%s`", CleanText(html)))
}

// ReplaceItalics finds all child "em" tags and replaces them in place with markdown italics.
func (t *Transformer) ReplaceItalics(elm *goquery.Selection) {
	t.Transform("em", elm, t.ReplaceItalic)
}

// ReplaceItalic replaces the DOM element in place with the text content wrapped in "_".
func (t *Transformer) ReplaceItalic(i int, s *goquery.Selection) {
	s.ReplaceWithHtml(fmt.Sprintf("_%s_", CleanText(s.Text())))
}

// ReplaceBolds finds all child "strong" tags and replaces them in place with markdown bold.
func (t *Transformer) ReplaceBolds(elm *goquery.Selection) {
	t.Transform("strong", elm, t.ReplaceBold)
}

// ReplaceBold replaces the DOM element in place with the text content wrapped in "**".
func (t *Transformer) ReplaceBold(i int, s *goquery.Selection) {
	s.ReplaceWithHtml(fmt.Sprintf("**%s**", CleanText(s.Text())))
}

// CleanText removes newlines, replaces common unicode characters with
// ascii, removes any other non-common ascii values, and trims any whitespace.
func CleanText(content string) string {
	lines := strings.Split(content, "\n")
	trimmed := make([]string, len(lines))
	for idx, line := range lines {
		// Replace invisible spaces
		line = strings.ReplaceAll(line, "\u00a0", " ")
		// Replace quotes
		line = strings.ReplaceAll(line, "\u201c", "\"")
		line = strings.ReplaceAll(line, "\u201d", "\"")
		line = strings.ReplaceAll(line, "\u2018", "'")
		line = strings.ReplaceAll(line, "\u2019", "'")
		// Remove any other non-ascii unicode character
		line = asciiFilter.ReplaceAllLiteralString(line, "")
		// Trim any remaining whitespace
		line = strings.Trim(line, " ")

		trimmed[idx] = line
	}

	return strings.Trim(strings.Join(trimmed, " "), " ")
}

// PrintUnicodeRunes finds all non-ascii characters in the string and prints out the unicode character point.
// This is useful for debugging to find unicode characters that need to be handled by the CleanText function.
func PrintUnicodeRunes(content string) {
	if content == "" {
		return
	}
	for _, char := range asciiFilter.FindAllString(content, -1) {
		if char == "" {
			continue
		}
		r := []rune(char)
		fmt.Printf("%U\n", r)
	}
}
