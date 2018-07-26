package cmd

import (
	"fmt"
	"io"

	"github.com/maorfr/orca/pkg/env"
	"github.com/maorfr/orca/pkg/resource"
	"github.com/spf13/cobra"
)

// NewDeleteCmd represents the get command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	}

	cmd.AddCommand(env.NewDeleteCmd(out))
	cmd.AddCommand(resource.NewDeleteCmd(out))

	return cmd
}
