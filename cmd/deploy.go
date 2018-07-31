package cmd

import (
	"io"

	"github.com/maorfr/orca/pkg/chart"
	"github.com/maorfr/orca/pkg/env"
	"github.com/spf13/cobra"
)

// NewDeployCmd represents the get command
func NewDeployCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deployment functions",
		Long:  ``,
	}

	cmd.AddCommand(chart.NewDeployCmd(out))
	cmd.AddCommand(env.NewDeployCmd(out))

	return cmd
}
