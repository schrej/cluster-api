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

package noderefutil

import (
	"testing"

	. "github.com/onsi/gomega"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestIndexMachineByNodeName(t *testing.T) {
	testCases := []struct {
		name     string
		object   client.Object
		expected []string
	}{
		{
			name:     "when the machine has no NodeRef",
			object:   &clusterv1.Machine{},
			expected: []string{},
		},
		{
			name: "when the machine has valid a NodeRef",
			object: &clusterv1.Machine{
				Status: clusterv1.MachineStatus{
					NodeRef: &clusterv1.PinnedObjectReference{
						Name: "node1",
					},
				},
			},
			expected: []string{"node1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			got := indexMachineByNodeName(tc.object)
			g.Expect(got).To(ConsistOf(tc.expected))
		})
	}
}
