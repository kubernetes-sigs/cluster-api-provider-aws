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

package kubeadm

import (
	"fmt"

	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
)

const (
	// localIPV4lookup looks up the instance's IP through the metadata service.
	// See https://cloudinit.readthedocs.io/en/latest/topics/instancedata.html
	localIPV4Lookup = "{{ ds.meta_data.local_ipv4 }}"

	// hostname lookup uses the instance metadata service to lookup its own hostname.
	hostnameLookup = "{{ ds.meta_data.hostname }}"

	// containerdSocket is the expected path to containerd socket.
	containerdSocket = "/var/run/containerd/containerd.sock"

	apiServerBindPort = 6443
)

// SetClusterConfigurationOverrides will modify the supplied configuration with certain values
// that cluster-api-provider-aws requires.
func SetClusterConfigurationOverrides(machine *actuators.MachineScope, base *v1beta1.ClusterConfiguration) {
	if base == nil {
		base = &v1beta1.ClusterConfiguration{}
	}
	s := machine.Scope

	// Set the apiserver cloud provider
	base.APIServer.CertSANs = append(base.APIServer.CertSANs, localIPV4Lookup, s.Network().APIServerELB.DNSName)

	if base.APIServer.ControlPlaneComponent.ExtraArgs == nil {
		base.APIServer.ControlPlaneComponent.ExtraArgs = map[string]string{}
	}
	if cloudProvider, ok := base.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"]; ok {
		klog.Infof("Overriding cloud provider %q with 'aws'", cloudProvider)
	}
	base.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"] = "aws"

	controlPlaneEndpoint := fmt.Sprintf("%s:%d", s.Network().APIServerELB.DNSName, 6443)
	if base.ControlPlaneEndpoint != "" {
		klog.Infof("Overriding control plane endpoint with %q", controlPlaneEndpoint)
	}
	base.ControlPlaneEndpoint = controlPlaneEndpoint

	if base.ClusterName != "" {
		klog.Infof("Overriding cluster name with %q", s.Name())
	}
	base.ClusterName = s.Name()

	base.Networking.DNSDomain = s.Cluster.Spec.ClusterNetwork.ServiceDomain
	base.Networking.PodSubnet = s.Cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]
	base.Networking.ServiceSubnet = s.Cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]
	base.KubernetesVersion = machine.Machine.Spec.Versions.ControlPlane
}

func SetInitConfigurationOverrides(base *v1beta1.InitConfiguration) {
	// Override critical variables
	if base == nil {
		base = &v1beta1.InitConfiguration{}
	}

	if base.NodeRegistration.Name != "" {
		klog.Infof("Overriding NodeRegistration name to %q", hostnameLookup)
	}
	base.NodeRegistration.Name = hostnameLookup

	if base.NodeRegistration.CRISocket != "" {
		klog.Infof("Overriding CRISocket to %q", containerdSocket)
	}
	base.NodeRegistration.CRISocket = containerdSocket

	if base.NodeRegistration.KubeletExtraArgs == nil {
		base.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}

	if base.NodeRegistration.KubeletExtraArgs == nil {
		base.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if _, ok := base.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok {
		klog.Infof("Overriding node's cloud-provider to 'aws'")
	}
	base.NodeRegistration.KubeletExtraArgs["cloud-provider"] = "aws"
}

func SetJoinNodeConfigurationOverrides(caCertHash, bootstrapToken string, machine *actuators.MachineScope, base *v1beta1.JoinConfiguration) {
	if base == nil {
		base = &v1beta1.JoinConfiguration{}
	}
	s := machine.Scope

	if base.Discovery.BootstrapToken == nil {
		base.Discovery.BootstrapToken = &v1beta1.BootstrapTokenDiscovery{}
	}
	base.Discovery.BootstrapToken.Token = bootstrapToken
	base.Discovery.BootstrapToken.APIServerEndpoint = s.Network().APIServerELB.DNSName
	base.Discovery.BootstrapToken.CACertHashes = append(base.Discovery.BootstrapToken.CACertHashes, caCertHash)
	if base.NodeRegistration.Name != "" {
		klog.Infof("Overriding NodeRegistration name to %q", hostnameLookup)
	}
	base.NodeRegistration.Name = hostnameLookup

	if base.NodeRegistration.CRISocket != "" {
		klog.Infof("Overriding CRISocket to %q", containerdSocket)
	}
	base.NodeRegistration.CRISocket = containerdSocket

	if base.NodeRegistration.KubeletExtraArgs == nil {
		base.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if _, ok := base.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok {
		klog.Infof("Overriding node's cloud-provider to 'aws'")
	}
	base.NodeRegistration.KubeletExtraArgs["cloud-provider"] = "aws"
}

func SetControlPlaneJoinConfigurationOverrides(base *v1beta1.JoinConfiguration) {
	if base == nil {
		base = &v1beta1.JoinConfiguration{}
	}
	base.ControlPlane.LocalAPIEndpoint.AdvertiseAddress = localIPV4Lookup
	base.ControlPlane.LocalAPIEndpoint.BindPort = apiServerBindPort
}

func DefaultClusterConfiguration(s *actuators.MachineScope) *v1beta1.ClusterConfiguration {
	return &v1beta1.ClusterConfiguration{
		APIServer:            v1beta1.APIServer{},
		ControlPlaneEndpoint: "",
		ClusterName:          "",
		Networking:           v1beta1.Networking{},
		KubernetesVersion:    "",
	}
}
