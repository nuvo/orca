package cmd

import (
	"io"

	"orca/pkg/chart"

	"github.com/spf13/cobra"
)

// NewPushCmd represents the get command
func NewPushCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push functions",
		Long:  ``,
	}

	cmd.AddCommand(chart.NewPushCmd(out))

	return cmd
}
