package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	internalVersioned "github.com/gaoxinge/website-v2-operator/pkg/client/clientset/versioned"
	"github.com/gaoxinge/website-v2-operator/pkg/controller/extensions.example.com/v2/website"
)

func main() {
	config := rest.Config{
		Host: "http://127.0.0.1:8001",
	}

	clientSet, err := kubernetes.NewForConfig(&config)
	if err != nil {
		log.Printf("new client set with error %v\n", err)
		return
	}

	internalClientSet, err := internalVersioned.NewForConfig(&config)
	if err != nil {
		log.Printf("new internal client set with error %v\n", err)
		return
	}

	list, err := internalClientSet.ExtensionsV2().Websites("default").List(context.TODO(), metav1.ListOptions{})
	fmt.Println(list, err)

	controller := website.NewController(clientSet, internalClientSet)

	stop := make(chan struct{})
	defer close(stop)

	err = controller.Run(stop)
	if err != nil {
		log.Printf("run with error %v\n", err)
		return
	}

	select {}
}
