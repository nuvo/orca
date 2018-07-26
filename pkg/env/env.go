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
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
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
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
