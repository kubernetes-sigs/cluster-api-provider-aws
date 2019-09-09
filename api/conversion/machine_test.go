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

package conversion

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	cabpkv1a2 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
)

const exampleMachineYAML = `
apiVersion: "cluster.k8s.io/v1alpha1"
kind: "Machine"
metadata:
  name: "rainbow"
  namespace: "equestria"
spec:
  metadata:
    name: "rainbow-"
    namespace: "equestria"
  versions:
    kubelet: "1.10.1"
    controlPlane: "1.11.2"
  providerID: "element://loyalty"
  providerSpec:
    value:
      ami:
        arn: "equestria://rainbow"
      imageLookupOrg: "w0nd3rb0lt5"
      instanceType: "pegasus.1xsmall"
      additionalTags:
        profession: "weather"
      iamInstanceProfile: "element-of-loyalty"
      publicIP: true
      additionalSecurityGroups:
      - id: "branch"
        arn: "eq://wb"
        filter:
        - name: "species"
          values:
          - "pegasus"
          - "alicorn"
      availabilityZone: "equestria-west2a"
      subnet:
        id: "weather"
        arn: "eq://ponyville"
      keyName: "loyalty"
      rootSize: 1234567
      kubeadmConfiguration:
        init:
          bootstrapTokens:
            - token: "horses.mylittlemagical0"
              description: "bootstrap element"
              ttl: "1234ms"
              expires:  "1996-12-19T16:39:57-08:00"
              usages:
              - friendship
              groups:
              - elements
          nodeRegistration:
            name: "rainbow"
            criSocket: "magic://cloud"
            taints:
            - key: "disposition"
              value: "lazy"
              effect: "PreferNoSchedule"
            kubeletExtraArgs:
              "--go-fast": "true"
          localAPIEndpoint:
            advertiseAddress: "10.100.0.1"
            bindPort: 65000
        join:
          nodeRegistration:
            name: "rainbow-dash"
            criSocket: "magic://strato/cloud"
            taints:
            - key: "job"
              value: "sleeping"
              effect: "NoSchedule"
            kubeletExtraArgs:
              "--wonderbolt": "true"
          caCertPath: "/tmp/wonderbolt.sock"
          discovery:
            bootstrapToken:
              token: "r41nb0w"
              apiServerEndpoint: "http://friendship.castle"
              caCertHashes:
              - "c313571a"
              unsafeSkipCAVerification: false
            fileDiscovery:
              kubeConfig: "/tmp/rainbow.conf"
            tlsBootstrapToken: "7w1l1ght"
            timeout: "8791ms"
          controlPlane:
            localAPIEndpoint:
              advertiseAddress: "friendship.castle"
              bindPort: 1235
      additionalUserDataFiles:
        - path: "/tmp/tank"
          owner: "rd:rd"
          permissions: "0640"
          content: "the best pet for me"
`

func getMachine(t *testing.T) (*capiv1a1.Machine, *capav1a1.AWSMachineProviderSpec) {
	scheme := runtime.NewScheme()
	capiv1a1.SchemeBuilder.AddToScheme(scheme)
	capav1a1.SchemeBuilder.AddToScheme(scheme)

	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder()

	var (
		machine    capiv1a1.Machine
		awsMachine capav1a1.AWSMachineProviderSpec
	)

	if _, _, err := decoder.Decode([]byte(exampleMachineYAML), nil, &machine); err != nil {
		t.Fatalf("failed to decode example: %v", err)
	}

	if machine.Spec.ProviderSpec.Value == nil {
		t.Fatalf("No providerspec found")
	}

	if _, _, err := decoder.Decode(machine.Spec.ProviderSpec.Value.Raw, nil, &awsMachine); err != nil {
		t.Fatalf("failed to decode example providerSpec: %v", err)
	}

	return &machine, &awsMachine
}

