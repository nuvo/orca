package orca

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nuvo/orca/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	annotationPrefix string = "orca.nuvocares.com"
	stateAnnotation  string = annotationPrefix + "/state"
	busyState        string = "busy"
	freeState        string = "free"
	deleteState      string = "delete"
	failedState      string = "failed"
	unknownState     string = "unknown"
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
	parallel                      int
	timeout                       int
	annotations                   []string
	labels                        []string
	validate                      bool

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
			releases, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
				KubeContext:   e.kubeContext,
				Namespace:     e.name,
				IncludeFailed: false,
			})
			if err != nil {
				log.Fatal(err)
			}

			switch e.output {
			case "yaml":
				utils.PrintReleasesYaml(releases)
			case "md":
				utils.PrintReleasesMarkdown(releases)
			case "table":
				utils.PrintReleasesTable(releases)
			case "":
				utils.PrintReleasesYaml(releases)
			}
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to get. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.StringVarP(&e.output, "output", "o", os.Getenv("ORCA_OUTPUT"), "output format (yaml, md, table). Overrides $ORCA_OUTPUT")

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
			if e.tls {
				if e.helmTLSStore == "" {
					return errors.New("tls is set to true and helm-tls-store is not defined")
				}
				if e.kubeContext == "" {
					return errors.New("kube-context has to be non-empty when tls is set to true")
				}
			}
			if len(e.override) == 0 {
				if e.chartsFile == "" {
					return errors.New("either charts-file or override has to be defined")
				}
				if e.deployOnlyOverrideIfEnvExists {
					return errors.New("override has to be defined when using deploy-only-override-if-env-exists")
				}
			}
			if e.chartsFile != "" && utils.CheckCircularDependencies(utils.InitReleasesFromChartsFile(e.chartsFile, e.name)) {
				return errors.New("Circular dependency found")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			log.Println("initializing chart repository configuration")
			if err := utils.AddRepository(utils.AddRepositoryOptions{
				Repo:  e.repo,
				Print: false,
			}); err != nil {
				log.Fatal(err)
			}
			if err := utils.UpdateRepositories(false); err != nil {
				log.Fatal(err)
			}

			log.Printf("deploying environment \"%s\"", e.name)
			nsPreExists, err := utils.NamespaceExists(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}
			if !nsPreExists {
				if err := utils.CreateNamespace(e.name, e.kubeContext, false); err != nil {
					log.Fatal(err)
				}
				log.Printf("created environment \"%s\"", e.name)
			}
			if err := lockEnvironment(e.name, e.kubeContext, true); err != nil {
				log.Fatal(err)
			}

			annotations := map[string]string{}
			for _, a := range e.annotations {
				k, v := utils.SplitInTwo(a, "=")
				annotations[k] = v
			}
			labels := map[string]string{}
			for _, a := range e.labels {
				k, v := utils.SplitInTwo(a, "=")
				labels[k] = v
			}
			if err := utils.UpdateNamespace(e.name, e.kubeContext, annotations, labels, true); err != nil {
				log.Fatal(err)
			}

			log.Print("initializing releases to deploy")
			var desiredReleases []utils.ReleaseSpec
			if nsPreExists && e.deployOnlyOverrideIfEnvExists {
				desiredReleases = utils.InitReleases(e.name, e.override)
			} else {
				if e.chartsFile != "" {
					desiredReleases = utils.InitReleasesFromChartsFile(e.chartsFile, e.name)
				}
				desiredReleases = utils.OverrideReleases(desiredReleases, e.override, e.name)
			}

			log.Print("getting currently deployed releases")
			installedReleases, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
				KubeContext:   e.kubeContext,
				Namespace:     e.name,
				IncludeFailed: false,
			})
			if err != nil {
				unlockEnvironment(e.name, e.kubeContext, true)
				log.Fatal(err)
			}
			log.Print("calculating delta between desired releases and currently deployed releases")
			releasesToInstall := utils.GetReleasesDelta(desiredReleases, installedReleases)

			log.Print("deploying releases")
			if err := utils.DeployChartsFromRepository(utils.DeployChartsFromRepositoryOptions{
				ReleasesToInstall: releasesToInstall,
				KubeContext:       e.kubeContext,
				Namespace:         e.name,
				Repo:              e.repo,
				TLS:               e.tls,
				HelmTLSStore:      e.helmTLSStore,
				PackedValues:      e.packedValues,
				SetValues:         e.set,
				Inject:            e.inject,
				Parallel:          e.parallel,
				Timeout:           e.timeout,
			}); err != nil {
				markEnvironmentAsFailed(e.name, e.kubeContext, true)
				log.Fatal(err)
			}

			if !e.deployOnlyOverrideIfEnvExists {
				log.Print("getting currently deployed releases")
				installedReleases, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
					KubeContext:   e.kubeContext,
					Namespace:     e.name,
					IncludeFailed: false,
				})
				if err != nil {
					markEnvironmentAsUnknown(e.name, e.kubeContext, true)
					log.Fatal(err)
				}
				log.Print("calculating delta between desired releases and currently deployed releases")
				releasesToDelete := utils.GetReleasesDelta(installedReleases, desiredReleases)
				log.Print("deleting undesired releases")
				if err := utils.DeleteReleases(utils.DeleteReleasesOptions{
					ReleasesToDelete: releasesToDelete,
					KubeContext:      e.kubeContext,
					TLS:              e.tls,
					HelmTLSStore:     e.helmTLSStore,
					Parallel:         e.parallel,
					Timeout:          e.timeout,
				}); err != nil {
					markEnvironmentAsFailed(e.name, e.kubeContext, true)
					log.Fatal(err)
				}
			}
			log.Printf("deployed environment \"%s\"", e.name)

			var envValid bool
			if e.validate {
				envValid, err = utils.IsEnvValidWithLoopBackOff(e.name, e.kubeContext)
			}

			unlockEnvironment(e.name, e.kubeContext, true)

			if !e.validate {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			if !envValid {
				markEnvironmentAsFailed(e.name, e.kubeContext, true)
				log.Fatalf("environment \"%s\" validation failed!", e.name)
			}
			// If we have made it so far, the environment is validated
			log.Printf("environment \"%s\" validated!", e.name)
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
	f.IntVarP(&e.parallel, "parallel", "p", utils.GetIntEnvVar("ORCA_PARALLEL", 1), "number of releases to act on in parallel. set this flag to 0 for full parallelism. Overrides $ORCA_PARALLEL")
	f.IntVar(&e.timeout, "timeout", utils.GetIntEnvVar("ORCA_TIMEOUT", 300), "time in seconds to wait for any individual Kubernetes operation (like Jobs for hooks). Overrides $ORCA_TIMEOUT")
	f.StringSliceVar(&e.annotations, "annotations", []string{}, "additional environment (namespace) annotations (can specify multiple): annotation=value")
	f.StringSliceVar(&e.labels, "labels", []string{}, "environment (namespace) labels (can specify multiple): label=value")
	f.BoolVar(&e.validate, "validate", utils.GetBoolEnvVar("ORCA_VALIDATE", false), "perform environment validation after deployment. Overrides $ORCA_VALIDATE")

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
			if e.tls && e.kubeContext == "" {
				return errors.New("kube-context has to be non-empty when tls is set to true")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			nsExists, err := utils.NamespaceExists(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}
			if nsExists {
				if err := markEnvironmentForDeletion(e.name, e.kubeContext, e.force, true); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("environment \"%s\" not found", e.name)
			}

			log.Print("getting currently deployed releases")
			releases, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
				KubeContext:   e.kubeContext,
				Namespace:     e.name,
				IncludeFailed: true,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Print("deleting releases")
			if err := utils.DeleteReleases(utils.DeleteReleasesOptions{
				ReleasesToDelete: releases,
				KubeContext:      e.kubeContext,
				TLS:              e.tls,
				HelmTLSStore:     e.helmTLSStore,
				Parallel:         e.parallel,
				Timeout:          e.timeout,
			}); err != nil {
				markEnvironmentAsFailed(e.name, e.kubeContext, true)
				log.Fatal(err)
			}

			if nsExists {
				if utils.Contains([]string{"default", "kube-system", "kube-public"}, e.name) {
					removeStateAnnotationsFromEnvironment(e.name, e.kubeContext, true)
				} else {
					utils.DeleteNamespace(e.name, e.kubeContext, false)
				}
			}
			log.Printf("deleted environment \"%s\"", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")
	f.BoolVar(&e.tls, "tls", utils.GetBoolEnvVar("ORCA_TLS", false), "enable TLS for request. Overrides $ORCA_TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", os.Getenv("HELM_TLS_STORE"), "path to TLS certs and keys. Overrides $HELM_TLS_STORE")
	f.BoolVar(&e.force, "force", utils.GetBoolEnvVar("ORCA_FORCE", false), "force environment deletion. Overrides $ORCA_FORCE")
	f.IntVarP(&e.parallel, "parallel", "p", utils.GetIntEnvVar("ORCA_PARALLEL", 1), "number of releases to act on in parallel. set this flag to 0 for full parallelism. Overrides $ORCA_PARALLEL")
	f.IntVar(&e.timeout, "timeout", utils.GetIntEnvVar("ORCA_TIMEOUT", 300), "time in seconds to wait for any individual Kubernetes operation (like Jobs for hooks). Overrides $ORCA_TIMEOUT")

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
			nsExists, err := utils.NamespaceExists(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}
			if !nsExists {
				log.Printf("environment \"%s\" not found", e.name)
				return
			}
			if err := lockEnvironment(e.name, e.kubeContext, false); err != nil {
				log.Fatal(err)
			}
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
			nsExists, err := utils.NamespaceExists(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}
			if !nsExists {
				log.Printf("environment \"%s\" not found", e.name)
				return
			}
			if err := unlockEnvironment(e.name, e.kubeContext, false); err != nil {
				log.Fatal(err)
			}
			log.Printf("unlocked environment \"%s\"", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")

	return cmd
}

type diffEnvCmd struct {
	nameLeft         string
	nameRight        string
	kubeContextLeft  string
	kubeContextRight string
	output           string

	out io.Writer
}

// NewDiffEnvCmd represents the diff env command
func NewDiffEnvCmd(out io.Writer) *cobra.Command {
	e := &diffEnvCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Show differences in Helm releases between environments (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.nameLeft == "" {
				return errors.New("name-left can not be empty")
			}
			if e.nameRight == "" {
				return errors.New("name-right can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			releasesLeft, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
				KubeContext:   e.kubeContextLeft,
				Namespace:     e.nameLeft,
				IncludeFailed: false,
			})
			if err != nil {
				log.Fatal(err)
			}
			releasesRight, err := utils.GetInstalledReleases(utils.GetInstalledReleasesOptions{
				KubeContext:   e.kubeContextRight,
				Namespace:     e.nameRight,
				IncludeFailed: false,
			})
			if err != nil {
				log.Fatal(err)
			}

			diffOptions := utils.DiffOptions{
				KubeContextLeft:   e.kubeContextLeft,
				KubeContextRight:  e.kubeContextRight,
				EnvNameLeft:       e.nameLeft,
				EnvNameRight:      e.nameRight,
				ReleasesSpecLeft:  releasesLeft,
				ReleasesSpecRight: releasesRight,
				Output:            e.output,
			}
			utils.PrintDiff(diffOptions)
		},
	}

	f := cmd.Flags()

	f.StringVar(&e.nameLeft, "name-left", os.Getenv("ORCA_NAME_LEFT"), "name of left environment to compare. Overrides $ORCA_NAME_LEFT")
	f.StringVar(&e.nameRight, "name-right", os.Getenv("ORCA_NAME_RIGHT"), "name of right environment to compare. Overrides $ORCA_NAME_RIGHT")
	f.StringVar(&e.kubeContextLeft, "kube-context-left", os.Getenv("ORCA_KUBE_CONTEXT_LEFT"), "name of the left kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT_LEFT")
	f.StringVar(&e.kubeContextRight, "kube-context-right", os.Getenv("ORCA_KUBE_CONTEXT_RIGHT"), "name of the right kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT_RIGHT")
	f.StringVarP(&e.output, "output", "o", os.Getenv("ORCA_OUTPUT"), "output format (yaml, md, table). Overrides $ORCA_OUTPUT")

	return cmd
}

// NewValidateEnvCmd represents the validate env command
func NewValidateEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Validate an environment (Kubernetes namespace)",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if e.name == "" {
				return errors.New("name can not be empty")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("validating environment \"%s\"", e.name)
			nsExists, err := utils.NamespaceExists(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}
			if !nsExists {
				log.Fatalf("environment \"%s\" not found", e.name)
			}

			envValid, err := utils.IsEnvValid(e.name, e.kubeContext)
			if err != nil {
				log.Fatal(err)
			}

			if !envValid {
				log.Fatalf("environment \"%s\" validation failed!", e.name)
			}
			// If we have made it so far, the environment is validated
			log.Printf("environment \"%s\" validated!", e.name)
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.name, "name", "n", os.Getenv("ORCA_NAME"), "name of environment (namespace) to delete. Overrides $ORCA_NAME")
	f.StringVar(&e.kubeContext, "kube-context", os.Getenv("ORCA_KUBE_CONTEXT"), "name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT")

	return cmd
}

