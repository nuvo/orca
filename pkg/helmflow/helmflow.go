package helmflow

import (
	"os"

	genutils "orca/pkg/utils/general"
	helmutils "orca/pkg/utils/helm"
)

// DeployChartFromMuseum deploys a Helm chart from a chart museum
func DeployChartFromMuseum(releaseName, name, version, kubeContext, namespace, museum, helmTLSStore string, tls bool, packedValues, set []string) {
	currDir, _ := os.Getwd()
	tempDir := genutils.MkRandomDir()
	os.Chdir(tempDir)

	if releaseName == "" {
		releaseName = name
	}
	helmutils.AddRepository(museum)
	helmutils.FetchChart(museum, name, version)
	helmutils.UpdateChartDependencies(name)
	valuesChain := helmutils.CreateValuesChain(name, packedValues)
	setChain := helmutils.CreateSetChain(name, set)

	helmutils.UpgradeRelease(name, releaseName, kubeContext, namespace, valuesChain, setChain, tls, helmTLSStore)

	os.Chdir(currDir)
	os.RemoveAll(tempDir)
}
