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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	cabpkv1a2 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
)

const exampleClusterYAML = `
apiVersion: cluster.k8s.io/v1alpha1
kind: Cluster
metadata:
  name: "ponyville"
  namespace: "equestria"
spec:
  clusterNetwork:
    services:
      cidrBlocks:
      - "10.100.0.0/24"
    pods:
      cidrBlocks:
      - "10.100.10.0/24"
    serviceDomain: "gov.ponyville.eq"
  providerSpec:
    value:
      networkSpec:
        vpc:
          id: "vpc-m4g1c"
          cidrBlock: "192.168.0.0/24"
          internetGatewayId: "i-shy"
          tags:
            vpc: "ismagic"
      subnets:
      - id: "twilight"
        cidrBlock: "192.168.10.0/24"
        availabilityZone: "equestria-west2a"
        isPublic: true,
        routeTableId: "friendshipMap"
        natGatewayId: "mirror"
        tags:
          twilight: alicorn
      region: "equestria-west2"
      sshKeyName: "harmony"
      caKeyPair:
         cert: "dHdpbGlnaHQK"
         key: "c3Rhcgo="
      etcdCAKeyPair:
         cert: "cGlua2llCg=="
         key: "YmFsbG9vbgo="
      frontProxyCAKeyPair:
         cert: "cmFpbmJvdwo="
         key: "Y2xvdWQK"
      saKeyPair:
         cert: "YXBwbGVqYWNrCg=="
         key: "YXBwbGUK"
      clusterConfiguration:
        etcd:
          local:
            imageRepository: "ponyville/library"
            imageTag: "goldenoaks"
            dataDir: "/var/lib/goldenoaks"
            extraArgs:
              "--window-open": "true"
            serverCertSANs:
            - "goldenoaks.ponyville.eq"
            peerCertSANs:
            - "your.local.library"
          external:
            endpoints:
            - "goldenoaks.ponyville.eq"
            caFile: "/etc/canterlot.cert"
            certFile: "/etc/oaks.cert"
            keyFile: "/etc/oaks.key"
        networking:
          serviceSubnet: "172.16.0.0/24"
          podSubnet: "172.16.80.0/24"
          dnsDomain: "ponyville.eq"
        kubernetesVersion: "v1.15.2"
        controlPlaneEndpoint: "castle.ponyville.eq"
        apiServer:
          extraArgs:
            "--run": "leaves"
          extraVolumes:
          - name: "magic"
            hostPath: "/etc/magic"
            mountPath: "/var/lib/magic"
            readOnly: true
            pathType: "BlockDevice"
          certSANs:
          - "door.castle.ponyville.eq"
          timeoutForControlPlane: "30s"
        controllerManager:
          extraArgs:
            "--winter-wrapup": "yes"
        scheduler:
          extraArgs:
            "--tardy": "false"
        dns:
          type: "kube-dns"
          imageRepository: "ponyville/phonebook"
          imageTag: "s4"
        certificatesDir: "/etc/scrolls"
        imageRepository: "canter.eq"
        useHyperKubeImage: true
        featureGates:
          "alicorn": true
        clusterName: "ponyville"
      additionalUserDataFiles:
      - path: "/lib/journal"
        owner: "twilight:harmony"
        permissions: "0777"
        content: "today I learned..."
      - path: "/var/spool/mail/celestia"
        owner: "twilight:twilight"
        permissions: "0660"
        content: "dear princess celestia..."
`

func getCluster(t *testing.T) (*capiv1a1.Cluster, *capav1a1.AWSClusterProviderSpec) {
	scheme := runtime.NewScheme()
	capiv1a1.SchemeBuilder.AddToScheme(scheme)
	capav1a1.SchemeBuilder.AddToScheme(scheme)

	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder()

	var (
		cluster    capiv1a1.Cluster
		awsCluster capav1a1.AWSClusterProviderSpec
	)

	if _, _, err := decoder.Decode([]byte(exampleClusterYAML), nil, &cluster); err != nil {
		t.Fatalf("failed to decode example: %v", err)
	}

	if _, _, err := decoder.Decode(cluster.Spec.ProviderSpec.Value.Raw, nil, &awsCluster); err != nil {
		t.Fatalf("failed to decode example providerSpce: %v", err)

	}

	return &cluster, &awsCluster
}

