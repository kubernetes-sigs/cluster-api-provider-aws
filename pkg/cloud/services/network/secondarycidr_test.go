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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
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
				SecondaryCidrBlock: pointer.String("secondary-cidr"),
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
		name              string
		haveSecondaryCIDR bool
		expect            func(m *mocks.MockEC2APIMockRecorder)
		wantErr           bool
	}{
		{
			name: "Should not associate secondary CIDR if no secondary cidr block info present in control plane",
		},
		{
			name:              "Should return error if unable to describe VPC",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:              "Should not associate secondary cidr block if already exist in VPC",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
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
			name:              "Should return error if no VPC found",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:              "Should return error if failed during associating secondary cidr block",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr-new")},
							},
						},
					}}, nil)
				m.AssociateVpcCidrBlock(gomock.AssignableToTypeOf(&ec2.AssociateVpcCidrBlockInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
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

			if !tt.haveSecondaryCIDR {
				mcpScope.ControlPlane.Spec.SecondaryCidrBlock = nil
			}

			s := NewService(mcpScope)
			s.EC2Client = ec2Mock

			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT())
			}

			err = s.associateSecondaryCidr()
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
		name              string
		haveSecondaryCIDR bool
		expect            func(m *mocks.MockEC2APIMockRecorder)
		wantErr           bool
	}{
		{
			name: "Should not disassociate secondary CIDR if no secondary cidr block info present in control plane",
		},
		{
			name:              "Should return error if unable to describe VPC",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name:              "Should return error if no VPC found",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:              "Should diassociate secondary cidr block if already exist in VPC",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr")},
							},
						},
					}}, nil)
				m.DisassociateVpcCidrBlock(gomock.AssignableToTypeOf(&ec2.DisassociateVpcCidrBlockInput{})).Return(nil, nil)
			},
		},
		{
			name:              "Should return error if failed to diassociate secondary cidr block",
			haveSecondaryCIDR: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
								{CidrBlock: aws.String("secondary-cidr")},
							},
						},
					}}, nil)
				m.DisassociateVpcCidrBlock(gomock.AssignableToTypeOf(&ec2.DisassociateVpcCidrBlockInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
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

			if !tt.haveSecondaryCIDR {
				mcpScope.ControlPlane.Spec.SecondaryCidrBlock = nil
			}

			s := NewService(mcpScope)
			s.EC2Client = ec2Mock

			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT())
			}

			err = s.disassociateSecondaryCidr()
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}
