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
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestMakeEKSEncryptionConfigs(t *testing.T) {
	providerOne := "provider"
	resourceOne := "resourceOne"
	resourceTwo := "resourceTwo"
	testCases := []struct {
		name   string
		input  *ekscontrolplanev1.EncryptionConfig
		expect []ekstypes.EncryptionConfig
	}{
		{
			name:   "nil input",
			input:  nil,
			expect: []ekstypes.EncryptionConfig{},
		},
		{
			name: "nil input",
			input: &ekscontrolplanev1.EncryptionConfig{
				Provider:  &providerOne,
				Resources: []*string{&resourceOne, &resourceTwo},
			},
			expect: []ekstypes.EncryptionConfig{{
				Provider:  &ekstypes.Provider{KeyArn: &providerOne},
				Resources: []string{resourceOne, resourceTwo},
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

func TestCompareEncryptionConfigIgnoresResources(t *testing.T) {
	g := NewWithT(t)
	// EKS encrypts all Kubernetes API data by default, so Resources no longer
	// determines which resources are encrypted and is not useful for comparison.
	updated := []ekstypes.EncryptionConfig{{
		Provider:  &ekstypes.Provider{KeyArn: ptr.To("key")},
		Resources: []string{"secrets"},
	}}
	existing := []ekstypes.EncryptionConfig{{
		Provider: &ekstypes.Provider{KeyArn: ptr.To("key")},
	}}

	g.Expect(compareEncryptionConfig(updated, existing)).To(BeTrue())
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
			v, err := parseEKSVersion(tc.input)
			g.Expect(err).To(BeNil())
			g.Expect(*v).To(Equal(tc.expect))
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
		expect *ekstypes.VpcConfigRequest
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
			expect: &ekstypes.VpcConfigRequest{
				SubnetIds: []string{idOne, idTwo},
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
			expect: &ekstypes.VpcConfigRequest{
				SubnetIds: []string{idOne, idTwo},
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
			expect: &ekstypes.VpcConfigRequest{
				SubnetIds:        []string{idOne, idTwo},
				SecurityGroupIds: []string{idOne},
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
			expect: &ekstypes.VpcConfigRequest{
				SubnetIds:         []string{idOne, idTwo},
				PublicAccessCidrs: []string{"10.0.0.0/24"},
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
		a      []string
		b      []string
		expect bool
	}{
		{
			name:   "no CIDRs",
			a:      nil,
			b:      nil,
			expect: true,
		},
		{
			name:   "every ipv4 address",
			a:      []string{"0.0.0.0/0"},
			b:      nil,
			expect: true,
		},
		{
			name:   "every ipv4 and ipv6 address",
			a:      []string{"0.0.0.0/0", "::/0"},
			b:      nil,
			expect: true,
		},
		{
			name:   "every address",
			a:      []string{"1.1.1.0/24"},
			b:      []string{"1.1.1.0/24"},
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
		expect *ekstypes.Logging
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
			expect: &ekstypes.Logging{
				ClusterLogging: []ekstypes.LogSetup{
					{
						Enabled: aws.Bool(true),
						Types:   []ekstypes.LogType{ekstypes.LogTypeApi},
					},
					{
						Enabled: aws.Bool(false),
						Types: []ekstypes.LogType{
							ekstypes.LogTypeAudit,
							ekstypes.LogTypeAuthenticator,
							ekstypes.LogTypeControllerManager,
							ekstypes.LogTypeScheduler,
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
		name                  string
		clusterVersion        string
		desiredVersion        string
		expectedUpdateVersion string
		updateError           error
		expectError           bool
	}{
		{
			name:           "no version update necessary",
			clusterVersion: "1.16",
			desiredVersion: "1.16",
			expectError:    false,
		},
		{
			name:                  "multi-minor upgrade targets the next minor version",
			clusterVersion:        "1.14",
			desiredVersion:        "1.16",
			expectedUpdateVersion: "1.15",
			expectError:           false,
		},
		{
			name:                  "rollback targets the desired previous minor version",
			clusterVersion:        "1.16",
			desiredVersion:        "1.15",
			expectedUpdateVersion: "1.15",
			expectError:           false,
		},
		{
			name:                  "api error",
			clusterVersion:        "1.14",
			desiredVersion:        "1.16",
			expectedUpdateVersion: "1.15",
			updateError:           errors.New(""),
			expectError:           true,
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
						EKSClusterName: clusterName,
						Version:        aws.String(tc.desiredVersion),
					},
				},
			})
			g.Expect(err).To(BeNil())

			eksMock.EXPECT().
				DescribeCluster(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
				Return(&eks.DescribeClusterOutput{
					Cluster: &ekstypes.Cluster{
						Name:    aws.String(clusterName),
						Version: aws.String(tc.clusterVersion),
					},
				}, nil)

			if tc.expectedUpdateVersion != "" {
				eksMock.EXPECT().
					UpdateClusterVersion(gomock.Eq(context.TODO()), gomock.Eq(&eks.UpdateClusterVersionInput{
						Name:    aws.String(clusterName),
						Version: aws.String(tc.expectedUpdateVersion),
					})).
					Return(&eks.UpdateClusterVersionOutput{}, tc.updateError)
				if tc.updateError == nil {
					eksMock.EXPECT().
						WaitUntilClusterUpdating(
							gomock.Eq(context.TODO()),
							gomock.Eq(&eks.DescribeClusterInput{Name: aws.String(clusterName)}),
							gomock.Any(),
						).
						Return(nil)
				}
			}

			s := NewService(scope)
			s.EKSClient = eksMock

			cluster, err := s.describeEKSCluster(context.TODO(), clusterName)
			g.Expect(err).To(BeNil())

			err = s.reconcileClusterVersion(context.TODO(), cluster)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestReconcileAccessConfig(t *testing.T) {
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
					DescribeCluster(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &ekstypes.Cluster{
							Name: aws.String("default.cluster"),
							AccessConfig: &ekstypes.AccessConfigResponse{
								AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
							},
						},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "needs upgrade",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DescribeCluster(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &ekstypes.Cluster{
							Name: aws.String("default.cluster"),
							AccessConfig: &ekstypes.AccessConfigResponse{
								AuthenticationMode: ekstypes.AuthenticationModeConfigMap,
							},
						},
					}, nil)
				m.WaitUntilClusterUpdating(
					gomock.Eq(context.TODO()),
					gomock.AssignableToTypeOf(&eks.DescribeClusterInput{}),
					gomock.Any(),
				).Return(nil)
				m.
					UpdateClusterConfig(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.UpdateClusterConfigInput{})).
					Return(&eks.UpdateClusterConfigOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "api error",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DescribeCluster(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.DescribeClusterInput{})).
					Return(&eks.DescribeClusterOutput{
						Cluster: &ekstypes.Cluster{
							Name: aws.String("default.cluster"),
							AccessConfig: &ekstypes.AccessConfigResponse{
								AuthenticationMode: ekstypes.AuthenticationModeApi,
							},
						},
					}, nil)
				m.
					UpdateClusterConfig(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.UpdateClusterConfigInput{})).
					Return(&eks.UpdateClusterConfigOutput{}, errors.New("Unsupported authentication mode update"))
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
						EKSClusterName: clusterName,
						AccessConfig: &ekscontrolplanev1.AccessConfig{
							AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
						},
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			cluster, err := s.describeEKSCluster(context.TODO(), clusterName)
			g.Expect(err).To(BeNil())

			err = s.reconcileAccessConfig(context.TODO(), cluster.AccessConfig)
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
		tags        map[string]string
		subnets     []infrav1.SubnetSpec
	}{
		{
			name:        "cluster create with 2 subnets",
			expectEKS:   func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
			role:        aws.String("arn:role"),
			tags: map[string]string{
				"kubernetes.io/cluster/" + clusterName: "owned",
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
						EKSClusterName:             clusterName,
						Version:                    version,
						RoleName:                   tc.role,
						NetworkSpec:                infrav1.NetworkSpec{Subnets: tc.subnets},
						BootstrapSelfManagedAddons: false,
						UpgradePolicy:              ekscontrolplanev1.UpgradePolicyStandard,
					},
				},
			})
			subnetIDs := make([]string, 0)
			for i := range tc.subnets {
				subnet := tc.subnets[i]
				subnetIDs = append(subnetIDs, subnet.ID)
			}

			if !tc.expectError {
				roleOutput := iam.GetRoleOutput{Role: &iamtypes.Role{Arn: tc.role}}
				iamMock.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(&roleOutput, nil)
				eksMock.EXPECT().CreateCluster(context.TODO(), &eks.CreateClusterInput{
					Name:             aws.String(clusterName),
					EncryptionConfig: []ekstypes.EncryptionConfig{},
					ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
						SubnetIds: subnetIDs,
					},
					RoleArn:                    tc.role,
					Tags:                       tc.tags,
					Version:                    version,
					BootstrapSelfManagedAddons: aws.Bool(false),
					UpgradePolicy: &ekstypes.UpgradePolicyRequest{
						SupportType: ekstypes.SupportTypeStandard,
					},
				}).Return(&eks.CreateClusterOutput{}, nil)
			}
			s := NewService(scope)
			s.IAMClient = iamMock
			s.EKSClient = eksMock

			_, err := s.createCluster(context.TODO(), clusterName)
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
				Provider:  ptr.To[string]("provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
			},
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  ptr.To[string]("provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
			},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
		},
		{
			name:                "needs upgrade",
			oldEncryptionConfig: nil,
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  ptr.To[string]("provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.WaitUntilClusterUpdating(
					gomock.Eq(context.TODO()),
					gomock.AssignableToTypeOf(&eks.DescribeClusterInput{}),
					gomock.Any(),
				).Return(nil)
				m.AssociateEncryptionConfig(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.AssociateEncryptionConfigInput{})).Return(&eks.AssociateEncryptionConfigOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "upgrade not allowed if encryption config updated as nil",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  ptr.To[string]("provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
			},
			newEncryptionConfig: nil,
			expect:              func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError:         true,
		},
		{
			name: "upgrade not allowed if encryption config exists",
			oldEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  ptr.To[string]("provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
			},
			newEncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
				Provider:  ptr.To[string]("new-provider"),
				Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
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

			err = s.reconcileEKSEncryptionConfig(context.TODO(), makeEksEncryptionConfigs(tc.oldEncryptionConfig))
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestReconcileUpgradePolicy(t *testing.T) {
	clusterName := "default.cluster"
	tests := []struct {
		name             string
		oldUpgradePolicy *ekstypes.UpgradePolicyResponse
		newUpgradePolicy ekscontrolplanev1.UpgradePolicy
		expect           *ekstypes.UpgradePolicyRequest
		expectError      bool
	}{
		{
			name: "no update necessary - upgrade policy omitted",
			oldUpgradePolicy: &ekstypes.UpgradePolicyResponse{
				SupportType: ekstypes.SupportTypeStandard,
			},
			expect:      nil,
			expectError: false,
		},
		{
			name:             "no update necessary - cannot get cluster upgrade policy",
			newUpgradePolicy: ekscontrolplanev1.UpgradePolicyStandard,
			expect:           nil,
			expectError:      false,
		},
		{
			name: "no update necessary - upgrade policy unchanged",
			oldUpgradePolicy: &ekstypes.UpgradePolicyResponse{
				SupportType: ekstypes.SupportTypeStandard,
			},
			newUpgradePolicy: ekscontrolplanev1.UpgradePolicyStandard,
			expect:           nil,
			expectError:      false,
		},
		{
			name: "needs update",
			oldUpgradePolicy: &ekstypes.UpgradePolicyResponse{
				SupportType: ekstypes.SupportTypeStandard,
			},
			newUpgradePolicy: ekscontrolplanev1.UpgradePolicyExtended,
			expect: &ekstypes.UpgradePolicyRequest{
				SupportType: ekstypes.SupportTypeExtended,
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

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
						Version:       aws.String("1.16"),
						UpgradePolicy: tc.newUpgradePolicy,
					},
				},
			})
			g.Expect(err).To(BeNil())

			s := NewService(scope)

			upgradePolicyRequest := s.reconcileUpgradePolicy(tc.oldUpgradePolicy)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(upgradePolicyRequest).To(Equal(tc.expect))
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
		Provider:  ptr.To[string]("new-provider"),
		Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
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
				RoleName: ptr.To[string]("arn-role"),
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
				EncryptionConfig:           encryptionConfig,
				BootstrapSelfManagedAddons: false,
			},
		},
	})
	g.Expect(err).To(BeNil())

	eksMock.EXPECT().CreateCluster(context.TODO(), &eks.CreateClusterInput{
		Name:    aws.String("cluster-name"),
		Version: aws.String("1.22"),
		EncryptionConfig: []ekstypes.EncryptionConfig{
			{
				Provider: &ekstypes.Provider{
					KeyArn: encryptionConfig.Provider,
				},
				Resources: aws.ToStringSlice(encryptionConfig.Resources),
			},
		},
		ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
			SubnetIds: []string{"sub-1", "sub-2"},
		},
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			IpFamily: ekstypes.IpFamilyIpv6,
		},
		Tags: map[string]string{
			"kubernetes.io/cluster/cluster-name": "owned",
		},
		BootstrapSelfManagedAddons: aws.Bool(false),
	}).Return(&eks.CreateClusterOutput{}, nil)
	iamMock.EXPECT().GetRole(gomock.Any(), &iam.GetRoleInput{
		RoleName: aws.String("arn-role"),
	}).Return(&iam.GetRoleOutput{
		Role: &iamtypes.Role{
			RoleName: aws.String("arn-role"),
		},
	}, nil)

	s := NewService(scope)
	s.EKSClient = eksMock
	s.IAMClient = iamMock

	_, err = s.createCluster(context.TODO(), "cluster-name")
	g.Expect(err).To(BeNil())
}

func TestCreateClusterWithBootstrapClusterCreatorAdminPermissions(t *testing.T) {
	g := NewWithT(t)

	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
	iamMock := mock_iamauth.NewMockIAMAPI(mockControl)

	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	clusterName := "test-cluster"
	scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
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
				Version:        aws.String("1.24"),
				RoleName:       aws.String("arn:role"),
				NetworkSpec: infrav1.NetworkSpec{
					Subnets: []infrav1.SubnetSpec{
						{ID: "1", AvailabilityZone: "us-west-2a"},
						{ID: "2", AvailabilityZone: "us-west-2b"},
					},
				},
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					BootstrapClusterCreatorAdminPermissions: ptr.To(false),
				},
			},
		},
	})
	g.Expect(err).To(BeNil())

	eksMock.EXPECT().CreateCluster(context.TODO(), &eks.CreateClusterInput{
		Name:    aws.String(clusterName),
		Version: aws.String("1.24"),
		ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
			SubnetIds: []string{"1", "2"},
		},
		RoleArn: aws.String("arn:role"),
		Tags: map[string]string{
			"kubernetes.io/cluster/test-cluster": "owned",
		},
		AccessConfig: &ekstypes.CreateAccessConfigRequest{
			BootstrapClusterCreatorAdminPermissions: ptr.To(false),
		},
		EncryptionConfig:           []ekstypes.EncryptionConfig{},
		BootstrapSelfManagedAddons: aws.Bool(false),
	}).Return(&eks.CreateClusterOutput{}, nil)

	iamMock.EXPECT().GetRole(gomock.Any(), gomock.Any()).Return(&iam.GetRoleOutput{
		Role: &iamtypes.Role{Arn: aws.String("arn:role")},
	}, nil)

	s := NewService(scope)
	s.EKSClient = eksMock
	s.IAMClient = iamMock

	_, err = s.createCluster(context.TODO(), clusterName)
	g.Expect(err).To(BeNil())
}

