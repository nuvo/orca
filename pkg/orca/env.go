package orca

import (
	"errors"
	"io"
	"os"

	"github.com/maorfr/orca/pkg/utils"

	"github.com/spf13/cobra"
)

type envCmd struct {
	chartsFile   string
	name         string
	override     []string
	packedValues []string
	set          []string
	kubeContext  string
	tls          bool
	helmTLSStore string
	repo         string
	createNS     bool
	onlyManaged  bool
	output       string
	inject       bool

	out io.Writer
}

// NewGetEnvCmd represents the get env command
func NewGetEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Get list of Helm releases in an environment (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.tls && e.helmTLSStore == "" {
				return errors.New("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			releases := utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, e.onlyManaged, true)

			switch e.output {
			case "yaml":
				utils.PrintReleasesYaml(releases)
			case "md":
				utils.PrintReleasesMarkdown(releases)
			case "":
				utils.PrintReleasesYaml(releases)
			}
		},
	}

	f := cmd.Flags()

	f.StringVar(&e.name, "name", os.Getenv("ORCA_NAME"), "name of environment (namespace) to get. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.onlyManaged, "only-managed", utils.GetBoolEnvVar("ORCA_ONLY_MANAGED", true), "list only releases managed by orca. Overrides $ORCA_ONLY_MANAGED")
	f.StringVarP(&e.output, "output", "o", os.Getenv("ORCA_OUTPUT"), "output format (yaml, md). Overrides $ORCA_OUTPUT")
	return cmd
}

// NewDeployEnvCmd represents the deploy env command
func NewDeployEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Deploy a list of Helm charts to an environment (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.tls && e.helmTLSStore == "" {
				return errors.New("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			if circular := utils.CheckCircularDependencies(utils.ChartsYamlToStruct(e.chartsFile, e.name)); circular {
				return errors.New("Circular dependency found")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if e.createNS {
				utils.CreateNamespace(e.name, e.kubeContext)
			}

			desiredReleases := utils.ChartsYamlToStruct(e.chartsFile, e.name)
			desiredReleases = utils.OverrideReleases(desiredReleases, e.override)
			installedReleases := utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, true, false)
			releasesToInstall := utils.GetReleasesDelta(desiredReleases, installedReleases)

			utils.AddRepository(e.repo, false)
			utils.UpdateRepositories(false)
			utils.DeployChartsFromRepository(releasesToInstall, e.kubeContext, e.name, e.repo, e.helmTLSStore, e.tls, e.packedValues, e.set, e.inject)

			installedReleases = utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, true, false)
			releasesToDelete := utils.GetReleasesDelta(installedReleases, desiredReleases)

			utils.DeleteReleases(releasesToDelete, e.kubeContext, e.helmTLSStore, e.tls)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.chartsFile, "charts-file", "c", os.Getenv("ORCA_CHARTS_FILE"), "path to file with list of Helm charts to install. Overrides $ORCA_CHARTS_FILE")
	f.StringSliceVar(&e.override, "override", []string{}, "chart to override with different version (can specify multiple): chart=version")
	f.StringVar(&e.name, "name", os.Getenv("ORCA_NAME"), "name of environment (namespace) to deploy to. Overrides $ORCA_NAME")
	f.StringVar(&e.repo, "repo", os.Getenv("ORCA_REPO"), "chart repository (name=url). Overrides $ORCA_REPO")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.StringSliceVarP(&e.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&e.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.createNS, "create-ns", utils.GetBoolEnvVar("ORCA_CREATE_NS", false), "should create new namespace. Overrides $ORCA_CREATE_NS")
	f.BoolVar(&e.inject, "inject", utils.GetBoolEnvVar("ORCA_INJECT", false), "enable injection during helm upgrade. Overrides $ORCA_INJECT")

	return cmd
}

// NewDeleteEnvCmd represents the delete env command
func NewDeleteEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Delete an environment (Kubernetes namespace) along with all Helm releases in it",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.tls && e.helmTLSStore == "" {
				return errors.New("TLS is set to true and HELM_TLS_STORE is not defined")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			releases := utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, true, true)
			utils.DeleteReleases(releases, e.kubeContext, e.helmTLSStore, e.tls)
			utils.DeleteNamespace(e.name, e.kubeContext)
		},
	}

	f := cmd.Flags()

	f.StringVar(&e.name, "name", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")

	return cmd
}
