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
	repo         string

	out io.Writer
}

// NewDeployChartCmd represents the deploy chart command
func NewDeployChartCmd(out io.Writer) *cobra.Command {
	c := &chartCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Deploy a Helm chart from chart repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if c.tls && c.helmTLSStore == "" {
				log.Fatal("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			utils.DeployChartFromRepository(c.releaseName, c.name, c.version, c.kubeContext, c.namespace, c.repo, c.helmTLSStore, c.tls, c.packedValues, c.set, true)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.name, "name", "", "name of chart to deploy")
	f.StringVar(&c.version, "version", "", "version of chart to deploy")
	f.StringVar(&c.repo, "repo", "", "chart repository (name=url)")
	f.StringVar(&c.releaseName, "release-name", "", "version of chart to deploy")
	f.StringVar(&c.kubeContext, "kube-context", "", "name of the kubeconfig context to use")
	f.StringVarP(&c.namespace, "namespace", "n", "", "kubernetes namespace to deploy to")
	f.StringSliceVarP(&c.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&c.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&c.tls, "tls", true, "enable TLS for request")
	f.StringVar(&c.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")

	return cmd
}

type chartPushCmd struct {
	path   string
	append string
	repo   string
	lint   bool

	out io.Writer
}

// NewPushChartCmd represents the push chart command
func NewPushChartCmd(out io.Writer) *cobra.Command {
	c := &chartPushCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Push Helm chart to chart repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PushChartToRepository(c.path, c.append, c.repo, c.lint, false)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.path, "path", "", "path to chart")
	f.StringVar(&c.append, "append", "", "string to append to version")
	f.StringVar(&c.repo, "repo", "", "chart repository (name=url)")
	f.BoolVar(&c.lint, "lint", false, "should perform lint before push")

	return cmd
}