func TestReconcileScalingConfig(t *testing.T) {
	testCases := []struct {
		name           string
		specConfig     *ekscontrolplanev1.ControlPlaneScalingConfig
		clusterConfig  *ekstypes.ControlPlaneScalingConfig
		expectUpdate   bool
		expectedConfig *ekstypes.ControlPlaneScalingConfig
	}{
		{
			name:           "spec has no scaling config - no update",
			specConfig:     nil,
			clusterConfig:  nil,
			expectUpdate:   false,
			expectedConfig: nil,
		},
		{
			name: "cluster has no scaling config but spec does - update needed",
			specConfig: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier2XL,
			},
			clusterConfig: nil,
			expectUpdate:  true,
			expectedConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
		},
		{
			name: "tier differs - update needed",
			specConfig: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier4XL,
			},
			clusterConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
			expectUpdate: true,
			expectedConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-4xl"),
			},
		},
		{
			name: "tier is same - no update",
			specConfig: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier2XL,
			},
			clusterConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
			expectUpdate:   false,
			expectedConfig: nil,
		},
		{
			name: "changing from standard to xl - update needed",
			specConfig: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTierXL,
			},
			clusterConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTierStandard,
			},
			expectUpdate: true,
			expectedConfig: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-xl"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      "capi-name",
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName:            "test-cluster",
						ControlPlaneScalingConfig: tc.specConfig,
					},
				},
			})
			g.Expect(err).To(BeNil())

			s := NewService(scope)
			result := s.reconcileScalingConfig(tc.clusterConfig)

			if tc.expectUpdate {
				g.Expect(result).ToNot(BeNil(), "expected update config to be returned")
				g.Expect(result.Tier).To(Equal(tc.expectedConfig.Tier))
			} else {
				g.Expect(result).To(BeNil(), "expected no update")
			}
		})
	}
}

