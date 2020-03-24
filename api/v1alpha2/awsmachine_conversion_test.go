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

func TestConvertAWSMachine(t *testing.T) {
	g := NewWithT(t)

	t.Run("from hub", func(t *testing.T) {
		t.Run("should restore SecretARN, assuming old version of object without field", func(t *testing.T) {
			src := &infrav1alpha3.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1alpha3.AWSMachineSpec{
					CloudInit: infrav1alpha3.CloudInit{
						InsecureSkipSecretsManager: true,
						SecretPrefix:               "something",
					},
				},
			}
			dst := &AWSMachine{}
			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &infrav1alpha3.AWSMachine{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())
			g.Expect(restored.Spec.CloudInit.SecretPrefix).To(Equal(src.Spec.CloudInit.SecretPrefix))
			g.Expect(restored.Spec.CloudInit.InsecureSkipSecretsManager).To(Equal(src.Spec.CloudInit.InsecureSkipSecretsManager))
		})
		t.Run("should convert rootVolume to rootDeviceSize", func(t *testing.T) {
			src := &infrav1alpha3.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: infrav1alpha3.AWSMachineSpec{
					RootVolume: &infrav1alpha3.RootVolume{
						Size:      10,
						Encrypted: true,
					},
				},
			}
			dst := &AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			g.Expect(dst.Spec.RootDeviceSize).To(Equal(int64(10)))
		})
		t.Run("should preserve fields", func(t *testing.T) {
			src := &infrav1alpha3.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-1",
					Annotations: map[string]string{},
				},
				Spec: infrav1alpha3.AWSMachineSpec{
					RootVolume: &infrav1alpha3.RootVolume{
						Size:      10,
						Encrypted: true,
					},
				},
			}
			dst := &AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &infrav1alpha3.AWSMachine{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())

			// Test field restored fields.
			g.Expect(restored).To(Equal(src))
		})
	})

	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert rootDeviceSize to RootVolume", func(t *testing.T) {
			src := &AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
				Spec: AWSMachineSpec{
					RootDeviceSize: 10,
				},
			}
			dst := &infrav1alpha3.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-1",
				},
			}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.RootVolume.Size).To(Equal(int64(10)))
		})

	})

	t.Run("should prefer newer cloudinit data on the v1alpha2 obj", func(t *testing.T) {
		src := &infrav1alpha3.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: infrav1alpha3.AWSMachineSpec{
				CloudInit: infrav1alpha3.CloudInit{
					SecretPrefix: "something",
				},
			},
		}
		dst := &AWSMachine{
			Spec: AWSMachineSpec{
				CloudInit: &CloudInit{
					EnableSecureSecretsManager: true,
					SecretPrefix:               "something-else",
				},
			},
		}
		g.Expect(dst.ConvertFrom(src)).To(Succeed())
		restored := &infrav1alpha3.AWSMachine{}
		g.Expect(dst.ConvertTo(restored)).To(Succeed())
		g.Expect(restored.Spec.CloudInit.SecretPrefix).To(Equal(src.Spec.CloudInit.SecretPrefix))
	})
	t.Run("should restore ImageLookupBaseOS", func(t *testing.T) {
		src := &infrav1alpha3.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: infrav1alpha3.AWSMachineSpec{
				ImageLookupBaseOS: "amazon-linux",
			},
		}
		dst := &AWSMachine{}
		g.Expect(dst.ConvertFrom(src)).To(Succeed())
		restored := &infrav1alpha3.AWSMachine{}
		g.Expect(dst.ConvertTo(restored)).To(Succeed())
		g.Expect(restored.Spec.ImageLookupBaseOS).To(Equal(src.Spec.ImageLookupBaseOS))
	})
	t.Run("should restore UncompressedUserData", func(t *testing.T) {
		flag := true
		src := &infrav1alpha3.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: infrav1alpha3.AWSMachineSpec{
				UncompressedUserData: &flag,
			},
		}
		dst := &AWSMachine{}
		g.Expect(dst.ConvertFrom(src)).To(Succeed())
		restored := &infrav1alpha3.AWSMachine{}
		g.Expect(dst.ConvertTo(restored)).To(Succeed())
		g.Expect(restored.Spec.UncompressedUserData).To(Equal(src.Spec.UncompressedUserData))
	})
}
