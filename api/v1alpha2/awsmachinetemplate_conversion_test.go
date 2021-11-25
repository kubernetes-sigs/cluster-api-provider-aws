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

package v1alpha2

import (
	"testing"

	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

func TestConvertAWSMachineTemplate(t *testing.T) {
	g := gomega.NewWithT(t)
	true := true

	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert rootDeviceSize to rootVolume", func(t *testing.T) {
			src := &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							RootDeviceSize: 10,
						},
					},
				},
			}
			dst := &infrav1beta1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(src.ConvertTo(dst)).To(gomega.Succeed())
			g.Expect(dst.Spec.Template.Spec.RootVolume.Size).To(gomega.Equal(int64(10)))
		})
	})

	t.Run("from hub", func(t *testing.T) {
		t.Run("should convert rootVolume to rootDeviceSize", func(t *testing.T) {
			src := &infrav1beta1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: infrav1beta1.AWSMachineTemplateSpec{
					Template: infrav1beta1.AWSMachineTemplateResource{
						Spec: infrav1beta1.AWSMachineSpec{
							RootVolume: &infrav1beta1.Volume{
								Size:      10,
								Encrypted: &true,
							},
						},
					},
				},
			}
			dst := &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(dst.ConvertFrom(src)).To(gomega.Succeed())
			g.Expect(dst.Spec.Template.Spec.RootDeviceSize).To(gomega.Equal(int64(10)))
		})

		t.Run("should preserve fields", func(t *testing.T) {
			src := &infrav1beta1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-1",
					Annotations: map[string]string{},
				},
				Spec: infrav1beta1.AWSMachineTemplateSpec{
					Template: infrav1beta1.AWSMachineTemplateResource{
						Spec: infrav1beta1.AWSMachineSpec{
							RootVolume: &infrav1beta1.Volume{
								Size:      10,
								Encrypted: &true,
							},
						},
					},
				},
			}
			dst := &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(dst.ConvertFrom(src)).To(gomega.Succeed())
			restored := &infrav1beta1.AWSMachineTemplate{}
			g.Expect(dst.ConvertTo(restored)).To(gomega.Succeed())

			// Test field restored fields.
			g.Expect(restored).To(gomega.Equal(src))
		})
	})
}
