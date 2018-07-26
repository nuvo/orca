package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orca",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(NewDeleteCmd(out))
	cmd.AddCommand(NewDeployCmd(out))
	cmd.AddCommand(NewDetermineCmd(out))
	cmd.AddCommand(NewGetCmd(out))
	cmd.AddCommand(NewPushCmd(out))

	return cmd
}
