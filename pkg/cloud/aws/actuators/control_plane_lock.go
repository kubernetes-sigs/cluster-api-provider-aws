/*
Copyright 2019 The Kubernetes Authors.

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

package actuators

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// ControlPlaneConfigMapName returns the name of the ConfigMap used to coordinate the bootstrapping of control plane
// nodes.
func ControlPlaneConfigMapName(cluster *v1alpha1.Cluster) string {
	return fmt.Sprintf("%s-controlplane", cluster.UID)
}

// ListOptionsForCluster returns a ListOptions with a label selector for clusterName.
func ListOptionsForCluster(clusterName string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", v1alpha1.MachineClusterLabelName, clusterName),
	}
}
