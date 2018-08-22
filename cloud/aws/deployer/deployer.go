/*
Copyright 2018 The Kubernetes Authors.

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

package deployer

import (
	"errors"

	clustercommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const ProviderName = "aws"

func init() {
	clustercommon.RegisterClusterProvisioner(ProviderName, &AWSDeployer{})
}

type AWSDeployer struct{}

func (*AWSDeployer) GetIP(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	return "", errors.New("Not implemented")
}

func (*AWSDeployer) GetKubeConfig(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	return "", errors.New("Not implemented")
}
