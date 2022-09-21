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

package v1alpha3

import (
	"testing"
	
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	
	"sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func fuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		AWSMachinePoolFuzzer,
	}
}

func AWSMachinePoolFuzzer(obj *AWSMachinePool, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSMachinePool.Spec.AWSLaunchTemplate.AMI.ARN and AWSMachinePool.Spec.AWSLaunchTemplate.AMI.Filters has been removed in v1beta1, so setting it to nil in order to avoid v1alpha3 --> v1beta1 --> v1alpha3 round trip errors.
	obj.Spec.AWSLaunchTemplate.AMI.ARN = nil
	obj.Spec.AWSLaunchTemplate.AMI.Filters = nil
}

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(v1beta1.AddToScheme(scheme)).To(Succeed())

	t.Run("for AWSMachinePool", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1beta1.AWSMachinePool{},
		Spoke:       &AWSMachinePool{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))

	t.Run("for AWSManagedMachinePool", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta1.AWSManagedMachinePool{},
		Spoke:  &AWSManagedMachinePool{},
	}))

	t.Run("for AWSFargateProfile", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &v1beta1.AWSFargateProfile{},
		Spoke:  &AWSFargateProfile{},
	}))
}
