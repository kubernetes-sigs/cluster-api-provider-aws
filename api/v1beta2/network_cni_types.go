/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

// CNISpec defines configuration for CNI.
type CNISpec struct {
	// CNIIngressRules specify rules to apply to control plane and worker node security groups.
	// The source for the rule will be set to control plane and worker security group IDs.
	CNIIngressRules CNIIngressRules `json:"cniIngressRules,omitempty"`
}

// CNIIngressRules is a slice of CNIIngressRule.
type CNIIngressRules []CNIIngressRule

// CNIIngressRule defines an AWS ingress rule for CNI requirements.
type CNIIngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`
	FromPort    int64                 `json:"fromPort"`
	ToPort      int64                 `json:"toPort"`
}
