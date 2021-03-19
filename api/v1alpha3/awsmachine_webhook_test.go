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
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestAWSMachine_Create(t *testing.T) {
	tests := []struct {
		name    string
		machine *AWSMachine
		wantErr bool
	}{
		{
			name: "ensure IOPS exists if type equal to io1",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Type: "io1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io2",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Type: "io2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume has no device name",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						DeviceName: "name",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ensure non root volume have device names",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []*Volume{
						{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ensure ensure IOPS exists if type equal to io1 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []*Volume{
						{
							DeviceName: "name",
							Type:       "io1",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ensure ensure IOPS exists if type equal to io2 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []*Volume{
						{
							DeviceName: "name",
							Type:       "io2",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "additional security groups may have id",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: aws.String("id"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "additional security groups may have filters",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							Filters: []Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "additional security groups can't have both id and filters",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: aws.String("id"),
							Filters: []Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := tt.machine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, machine); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			testEnv.Delete(ctx, machine)
		})
	}
}

func TestAWSMachine_Update(t *testing.T) {
	tests := []struct {
		name       string
		oldMachine *AWSMachine
		newMachine *AWSMachine
		wantErr    bool
	}{
		{
			name: "change in providerid, cloudinit, tags and securitygroups",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID: pointer.StringPtr("ID"),
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: pointer.StringPtr("ID"),
						},
					},
					CloudInit: CloudInit{
						SecretPrefix: "test",
						SecretCount:  5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "change in fields other than providerid, tags and securitygroups",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:               nil,
					AdditionalTags:           nil,
					AdditionalSecurityGroups: nil,
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ImageLookupOrg: "test",
					InstanceType:   "test",
					ProviderID:     pointer.StringPtr("ID"),
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: pointer.StringPtr("ID"),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			machine := tt.oldMachine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			if err := testEnv.Create(ctx, machine); err != nil {
				t.Errorf("failed to create machine: %v", err)
			}
			machine.Spec = tt.newMachine.Spec
			if err := testEnv.Update(ctx, machine); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSMachine_SecretsBackend(t *testing.T) {
	baseMachine := &AWSMachine{
		Spec: AWSMachineSpec{
			ProviderID:               nil,
			AdditionalTags:           nil,
			AdditionalSecurityGroups: nil,
		},
	}

	tests := []struct {
		name                   string
		cloudInit              CloudInit
		expectedSecretsBackend string
	}{
		{
			name:                   "with insecure skip secrets manager unset",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: false},
			expectedSecretsBackend: "secrets-manager",
		},
		{
			name:                   "with insecure skip secrets manager unset and secrets backend set",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: false, SecureSecretsBackend: "ssm-parameter-store"},
			expectedSecretsBackend: "ssm-parameter-store",
		},
		{
			name:                   "with insecure skip secrets manager set",
			cloudInit:              CloudInit{InsecureSkipSecretsManager: true},
			expectedSecretsBackend: "",
		},
	}

	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			machine := baseMachine.DeepCopy()
			machine.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "machine-",
				Namespace:    "default",
			}
			machine.Spec.CloudInit = tt.cloudInit
			if err := testEnv.Create(ctx, machine); err != nil {
				t.Errorf("failed to create machine: %v", err)
			}
			g := NewWithT(t)
			g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(SecretBackend(tt.expectedSecretsBackend)))
		})
	}
}
