package orca

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/maorfr/orca/pkg/utils"

	"github.com/spf13/cobra"
)

type envCmd struct {
	chartsFile                    string
	name                          string
	override                      []string
	packedValues                  []string
	set                           []string
	kubeContext                   string
	tls                           bool
	helmTLSStore                  string
	repo                          string
	createNS                      bool
	onlyManaged                   bool
	output                        string
	inject                        bool
	force                         bool
	deployOnlyOverrideIfEnvExists bool

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
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			releases := utils.GetInstalledReleases(e.kubeContext, e.name, false)

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

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to get. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.StringVarP(&e.output, "output", "o", os.Getenv("ORCA_OUTPUT"), "output format (yaml, md). Overrides $ORCA_OUTPUT")

	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.onlyManaged, "only-managed", utils.GetBoolEnvVar("ORCA_ONLY_MANAGED", true), "list only releases managed by orca. Overrides $ORCA_ONLY_MANAGED")

	f.MarkDeprecated("tls", "this command is no longer executed using helm")
	f.MarkDeprecated("helm-tls-store", "this command is no longer executed using helm")
	f.MarkDeprecated("only-managed", "environment is considered managed in any case")
	return cmd
}

// NewDeployEnvCmd represents the deploy env command
func NewDeployEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"environment"},
		Short:   "Deploy a list of Helm charts to an environment (Kubernetes namespace)",
		Long:    ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			if e.repo == "" {
				return errors.New("repo can not be empty")
			}
			if e.tls && e.helmTLSStore == "" {
				return errors.New("tls is set to true and helm-tls-store is not defined")
			}
			if e.chartsFile == "" && len(e.override) == 0 {
				return errors.New("either charts-file or override has to be defined")
			}
			if len(e.override) == 0 && e.deployOnlyOverrideIfEnvExists {
				return errors.New("override has to be defined when using deploy-only-override-if-env-exists")
			}
			if circular := utils.CheckCircularDependencies(utils.InitReleasesFromChartsFile(e.chartsFile, e.name)); circular {
				return errors.New("Circular dependency found")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			print := false

			utils.AddRepository(e.repo, print)
			utils.UpdateRepositories(print)

			nsPreExists := true
			if !utils.NamespaceExists(e.name, e.kubeContext) {
				nsPreExists = false
				utils.CreateNamespace(e.name, e.kubeContext, print)
				log.Printf("created environment \"%s\"", e.name)
			}
			lockEnvironment(e.name, e.kubeContext, print)

			var desiredReleases []utils.ReleaseSpec
			if nsPreExists && e.deployOnlyOverrideIfEnvExists {
				desiredReleases = utils.InitReleases(e.name, e.override)
			} else {
				desiredReleases = utils.InitReleasesFromChartsFile(e.chartsFile, e.name)
				desiredReleases = utils.OverrideReleases(desiredReleases, e.override, e.name)
			}

			includeFailed := false
			installedReleases := utils.GetInstalledReleases(e.kubeContext, e.name, includeFailed)
			releasesToInstall := utils.GetReleasesDelta(desiredReleases, installedReleases)

			utils.DeployChartsFromRepository(releasesToInstall, e.kubeContext, e.name, e.repo, e.helmTLSStore, e.tls, e.packedValues, e.set, e.inject)

			if !e.deployOnlyOverrideIfEnvExists {
				installedReleases = utils.GetInstalledReleases(e.kubeContext, e.name, includeFailed)
				releasesToDelete := utils.GetReleasesDelta(installedReleases, desiredReleases)
				utils.DeleteReleases(releasesToDelete, e.kubeContext, e.helmTLSStore, e.tls)
			}
			unlockEnvironment(e.name, e.kubeContext, print)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.chartsFile, "charts-file", "c", os.Getenv("ORCA_CHARTS_FILE"), "path to file with list of Helm charts to install. Overrides $ORCA_CHARTS_FILE")
	f.StringSliceVar(&e.override, "override", []string{}, "chart to override with different version (can specify multiple): chart=version")
	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to deploy to. Overrides $ORCA_NAME")
	f.StringVar(&e.repo, "repo", os.Getenv("ORCA_REPO"), "chart repository (name=url). Overrides $ORCA_REPO")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.StringSliceVarP(&e.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&e.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.inject, "inject", utils.GetBoolEnvVar("ORCA_INJECT", false), "enable injection during helm upgrade. Overrides $ORCA_INJECT (requires helm inject plugin: https://github.com/maorfr/helm-inject)")
	f.BoolVarP(&e.deployOnlyOverrideIfEnvExists, "deploy-only-override-if-env-exists", "x", false, "if environment exists - deploy only override(s) (support for features spanning multiple services). Overrides $ORCA_DEPLOY_ONLY_OVERRIDE_IF_ENV_EXISTS")

	f.BoolVar(&e.createNS, "create-ns", utils.GetBoolEnvVar("ORCA_CREATE_NS", false), "should create new namespace. Overrides $ORCA_CREATE_NS")
	f.MarkDeprecated("create-ns", "namespace will be created if it does not exist")
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
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			if e.tls && e.helmTLSStore == "" {
				return errors.New("tls is set to true and helm-tls-store is not defined")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.NamespaceExists(e.name, e.kubeContext) {
				log.Printf("environment \"%s\" not found", e.name)
				return
			}
			print := false
			includeFailed := true
			markEnvironmentForDeletion(e.name, e.kubeContext, e.force, print)
			releases := utils.GetInstalledReleases(e.kubeContext, e.name, includeFailed)
			utils.DeleteReleases(releases, e.kubeContext, e.helmTLSStore, e.tls)
			utils.DeleteNamespace(e.name, e.kubeContext, print)
			log.Printf("deleted environment \"%s\"", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.force, "force", utils.GetBoolEnvVar("ORCA_FORCE", false), "force environment deletion. Overrides $ORCA_FORCE")

	return cmd
}

