package orca

import (
	"errors"
	"io"
	"os"

	"github.com/maorfr/orca/pkg/utils"

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
	inject       bool

	out io.Writer
}

// NewDeployChartCmd represents the deploy chart command
func NewDeployChartCmd(out io.Writer) *cobra.Command {
	c := &chartCmd{out: out}

	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Deploy a Helm chart from chart repository",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if c.tls && c.helmTLSStore == "" {
				return errors.New("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			if c.repo == "" {
				return errors.New("repo can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			utils.DeployChartFromRepository(c.releaseName, c.name, c.version, c.kubeContext, c.namespace, c.repo, c.helmTLSStore, c.tls, c.packedValues, c.set, true, c.inject)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.name, "name", os.Getenv("ORCA_NAME"), "name of chart to deploy. Overrides $ORCA_NAME")
	f.StringVar(&c.version, "version", os.Getenv("ORCA_VERSION"), "version of chart to deploy. Overrides $ORCA_VERSION")
	f.StringVar(&c.repo, "repo", os.Getenv("ORCA_REPO"), "chart repository (name=url). Overrides $ORCA_REPO")
	f.StringVar(&c.releaseName, "release-name", os.Getenv("ORCA_RELEASE_NAME"), "release name. Overrides $ORCA_RELEASE_NAME")
	f.StringVar(&c.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.StringVarP(&c.namespace, "namespace", "n", os.Getenv("ORCA_NAMESPACE"), "kubernetes namespace to deploy to. Overrides $ORCA_NAMESPACE")
	f.StringSliceVarP(&c.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&c.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&c.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&c.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&c.inject, "inject", utils.GetBoolEnvVar("ORCA_INJECT", false), "enable injection during helm upgrade. Overrides $ORCA_INJECT")

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
		Args: func(cmd *cobra.Command, args []string) error {
			if c.repo == "" {
				return errors.New("repo can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			utils.PushChartToRepository(c.path, c.append, c.repo, c.lint, false)
		},
	}

	f := cmd.Flags()

	f.StringVar(&c.path, "path", os.Getenv("HELM_PATH"), "path to chart. Overrides $ORCA_PATH")
	f.StringVar(&c.append, "append", os.Getenv("HELM_APPEND"), "string to append to version. Overrides $ORCA_APPEND")
	f.StringVar(&c.repo, "repo", os.Getenv("HELM_REPO"), "chart repository (name=url). Overrides $ORCA_REPO")
	f.BoolVar(&c.lint, "lint", utils.GetBoolEnvVar("ORCA_LINT", false), "should perform lint before push. Overrides $ORCA_LINT")

	return cmd
}
