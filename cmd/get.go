package cmd

import (
	"io"

	"orca/pkg/env"
	"orca/pkg/resource"

	"github.com/spf13/cobra"
)

// NewGetCmd represents the get command
func NewGetCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get functions",
		Long:  ``,
	}

	cmd.AddCommand(env.NewGetCmd(out))
	cmd.AddCommand(resource.NewGetCmd(out))

	return cmd
}
