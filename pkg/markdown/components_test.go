package markdown

import (
	"testing"
)

func TestUnorderedListToString(t *testing.T) {
	items := []string{"item 1", "item 2"}
	list := NewUnorderedList(items)

	result := list.String()
	expected := "* item 1\n* item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestOrderedListToString(t *testing.T) {
	items := []string{"item 1", "item 2"}
	list := NewOrderedList(items)

	result := list.String()
	expected := "1. item 1\n1. item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH1HeaderToString(t *testing.T) {
	h := Header{headerType: h1, Content: "Title"}

	result := h.String()
	expected := "# Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH2HeaderToString(t *testing.T) {
	h := Header{headerType: h2, Content: "Title"}

	result := h.String()
	expected := "## Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH3HeaderToString(t *testing.T) {
	h := Header{headerType: h3, Content: "Title"}

	result := h.String()
	expected := "### Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH4HeaderToString(t *testing.T) {
	h := Header{headerType: h4, Content: "Title"}

	result := h.String()
	expected := "#### Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH5HeaderToString(t *testing.T) {
	h := Header{headerType: h5, Content: "Title"}

	result := h.String()
	expected := "##### Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestH6HeaderToString(t *testing.T) {
	h := Header{headerType: h6, Content: "Title"}

	result := h.String()
	expected := "###### Title"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestParagraphToString(t *testing.T) {
	p := Paragraph{Content: "content\ncan contain\nnewlines"}

	result := p.String()
	expected := "content\ncan contain\nnewlines"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestCodeBlockToString(t *testing.T) {
	cb := CodeBlock{Lang: "go", Code: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}"}

	result := cb.String()
	expected := "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}\n```"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestHorizontalRuleToString(t *testing.T) {
	hr := HorizontalRule{}

	result := hr.String()
	expected := "---"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestTableToStringWithHeaders(t *testing.T) {
	headers := []string{"Column 1", "Column 2"}
	rows := [][]string{
		{"data 1,1", "data 1,2"},
		{"data 2,1", "data 2,2"},
	}
	table := Table{Headers: headers, Rows: rows}

	result := table.String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| data 1,1 | data 1,2 |\n| data 2,1 | data 2,2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}
