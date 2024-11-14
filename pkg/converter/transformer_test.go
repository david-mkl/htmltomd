package converter

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func newTestDoc(content string) *goquery.Document {
	r := strings.NewReader(content)
	doc, _ := goquery.NewDocumentFromReader(r)
	return doc
}

func clean(str string) string {
	return strings.Trim(str, " \n\t")
}

func deepClean(str string) string {
	return clean(strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\t", ""))
}

func TestRemoveScripts(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<head>
		<link rel="stylesheet" href="styles.css" type="text/css" />
		<style type="text/css">
		html {
			font-size: 14px;
		}
		</style>
		<script type="text/javascript" src="script.js"></script>
	</head>
	<body>
		<h1>Hello!</h1>
		<script type="text/javascript">
		var i = 0;
		</script>
		<script type="text/javascript" src="script.js"></script>
	</body>
</html>
`)

	html := doc.Find("html").First()
	tr.RemoveScripts(html)

	// head should be empty
	if len(html.Find("head").First().Children().Nodes) != 0 {
		html.Find("head").First().Children().Each(func(i int, s *goquery.Selection) {
			t.Errorf("Unexpected node in head: %s", s.Nodes[0].Data)
		})
	}

	// body should contain no scripts
	if len(html.Find("body").ChildrenFiltered("script").Nodes) != 0 {
		html.Find("head").First().Children().Each(func(i int, s *goquery.Selection) {
			t.Errorf("Unexpected script node in body: %s", s.Nodes[0].Data)
		})
	}

	result := clean(html.Text())
	expected := "Hello!"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestTransform(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<h1>Hello World!</h1>
	</body>
</html>
`)

	capitalize := func(i int, s *goquery.Selection) {
		s.SetText(strings.ToUpper(s.Text()))
	}
	hyphenize := func(i int, s *goquery.Selection) {
		s.SetText(strings.ReplaceAll(s.Text(), " ", "-"))
	}

	tr.Transform("h1", doc.Find("body"), hyphenize, capitalize)

	result := clean(doc.Find("body").Text())
	expected := "HELLO-WORLD!"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestToUnorderedList(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<ul>
			<li>unordered item 1</li>
			<li>unordered item 2</li>
		</ul>
	</body>
</html>
`)
	result := tr.ToList(doc.Find("ul")).String()
	expected := "* unordered item 1\n* unordered item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestToOrderedList(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<ol>
			<li>ordered item 1</li>
			<li>ordered item 2</li>
		</ol>
	</body>
</html>
`)

	result := tr.ToList(doc.Find("ol")).String()
	expected := "1. ordered item 1\n1. ordered item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestToTableWithoutHeadAndBody(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<tr>
				<th>Column 1</th>
				<th>Column 2</th>
			</tr>
			<tr>
				<td>Data 1</td>
				<td>Data 2</td>
			</tr>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| Data 1 | Data 2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestToTableWithoutHeaders(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<tr>
				<td>Data 1,1</td>
				<td>Data 1,2</td>
			</tr>
			<tr>
				<td>Data 2,1</td>
				<td>Data 2,2</td>
			</tr>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	expected := "| Data 1,1 | Data 1,2 |\n| Data 2,1 | Data 2,2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestToTableWithHeadAndBody(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<thead>
				<tr>
					<th>Column 1</th>
					<th>Column 2</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>Data 1</td>
					<td>Data 2</td>
				</tr>
			</tbody>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| Data 1 | Data 2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
func TestToTableWithHeadWithoutRow(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<thead>
				<th>Column 1</th>
				<th>Column 2</th>
			</thead>
			<tr>
				<td>Data 1</td>
				<td>Data 2</td>
			</tr>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| Data 1 | Data 2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestToTableWithTdAsHeadersInThead(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<thead>
				<td>Column 1</td>
				<td>Column 2</td>
			</thead>
			<tr>
				<td>Data 1</td>
				<td>Data 2</td>
			</tr>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| Data 1 | Data 2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestToTableWithMultipleHeaderRows(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<table>
			<thead>
				<tr>
					<th>Column 1</th>
					<th>Column 2</th>
				</tr>
				<tr>
					<th>Column 3</th>
					<th>Column 4</th>
				</tr>
			</thead>
			<tr>
				<td>Data 1</td>
				<td>Data 2</td>
			</tr>
		</table>
	</body>
</html>
`)

	result := tr.ToTable(doc.Find("table")).String()
	// The second row in thead will be ignored, as markdown only supports one header row
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| Data 1 | Data 2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestReplaceAll(t *testing.T) {
	tr := NewTransformer(nil)
	doc := newTestDoc(`
<html>
	<body>
		<h1>This is a <a href="mock://example.com">Link</a></h1>
		<img src="mock://example.com" alt="Test Image" />
		<p>
			This is a <strong>bold statement</strong>. This is an <em>italicized statement</em>. This is <code>inline code</code>.
		</p>
	</body>
</html>
`)

	tr.ReplaceAll(doc.Find("body"))

	result := deepClean(doc.Text())
	expected := "This is a [Link](mock://example.com)![Test Image](mock://example.com)This is a **bold statement**. This is an _italicized statement_. This is `inline code`."

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestDefaultTextCleaner(t *testing.T) {
	tc := NewTextCleaner(nil)

	dirty := "  \u00b6\u2018Hello\u2019\u00a0\u201cWorld\u201d! "

	result := tc.CleanText(dirty)
	expected := "'Hello' \"World\"!"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAsciiOnlyTextCleaner(t *testing.T) {
	tc := NewTextCleaner(&TextCleanerConf{
		AsciiOnly: true,
	})

	dirty := "  \u00b6\u2018Hello\u2019\u00a0\u201cWorld\u201d! "

	result := tc.CleanText(dirty)
	expected := "'Hello' \"World\"!"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestUnicodeTextCleaner(t *testing.T) {
	tc := NewTextCleaner(&TextCleanerConf{
		AsciiOnly: false,
	})

	dirty := "  \u00b6\u2018Ħëlľō\u2019\u00a0\u201cŴórłď\u201d! "

	result := tc.CleanText(dirty)
	expected := "'Ħëlľō' \"Ŵórłď\"!"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}
