package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type DeployChartsFromRepositoryOptions struct {
	ReleasesToInstall []ReleaseSpec
	KubeContext       string
	Namespace         string
	Repo              string
	TLS               bool
	HelmTLSStore      string
	PackedValues      []string
	SetValues         []string
	Inject            bool
	Parallel          int
	Timeout           int
}

// DeployChartsFromRepository deploys a list of Helm charts from a repository in parallel
func DeployChartsFromRepository(o DeployChartsFromRepositoryOptions) {

	releasesToInstall := o.ReleasesToInstall
	parallel := o.Parallel

	totalReleases := len(releasesToInstall)
	if parallel == 0 {
		parallel = totalReleases
	}
	bwgSize := int(math.Min(float64(parallel), float64(totalReleases))) // Very stingy :)
	bwg := NewBoundedWaitGroup(bwgSize)
	var mutex = &sync.Mutex{}

	for len(releasesToInstall) > 0 {

		for _, c := range releasesToInstall {

			bwg.Add(1)
			go func(c ReleaseSpec) {
				defer bwg.Done()

				// If there are (still) any dependencies - leave this chart for a later iteration
				if len(c.Dependencies) != 0 {
					return
				}

				// Find index of chart in slice
				// may have changed by now since we are using go routines
				// If chart was not found - another routine is taking care of it
				mutex.Lock()
				index := GetChartIndex(releasesToInstall, c.ChartName)
				if index == -1 {
					mutex.Unlock()
					return
				}
				releasesToInstall = RemoveChartFromCharts(releasesToInstall, index)
				mutex.Unlock()

				// deploy chart
				log.Println("deploying chart", c.ChartName, "version", c.ChartVersion)
				DeployChartFromRepository(DeployChartFromRepositoryOptions{
					ReleaseName:  c.ReleaseName,
					Name:         c.ChartName,
					Version:      c.ChartVersion,
					KubeContext:  o.KubeContext,
					Namespace:    o.Namespace,
					Repo:         o.Repo,
					TLS:          o.TLS,
					HelmTLSStore: o.HelmTLSStore,
					PackedValues: o.PackedValues,
					SetValues:    o.SetValues,
					IsIsolated:   false,
					Inject:       o.Inject,
					Timeout:      o.Timeout,
				})
				log.Println("deployed chart", c.ChartName, "version", c.ChartVersion)

				// Deployment is done, remove chart from dependencies
				mutex.Lock()
				releasesToInstall = RemoveChartFromDependencies(releasesToInstall, c.ChartName)
				mutex.Unlock()
			}(c)
		}
		time.Sleep(5 * time.Second)
	}
	bwg.Wait()
}

type DeployChartFromRepositoryOptions struct {
	ReleaseName  string
	Name         string
	Version      string
	KubeContext  string
	Namespace    string
	Repo         string
	TLS          bool
	HelmTLSStore string
	PackedValues []string
	SetValues    []string
	IsIsolated   bool
	Inject       bool
	Timeout      int
}

// DeployChartFromRepository deploys a Helm chart from a chart repository
func DeployChartFromRepository(o DeployChartFromRepositoryOptions) {
	tempDir := MkRandomDir()

	if o.ReleaseName == "" {
		o.ReleaseName = o.Name
	}
	if o.IsIsolated {
		AddRepository(AddRepositoryOptions{
			Repo:  o.Repo,
			Print: o.IsIsolated,
		})
		UpdateRepositories(o.IsIsolated)
	}
	FetchChart(FetchChartOptions{
		Repo:    o.Repo,
		Name:    o.Name,
		Version: o.Version,
		Dir:     tempDir,
		Print:   o.IsIsolated,
	})
	path := fmt.Sprintf("%s/%s", tempDir, o.Name)
	UpdateChartDependencies(UpdateChartDependenciesOptions{
		Path:  path,
		Print: o.IsIsolated,
	})
	valuesChain := CreateValuesChain(o.Name, tempDir, o.PackedValues)
	setChain := CreateSetChain(o.Name, o.SetValues)

	UpgradeRelease(UpgradeReleaseOptions{
		Name:         o.Name,
		ReleaseName:  o.ReleaseName,
		KubeContext:  o.KubeContext,
		Namespace:    o.Namespace,
		Values:       valuesChain,
		Set:          setChain,
		TLS:          o.TLS,
		HelmTLSStore: o.HelmTLSStore,
		Dir:          tempDir,
		Print:        o.IsIsolated,
		Inject:       o.Inject,
		Timeout:      o.Timeout,
	})

	os.RemoveAll(tempDir)
}

