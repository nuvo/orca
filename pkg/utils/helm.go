package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// GetInstalledReleases gets the installed Helm releases in a given namespace
func GetInstalledReleases(kubeContext, namespace, helmTLSStore string, tls, onlyManaged, includeFailed bool) []ReleaseSpec {

	const ReleaseNameCol = 0
	const statusCol = 7
	const VersionCol = 8

	var releaseSpecs []ReleaseSpec
	list := List(kubeContext, namespace, helmTLSStore, tls)

	didHeadersRowPass := false
	for _, line := range strings.Split(list, "\n") {
		if strings.HasPrefix(line, "NAME") {
			didHeadersRowPass = true
			continue
		}
		if didHeadersRowPass && strings.Trim(line, " ") != "" {
			if !(onlyManaged && strings.HasPrefix(line, namespace)) {
				continue
			}

			words := strings.Fields(line)

			if words[statusCol] == "FAILED" && !includeFailed {
				continue
			}

			var releaseSpec ReleaseSpec
			releaseSpec.ReleaseName = words[ReleaseNameCol]
			releaseSpec.ChartName = strings.TrimLeft(releaseSpec.ReleaseName, namespace) // Strange behavior when replacing namespace+"-"
			releaseSpec.ChartName = strings.TrimLeft(releaseSpec.ChartName, "-")
			releaseSpec.ChartVersion = strings.TrimLeft(words[VersionCol], releaseSpec.ChartName+"-")

			releaseSpecs = append(releaseSpecs, releaseSpec)
		}
	}

	return releaseSpecs
}

// DeployChartsFromRepository deploys a list of Helm charts from a repository in parallel
func DeployChartsFromRepository(releasesToInstall []ReleaseSpec, kubeContext, namespace, repo, helmTLSStore string, tls bool, packedValues, set []string) {

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for len(releasesToInstall) > 0 {

		mutex.Lock()
		for _, c := range releasesToInstall {

			wg.Add(1)
			go func(c ReleaseSpec) {
				defer wg.Done()

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
				DeployChartFromRepository(c.ReleaseName, c.ChartName, c.ChartVersion, kubeContext, namespace, repo, helmTLSStore, tls, packedValues, set, false)
				log.Println("deployed chart", c.ChartName, "version", c.ChartVersion)

				// Deployment is done, remove chart from dependencies
				mutex.Lock()
				releasesToInstall = RemoveChartFromDependencies(releasesToInstall, c.ChartName)
				mutex.Unlock()

			}(c)
		}
		mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
	wg.Wait()
}

// DeployChartFromRepository deploys a Helm chart from a chart repository
func DeployChartFromRepository(releaseName, name, version, kubeContext, namespace, repo, helmTLSStore string, tls bool, packedValues, set []string, isIsolated bool) {
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

	UpgradeRelease(name, releaseName, kubeContext, namespace, valuesChain, setChain, tls, helmTLSStore, tempDir, isIsolated)

	os.RemoveAll(tempDir)
}

// List get a list of installed releases in a given namespace
func List(kubeContext, namespace, helmTLSStore string, tls bool) string {
	cmd := fmt.Sprintf("helm ls%s --kube-context %s --namespace %s", getTLS(tls, kubeContext, helmTLSStore), kubeContext, namespace)
	list := Exec(cmd)

	return list
}

// Lint takes a path to a chart and runs a series of tests to verify that the chart is well-formed
func Lint(path string, print bool) {
	cmd := fmt.Sprintf("helm lint %s", path)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// AddRepository adds a chart repository to the repositories file
func AddRepository(repo string, print bool) {
	repoName, repoURL := SplitInTwo(repo, "=")

	cmd := fmt.Sprintf("helm repo add %s %s", repoName, repoURL)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// UpdateRepositories updates helm repositories
func UpdateRepositories(print bool) {
	cmd := fmt.Sprintf("helm repo update")
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// FetchChart fetches a chart from chart repository by name and version and untars it in the local directory
func FetchChart(repo, name, version, dir string, print bool) {
	repoName, _ := SplitInTwo(repo, "=")

	cmd := fmt.Sprintf("helm fetch %s/%s --version %s --untar -d %s", repoName, name, version, dir)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
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

	cmd := fmt.Sprintf("helm push %s %s", path, repoName)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// UpdateChartDependencies performs helm dependency update
func UpdateChartDependencies(path string, print bool) {
	cmd := fmt.Sprintf("helm dependency update %s", path)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// CreateValuesChain will create a chain of values files to use
func CreateValuesChain(name, dir string, packedValues []string) string {
	values := " "
	format := "%s/%s/%s"
	fileToTest := fmt.Sprintf(format, dir, name, "values.yaml")
	if _, err := os.Stat(fileToTest); err == nil {
		values = values + fmt.Sprintf("-f %s", fileToTest)
	}

	for _, v := range packedValues {
		fileToTest = fmt.Sprintf(format, dir, name, v)
		if _, err := os.Stat(fileToTest); err == nil {
			if !strings.Contains(values, " "+fileToTest) {
				values = values + fmt.Sprintf(" -f %s", fileToTest)
			}
		}
	}

	return values
}

// CreateSetChain will create a chain of sets to use
func CreateSetChain(name string, inputSet []string) string {
	set := fmt.Sprintf(" --set fullnameOverride=%s", name)

	for _, s := range inputSet {
		set = set + fmt.Sprintf(" --set %s", s)
	}

	return set
}

// UpgradeRelease performs helm upgrade -i
func UpgradeRelease(name, releaseName, kubeContext, namespace, values, set string, tls bool, helmTLSStore, dir string, print bool) {
	cmd := fmt.Sprintf("helm upgrade%s -i %s --kube-context %s --namespace %s%s%s %s/%s", getTLS(tls, kubeContext, helmTLSStore), releaseName, kubeContext, namespace, values, set, dir, name)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// DeleteReleases deletes a list of releases in parallel
func DeleteReleases(releasesToDelete []ReleaseSpec, kubeContext, helmTLSStore string, tls bool) {
	var wg sync.WaitGroup

	for _, c := range releasesToDelete {
		wg.Add(1)
		go func(c ReleaseSpec) {
			defer wg.Done()
			log.Println("deleting", c.ReleaseName)
			DeleteRelease(c.ReleaseName, kubeContext, tls, helmTLSStore, false)
			log.Println("deleted", c.ReleaseName)
		}(c)
	}
	wg.Wait()
}

// DeleteRelease deletes a release from Kubernetes
func DeleteRelease(releaseName, kubeContext string, tls bool, helmTLSStore string, print bool) {
	cmd := fmt.Sprintf("helm delete%s --purge %s --kube-context %s", getTLS(tls, kubeContext, helmTLSStore), releaseName, kubeContext)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

func getTLS(tls bool, kubeContext, helmTLSStore string) string {
	tlsStr := ""
	if tls == true {
		tlsStr = fmt.Sprintf(" --tls --tls-cert %s/%s.cert.pem --tls-key %s/%s.key.pem", helmTLSStore, kubeContext, helmTLSStore, kubeContext)
	}
	return tlsStr
}
