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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestService_releaseAddresses(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name    string
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name: "Should return error if failed to describe IP addresses",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
			wantErr: true,
		},
		{
			name: "Should ignore releasing elastic IP addresses if not found",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(nil, nil)
			},
		},
		{
			name: "Should return error if failed to disassociate IP address",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(&ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{
							AssociationId: aws.String("association-id-1"),
							PublicIp:      aws.String("public-ip"),
							AllocationId:  aws.String("allocation-id"),
						},
					},
				}, nil)
				m.DisassociateAddress(gomock.AssignableToTypeOf(&ec2.DisassociateAddressInput{})).Return(nil, awserrors.NewFailedDependency("dependency-failure"))
			},
			wantErr: true,
		},
		{
			name: "Should be able to release the IP address",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(&ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{
							AssociationId: aws.String("association-id-1"),
							PublicIp:      aws.String("public-ip"),
							AllocationId:  aws.String("allocation-id"),
						},
					},
				}, nil)
				m.DisassociateAddress(gomock.AssignableToTypeOf(&ec2.DisassociateAddressInput{})).Return(nil, nil)
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, nil)
			},
		},
		{
			name: "Should retry if unable to release the IP address because of Auth Failure",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(&ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{
							AssociationId: aws.String("association-id-1"),
							PublicIp:      aws.String("public-ip"),
							AllocationId:  aws.String("allocation-id"),
						},
					},
				}, nil)
				m.DisassociateAddress(gomock.AssignableToTypeOf(&ec2.DisassociateAddressInput{})).Return(nil, nil).Times(2)
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, awserr.New(awserrors.AuthFailure, awserrors.AuthFailure, errors.Errorf(awserrors.AuthFailure)))
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, nil)
			},
		},
		{
			name: "Should retry if unable to release the IP address because IP is already in use",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(&ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{
							AssociationId: aws.String("association-id-1"),
							PublicIp:      aws.String("public-ip"),
							AllocationId:  aws.String("allocation-id"),
						},
					},
				}, nil)
				m.DisassociateAddress(gomock.AssignableToTypeOf(&ec2.DisassociateAddressInput{})).Return(nil, nil).Times(2)
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, awserr.New(awserrors.InUseIPAddress, awserrors.InUseIPAddress, errors.Errorf(awserrors.InUseIPAddress)))
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, nil)
			},
		},
		{
			name: "Should not retry if unable to release the IP address due to dependency failure",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAddresses(gomock.AssignableToTypeOf(&ec2.DescribeAddressesInput{})).Return(&ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{
							AssociationId: aws.String("association-id-1"),
							PublicIp:      aws.String("public-ip"),
							AllocationId:  aws.String("allocation-id"),
						},
					},
				}, nil)
				m.DisassociateAddress(gomock.AssignableToTypeOf(&ec2.DisassociateAddressInput{})).Return(nil, nil).Times(2)
				m.ReleaseAddress(gomock.AssignableToTypeOf(&ec2.ReleaseAddressInput{})).Return(nil, awserr.New("dependency-failure", "dependency-failure", errors.Errorf("dependency-failure")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme := runtime.NewScheme()
			err := infrav1.AddToScheme(scheme)
			g.Expect(err).NotTo(HaveOccurred())
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client:     client,
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			})
			g.Expect(err).NotTo(HaveOccurred())

			s := NewService(cs)
			s.EC2Client = ec2Mock

			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT())
			}

			if err := s.releaseAddresses(); (err != nil) != tt.wantErr {
				t.Errorf("releaseAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
