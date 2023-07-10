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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
)

func TestMachineDefault(t *testing.T) {
	machine := &AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSMachine", utildefaulting.DefaultValidateTest(machine))
	machine.Default()
	g := NewWithT(t)
	g.Expect(machine.Spec.CloudInit.SecureSecretsBackend).To(Equal(SecretBackendSecretsManager))
}

func TestAWSMachineCreate(t *testing.T) {
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
					InstanceType: "test",
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
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume throughput is nonnegative",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						Throughput: aws.Int64(-125),
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure root volume with device name works (for clusterctl move)",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					RootVolume: &Volume{
						DeviceName: "name",
						Type:       "gp2",
						Size:       *aws.Int64(8),
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "ensure non root volume have device names",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io1 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							DeviceName: "name",
							Type:       "io1",
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure IOPS exists if type equal to io2 for non root volumes",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							DeviceName: "name",
							Type:       "io2",
						},
					},
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "ensure non root volume throughput is nonnegative",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					NonRootVolumes: []Volume{
						{
							Throughput: aws.Int64(-125),
						},
					},
					InstanceType: "test",
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
					InstanceType: "test",
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
					InstanceType: "test",
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
					InstanceType: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "valid additional tags are accepted",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalTags: Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					InstanceType: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "empty instance type not allowed",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "",
				},
			},
			wantErr: true,
		},
		{
			name: "instance type minimum length is 2",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					InstanceType: "t",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tags return error",
			machine: &AWSMachine{
				Spec: AWSMachineSpec{
					AdditionalTags: Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
					InstanceType: "test",
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

func TestAWSMachineUpdate(t *testing.T) {
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
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID:   pointer.String("ID"),
					InstanceType: "test",
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: pointer.String("ID"),
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
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ImageLookupOrg: "test",
					InstanceType:   "test",
					ProviderID:     pointer.String("ID"),
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: []AWSResourceReference{
						{
							ID: pointer.String("ID"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change in tags adding invalid ones",
			oldMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: Tags{
						"key-1": "value-1",
					},
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
				},
			},
			newMachine: &AWSMachine{
				Spec: AWSMachineSpec{
					ProviderID: nil,
					AdditionalTags: Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
					AdditionalSecurityGroups: nil,
					InstanceType:             "test",
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

func TestAWSMachineSecretsBackend(t *testing.T) {
	baseMachine := &AWSMachine{
		Spec: AWSMachineSpec{
			ProviderID:               nil,
			AdditionalTags:           nil,
			AdditionalSecurityGroups: nil,
			InstanceType:             "test",
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
