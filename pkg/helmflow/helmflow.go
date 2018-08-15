package helmflow

import (
	"os"

	"orca/pkg/helm"
	genutils "orca/pkg/utils/general"
)

// DeployChartFromMuseum deploys a Helm chart from a chart museum
func DeployChartFromMuseum(releaseName, name, version, kubeContext, namespace, museum, helmTLSStore string, tls bool, packedValues, set []string) {
	currDir, _ := os.Getwd()
	tempDir := genutils.MkRandomDir()
	os.Chdir(tempDir)

	if releaseName == "" {
		releaseName = name
	}
	helm.AddRepository(museum)
	helm.FetchChart(museum, name, version)
	helm.UpdateChartDependencies(name)
	valuesChain := helm.CreateValuesChain(name, packedValues)
	setChain := helm.CreateSetChain(name, set)

	helm.UpgradeRelease(name, releaseName, kubeContext, namespace, valuesChain, setChain, tls, helmTLSStore)

	os.Chdir(currDir)
	os.RemoveAll(tempDir)
}
