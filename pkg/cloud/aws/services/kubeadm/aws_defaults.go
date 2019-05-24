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
	"strings"

	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"

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

// joinMachine is a local interface to scope down dependencies
type joinMachine interface {
	GetScope() *actuators.Scope
	GetMachine() *v1alpha1.Machine
}

// SetDefaultClusterConfiguration sets default dynamic values without overriding
// user specified values.
func SetDefaultClusterConfiguration(machine joinMachine, base *kubeadmv1beta1.ClusterConfiguration) *kubeadmv1beta1.ClusterConfiguration {
	if base == nil {
		base = &kubeadmv1beta1.ClusterConfiguration{}
	}
	out := base.DeepCopy()

	s := machine.GetScope()

	// Only set the control plane endpoint if the user hasn't specified one.
	if out.ControlPlaneEndpoint == "" {
		out.ControlPlaneEndpoint = fmt.Sprintf("%s:%d", s.Network().APIServerELB.DNSName, APIServerBindPort)
	}
	// Add the control plane endpoint to the list of cert SAN
	out.APIServer.CertSANs = append(out.APIServer.CertSANs, localIPV4Lookup, s.Network().APIServerELB.DNSName)
	return out
}

// SetClusterConfigurationOverrides will modify the supplied configuration with certain values
// that cluster-api-provider-aws requires overriding user specified input.
func SetClusterConfigurationOverrides(machine joinMachine, base *kubeadmv1beta1.ClusterConfiguration) *kubeadmv1beta1.ClusterConfiguration {
	if base == nil {
		base = &kubeadmv1beta1.ClusterConfiguration{}
	}
	s := machine.GetScope()

	out := SetDefaultClusterConfiguration(machine, base.DeepCopy())

	// cloud-provider for the APIServer must be set to 'aws'.
	if out.APIServer.ControlPlaneComponent.ExtraArgs == nil {
		out.APIServer.ControlPlaneComponent.ExtraArgs = map[string]string{}
	}
	if cp, ok := out.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		machine.GetScope().Logger.Info("Overriding cloud provider with required value", "provided-cloud-provider", cp, "required-cloud-provider", CloudProvider)
	}
	out.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"] = CloudProvider

	if out.ControllerManager.ExtraArgs == nil {
		out.ControllerManager.ExtraArgs = map[string]string{}
	}
	if cp, ok := out.ControllerManager.ExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		machine.GetScope().Logger.Info("Overriding cloud provider with required value", "provided-cloud-provider", cp, "required-cloud-provider", CloudProvider)
	}
	out.ControllerManager.ExtraArgs["cloud-provider"] = CloudProvider

	// The kubeadm config clustername must match the provided name of the cluster.
	if out.ClusterName != "" && out.ClusterName != s.Name() {
		machine.GetScope().Logger.Info("Overriding provided cluster name. The kubeadm cluster name and cluster-api name must match.",
			"provided-cluster-name", out.ClusterName,
			"required-cluster-name", s.Name())
	}
	out.ClusterName = s.Name()

	// The networking values provided by the Cluster object must equal the
	// kubeadm networking configuration.
	out.Networking.DNSDomain = s.Cluster.Spec.ClusterNetwork.ServiceDomain
	out.Networking.PodSubnet = s.Cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]
	out.Networking.ServiceSubnet = s.Cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]

	// The kubernetes version that kubeadm is using must be the same as the
	// requested version in the config
	out.KubernetesVersion = machine.GetMachine().Spec.Versions.ControlPlane
	return out
}

