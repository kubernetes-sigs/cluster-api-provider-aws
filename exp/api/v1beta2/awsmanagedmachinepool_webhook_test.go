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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
)

func TestAWSManagedMachinePoolDefault(t *testing.T) {
	fargate := &AWSManagedMachinePool{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSManagedMachinePool", utildefaulting.DefaultValidateTest(fargate))
	fargate.Default()
}

func TestAWSManagedMachinePoolValidateCreate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		pool    *AWSManagedMachinePool
		wantErr bool
	}{
		{
			name: "pool requires a EKS Node group name",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "",
				},
			},

			wantErr: true,
		},
		{
			name: "pool with valid EKS Node group name",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},

			wantErr: false,
		},
		{
			name: "pool with valid tags is accepted",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-2",
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			},

			wantErr: false,
		},
		{
			name: "invalid tags are rejected",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid update config",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					UpdateConfig: &UpdateConfig{
						MaxUnavailable: aws.Int(1),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "update config with no values",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					UpdateConfig:     &UpdateConfig{},
				},
			},
			wantErr: true,
		},
		{
			name: "update config with both values",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					UpdateConfig: &UpdateConfig{
						MaxUnavailable:           aws.Int(1),
						MaxUnavailablePercentage: aws.Int(10),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "minSize 0 is accepted",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					Scaling: &ManagedMachinePoolScaling{
						MinSize: pointer.Int32(0),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pool.ValidateCreate()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
		})
	}
}

func TestAWSManagedMachinePoolValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		new     *AWSManagedMachinePool
		old     *AWSManagedMachinePool
		wantErr bool
	}{
		{
			name: "update EKS node groups name is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-2",
				},
			},
			wantErr: true,
		},
		{
			name: "adding tags is accepted",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "adding invalid tags is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "adding update config is accepted",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					UpdateConfig: &UpdateConfig{
						MaxUnavailablePercentage: aws.Int(10),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "removing update config is accepted",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					UpdateConfig: &UpdateConfig{
						MaxUnavailablePercentage: aws.Int(10),
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.new.ValidateUpdate(tt.old.DeepCopy())
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
		})
	}
}
