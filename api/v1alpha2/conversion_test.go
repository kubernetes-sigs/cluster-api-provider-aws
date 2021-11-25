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

package v1alpha2

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	v1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(v1beta1.AddToScheme(scheme)).To(Succeed())

	t.Run("for AWSMachine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1beta1.AWSMachine{},
		Spoke:       &AWSMachine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs, CustomObjectMetaFuzzFunc},
	}))
	t.Run("for AWSMachineTemplate", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme:      scheme,
		Hub:         &v1beta1.AWSMachineTemplate{},
		Spoke:       &AWSMachineTemplate{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{fuzzFuncs, CustomObjectMetaFuzzFunc},
	}))
}
func fuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		AWSMachineFuzzer,
		AWSMachineTemplateFuzzer,
		AMIReferenceFuzzer,
		Beta1AWSMachineTemplateFuzzer,
	}
}

func AWSMachineFuzzer(obj *AWSMachine, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSMachine.Spec.AMI.ARN and AWSMachine.Spec.AMI.Filters has been removed in v1beta1, so setting it to nil in order to avoid v1alpha3 --> <hub> --> v1alpha3 round trip errors.
	obj.Spec.AMI.ARN = nil
	obj.Spec.AMI.Filters = nil
}

func AWSMachineTemplateFuzzer(obj *AWSMachineTemplate, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// AWSMachineTemplate.Spec.Template.Spec.AMI.ARN and AWSMachineTemplate.Spec.Template.Spec.AMI.Filters has been removed in v1beta1, so setting it to nil in order to avoid v1alpha3 --> v1beta1 --> v1alpha3 round trip errors.
	obj.Spec.Template.Spec.AMI.ARN = nil
	obj.Spec.Template.Spec.AMI.Filters = nil
}

func Beta1AWSMachineTemplateFuzzer(obj *v1beta1.AWSMachineTemplate, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.Spec.Template.Spec.NonRootVolumes = nil
}

func AMIReferenceFuzzer(obj *v1beta1.AMIReference, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.EKSOptimizedLookupType = nil
}

func CustomObjectMetaFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		ObjectMetaFuzzer,
		ClusterObjectMetaFuzzer,
	}
}

func ObjectMetaFuzzer(in *metav1.ObjectMeta, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	in.Annotations = make(map[string]string)
	in.Labels = make(map[string]string)
}

func ClusterObjectMetaFuzzer(in *clusterv1.ObjectMeta, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	in.Annotations = make(map[string]string)
	in.Labels = make(map[string]string)
}
