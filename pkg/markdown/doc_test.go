package markdown

import (
	"testing"

	"github.com/david-mk-lawrence/htmltomd/pkg/util"
)

type testContent struct {
}

func (tc testContent) String() string {
	return "test"
}

func TestGetConfigWithDefaults(t *testing.T) {
	doc := NewDoc(DocConfig{})

	conf := doc.GetConfig()

	if conf.Title != nil {
		t.Errorf("Expected nil")
	}
	if *conf.ReduceHeaders != true {
		t.Errorf("Expected %t. Got %t.", true, *conf.ReduceHeaders)
	}
	if *conf.Separator != "\n\n" {
		t.Errorf("Expected %s. Got %s.", "\\n\\n", *conf.Separator)
	}
}

func TestGetConfigWithSetValue(t *testing.T) {
	doc := NewDoc(DocConfig{Title: util.String("title"), ReduceHeaders: util.Bool(false), Separator: util.String("-")})

	conf := doc.GetConfig()

	if *conf.Title != "title" {
		t.Errorf("Expected %s. Got %s.", "\"\"", *conf.Title)
	}
	if *conf.ReduceHeaders != false {
		t.Errorf("Expected %t. Got %t.", true, *conf.ReduceHeaders)
	}
	if *conf.Separator != "-" {
		t.Errorf("Expected %s. Got %s.", "\\n\\n", *conf.Separator)
	}
}

func TestGetRenderConfig(t *testing.T) {
	doc := NewDoc(DocConfig{Title: util.String("title"), ReduceHeaders: util.Bool(false), Separator: util.String("-")})

	conf := doc.GetRenderConfig()

	if conf.Title != nil {
		t.Error("Expected nil")
	}
	if *conf.ReduceHeaders != false {
		t.Errorf("Expected %t. Got %t.", true, *conf.ReduceHeaders)
	}
	if *conf.Separator != "-" {
		t.Errorf("Expected %s. Got %s.", "\\n\\n", *conf.Separator)
	}
}

func TestAddDoc(t *testing.T) {
	subdoc := NewDoc(DocConfig{})
	subdoc.AddContent(testContent{})

	doc := NewDoc(DocConfig{})
	doc.AddDoc(subdoc)

	result := doc.String()
	expected := "test"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddHeaderWithReducedHeaders(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddHeader("h1", "Header")

	result := doc.String()
	expected := "## Header"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddHeaderWithUnreducedHeaders(t *testing.T) {
	doc := NewDoc(DocConfig{ReduceHeaders: util.Bool(false)})
	doc.AddHeader("h1", "Header")

	result := doc.String()
	expected := "# Header"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddUnorderedList(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddUnorderedList([]string{"item 1", "item 2"})

	result := doc.String()
	expected := "* item 1\n* item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddOrderedList(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddOrderedList([]string{"item 1", "item 2"})

	result := doc.String()
	expected := "1. item 1\n1. item 2"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddParagraph(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddParagraph("Paragraph")

	result := doc.String()
	expected := "Paragraph"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddEmptyParagraph(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddParagraph("Paragraph")
	doc.AddParagraph("")

	result := doc.String()
	expected := "Paragraph" // Should not see a newline. Empty content should have been ignored.

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddCodeBlock(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddCodeBlock("go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}")

	result := doc.String()
	expected := "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}\n```"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddHorizontalRule(t *testing.T) {
	doc := NewDoc(DocConfig{})
	doc.AddHorizontalRule()

	result := doc.String()
	expected := "---"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestAddTable(t *testing.T) {
	doc := NewDoc(DocConfig{})

	headers := []string{"Column 1", "Column 2"}
	rows := [][]string{
		{"data 1,1", "data 1,2"},
		{"data 2,1", "data 2,2"},
	}
	doc.AddTable(headers, rows)

	result := doc.String()
	expected := "| Column 1 | Column 2 |\n| --- | --- |\n| data 1,1 | data 1,2 |\n| data 2,1 | data 2,2 |"

	if result != expected {
		t.Errorf("Expected\n%s\nGot\n%s", expected, result)
	}
}

func TestTitle(t *testing.T) {
	title := "Test Doc"
	doc := NewDoc(DocConfig{Title: &title})

	result := doc.Title()
	expected := "# Test Doc"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestTitleNotSet(t *testing.T) {
	doc := NewDoc(DocConfig{})

	result := doc.Title()
	expected := ""

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}

func TestContentWithTitle(t *testing.T) {
	title := "Test Doc"
	doc := NewDoc(DocConfig{Title: &title})
	doc.AddParagraph("Paragraph")

	result := doc.String()
	expected := "# Test Doc\n\nParagraph"

	if result != expected {
		t.Errorf("Expected %s. Got %s", expected, result)
	}
}
