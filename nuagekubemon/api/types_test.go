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

package api

import (
	"github.com/google/gofuzz"
	metav1 "k8s.io/kubernetes/pkg/api/meta"
	apitesting "k8s.io/kubernetes/pkg/api/testing"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"
	"math/rand"
	"testing"
)

var _ runtime.Object = &NuageNetworkPolicy{}
var _ metav1.ObjectMetaAccessor = &NuageNetworkPolicy{}

var _ runtime.Object = &NuageNetworkPolicyList{}
var _ unversioned.ListMetaAccessor = &NuageNetworkPolicyList{}

func exampleFuzzerFuncs() []interface{} {
	return []interface{}{
		func(obj *NuageNetworkPolicyList, c fuzz.Continue) {
			c.FuzzNoCustom(obj)
			obj.Items = make([]NuageNetworkPolicy, c.Intn(10))
			for i := range obj.Items {
				c.Fuzz(&obj.Items[i])
			}
		},
	}
}

// TestRoundTrip tests that the third-party kinds can be marshaled and unmarshaled correctly to/from JSON
// without the loss of information. Moreover, deep copy is tested.
func TestRoundTrip(t *testing.T) {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)

	AddToScheme(scheme)

	seed := rand.Int63()
	//fuzzerFuncs := apitesting.MergeFuzzerFuncs(t, apitesting.GenericFuzzerFuncs(t, codecs), exampleFuzzerFuncs())
	fuzzer := apitesting.FuzzerFor(t, SchemeGroupVersion, rand.NewSource(seed))

	apitesting.RoundTripSpecificKindWithoutProtobuf(t, SchemeGroupVersion.WithKind("NuageNetworkPolicy"), scheme, codecs, fuzzer, nil)
	apitesting.RoundTripSpecificKindWithoutProtobuf(t, SchemeGroupVersion.WithKind("NuageNetworkPolicyList"), scheme, codecs, fuzzer, nil)
}
