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
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	clustercommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

// ProviderName is the name of the cloud provider
const ProviderName = "aws"

func init() {
	clustercommon.RegisterClusterProvisioner(ProviderName, &AWSDeployer{})
}

// AWSDeployer implements the cluster-api Deployer interface.
type AWSDeployer struct{}

// GetIP returns the IP of a machine, but this is going away.
func (*AWSDeployer) GetIP(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create codec in deployer")
	}

	status := &v1alpha1.AWSMachineProviderStatus{}
	if err := codec.DecodeProviderStatus(machine.Status.ProviderStatus, status); err != nil {
		return "", errors.Wrap(err, "failed to decode machine provider status in deployer")
	}

	// will have credentials in environment
	sess := session.Must(session.NewSession())
	ec2client := ec2.New(sess)

	dio, err := ec2client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{status.InstanceID},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get instance description in deployer")
	}

	if len(dio.Reservations) == 0 {
		return "", errors.New(fmt.Sprintf("no reservartions found for instance id %v", *status.InstanceID))
	}
	if len(dio.Reservations[0].Instances) == 0 {
		return "", errors.New(fmt.Sprintf("instance was not found in a reserver for instance id %v", *status.InstanceID))
	}
	return *dio.Reservations[0].Instances[0].PublicIpAddress, nil
}

// GetKubeConfig returns the kubeconfig after the bootstrap process is complete.
func (*AWSDeployer) GetKubeConfig(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	return "", errors.New("Not implemented")
}
