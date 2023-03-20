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
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestMakeEKSEncryptionConfigs(t *testing.T) {
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
				subnets: []infrav1.SubnetSpec{
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
			name: "ipv6 subnets",
			input: input{
				subnets: []infrav1.SubnetSpec{
					{
						ID:               idOne,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2a",
						IsPublic:         true,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:85a3:1::/64",
					},
					{
						ID:               idTwo,
						CidrBlock:        "10.0.10.0/24",
						AvailabilityZone: "us-west-2b",
						IsPublic:         false,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:85a3:2::/64",
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
				subnets: []infrav1.SubnetSpec{
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
				subnets: []infrav1.SubnetSpec{
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

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
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

func TestCreateCluster(t *testing.T) {
	clusterName := "cluster.default"
	version := aws.String("1.24")
	tests := []struct {
		name        string
		expectEKS   func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
		role        *string
		tags        map[string]*string
		subnets     []infrav1.SubnetSpec
	}{
		{
			name:        "cluster create with 2 subnets",
			expectEKS:   func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
			role:        aws.String("arn:role"),
			tags: map[string]*string{
				"kubernetes.io/cluster/" + clusterName: aws.String("owned"),
			},
			subnets: []infrav1.SubnetSpec{
				{ID: "1", AvailabilityZone: "us-west-2a"}, {ID: "2", AvailabilityZone: "us-west-2b"},
			},
		},
		{
			name:        "cluster create without subnets",
			expectEKS:   func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: true,
			role:        aws.String("arn:role"),
			subnets:     []infrav1.SubnetSpec{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			iamMock := mock_iamauth.NewMockIAMAPI(mockControl)
			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, _ := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      "capi-name",
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
						Version:        version,
						RoleName:       tc.role,
						NetworkSpec:    infrav1.NetworkSpec{Subnets: tc.subnets},
					},
				},
			})
			subnetIds := make([]*string, 0)
			for i := range tc.subnets {
				subnet := tc.subnets[i]
				subnetIds = append(subnetIds, &subnet.ID)
			}

			if !tc.expectError {
				roleOutput := iam.GetRoleOutput{Role: &iam.Role{Arn: tc.role}}
				iamMock.EXPECT().GetRole(gomock.Any()).Return(&roleOutput, nil)
				eksMock.EXPECT().CreateCluster(&eks.CreateClusterInput{
					Name:             aws.String(clusterName),
					EncryptionConfig: []*eks.EncryptionConfig{},
					ResourcesVpcConfig: &eks.VpcConfigRequest{
						SubnetIds: subnetIds,
					},
					RoleArn: tc.role,
					Tags:    tc.tags,
					Version: version,
				}).Return(&eks.CreateClusterOutput{}, nil)
			}
			s := NewService(scope)
			s.IAMClient = iamMock
			s.EKSClient = eksMock

			_, err := s.createCluster(clusterName)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestReconcileEKSEncryptionConfig(t *testing.T) {
	clusterName := "default.cluster"
	tests := []struct {
		name                string
		oldEncryptionConfig *ekscontrolplanev1.EncryptionConfig
		newEncryptionConfig *ekscontrolplanev1.EncryptionConfig
		expect              func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError         bool
	}{
		{
			name:                "no upgrade necessary - encryption disabled",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{},
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{},
			expect:              func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError:         false,
		},
		{
			name: "no upgrade necessary - encryption config unchanged",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
		},
		{
			name:                "needs upgrade",
			oldEncryptionConfig: nil,
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.WaitUntilClusterUpdating(
					gomock.AssignableToTypeOf(&eks.DescribeClusterInput{}), gomock.Any(),
				).Return(nil)
				m.AssociateEncryptionConfig(gomock.AssignableToTypeOf(&eks.AssociateEncryptionConfigInput{})).Return(&eks.AssociateEncryptionConfigOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "upgrade not allowed if encryption config updated as nil",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			newEncryptionConfig: nil,
			expect:              func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError:         true,
		},
		{
			name: "upgrade not allowed if encryption config exists",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  pointer.String("new-provider"),
				Resources: []*string{pointer.String("foo"), pointer.String("bar")},
			},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						Version:          aws.String("1.16"),
						EncryptionConfig: tc.newEncryptionConfig,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.reconcileEKSEncryptionConfig(makeEksEncryptionConfigs(tc.oldEncryptionConfig))
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestCreateIPv6Cluster(t *testing.T) {
	g := NewWithT(t)

	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
	iamMock := mock_iamauth.NewMockIAMAPI(mockControl)

	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	encryptionConfig := &ekscontrolplanev1.EncryptionConfig{
		Provider:  pointer.String("new-provider"),
		Resources: []*string{pointer.String("foo"), pointer.String("bar")},
	}
	vpcSpec := infrav1.VPCSpec{
		IPv6: &infrav1.IPv6{
			CidrBlock: "2001:db8:85a3::/56",
		},
	}
	scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "cluster-name",
			},
		},
		ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
			Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				RoleName: pointer.String("arn-role"),
				Version:  aws.String("1.22"),
				NetworkSpec: infrav1.NetworkSpec{
					Subnets: []infrav1.SubnetSpec{
						{
							ID:               "sub-1",
							CidrBlock:        "10.0.10.0/24",
							AvailabilityZone: "us-west-2a",
							IsPublic:         true,
							IsIPv6:           true,
							IPv6CidrBlock:    "2001:db8:85a3:1::/64",
						},
						{
							ID:               "sub-2",
							CidrBlock:        "10.0.10.0/24",
							AvailabilityZone: "us-west-2b",
							IsPublic:         false,
							IsIPv6:           true,
							IPv6CidrBlock:    "2001:db8:85a3:2::/64",
						},
					},
					VPC: vpcSpec,
				},
				EncryptionConfig: encryptionConfig,
			},
		},
	})
	g.Expect(err).To(BeNil())

	eksMock.EXPECT().CreateCluster(&eks.CreateClusterInput{
		Name:    aws.String("cluster-name"),
		Version: aws.String("1.22"),
		EncryptionConfig: []*eks.EncryptionConfig{
			{
				Provider: &eks.Provider{
					KeyArn: encryptionConfig.Provider,
				},
				Resources: encryptionConfig.Resources,
			},
		},
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds: []*string{pointer.String("sub-1"), pointer.String("sub-2")},
		},
		KubernetesNetworkConfig: &eks.KubernetesNetworkConfigRequest{
			IpFamily: pointer.String("ipv6"),
		},
		Tags: map[string]*string{
			"kubernetes.io/cluster/cluster-name": pointer.String("owned"),
		},
	}).Return(&eks.CreateClusterOutput{}, nil)
	iamMock.EXPECT().GetRole(&iam.GetRoleInput{
		RoleName: aws.String("arn-role"),
	}).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			RoleName: pointer.String("arn-role"),
		},
	}, nil)

	s := NewService(scope)
	s.EKSClient = eksMock
	s.IAMClient = iamMock

	_, err = s.createCluster("cluster-name")
	g.Expect(err).To(BeNil())
}