func TestConvertMachine(t *testing.T) {
	var (
		newMachine       capiv1a2.Machine
		newAWSMachine    capav1a2.AWSMachine
		newKubeadmConfig cabpkv1a2.KubeadmConfig
	)

	oldMachine, oldAWSMachine := getMachine(t)
	oldCluster, oldAWSCluster := getCluster(t)

	converter := NewMachineConverter(oldCluster, oldMachine)

	if err := converter.GetMachine(&newMachine); err != nil {
		t.Fatalf("Unexpected error converting machine: %v", err)
	}

	assert := asserter{t}

	if oldMachine == nil {
		t.Fatalf("Unexpectedly nil machine")
	}

	assert.stringEqual(oldMachine.ObjectMeta.Name, newMachine.ObjectMeta.Name, "machine name")
	assert.stringEqual(oldMachine.ObjectMeta.Namespace, newMachine.ObjectMeta.Namespace, "machine namespace")

	assert.stringEqual(oldMachine.Spec.ObjectMeta.Name, newMachine.Spec.ObjectMeta.Name, "node name")
	assert.stringEqual(oldMachine.Spec.ObjectMeta.Namespace, newMachine.Spec.ObjectMeta.Namespace, "node namespace")

	assert.stringPtrEqual(&oldMachine.Spec.Versions.ControlPlane, newMachine.Spec.Version, "version")

	assert.stringPtrEqual(oldMachine.Spec.ProviderID, newMachine.Spec.ProviderID, "provider ID")

	if err := converter.GetAWSMachine(&newAWSMachine); err != nil {
		t.Fatalf("Unexpected error converting AWSMachine: %v", err)
	}

	t.Logf("converted machine: %+v", newAWSMachine)

	// Pull the provider ID from the machine
	// It's in both places because the infra provider is the primary component responsible for the value, but CAPI needs it
	// for setting node refs, so CAPI copies it from the infra resource to the machine.
	assert.stringPtrEqual(newAWSMachine.Spec.ProviderID, newMachine.Spec.ProviderID, "aws machine provider ID")

	assert.stringEqual(newAWSMachine.Name, oldMachine.Name, "aws machine name")
	assert.stringEqual(newAWSMachine.Namespace, oldMachine.Namespace, "aws machine namespace")

	assert.stringEqual(newAWSMachine.Name, newMachine.Spec.InfrastructureRef.Name, "infra ref name")
	assert.stringEqual(newAWSMachine.Namespace, newMachine.Spec.InfrastructureRef.Namespace, "infra ref namespace")
	assert.stringEqual("AWSMachine", newMachine.Spec.InfrastructureRef.Kind, "infra ref kind")
	assert.stringEqual("infrastructure.cluster.x-k8s.io/v1alpha2", newMachine.Spec.InfrastructureRef.APIVersion, "infra ref APIVersion")

	assert.awsRefEqual(&oldAWSMachine.AMI, &newAWSMachine.Spec.AMI, "AMI")
	assert.stringEqual(oldAWSMachine.ImageLookupOrg, newAWSMachine.Spec.ImageLookupOrg, "image lookup org")
	assert.stringEqual(oldAWSMachine.InstanceType, newAWSMachine.Spec.InstanceType, "instance type")

	assert.stringMapEqual(oldAWSMachine.AdditionalTags, newAWSMachine.Spec.AdditionalTags, "additional tags")

	if newAWSMachine.Spec.PublicIP == nil {
		t.Errorf("public ip should be %v, was nil", *oldAWSMachine.PublicIP)
	}

	if len(oldAWSMachine.AdditionalSecurityGroups) == len(newAWSMachine.Spec.AdditionalSecurityGroups) {
		for i := range oldAWSMachine.AdditionalSecurityGroups {
			assert.awsRefEqual(&oldAWSMachine.AdditionalSecurityGroups[i], &newAWSMachine.Spec.AdditionalSecurityGroups[i], fmt.Sprintf("AdditionalSecurityGroups[%d]", i))
		}

	} else {
		t.Errorf(
			"AdditionalSecurityGroups has length %d, expected %d",
			len(newAWSMachine.Spec.AdditionalSecurityGroups),
			len(oldAWSMachine.AdditionalSecurityGroups),
		)
	}

	assert.stringPtrEqual(oldAWSMachine.AvailabilityZone, newAWSMachine.Spec.AvailabilityZone, "availability zone")
	assert.awsRefEqual(oldAWSMachine.Subnet, newAWSMachine.Spec.Subnet, "subnet")

	assert.stringEqual(oldAWSMachine.KeyName, newAWSMachine.Spec.SSHKeyName, "KeyName")

	if oldAWSMachine.RootDeviceSize != newAWSMachine.Spec.RootDeviceSize {
		t.Errorf("Expected RoodDeviceSize %d, got %d", oldAWSMachine.RootDeviceSize, newAWSMachine.Spec.RootDeviceSize)
	}

	if err := converter.GetKubeadmConfig(&newKubeadmConfig); err != nil {
		t.Fatalf("Unexpected error getting kubeadm config: %v", err)
	}

	t.Logf("KubeadmConfig: %+v", newKubeadmConfig)

	oldInit := oldAWSMachine.KubeadmConfiguration.Init
	newInit := newKubeadmConfig.Spec.InitConfiguration

	if len(oldInit.BootstrapTokens) == len(newInit.BootstrapTokens) {
		for i, oldToken := range oldInit.BootstrapTokens {
			newToken := newInit.BootstrapTokens[i]

			assert.stringEqual(oldToken.Token.String(), newToken.Token.String(), fmt.Sprintf("token[%d] token", i))
			assert.stringEqual(oldToken.Description, newToken.Description, fmt.Sprintf("token[%d] description", i))
			assert.stringEqual(oldToken.TTL.String(), newToken.TTL.String(), fmt.Sprintf("token[%d] TTL", i))
			assert.stringEqual(oldToken.Expires.String(), newToken.Expires.String(), fmt.Sprintf("token[%d] TTL", i))
			assert.stringArrayEqual(oldToken.Groups, newToken.Groups, fmt.Sprintf("token[%d] Groups", i))
			assert.stringArrayEqual(oldToken.Groups, newToken.Groups, fmt.Sprintf("token[%d] Groups", i))
		}
	} else {
		t.Errorf("BootstrapTokens has length %d, expected %d", len(newInit.BootstrapTokens), len(oldInit.BootstrapTokens))
	}

	assert.nodeRegistrationEqual(&oldInit.NodeRegistration, &newInit.NodeRegistration, "init node registration")

	assert.stringEqual(oldInit.LocalAPIEndpoint.AdvertiseAddress, newInit.LocalAPIEndpoint.AdvertiseAddress, "init API address")
	// :shrug:
	assert.stringEqual(string(oldInit.LocalAPIEndpoint.BindPort), string(newInit.LocalAPIEndpoint.BindPort), "init API port")

	oldJoin := oldAWSMachine.KubeadmConfiguration.Join
	newJoin := newKubeadmConfig.Spec.JoinConfiguration

	if newJoin == nil {
		t.Fatal("Join information unexpectedly nil")
	}

	assert.nodeRegistrationEqual(&oldJoin.NodeRegistration, &newJoin.NodeRegistration, "join node registration")
	assert.stringEqual(oldJoin.CACertPath, newJoin.CACertPath, "join CA cert path")

	assert.stringEqual(oldJoin.Discovery.BootstrapToken.Token, newJoin.Discovery.BootstrapToken.Token, "join discovery bootstrap token")
	assert.stringEqual(oldJoin.Discovery.BootstrapToken.APIServerEndpoint, newJoin.Discovery.BootstrapToken.APIServerEndpoint, "join discovery bootstrap APIServerEndpoint")
	assert.stringArrayEqual(oldJoin.Discovery.BootstrapToken.CACertHashes, newJoin.Discovery.BootstrapToken.CACertHashes, "join discovery bootstrap token ca cert hashes")

	if oldJoin.Discovery.BootstrapToken.UnsafeSkipCAVerification != newJoin.Discovery.BootstrapToken.UnsafeSkipCAVerification {
		t.Errorf("Expected join discovery bootstrap token validation to be %v, got %v",
			oldJoin.Discovery.BootstrapToken.UnsafeSkipCAVerification,
			newJoin.Discovery.BootstrapToken.UnsafeSkipCAVerification,
		)
	}

	if oldJoin.Discovery.File == nil {
		if newJoin.Discovery.File != nil {
			t.Errorf("Expected discovery file to be nil, but was %s", newJoin.Discovery.File.KubeConfigPath)
		}
	} else {
		if newJoin.Discovery.File == nil {
			t.Errorf("Expected discovery file to be %s, but was nil", oldJoin.Discovery.File.KubeConfigPath)
		} else {
			assert.stringEqual(oldJoin.Discovery.File.KubeConfigPath, newJoin.Discovery.File.KubeConfigPath, "join discovery file kubeconfig path")
		}
	}

	assert.stringEqual(oldJoin.Discovery.TLSBootstrapToken, newJoin.Discovery.TLSBootstrapToken, "join discovery TLS bootsrap token")
	assert.stringEqual(oldJoin.Discovery.Timeout.String(), newJoin.Discovery.Timeout.String(), "join discovery timeout")

	assert.stringEqual(oldJoin.ControlPlane.LocalAPIEndpoint.AdvertiseAddress, newJoin.ControlPlane.LocalAPIEndpoint.AdvertiseAddress, "local API endpoint advertise address")
	assert.stringEqual(oldJoin.ControlPlane.LocalAPIEndpoint.AdvertiseAddress, newJoin.ControlPlane.LocalAPIEndpoint.AdvertiseAddress, "local API endpoint advertise address")

	oldClusterCfg := oldAWSCluster.ClusterConfiguration
	newClusterCfg := newKubeadmConfig.Spec.ClusterConfiguration

	if newClusterCfg == nil {
		t.Fatalf("ClusterConfiguration is unexpectedly nil")
	}

	if newClusterCfg.Etcd.Local == nil {
		t.Errorf("etcd local unexpectedly nil")
	} else {

		assert.stringEqual(
			oldClusterCfg.Etcd.Local.ImageRepository,
			newClusterCfg.Etcd.Local.ImageRepository,
			"local etcd image repo",
		)
		assert.stringEqual(
			oldClusterCfg.Etcd.Local.ImageTag,
			newClusterCfg.Etcd.Local.ImageTag,
			"local etcd image tag",
		)
		assert.stringEqual(
			oldClusterCfg.Etcd.Local.DataDir,
			newClusterCfg.Etcd.Local.DataDir,
			"local etcd data dir",
		)

		assert.stringMapEqual(
			oldClusterCfg.Etcd.Local.ExtraArgs,
			newClusterCfg.Etcd.Local.ExtraArgs,
			"local etcd extra args",
		)
		assert.stringArrayEqual(
			oldClusterCfg.Etcd.Local.ServerCertSANs,
			newClusterCfg.Etcd.Local.ServerCertSANs,
			"local etcd server cert sans",
		)
		assert.stringArrayEqual(
			oldClusterCfg.Etcd.Local.PeerCertSANs,
			newClusterCfg.Etcd.Local.PeerCertSANs,
			"local etcd peer cert sans",
		)
	}

	if newClusterCfg.Etcd.External == nil {
		t.Errorf("etcd external unexpectedly nil")
	} else {

		assert.stringArrayEqual(
			oldClusterCfg.Etcd.External.Endpoints,
			newClusterCfg.Etcd.External.Endpoints,
			"external etcd endpoints",
		)
		assert.stringEqual(
			oldClusterCfg.Etcd.External.CAFile,
			newClusterCfg.Etcd.External.CAFile,
			"external etcd CAFile",
		)
		assert.stringEqual(
			oldClusterCfg.Etcd.External.CertFile,
			newClusterCfg.Etcd.External.CertFile,
			"external etcd CertFile",
		)
		assert.stringEqual(
			oldClusterCfg.Etcd.External.KeyFile,
			newClusterCfg.Etcd.External.KeyFile,
			"external etcd KeyFile",
		)
	}

	assert.stringEqual(
		oldClusterCfg.Networking.ServiceSubnet,
		newClusterCfg.Networking.ServiceSubnet,
		"Networking ServiceSubnet",
	)
	assert.stringEqual(
		oldClusterCfg.Networking.PodSubnet,
		newClusterCfg.Networking.PodSubnet,
		"Networking PodSubnet",
	)
	assert.stringEqual(
		oldClusterCfg.Networking.DNSDomain,
		newClusterCfg.Networking.DNSDomain,
		"Networking DNSDomain",
	)

	assert.stringEqual(oldClusterCfg.KubernetesVersion, newClusterCfg.KubernetesVersion, "kubernetes version")

	assert.stringEqual(oldClusterCfg.ControlPlaneEndpoint, newClusterCfg.ControlPlaneEndpoint, "control plane endpoint")

	assert.stringMapEqual(
		oldClusterCfg.APIServer.ExtraArgs,
		newClusterCfg.APIServer.ExtraArgs,
		"apiserver extra args",
	)
	if len(newClusterCfg.APIServer.ExtraVolumes) == 1 {
		oldV := newClusterCfg.APIServer.ExtraVolumes[0]
		newV := newClusterCfg.APIServer.ExtraVolumes[0]

		assert.stringEqual(oldV.Name, newV.Name, "APIServer extra volume name")
		assert.stringEqual(oldV.HostPath, newV.HostPath, "APIServer extra volume host path")
		assert.stringEqual(oldV.MountPath, newV.MountPath, "APIServer extra volume mount path")
		if oldV.ReadOnly != newV.ReadOnly {
			t.Errorf("Expected APIServer extra volume readOnly to be %v, was %v", oldV.ReadOnly, newV.ReadOnly)
		}
		assert.stringEqual(string(oldV.PathType), string(newV.PathType), "APIServer extra volume mount path")
	} else {
		t.Errorf(
			"Expcted APIServer ExtraVolumes to have length 1, got %d",
			len(newClusterCfg.APIServer.ExtraVolumes),
		)
	}

	assert.stringArrayEqual(oldClusterCfg.APIServer.CertSANs, newClusterCfg.APIServer.CertSANs, "APIServer cert SANs")
	assert.stringMapEqual(oldClusterCfg.ControllerManager.ExtraArgs, newClusterCfg.ControllerManager.ExtraArgs, "ControllerManager extra args")
	assert.stringMapEqual(oldClusterCfg.Scheduler.ExtraArgs, newClusterCfg.Scheduler.ExtraArgs, "scheduler extra args")

	assert.stringEqual(string(oldClusterCfg.DNS.Type), string(newClusterCfg.DNS.Type), "DNS Type")
	assert.stringEqual(oldClusterCfg.DNS.ImageRepository, newClusterCfg.DNS.ImageRepository, "DNS image repo")
	assert.stringEqual(oldClusterCfg.DNS.ImageTag, newClusterCfg.DNS.ImageTag, "DNS image tag")

	assert.stringEqual(oldClusterCfg.CertificatesDir, newClusterCfg.CertificatesDir, "Certificates dir")
	assert.stringEqual(oldClusterCfg.ImageRepository, newClusterCfg.ImageRepository, "ImageRepository")

	if oldClusterCfg.UseHyperKubeImage != newClusterCfg.UseHyperKubeImage {
		t.Errorf("Expected ClusterCfg hyperkube image to be %v, got %v",
			oldClusterCfg.UseHyperKubeImage,
			newClusterCfg.UseHyperKubeImage,
		)
	}

	if len(oldClusterCfg.FeatureGates) == len(newClusterCfg.FeatureGates) {
		for key := range oldClusterCfg.FeatureGates {
			if oldClusterCfg.FeatureGates[key] != newClusterCfg.FeatureGates[key] {
				t.Errorf("Expected ClusterCfg feature gate %s image to be %v, got %v",
					key,
					oldClusterCfg.FeatureGates,
					newClusterCfg.FeatureGates,
				)
			}
		}
	} else {
		t.Errorf("Expected Feature gates ho have length %d, has length %d", len(oldClusterCfg.FeatureGates), len(newClusterCfg.FeatureGates))
	}
	assert.stringEqual(oldClusterCfg.ClusterName, newClusterCfg.ClusterName, "Cluster name")

	if len(oldAWSCluster.AdditionalUserDataFiles) == len(newKubeadmConfig.Spec.Files) {
		for i, oldFile := range oldAWSCluster.AdditionalUserDataFiles {
			newFile := newKubeadmConfig.Spec.Files[i]

			assert.stringEqual(oldFile.Path, newFile.Path, fmt.Sprintf("additional file[%d] path", i))
			assert.stringEqual(oldFile.Owner, newFile.Owner, fmt.Sprintf("additional file[%d] owner", i))
			assert.stringEqual(oldFile.Permissions, newFile.Permissions, fmt.Sprintf("additional file[%d] permissions", i))
			assert.stringEqual(oldFile.Content, newFile.Content, fmt.Sprintf("additional file[%d] content", i))
		}
	} else {
		assert.Errorf("expected AdditionalUserDataFiles to have length %d, has %d", len(oldAWSCluster.AdditionalUserDataFiles), len(newKubeadmConfig.Spec.Files))
	}
}
