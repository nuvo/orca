package cmd

import (
	"fmt"
	"io"

	"github.com/maorfr/orca/pkg/chart"
	"github.com/maorfr/orca/pkg/env"
	"github.com/spf13/cobra"
)

// NewDeployCmd represents the get command
func NewDeployCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy called")
		},
	}

	cmd.AddCommand(chart.NewDeployCmd(out))
	cmd.AddCommand(env.NewDeployCmd(out))

	return cmd
}
