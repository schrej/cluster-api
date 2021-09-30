/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	clusterv1exp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	t.Run("for MachinePool", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &clusterv1exp.MachinePool{},
		Spoke:       &MachinePool{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachinePoolFuzzFunc, ObjectReferenceFuzzFunc},
	}))
}

func ObjectReferenceFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *corev1.ObjectReference, c fuzz.Continue) {
			c.FuzzNoCustom(obj)

			obj.FieldPath = ""
			obj.ResourceVersion = ""
			obj.UID = ""
		},
	}
}

func MachinePoolFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(m *MachinePool, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(m)

			setRefNamespace(&m.Spec.Template.Spec.InfrastructureRef, m.Namespace)
			setRefNamespace(m.Spec.Template.Spec.Bootstrap.ConfigRef, m.Namespace)
		},
	}
}
