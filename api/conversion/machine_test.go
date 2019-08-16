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
	oldMachine, oldAWSMachine := getMachine(t)

	newMachine, newAWSMachine, err := ConvertMachine(oldMachine)

	if err != nil {
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

	if newAWSMachine == nil {
		t.Fatalf("AWSMachine unexpectedly nil")
	}

	t.Logf("converted machine: %+v", newAWSMachine)

	// Pull the provider ID from the machine (why is it in both places?)
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

	oldTags := oldAWSMachine.AdditionalTags
	newTags := newAWSMachine.Spec.AdditionalTags

	if len(oldTags) == len(newTags) {
		for key := range oldTags {
			assert.stringEqual(oldTags[key], newTags[key], fmt.Sprintf("machine tag %s", key))
		}
	} else {
		t.Errorf("Machine tags has length %d, expected %d", len(newTags), len(oldTags))
	}

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

	assert.stringEqual(oldAWSMachine.KeyName, newAWSMachine.Spec.KeyName, "KeyName")

	if oldAWSMachine.RootDeviceSize != newAWSMachine.Spec.RootDeviceSize {
		t.Errorf("Expected RoodDeviceSize %d, got %d", oldAWSMachine.RootDeviceSize, newAWSMachine.Spec.RootDeviceSize)
	}
}
