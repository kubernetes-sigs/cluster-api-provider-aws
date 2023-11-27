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

var (
	oldDiskSize                   = int32(50)
	newDiskSize                   = int32(100)
	oldAmiType                    = Al2x86_64
	newAmiType                    = Al2x86_64GPU
	oldCapacityType               = ManagedMachinePoolCapacityTypeOnDemand
	newCapacityType               = ManagedMachinePoolCapacityTypeSpot
	oldAvailabilityZoneSubnetType = AZSubnetTypePublic
	newAvailabilityZoneSubnetType = AZSubnetTypePrivate
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
		{
			name: "With Launch Template/AMI type and ami id both are not specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AMIType:          nil,
					AWSLaunchTemplate: &AWSLaunchTemplate{
						AMI: infrav1.AMIReference{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "With Launch Template/AMI type not specific, ami id specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AMIType:          nil,
					AWSLaunchTemplate: &AWSLaunchTemplate{
						AMI: infrav1.AMIReference{
							ID: aws.String("test-ami-id"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "With Launch Template/AMI type not specific, ami EKSAMILookupType specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AMIType:          nil,
					AWSLaunchTemplate: &AWSLaunchTemplate{
						AMI: infrav1.AMIReference{
							EKSOptimizedLookupType: infrav1.AmazonLinux.GetPtr(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "With Launch Template/AMI type specific, ami not specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:  "eks-node-group-3",
					AMIType:           Al2x86_64.GetPtr(),
					AWSLaunchTemplate: &AWSLaunchTemplate{},
				},
			},
			wantErr: false,
		},
		{
			name: "With Launch Template/AMI type and ami id are specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AMIType:          Al2x86_64.GetPtr(),
					AWSLaunchTemplate: &AWSLaunchTemplate{
						AMI: infrav1.AMIReference{
							ID: aws.String("test-ami-id"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "With Launch Template/AMI type and ami EKSAMILookupType are specific",
			pool: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-3",
					AMIType:          Al2x86_64.GetPtr(),
					AWSLaunchTemplate: &AWSLaunchTemplate{
						AMI: infrav1.AMIReference{
							EKSOptimizedLookupType: infrav1.AmazonLinux.GetPtr(),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := tt.pool.ValidateCreate()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
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
		{
			name: "adding subnet id is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					SubnetIDs:        []string{"subnet-1"},
				},
			},
			wantErr: true,
		},
		{
			name: "removing subnet id is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					SubnetIDs:        []string{"subnet-1"},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing subnet id is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					SubnetIDs:        []string{"subnet-1"},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					SubnetIDs:        []string{"subnet-2"},
				},
			},
			wantErr: true,
		},
		{
			name: "removing role name is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RoleName:         "role-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing role name is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RoleName:         "role-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RoleName:         "role-2",
				},
			},
			wantErr: true,
		},
		{
			name: "adding disk size is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					DiskSize:         &newDiskSize,
				},
			},
			wantErr: true,
		},
		{
			name: "removing disk size is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					DiskSize:         &oldDiskSize,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing disk size is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					DiskSize:         &oldDiskSize,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					DiskSize:         &newDiskSize,
				},
			},
			wantErr: true,
		},
		{
			name: "adding ami type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          &newAmiType,
				},
			},
			wantErr: true,
		},
		{
			name: "removing ami type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          &oldAmiType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing ami type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          &oldAmiType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          &newAmiType,
				},
			},
			wantErr: true,
		},
		{
			name: "adding remote access is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RemoteAccess: &ManagedRemoteAccess{
						Public: false,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "removing remote access is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RemoteAccess: &ManagedRemoteAccess{
						Public: false,
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing remote access is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RemoteAccess: &ManagedRemoteAccess{
						Public: false,
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					RemoteAccess: &ManagedRemoteAccess{
						Public: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "removing capacity type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					CapacityType:     &oldCapacityType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing capacity type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					CapacityType:     &oldCapacityType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					CapacityType:     &newCapacityType,
				},
			},
			wantErr: true,
		},
		{
			name: "adding availability zones is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:  "eks-node-group-1",
					AvailabilityZones: []string{"us-east-1a"},
				},
			},
			wantErr: true,
		},
		{
			name: "removing availability zones is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:  "eks-node-group-1",
					AvailabilityZones: []string{"us-east-1a"},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing availability zones is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:  "eks-node-group-1",
					AvailabilityZones: []string{"us-east-1a"},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:  "eks-node-group-1",
					AvailabilityZones: []string{"us-east-1a", "us-east-1b"},
				},
			},
			wantErr: true,
		},
		{
			name: "adding availability zone subnet type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:           "eks-node-group-1",
					AvailabilityZoneSubnetType: &newAvailabilityZoneSubnetType,
				},
			},
			wantErr: true,
		},
		{
			name: "removing availability zone subnet type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:           "eks-node-group-1",
					AvailabilityZoneSubnetType: &oldAvailabilityZoneSubnetType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing availability zone subnet type is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:           "eks-node-group-1",
					AvailabilityZoneSubnetType: &oldAvailabilityZoneSubnetType,
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName:           "eks-node-group-1",
					AvailabilityZoneSubnetType: &newAvailabilityZoneSubnetType,
				},
			},
			wantErr: true,
		},
		{
			name: "adding launch template is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "removing launch template is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name: "test",
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
				},
			},
			wantErr: true,
		},
		{
			name: "changing launch template name is rejected",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name: "test",
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name: "test2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "changing launch template fields other than name is accepted",
			old: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          Al2x86_64.GetPtr(),
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name:              "test",
						ImageLookupFormat: "test",
					},
				},
			},
			new: &AWSManagedMachinePool{
				Spec: AWSManagedMachinePoolSpec{
					EKSNodegroupName: "eks-node-group-1",
					AMIType:          Al2x86_64.GetPtr(),
					AWSLaunchTemplate: &AWSLaunchTemplate{
						Name:              "test",
						ImageLookupFormat: "test2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := tt.new.ValidateUpdate(tt.old.DeepCopy())
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}
