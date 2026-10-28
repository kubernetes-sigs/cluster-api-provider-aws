/*
Copyright 2026 The Kubernetes Authors.

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

package webhooks

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

func TestAWSMachinePoolTemplateDefault(t *testing.T) {
	g := NewWithT(t)

	template := &expinfrav1.AWSMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
	}
	t.Run("for AWSMachinePoolTemplate", utildefaulting.DefaultValidateTest(context.Background(), template, &AWSMachinePoolTemplate{}))

	err := (&AWSMachinePoolTemplate{}).Default(context.Background(), template)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(template.Spec.Template.Spec.DefaultCoolDown.Duration).To(Equal(300 * time.Second))
	g.Expect(template.Spec.Template.Spec.DefaultInstanceWarmup.Duration).To(Equal(300 * time.Second))
}

func TestAWSMachinePoolTemplateValidateCreate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name     string
		template *expinfrav1.AWSMachinePoolTemplate
		wantErr  bool
	}{
		{
			name: "empty spec is valid",
			template: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "providerID is rejected",
			template: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							ProviderID: "aws:///us-east-1a/asg-name",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "providerIDList is rejected",
			template: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							ProviderIDList: []string{"aws:///us-east-1a/i-01234567890123456"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "root volume io1 without iops is rejected",
			template: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
								RootVolume: &infrav1.Volume{Type: infrav1.VolumeTypeIO1},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "subnet with both ID and filters is rejected",
			template: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							Subnets: []infrav1.AWSResourceReference{
								{ID: aws.String("subnet-1"), Filters: []infrav1.Filter{{Name: "tag", Values: []string{"v"}}}},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			webhook := &AWSMachinePoolTemplate{}
			warn, err := webhook.ValidateCreate(context.Background(), tt.template)

			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(warn).To(BeEmpty())
		})
	}
}

func TestAWSMachinePoolTemplateValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		old     *expinfrav1.AWSMachinePoolTemplate
		new     *expinfrav1.AWSMachinePoolTemplate
		wantErr bool
	}{
		{
			name: "arbitrary spec changes are allowed (no immutability)",
			old: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MinSize: 0,
							MaxSize: 1,
							AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
								InstanceType: "t3.medium",
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MinSize:           1,
							MaxSize:           5,
							AvailabilityZones: []string{"us-east-1a"},
							AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
								InstanceType: "m5.large",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "setting providerIDList on update is rejected",
			old: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{},
					},
				},
			},
			new: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							ProviderIDList: []string{"aws:///us-east-1a/i-01234567890123456"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid refreshPreferences on update is rejected",
			old: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{},
					},
				},
			},
			new: &expinfrav1.AWSMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: expinfrav1.AWSMachinePoolTemplateSpec{
					Template: expinfrav1.AWSMachinePoolTemplateResource{
						Spec: expinfrav1.AWSMachinePoolSpec{
							RefreshPreferences: &expinfrav1.RefreshPreferences{
								MaxHealthyPercentage: aws.Int64(110),
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			webhook := &AWSMachinePoolTemplate{}
			warn, err := webhook.ValidateUpdate(context.Background(), tt.old, tt.new)

			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(warn).To(BeEmpty())
		})
	}
}

func TestAWSMachinePoolTemplateSpecValidationEnforcedOnUpdate(t *testing.T) {
	g := NewWithT(t)

	template := &expinfrav1.AWSMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: expinfrav1.AWSMachinePoolTemplateSpec{
			Template: expinfrav1.AWSMachinePoolTemplateResource{
				Spec: expinfrav1.AWSMachinePoolSpec{
					AWSLaunchTemplate: expinfrav1.AWSLaunchTemplate{
						RootVolume: &infrav1.Volume{Type: infrav1.VolumeTypeIO1}, // io1 without IOPS: invalid
					},
				},
			},
		},
	}

	webhook := &AWSMachinePoolTemplate{}

	// The same spec validation applies to create and update: an invalid spec is
	// rejected either way, even though fields remain mutable across updates.
	_, err := webhook.ValidateCreate(context.Background(), template)
	g.Expect(err).To(HaveOccurred())

	_, err = webhook.ValidateUpdate(context.Background(), template.DeepCopy(), template)
	g.Expect(err).To(HaveOccurred())
}
