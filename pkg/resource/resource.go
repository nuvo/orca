package resource

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type resourceCmd struct {
	url string

	out io.Writer
}

// NewGetCmd represents the get resource command
func NewGetCmd(out io.Writer) *cobra.Command {
	g := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get resource called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&g.url, "url", "", "url help")

	return cmd
}

// NewDeleteCmd represents the delete resource command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	g := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete resource called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&g.url, "url", "", "url help")

	return cmd
}
