package main

import (
	"fmt"

	kube "github.com/leonardodelira/go-lib-http/kubeclient"
)

func main() {
	kubeClient, err := kube.NewClient(
		kube.WithURL("http://localhost:8080"),
		kube.WithTimeOut(5000),
	)
	if err != nil {
		panic(err)
	}
	fmt.Print(kubeClient)
}
