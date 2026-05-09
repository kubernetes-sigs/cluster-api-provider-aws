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

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

func TestAWSManagedMachinePoolTemplateDefault(t *testing.T) {
	g := NewWithT(t)

	template := &expinfrav1.AWSManagedMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
	}
	t.Run("for AWSManagedMachinePoolTemplate", utildefaulting.DefaultValidateTest(context.Background(), template, &AWSManagedMachinePoolTemplate{}))

	err := (&AWSManagedMachinePoolTemplate{}).Default(context.Background(), template)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(template.Spec.Template.Spec.UpdateConfig).NotTo(BeNil())
	g.Expect(template.Spec.Template.Spec.UpdateConfig.MaxUnavailable).To(Equal(ptr.To[int](1)))
}

func TestAWSManagedMachinePoolTemplateDefault_EKSNodegroupNameNotGenerated(t *testing.T) {
	g := NewWithT(t)

	template := &expinfrav1.AWSManagedMachinePoolTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
			Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
				Spec: expinfrav1.AWSManagedMachinePoolSpec{
					EKSNodegroupName: "",
				},
			},
		},
	}

	err := (&AWSManagedMachinePoolTemplate{}).Default(context.Background(), template)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(template.Spec.Template.Spec.EKSNodegroupName).To(Equal(""))
}

func TestAWSManagedMachinePoolTemplateValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		old     *expinfrav1.AWSManagedMachinePoolTemplate
		new     *expinfrav1.AWSManagedMachinePoolTemplate
		wantErr bool
	}{
		{
			name: "metadata changes are allowed",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"old": "label",
					},
				},
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"new": "label",
					},
				},
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{},
					},
				},
			},
			wantErr: false,
		},
		// Mutable fields - changes should be allowed
		{
			name: "mutable field - scaling changes are allowed",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							Scaling: &expinfrav1.ManagedMachinePoolScaling{
								MinSize: ptr.To[int32](1),
								MaxSize: ptr.To[int32](5),
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							Scaling: &expinfrav1.ManagedMachinePoolScaling{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							UpdateConfig: &expinfrav1.UpdateConfig{
								MaxUnavailable: ptr.To[int](1),
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							UpdateConfig: &expinfrav1.UpdateConfig{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							Labels: map[string]string{
								"old": "label",
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							Taints: expinfrav1.Taints{
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
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							Taints: expinfrav1.Taints{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							InstanceType: ptr.To[string]("t3.medium"),
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							InstanceType: ptr.To[string]("t3.large"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mutable field - lifecycleHooks changes are allowed",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLifecycleHooks: []expinfrav1.AWSLifecycleHook{
								{
									Name:                "hook1",
									LifecycleTransition: expinfrav1.LifecycleHookTransitionInstanceLaunching,
								},
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLifecycleHooks: []expinfrav1.AWSLifecycleHook{
								{
									Name:                "hook2",
									LifecycleTransition: expinfrav1.LifecycleHookTransitionInstanceLaunching,
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							EKSNodegroupName: "old-name",
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							EKSNodegroupName: "new-name",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - subnetIDs change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							SubnetIDs: []string{"subnet-1"},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							SubnetIDs: []string{"subnet-2"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - roleName change is rejected once set",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							RoleName: "old-role",
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							RoleName: "new-role",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - diskSize change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							DiskSize: ptr.To[int32](100),
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							DiskSize: ptr.To[int32](200),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - amiType change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AMIType: ptr.To[expinfrav1.ManagedMachineAMIType](expinfrav1.Al2x86_64),
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AMIType: ptr.To[expinfrav1.ManagedMachineAMIType](expinfrav1.Al2x86_64GPU),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - remoteAccess change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							RemoteAccess: &expinfrav1.ManagedRemoteAccess{
								SSHKeyName: ptr.To[string]("key1"),
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							RemoteAccess: &expinfrav1.ManagedRemoteAccess{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							CapacityType: ptr.To[expinfrav1.ManagedMachinePoolCapacityType](expinfrav1.ManagedMachinePoolCapacityTypeOnDemand),
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							CapacityType: ptr.To[expinfrav1.ManagedMachinePoolCapacityType](expinfrav1.ManagedMachinePoolCapacityTypeSpot),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - availabilityZones change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AvailabilityZones: []string{"us-east-1a"},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AvailabilityZones: []string{"us-east-1b"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - availabilityZoneSubnetType change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AvailabilityZoneSubnetType: ptr.To[expinfrav1.AZSubnetType](expinfrav1.AZSubnetTypePublic),
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AvailabilityZoneSubnetType: ptr.To[expinfrav1.AZSubnetType](expinfrav1.AZSubnetTypePrivate),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - adding awsLaunchTemplate is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
								Name: "template",
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "immutable field - awsLaunchTemplate.name change is rejected",
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
								Name: "old-template",
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
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
			old: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
								Name:          "template",
								VersionNumber: ptr.To[int64](1),
							},
						},
					},
				},
			},
			new: &expinfrav1.AWSManagedMachinePoolTemplate{
				Spec: expinfrav1.AWSManagedMachinePoolTemplateSpec{
					Template: expinfrav1.AWSManagedMachinePoolTemplateResource{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{
							AWSLaunchTemplate: &expinfrav1.AWSLaunchTemplate{
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
			warn, err := (&AWSManagedMachinePoolTemplate{}).ValidateUpdate(context.Background(), tt.old, tt.new)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
			g.Expect(warn).To(BeEmpty())
		})
	}
}
