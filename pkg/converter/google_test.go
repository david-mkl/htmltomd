package converter

import "testing"

func TestDefaultGoogleConverter(t *testing.T) {
	doc := newTestDoc(`
<html>
	<head>
		<title>Test Doc</title>
	</head>
	<body>
		<h1>Section Title</h1>
		<span>Subtitle</span>
		<p>Test Paragraph 1</p>
		<div>
			<p>Test Paragraph 2</p>
			<p>Test Paragraph 3</p>
		</div>
		<hr>
		<div>
			<ul>
				<li>Item 1</li>
				<li>Item 2</li>
			</ul>
			<ol>
				<li>Item 1</li>
				<li>Item 2</li>
			</ol>
		</div>
		<hr style="page-break-after: auto">
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

	s := NewGoogleSelectionConverter(SelectionConverterConfig{})
	c := DocumentConverter{SelectionConv: s}

	result := c.DocumentToMarkdown(doc).String()
	expected := `# Test Doc

## Section Title

Subtitle

Test Paragraph 1

Test Paragraph 2

Test Paragraph 3

---

* Item 1
* Item 2

1. Item 1
1. Item 2

| Column 1 | Column 2 |
| --- | --- |
| Data 1 | Data 2 |`

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
