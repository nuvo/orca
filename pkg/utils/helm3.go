package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	rspb "k8s.io/helm/pkg/proto/hapi/release"
)

// GetInstalledReleasesOptions are options passed to GetInstalledReleases
type GetInstalledReleasesOptions struct {
	KubeContext   string
	Namespace     string
	IncludeFailed bool
}

// GetInstalledReleases gets the installed Helm releases in a given namespace
func GetInstalledReleases(o GetInstalledReleasesOptions) ([]ReleaseSpec, error) {

	tillerNamespace := "kube-system"
	labels := "OWNER=TILLER,STATUS in (DEPLOYED,FAILED)"
	if !o.IncludeFailed {
		labels = strings.Replace(labels, "FAILED", "", -1)
	}
	storage, err := getTillerStorage(o.KubeContext, tillerNamespace)
	if err != nil {
		return nil, err
	}

	var releaseSpecs []ReleaseSpec
	list, err := listReleases(o.KubeContext, o.Namespace, storage, tillerNamespace, labels)
	if err != nil {
		return nil, err
	}

	for _, releaseData := range list {

		if releaseData.status != "DEPLOYED" {
			continue
		}

		var releaseSpec ReleaseSpec
		releaseSpec.ReleaseName = releaseData.name
		releaseSpec.ChartName = releaseData.chart
		releaseSpec.ChartVersion = releaseData.version

		releaseSpecs = append(releaseSpecs, releaseSpec)
	}

	if !o.IncludeFailed {
		return releaseSpecs, nil
	}

	for _, releaseData := range list {
		if releaseData.status != "FAILED" {
			continue
		}

		exists := false
		for _, rs := range releaseSpecs {
			if releaseData.name == rs.ReleaseName {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		var releaseSpec ReleaseSpec
		releaseSpec.ReleaseName = releaseData.name
		releaseSpec.ChartName = releaseData.chart
		releaseSpec.ChartVersion = releaseData.version

		releaseSpecs = append(releaseSpecs, releaseSpec)
	}

	return releaseSpecs, nil
}

func getTillerStorage(kubeContext, tillerNamespace string) (string, error) {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return "", err
	}
	coreV1 := clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "name=tiller",
	}
	pods, err := coreV1.Pods(tillerNamespace).List(listOptions)
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("Found 0 tiller pods")
	}

	storage := "cfgmaps"
	for _, c := range pods.Items[0].Spec.Containers[0].Command {
		if strings.Contains(c, "secret") {
			storage = "secrets"
		}
	}

	return storage, nil
}

type releaseData struct {
	name      string
	revision  int32
	updated   string
	status    string
	chart     string
	version   string
	namespace string
	time      time.Time
}

func listReleases(kubeContext, namespace, storage, tillerNamespace, labels string) ([]releaseData, error) {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return nil, err
	}
	var releasesData []releaseData
	coreV1 := clientset.CoreV1()
	switch storage {
	case "secrets":
		secrets, err := coreV1.Secrets(tillerNamespace).List(metav1.ListOptions{
			LabelSelector: labels,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range secrets.Items {
			releaseData := getReleaseData(namespace, (string)(item.Data["release"]))
			if releaseData == nil {
				continue
			}
			releasesData = append(releasesData, *releaseData)
		}
	case "cfgmaps":
		configMaps, err := coreV1.ConfigMaps(tillerNamespace).List(metav1.ListOptions{
			LabelSelector: labels,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range configMaps.Items {
			releaseData := getReleaseData(namespace, item.Data["release"])
			if releaseData == nil {
				continue
			}
			releasesData = append(releasesData, *releaseData)
		}
	}

	sort.Slice(releasesData[:], func(i, j int) bool {
		return strings.Compare(releasesData[i].chart, releasesData[j].chart) <= 0
	})

	return releasesData, nil
}

func getReleaseData(namespace, itemReleaseData string) *releaseData {

	data, err := decodeRelease(itemReleaseData)

	if err != nil {
		return nil
	}
	
	if namespace != "" && data.Namespace != namespace {
		return nil
	}

	deployTime := time.Unix(data.Info.LastDeployed.Seconds, 0)
	chartMeta := data.GetChart().Metadata
	releaseData := releaseData{
		time:      deployTime,
		name:      data.Name,
		revision:  data.Version,
		updated:   deployTime.Format("Mon Jan _2 15:04:05 2006"),
		status:    data.GetInfo().Status.Code.String(),
		chart:     chartMeta.Name,
		version:   chartMeta.Version,
		namespace: data.Namespace,
	}
	return &releaseData
}

// GetClientToK8s returns a k8s ClientSet
func GetClientToK8s() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if kubeConfigPath := os.Getenv("KUBECONFIG"); kubeConfigPath != "" {
		kubeconfig = kubeConfigPath // CI process
	} else {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config") // Development environment
	}

	var config *rest.Config

	_, err := os.Stat(kubeconfig)
	if err != nil {
		// In cluster configuration
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// Out of cluster configuration
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

var b64 = base64.StdEncoding
var magicGzip = []byte{0x1f, 0x8b, 0x08}

func decodeRelease(data string) (*rspb.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls rspb.Release
	// unmarshal protobuf bytes
	if err := proto.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}
