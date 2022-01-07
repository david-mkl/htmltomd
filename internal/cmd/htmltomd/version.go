package htmltomd

import (
	"github.com/david-mk-lawrence/html-to-md/internal/version"
	"github.com/spf13/cobra"
)

type versionCmd struct {
}

func init() {
	c := versionCmd{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get the current version",
		Run:   c.version,
		Args:  cobra.NoArgs,
	}

	rootCmd.AddCommand(cmd)
}

func (c *versionCmd) version(cmd *cobra.Command, args []string) {
	out("%s", version.GetVersion())
}
