package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func createClient(kubeconfigPath string, contextName string) (kubernetes.Interface, error) {
	var kubeconfig *rest.Config

	if kubeconfigPath != "" {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
		configOverrides := &clientcmd.ConfigOverrides{}
		if contextName != "" {
			configOverrides.CurrentContext = contextName
		}
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err := clientConfig.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
		}
		kubeconfig = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
		}
		kubeconfig = config
	}
	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	return client, nil
}

func main() {
	var client kubernetes.Interface
	var err error

	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	contextName := flag.String("context", "", "context to use in the kubeconfig file")
	flag.Parse()

	client, err = createClient(*kubeconfig, *contextName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	version, err := client.Discovery().ServerVersion()
	if err != nil {
		fmt.Printf("unable to determine Kubernetes version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Kubernetes version: %s\n", version)
}
