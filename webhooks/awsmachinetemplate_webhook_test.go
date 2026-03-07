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

	"github.com/aws/aws-sdk-go-v2/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func TestAWSMachineTemplateValidateCreate(t *testing.T) {
	tests := []struct {
		name          string
		inputTemplate *infrav1.AWSMachineTemplate
		wantError     bool
	}{
		{
			name: "don't allow providerID",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							ProviderID: ptr.To[string]("something"),
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "don't allow secretARN",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							CloudInit: infrav1.CloudInit{
								SecretPrefix: "something",
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "ensure RootVolume DeviceName can be set for use with clusterctl move",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							RootVolume: &infrav1.Volume{
								DeviceName: "name",
								Type:       "gp2",
								Size:       *aws.Int64(8),
							},
							InstanceType: "test",
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "hostID and dynamicHostAllocation are mutually exclusive",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "host",
							HostID:       aws.String("h-1234567890abcdef0"),
							DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
								Tags: map[string]string{
									"Environment": "test",
								},
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "hostAffinity=host requires hostID or dynamicHostAllocation",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "host",
							HostAffinity: ptr.To("host"),
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "hostAffinity=host with hostID is valid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "host",
							HostAffinity: ptr.To("host"),
							HostID:       ptr.To("h-09dcf61cb388b0149"),
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "hostAffinity=host with dynamicHostAllocation is valid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "host",
							HostAffinity: ptr.To("host"),
							DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
								Tags: map[string]string{"env": "test"},
							},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "hostAffinity=default without hostID is valid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							HostAffinity: ptr.To("default"),
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "hostID without tenancy=host is invalid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "default",
							HostID:       ptr.To("h-09dcf61cb388b0149"),
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "hostAffinity=host without tenancy=host is invalid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "default",
							HostAffinity: ptr.To("host"),
							HostID:       ptr.To("h-09dcf61cb388b0149"),
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "dynamicHostAllocation without tenancy=host is invalid",
			inputTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test",
							Tenancy:      "dedicated",
							DynamicHostAllocation: &infrav1.DynamicHostAllocationSpec{
								Tags: map[string]string{"env": "test"},
							},
						},
					},
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := tt.inputTemplate.DeepCopy()
			template.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "template-",
				Namespace:    "default",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, template); (err != nil) != tt.wantError {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestAWSMachineTemplateValidateUpdate(t *testing.T) {
	tests := []struct {
		name             string
		modifiedTemplate *infrav1.AWSMachineTemplate
		wantError        bool
	}{
		{
			name: "don't allow updates",
			modifiedTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							InstanceType: "test2",
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "allow defaulted values to update",
			modifiedTemplate: &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							CloudInit:    infrav1.CloudInit{},
							InstanceType: "test",
							InstanceMetadataOptions: &infrav1.InstanceMetadataOptions{
								HTTPEndpoint:            infrav1.InstanceMetadataEndpointStateEnabled,
								HTTPProtocolIPv6:        infrav1.InstanceMetadataEndpointStateDisabled,
								HTTPPutResponseHopLimit: 1,
								HTTPTokens:              infrav1.HTTPTokensStateOptional,
								InstanceMetadataTags:    infrav1.InstanceMetadataEndpointStateDisabled,
							},
						},
					},
				},
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			template := &infrav1.AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "template-",
					Namespace:    "default",
				},
				Spec: infrav1.AWSMachineTemplateSpec{
					Template: infrav1.AWSMachineTemplateResource{
						Spec: infrav1.AWSMachineSpec{
							CloudInit:    infrav1.CloudInit{},
							InstanceType: "test",
						},
					},
				},
			}

			if err := testEnv.Create(ctx, template); err != nil {
				t.Errorf("failed to create template: %v", err)
			}
			template.Spec = tt.modifiedTemplate.Spec
			if err := testEnv.Update(ctx, template); (err != nil) != tt.wantError {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantError)
			}
		},
		)
	}
}