func newComponentConfigTestService(t *testing.T, spec ekscontrolplanev1.AWSManagedControlPlaneSpec) *Service {
	t.Helper()

	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	if spec.EKSClusterName == "" {
		spec.EKSClusterName = "test-cluster"
	}

	managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "capi-name",
			},
		},
		ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
			Spec: spec,
		},
	})
	if err != nil {
		t.Fatalf("failed to create managed control plane scope: %v", err)
	}

	return NewService(managedScope)
}

func schedulerResponse(strategyType ekstypes.ScoringStrategyType, weights ...ekstypes.ResourceWeight) *ekstypes.KubeSchedulerConfigResponse {
	return &ekstypes.KubeSchedulerConfigResponse{
		NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
			ScoringStrategy: &ekstypes.ScoringStrategy{
				Type:      strategyType,
				Resources: weights,
			},
		},
	}
}

func TestReconcileKubeSchedulerConfig(t *testing.T) {
	defaultWeights := []ekstypes.ResourceWeight{
		{Name: aws.String("cpu"), Weight: aws.Int32(1)},
		{Name: aws.String("memory"), Weight: aws.Int32(1)},
	}
	mostAllocated := func(weights ...ekscontrolplanev1.ResourceWeight) *ekscontrolplanev1.KubeSchedulerConfig {
		return &ekscontrolplanev1.KubeSchedulerConfig{
			NodeResourcesFit: &ekscontrolplanev1.NodeResourcesFitConfig{
				ScoringStrategy: &ekscontrolplanev1.ScoringStrategy{
					Type:      ekscontrolplanev1.ScoringStrategyTypeMostAllocated,
					Resources: weights,
				},
			},
		}
	}

	testCases := []struct {
		name          string
		specConfig    *ekscontrolplanev1.KubeSchedulerConfig
		clusterConfig *ekstypes.KubeSchedulerConfigResponse
		expectUpdate  bool
	}{
		{
			name:          "spec has no scheduler config - no update",
			specConfig:    nil,
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeLeastAllocated, defaultWeights...),
			expectUpdate:  false,
		},
		{
			name:          "cluster reports no scheduler config - update needed",
			specConfig:    mostAllocated(),
			clusterConfig: nil,
			expectUpdate:  true,
		},
		{
			name:          "strategy type differs - update needed",
			specConfig:    mostAllocated(),
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeLeastAllocated, defaultWeights...),
			expectUpdate:  true,
		},
		{
			name:          "strategy type matches and spec has no resources - no update",
			specConfig:    mostAllocated(),
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated, defaultWeights...),
			expectUpdate:  false,
		},
		{
			name:          "resource weights differ - update needed",
			specConfig:    mostAllocated(ekscontrolplanev1.ResourceWeight{Name: "cpu", Weight: 2}, ekscontrolplanev1.ResourceWeight{Name: "memory", Weight: 2}),
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated, defaultWeights...),
			expectUpdate:  true,
		},
		{
			name:       "resource weights match in a different order - no update",
			specConfig: mostAllocated(ekscontrolplanev1.ResourceWeight{Name: "memory", Weight: 2}, ekscontrolplanev1.ResourceWeight{Name: "cpu", Weight: 2}),
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated,
				ekstypes.ResourceWeight{Name: aws.String("cpu"), Weight: aws.Int32(2)},
				ekstypes.ResourceWeight{Name: aws.String("memory"), Weight: aws.Int32(2)},
			),
			expectUpdate: false,
		},
		{
			name:          "spec lists fewer resources than the cluster - update needed",
			specConfig:    mostAllocated(ekscontrolplanev1.ResourceWeight{Name: "cpu", Weight: 1}),
			clusterConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated, defaultWeights...),
			expectUpdate:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			s := newComponentConfigTestService(t, ekscontrolplanev1.AWSManagedControlPlaneSpec{KubeSchedulerConfig: tc.specConfig})

			result := s.reconcileKubeSchedulerConfig(tc.clusterConfig)
			if tc.expectUpdate {
				g.Expect(result).ToNot(BeNil(), "expected update config to be returned")
				g.Expect(result.NodeResourcesFit.ScoringStrategy.Type).To(Equal(ekstypes.ScoringStrategyType(tc.specConfig.NodeResourcesFit.ScoringStrategy.Type)))
				g.Expect(result.NodeResourcesFit.ScoringStrategy.Resources).To(HaveLen(len(tc.specConfig.NodeResourcesFit.ScoringStrategy.Resources)))
			} else {
				g.Expect(result).To(BeNil(), "expected no update")
			}
		})
	}
}

