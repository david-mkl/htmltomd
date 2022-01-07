package htmltomd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "htmltomd",
		Short: "converts HTML to Markdown",
		Long:  `Reads a local HTML file and converts it into a Markdown document`,
	}
	quiet   bool = false
	verbose bool = false
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "make output more verbose")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func out(format string, a ...interface{}) {
	if !quiet {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

func outV(format string, a ...interface{}) {
	if verbose {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}