func TestConvertCluster(t *testing.T) {
	var (
		newCluster    capiv1a2.Cluster
		newAWSCluster capav1a2.AWSCluster
	)

	oldCluster, oldAWSCluster := getCluster(t)

	converter := NewClusterConverter(oldCluster)

	if err := converter.GetCluster(&newCluster); err != nil {
		t.Fatalf("Unexpected error converting cluster: %v", err)
	}

	assert := asserter{t}

	assert.stringEqual(oldCluster.Name, newCluster.Name, "name")
	assert.stringEqual(oldCluster.Namespace, newCluster.Namespace, "namespace")

	assert.stringArrayEqual(
		oldCluster.Spec.ClusterNetwork.Services.CIDRBlocks,
		newCluster.Spec.ClusterNetwork.Services.CIDRBlocks,
		"services CIDR blocks",
	)
	assert.stringArrayEqual(
		oldCluster.Spec.ClusterNetwork.Pods.CIDRBlocks,
		newCluster.Spec.ClusterNetwork.Pods.CIDRBlocks,
		"pods CIDR blocks",
	)
	assert.stringEqual(oldCluster.Spec.ClusterNetwork.ServiceDomain, oldCluster.Spec.ClusterNetwork.ServiceDomain, "service domain")

	if err := converter.GetAWSCluster(&newAWSCluster); err != nil {
		t.Fatalf("Unexpected error converting AWS Cluster: %v", err)
	}

	t.Logf("converted cluster: %+v", newAWSCluster)

	assert.stringEqual(oldAWSCluster.NetworkSpec.VPC.ID, newAWSCluster.Spec.NetworkSpec.VPC.ID, "vpc ID")
	assert.stringEqual(oldAWSCluster.NetworkSpec.VPC.CidrBlock, newAWSCluster.Spec.NetworkSpec.VPC.CidrBlock, "VPC cidr block")

	if newAWSCluster.Spec.NetworkSpec.VPC.InternetGatewayID == nil {
		t.Errorf("Expected InternetGatewayID %q, got nil", *oldAWSCluster.NetworkSpec.VPC.InternetGatewayID)
	} else {
		assert.stringEqual(*oldAWSCluster.NetworkSpec.VPC.InternetGatewayID, *newAWSCluster.Spec.NetworkSpec.VPC.InternetGatewayID, "VPC gateway ID")
	}

	oldTags := oldAWSCluster.NetworkSpec.VPC.Tags
	newTags := newAWSCluster.Spec.NetworkSpec.VPC.Tags

	if len(oldTags) == len(newTags) {
		for key := range oldAWSCluster.NetworkSpec.VPC.Tags {
			assert.stringEqual(oldTags[key], newTags[key], fmt.Sprintf("VPC tag %s", key))
		}
	} else {
		t.Errorf("VPC tags has length %d, expected %d", len(newTags), len(oldTags))
	}

	if len(oldAWSCluster.NetworkSpec.Subnets) == len(newAWSCluster.Spec.NetworkSpec.Subnets) {
		for i, subnet := range oldAWSCluster.NetworkSpec.Subnets {
			assert.stringEqual(subnet.String(), newAWSCluster.Spec.NetworkSpec.Subnets[i].String(), fmt.Sprintf("subnet[%d]", i))
		}
	} else {
		assert.Errorf(
			"Subnet has length %d, expected %d",
			len(newAWSCluster.Spec.NetworkSpec.Subnets),
			len(oldAWSCluster.NetworkSpec.Subnets),
		)
	}

	assert.stringEqual(oldAWSCluster.Region, newAWSCluster.Spec.Region, "region")
	assert.stringEqual(oldAWSCluster.SSHKeyName, newAWSCluster.Spec.SSHKeyName, "sshkey")

	if newCluster.Spec.InfrastructureRef == nil {
		t.Error("Unexpectedly nil infrastructure ref")
	} else {
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Name, newAWSCluster.Name, "aws cluster ref name")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Namespace, newAWSCluster.Namespace, "aws cluster ref namespace")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Kind, "AWSCluster", "aws cluster ref kind")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.APIVersion, "infrastructure.cluster.x-k8s.io/v1alpha2", "aws cluster ref apiversion")
	}

	uid := "bf022564-a1ff-4c72-b211-c9ae0082b46b"

	kubeadmCfg := &cabpkv1a2.KubeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ponyville-config",
			UID:  types.UID(uid),
		},
	}

	secrets, err := converter.GetSecrets(&newCluster, kubeadmCfg)
	if err != nil {
		t.Fatalf("Unexpected error getting secrets: %v", err)
	}

	if len(secrets) == 4 {
		expected := []struct {
			name string
			cert []byte
			key  []byte
		}{
			{
				name: "ponyville-ca",
				cert: oldAWSCluster.CAKeyPair.Cert,
				key:  oldAWSCluster.CAKeyPair.Key,
			},
			{
				name: "ponyville-etcd",
				cert: oldAWSCluster.EtcdCAKeyPair.Cert,
				key:  oldAWSCluster.EtcdCAKeyPair.Key,
			},
			{
				name: "ponyville-proxy",
				cert: oldAWSCluster.FrontProxyCAKeyPair.Cert,
				key:  oldAWSCluster.FrontProxyCAKeyPair.Key,
			},
			{
				name: "ponyville-sa",
				cert: oldAWSCluster.SAKeyPair.Cert,
				key:  oldAWSCluster.SAKeyPair.Key,
			},
		}

		for i, pair := range expected {
			actual := secrets[i]
			assert.stringEqual(pair.name, actual.Name, fmt.Sprintf("secret[%d] name", i))
			assert.stringEqual(oldCluster.Namespace, actual.Namespace, fmt.Sprintf("secret[%d] namespace", i))
			assert.stringEqual(string(pair.cert), string(actual.Data["tls.crt"]), fmt.Sprintf("secret[%d] cert", i))
			assert.stringEqual(string(pair.key), string(actual.Data["tls.key"]), fmt.Sprintf("secret[%d] key", i))
			assert.stringEqual(
				oldCluster.Name,
				actual.Labels[capiv1a2.MachineClusterLabelName],
				fmt.Sprintf("secret[%d] label name", i))

			if len(actual.OwnerReferences) == 1 {
				ref := actual.OwnerReferences[0]
				assert.stringEqual(kubeadmCfg.Name, ref.Name, fmt.Sprintf("secret[%d] ownerref name", i))
				assert.stringEqual(uid, string(ref.UID), fmt.Sprintf("secret[%d] ownerref uid", i))
				assert.stringEqual(cabpkv1a2.GroupVersion.String(), ref.APIVersion, fmt.Sprintf("secret[%d] ownerref uid", i))
				assert.stringEqual("KubeadmConfig", ref.Kind, fmt.Sprintf("secret[%d] ownerref uid", i))
			} else {
				t.Errorf("Expected 1 owner reference, got %d", len(actual.OwnerReferences))
			}
		}
	} else {
		t.Errorf("Expected 4 secrets, got %d", len(secrets))
	}
}
