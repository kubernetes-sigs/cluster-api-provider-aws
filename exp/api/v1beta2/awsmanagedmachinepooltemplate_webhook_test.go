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

package v1beta2

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

func TestAWSManagedMachinePoolTemplateDefault(t *testing.T) {
	g := NewWithT(t)

	template := &AWSManagedMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
	}
	t.Run("for AWSManagedMachinePoolTemplate", utildefaulting.DefaultValidateTest(context.Background(), template, &awsManagedMachinePoolTemplateWebhook{}))

	err := (&awsManagedMachinePoolTemplateWebhook{}).Default(context.Background(), template)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(template.Spec.Template.Spec.UpdateConfig).NotTo(BeNil())
	g.Expect(template.Spec.Template.Spec.UpdateConfig.MaxUnavailable).To(Equal(ptr.To[int](1)))
}

func TestAWSManagedMachinePoolTemplateDefault_EKSNodegroupNameNotGenerated(t *testing.T) {
	g := NewWithT(t)

	template := &AWSManagedMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: AWSManagedMachinePoolTemplateSpec{
			Template: AWSManagedMachinePoolTemplateResource{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "",
				},
			},
		},
	}

	err := (&awsManagedMachinePoolTemplateWebhook{}).Default(context.Background(), template)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(template.Spec.Template.Spec.EKSNodegroupName).To(Equal(""))
}

func TestAWSManagedMachinePoolTemplateValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		old     *AWSManagedMachinePoolTemplate
		new     *AWSManagedMachinePoolTemplate
		wantErr bool
	}{
		{
			name: "metadata changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"old": "label",
					},
				},
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"new": "label",
					},
				},
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{},
					},
				},
			},
			wantErr: false,
		},
		// Mutable fields - changes should be allowed
		{
			name: "mutable field - scaling changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Scaling: &ManagedMachinePoolScaling{
								MinSize: ptr.To[int32](1),
								MaxSize: ptr.To[int32](5),
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Scaling: &ManagedMachinePoolScaling{
								MinSize: ptr.To[int32](1),
								MaxSize: ptr.To[int32](10),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - updateConfig changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							UpdateConfig: &UpdateConfig{
								MaxUnavailable: ptr.To[int](1),
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							UpdateConfig: &UpdateConfig{
								MaxUnavailable: ptr.To[int](2),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - labels changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Labels: map[string]string{
								"old": "label",
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Labels: map[string]string{
								"new": "label",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - taints changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Taints: Taints{
								{
									Key:    "key1",
									Value:  "value1",
									Effect: "NoSchedule",
								},
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							Taints: Taints{
								{
									Key:    "key2",
									Value:  "value2",
									Effect: "NoSchedule",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - instanceType changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							InstanceType: ptr.To[string]("t3.medium"),
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							InstanceType: ptr.To[string]("t3.large"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - lifecycleHooks changes are allowed",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLifecycleHooks: []AWSLifecycleHook{
								{
									Name:                "hook1",
									LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
								},
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLifecycleHooks: []AWSLifecycleHook{
								{
									Name:                "hook2",
									LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		// Immutable fields - changes should be rejected
		{
			name: "immutable field - eksNodegroupName change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							EKSNodegroupName: "old-name",
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							EKSNodegroupName: "new-name",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - subnetIDs change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							SubnetIDs: []string{"subnet-1"},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							SubnetIDs: []string{"subnet-2"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - roleName change is rejected once set",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							RoleName: "old-role",
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							RoleName: "new-role",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - diskSize change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							DiskSize: ptr.To[int32](100),
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							DiskSize: ptr.To[int32](200),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - amiType change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AMIType: ptr.To[ManagedMachineAMIType](Al2x86_64),
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AMIType: ptr.To[ManagedMachineAMIType](Al2x86_64GPU),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - remoteAccess change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							RemoteAccess: &ManagedRemoteAccess{
								SSHKeyName: ptr.To[string]("key1"),
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							RemoteAccess: &ManagedRemoteAccess{
								SSHKeyName: ptr.To[string]("key2"),
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - capacityType change is rejected once set",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							CapacityType: ptr.To[ManagedMachinePoolCapacityType](ManagedMachinePoolCapacityTypeOnDemand),
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							CapacityType: ptr.To[ManagedMachinePoolCapacityType](ManagedMachinePoolCapacityTypeSpot),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - availabilityZones change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AvailabilityZones: []string{"us-east-1a"},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AvailabilityZones: []string{"us-east-1b"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - availabilityZoneSubnetType change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AvailabilityZoneSubnetType: ptr.To[AZSubnetType](AZSubnetTypePublic),
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AvailabilityZoneSubnetType: ptr.To[AZSubnetType](AZSubnetTypePrivate),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - adding awsLaunchTemplate is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name: "template",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - removing awsLaunchTemplate is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name: "template",
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - awsLaunchTemplate.name change is rejected",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name: "old-template",
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name: "new-template",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "mutable launchTemplate fields - other fields can change",
			old: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name:          "template",
								VersionNumber: ptr.To[int64](1),
							},
						},
					},
				},
			},
			new: &AWSManagedMachinePoolTemplate{
				Spec: AWSManagedMachinePoolTemplateSpec{
					Template: AWSManagedMachinePoolTemplateResource{
						Spec: AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &AWSLaunchTemplate{
								Name:          "template",
								VersionNumber: ptr.To[int64](2),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := (&awsManagedMachinePoolTemplateWebhook{}).ValidateUpdate(context.Background(), tt.old, tt.new)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
			g.Expect(warn).To(BeEmpty())
		})
	}
}
