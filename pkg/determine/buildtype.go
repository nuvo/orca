package determine

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type determineCmd struct {
	nada string

	out io.Writer
}

// Buildtype represents the determine buildtype command
func Buildtype(out io.Writer) *cobra.Command {
	s := &determineCmd{out: out}

	cmd := &cobra.Command{
		Use:   "buildtype",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("determine buildtype called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
