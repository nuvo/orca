package cmd

import (
	"io"

	"github.com/maorfr/orca/pkg/env"
	"github.com/maorfr/orca/pkg/resource"
	"github.com/spf13/cobra"
)

// NewDeleteCmd represents the get command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletion functions",
		Long:  ``,
	}

	cmd.AddCommand(env.NewDeleteCmd(out))
	cmd.AddCommand(resource.NewDeleteCmd(out))

	return cmd
}
