package cmd

import (
	"io"
	"orca/pkg/orca"

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

// NewDeleteCmd represents the get command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletion functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewDeleteEnvCmd(out))
	cmd.AddCommand(orca.NewDeleteResourceCmd(out))

	return cmd
}

// NewDeployCmd represents the get command
func NewDeployCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deployment functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewDeployChartCmd(out))
	cmd.AddCommand(orca.NewDeployEnvCmd(out))

	return cmd
}

// NewDetermineCmd represents the get command
func NewDetermineCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "determine",
		Short: "Determination functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewDetermineBuildtype(out))

	return cmd
}

// NewGetCmd represents the get command
func NewGetCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewGetEnvCmd(out))
	cmd.AddCommand(orca.NewGetResourceCmd(out))

	return cmd
}

// NewPushCmd represents the get command
func NewPushCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewPushChartCmd(out))

	return cmd
}
