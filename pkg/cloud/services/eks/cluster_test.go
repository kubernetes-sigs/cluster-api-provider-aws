/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/version"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks/mock_eksiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestMakeEksEncryptionConfigs(t *testing.T) {
	providerOne := "provider"
	resourceOne := "resourceOne"
	resourceTwo := "resourceTwo"
	testCases := []struct {
		name   string
		input  *ekscontrolplanev1.EncryptionConfig
		expect []*eks.EncryptionConfig
	}{
		{
			name:   "nil input",
			input:  nil,
			expect: []*eks.EncryptionConfig{},
		},
		{
			name: "nil input",
			input: &ekscontrolplanev1.EncryptionConfig{
				Provider:  &providerOne,
				Resources: []*string{&resourceOne, &resourceTwo},
			},
			expect: []*eks.EncryptionConfig{{
				Provider:  &eks.Provider{KeyArn: &providerOne},
				Resources: []*string{&resourceOne, &resourceTwo},
			}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(makeEksEncryptionConfigs(tc.input)).To(Equal(tc.expect))
		})
	}
}

func TestParseEKSVersion(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect version.Version
	}{
		{
			name:   "with patch",
			input:  "1.17.8",
			expect: *version.MustParseGeneric("1.17"),
		},
		{
			name:   "with v",
			input:  "v1.17.8",
			expect: *version.MustParseGeneric("1.17"),
		},
		{
			name:   "no patch",
			input:  "1.17",
			expect: *version.MustParseGeneric("1.17"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(*parseEKSVersion(tc.input)).To(Equal(tc.expect))
		})
	}
}
func TestVersionToEKS(t *testing.T) {
	testCases := []struct {
		name   string
		input  *version.Version
		expect string
	}{
		{
			name:   "with patch",
			input:  version.MustParseGeneric("1.17.8"),
			expect: "1.17",
		},
		{
			name:   "no patch",
			input:  version.MustParseGeneric("1.17"),
			expect: "1.17",
		},
		{
			name:   "with extra data",
			input:  version.MustParseGeneric("1.17-alpha"),
			expect: "1.17",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(versionToEKS(tc.input)).To(Equal(tc.expect))
		})
	}
}

func TestMakeVPCConfig(t *testing.T) {
	type input struct {
		subnets        infrav1.Subnets
		endpointAccess ekscontrolplanev1.EndpointAccess
		securityGroups map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	}

	idOne := "one"
	idTwo := "two"
	testCases := []struct {
		name   string
		input  input
		err    bool
		expect *eks.VpcConfigRequest
	}{
		{
			name: "no subnets",
			input: input{
				subnets:        nil,
				endpointAccess: ekscontrolplanev1.EndpointAccess{},
			},
			err:    true,
			expect: nil,
		},
		{
			name: "enough subnets",
			input: input{
				subnets: []*infrav1.SubnetSpec{
					{
						ID:               idOne,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2a",
						IsPublic:         true,
					},
					{
						ID:               idTwo,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2b",
						IsPublic:         false,
					},
				},
				endpointAccess: ekscontrolplanev1.EndpointAccess{},
			},
			expect: &eks.VpcConfigRequest{
				SubnetIds: []*string{&idOne, &idTwo},
			},
		},
		{
			name: "security groups",
			input: input{
				subnets: []*infrav1.SubnetSpec{
					{
						ID:               idOne,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2a",
						IsPublic:         true,
					},
					{
						ID:               idTwo,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2b",
						IsPublic:         false,
					},
				},
				endpointAccess: ekscontrolplanev1.EndpointAccess{},
				securityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupEKSNodeAdditional: {
						ID: idOne,
					},
				},
			},
			expect: &eks.VpcConfigRequest{
				SubnetIds:        []*string{&idOne, &idTwo},
				SecurityGroupIds: []*string{&idOne},
			},
		},
		{
			name: "non canonical public access CIDR",
			input: input{
				subnets: []*infrav1.SubnetSpec{
					{
						ID:               idOne,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2a",
						IsPublic:         true,
					},
					{
						ID:               idTwo,
						CidrBlock:        "10.0.10.1/24",
						AvailabilityZone: "us-west-2b",
						IsPublic:         false,
					},
				},
				endpointAccess: ekscontrolplanev1.EndpointAccess{
					PublicCIDRs: []*string{aws.String("10.0.0.1/24")},
				},
			},
			expect: &eks.VpcConfigRequest{
				SubnetIds:         []*string{&idOne, &idTwo},
				PublicAccessCidrs: []*string{aws.String("10.0.0.0/24")},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			config, err := makeVpcConfig(tc.input.subnets, tc.input.endpointAccess, tc.input.securityGroups)
			if tc.err {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(config).To(Equal(tc.expect))
			}
		})
	}

}

