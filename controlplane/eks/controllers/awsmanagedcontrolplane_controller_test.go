/*
Copyright 2024 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	stsrequest "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	ec2Service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	eksService "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/network"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/securitygroup"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
)

const (
	maxActiveUpdateDeleteWait = 30 * time.Minute
)

func TestAWSManagedControlPlaneReconcilerIntegrationTests(t *testing.T) {
	var (
		reconciler AWSManagedControlPlaneReconciler
		mockCtrl   *gomock.Controller
		recorder   *record.FakeRecorder
		ctx        context.Context

		ec2Mock              *mocks.MockEC2API
		eksMock              *mock_eksiface.MockEKSAPI
		iamMock              *mock_iamauth.MockIAMAPI
		stsMock              *mock_stsiface.MockSTSAPI
		awsNodeMock          *mock_services.MockAWSNodeInterface
		iamAuthenticatorMock *mock_services.MockIAMAuthenticatorInterface
		kubeProxyMock        *mock_services.MockKubeProxyInterface
	)

	setup := func(t *testing.T) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		recorder = record.NewFakeRecorder(10)
		reconciler = AWSManagedControlPlaneReconciler{
			Client:    testEnv.Client,
			Recorder:  recorder,
			EnableIAM: true,
		}
		ctx = context.TODO()

		ec2Mock = mocks.NewMockEC2API(mockCtrl)
		eksMock = mock_eksiface.NewMockEKSAPI(mockCtrl)
		iamMock = mock_iamauth.NewMockIAMAPI(mockCtrl)
		stsMock = mock_stsiface.NewMockSTSAPI(mockCtrl)

		// Mocking these as well, since the actual implementation requires a remote client to an actual cluster
		awsNodeMock = mock_services.NewMockAWSNodeInterface(mockCtrl)
		iamAuthenticatorMock = mock_services.NewMockIAMAuthenticatorInterface(mockCtrl)
		kubeProxyMock = mock_services.NewMockKubeProxyInterface(mockCtrl)
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	t.Run("Should successfully reconcile AWSManagedControlPlane creation with managed VPC", func(t *testing.T) {
		g := NewWithT(t)
		setup(t)
		defer teardown()

		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		cluster, awsManagedCluster, awsManagedControlPlane := getManagedClusterObjects("test-cluster", ns.Name)

		// Make controller manage resources
		awsManagedControlPlane.Spec.NetworkSpec.VPC.ID = ""
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[0].ID = "my-managed-subnet-priv"
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[1].ID = "my-managed-subnet-pub1"
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[2].ID = "my-managed-subnet-pub2"

		// NAT gateway of the public subnet will be accessed by the private subnet in the same zone,
		// so use same zone for the 2 test subnets
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[0].AvailabilityZone = "us-east-1a"
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[1].AvailabilityZone = "us-east-1a"
		// Our EKS code currently requires at least 2 different AZs
		awsManagedControlPlane.Spec.NetworkSpec.Subnets[2].AvailabilityZone = "us-east-1c"

		mockedCallsForMissingEverything(ec2Mock.EXPECT(), awsManagedControlPlane.Spec.NetworkSpec.Subnets)
		mockedCreateSGCalls(ec2Mock.EXPECT())
		mockedDescribeInstanceCall(ec2Mock.EXPECT())
		mockedEKSControlPlaneIAMRole(g, iamMock.EXPECT())
		mockedEKSCluster(ctx, g, eksMock.EXPECT(), iamMock.EXPECT(), ec2Mock.EXPECT(), stsMock.EXPECT(), awsNodeMock.EXPECT(), kubeProxyMock.EXPECT(), iamAuthenticatorMock.EXPECT())

		g.Expect(testEnv.Create(ctx, &cluster)).To(Succeed())
		cluster.Status.InfrastructureReady = true
		g.Expect(testEnv.Client.Status().Update(ctx, &cluster)).To(Succeed())
		g.Expect(testEnv.Create(ctx, &awsManagedCluster)).To(Succeed())
		g.Expect(testEnv.Create(ctx, &awsManagedControlPlane)).To(Succeed())
		g.Eventually(func() bool {
			controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
			key := client.ObjectKey{
				Name:      awsManagedControlPlane.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, controlPlane)
			return err == nil
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created AWSManagedControlPlane %q", awsManagedControlPlane.Name))

		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &cluster, &awsManagedCluster, &awsManagedControlPlane, controllerIdentity, ns)).To(Succeed())
		})

		// patch the paused condition
		awsManagedControlPlanePatcher, err := patch.NewHelper(&awsManagedControlPlane, testEnv)
		awsManagedControlPlane.Status.Conditions = clusterv1.Conditions{
			{
				Type:   "Paused",
				Status: corev1.ConditionFalse,
				Reason: "NotPaused",
			},
		}

		g.Expect(awsManagedControlPlanePatcher.Patch(ctx, &awsManagedControlPlane)).To(Succeed())
		g.Expect(err).ShouldNot(HaveOccurred())

		managedScope := getAWSManagedControlPlaneScope(&cluster, &awsManagedControlPlane)

		reconciler.awsNodeServiceFactory = func(scope scope.AWSNodeScope) services.AWSNodeInterface {
			return awsNodeMock
		}

		ec2Svc := ec2Service.NewService(managedScope)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		eksSvc := eksService.NewService(managedScope)
		eksSvc.EC2Client = ec2Mock
		eksSvc.EKSClient = eksMock
		eksSvc.IAMService.IAMClient = iamMock
		eksSvc.STSClient = stsMock
		reconciler.eksServiceFactory = func(scope *scope.ManagedControlPlaneScope) *eksService.Service {
			return eksSvc
		}

		reconciler.iamAuthenticatorServiceFactory = func(scope.IAMAuthScope, iamauth.BackendType, client.Client) services.IAMAuthenticatorInterface {
			return iamAuthenticatorMock
		}
		reconciler.kubeProxyServiceFactory = func(scope scope.KubeProxyScope) services.KubeProxyInterface {
			return kubeProxyMock
		}

		networkSvc := network.NewService(managedScope)
		networkSvc.EC2Client = ec2Mock
		reconciler.networkServiceFactory = func(clusterScope scope.NetworkScope) services.NetworkInterface {
			return networkSvc
		}

		testSecurityGroupRoles := []infrav1.SecurityGroupRole{
			infrav1.SecurityGroupEKSNodeAdditional,
			infrav1.SecurityGroupBastion,
		}
		sgSvc := securitygroup.NewService(managedScope, testSecurityGroupRoles)
		sgSvc.EC2Client = ec2Mock

		reconciler.securityGroupServiceFactory = func(scope *scope.ManagedControlPlaneScope) services.SecurityGroupInterface {
			return sgSvc
		}

		_, err = reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Namespace: awsManagedControlPlane.Namespace,
				Name:      awsManagedControlPlane.Name,
			},
		})
		g.Expect(err).To(BeNil())

		g.Expect(testEnv.Get(ctx, client.ObjectKeyFromObject(&awsManagedControlPlane), &awsManagedControlPlane)).To(Succeed())
		g.Expect(awsManagedControlPlane.Finalizers).To(ContainElement(ekscontrolplanev1.ManagedControlPlaneFinalizer))
	})
}

func createControllerIdentity(g *WithT) *infrav1.AWSClusterControllerIdentity {
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			Kind: string(infrav1.ControllerIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	g.Expect(testEnv.Create(ctx, controllerIdentity)).To(Succeed())
	return controllerIdentity
}

// mockedCallsForMissingEverything mocks most of the AWSManagedControlPlane reconciliation calls to the AWS API,
// except for what other functions provide (see `mockedCreateSGCalls` and `mockedDescribeInstanceCall`).
func mockedCallsForMissingEverything(ec2Rec *mocks.MockEC2APIMockRecorder, subnets infrav1.Subnets) {
	describeVPCByNameCall := ec2Rec.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{"test-cluster-vpc"},
			},
		},
	})).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{},
	}, nil)

	ec2Rec.CreateVpc(context.TODO(), gomock.Eq(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/8"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVpc,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-vpc"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
				},
			},
		},
	})).After(describeVPCByNameCall).Return(&ec2.CreateVpcOutput{
		Vpc: &ec2types.Vpc{
			State:     ec2types.VpcStateAvailable,
			VpcId:     aws.String("vpc-new"),
			CidrBlock: aws.String("10.0.0.0/8"),
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-cluster-vpc"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
					Value: aws.String("common"),
				},
			},
		},
	}, nil)

	ec2Rec.DescribeVpcAttribute(context.TODO(), gomock.Eq(&ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String("vpc-new"),
		Attribute: ec2types.VpcAttributeNameEnableDnsHostnames,
	})).Return(&ec2.DescribeVpcAttributeOutput{
		EnableDnsHostnames: &ec2types.AttributeBooleanValue{Value: aws.Bool(true)},
	}, nil)

	ec2Rec.DescribeVpcAttribute(context.TODO(), gomock.Eq(&ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String("vpc-new"),
		Attribute: ec2types.VpcAttributeNameEnableDnsSupport,
	})).Return(&ec2.DescribeVpcAttributeOutput{
		EnableDnsSupport: &ec2types.AttributeBooleanValue{Value: aws.Bool(true)},
	}, nil)

	ec2Rec.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
		},
	})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{},
	}, nil)

	zones := []ec2types.AvailabilityZone{}
	for _, subnet := range subnets {
		zones = append(zones, ec2types.AvailabilityZone{
			ZoneName: aws.String(subnet.AvailabilityZone),
			ZoneType: aws.String("availability-zone"),
		})
	}
	ec2Rec.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
		Return(&ec2.DescribeAvailabilityZonesOutput{
			AvailabilityZones: zones,
		}, nil).MaxTimes(2)

	for subnetIndex, subnet := range subnets {
		subnetID := fmt.Sprintf("subnet-%d", subnetIndex+1)
		var kubernetesRoleTagKey string
		var capaRoleTagValue string
		if subnet.IsPublic {
			kubernetesRoleTagKey = "kubernetes.io/role/elb"
			capaRoleTagValue = "public"
		} else {
			kubernetesRoleTagKey = "kubernetes.io/role/internal-elb"
			capaRoleTagValue = "private"
		}
		ec2Rec.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
			ZoneNames: []string{subnet.AvailabilityZone},
		}).
			Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					{
						ZoneName: aws.String(subnet.AvailabilityZone),
						ZoneType: aws.String("availability-zone"),
					},
				},
			}, nil).MaxTimes(1)
		ec2Rec.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
			VpcId:            aws.String("vpc-new"),
			CidrBlock:        aws.String(subnet.CidrBlock),
			AvailabilityZone: aws.String(subnet.AvailabilityZone),
			TagSpecifications: []ec2types.TagSpecification{
				{
					ResourceType: ec2types.ResourceTypeSubnet,
					Tags: []ec2types.Tag{
						{
							Key: aws.String("Name"),
							// Assume that `ID` doesn't start with `subnet-` so that it becomes managed and `ID` denotes the desired name
							Value: aws.String(subnet.ID),
						},
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String(kubernetesRoleTagKey),
							Value: aws.String("1"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String(capaRoleTagValue),
						},
					},
				},
			},
		})).Return(&ec2.CreateSubnetOutput{
			Subnet: &ec2types.Subnet{
				VpcId:               aws.String("vpc-new"),
				SubnetId:            aws.String(subnetID),
				CidrBlock:           aws.String(subnet.CidrBlock),
				AvailabilityZone:    aws.String(subnet.AvailabilityZone),
				MapPublicIpOnLaunch: aws.Bool(false),
				Tags: []ec2types.Tag{
					{
						Key: aws.String("Name"),
						// Assume that `ID` doesn't start with `subnet-` so that it becomes managed and `ID` denotes the desired name
						Value: aws.String(subnet.ID),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("kubernetes.io/role/elb"),
						Value: aws.String("1"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("public"),
					},
				},
			},
		}, nil)

		ec2Rec.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
			SubnetIds: []string{subnetID},
		}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
			Subnets: []ec2types.Subnet{
				{
					VpcId:               aws.String("vpc-new"),
					SubnetId:            aws.String(subnetID),
					CidrBlock:           aws.String(subnet.CidrBlock),
					AvailabilityZone:    aws.String(subnet.AvailabilityZone),
					MapPublicIpOnLaunch: aws.Bool(false),
					State:               ec2types.SubnetStateAvailable,
				},
			},
		}, nil)

		if subnet.IsPublic {
			ec2Rec.ModifySubnetAttribute(context.TODO(), gomock.Eq(&ec2.ModifySubnetAttributeInput{
				SubnetId: aws.String(subnetID),
				MapPublicIpOnLaunch: &ec2types.AttributeBooleanValue{
					Value: aws.Bool(true),
				},
			})).Return(&ec2.ModifySubnetAttributeOutput{}, nil)
		}
	}

	ec2Rec.DescribeRouteTables(context.TODO(), gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		}})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{
			{
				Routes: []ec2types.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
			},
		},
	}, nil).MinTimes(1).MaxTimes(2)

	ec2Rec.DescribeInternetGateways(context.TODO(), gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{"vpc-new"},
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []ec2types.InternetGateway{},
	}, nil)

	ec2Rec.CreateInternetGateway(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateInternetGatewayInput{})).
		Return(&ec2.CreateInternetGatewayOutput{
			InternetGateway: &ec2types.InternetGateway{
				InternetGatewayId: aws.String("igw-1"),
				Tags: []ec2types.Tag{
					{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster")),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-igw"),
					},
				},
			},
		}, nil)

	ec2Rec.AttachInternetGateway(context.TODO(), gomock.Eq(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String("igw-1"),
		VpcId:             aws.String("vpc-new"),
	})).
		Return(&ec2.AttachInternetGatewayOutput{}, nil)

	ec2Rec.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		}}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil).MinTimes(1).MaxTimes(2)

	ec2Rec.DescribeAddresses(context.TODO(), gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: []string{"common"},
			},
		},
	})).Return(&ec2.DescribeAddressesOutput{
		Addresses: []ec2types.Address{},
	}, nil)

	for subnetIndex, subnet := range subnets {
		subnetID := fmt.Sprintf("subnet-%d", subnetIndex+1)

		// NAT gateways are attached to public subnets
		if subnet.IsPublic {
			eipAllocationID := strconv.Itoa(1234 + subnetIndex)
			natGatewayID := fmt.Sprintf("nat-%d", subnetIndex+1)

			ec2Rec.AllocateAddress(context.TODO(), gomock.Eq(&ec2.AllocateAddressInput{
				Domain: ec2types.DomainTypeVpc,
				TagSpecifications: []ec2types.TagSpecification{
					{
						ResourceType: ec2types.ResourceTypeElasticIp,
						Tags: []ec2types.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("test-cluster-eip-common"),
							},
							{
								Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
								Value: aws.String("owned"),
							},
							{
								Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
								Value: aws.String("common"),
							},
						},
					},
				},
			})).Return(&ec2.AllocateAddressOutput{
				AllocationId: aws.String(eipAllocationID),
			}, nil)

			ec2Rec.CreateNatGateway(context.TODO(), gomock.Eq(&ec2.CreateNatGatewayInput{
				AllocationId: aws.String(eipAllocationID),
				SubnetId:     aws.String(subnetID),
				TagSpecifications: []ec2types.TagSpecification{
					{
						ResourceType: ec2types.ResourceTypeNatgateway,
						Tags: []ec2types.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("test-cluster-nat"),
							},
							{
								Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
								Value: aws.String("owned"),
							},
							{
								Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
								Value: aws.String("common"),
							},
						},
					},
				},
			})).Return(&ec2.CreateNatGatewayOutput{
				NatGateway: &ec2types.NatGateway{
					NatGatewayId: aws.String(natGatewayID),
					SubnetId:     aws.String(subnetID),
				},
			}, nil)

			ec2Rec.DescribeNatGateways(gomock.Any(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
				NatGatewayIds: []string{natGatewayID},
			}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
				NatGateways: []ec2types.NatGateway{
					{
						NatGatewayId: aws.String(natGatewayID),
						SubnetId:     aws.String(subnetID),
						State:        ec2types.NatGatewayStateAvailable,
					},
				},
			}, nil)
		}

		routeTableID := fmt.Sprintf("rtb-%d", subnetIndex+1)
		var routeTablePublicPrivate string
		if subnet.IsPublic {
			routeTablePublicPrivate = "public"
		} else {
			routeTablePublicPrivate = "private"
		}
		ec2Rec.CreateRouteTable(context.TODO(), gomock.Eq(&ec2.CreateRouteTableInput{
			TagSpecifications: []ec2types.TagSpecification{
				{
					ResourceType: ec2types.ResourceTypeRouteTable,
					Tags: []ec2types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String(fmt.Sprintf("test-cluster-rt-%s-%s", routeTablePublicPrivate, subnet.AvailabilityZone)),
						},
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("common"),
						},
					},
				},
			},
			VpcId: aws.String("vpc-new"),
		})).Return(&ec2.CreateRouteTableOutput{
			RouteTable: &ec2types.RouteTable{
				RouteTableId: aws.String(routeTableID),
			},
		}, nil)

		if subnet.IsPublic {
			ec2Rec.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				GatewayId:            aws.String("igw-1"),
				RouteTableId:         aws.String(routeTableID),
			})).Return(&ec2.CreateRouteOutput{}, nil)
		} else {
			// Private subnet uses a NAT gateway attached to a public subnet in the same AZ
			var natGatewayID string
			for otherSubnetIndex, otherSubnet := range subnets {
				if otherSubnet.IsPublic && subnet.AvailabilityZone == otherSubnet.AvailabilityZone {
					natGatewayID = fmt.Sprintf("nat-%d", otherSubnetIndex+1)
					break
				}
			}
			if natGatewayID == "" {
				panic("Could not find NAT gateway from public subnet of same AZ")
			}
			ec2Rec.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				NatGatewayId:         aws.String(natGatewayID),
				RouteTableId:         aws.String(routeTableID),
			})).Return(&ec2.CreateRouteOutput{}, nil)
		}

		ec2Rec.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
			RouteTableId: aws.String(routeTableID),
			SubnetId:     aws.String(subnetID),
		})).Return(&ec2.AssociateRouteTableOutput{}, nil)
	}
}

func mockedCreateSGCalls(ec2Rec *mocks.MockEC2APIMockRecorder) {
	ec2Rec.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(
		&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []ec2types.SecurityGroup{
				{
					GroupId:   aws.String("1"),
					GroupName: aws.String("test-sg"),
				},
			},
		}, nil)
	securityGroupAdditionalCall := ec2Rec.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-new"),
		GroupName:   aws.String("test-cluster-node-eks-additional"),
		Description: aws.String("Kubernetes cluster test-cluster: node-eks-additional"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-node-eks-additional"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("node-eks-additional"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-node-eks-additional")}, nil)
	ec2Rec.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-new"),
		GroupName:   aws.String("test-cluster-bastion"),
		Description: aws.String("Kubernetes cluster test-cluster: bastion"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-bastion"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("bastion"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-bastion")}, nil)
	ec2Rec.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-node-eks-additional"),
	})).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
		After(securityGroupAdditionalCall).Times(2)
}

func mockedDescribeInstanceCall(ec2Rec *mocks.MockEC2APIMockRecorder) {
	ec2Rec.DescribeInstances(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: []string{"bastion"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"pending", "running", "stopping", "stopped"},
			},
		},
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{
			{
				Instances: []ec2types.Instance{
					{
						InstanceId:   aws.String("id-1"),
						InstanceType: ec2types.InstanceTypeM5Large,
						SubnetId:     aws.String("subnet-1"),
						ImageId:      aws.String("ami-1"),
						IamInstanceProfile: &ec2types.IamInstanceProfile{
							Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
						},
						State: &ec2types.InstanceState{
							Code: aws.Int32(16),
							Name: ec2types.InstanceStateNameRunning,
						},
						RootDeviceName: aws.String("device-1"),
						BlockDeviceMappings: []ec2types.InstanceBlockDeviceMapping{
							{
								DeviceName: aws.String("device-1"),
								Ebs: &ec2types.EbsInstanceBlockDevice{
									VolumeId: aws.String("volume-1"),
								},
							},
						},
						Placement: &ec2types.Placement{
							AvailabilityZone: aws.String("us-east-1a"),
						},
					},
				},
			},
		},
	}, nil)
}

func mockedEKSControlPlaneIAMRole(g *WithT, iamRec *mock_iamauth.MockIAMAPIMockRecorder) {
	getRoleCall := iamRec.GetRole(gomock.Any(), &iam.GetRoleInput{
		RoleName: aws.String("test-cluster-iam-service-role"),
	}).Return(nil, &smithy.GenericAPIError{Code: "NoSuchEntity", Message: ""})

	createRoleCall := iamRec.CreateRole(gomock.Any(), gomock.Any()).After(getRoleCall).DoAndReturn(func(ctx context.Context, input *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
		g.Expect(input.RoleName).To(BeComparableTo(aws.String("test-cluster-iam-service-role")))
		return &iam.CreateRoleOutput{
			Role: &iamtypes.Role{
				RoleName: aws.String("test-cluster-iam-service-role"),
				Arn:      aws.String("arn:aws:iam::123456789012:role/test-cluster-iam-service-role"),
				Tags:     input.Tags,
			},
		}, nil
	})

	iamRec.ListAttachedRolePolicies(gomock.Any(), &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String("test-cluster-iam-service-role"),
	}).After(createRoleCall).Return(&iam.ListAttachedRolePoliciesOutput{
		AttachedPolicies: []iamtypes.AttachedPolicy{},
	}, nil)

	getPolicyCall := iamRec.GetPolicy(gomock.Any(), &iam.GetPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
	}).Return(&iam.GetPolicyOutput{
		// This policy is predefined by AWS
		Policy: &iamtypes.Policy{
			// Fields are not used. Our code only checks for existence of the policy.
		},
	}, nil)

	iamRec.AttachRolePolicy(gomock.Any(), &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
		RoleName:  aws.String("test-cluster-iam-service-role"),
	}).After(getPolicyCall).Return(&iam.AttachRolePolicyOutput{}, nil)
}

func mockedEKSCluster(ctx context.Context, g *WithT, eksRec *mock_eksiface.MockEKSAPIMockRecorder, iamRec *mock_iamauth.MockIAMAPIMockRecorder, ec2Rec *mocks.MockEC2APIMockRecorder, stsRec *mock_stsiface.MockSTSAPIMockRecorder, awsNodeRec *mock_services.MockAWSNodeInterfaceMockRecorder, kubeProxyRec *mock_services.MockKubeProxyInterfaceMockRecorder, iamAuthenticatorRec *mock_services.MockIAMAuthenticatorInterfaceMockRecorder) {
	describeClusterCall := eksRec.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String("test-cluster"),
	}).Return(nil, &ekstypes.ResourceNotFoundException{
		Message: aws.String("cluster not found"),
	})

	getRoleCall := iamRec.GetRole(gomock.Any(), &iam.GetRoleInput{
		RoleName: aws.String("test-cluster-iam-service-role"),
	}).After(describeClusterCall).Return(&iam.GetRoleOutput{
		Role: &iamtypes.Role{
			RoleName: aws.String("test-cluster-iam-service-role"),
			Arn:      aws.String("arn:aws:iam::123456789012:role/test-cluster-iam-service-role"),
		},
	}, nil)

	resourcesVpcConfig := &ekstypes.VpcConfigResponse{
		ClusterSecurityGroupId: aws.String("eks-cluster-sg-test-cluster-44556677"),
	}

	clusterARN := aws.String("arn:aws:eks:us-east-1:1133557799:cluster/test-cluster")
	clusterCreating := ekstypes.Cluster{
		Arn:                clusterARN,
		Name:               aws.String("test-cluster"),
		Status:             ekstypes.ClusterStatusCreating,
		ResourcesVpcConfig: resourcesVpcConfig,
		CertificateAuthority: &ekstypes.Certificate{
			Data: aws.String(base64.StdEncoding.EncodeToString([]byte("foobar"))),
		},
		Logging: &ekstypes.Logging{
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
	}

	createClusterCall := eksRec.CreateCluster(ctx, gomock.Any()).After(getRoleCall).DoAndReturn(func(ctx context.Context, input *eks.CreateClusterInput, optFns ...func(*eks.Options)) (*eks.CreateClusterOutput, error) {
		g.Expect(input.Name).To(BeComparableTo(aws.String("test-cluster")))
		return &eks.CreateClusterOutput{
			Cluster: &clusterCreating,
		}, nil
	})

	waitUntilClusterActiveCall := eksRec.WaitUntilClusterActive(ctx, &eks.DescribeClusterInput{
		Name: aws.String("test-cluster"),
	}, maxActiveUpdateDeleteWait).After(createClusterCall).Return(nil)

	clusterActive := clusterCreating // copy
	clusterActive.Status = ekstypes.ClusterStatusActive
	clusterActive.Endpoint = aws.String("https://F00D133712341337.gr7.us-east-1.eks.amazonaws.com")
	clusterActive.Version = aws.String("1.24")

	eksRec.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String("test-cluster"),
	}).After(waitUntilClusterActiveCall).Return(&eks.DescribeClusterOutput{
		Cluster: &clusterActive,
	}, nil)

	// AWS precreates a default security group together with the cluster
	// (https://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html)
	clusterSgDesc := &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []ec2types.SecurityGroup{
			{
				GroupId:   aws.String("sg-11223344"),
				GroupName: aws.String("eks-cluster-sg-test-cluster-44556677"),
			},
		},
	}
	ec2Rec.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:aws:eks:cluster-name"),
				Values: []string{"test-cluster"},
			},
		},
	})).Return(
		clusterSgDesc, nil)
	ec2Rec.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{"eks-cluster-sg-test-cluster-44556677"},
	})).Return(
		clusterSgDesc, nil)

	req, err := http.NewRequest(http.MethodGet, "foobar", http.NoBody)
	g.Expect(err).To(BeNil())
	stsRec.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{}).Return(&stsrequest.Request{
		HTTPRequest: req,
		Operation:   &stsrequest.Operation{},
	}, &sts.GetCallerIdentityOutput{})

	eksRec.TagResource(ctx, &eks.TagResourceInput{
		ResourceArn: clusterARN,
		Tags: map[string]string{
			"Name": "test-cluster",
			"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
			"sigs.k8s.io/cluster-api-provider-aws/role":                 "common",
		},
	}).Return(&eks.TagResourceOutput{}, nil)

	eksRec.ListAddons(ctx, &eks.ListAddonsInput{
		ClusterName: aws.String("test-cluster"),
	}).Return(&eks.ListAddonsOutput{}, nil)

	eksRec.UpdateClusterConfig(ctx, gomock.AssignableToTypeOf(&eks.UpdateClusterConfigInput{})).After(waitUntilClusterActiveCall).Return(&eks.UpdateClusterConfigOutput{}, nil)

	awsNodeRec.ReconcileCNI(gomock.Any()).Return(nil)
	kubeProxyRec.ReconcileKubeProxy(gomock.Any()).Return(nil)
	iamAuthenticatorRec.ReconcileIAMAuthenticator(gomock.Any()).Return(nil)
}