// SetInitConfigurationOverrides overrides user input on particular fields for
// the kubeadm InitConfiguration.
func SetInitConfigurationOverrides(machine joinMachine, base *kubeadmv1beta1.InitConfiguration) *kubeadmv1beta1.InitConfiguration {
	if base == nil {
		base = &kubeadmv1beta1.InitConfiguration{}
	}
	out := base.DeepCopy()

	if out.NodeRegistration.Name != "" && out.NodeRegistration.Name != HostnameLookup {
		machine.GetScope().Info("Overriding NodeRegistration name. The node registration needs to be dynamically generated in aws.",
			"provided-node-registration-name", out.NodeRegistration.Name,
			"required-node-registration-name", HostnameLookup)
	}
	out.NodeRegistration.Name = HostnameLookup

	// TODO(chuckha): This may become a default instead of an override.
	if out.NodeRegistration.CRISocket != "" && out.NodeRegistration.CRISocket != ContainerdSocket {
		machine.GetScope().Info("Overriding CRISocket. Containerd is only supported container runtime.",
			"provided-container-runtime-socket", out.NodeRegistration.CRISocket,
			"required-container-runtime-socket", ContainerdSocket)
	}
	out.NodeRegistration.CRISocket = ContainerdSocket

	if out.NodeRegistration.KubeletExtraArgs == nil {
		out.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if cp, ok := out.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		machine.GetScope().Info("Overriding node's cloud-provider", "provided-cloud-provider", cp, "required-cloud-provider", CloudProvider)
	}
	out.NodeRegistration.KubeletExtraArgs["cloud-provider"] = CloudProvider
	if machine != nil && machine.GetMachine() != nil {
		out.NodeRegistration.Taints = machine.GetMachine().Spec.Taints
	}
	return out
}

// SetJoinNodeConfigurationOverrides overrides user input for certain fields of
// the kubeadm JoinConfiguration during a worker node join.
func SetJoinNodeConfigurationOverrides(caCertHash, bootstrapToken string, machine joinMachine, base *kubeadmv1beta1.JoinConfiguration) *kubeadmv1beta1.JoinConfiguration {
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
		machine.GetScope().Logger.Info("Overriding NodeRegistration name . The node registration needs to be dynamically generated in aws.",
			"provided-node-registration-name", out.NodeRegistration.Name,
			"required-node-registration-name", HostnameLookup)
	}
	out.NodeRegistration.Name = HostnameLookup

	// TODO(chuckha): This may become a default instead of an override.
	if out.NodeRegistration.CRISocket != "" && out.NodeRegistration.CRISocket != ContainerdSocket {
		machine.GetScope().Logger.Info("Overriding CRISocket. Containerd is only supported container runtime.",
			"provided-container-runtime-socket", out.NodeRegistration.CRISocket,
			"required-container-runtime-socket", ContainerdSocket)
	}
	out.NodeRegistration.CRISocket = ContainerdSocket

	if out.NodeRegistration.KubeletExtraArgs == nil {
		out.NodeRegistration.KubeletExtraArgs = map[string]string{}
	}
	if cp, ok := out.NodeRegistration.KubeletExtraArgs["cloud-provider"]; ok && cp != CloudProvider {
		machine.GetScope().Logger.Info("Overriding node's cloud-provider to the required value",
			"provided-cloud-provider", cp,
			"required-cloud-provider", CloudProvider)
	}
	out.NodeRegistration.KubeletExtraArgs["cloud-provider"] = CloudProvider
	if !util.IsControlPlaneMachine(machine.GetMachine()) {
		if val, ok := out.NodeRegistration.KubeletExtraArgs["node-labels"]; ok {
			labels := append(strings.Split(val, ","), nodeRole)
			out.NodeRegistration.KubeletExtraArgs["node-labels"] = strings.Join(labels, ",")
		} else {
			out.NodeRegistration.KubeletExtraArgs["node-labels"] = nodeRole
		}
	}
	if machine != nil && machine.GetMachine() != nil {
		out.NodeRegistration.Taints = machine.GetMachine().Spec.Taints
	}
	return out
}

// SetControlPlaneJoinConfigurationOverrides user input for kubeadm join
// configuration during a control plane join action.
func SetControlPlaneJoinConfigurationOverrides(base *kubeadmv1beta1.JoinConfiguration) *kubeadmv1beta1.JoinConfiguration {
	if base == nil {
		base = &kubeadmv1beta1.JoinConfiguration{}
	}
	out := base.DeepCopy()

	if out.ControlPlane == nil {
		out.ControlPlane = &kubeadmv1beta1.JoinControlPlane{}
	}
	out.ControlPlane.LocalAPIEndpoint.AdvertiseAddress = localIPV4Lookup
	out.ControlPlane.LocalAPIEndpoint.BindPort = APIServerBindPort
	return out
}