func TestReconcileKubeAPIServerConfig(t *testing.T) {
	testCases := []struct {
		name          string
		specConfig    *ekscontrolplanev1.KubeAPIServerConfig
		clusterConfig *ekstypes.KubeApiServerConfigResponse
		expectUpdate  bool
	}{
		{
			name:          "spec has no api server config - no update",
			specConfig:    nil,
			clusterConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("60m")},
			expectUpdate:  false,
		},
		{
			name:          "cluster reports no api server config - update needed",
			specConfig:    &ekscontrolplanev1.KubeAPIServerConfig{EventTTL: "15m"},
			clusterConfig: nil,
			expectUpdate:  true,
		},
		{
			name:          "event ttl differs - update needed",
			specConfig:    &ekscontrolplanev1.KubeAPIServerConfig{EventTTL: "15m"},
			clusterConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("60m")},
			expectUpdate:  true,
		},
		{
			name:          "event ttl is equal in a different unit - no update",
			specConfig:    &ekscontrolplanev1.KubeAPIServerConfig{EventTTL: "60m"},
			clusterConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("1h")},
			expectUpdate:  false,
		},
		{
			name:       "service node port range matches and event ttl is not set - no update",
			specConfig: &ekscontrolplanev1.KubeAPIServerConfig{ServiceNodePortRange: &ekscontrolplanev1.ServiceNodePortRange{MinPort: 30000, MaxPort: 32767}},
			clusterConfig: &ekstypes.KubeApiServerConfigResponse{
				EventTtl:             aws.String("60m"),
				ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 30000, MaxPort: 32767},
			},
			expectUpdate: false,
		},
		{
			name:       "service node port range differs - update needed",
			specConfig: &ekscontrolplanev1.KubeAPIServerConfig{ServiceNodePortRange: &ekscontrolplanev1.ServiceNodePortRange{MinPort: 20000, MaxPort: 32767}},
			clusterConfig: &ekstypes.KubeApiServerConfigResponse{
				ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 30000, MaxPort: 32767},
			},
			expectUpdate: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			s := newComponentConfigTestService(t, ekscontrolplanev1.AWSManagedControlPlaneSpec{KubeAPIServerConfig: tc.specConfig})

			result := s.reconcileKubeAPIServerConfig(tc.clusterConfig)
			if tc.expectUpdate {
				g.Expect(result).ToNot(BeNil(), "expected update config to be returned")
				if tc.specConfig.EventTTL != "" {
					g.Expect(aws.ToString(result.EventTtl)).To(Equal(tc.specConfig.EventTTL))
				}
			} else {
				g.Expect(result).To(BeNil(), "expected no update")
			}
		})
	}
}

