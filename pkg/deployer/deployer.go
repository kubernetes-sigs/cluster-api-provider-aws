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

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// Deployer satisfies the ProviderDeployer(https://github.com/kubernetes-sigs/cluster-api/blob/master/cmd/clusterctl/clusterdeployer/clusterdeployer.go) interface.
type Deployer struct {
	scopeGetter actuators.ScopeGetter
}

// Params is used to create a new deployer.
type Params struct {
	ScopeGetter actuators.ScopeGetter
}

// New returns a new Deployer.
func New(params Params) *Deployer {
	return &Deployer{
		scopeGetter: params.ScopeGetter,
	}
}

// GetIP returns the IP of a machine, but this is going away.
func (d *Deployer) GetIP(cluster *clusterv1.Cluster, _ *clusterv1.Machine) (string, error) {
	scope, err := d.scopeGetter.GetScope(actuators.ScopeParams{Cluster: cluster})
	if err != nil {
		return "", err
	}

	if scope.ClusterStatus != nil && scope.ClusterStatus.Network.APIServerELB.DNSName != "" {
		return scope.ClusterStatus.Network.APIServerELB.DNSName, nil
	}

	elbsvc := elb.NewService(scope)
	return elbsvc.GetAPIServerDNSName()
}

// GetKubeConfig returns the kubeconfig after the bootstrap process is complete.
func (d *Deployer) GetKubeConfig(cluster *clusterv1.Cluster, _ *clusterv1.Machine) (string, error) {

	// Load provider config.
	config, err := providerv1.ClusterConfigFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return "", errors.Errorf("failed to load cluster provider status: %v", err)
	}

	cert, err := certificates.DecodeCertPEM(config.CAKeyPair.Cert)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode CA Cert")
	} else if cert == nil {
		return "", errors.New("certificate not found in config")
	}

	key, err := certificates.DecodePrivateKeyPEM(config.CAKeyPair.Key)
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

	cfg, err := certificates.NewKubeconfig(cluster.Name, server, cert, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate a kubeconfig")
	}

	yaml, err := clientcmd.Write(*cfg)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize config to yaml")
	}

	return string(yaml), nil
}
