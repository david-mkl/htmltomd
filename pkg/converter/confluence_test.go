package converter

import "testing"

func TestDefaultConfluenceConverter(t *testing.T) {
	doc := newTestDoc(`
<html>
	<head>
		<title>Test Doc</title>
	</head>
	<body>
		<div>
			<p>Ignored Content</p>
			<span id="title-text">Test Doc</span>
		</div>
		<div id="main-content">
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
			<div class="code">
				<div>
					<pre>Pre without highlight params</pre>
				</div>
			</div>
			<div class="code">
				<div>
					<pre data-syntaxhighlighter-params="brush: php; gutter: false; theme: Confluence">Pre with default highlight params</pre>
				</div>
			</div>
			<div class="code">
				<div>
					<pre data-syntaxhighlighter-params="brush: python; gutter: false; theme: Confluence">import math
print(math.pi)</pre>
				</div>
			</div>
			<div class="panel">
				<div class="confluence-information-macro-body">
					<p>Note Panel</p>
				</div>
			</div>
			<div class="confluence-information-macro confluence-information-macro-information">
				<div class="confluence-information-macro-body">
					<p>Info Panel</p>
				</div>
			</div>
			<div class="confluence-information-macro confluence-information-macro-note">
				<div class="confluence-information-macro-body">
					<p>Warning Panel</p>
				</div>
			</div>
			<div class="confluence-information-macro confluence-information-macro-tip">
				<div class="confluence-information-macro-body">
					<p>Tip Panel</p>
				</div>
			</div>
			<div class="confluence-information-macro confluence-information-macro-warning">
				<div class="confluence-information-macro-body">
					<p>Error Panel</p>
				</div>
			</div>
		</div>
	</body>
</html>
`)

	s := NewConfluenceSelectionConverter(SelectionConverterConfig{})
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
| Data 1 | Data 2 |

` + "```txt\nPre without highlight params\n```\n\n" + "```txt\nPre with default highlight params\n```\n\n" + "```python\nimport math\nprint(math.pi)\n```\n\n" + `Note Panel

Info Panel

Warning Panel

Tip Panel

Error Panel`

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
