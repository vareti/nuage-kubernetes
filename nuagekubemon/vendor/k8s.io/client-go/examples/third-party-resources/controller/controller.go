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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	tprv1 "k8s.io/client-go/examples/third-party-resources/apis/tpr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type NuageNetworkPolicyController struct {
	NuageNetworkPolicyClient *rest.RESTClient
	NuageNetworkPolicyScheme *runtime.Scheme
}

// Run starts an NuageNetworkPolicy resource controller
func (c *NuageNetworkPolicyController) Run(ctx context.Context) error {
	fmt.Print("Watch NuageNetworkPolicy objects\n")

	// Watch NuageNetworkPolicy objects
	_, err := c.watchNuageNetworkPolicies(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for NuageNetworkPolicy resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *NuageNetworkPolicyController) watchNuageNetworkPolicies(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.NuageNetworkPolicyClient,
		tprv1.NuageNetworkPolicyResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&tprv1.NuageNetworkPolicy{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *NuageNetworkPolicyController) onAdd(obj interface{}) {
	policy := obj.(*tprv1.NuageNetworkPolicy)
	fmt.Printf("[CONTROLLER] OnAdd %s\n", policy.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.NuageNetworkPolicyScheme.Copy(policy)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}
	policyCopy := copyObj.(*tprv1.NuageNetworkPolicy)

	err = c.NuageNetworkPolicyClient.Put().
		Name(policy.ObjectMeta.Name).
		Namespace(policy.ObjectMeta.Namespace).
		Resource(tprv1.NuageNetworkPolicyResourcePlural).
		Body(policyCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", policyCopy)
	}
}

func (c *NuageNetworkPolicyController) onUpdate(oldObj, newObj interface{}) {
	oldNuageNetworkPolicy := oldObj.(*tprv1.NuageNetworkPolicy)
	newNuageNetworkPolicy := newObj.(*tprv1.NuageNetworkPolicy)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %s\n", oldNuageNetworkPolicy.ObjectMeta.SelfLink)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %s\n", newNuageNetworkPolicy.ObjectMeta.SelfLink)
}

func (c *NuageNetworkPolicyController) onDelete(obj interface{}) {
	policy := obj.(*tprv1.NuageNetworkPolicy)
	fmt.Printf("[CONTROLLER] OnDelete %s\n", policy.ObjectMeta.SelfLink)
}