func TestReconcileKubeControllerManagerConfig(t *testing.T) {
	testCases := []struct {
		name          string
		specConfig    *ekscontrolplanev1.KubeControllerManagerConfig
		clusterConfig *ekstypes.KubeControllerManagerConfigResponse
		expectUpdate  bool
	}{
		{
			name:         "spec has no controller manager config - no update",
			specConfig:   nil,
			expectUpdate: false,
		},
		{
			name: "sync period differs - update needed",
			specConfig: &ekscontrolplanev1.KubeControllerManagerConfig{
				HorizontalPodAutoscalerControllerConfig: &ekscontrolplanev1.HorizontalPodAutoscalerControllerConfig{HorizontalPodAutoscalerSyncPeriod: "10s"},
			},
			clusterConfig: &ekstypes.KubeControllerManagerConfigResponse{
				HorizontalPodAutoscalerControllerConfig: &ekstypes.HorizontalPodAutoscalerControllerConfigResponse{HorizontalPodAutoscalerSyncPeriod: aws.String("15s")},
			},
			expectUpdate: true,
		},
		{
			name: "terminated pod gc threshold differs - update needed",
			specConfig: &ekscontrolplanev1.KubeControllerManagerConfig{
				PodGCControllerConfig: &ekscontrolplanev1.PodGCControllerConfig{TerminatedPodGCThreshold: aws.Int32(10000)},
			},
			clusterConfig: &ekstypes.KubeControllerManagerConfigResponse{
				PodGcControllerConfig: &ekstypes.PodGcControllerConfigResponse{TerminatedPodGcThreshold: aws.Int32(12500)},
			},
			expectUpdate: true,
		},
		{
			name: "all parameters match - no update",
			specConfig: &ekscontrolplanev1.KubeControllerManagerConfig{
				HorizontalPodAutoscalerControllerConfig: &ekscontrolplanev1.HorizontalPodAutoscalerControllerConfig{HorizontalPodAutoscalerSyncPeriod: "10s"},
				PodGCControllerConfig:                   &ekscontrolplanev1.PodGCControllerConfig{TerminatedPodGCThreshold: aws.Int32(10000)},
			},
			clusterConfig: &ekstypes.KubeControllerManagerConfigResponse{
				HorizontalPodAutoscalerControllerConfig: &ekstypes.HorizontalPodAutoscalerControllerConfigResponse{HorizontalPodAutoscalerSyncPeriod: aws.String("10s")},
				PodGcControllerConfig:                   &ekstypes.PodGcControllerConfigResponse{TerminatedPodGcThreshold: aws.Int32(10000)},
			},
			expectUpdate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			s := newComponentConfigTestService(t, ekscontrolplanev1.AWSManagedControlPlaneSpec{KubeControllerManagerConfig: tc.specConfig})

			result := s.reconcileKubeControllerManagerConfig(tc.clusterConfig)
			if tc.expectUpdate {
				g.Expect(result).ToNot(BeNil(), "expected update config to be returned")
			} else {
				g.Expect(result).To(BeNil(), "expected no update")
			}
		})
	}
}

