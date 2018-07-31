package cmd

import (
	"io"

	"github.com/maorfr/orca/pkg/determine"
	"github.com/spf13/cobra"
)

// NewDetermineCmd represents the get command
func NewDetermineCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "determine",
		Short: "Determination functions",
		Long:  ``,
	}

	cmd.AddCommand(determine.Buildtype(out))

	return cmd
}
