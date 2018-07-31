package env

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type envCmd struct {
	nada string

	out io.Writer
}

// NewGetCmd represents the get env command
func NewGetCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Get list of Helm releases in an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}

// NewDeployCmd represents the deploy env command
func NewDeployCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Deploy a list of Helm charts to an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}

// NewDeleteCmd represents the delete env command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Delete an environment (Kubernetes namespace) along with all Helm releases in it",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
