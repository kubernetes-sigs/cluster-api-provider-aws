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

package network

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	"sigs.k8s.io/cluster-api/api/v1beta1"
)

func setupNewManagedControlPlaneScope(cl client.Client) (*scope.ManagedControlPlaneScope, error) {
	return scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:  cl,
		Cluster: &v1beta1.Cluster{},
		ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
			Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				SecondaryCidrBlock: ptr.To[string]("secondary-cidr"),
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{ID: "vpc-id"},
				},
			},
		},
	})
}

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := ekscontrolplanev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func TestServiceAssociateSecondaryCidr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                                    string
		fillAWSManagedControlPlaneSecondaryCIDR bool
		networkSecondaryCIDRBlocks              []infrav1.VpcCidrBlock
		expect                                  func(m *mocks.MockEC2APIMockRecorder)
		wantErr                                 bool
	}{
		{
			name:                                    "Should not associate secondary CIDR if no secondary cidr block info present in control plane",
			fillAWSManagedControlPlaneSecondaryCIDR: false,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				// No calls expected
				m.DescribeVpcsWithContext(context.TODO(), gomock.Any()).Times(0)
				m.AssociateVpcCidrBlockWithContext(context.TODO(), gomock.Any()).Times(0)
			},
			wantErr: false,
		},
		{
			name:                                    "Should return error if unable to describe VPC",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:                                    "Should not associate secondary cidr block if already exist in VPC",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr")},
							},
						},
					}}, nil)
			},
		},
		{
			name:                                    "Should return error if no VPC found",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{},
				}, nil)
			},
			wantErr: true,
		},
		{
			name:                                    "Should return error if failed during associating secondary cidr block",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{},
						},
					}}, nil)
				m.AssociateVpcCidrBlockWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.AssociateVpcCidrBlockInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:                                    "Should successfully associate secondary CIDR block if none is associated yet",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{},
						},
					}}, nil)
				m.AssociateVpcCidrBlockWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.AssociateVpcCidrBlockInput{})).Return(&ec2.AssociateVpcCidrBlockOutput{
					CidrBlockAssociation: &ec2.VpcCidrBlockAssociation{
						AssociationId: ptr.To[string]("association-id-success"),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:                                    "Should successfully associate missing secondary CIDR blocks",
			fillAWSManagedControlPlaneSecondaryCIDR: false,
			networkSecondaryCIDRBlocks: []infrav1.VpcCidrBlock{
				{
					IPv4CidrBlock: "10.0.1.0/24",
				},
				{
					IPv4CidrBlock: "10.0.2.0/24",
				},
				{
					IPv4CidrBlock: "10.0.3.0/24",
				},
				{
					IPv4CidrBlock: "10.0.4.0/24",
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				// Two are simulated to exist...
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{
									AssociationId: ptr.To[string]("association-id-existing-1"),
									CidrBlock:     ptr.To[string]("10.0.1.0/24"),
								},
								{
									AssociationId: ptr.To[string]("association-id-existing-3"),
									CidrBlock:     ptr.To[string]("10.0.3.0/24"),
								},
							},
						},
					}}, nil)

				// ...the other two should be created
				m.AssociateVpcCidrBlockWithContext(context.TODO(), gomock.Eq(&ec2.AssociateVpcCidrBlockInput{
					CidrBlock: ptr.To[string]("10.0.2.0/24"),
					VpcId:     ptr.To[string]("vpc-id"),
				})).Return(&ec2.AssociateVpcCidrBlockOutput{
					CidrBlockAssociation: &ec2.VpcCidrBlockAssociation{
						AssociationId: ptr.To[string]("association-id-success-2"),
					},
				}, nil)
				m.AssociateVpcCidrBlockWithContext(context.TODO(), gomock.Eq(&ec2.AssociateVpcCidrBlockInput{
					CidrBlock: ptr.To[string]("10.0.4.0/24"),
					VpcId:     ptr.To[string]("vpc-id"),
				})).Return(&ec2.AssociateVpcCidrBlockOutput{
					CidrBlockAssociation: &ec2.VpcCidrBlockAssociation{
						AssociationId: ptr.To[string]("association-id-success-4"),
					},
				}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			cl := fake.NewClientBuilder().WithScheme(scheme).Build()

			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			mcpScope, err := setupNewManagedControlPlaneScope(cl)
			g.Expect(err).NotTo(HaveOccurred())

			if !tt.fillAWSManagedControlPlaneSecondaryCIDR {
				mcpScope.ControlPlane.Spec.SecondaryCidrBlock = nil
			}
			mcpScope.ControlPlane.Spec.NetworkSpec.VPC.SecondaryCidrBlocks = tt.networkSecondaryCIDRBlocks

			s := NewService(mcpScope)
			s.EC2Client = ec2Mock

			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT())
			}

			err = s.associateSecondaryCidrs()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestServiceDiassociateSecondaryCidr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                                    string
		fillAWSManagedControlPlaneSecondaryCIDR bool
		networkSecondaryCIDRBlocks              []infrav1.VpcCidrBlock
		expect                                  func(m *mocks.MockEC2APIMockRecorder)
		wantErr                                 bool
	}{
		{
			name:                                    "Should not disassociate secondary CIDR if no secondary cidr block info present in control plane",
			fillAWSManagedControlPlaneSecondaryCIDR: false,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				// No calls expected
				m.DescribeVpcsWithContext(context.TODO(), gomock.Any()).Times(0)
				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.Any()).Times(0)
			},
			wantErr: false,
		},
		{
			name:                                    "Should return error if unable to describe VPC",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:                                    "Should return error if no VPC found",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:                                    "Should diassociate secondary cidr block if already exist in VPC",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr")},
							},
						},
					}}, nil)
				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DisassociateVpcCidrBlockInput{})).Return(nil, nil)
			},
		},
		{
			name:                                    "Should return error if failed to diassociate secondary cidr block",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr")},
							},
						},
					}}, nil)
				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DisassociateVpcCidrBlockInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:                                    "Should successfully return from disassociating secondary CIDR blocks if none is currently associated",
			fillAWSManagedControlPlaneSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{},
						},
					}}, nil)

				// No calls expected
				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.Any()).Times(0)
			},
			wantErr: false,
		},
		{
			name:                                    "Should successfully disassociate existing secondary CIDR blocks",
			fillAWSManagedControlPlaneSecondaryCIDR: false,
			networkSecondaryCIDRBlocks: []infrav1.VpcCidrBlock{
				{
					IPv4CidrBlock: "10.0.1.0/24",
				},
				{
					IPv4CidrBlock: "10.0.2.0/24",
				},
				{
					IPv4CidrBlock: "10.0.3.0/24",
				},
				{
					IPv4CidrBlock: "10.0.4.0/24",
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				// Two are simulated to exist...
				m.DescribeVpcsWithContext(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{
									AssociationId: ptr.To[string]("association-id-existing-1"),
									CidrBlock:     ptr.To[string]("10.0.1.0/24"),
								},
								{
									AssociationId: ptr.To[string]("association-id-existing-3"),
									CidrBlock:     ptr.To[string]("10.0.3.0/24"),
								},
							},
						},
					}}, nil)

				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.Eq(&ec2.DisassociateVpcCidrBlockInput{
					AssociationId: ptr.To[string]("association-id-existing-1"), // 10.0.1.0/24 (see above)
				})).Return(&ec2.DisassociateVpcCidrBlockOutput{}, nil)
				m.DisassociateVpcCidrBlockWithContext(context.TODO(), gomock.Eq(&ec2.DisassociateVpcCidrBlockInput{
					AssociationId: ptr.To[string]("association-id-existing-3"), // 10.0.3.0/24 (see above)
				})).Return(&ec2.DisassociateVpcCidrBlockOutput{}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			g.Expect(err).NotTo(HaveOccurred())
			cl := fake.NewClientBuilder().WithScheme(scheme).Build()

			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			mcpScope, err := setupNewManagedControlPlaneScope(cl)
			g.Expect(err).NotTo(HaveOccurred())

			if !tt.fillAWSManagedControlPlaneSecondaryCIDR {
				mcpScope.ControlPlane.Spec.SecondaryCidrBlock = nil
			}
			mcpScope.ControlPlane.Spec.NetworkSpec.VPC.SecondaryCidrBlocks = tt.networkSecondaryCIDRBlocks

			s := NewService(mcpScope)
			s.EC2Client = ec2Mock

			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT())
			}

			err = s.disassociateSecondaryCidrs()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}