func TestReconcileControlPlaneComponentConfig(t *testing.T) {
	spec := ekscontrolplanev1.AWSManagedControlPlaneSpec{
		KubeSchedulerConfig: &ekscontrolplanev1.KubeSchedulerConfig{
			NodeResourcesFit: &ekscontrolplanev1.NodeResourcesFitConfig{
				ScoringStrategy: &ekscontrolplanev1.ScoringStrategy{Type: ekscontrolplanev1.ScoringStrategyTypeMostAllocated},
			},
		},
		KubeAPIServerConfig: &ekscontrolplanev1.KubeAPIServerConfig{EventTTL: "15m"},
	}

	tests := []struct {
		name        string
		cluster     *ekstypes.Cluster
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "no update necessary",
			cluster: &ekstypes.Cluster{
				KubeSchedulerConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated),
				KubeApiServerConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("15m")},
			},
			expect:      func(_ *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
		},
		{
			name: "only the api server config differs",
			cluster: &ekstypes.Cluster{
				KubeSchedulerConfig: schedulerResponse(ekstypes.ScoringStrategyTypeMostAllocated),
				KubeApiServerConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("60m")},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					UpdateClusterConfig(gomock.Eq(context.TODO()), gomock.Eq(&eks.UpdateClusterConfigInput{
						Name:                aws.String("test-cluster"),
						KubeApiServerConfig: &ekstypes.KubeApiServerConfigRequest{EventTtl: aws.String("15m")},
					})).
					Return(&eks.UpdateClusterConfigOutput{}, nil)
				m.WaitUntilClusterUpdating(
					gomock.Eq(context.TODO()),
					gomock.AssignableToTypeOf(&eks.DescribeClusterInput{}),
					gomock.Any(),
				).Return(nil)
			},
			expectError: false,
		},
		{
			name: "api error",
			cluster: &ekstypes.Cluster{
				KubeSchedulerConfig: schedulerResponse(ekstypes.ScoringStrategyTypeLeastAllocated),
				KubeApiServerConfig: &ekstypes.KubeApiServerConfigResponse{EventTtl: aws.String("60m")},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					UpdateClusterConfig(gomock.Eq(context.TODO()), gomock.AssignableToTypeOf(&eks.UpdateClusterConfigInput{})).
					Return(nil, errors.New("invalid parameter"))
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
			tc.expect(eksMock.EXPECT())

			s := newComponentConfigTestService(t, spec)
			s.EKSClient = eksMock

			err := s.reconcileControlPlaneComponentConfig(context.TODO(), tc.cluster)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestDurationsEqual(t *testing.T) {
	testCases := []struct {
		a, b     string
		expected bool
	}{
		{a: "15m", b: "15m", expected: true},
		{a: "60m", b: "1h", expected: true},
		{a: "15s", b: "10s", expected: false},
		{a: "15m", b: "", expected: false},
		{a: "invalid", b: "invalid", expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.a+"/"+tc.b, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(durationsEqual(tc.a, tc.b)).To(Equal(tc.expected))
		})
	}
}
