package chart

import (
	"fmt"
	"io"
	"orca/pkg/helmflow"

	"github.com/spf13/cobra"
)

type chartCmd struct {
	name         string
	version      string
	releaseName  string
	packedValues []string
	set          []string
	kubeContext  string
	namespace    string
	tls          bool
	helmTLSStore string
	museum       string

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
			helmflow.DeployChartFromMuseum(s.releaseName, s.name, s.version, s.kubeContext, s.namespace, s.museum, s.helmTLSStore, s.tls, s.packedValues, s.set)
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.name, "name", "", "name of chart to deploy")
	f.StringVar(&s.version, "version", "", "version of chart to deploy")
	f.StringVar(&s.museum, "museum", "", "chart museum instance (name=url)")
	f.StringVar(&s.releaseName, "release-name", "", "version of chart to deploy")
	f.StringVar(&s.kubeContext, "kube-context", "", "kubernetes context to deploy to")
	f.StringVarP(&s.namespace, "namespace", "n", "", "kubernetes namespace to deploy to")
	f.StringSliceVarP(&s.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&s.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&s.tls, "tls", false, "should use communication over TLS")
	f.StringVar(&s.helmTLSStore, "helm-tls-store", "", "directory with TLS certs and keys")

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

	f.StringVar(&s.name, "name", "", "name help")

	return cmd
}
