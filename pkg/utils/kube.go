package utils

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateNamespace creates a namespace
func CreateNamespace(name, kubeContext string) {

	clientset := getClientSet(kubeContext)

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_, err := clientset.Core().Namespaces().Create(nsSpec)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("created namespace:", name)
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(name, kubeContext string) {
	clientset := getClientSet(kubeContext)

	deleteOptions := &metav1.DeleteOptions{}
	err := clientset.Core().Namespaces().Delete(name, deleteOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("deleted namespace:", name)
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func getClientSet(kubeContext string) *kubernetes.Clientset {
	var kubeconfig string
	if kubeConfigPath := os.Getenv("KUBECONFIG"); kubeConfigPath != "" {
		kubeconfig = kubeConfigPath // CI process
	} else {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config") // Development environment
	}

	// use the current context in kubeconfig
	config, err := buildConfigFromFlags(kubeContext, kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return clientset
}
