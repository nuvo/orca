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
		Short: "Determine build type based on path filters",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("determine buildtype called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
