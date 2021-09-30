/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	t.Run("for Cluster", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:                &v1beta1.Cluster{},
		Spoke:              &Cluster{},
		SpokeAfterMutation: clusterSpokeAfterMutation,
		FuzzerFuncs:        []fuzzer.FuzzerFuncs{ClusterFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for Machine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.Machine{},
		Spoke:       &Machine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineFuzzFunc, BootstrapFuzzFuncs, MachineStatusFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for MachineSet", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineSet{},
		Spoke:       &MachineSet{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineSetFuzzFunc, BootstrapFuzzFuncs, CustomObjectMetaFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for MachineDeployment", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineDeployment{},
		Spoke:       &MachineDeployment{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineDeploymentFuzzFunc, BootstrapFuzzFuncs, CustomObjectMetaFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for MachineHealthCheck", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineHealthCheck{},
		Spoke:       &MachineHealthCheck{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineHealthCheckFuzzFunc, ObjectReferenceFuzzFunc},
	}))
}

func MachineStatusFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		MachineStatusFuzzer,
	}
}

func MachineStatusFuzzer(in *MachineStatus, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// These fields have been removed in v1beta1
	// data is going to be lost, so we're forcing zero values to avoid round trip errors.
	in.Version = nil
}

func CustomObjectMetaFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		CustomObjectMetaFuzzer,
	}
}

func CustomObjectMetaFuzzer(in *ObjectMeta, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// These fields have been removed in v1alpha4
	// data is going to be lost, so we're forcing zero values here.
	in.Name = ""
	in.GenerateName = ""
	in.Namespace = ""
	in.OwnerReferences = nil
}

func BootstrapFuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		BootstrapFuzzer,
	}
}

func BootstrapFuzzer(obj *Bootstrap, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Bootstrap.Data has been removed in v1alpha4, so setting it to nil in order to avoid v1alpha3 --> <hub> --> v1alpha3 round trip errors.
	obj.Data = nil
}

// clusterSpokeAfterMutation modifies the spoke version of the Cluster such that it can pass an equality test in the
// spoke-hub-spoke conversion scenario.
func clusterSpokeAfterMutation(c conversion.Convertible) {
	cluster := c.(*Cluster)

	// Create a temporary 0-length slice using the same underlying array as cluster.Status.Conditions to avoid
	// allocations.
	tmp := cluster.Status.Conditions[:0]

	for i := range cluster.Status.Conditions {
		condition := cluster.Status.Conditions[i]

		// Keep everything that is not ControlPlaneInitializedCondition
		if condition.Type != ConditionType(v1beta1.ControlPlaneInitializedCondition) {
			tmp = append(tmp, condition)
		}
	}

	// Point cluster.Status.Conditions and our slice that does not have ControlPlaneInitializedCondition
	cluster.Status.Conditions = tmp
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

func ClusterFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(c *Cluster, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(c)

			setRefNamespace(c.Spec.ControlPlaneRef, c.Namespace)
			setRefNamespace(c.Spec.InfrastructureRef, c.Namespace)
		},
	}
}

func MachineFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(m *Machine, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(m)

			setRefNamespace(&m.Spec.InfrastructureRef, m.Namespace)
			setRefNamespace(m.Spec.Bootstrap.ConfigRef, m.Namespace)
			if m.Status.NodeRef != nil {
				m.Status.NodeRef.Namespace = m.Namespace
				fuzzer.Fuzz(&m.Status.NodeRef.UID)
			}
		},
	}
}

func MachineSetFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(ms *MachineSet, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(ms)

			setRefNamespace(&ms.Spec.Template.Spec.InfrastructureRef, ms.Namespace)
			setRefNamespace(ms.Spec.Template.Spec.Bootstrap.ConfigRef, ms.Namespace)
		},
	}
}

func MachineDeploymentFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(md *MachineDeployment, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(md)

			setRefNamespace(&md.Spec.Template.Spec.InfrastructureRef, md.Namespace)
			setRefNamespace(md.Spec.Template.Spec.Bootstrap.ConfigRef, md.Namespace)
		},
	}
}

func MachineHealthCheckFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(mh *MachineHealthCheck, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(mh)

			setRefNamespace(mh.Spec.RemediationTemplate, mh.Namespace)
		},
	}
}