// NewLockEnvCmd represents the lock env command
func NewLockEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Lock an environment (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.NamespaceExists(e.name, e.kubeContext) {
				log.Printf("environment \"%s\" not found", e.name)
				return
			}
			print := false
			lockEnvironment(e.name, e.kubeContext, print)
			log.Printf("locked environment \"%s\"", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")

	return cmd
}

// NewUnlockEnvCmd represents the unlock env command
func NewUnlockEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Unlock an environment (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.NamespaceExists(e.name, e.kubeContext) {
				log.Printf("environment \"%s\" not found", e.name)
				return
			}
			print := false
			unlockEnvironment(e.name, e.kubeContext, print)
			log.Printf("unlocked environment \"%s\"", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")

	return cmd
}

const (
	annotationPrefix string = "orca.nuvocares.com"
	stateAnnotation  string = annotationPrefix + "/state"
	busyState        string = "busy"
	freeState        string = "free"
	deleteState      string = "delete"
)

// lockEnvironment annotates a namespace with "busy"
func lockEnvironment(name, kubeContext string, print bool) {
	sleepPeriod := 5 * time.Second
	state := utils.GetNamespace(name, kubeContext).Annotations[stateAnnotation]
	if state != "" {
		if state != freeState && state != busyState {
			log.Fatal("Environment state is ", state)
		}
		for state == busyState {
			log.Printf("environment \"%s\" %s, backing off for %d seconds", name, busyState, int(sleepPeriod.Seconds()))
			time.Sleep(sleepPeriod)
			sleepPeriod += 5 * time.Second
			state = utils.GetNamespace(name, kubeContext).Annotations[stateAnnotation]
		}
	}
	// There is a race condition here, may need to attend to it in the future
	annotations := map[string]string{stateAnnotation: busyState}
	utils.UpdateNamespace(name, kubeContext, annotations, print)
}

// unlockEnvironment annotates a namespace with "free"
func unlockEnvironment(name, kubeContext string, print bool) {
	annotations := map[string]string{stateAnnotation: freeState}
	utils.UpdateNamespace(name, kubeContext, annotations, print)
}

// markEnvironmentForDeletion annotates a namespace with "delete"
func markEnvironmentForDeletion(name, kubeContext string, force, print bool) {
	if !force {
		lockEnvironment(name, kubeContext, print)
	}
	annotations := map[string]string{stateAnnotation: deleteState}
	utils.UpdateNamespace(name, kubeContext, annotations, print)
}
