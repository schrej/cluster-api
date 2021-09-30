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
	"sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	t.Run("for Cluster", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.Cluster{},
		Spoke:       &Cluster{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{ClusterFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for ClusterClass", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.ClusterClass{},
		Spoke:       &ClusterClass{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{ClusterClassFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for Machine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.Machine{},
		Spoke:       &Machine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineFuzzFunc, MachineStatusFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for MachineSet", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineSet{},
		Spoke:       &MachineSet{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineSetFuzzFunc, ObjectReferenceFuzzFunc},
	}))

	t.Run("for MachineDeployment", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineDeployment{},
		Spoke:       &MachineDeployment{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{MachineDeploymentFuzzFunc, ObjectReferenceFuzzFunc},
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

func ClusterClassFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(c *ClusterClass, fuzzer fuzz.Continue) {
			fuzzer.FuzzNoCustom(c)

			setRefNamespace(c.Spec.ControlPlane.Ref, c.Namespace)
			setRefNamespace(c.Spec.Infrastructure.Ref, c.Namespace)
			if c.Spec.ControlPlane.MachineInfrastructure != nil {
				setRefNamespace(c.Spec.ControlPlane.MachineInfrastructure.Ref, c.Namespace)
			}
			for i, w := range c.Spec.Workers.MachineDeployments {
				setRefNamespace(w.Template.Bootstrap.Ref, c.Namespace)
				setRefNamespace(w.Template.Infrastructure.Ref, c.Namespace)
				c.Spec.Workers.MachineDeployments[i] = w
			}
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
