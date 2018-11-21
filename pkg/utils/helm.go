package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

// DeployChartsFromRepositoryOptions are options passed to DeployChartsFromRepository
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
func DeployChartsFromRepository(o DeployChartsFromRepositoryOptions) error {
	releasesToInstall := o.ReleasesToInstall
	if len(releasesToInstall) == 0 {
		return nil
	}
	parallel := o.Parallel

	totalReleases := len(releasesToInstall)
	if parallel == 0 {
		parallel = totalReleases
	}
	bwgSize := int(math.Min(float64(parallel), float64(totalReleases))) // Very stingy :)
	bwg := NewBoundedWaitGroup(bwgSize)
	errc := make(chan error, 1)
	var mutex = &sync.Mutex{}

	for len(releasesToInstall) > 0 && len(errc) == 0 {

		for _, r := range releasesToInstall {

			if len(errc) != 0 {
				break
			}

			bwg.Add(1)
			go func(r ReleaseSpec) {
				defer bwg.Done()

				// If there has been an error in a concurrent deployment - don`t deploy anymore
				if len(errc) != 0 {
					return
				}

				// If there are (still) any dependencies - leave this chart for a later iteration
				if len(r.Dependencies) != 0 {
					return
				}

				// Find index of chart in slice
				// may have changed by now since we are using go routines
				// If chart was not found - another routine is taking care of it
				mutex.Lock()
				index := GetChartIndex(releasesToInstall, r.ChartName)
				if index == -1 {
					mutex.Unlock()
					return
				}
				releasesToInstall = RemoveChartFromCharts(releasesToInstall, index)
				mutex.Unlock()

				// deploy chart
				log.Println("deploying chart", r.ChartName, "version", r.ChartVersion)
				if err := DeployChartFromRepository(DeployChartFromRepositoryOptions{
					ReleaseName:  r.ReleaseName,
					Name:         r.ChartName,
					Version:      r.ChartVersion,
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
				}); err != nil {
					log.Println("failed deploying chart", r.ChartName, "version", r.ChartVersion)
					errc <- err
					return
				}
				log.Println("deployed chart", r.ChartName, "version", r.ChartVersion)

				// Deployment is done, remove chart from dependencies
				mutex.Lock()
				releasesToInstall = RemoveChartFromDependencies(releasesToInstall, r.ChartName)
				mutex.Unlock()
			}(r)

		}
		time.Sleep(5 * time.Second)
	}
	bwg.Wait()

	if len(errc) != 0 {
		// This is not exactly the correct behavior
		// There may be more than 1 error in the channel
		// But first let's make it work
		err := <-errc
		close(errc)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteReleasesOptions are options passed to DeleteReleases
type DeleteReleasesOptions struct {
	ReleasesToDelete []ReleaseSpec
	KubeContext      string
	TLS              bool
	HelmTLSStore     string
	Parallel         int
	Timeout          int
}

// DeleteReleases deletes a list of releases in parallel
func DeleteReleases(o DeleteReleasesOptions) error {
	releasesToDelete := o.ReleasesToDelete
	if len(releasesToDelete) == 0 {
		return nil
	}
	parallel := o.Parallel

	print := false
	totalReleases := len(releasesToDelete)
	if parallel == 0 {
		parallel = totalReleases
	}
	bwgSize := int(math.Min(float64(parallel), float64(totalReleases))) // Very stingy :)
	bwg := NewBoundedWaitGroup(bwgSize)
	errc := make(chan error, 1)

	for _, r := range releasesToDelete {
		bwg.Add(1)
		go func(r ReleaseSpec) {
			defer bwg.Done()
			log.Println("deleting", r.ReleaseName)
			if err := DeleteRelease(DeleteReleaseOptions{
				ReleaseName:  r.ReleaseName,
				KubeContext:  o.KubeContext,
				TLS:          o.TLS,
				HelmTLSStore: o.HelmTLSStore,
				Timeout:      o.Timeout,
				Print:        print,
			}); err != nil {
				log.Println("failed deleting chart", r.ReleaseName)
				errc <- err
				return
			}
			log.Println("deleted", r.ReleaseName)
		}(r)
	}
	bwg.Wait()

	if len(errc) != 0 {
		// This is not exactly the correct behavior
		// There may be more than 1 error in the channel
		// But first let's make it work
		err := <-errc
		close(errc)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeployChartFromRepositoryOptions are options passed to DeployChartFromRepository
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
func DeployChartFromRepository(o DeployChartFromRepositoryOptions) error {
	tempDir := MkRandomDir()

	if o.ReleaseName == "" {
		o.ReleaseName = o.Name
	}
	if o.IsIsolated {
		if err := AddRepository(AddRepositoryOptions{
			Repo:  o.Repo,
			Print: o.IsIsolated,
		}); err != nil {
			return err
		}
		if err := UpdateRepositories(o.IsIsolated); err != nil {
			return err
		}
	}
	if err := FetchChart(FetchChartOptions{
		Repo:    o.Repo,
		Name:    o.Name,
		Version: o.Version,
		Dir:     tempDir,
		Print:   o.IsIsolated,
	}); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", tempDir, o.Name)
	if err := UpdateChartDependencies(UpdateChartDependenciesOptions{
		Path:  path,
		Print: o.IsIsolated,
	}); err != nil {
		return err
	}
	valuesChain := createValuesChain(o.Name, tempDir, o.PackedValues)
	setChain := createSetChain(o.Name, o.SetValues)

	if err := UpgradeRelease(UpgradeReleaseOptions{
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
	}); err != nil {
		return err
	}

	os.RemoveAll(tempDir)
	return nil
}

// PushChartToRepositoryOptions are options passed to PushChartToRepository
type PushChartToRepositoryOptions struct {
	Path   string
	Append string
	Repo   string
	Lint   bool
	Print  bool
}

// PushChartToRepository packages and pushes a Helm chart to a chart repository
func PushChartToRepository(o PushChartToRepositoryOptions) error {
	newVersion := UpdateChartVersion(o.Path, o.Append)
	if o.Lint {
		if err := Lint(LintOptions{
			Path:  o.Path,
			Print: o.Print,
		}); err != nil {
			return err
		}
	}
	if err := AddRepository(AddRepositoryOptions{
		Repo:  o.Repo,
		Print: o.Print,
	}); err != nil {
		return err
	}
	if err := UpdateChartDependencies(UpdateChartDependenciesOptions{
		Path:  o.Path,
		Print: o.Print,
	}); err != nil {
		return err
	}
	if err := PushChart(PushChartOptions{
		Repo:  o.Repo,
		Path:  o.Path,
		Print: o.Print,
	}); err != nil {
		return err
	}
	fmt.Println(newVersion)
	return nil
}

// LintOptions are options passed to Lint
type LintOptions struct {
	Path  string
	Print bool
}

// Lint takes a path to a chart and runs a series of tests to verify that the chart is well-formed
func Lint(o LintOptions) error {
	cmd := []string{"helm", "lint", o.Path}
	err := PrintExec(cmd, o.Print)

	return err
}

// AddRepositoryOptions are options passed to AddRepository
type AddRepositoryOptions struct {
	Repo  string
	Print bool
}

// AddRepository adds a chart repository to the repositories file
func AddRepository(o AddRepositoryOptions) error {
	repoName, repoURL := SplitInTwo(o.Repo, "=")

	cmd := []string{
		"helm", "repo",
		"add", repoName, repoURL,
	}
	err := PrintExec(cmd, o.Print)

	return err
}

// UpdateRepositories updates helm repositories
func UpdateRepositories(print bool) error {
	cmd := []string{"helm", "repo", "update"}
	err := PrintExec(cmd, print)

	return err
}

// FetchChartOptions are options passed to FetchChart
type FetchChartOptions struct {
	Repo    string
	Name    string
	Version string
	Dir     string
	Print   bool
}

// FetchChart fetches a chart from chart repository by name and version and untars it in the local directory
func FetchChart(o FetchChartOptions) error {
	repoName, _ := SplitInTwo(o.Repo, "=")

	cmd := []string{
		"helm", "fetch",
		fmt.Sprintf("%s/%s", repoName, o.Name),
		"--version", o.Version,
		"--untar",
		"-d", o.Dir,
	}
	err := PrintExec(cmd, o.Print)

	return err
}

// PushChartOptions are options passed to PushChart
type PushChartOptions struct {
	Repo  string
	Path  string
	Print bool
}

// PushChart pushes a helm chart to a chart repository
func PushChart(o PushChartOptions) error {
	repoName, _ := SplitInTwo(o.Repo, "=")

	cmd := []string{"helm", "push", o.Path, repoName}
	err := PrintExec(cmd, o.Print)

	return err
}

// UpdateChartDependenciesOptions are options passed to UpdateChartDependencies
type UpdateChartDependenciesOptions struct {
	Path  string
	Print bool
}

// UpdateChartDependencies performs helm dependency update
func UpdateChartDependencies(o UpdateChartDependenciesOptions) error {
	cmd := []string{"helm", "dependency", "update", o.Path}
	err := PrintExec(cmd, o.Print)

	return err
}

// UpgradeReleaseOptions are options passed to UpgradeRelease
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
func UpgradeRelease(o UpgradeReleaseOptions) error {
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
	err := PrintExec(cmd, o.Print)

	return err
}

// DeleteReleaseOptions are options passed to DeleteRelease
type DeleteReleaseOptions struct {
	ReleaseName  string
	KubeContext  string
	TLS          bool
	HelmTLSStore string
	Timeout      int
	Print        bool
}

// DeleteRelease deletes a release from Kubernetes
func DeleteRelease(o DeleteReleaseOptions) error {
	cmd := []string{
		"helm", "delete", o.ReleaseName, "--purge",
		"--timeout", fmt.Sprintf("%d", o.Timeout),
	}
	if o.KubeContext != "" {
		cmd = append(cmd, "--kube-context", o.KubeContext)
	}
	cmd = append(cmd, getTLS(o.TLS, o.KubeContext, o.HelmTLSStore)...)
	err := PrintExec(cmd, o.Print)

	return err
}

// createValuesChain will create a chain of values files to use
func createValuesChain(name, dir string, packedValues []string) []string {
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

// createSetChain will create a chain of sets to use
func createSetChain(name string, inputSet []string) []string {
	set := []string{"--set", fmt.Sprintf("fullnameOverride=%s", name)}
	for _, s := range inputSet {
		set = append(set, "--set", s)
	}
	return set
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
