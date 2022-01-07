package htmltomd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/david-mk-lawrence/html-to-md/pkg/converter"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

type convertCmd struct {
	outputDir    string
	outputFormat string
	inputFormat  string
}

func init() {
	c := convertCmd{}

	cmd := &cobra.Command{
		Use:   "convert [input.html|input_directory]",
		Short: "convert HTML file(s) to markdown",
		Long: `Input may be specified as either a directory or file. If a directory is given,
then all ".html" files in the directory will be converted.

If an output directory is specified, then the converted markdown files will be placed
there. Otherwise, a directory will be created called "html_to_md_converted".`,
		RunE: c.convert,
		Args: cobra.ExactArgs(1),
	}
	cmd.PersistentFlags().StringVar(&c.inputFormat, "input-format", "html", "source of html file. Can be 'html', 'confluence', or 'google'.")
	cmd.PersistentFlags().StringVar(&c.outputFormat, "output-format", "md", "style of markdown output. Can be 'hugo' or 'md'.")
	cmd.PersistentFlags().StringVarP(&c.outputDir, "out", "o", "./html_to_md_converted", "output directory")

	rootCmd.AddCommand(cmd)
}

func (c *convertCmd) convert(cmd *cobra.Command, args []string) (err error) {
	htmlPath, err := filepath.Abs(args[0])
	if err != nil {
		return
	}
	outV("Reading from %s", htmlPath)

	htmlFiles, err := c.getInputFiles(htmlPath)
	if err != nil {
		return
	}
	outV("Found %d html files", len(htmlFiles))

	if err = os.MkdirAll(c.outputDir, 0755); err != nil {
		return
	}
	outV("Placing markdown files in %s", c.outputDir)

	var wg sync.WaitGroup
	wgDoneChan := make(chan bool)
	errChan := make(chan error)

	conf := converter.SelectionConverterConfig{Transformer: &converter.Transformer{Format: c.outputFormat}}
	var selConv converter.SelectionConverter
	if c.inputFormat == "confluence" {
		selConv = converter.NewConfluenceSelectionConverter(conf)
	} else if c.inputFormat == "google" {
		selConv = converter.NewGoogleSelectionConverter(conf)
	} else {
		selConv = converter.NewHTMLSelectionConverter(conf)
	}
	conv := converter.DocumentConverter{SelectionConv: selConv}

	for _, htmlFile := range htmlFiles {
		wg.Add(1)
		go c.convertFile(conv, htmlFile, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(wgDoneChan)
		close(errChan)
	}()

	select {
	case <-wgDoneChan:
		break
	case convErr := <-errChan:
		err = multierror.Append(err, convErr)
	}

	return
}

func (c *convertCmd) getInputFiles(htmlPath string) (htmlFiles []string, err error) {
	info, err := os.Stat(htmlPath)
	if err != nil {
		return
	}

	if info.IsDir() {
		searchDir := filepath.Join(htmlPath, "*.html")
		outV("Input directory search path: %s", searchDir)
		htmlFiles, err = filepath.Glob(searchDir)
		return
	}

	if ext := filepath.Ext(htmlPath); ext != ".html" {
		err = fmt.Errorf("only html files can be used as input. Got %s", ext)
		return
	}

	htmlFiles = []string{htmlPath}
	return
}

func (c *convertCmd) getOutputFile(htmlFile string) string {
	return filepath.Join(c.outputDir, strings.TrimSuffix(filepath.Base(htmlFile), filepath.Ext(htmlFile))+".md")
}

func (c *convertCmd) convertFile(conv converter.DocumentConverter, htmlPath string, wg *sync.WaitGroup, errs chan<- error) {
	defer wg.Done()

	outFile := c.getOutputFile(htmlPath)
	out("Converting %s to %s", htmlPath, outFile)

	f, err := os.Open(htmlPath)
	if err != nil {
		errs <- err
		return
	}

	htmlDoc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		errs <- err
		return
	}

	mdContent := conv.DocumentToMarkdown(htmlDoc).String() + "\n"

	err = ioutil.WriteFile(outFile, []byte(mdContent), 0755)
	if err != nil {
		errs <- err
		return
	}
}