func TestPublicAccessCIDRsEqual(t *testing.T) {
	testCases := []struct {
		name   string
		a      []*string
		b      []*string
		expect bool
	}{
		{
			name:   "no CIDRs",
			a:      nil,
			b:      nil,
			expect: true,
		},
		{
			name:   "every address",
			a:      []*string{aws.String("0.0.0.0/0")},
			b:      nil,
			expect: true,
		},
		{
			name:   "every address",
			a:      []*string{aws.String("1.1.1.0/24")},
			b:      []*string{aws.String("1.1.1.0/24")},
			expect: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(publicAccessCIDRsEqual(tc.a, tc.b)).To(Equal(tc.expect))
		})
	}
}

func TestMakeEKSLogging(t *testing.T) {
	testCases := []struct {
		name   string
		input  *ekscontrolplanev1.ControlPlaneLoggingSpec
		expect *eks.Logging
	}{
		{
			name:   "no subnets",
			input:  nil,
			expect: nil,
		},
		{
			name: "some enabled, some disabled",
			input: &ekscontrolplanev1.ControlPlaneLoggingSpec{
				APIServer: true,
				Audit:     false,
			},
			expect: &eks.Logging{
				ClusterLogging: []*eks.LogSetup{
					{
						Enabled: aws.Bool(true),
						Types:   []*string{aws.String(eks.LogTypeApi)},
					},
					{
						Enabled: aws.Bool(false),
						Types: []*string{
							aws.String(eks.LogTypeAudit),
							aws.String(eks.LogTypeAuthenticator),
							aws.String(eks.LogTypeControllerManager),
							aws.String(eks.LogTypeScheduler),
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			logging := makeEksLogging(tc.input)
			g.Expect(logging).To(Equal(tc.expect))
		})
	}
}

func TestReconcileClusterVersion(t *testing.T) {
	clusterName := "default.cluster"
	tests := []struct {
		name        string
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "no upgrade necessary",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DescribeCluster(gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &eks.Cluster{
							Name:    aws.String("default.cluster"),
							Version: aws.String("1.16"),
						},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "needs upgrade",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DescribeCluster(gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &eks.Cluster{
							Name:    aws.String("default.cluster"),
							Version: aws.String("1.14"),
						},
					}, nil)
				m.WaitUntilClusterUpdating(
					gomock.AssignableToTypeOf(&eks.DescribeClusterInput{}), gomock.Any(),
				).Return(nil)
				m.
					UpdateClusterVersion(gomock.AssignableToTypeOf(&eks.UpdateClusterVersionInput{})).
					Return(&eks.UpdateClusterVersionOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "api error",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DescribeCluster(gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &eks.Cluster{
							Name:    aws.String("default.cluster"),
							Version: aws.String("1.14"),
						},
					}, nil)
				m.
					UpdateClusterVersion(gomock.AssignableToTypeOf(&eks.UpdateClusterVersionInput{})).
					Return(&eks.UpdateClusterVersionOutput{}, errors.New(""))
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						Version: aws.String("1.16"),
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			cluster, err := s.describeEKSCluster(clusterName)
			g.Expect(err).To(BeNil())

			err = s.reconcileClusterVersion(cluster)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}
