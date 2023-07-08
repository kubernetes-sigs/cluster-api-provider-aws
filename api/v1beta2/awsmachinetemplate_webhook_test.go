/*
Copyright 2022 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestAWSMachineTemplateValidateCreate(t *testing.T) {
	tests := []struct {
		name          string
		inputTemplate *AWSMachineTemplate
		wantError     bool
	}{
		{
			name: "don't allow providerID",
			inputTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							ProviderID: pointer.String("something"),
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "don't allow secretARN",
			inputTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit: CloudInit{
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
			inputTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							RootVolume: &Volume{
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
		modifiedTemplate *AWSMachineTemplate
		wantError        bool
	}{
		{
			name: "don't allow updates",
			modifiedTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							InstanceType: "test2",
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "allow defaulted values to update",
			modifiedTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit:    CloudInit{},
							InstanceType: "test",
							InstanceMetadataOptions: &InstanceMetadataOptions{
								HTTPEndpoint:            InstanceMetadataEndpointStateEnabled,
								HTTPPutResponseHopLimit: 1,
								HTTPTokens:              HTTPTokensStateOptional,
								InstanceMetadataTags:    InstanceMetadataEndpointStateDisabled,
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
			template := &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "template-",
					Namespace:    "default",
				},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit:    CloudInit{},
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
