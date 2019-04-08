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

	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"

	"k8s.io/klog"
	kubeadmv1beta1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
)

const (
	// localIPV4lookup looks up the instance's IP through the metadata service.
	// See https://cloudinit.readthedocs.io/en/latest/topics/instancedata.html
	localIPV4Lookup = "{{ ds.meta_data.local_ipv4 }}"

	// HostnameLookup uses the instance metadata service to lookup its own hostname.
	HostnameLookup = "{{ ds.meta_data.hostname }}"

	// ContainerdSocket is the expected path to containerd socket.
	ContainerdSocket = "/var/run/containerd/containerd.sock"

	// APIServerBindPort is the default port for the kube-apiserver to bind to.
	APIServerBindPort = 6443

	// CloudProvider is the name of the cloud provider passed to various
	// kubernetes components.
	CloudProvider = "aws"

	nodeRole = "node-role.kubernetes.io/node="
)

// SetDefaultClusterConfiguration sets default dynamic values without overriding
// user specified values.
func SetDefaultClusterConfiguration(machine *actuators.MachineScope, base *kubeadmv1beta1.ClusterConfiguration) {
	if base == nil {
		base = &kubeadmv1beta1.ClusterConfiguration{}
	}
	s := machine.Scope

	// Only set the control plane endpoint if the user hasn't specified one.
	if base.ControlPlaneEndpoint == "" {
		base.ControlPlaneEndpoint = fmt.Sprintf("%s:%d", s.Network().APIServerELB.DNSName, APIServerBindPort)
	}
	// Add the control plane endpoint to the list of cert SAN
	base.APIServer.CertSANs = append(base.APIServer.CertSANs, localIPV4Lookup, s.Network().APIServerELB.DNSName)
}

// SetClusterConfigurationOverrides will modify the supplied configuration with certain values
// that cluster-api-provider-aws requires overriding user specified input.
func SetClusterConfigurationOverrides(machine *actuators.MachineScope, base *kubeadmv1beta1.ClusterConfiguration) {
	if base == nil {
		base = &kubeadmv1beta1.ClusterConfiguration{}
	}
	s := machine.Scope

	SetDefaultClusterConfiguration(machine, base)

	// cloud-provider for the APIServer must be set to 'aws'.
	if base.APIServer.ControlPlaneComponent.ExtraArgs == nil {
		base.APIServer.ControlPlaneComponent.ExtraArgs = map[string]string{}
	}
	if cp, ok := base.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		klog.Infof("Overriding cloud provider %q with required value %q", cp, CloudProvider)
	}
	base.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"] = CloudProvider

	if base.ControllerManager.ExtraArgs == nil {
		base.ControllerManager.ExtraArgs = map[string]string{}
	}
	if cp, ok := base.ControllerManager.ExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		klog.Infof("Overriding cloud provider %q with required value %q", cp, CloudProvider)
	}
	base.ControllerManager.ExtraArgs["cloud-provider"] = CloudProvider

	// The kubeadm config clustername must match the provided name of the cluster.
	if base.ClusterName != "" && base.ClusterName != s.Name() {
		klog.Infof("Overriding provided cluster name %q with %q. The kubeadm cluster name and cluster-api name must match.", base.ClusterName, s.Name())
	}
	base.ClusterName = s.Name()

	// The networking values provided by the Cluster object must equal the
	// kubeadm networking configuration.
	base.Networking.DNSDomain = s.Cluster.Spec.ClusterNetwork.ServiceDomain
	base.Networking.PodSubnet = s.Cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]
	base.Networking.ServiceSubnet = s.Cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]

	// The kubernetes version that kubeadm is using must be the same as the
	// requested version in the config
	base.KubernetesVersion = machine.Machine.Spec.Versions.ControlPlane
}

