package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orca",
		Short: "CI\\CD simplifier",
		Long: `Orca is a CI\CD simplifier, the glue behind the process.
Instead of writing scripts on top of scripts, Orca holds all the logic.
Use it wisely...`,
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(NewDeleteCmd(out))
	cmd.AddCommand(NewDeployCmd(out))
	cmd.AddCommand(NewDetermineCmd(out))
	cmd.AddCommand(NewGetCmd(out))
	cmd.AddCommand(NewPushCmd(out))

	return cmd
}
