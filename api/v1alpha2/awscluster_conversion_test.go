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

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

func TestConvertAWSCluster(t *testing.T) {
	g := NewWithT(t)

	t.Run("from hub", func(t *testing.T) {
		t.Run("should restore SSHKeyName, retaining a nil value", func(t *testing.T) {
			src := &infrav1alpha3.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1alpha3.AWSClusterSpec{
					SSHKeyName: nil,
				},
			}
			dst := &AWSCluster{}
			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &infrav1alpha3.AWSCluster{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())
			g.Expect(restored.Spec.SSHKeyName).To(BeNil())
		})
	})

}
