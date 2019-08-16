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
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
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
      - "10.100.0.0/24"
    serviceDomain: "gov.ponyville.eq"
  providerSpec:
    value:
      networkSpec:
        vpc:
          id: "vpc-m4g1c"
          cidrBlock: "192.168.0.0./24"
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
	oldCluster, oldAWSCluster := getCluster(t)

	newCluster, newAWSCluster, err := ConvertCluster(oldCluster)
	if err != nil {
		t.Fatalf("Unexpected error running convertCluster: %v", err)
	}

	assert := asserter{t}

	if newCluster == nil {
		t.Fatal("Unexepctedly nil cluster")
	}

	assert.stringEqual(oldCluster.Name, newCluster.Name, "name")
	assert.stringEqual(oldCluster.Namespace, newCluster.Namespace, "namespace")
	assert.stringEqual(string(oldCluster.UID), string(newCluster.UID), "UID")

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

	if newAWSCluster == nil {
		t.Fatalf("unexpectedly nil provider spec")
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
		assert.notEmpty(newAWSCluster.Name, "awscluster name")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Name, newAWSCluster.Name, "aws cluster ref name")
		assert.notEmpty(newAWSCluster.Namespace, "awscluster namespace")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Namespace, newAWSCluster.Namespace, "aws cluster ref namespace")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.Kind, "AWSCluster", "aws cluster ref kind")
		assert.stringEqual(newCluster.Spec.InfrastructureRef.APIVersion, "infrastructure.cluster.x-k8s.io/v1alpha2", "aws cluster ref apiversion")
	}

}