// SetInitConfigurationOverrides overrides user input on particular fields for
// the kubeadm InitConfiguration.
func SetInitConfigurationOverrides(base *kubeadmv1beta1.InitConfiguration) {
	if base == nil {
		base = &kubeadmv1beta1.InitConfiguration{}
	}

	if base.NodeRegistration.Name != "" && base.NodeRegistration.Name != HostnameLookup {
		klog.Infof("Overriding NodeRegistration name from %q to %q. The node registration needs to be dynamically generated in aws.", base.NodeRegistration.Name, HostnameLookup)
	}
	base.NodeRegistration.Name = HostnameLookup

	// TODO(chuckha): This may become a default instead of an override.
	if base.NodeRegistration.CRISocket != "" && base.NodeRegistration.CRISocket != ContainerdSocket {
		klog.Infof("Overriding CRISocket from %q to %q. Containerd is only supported container runtime.", base.NodeRegistration.CRISocket, ContainerdSocket)
	}
	base.NodeRegistration.CRISocket = ContainerdSocket

	if base.NodeRegistration.KubeletExtraArgs == nil {
		base.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if cp, ok := base.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		klog.Infof("Overriding node's cloud-provider to the required value of %q.", CloudProvider)
	}
	base.NodeRegistration.KubeletExtraArgs["cloud-provider"] = CloudProvider
}

// joinMachine is a local interface scoping down exactly what SetJoinNodeConfigurationOverrides needs
type joinMachine interface {
	GetScope() *actuators.Scope
	GetMachine() *v1alpha1.Machine
}

// SetJoinNodeConfigurationOverrides overrides user input for certain fields of
// the kubeadm JoinConfiguration during a worker node join.
func SetJoinNodeConfigurationOverrides(caCertHash, bootstrapToken string, machine joinMachine, base *kubeadmv1beta1.JoinConfiguration) kubeadmv1beta1.JoinConfiguration {
	if base == nil {
		base = &kubeadmv1beta1.JoinConfiguration{}
	}
	out := base.DeepCopy()

	if out.Discovery.BootstrapToken == nil {
		out.Discovery.BootstrapToken = &kubeadmv1beta1.BootstrapTokenDiscovery{}
	}
	// TODO: should this actually be the cluster's ContolPlaneEndpoint?
	out.Discovery.BootstrapToken.APIServerEndpoint = fmt.Sprintf("%s:%d", machine.GetScope().Network().APIServerELB.DNSName, APIServerBindPort)
	out.Discovery.BootstrapToken.Token = bootstrapToken
	out.Discovery.BootstrapToken.CACertHashes = append(out.Discovery.BootstrapToken.CACertHashes, caCertHash)

	if out.NodeRegistration.Name != "" && out.NodeRegistration.Name != HostnameLookup {
		klog.Infof("Overriding NodeRegistration name from %q to %q. The node registration needs to be dynamically generated in aws.", out.NodeRegistration.Name, HostnameLookup)
	}
	out.NodeRegistration.Name = HostnameLookup

	// TODO(chuckha): This may become a default instead of an override.
	if out.NodeRegistration.CRISocket != "" && out.NodeRegistration.CRISocket != ContainerdSocket {
		klog.Infof("Overriding CRISocket from %q to %q. Containerd is only supported container runtime.", out.NodeRegistration.CRISocket, ContainerdSocket)
	}
	out.NodeRegistration.CRISocket = ContainerdSocket

	if out.NodeRegistration.KubeletExtraArgs == nil {
		out.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if cp, ok := out.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		klog.Infof("Overriding node's cloud-provider to the required value of %q.", CloudProvider)
	}
	out.NodeRegistration.KubeletExtraArgs["cloud-provider"] = CloudProvider
	if !util.IsControlPlaneMachine(machine.GetMachine()) {
		out.NodeRegistration.KubeletExtraArgs["node-labels"] = nodeRole
	}
	return *out
}

// SetControlPlaneJoinConfigurationOverrides user input for kubeadm join
// configuration during a control plane join action.
func SetControlPlaneJoinConfigurationOverrides(base *kubeadmv1beta1.JoinConfiguration) {
	if base == nil {
		base = &kubeadmv1beta1.JoinConfiguration{}
	}
	if base.ControlPlane == nil {
		base.ControlPlane = &kubeadmv1beta1.JoinControlPlane{}
	}
	base.ControlPlane.LocalAPIEndpoint.AdvertiseAddress = localIPV4Lookup
	base.ControlPlane.LocalAPIEndpoint.BindPort = APIServerBindPort
}
