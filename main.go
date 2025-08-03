package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Initialize the Kubernetes client with the specified kubeconfig and context, exit on failure.
	client, err = createClient(*kubeconfig, *contextName)
	if err != nil {
		fmt.Printf("unable to create kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Get and print the Kubernetes server version from the cluster.
	version, err := client.Discovery().ServerVersion()
	if err != nil {
		fmt.Printf("unable to determine Kubernetes version: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Kubernetes version: %s\n", version)

	// List and print all Kubernetes namespaces in the current cluster context.
	ctx := context.Background()
	namespaces, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("unable to get namespaces: %v\n", err)
		os.Exit(1)
	}
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
}
