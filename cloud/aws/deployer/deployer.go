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
	"github.com/aws/aws-sdk-go/aws/session"
	awsSSM "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ssm"
	clustercommon "sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
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

	status := &v1alpha1.AWSClusterProviderStatus{}
	if err := codec.DecodeProviderStatus(cluster.Status.ProviderStatus, status); err != nil {
		return "", errors.Wrap(err, "failed to decode cluster provider status in deployer")
	}
	if status.Network.APIServerELB.DNSName == "" {
		return "", errors.New("ELB has no DNS name")
	}
	return status.Network.APIServerELB.DNSName, nil
}

// GetKubeConfig returns the kubeconfig after the bootstrap process is complete.
func (*AWSDeployer) GetKubeConfig(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	codec, err := v1alpha1.NewCodec()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create codec in deployer")
	}

	sess, err := session.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "unable to create AWS session")
	}

	svc := ssm.NewService(awsSSM.New(sess))
	if err != nil {
		return "", errors.Wrap(err, "unable to create AWS SSM client")
	}

	status := &v1alpha1.AWSClusterProviderStatus{}
	if err := codec.DecodeProviderStatus(cluster.Status.ProviderStatus, status); err != nil {
		return "", errors.Wrap(err, "failed to decode machine provider status in deployer")
	}

	certPem, err := svc.GetParameter(cluster.Name, certificates.SSMCACertificatePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve CA certificate")
	}

	cert, err := certificates.DecodeCertPEM([]byte(certPem))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode CA Cert")
	}

	keyPem, err := svc.GetParameter(cluster.Name, certificates.SSMCAPrivateKeyPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve CA private-key")
	}

	key, err := certificates.DecodePrivateKeyPEM([]byte(keyPem))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode private key")
	}

	cfg, err := certificates.NewKubeconfig(status.Network.APIServerELB.DNSName, cert, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate a kubeconfig")
	}
	yaml, err := clientcmd.Write(*cfg)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize config to yaml")
	}
	return string(yaml), nil
}
