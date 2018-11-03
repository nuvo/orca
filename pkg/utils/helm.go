package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

// DeployChartsFromRepository deploys a list of Helm charts from a repository in parallel
func DeployChartsFromRepository(releasesToInstall []ReleaseSpec, kubeContext, namespace, repo, helmTLSStore string, tls bool, packedValues, set []string, inject bool, parallel, timeout int) {

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
				DeployChartFromRepository(c.ReleaseName, c.ChartName, c.ChartVersion, kubeContext, namespace, repo, helmTLSStore, tls, packedValues, set, false, inject, timeout)
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

// DeployChartFromRepository deploys a Helm chart from a chart repository
func DeployChartFromRepository(releaseName, name, version, kubeContext, namespace, repo, helmTLSStore string, tls bool, packedValues, set []string, isIsolated, inject bool, timeout int) {
	tempDir := MkRandomDir()

	if releaseName == "" {
		releaseName = name
	}
	if isIsolated {
		AddRepository(repo, isIsolated)
		UpdateRepositories(isIsolated)
	}
	FetchChart(repo, name, version, tempDir, isIsolated)
	path := fmt.Sprintf("%s/%s", tempDir, name)
	UpdateChartDependencies(path, isIsolated)
	valuesChain := CreateValuesChain(name, tempDir, packedValues)
	setChain := CreateSetChain(name, set)

	UpgradeRelease(name, releaseName, kubeContext, namespace, valuesChain, setChain, tls, helmTLSStore, tempDir, isIsolated, inject, timeout)

	os.RemoveAll(tempDir)
}

// Lint takes a path to a chart and runs a series of tests to verify that the chart is well-formed
func Lint(path string, print bool) {
	cmd := []string{"helm", "lint", path}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// AddRepository adds a chart repository to the repositories file
func AddRepository(repo string, print bool) {
	repoName, repoURL := SplitInTwo(repo, "=")

	cmd := []string{
		"helm", "repo",
		"add", repoName, repoURL,
	}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// UpdateRepositories updates helm repositories
func UpdateRepositories(print bool) {
	cmd := []string{"helm", "repo", "update"}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// FetchChart fetches a chart from chart repository by name and version and untars it in the local directory
func FetchChart(repo, name, version, dir string, print bool) {
	repoName, _ := SplitInTwo(repo, "=")

	cmd := []string{
		"helm", "fetch",
		fmt.Sprintf("%s/%s", repoName, name),
		"--version", version,
		"--untar",
		"-d", dir,
	}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// PushChartToRepository packages and pushes a Helm chart to a chart repository
func PushChartToRepository(path, append, repo string, lint, print bool) {
	newVersion := UpdateChartVersion(path, append)
	if lint {
		Lint(path, print)
	}
	AddRepository(repo, print)
	UpdateChartDependencies(path, print)
	PushChart(repo, path, print)
	fmt.Println(newVersion)
}

// PushChart pushes a helm chart to a chart repository
func PushChart(repo, path string, print bool) {
	repoName, _ := SplitInTwo(repo, "=")

	cmd := []string{"helm", "push", path, repoName}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// UpdateChartDependencies performs helm dependency update
func UpdateChartDependencies(path string, print bool) {
	cmd := []string{"helm", "dependency", "update", path}
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
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

// UpgradeRelease performs helm upgrade -i
func UpgradeRelease(name, releaseName, kubeContext, namespace string, values, set []string, tls bool, helmTLSStore, dir string, print, inject bool, timeout int) {
	cmd := []string{"helm"}
	kubeContextFlag := "--kube-context"
	if inject {
		kubeContextFlag = "--kubecontext"
		cmd = append(cmd, "inject")
	}
	cmd = append(cmd, "upgrade", "-i", releaseName, fmt.Sprintf("%s/%s", dir, name))
	if kubeContext != "" {
		cmd = append(cmd, kubeContextFlag, kubeContext)
	}
	if namespace != "" {
		cmd = append(cmd, "--namespace", namespace)
	}
	cmd = append(cmd, values...)
	cmd = append(cmd, set...)
	cmd = append(cmd, "--timeout", fmt.Sprintf("%d", timeout))
	cmd = append(cmd, getTLS(tls, kubeContext, helmTLSStore)...)

	// cmd := fmt.Sprintf("helm %supgrade%s -i %s --%s %s --namespace %s%s%s %s/%s --timeout %d", injectStr, getTLS(tls, kubeContext, helmTLSStore), releaseName, kubeContextFlag, kubeContext, namespace, values, set, dir, name, timeout)
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
}

// DeleteReleases deletes a list of releases in parallel
func DeleteReleases(releasesToDelete []ReleaseSpec, kubeContext, helmTLSStore string, tls bool, parallel, timeout int) {
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
			DeleteRelease(c.ReleaseName, kubeContext, tls, helmTLSStore, timeout, print)
			log.Println("deleted", c.ReleaseName)
		}(c)
	}
	bwg.Wait()
}

// DeleteRelease deletes a release from Kubernetes
func DeleteRelease(releaseName, kubeContext string, tls bool, helmTLSStore string, timeout int, print bool) {
	cmd := []string{
		"helm", "delete", releaseName, "--purge",
		"--timeout", fmt.Sprintf("%d", timeout),
	}
	if kubeContext != "" {
		cmd = append(cmd, "--kube-context", kubeContext)
	}
	cmd = append(cmd, getTLS(tls, kubeContext, helmTLSStore)...)
	if print {
		fmt.Println(cmd)
	}
	output := Exec(cmd)
	if print {
		fmt.Print(output)
	}
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
