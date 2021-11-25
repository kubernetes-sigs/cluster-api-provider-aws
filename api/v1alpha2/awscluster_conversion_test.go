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
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

func TestConvertAWSCluster(t *testing.T) {
	g := NewWithT(t)

	t.Run("from hub", func(t *testing.T) {
		t.Run("should restore SSHKeyName, retaining a nil value", func(t *testing.T) {
			src := &infrav1beta1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1beta1.AWSClusterSpec{
					SSHKeyName: nil,
				},
			}
			dst := &AWSCluster{}
			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &infrav1beta1.AWSCluster{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())
			g.Expect(restored.Spec.SSHKeyName).To(BeNil())
		})

		t.Run("should convert subnets to pointers", func(t *testing.T) {
			src := &infrav1beta1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1beta1.AWSClusterSpec{
					NetworkSpec: infrav1beta1.NetworkSpec{
						Subnets: infrav1beta1.Subnets{
							infrav1beta1.SubnetSpec{
								ID: "test",
							},
						},
					},
				},
			}
			dst := &AWSCluster{}
			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			g.Expect(dst.Spec.NetworkSpec.Subnets).ToNot(BeEmpty())
			g.Expect(dst.Spec.NetworkSpec.Subnets[0].ID).To(Equal("test"))
		})
	})

	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert subnets from pointers", func(t *testing.T) {
			src := &AWSCluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						Subnets: Subnets{
							&SubnetSpec{
								ID: "test",
							},
						},
					},
				},
			}
			dst := &infrav1beta1.AWSCluster{}
			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.NetworkSpec.Subnets).ToNot(BeEmpty())
			g.Expect(dst.Spec.NetworkSpec.Subnets[0].ID).To(Equal("test"))
		})
	})
}
