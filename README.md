# html-to-md

CLI tool and library to Convert HTML to Markdown with support for Confluence, and Google Docs.

## Input Sources

In addition to arbitrary HTML, `html-to-md` can also handle HTML files that have been exported from Confluence and Google Docs. In these cases, `html-to-md` will search for specific known elements that can be converted into markdown.

For example, Confluence expresses code fences with HTML and CSS that have a known structure and CSS classes. `html-to-md` will search for these elements and convert them to markdown.

## Output Formats

`html-to-md` can output markdown in specific formats, such as for a [Hugo](https://gohugo.io/) website.

For example, an image in normal markdown is expressed as

```markdown
![Alt Text](https://source.png)
```

`html-to-md` can be configured to instead output image references as a [Hugo figure shortcode](https://gohugo.io/content-management/shortcodes/#figure) like

```go
{{< figure src="https://source.png" alt="Alt Text" >}}
```

## Default Conversions

`html-to-md` will search for the elements below and convert them to markdown format.

|  | From | To |
| --- | --- | --- |
| Links | `<a href="https://link">Link</a>` |  `[Link](https://link)` |
| Bold | `<strong>Bold Text</strong>` |  `**Bold Text**` |
| Italics | `<em>Italics</em>` |  `_Italics_` |
| Images | `<img src="https://source.png" alt="Alt Text" />` |  `![Alt Text](https://source.png)` |
| Code | `<code>Code</code>` |  `` ` ``Code`` ` `` |

### Preformatted Text

```html
<pre>
def func():
    print("Hello World")
</pre>
```

to

    ```
    def func():
        print("Hello World")
    ```

### Tables

```html
<table>
    <tr>
        <th>Header 1</th>
        <th>Header 2</th>
    </tr>
    <tr>
        <td>Data 1,1</td>
        <td>Data 1,2</td>
    </tr>
    <tr>
        <td>Data 2,1</td>
        <td>Data 2,2</td>
    </tr>
</table>
```

to

```txt
| Header 1 | Header 2 |
| --- | --- |
| Data 1,1 | Data 1,2 |
| Data 2,1 | Data 2,2 |
```

## Install

The binary can be built from source and requires go >= 1.17 to be installed on your system. (The `build` step assumes you have appropriate values for `GOOS` and `GOARCH` set for your system).

Build the binary with

```sh
make install
make build
```

This creates an executable called `html2md`.

## Usage

```sh
./html2md -h
```

### Convert Command

Converts `.html` files to `.md` files.

```sh
./html2md convert <file.html|directory>
```

The argument can be an HTML file or a directory containing HTML files.

#### Flags

An optional `--out` flag can be specified to indicate the directory where converted files should be placed (the directory will be created if it doesn't exist). If not specified, a directory called `html_to_md_converted` will be created for the converted files.

The input source can be specified with a `--input-format` flag to handle specific kinds of input HTML files. Supported values are

* `html` - Arbitrary HTML. This is the default value.
* `confluence` - Confluence Docs that have been converted to HTML
* `google` - Google Docs that have been converted to HTML

For example

```sh
./html2md convert --input-format confluence path/to/confluence/files
```

Specify the output format with a `--output-format` flag. Supported values are

* `md` - Renders markdown elements normally. This is the default value.
* `hugo` - Renders markdown elements as shortcodes for a Hugo website

```txt
./html2md convert --output-format hugo path/to/files
```

## Usage as a Library

You may also install the components of this tool to use in your own Go code for further customization.

```sh
go get github.com/david-mk-lawrence/html-to-md
```

Then import the converter package

```go
import "github.com/david-mk-lawrence/html-to-md/pkg/converter"
```

Two structs are needed to convert documents. A DocumentConverter and a struct that implements a SelectionConverter interface. A DocumentConverter is what handles the HTML document itself. A SelectionConverter is an interface that handles and converts specific elements in the document. This library provides a

* HTMLSelectionConverter
* GoogleSelectionConverter
* ConfluenceSelectionConverter

As an example, initialize a standard HTML converter with

```go
s := NewHTMLSelectionConverter(SelectionConverterConfig{})
c := DocumentConverter{SelectionConv: &s}
markdownContent := c.DocumentToMarkdown(doc).String()
```

Where `doc` is a `*goquery.Document` (see [goquery](https://github.com/PuerkitoBio/goquery) for more information).

### Customization

Since SelectionConverter is an interface, you may write your own implementation. The provided SelectionConverters are also customizable, so you may also just override specific hooks.

#### Customize Existing Converters

For example, the ConfluenceSelectionConverter will search for an element in the document with an ID of "title-text" in order to obtain the title of the document. You may override this behavior by providing a custom function to obtain the title element in a different location.

```go
conf := SelectionConverterConfig{
    TitleFinder: func(doc *goquery.Document) string {
        return doc.Find(".custom-title-location").First().Text()
    },
}
s := NewConfluenceSelectionConverter(conf)
c := DocumentConverter{SelectionConv: &s}
markdownContent := c.DocumentToMarkdown(doc).String()
```

The following hooks may be configured in the SelectionConverterConfig

| Signature | Description  |
| --- | ---  |
| FindRootElement(*goquery.Document) *goquery.Selection | This defines where the SelectionConverter will begin looking for content. For example, to crawl the entire document, `return doc.Find("html")` |
| FindTitle(*goquery.Document) string | This defines what the title of the final Markdown document will be. |
| FindContentElements(*goquery.Selection) *goquery.Selection | As the SelectionConverter crawls down the document from the root, it will only continue to crawl selections returned from this function. Generally this is a good way to filter on specific HTML tags. For example, if content is only in `p` and `span` tags, then `return s.ChildrenFiltered("p,span")` |
| HandleMatchedSelection(int, *goquery.Selection, *markdown.Doc, SelectionToMD) | This function will be called for every matched element returned by `FindContentElements`. This function is where content should be extracted from the element and added to the markdown document. `SelectionToMD` is a callable that enables the converted to recursively crawl through the document. It should be called on elements that have children. |

#### Create New Custom Converter

If the provided SelectionConverters do not handle your documents properly, or cannot effectively be overwritten, you can write your own entirely custom SelectionConverter. Simply implement the functions in the table above.

For example

```go
type CustomSelectionConverter struct {
}

func (c *CustomSelectionConverter) FindRootElement(doc *goquery.Document) *goquery.Selection {
    return doc.Find("body").First()
}

func (c *CustomSelectionConverter) FindTitle(doc *goquery.Document) string {
    return doc.Find("h1").First().Text()
}

func (c *CustomSelectionConverter) FindContentElements(s *goquery.Selection) *goquery.Selection {
    return s.Find("p")
}

func (c *CustomSelectionConverter) HandleMatchedSelection(i int, s *goquery.Selection, md *markdown.Doc, toMD SelectionToMD) {
    md.AddParagraph(s.Text())
}
```

## Exporting HTML

### Exporting Confluence Docs to HTML

Confluence only supports exporting entire spaces to HTML. To export a space, go to "Space Settings" and select "Export Space".

### Exporting Google Docs to HTML

With the Google Doc open, select File -> Download -> Web Page. This will download the HTML as a zip archive. Unzip the archive which will contain the HTML file and other resources like images.