// lockEnvironment annotates a namespace with "busy"
func lockEnvironment(name, kubeContext string, print bool) error {
	sleepPeriod := 5 * time.Second
	ns, err := utils.GetNamespace(name, kubeContext)
	if err != nil {
		return err
	}
	state := ns.Annotations[stateAnnotation]
	if state != "" {
		if state != freeState && state != busyState {
			return fmt.Errorf("Environment state is %s", state)
		}
		for state == busyState {
			log.Printf("environment \"%s\" %s, backing off for %d seconds", name, busyState, int(sleepPeriod.Seconds()))
			time.Sleep(sleepPeriod)
			sleepPeriod += 5 * time.Second
			ns, err := utils.GetNamespace(name, kubeContext)
			if err != nil {
				return err
			}
			state = ns.Annotations[stateAnnotation]
		}
	}
	// There is a race condition here, may need to attend to it in the future
	annotations := map[string]string{stateAnnotation: busyState}
	err = utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}

// unlockEnvironment annotates a namespace with "free"
func unlockEnvironment(name, kubeContext string, print bool) error {
	ns, err := utils.GetNamespace(name, kubeContext)
	if err != nil {
		return err
	}
	state := ns.Annotations[stateAnnotation]
	if state != "" {
		if state != freeState && state != busyState {
			return fmt.Errorf("Environment state is %s", state)
		}
	}
	annotations := map[string]string{stateAnnotation: freeState}
	err = utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}

// markEnvironmentForDeletion annotates a namespace with "delete"
func markEnvironmentForDeletion(name, kubeContext string, force, print bool) error {
	if !force {
		if err := lockEnvironment(name, kubeContext, print); err != nil {
			return err
		}
	}
	annotations := map[string]string{stateAnnotation: deleteState}
	err := utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}

// markEnvironmentAsFailed annotates a namespace with "failed"
func markEnvironmentAsFailed(name, kubeContext string, print bool) error {
	annotations := map[string]string{stateAnnotation: failedState}
	err := utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}

// markEnvironmentAsUnknown annotates a namespace with "unknown"
func markEnvironmentAsUnknown(name, kubeContext string, print bool) error {
	annotations := map[string]string{stateAnnotation: unknownState}
	err := utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}

// unlockEnvironment annotates a namespace with "unknown"
func removeStateAnnotationsFromEnvironment(name, kubeContext string, print bool) error {
	annotations := map[string]string{}
	err := utils.UpdateNamespace(name, kubeContext, annotations, map[string]string{}, print)

	return err
}
