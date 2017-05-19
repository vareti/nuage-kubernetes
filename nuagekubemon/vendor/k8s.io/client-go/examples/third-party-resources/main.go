/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	tprv1 "k8s.io/client-go/examples/third-party-resources/apis/tpr/v1"
	exampleclient "k8s.io/client-go/examples/third-party-resources/client"
	examplecontroller "k8s.io/client-go/examples/third-party-resources/controller"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize third party resource if it does not exist
	err = exampleclient.CreateTPR(clientset)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}

	// make a new config for our extension's API group, using the first config as a baseline
	exampleClient, exampleScheme, err := exampleclient.NewClient(config)
	if err != nil {
		panic(err)
	}

	// wait until TPR gets processed
	err = exampleclient.WaitForNuageNetworkPolicyResource(exampleClient)
	if err != nil {
		panic(err)
	}

	// start a controller on instances of our TPR
	controller := examplecontroller.NuageNetworkPolicyController{
		NuageNetworkPolicyClient: exampleClient,
		NuageNetworkPolicyScheme: exampleScheme,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	// Create an instance of our TPR
	policy := &tprv1.NuageNetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "allow-tcp-80",
		},
		Spec: tprv1.NuageNetworkPolicySpec{},
	}
	var result tprv1.NuageNetworkPolicy
	err = exampleClient.Post().
		Resource(tprv1.NuageNetworkPolicyResourcePlural).
		Namespace(apiv1.NamespaceDefault).
		Body(policy).
		Do().Into(&result)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	// Poll until Example object is handled by controller and gets status updated to "Processed"
	err = exampleclient.WaitForNuageNetworkPolicyInstanceProcessed(exampleClient, "allow-tcp-80")
	if err != nil {
		panic(err)
	}
	fmt.Print("PROCESSED\n")

	// Fetch a list of our TPRs
	policyList := tprv1.NuageNetworkPolicyList{}
	err = exampleClient.Get().Resource(tprv1.NuageNetworkPolicyResourcePlural).Do().Into(&policyList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", policyList)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
