package orca

import (
	"io"
	"log"
	"os"

	"orca/pkg/utils"

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

// NewDeployChartCmd represents the deploy chart command
func NewDeployChartCmd(out io.Writer) *cobra.Command {
	c := &chartCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Deploy a Helm chart from chart museum",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if c.tls && c.helmTLSStore == "" {
				log.Fatal("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			utils.DeployChartFromMuseum(c.releaseName, c.name, c.version, c.kubeContext, c.namespace, c.museum, c.helmTLSStore, c.tls, c.packedValues, c.set, true)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.name, "name", "", "name of chart to deploy")
	f.StringVar(&c.version, "version", "", "version of chart to deploy")
	f.StringVar(&c.museum, "museum", "", "chart museum instance (name=url)")
	f.StringVar(&c.releaseName, "release-name", "", "version of chart to deploy")
	f.StringVar(&c.kubeContext, "kube-context", "", "kubernetes context to deploy to")
	f.StringVarP(&c.namespace, "namespace", "n", "", "kubernetes namespace to deploy to")
	f.StringSliceVarP(&c.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&c.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&c.tls, "tls", true, "should use communication over TLS")
	f.StringVar(&c.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "directory with TLS certs and keys")

	return cmd
}

type chartPushCmd struct {
	path   string
	append string
	museum string
	lint   bool

	out io.Writer
}

// NewPushChartCmd represents the push chart command
func NewPushChartCmd(out io.Writer) *cobra.Command {
	c := &chartPushCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Push Helm chart to chart museum",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PushChartToMuseum(c.path, c.append, c.museum, c.lint, false)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.path, "path", "", "path to chart")
	f.StringVar(&c.append, "append", "", "string to append to version")
	f.StringVar(&c.museum, "museum", "", "chart museum instance (name=url)")
	f.BoolVar(&c.lint, "lint", false, "should perform lint before push")

	return cmd
}
