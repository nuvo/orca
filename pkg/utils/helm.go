package utils

import (
	"fmt"
	"os"
	"strings"
)

// GetInstalledReleases gets the installed Helm releases in a given namespace
func GetInstalledReleases(kubeContext, namespace, helmTLSStore string, tls, onlyManaged bool) []ReleaseSpec {

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

			if words[statusCol] == "FAILED" {
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

// DeployChartFromMuseum deploys a Helm chart from a chart museum
func DeployChartFromMuseum(releaseName, name, version, kubeContext, namespace, museum, helmTLSStore string, tls bool, packedValues, set []string, isIsolated bool) {
	tempDir := MkRandomDir()

	if releaseName == "" {
		releaseName = name
	}
	if isIsolated {
		AddRepository(museum, isIsolated)
		UpdateRepository(museum, isIsolated)
	}
	FetchChart(museum, name, version, tempDir, isIsolated)
	path := fmt.Sprintf("%s/%s", tempDir, name)
	UpdateChartDependencies(path, isIsolated)
	valuesChain := CreateValuesChain(name, tempDir, packedValues)
	setChain := CreateSetChain(name, set)

	UpgradeRelease(name, releaseName, kubeContext, namespace, valuesChain, setChain, tls, helmTLSStore, tempDir, isIsolated)

	os.RemoveAll(tempDir)
}

// PushChartToMuseum packages and pushes a Helm chart to a chart repository
func PushChartToMuseum(path, append, museum string, lint, print bool) {
	newVersion := UpdateChartVersion(path, append)
	if lint {
		Lint(path, print)
	}
	AddRepository(museum, print)
	UpdateChartDependencies(path, print)
	PushChart(museum, path, print)
	fmt.Println(newVersion)
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
func AddRepository(museum string, print bool) {
	museumSplit := strings.Split(museum, "=")
	museumName := museumSplit[0]
	museumURL := museumSplit[1]

	cmd := fmt.Sprintf("helm repo add %s %s", museumName, museumURL)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// UpdateRepository updates helm repositories
func UpdateRepository(museum string, print bool) {
	cmd := fmt.Sprintf("helm repo update")
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// FetchChart fetches a chart from museum by name and version and untars it in the local directory
func FetchChart(museum, name, version, dir string, print bool) {
	museumSplit := strings.Split(museum, "=")
	museumName := museumSplit[0]

	cmd := fmt.Sprintf("helm fetch %s/%s --version %s --untar -d %s", museumName, name, version, dir)
	output := Exec(cmd)
	if print {
		fmt.Println(cmd)
		fmt.Print(output)
	}
}

// PushChart pushes a helm chart to a chart repository
func PushChart(museum, path string, print bool) {
	museumSplit := strings.Split(museum, "=")
	museumName := museumSplit[0]

	cmd := fmt.Sprintf("helm push %s %s", path, museumName)
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
