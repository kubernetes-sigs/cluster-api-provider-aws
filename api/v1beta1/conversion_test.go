/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"

	fuzz "github.com/google/gofuzz"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func fuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		AWSMachineFuzzer,
		AWSMachineTemplateFuzzer,
	}
}

func AWSMachineFuzzer(obj *AWSMachine, c fuzz.Continue) {
	c.FuzzNoCustom(obj)
	
	// AWSMachine.Spec.FailureDomain, AWSMachine.Spec.Subnet.ARN and AWSMachine.Spec.AdditionalSecurityGroups.ARN has been removed in v1beta2, so setting it to nil in order to avoid v1beta1 --> v1beta2 --> v1beta1 round trip errors.
	if obj.Spec.Subnet != nil {
		obj.Spec.Subnet.ARN = nil
	}
	restored := make([]AWSResourceReference, len(obj.Spec.AdditionalSecurityGroups))
	for _, sg := range obj.Spec.AdditionalSecurityGroups {
		sg.ARN = nil
		restored = append(restored, sg)
	}
	obj.Spec.AdditionalSecurityGroups = restored
	obj.Spec.FailureDomain = nil
}

func AWSMachineTemplateFuzzer(obj *AWSMachineTemplate, c fuzz.Continue) {
	c.FuzzNoCustom(obj)
	
	// AWSMachineTemplate.Spec.Template.Spec.FailureDomain, AWSMachineTemplate.Spec.Template.Spec.Subnet.ARN and AWSMachineTemplate.Spec.Template.Spec.AdditionalSecurityGroups.ARN has been removed in v1beta2, so setting it to nil in order to avoid  v1beta1 --> v1beta2 --> v1beta round trip errors.
	if obj.Spec.Template.Spec.Subnet != nil {
		obj.Spec.Template.Spec.Subnet.ARN = nil
	}
	restored := make([]AWSResourceReference, len(obj.Spec.Template.Spec.AdditionalSecurityGroups))
	for _, sg := range obj.Spec.Template.Spec.AdditionalSecurityGroups {
		sg.ARN = nil
		restored = append(restored, sg)
	}
	obj.Spec.Template.Spec.AdditionalSecurityGroups = restored
	obj.Spec.Template.Spec.FailureDomain = nil
}

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(v1beta2.AddToScheme(scheme)).To(Succeed())

	t.Run("for AWSCluster", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSCluster{},
		Spoke:  &AWSCluster{},
	}))

	t.Run("for AWSMachine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSMachine{},
		Spoke:  &AWSMachine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSMachineTemplate", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSMachineTemplate{},
		Spoke:  &AWSMachineTemplate{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSClusterStaticIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSClusterStaticIdentity{},
		Spoke:  &AWSClusterStaticIdentity{},
	}))

	t.Run("for AWSClusterControllerIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSClusterControllerIdentity{},
		Spoke:  &AWSClusterControllerIdentity{},
	}))

	t.Run("for AWSClusterRoleIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta2.AWSClusterRoleIdentity{},
		Spoke:  &AWSClusterRoleIdentity{},
	}))
}