type LintOptions struct {
	Path  string
	Print bool
}

// Lint takes a path to a chart and runs a series of tests to verify that the chart is well-formed
func Lint(o LintOptions) {
	cmd := []string{"helm", "lint", o.Path}
	PrintExec(cmd, o.Print)
}

type AddRepositoryOptions struct {
	Repo  string
	Print bool
}

// AddRepository adds a chart repository to the repositories file
func AddRepository(o AddRepositoryOptions) {
	repoName, repoURL := SplitInTwo(o.Repo, "=")

	cmd := []string{
		"helm", "repo",
		"add", repoName, repoURL,
	}
	PrintExec(cmd, o.Print)
}

// UpdateRepositories updates helm repositories
func UpdateRepositories(print bool) {
	cmd := []string{"helm", "repo", "update"}
	PrintExec(cmd, print)
}

type FetchChartOptions struct {
	Repo    string
	Name    string
	Version string
	Dir     string
	Print   bool
}

// FetchChart fetches a chart from chart repository by name and version and untars it in the local directory
func FetchChart(o FetchChartOptions) {
	repoName, _ := SplitInTwo(o.Repo, "=")

	cmd := []string{
		"helm", "fetch",
		fmt.Sprintf("%s/%s", repoName, o.Name),
		"--version", o.Version,
		"--untar",
		"-d", o.Dir,
	}
	PrintExec(cmd, o.Print)
}

type PushChartToRepositoryOptions struct {
	Path   string
	Append string
	Repo   string
	Lint   bool
	Print  bool
}

// PushChartToRepository packages and pushes a Helm chart to a chart repository
func PushChartToRepository(o PushChartToRepositoryOptions) {
	newVersion := UpdateChartVersion(o.Path, o.Append)
	if o.Lint {
		Lint(LintOptions{
			Path:  o.Path,
			Print: o.Print,
		})
	}
	AddRepository(AddRepositoryOptions{
		Repo:  o.Repo,
		Print: o.Print,
	})
	UpdateChartDependencies(UpdateChartDependenciesOptions{
		Path:  o.Path,
		Print: o.Print,
	})
	PushChart(PushChartOptions{
		Repo:  o.Repo,
		Path:  o.Path,
		Print: o.Print,
	})
	fmt.Println(newVersion)
}

type PushChartOptions struct {
	Repo  string
	Path  string
	Print bool
}

// PushChart pushes a helm chart to a chart repository
func PushChart(o PushChartOptions) {
	repoName, _ := SplitInTwo(o.Repo, "=")

	cmd := []string{"helm", "push", o.Path, repoName}
	PrintExec(cmd, o.Print)
}

type UpdateChartDependenciesOptions struct {
	Path  string
	Print bool
}

// UpdateChartDependencies performs helm dependency update
func UpdateChartDependencies(o UpdateChartDependenciesOptions) {
	cmd := []string{"helm", "dependency", "update", o.Path}
	PrintExec(cmd, o.Print)
}

// CreateValuesChain will create a chain of values files to use
func CreateValuesChain(name, dir string, packedValues []string) []string {
	var values []string
	format := "%s/%s/%s"
	fileToTest := fmt.Sprintf(format, dir, name, "values.yaml")
	if _, err := os.Stat(fileToTest); err == nil {
		values = append(values, "-f", fileToTest)
	}
	for _, v := range packedValues {
		fileToTest = fmt.Sprintf(format, dir, name, v)
		_, err := os.Stat(fileToTest)
		if err != nil {
			continue
		}
		if Contains(values, fileToTest) {
			continue
		}
		values = append(values, "-f", fileToTest)
	}
	return values
}

