/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	"testing"

	. "github.com/onsi/gomega"

	fuzz "github.com/google/gofuzz"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtime "k8s.io/apimachinery/pkg/runtime"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	v1alpha4 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func fuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		AWSClusterStaticIdentityFuzzer,
		AWSMachineFuzzer,
		AWSMachineTemplateFuzzer,
	}
}

func AWSClusterStaticIdentityFuzzer(obj *AWSClusterStaticIdentity, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSClusterStaticIdentity.Spec.SecretRef.Namespace has been removed in v1alpha4, so setting it to nil in order to avoid v1alpha3 --> v1alpha4 --> v1alpha3 round trip errors.
	obj.Spec.SecretRef.Namespace = ""
}

func AWSMachineFuzzer(obj *AWSMachine, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSMachine.Spec.AMI.ARN and AWSMachine.Spec.AMI.Filters has been removed in v1alpha4, so setting it to nil in order to avoid v1alpha3 --> v1alpha4 --> v1alpha3 round trip errors.
	obj.Spec.AMI.ARN = nil
	obj.Spec.AMI.Filters = nil
}

func AWSMachineTemplateFuzzer(obj *AWSMachineTemplate, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSMachineTemplate.Spec.Template.Spec.AMI.ARN and AWSMachineTemplate.Spec.Template.Spec.AMI.Filters has been removed in v1alpha4, so setting it to nil in order to avoid v1alpha3 --> v1alpha4 --> v1alpha3 round trip errors.
	obj.Spec.Template.Spec.AMI.ARN = nil
	obj.Spec.Template.Spec.AMI.Filters = nil
}

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(v1alpha4.AddToScheme(scheme)).To(Succeed())

	t.Run("for AWSCluster", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1alpha4.AWSCluster{},
		Spoke:  &AWSCluster{},
	}))

	t.Run("for AWSMachine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1alpha4.AWSMachine{},
		Spoke:       &AWSMachine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSMachineTemplate", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1alpha4.AWSMachineTemplate{},
		Spoke:       &AWSMachineTemplate{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSClusterStaticIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1alpha4.AWSClusterStaticIdentity{},
		Spoke:       &AWSClusterStaticIdentity{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSClusterControllerIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1alpha4.AWSClusterControllerIdentity{},
		Spoke:  &AWSClusterControllerIdentity{},
	}))

	t.Run("for AWSClusterRoleIdentity", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1alpha4.AWSClusterRoleIdentity{},
		Spoke:  &AWSClusterRoleIdentity{},
	}))
}
