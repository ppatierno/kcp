/*
Copyright 2022 The KCP Authors.

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

package clusterworkspace

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
)

func TestReconcileMetadata(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    *tenancyv1alpha1.ClusterWorkspace
		expected metav1.ObjectMeta
	}{
		{
			name: "adds entirely missing labels",
			input: &tenancyv1alpha1.ClusterWorkspace{
				Status: tenancyv1alpha1.ClusterWorkspaceStatus{
					Phase: tenancyv1alpha1.ClusterWorkspacePhaseReady,
					Initializers: []tenancyv1alpha1.ClusterWorkspaceInitializer{
						"pluto", "venus", "apollo",
					},
				},
			},
			expected: metav1.ObjectMeta{
				Labels: map[string]string{
					"internal.kcp.dev/phase": "Ready",
					"initializer.internal.kcp.dev/2eadcbf778956517ec99fd1c1c32a9b13c": "2eadcbf778956517ec99fd1c1c32a9b13cbae759770fc37c341c7fe8",
					"initializer.internal.kcp.dev/aceeb26461953562d30366db65b200f642": "aceeb26461953562d30366db65b200f64241f9e5fe888892d52eea5c",
					"initializer.internal.kcp.dev/ccf53a4988ae8515ee77131ef507cabaf1": "ccf53a4988ae8515ee77131ef507cabaf18822766c2a4cff33b24eb8",
				},
			},
		},
		{
			name: "adds partially missing labels",
			input: &tenancyv1alpha1.ClusterWorkspace{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"internal.kcp.dev/phase": "Ready",
						"initializer.internal.kcp.dev/2eadcbf778956517ec99fd1c1c32a9b13c": "2eadcbf778956517ec99fd1c1c32a9b13cbae759770fc37c341c7fe8",
						"initializer.internal.kcp.dev/aceeb26461953562d30366db65b200f642": "aceeb26461953562d30366db65b200f64241f9e5fe888892d52eea5c",
					},
				},
				Status: tenancyv1alpha1.ClusterWorkspaceStatus{
					Phase: tenancyv1alpha1.ClusterWorkspacePhaseReady,
					Initializers: []tenancyv1alpha1.ClusterWorkspaceInitializer{
						"pluto", "venus", "apollo",
					},
				},
			},
			expected: metav1.ObjectMeta{
				Labels: map[string]string{
					"internal.kcp.dev/phase": "Ready",
					"initializer.internal.kcp.dev/2eadcbf778956517ec99fd1c1c32a9b13c": "2eadcbf778956517ec99fd1c1c32a9b13cbae759770fc37c341c7fe8",
					"initializer.internal.kcp.dev/aceeb26461953562d30366db65b200f642": "aceeb26461953562d30366db65b200f64241f9e5fe888892d52eea5c",
					"initializer.internal.kcp.dev/ccf53a4988ae8515ee77131ef507cabaf1": "ccf53a4988ae8515ee77131ef507cabaf18822766c2a4cff33b24eb8",
				},
			},
		},
		{
			name: "removes previously-needed labels removed on mutation that removes initializer",
			input: &tenancyv1alpha1.ClusterWorkspace{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"internal.kcp.dev/phase": "Ready",
						"initializer.internal.kcp.dev/2eadcbf778956517ec99fd1c1c32a9b13c": "2eadcbf778956517ec99fd1c1c32a9b13cbae759770fc37c341c7fe8",
						"initializer.internal.kcp.dev/aceeb26461953562d30366db65b200f642": "aceeb26461953562d30366db65b200f64241f9e5fe888892d52eea5c",
					},
				},
				Status: tenancyv1alpha1.ClusterWorkspaceStatus{
					Phase: tenancyv1alpha1.ClusterWorkspacePhaseReady,
					Initializers: []tenancyv1alpha1.ClusterWorkspaceInitializer{
						"pluto",
					},
				},
			},
			expected: metav1.ObjectMeta{
				Labels: map[string]string{
					"internal.kcp.dev/phase": "Ready",
					"initializer.internal.kcp.dev/2eadcbf778956517ec99fd1c1c32a9b13c": "2eadcbf778956517ec99fd1c1c32a9b13cbae759770fc37c341c7fe8",
				},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			reconcileMetadata(testCase.input)
			if diff := cmp.Diff(testCase.input.ObjectMeta, testCase.expected); diff != "" {
				t.Errorf("invalid output after reconciling metadata: %v", diff)
			}
		})
	}
}
