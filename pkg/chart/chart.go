package chart

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type chartCmd struct {
	nada string

	out io.Writer
}

// NewDeployCmd represents the deploy chart command
func NewDeployCmd(out io.Writer) *cobra.Command {
	s := &chartCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy chart called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}

// NewPushCmd represents the push chart command
func NewPushCmd(out io.Writer) *cobra.Command {
	s := &chartCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("push chart called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
