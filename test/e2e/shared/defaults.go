// +build e2e

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

package shared

const (
	DefaultSSHKeyPairName        = "cluster-api-provider-aws-sigs-k8s-io"
	AMIPrefix                    = "capa-ami-ubuntu-18.04-"
	DefaultImageLookupOrg        = "258751437250"
	KubernetesVersion            = "KUBERNETES_VERSION"
	CNIPath                      = "CNI"
	CNIResources                 = "CNI_RESOURCES"
	AwsNodeMachineType           = "AWS_NODE_MACHINE_TYPE"
	AwsAvailabilityZone1         = "AWS_AVAILABILITY_ZONE_1"
	AwsAvailabilityZone2         = "AWS_AVAILABILITY_ZONE_2"
	MultiAzFlavor                = "multi-az"
	LimitAzFlavor                = "limit-az"
	SpotInstancesFlavor          = "spot-instances"
	SSMFlavor                    = "ssm"
	StorageClassFailureZoneLabel = "failure-domain.beta.kubernetes.io/zone"
)
