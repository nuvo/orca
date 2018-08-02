package helm

import (
	"fmt"
	"os"
	"strings"

	genutils "orca/pkg/utils/general"
)

func AddRepository(museum string) {
	museumSplit := strings.Split(museum, "=")
	museumName := museumSplit[0]
	museumURL := museumSplit[1]

	cmd := fmt.Sprintf("helm repo add %s %s", museumName, museumURL)
	genutils.Exec(cmd)
}

func FetchChart(museum, name, version string) {
	museumSplit := strings.Split(museum, "=")
	museumName := museumSplit[0]

	cmd := fmt.Sprintf("helm fetch %s/%s --version %s --untar", museumName, name, version)
	genutils.Exec(cmd)
}

func UpdateChartDependencies(name string) {
	currDir, _ := os.Getwd()
	os.Chdir(name)

	cmd := fmt.Sprintf("helm dependency update")
	genutils.Exec(cmd)

	os.Chdir(currDir)
}

func CreateValuesChain(name string, packedValues []string) string {
	currDir, _ := os.Getwd()
	os.Chdir(name)

	values := " "
	if _, err := os.Stat("values.yaml"); err == nil {
		values = values + "-f values.yaml"
	}

	for _, v := range packedValues {
		if _, err := os.Stat(v); err == nil {
			if !strings.Contains(values, " "+v) {
				values = values + " -f " + v
			}
		}
	}

	os.Chdir(currDir)
	return values
}

func CreateSetChain(name string, inputSet []string) string {
	set := fmt.Sprintf(" --set fullnameOverride=%s", name)

	for _, s := range inputSet {
		set = set + fmt.Sprintf(" --set %s", s)
	}

	return set
}

func UpgradeRelease(name, releaseName, kubeContext, namespace, values, set string, tls bool, helmTLSStore string) {

	currDir, _ := os.Getwd()
	fmt.Println(currDir, name)
	os.Chdir(name)

	cmd := fmt.Sprintf("helm upgrade%s -i %s --kube-context %s --namespace %s%s%s .", getTLS(tls, kubeContext, helmTLSStore), releaseName, kubeContext, namespace, values, set)
	fmt.Println(cmd)
	genutils.Exec(cmd)

	os.Chdir(currDir)
}

func getTLS(tls bool, kubeContext, helmTLSStore string) string {
	tlsStr := ""
	if tls == true {
		tlsStr = fmt.Sprintf(" --tls --tls-cert %s/%s.cert.pem --tls-key %s/%s.key.pem", helmTLSStore, kubeContext, helmTLSStore, kubeContext)
	}
	return tlsStr
}
