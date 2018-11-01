// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package deployer

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// New returns a new deployer.
func New(servicesGetter service.Getter) *Deployer {
	return &Deployer{servicesGetter}
}

// Deployer satisfies the ProviderDeployer(https://github.com/kubernetes-sigs/cluster-api/blob/master/cmd/clusterctl/clusterdeployer/clusterdeployer.go) inteface.
type Deployer struct {
	servicesGetter service.Getter
}

// GetIP returns the IP of a machine, but this is going away.
func (d *Deployer) GetIP(cluster *clusterv1.Cluster, _ *clusterv1.Machine) (string, error) {
	if cluster.Status.ProviderStatus != nil {

		// Load provider status.
		status, err := providerv1.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
		if err != nil {
			return "", errors.Errorf("failed to load cluster provider status: %v", err)
		}

		if status.Network.APIServerELB.DNSName != "" {
			return status.Network.APIServerELB.DNSName, nil
		}
	}

	// Load provider config.
	config, err := providerv1.ClusterConfigFromProviderConfig(cluster.Spec.ProviderConfig)
	if err != nil {
		return "", errors.Errorf("failed to load cluster provider config: %v", err)
	}

	sess := d.servicesGetter.Session(config)
	elb := d.servicesGetter.ELB(sess)
	return elb.GetAPIServerDNSName(cluster.Name)
}

// GetKubeConfig returns the kubeconfig after the bootstrap process is complete.
func (d *Deployer) GetKubeConfig(cluster *clusterv1.Cluster, _ *clusterv1.Machine) (string, error) {

	// Load provider config.
	config, err := providerv1.ClusterConfigFromProviderConfig(cluster.Spec.ProviderConfig)
	if err != nil {
		return "", errors.Errorf("failed to load cluster provider status: %v", err)
	}

	cert, err := certificates.DecodeCertPEM(config.CACertificate)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode CA Cert")
	} else if cert == nil {
		return "", errors.New("certificate not found in status")
	}

	key, err := certificates.DecodePrivateKeyPEM(config.CAPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode private key")
	} else if key == nil {
		return "", errors.New("key not found in status")
	}

	dnsName, err := d.GetIP(cluster, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to get DNS address")
	}

	server := fmt.Sprintf("https://%s:6443", dnsName)

	cfg, err := certificates.NewKubeconfig(server, cert, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate a kubeconfig")
	}

	yaml, err := clientcmd.Write(*cfg)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize config to yaml")
	}

	return string(yaml), nil
}
