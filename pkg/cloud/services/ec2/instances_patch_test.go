package ec2

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestCreateInstance_Patch(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bootstrap-data",
		},
		Data: map[string][]byte{
			"value": []byte("data"),
		},
	}

	az := "test-zone-1a"

	data := []byte("userData")

	userDataCompressed, err := userdata.GzipBytes(data)
	if err != nil {
		t.Fatal("Failed to gzip test user data")
	}

	isUncompressedFalse := false

	testcases := []struct {
		name          string
		machine       *clusterv1.Machine
		machineConfig *infrav1.AWSMachineSpec
		awsCluster    *infrav1.AWSCluster
		expect        func(m *mocks.MockEC2APIMockRecorder)
		check         func(instance *infrav1.Instance, err error)
	}{
		{
			name: "when runInstances fails due to IAM issue in DryRun=true, we don't tag the NetworkInterface",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: aws.String("bootstrap-data"),
					},
					FailureDomain: aws.String("us-east-1c"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.2xlarge",
				Subnet: &infrav1.AWSResourceReference{
					Filters: []infrav1.Filter{
						{
							Name:   "availability-zone",
							Values: []string{"us-east-1c"},
						},
					},
				},
				UncompressedUserData: &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-foo",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:               "subnet-1",
								AvailabilityZone: "us-east-1a",
								IsPublic:         false,
							},
							infrav1.SubnetSpec{
								ID:               "subnet-2",
								AvailabilityZone: "us-east-1b",
								IsPublic:         false,
							},
							infrav1.SubnetSpec{
								ID:               "subnet-3",
								AvailabilityZone: "us-east-1c",
								IsPublic:         false,
							},
						},
					},
				},
				Status: infrav1.AWSClusterStatus{
					Network: infrav1.NetworkStatus{
						SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
							infrav1.SecurityGroupControlPlane: {
								ID: "1",
							},
							infrav1.SecurityGroupNode: {
								ID: "2",
							},
							infrav1.SecurityGroupLB: {
								ID: "3",
							},
						},
						APIServerELB: infrav1.LoadBalancer{
							DNSName: "test-apiserver.us-east-1.aws",
						},
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.2xlarge"),
						},
					})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("x86_64"),
									},
								},
							},
						},
					}, nil)
				m.DescribeSubnetsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
						{
							Name:   aws.String("availability-zone"),
							Values: aws.StringSlice([]string{"us-east-1c"}),
						},
					}})).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{
						{
							VpcId:               aws.String("vpc-bar"),
							SubnetId:            aws.String("subnet-4"),
							AvailabilityZone:    aws.String("us-east-1c"),
							CidrBlock:           aws.String("10.0.10.0/24"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
						{
							VpcId:            aws.String("vpc-foo"),
							SubnetId:         aws.String("subnet-3"),
							AvailabilityZone: aws.String("us-east-1c"),
							CidrBlock:        aws.String("10.0.11.0/24"),
						},
					},
				}, nil)
				m.
					// First call returns UnauthorizedOperation, forcing us to drop network interface tagging
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(nil, awserr.New("UnauthorizedOperation", "", nil))
				m.
					// Second call to create the instance without trying to include NetworkInterface tags
					RunInstancesWithContext(context.TODO(), &ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.2xlarge"),
						KeyName:      aws.String("default"),
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-3"),
								Groups:      aws.StringSlice([]string{"2", "3"}),
							},
						},
						TagSpecifications: []*ec2.TagSpecification{
							{
								ResourceType: aws.String("instance"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("/"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("aws-test1"),
									},
									{
										Key:   aws.String("kubernetes.io/cluster/test1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("node"),
									},
								},
							},
							{
								ResourceType: aws.String("volume"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("/"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("aws-test1"),
									},
									{
										Key:   aws.String("kubernetes.io/cluster/test1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test1"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("node"),
									},
								},
							},
						},
						UserData: aws.String(base64.StdEncoding.EncodeToString(userDataCompressed)),
						MaxCount: aws.Int64(1),
						MinCount: aws.Int64(1),
					}).Return(&ec2.Reservation{
					Instances: []*ec2.Instance{
						{
							State: &ec2.InstanceState{
								Name: aws.String(ec2.InstanceStateNamePending),
							},
							IamInstanceProfile: &ec2.IamInstanceProfile{
								Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
							},
							InstanceId:     aws.String("two"),
							InstanceType:   aws.String("m5.large"),
							SubnetId:       aws.String("subnet-3"),
							ImageId:        aws.String("ami-1"),
							RootDeviceName: aws.String("device-1"),
							BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
								{
									DeviceName: aws.String("device-1"),
									Ebs: &ec2.EbsInstanceBlockDevice{
										VolumeId: aws.String("volume-1"),
									},
								},
							},
							Placement: &ec2.Placement{
								AvailabilityZone: &az,
							},
						},
					},
				}, nil)
				m.
					DescribeNetworkInterfacesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []*ec2.NetworkInterface{},
						NextToken:         nil,
					}, nil)
			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatalf("failed to create scheme: %v", err)
			}

			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test1",
				},
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: &clusterv1.ClusterNetwork{
						ServiceDomain: "cluster.local",
						Services: &clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
						Pods: &clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
					},
				},
			}

			machine := tc.machine

			awsMachine := &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test1",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "test1",
						},
					},
				},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, cluster, machine).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client:     client,
				Cluster:    cluster,
				AWSCluster: tc.awsCluster,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
				Client:       client,
				Cluster:      cluster,
				Machine:      machine,
				AWSMachine:   awsMachine,
				InfraCluster: clusterScope,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}
			machineScope.AWSMachine.Spec = *tc.machineConfig
			tc.expect(ec2Mock.EXPECT())

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			instance, err := s.CreateInstance(machineScope, data, "")
			tc.check(instance, err)
		})
	}
}
