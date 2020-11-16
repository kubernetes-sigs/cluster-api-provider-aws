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

package v1alpha3

import (
	"testing"

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
							ProviderID: pointer.StringPtr("something"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inputTemplate.ValidateCreate()
			if (err != nil) != tt.wantError {
				t.Errorf("unexpected result - wanted %+v, got %+v", tt.wantError, err)
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
			name: "don't allow ssm parameter store",
			modifiedTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit: CloudInit{
								SecureSecretsBackend: SecretBackendSSMParameterStore,
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "allow secrets manager",
			modifiedTemplate: &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit: CloudInit{
								SecureSecretsBackend: SecretBackendSecretsManager,
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
			oldTemplate := &AWSMachineTemplate{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: AWSMachineTemplateSpec{
					Template: AWSMachineTemplateResource{
						Spec: AWSMachineSpec{
							CloudInit: CloudInit{},
						},
					},
				},
			}

			err := tt.modifiedTemplate.ValidateUpdate(oldTemplate)
			if (err != nil) != tt.wantError {
				t.Errorf("unexpected result - wanted %+v, got %+v", tt.wantError, err)
			}
		})
	}
}