// CreateSetChain will create a chain of sets to use
func CreateSetChain(name string, inputSet []string) []string {
	set := []string{"--set", fmt.Sprintf("fullnameOverride=%s", name)}
	for _, s := range inputSet {
		set = append(set, "--set", s)
	}
	return set
}

type UpgradeReleaseOptions struct {
	Name         string
	ReleaseName  string
	KubeContext  string
	Namespace    string
	Values       []string
	Set          []string
	TLS          bool
	HelmTLSStore string
	Dir          string
	Print        bool
	Inject       bool
	Timeout      int
}

// UpgradeRelease performs helm upgrade -i
func UpgradeRelease(o UpgradeReleaseOptions) {
	cmd := []string{"helm"}
	kubeContextFlag := "--kube-context"
	if o.Inject {
		kubeContextFlag = "--kubecontext"
		cmd = append(cmd, "inject")
	}
	cmd = append(cmd, "upgrade", "-i", o.ReleaseName, fmt.Sprintf("%s/%s", o.Dir, o.Name))
	if o.KubeContext != "" {
		cmd = append(cmd, kubeContextFlag, o.KubeContext)
	}
	if o.Namespace != "" {
		cmd = append(cmd, "--namespace", o.Namespace)
	}
	cmd = append(cmd, o.Values...)
	cmd = append(cmd, o.Set...)
	cmd = append(cmd, "--timeout", fmt.Sprintf("%d", o.Timeout))
	cmd = append(cmd, getTLS(o.TLS, o.KubeContext, o.HelmTLSStore)...)
	PrintExec(cmd, o.Print)
}

type DeleteReleasesOptions struct {
	ReleasesToDelete []ReleaseSpec
	KubeContext      string
	TLS              bool
	HelmTLSStore     string
	Parallel         int
	Timeout          int
}

// DeleteReleases deletes a list of releases in parallel
func DeleteReleases(o DeleteReleasesOptions) {
	releasesToDelete := o.ReleasesToDelete
	parallel := o.Parallel

	print := false
	totalReleases := len(releasesToDelete)
	if parallel == 0 {
		parallel = totalReleases
	}
	bwgSize := int(math.Min(float64(parallel), float64(totalReleases))) // Very stingy :)
	bwg := NewBoundedWaitGroup(bwgSize)

	for _, c := range releasesToDelete {
		bwg.Add(1)
		go func(c ReleaseSpec) {
			defer bwg.Done()
			log.Println("deleting", c.ReleaseName)
			DeleteRelease(DeleteReleaseOptions{
				ReleaseName:  c.ReleaseName,
				KubeContext:  o.KubeContext,
				TLS:          o.TLS,
				HelmTLSStore: o.HelmTLSStore,
				Timeout:      o.Timeout,
				Print:        print,
			})
			log.Println("deleted", c.ReleaseName)
		}(c)
	}
	bwg.Wait()
}

type DeleteReleaseOptions struct {
	ReleaseName  string
	KubeContext  string
	TLS          bool
	HelmTLSStore string
	Timeout      int
	Print        bool
}

// DeleteRelease deletes a release from Kubernetes
func DeleteRelease(o DeleteReleaseOptions) {
	cmd := []string{
		"helm", "delete", o.ReleaseName, "--purge",
		"--timeout", fmt.Sprintf("%d", o.Timeout),
	}
	if o.KubeContext != "" {
		cmd = append(cmd, "--kube-context", o.KubeContext)
	}
	cmd = append(cmd, getTLS(o.TLS, o.KubeContext, o.HelmTLSStore)...)
	PrintExec(cmd, o.Print)
}

func getTLS(tls bool, kubeContext, helmTLSStore string) []string {
	var tlsStr []string
	if tls == true {
		tlsStr = []string{
			"--tls",
			"--tls-cert", fmt.Sprintf("%s/%s.cert.pem", helmTLSStore, kubeContext),
			"--tls-key", fmt.Sprintf("%s/%s.key.pem", helmTLSStore, kubeContext),
		}
	}
	return tlsStr
}
