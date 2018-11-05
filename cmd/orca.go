package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/maorfr/orca/pkg/orca"

	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		log.Fatal("Failed to execute command")
	}
}

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orca",
		Short: "CI\\CD simplifier",
		Long: `Orca is a CI\CD simplifier, the glue behind the process.
Instead of writing scripts on top of scripts, Orca holds all the logic.
`,
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(NewDeleteCmd(out))
	cmd.AddCommand(NewDeployCmd(out))
	cmd.AddCommand(NewDetermineCmd(out))
	cmd.AddCommand(NewGetCmd(out))
	cmd.AddCommand(NewPushCmd(out))
	cmd.AddCommand(NewCreateCmd(out))
	cmd.AddCommand(NewVersionCmd(out))
	cmd.AddCommand(NewLockCmd(out))
	cmd.AddCommand(NewUnlockCmd(out))
	cmd.AddCommand(NewDiffCmd(out))

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
	cmd.AddCommand(orca.NewDeployArtifactCmd(out))

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

// NewLockCmd represents the lock command
func NewLockCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Lock functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewLockEnvCmd(out))

	return cmd
}

// NewUnlockCmd represents the unlock command
func NewUnlockCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Unlock functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewUnlockEnvCmd(out))

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

// NewCreateCmd represents the create command
func NewCreateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creation functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewCreateResourceCmd(out))

	return cmd
}

// NewDiffCmd represents the create command
func NewDiffCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Differentiation functions",
		Long:  ``,
	}

	cmd.AddCommand(orca.NewDiffEnvCmd(out))

	return cmd
}

var GitTag, GitCommit string

// NewVersionCmd prints version information
func NewVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version %s (git-%s)\n", GitTag, GitCommit)
		},
	}

	return cmd
}
