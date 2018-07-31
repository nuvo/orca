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
		Short: "Deploy a Helm chart from chart museum",
		Long:  ``,
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
		Short: "Push Helm chart to chart museum",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("push chart called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
