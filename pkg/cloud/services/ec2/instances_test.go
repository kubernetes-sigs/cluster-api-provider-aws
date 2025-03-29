/*
Copyright 2018 The Kubernetes Authors.

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

package ec2

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestInstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mocks.MockEC2APIMockRecorder)
		check      func(instance *infrav1.Instance, err error)
	}{
		{
			name:       "does not exist",
			instanceID: "hello",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("hello")},
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
			check: func(instance *infrav1.Instance, err error) {
				if err == nil {
					t.Fatalf("expects error when instance could not be found: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "does not exist with bad request error",
			instanceID: "hello-does-not-exist",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("hello-does-not-exist")},
				})).
					Return(nil, awserr.New(awserrors.InvalidInstanceID, "does not exist", nil))
			},
			check: func(instance *infrav1.Instance, err error) {
				if err == nil {
					t.Fatalf("expects error when DescribeInstances returns error: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "instance exists",
			instanceID: "id-1",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				az := "test-zone-1a"
				m.DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("id-1")},
				})).
					Return(&ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										InstanceId:   aws.String("id-1"),
										InstanceType: aws.String("m5.large"),
										SubnetId:     aws.String("subnet-1"),
										ImageId:      aws.String("ami-1"),
										IamInstanceProfile: &ec2.IamInstanceProfile{
											Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
										},
										State: &ec2.InstanceState{
											Code: aws.Int64(16),
											Name: aws.String(ec2.StateAvailable),
										},
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
							},
						},
					}, nil)
			},
			check: func(instance *infrav1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance == nil {
					t.Fatalf("expected instance but got nothing")
				}

				if instance != nil && instance.ID != "id-1" {
					t.Fatalf("expected id-1 but got: %v", instance.ID)
				}
			},
		},
		{
			name:       "error describing instances",
			instanceID: "one",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeInstancesWithContext(context.TODO(), &ec2.DescribeInstancesInput{
					InstanceIds: []*string{aws.String("one")},
				}).
					Return(nil, errors.New("some unknown error"))
			},
			check: func(i *infrav1.Instance, err error) {
				if err == nil {
					t.Fatalf("expected an error but got none.")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client:  client,
				Cluster: &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "test-vpc",
							},
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			instance, err := s.InstanceIfExists(aws.String(tc.instanceID))
			tc.check(instance, err)
		})
	}
}

func TestTerminateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	instanceNotFoundError := errors.New("instance not found")

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mocks.MockEC2APIMockRecorder)
		check      func(err error)
	}{
		{
			name:       "instance exists",
			instanceID: "i-exist",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.TerminateInstancesWithContext(context.TODO(), gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-exist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, nil)
			},
			check: func(err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name:       "instance does not exist",
			instanceID: "i-donotexist",
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.TerminateInstancesWithContext(context.TODO(), gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []*string{aws.String("i-donotexist")},
				})).
					Return(&ec2.TerminateInstancesOutput{}, instanceNotFoundError)
			},
			check: func(err error) {
				if err == nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client:     client,
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			err = s.TerminateInstance(tc.instanceID)
			tc.check(err)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bootstrap-data",
		},
		Data: map[string][]byte{
			"value": []byte("data"),
		},
	}

	az := "test-zone-1a"
	tenancy := "dedicated"

	data := []byte("userData")

	userDataCompressed, err := userdata.GzipBytes(data)
	if err != nil {
		t.Fatal("Failed to gzip test user data")
	}

	isUncompressedFalse := false
	isUncompressedTrue := true

	testcases := []struct {
		name          string
		machine       *clusterv1.Machine
		machineConfig *infrav1.AWSMachineSpec
		awsCluster    *infrav1.AWSCluster
		expect        func(m *mocks.MockEC2APIMockRecorder)
		check         func(instance *infrav1.Instance, err error)
	}{
		{
			name: "simple",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "with availability zone",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					FailureDomain: aws.String("us-east-1c"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.2xlarge",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
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
							infrav1.SubnetSpec{
								ID:               "subnet-3-public",
								AvailabilityZone: "us-east-1c",
								IsPublic:         true,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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

				if instance.SubnetID != "subnet-3" {
					t.Fatalf("expected subnet-3 from availability zone us-east-1c, got %q", instance.SubnetID)
				}
			},
		},
		{
			name: "when multiple subnets match filters, subnets in the cluster vpc are preferred",
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
							VpcId:               aws.String("vpc-incorrect-1"),
							SubnetId:            aws.String("subnet-5"),
							AvailabilityZone:    aws.String("us-east-1c"),
							CidrBlock:           aws.String("10.0.12.0/24"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
						{
							VpcId:               aws.String("vpc-incorrect-2"),
							SubnetId:            aws.String("subnet-4"),
							AvailabilityZone:    aws.String("us-east-1c"),
							CidrBlock:           aws.String("10.0.10.0/24"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
						{
							VpcId:               aws.String("vpc-foo"),
							SubnetId:            aws.String("subnet-3"),
							AvailabilityZone:    aws.String("us-east-1c"),
							CidrBlock:           aws.String("10.0.11.0/24"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					},
				}, nil)
				m.
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
							{
								ResourceType: aws.String("network-interface"),
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

				if instance.SubnetID != "subnet-3" {
					t.Fatalf("expected subnet-3 from availability zone us-east-1c, got %q", instance.SubnetID)
				}
			},
		},
		{
			name: "with a subnet outside the cluster vpc",
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
							Name:   "vpc-id",
							Values: []string{"vpc-bar"},
						},
						{
							Name:   "availability-zone",
							Values: []string{"us-east-1c"},
						},
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{
					infrav1.SecurityGroupNode: "4",
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
						filter.EC2.VPC("vpc-bar"),
						{
							Name:   aws.String("availability-zone"),
							Values: aws.StringSlice([]string{"us-east-1c"}),
						},
					}})).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{
						{
							VpcId:            aws.String("vpc-bar"),
							SubnetId:         aws.String("subnet-5"),
							AvailabilityZone: aws.String("us-east-1c"),
							CidrBlock:        aws.String("10.0.11.0/24"),
						},
					},
				}, nil)
				m.
					RunInstancesWithContext(context.TODO(), &ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.2xlarge"),
						KeyName:      aws.String("default"),
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-5"),
								Groups:      aws.StringSlice([]string{"4", "3"}),
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
							{
								ResourceType: aws.String("network-interface"),
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
							SubnetId:       aws.String("subnet-5"),
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

				if instance.SubnetID != "subnet-5" {
					t.Fatalf("expected subnet-5 from availability zone us-east-1c, got %q", instance.SubnetID)
				}
			},
		},
		{
			name: "with ImageLookupOrg specified at the machine level",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				ImageLookupOrg: "test-org-123",
				InstanceType:   "m6g.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-24.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m6g.large"),
						},
					})).
					Return(&ec2.DescribeInstanceTypesOutput{
						InstanceTypes: []*ec2.InstanceTypeInfo{
							{
								ProcessorInfo: &ec2.ProcessorInfo{
									SupportedArchitectures: []*string{
										aws.String("arm64"),
									},
								},
							},
						},
					}, nil)
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("test-org-123")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName)},
							},
							{
								Name:   aws.String("architecture"),
								Values: []*string{aws.String("arm64")},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("available")},
							},
							{
								Name:   aws.String("virtualization-type"),
								Values: []*string{aws.String("hvm")},
							},
						},
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "with ImageLookupOrg specified at the cluster-level",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					ImageLookupOrg: "cluster-level-image-lookup-org",
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-24.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("cluster-level-image-lookup-org")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName)},
							},
							{
								Name:   aws.String("architecture"),
								Values: []*string{aws.String("x86_64")},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("available")},
							},
							{
								Name:   aws.String("virtualization-type"),
								Values: []*string{aws.String("hvm")},
							},
						},
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "AWSMachine ImageLookupOrg overrides AWSCluster ImageLookupOrg",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType:   "m5.large",
				ImageLookupOrg: "machine-level-image-lookup-org",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					ImageLookupOrg: "cluster-level-image-lookup-org",
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-24.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				// verify that the ImageLookupOrg is used when finding AMIs
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("owner-id"),
								Values: []*string{aws.String("machine-level-image-lookup-org")},
							},
							{
								Name:   aws.String("name"),
								Values: []*string{aws.String(amiName)},
							},
							{
								Name:   aws.String("architecture"),
								Values: []*string{aws.String("x86_64")},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("available")},
							},
							{
								Name:   aws.String("virtualization-type"),
								Values: []*string{aws.String("hvm")},
							},
						},
					})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2006-01-02T15:04:05.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "subnet filter and failureDomain defined",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					FailureDomain: aws.String("us-east-1b"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					Filters: []infrav1.Filter{{
						Name:   "tag:some-tag",
						Values: []string{"some-value"},
					}},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
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
							aws.String("m5.large"),
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
				m.
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("tag:some-tag"), Values: aws.StringSlice([]string{"some-value"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:         aws.String("filtered-subnet-1"),
							AvailabilityZone: aws.String("us-east-1b"),
						}},
					}, nil)
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "with subnet ID that belongs to Cluster",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("matching-subnet"),
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID: "matching-subnet",
						}},
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
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"matching-subnet"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:         aws.String("matching-subnet"),
							AvailabilityZone: aws.String("us-east-1b"),
						}},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("matching-subnet"),
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
		{
			name: "with subnet ID that does not exist",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("non-matching-subnet"),
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID: "subnet-1",
						}},
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
							aws.String("m5.large"),
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
				m.
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"non-matching-subnet"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{},
					}, nil)
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "failed to run machine \"aws-test1\", no subnets available matching criteria"
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: %s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "with subnet ID that does not belong to Cluster",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("matching-subnet"),
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID: "subnet-1",
						}},
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
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("matching-subnet"),
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
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				m.
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"matching-subnet"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId: aws.String("matching-subnet"),
						}},
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
		{
			name: "subnet id and failureDomain don't match",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					FailureDomain: aws.String("us-east-1b"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("subnet-1"),
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID:               "subnet-1",
							AvailabilityZone: "us-west-1b",
						}},
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
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"subnet-1"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:         aws.String("subnet-1"),
							AvailabilityZone: aws.String("us-west-1b"),
						}},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "failed to run machine \"aws-test1\", found 1 subnets matching criteria but post-filtering failed. subnet \"subnet-1\" availability zone \"us-west-1b\" does not match failure domain \"us-east-1b\""
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: `%s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "public IP true and failureDomain doesn't have public subnet",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					FailureDomain: aws.String("us-east-1b"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				PublicIP:     aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID:               "private-subnet-1",
							AvailabilityZone: "us-east-1b",
							IsPublic:         false,
						}},
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
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "failed to run machine \"aws-test1\" with public IP, no public subnets available in availability zone \"us-east-1b\""
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: `%s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "public IP true and public subnet ID given",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("public-subnet-1"),
				},
				PublicIP: aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID:       "public-subnet-1",
							IsPublic: true,
						}},
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
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"public-subnet-1"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("public-subnet-1"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(true),
						}},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "public IP true, public subnet ID given and MapPublicIpOnLaunch is false",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("public-subnet-1"),
				},
				PublicIP: aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID:       "public-subnet-1",
							IsPublic: true,
						}},
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
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"public-subnet-1"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("public-subnet-1"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						}},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Do(func(_ context.Context, in *ec2.RunInstancesInput, _ ...request.Option) {
						if len(in.NetworkInterfaces) == 0 {
							t.Fatalf("expected a NetworkInterface to be defined")
						}
						if !aws.BoolValue(in.NetworkInterfaces[0].AssociatePublicIpAddress) {
							t.Fatalf("expected AssociatePublicIpAddress to be set and true")
						}
						if subnet := aws.StringValue(in.NetworkInterfaces[0].SubnetId); subnet != "public-subnet-1" {
							t.Fatalf("expected subnet ID to be \"public-subnet-1\", got %q", subnet)
						}
						if in.NetworkInterfaces[0].Groups == nil {
							t.Fatalf("expected security groups to be set")
						}
					}).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "efa interface type",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:         "m5.large",
				NetworkInterfaceType: infrav1.NetworkInterfaceTypeEFAWithENAInterface,
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Do(func(_ context.Context, in *ec2.RunInstancesInput, _ ...request.Option) {
						if len(in.NetworkInterfaces) == 0 {
							t.Fatalf("expected a NetworkInterface to be defined")
						}
						if in.NetworkInterfaces[0].Groups == nil {
							t.Fatalf("expected security groups to be set")
						}
						if interfaceType := aws.StringValue(in.NetworkInterfaces[0].InterfaceType); interfaceType != "efa" {
							t.Fatalf("expected interface type to be \"efa\": got %q", interfaceType)
						}
					}).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
		{
			name: "public IP true and private subnet ID given",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					ID: aws.String("private-subnet-1"),
				},
				PublicIP: aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{{
							ID:       "private-subnet-1",
							IsPublic: false,
						}},
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
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"private-subnet-1"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("private-subnet-1"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						}},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "failed to run machine \"aws-test1\", found 1 subnets matching criteria but post-filtering failed. subnet \"private-subnet-1\" is a private subnet."
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: `%s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "both public IP and subnet filter defined",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					Filters: []infrav1.Filter{{
						Name:   "tag:some-tag",
						Values: []string{"some-value"},
					}},
				},
				PublicIP: aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "private-subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								ID:       "public-subnet-1",
								IsPublic: true,
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
							aws.String("m5.large"),
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
				m.
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("tag:some-tag"), Values: aws.StringSlice([]string{"some-value"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("public-subnet-1"),
							MapPublicIpOnLaunch: aws.Bool(true),
						}},
					}, nil)
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "both public IP, subnet filter defined and MapPublicIpOnLaunch is false",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				Subnet: &infrav1.AWSResourceReference{
					Filters: []infrav1.Filter{{
						Name:   "tag:some-tag",
						Values: []string{"some-value"},
					}},
				},
				PublicIP: aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "private-subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								ID:       "public-subnet-1",
								IsPublic: true,
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
							aws.String("m5.large"),
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
				m.
					DescribeSubnetsWithContext(context.TODO(), &ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							{Name: aws.String("tag:some-tag"), Values: aws.StringSlice([]string{"some-value"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("public-subnet-1"),
							MapPublicIpOnLaunch: aws.Bool(false),
						}},
					}, nil)
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Do(func(_ context.Context, in *ec2.RunInstancesInput, _ ...request.Option) {
						if len(in.NetworkInterfaces) == 0 {
							t.Fatalf("expected a NetworkInterface to be defined")
						}
						if !aws.BoolValue(in.NetworkInterfaces[0].AssociatePublicIpAddress) {
							t.Fatalf("expected AssociatePublicIpAddress to be set and true")
						}
						if subnet := aws.StringValue(in.NetworkInterfaces[0].SubnetId); subnet != "public-subnet-1" {
							t.Fatalf("expected subnet ID to be \"public-subnet-1\", got %q", subnet)
						}
						if in.NetworkInterfaces[0].Groups == nil {
							t.Fatalf("expected security groups to be set")
						}
					}).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "public IP true and public subnet exists",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				PublicIP:     aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "private-subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								ID:       "public-subnet-1",
								IsPublic: true,
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
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "public IP true, public subnet exists and MapPublicIpOnLaunch is false",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				PublicIP:     aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "private-subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								ID:       "public-subnet-1",
								IsPublic: true,
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
							aws.String("m5.large"),
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
				m.
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Do(func(_ context.Context, in *ec2.RunInstancesInput, _ ...request.Option) {
						if len(in.NetworkInterfaces) == 0 {
							t.Fatalf("expected a NetworkInterface to be defined")
						}
						if !aws.BoolValue(in.NetworkInterfaces[0].AssociatePublicIpAddress) {
							t.Fatalf("expected AssociatePublicIpAddress to be set and true")
						}
						if subnet := aws.StringValue(in.NetworkInterfaces[0].SubnetId); subnet != "public-subnet-1" {
							t.Fatalf("expected subnet ID to be \"public-subnet-1\", got %q", subnet)
						}
						if in.NetworkInterfaces[0].Groups == nil {
							t.Fatalf("expected security groups to be set")
						}
					}).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("public-subnet-1"),
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
		{
			name: "public IP true and no public subnet exists",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				PublicIP:     aws.Bool(true),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "vpc-id",
						},
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "private-subnet-1",
								IsPublic: false,
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
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "failed to run machine \"aws-test1\" with public IP, no public subnets available"
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}

				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: %s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "with multiple block device mappings",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
				NonRootVolumes: []infrav1.Volume{{
					DeviceName: "device-2",
					Size:       8,
				}},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
								ImageId:        aws.String("ami-1"),
								RootDeviceName: aws.String("device-1"),
								BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
									{
										DeviceName: aws.String("device-1"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-1"),
										},
									},
									{
										DeviceName: aws.String("device-2"),
										Ebs: &ec2.EbsInstanceBlockDevice{
											VolumeId: aws.String("volume-2"),
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
		{
			name: "with dedicated tenancy cloud-config",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:         "m5.large",
				Tenancy:              "dedicated",
				UncompressedUserData: &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							Tenancy: &tenancy,
						},
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-1"),
								Groups:      []*string{aws.String("2"), aws.String("3")},
							},
						},
						TagSpecifications: []*ec2.TagSpecification{
							{
								ResourceType: aws.String("instance"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
										Value: aws.String("default/machine-aws-test1"),
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
								ResourceType: aws.String("network-interface"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
					})).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
									Tenancy:          &tenancy,
								},
							},
						},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
		{
			name: "with custom placement group cloud-config",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:         "m5.large",
				PlacementGroupName:   "placement-group1",
				UncompressedUserData: &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							GroupName: aws.String("placement-group1"),
						},
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-1"),
								Groups:      []*string{aws.String("2"), aws.String("3")},
							},
						},
						TagSpecifications: []*ec2.TagSpecification{
							{
								ResourceType: aws.String("instance"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
										Value: aws.String("default/machine-aws-test1"),
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
								ResourceType: aws.String("network-interface"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
					})).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
									GroupName:        aws.String("placement-group1"),
								},
							},
						},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
		{
			name: "with dedicated tenancy and placement group ignition",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:         "m5.large",
				Tenancy:              "dedicated",
				PlacementGroupName:   "placement-group1",
				UncompressedUserData: &isUncompressedTrue,
				Ignition:             &infrav1.Ignition{},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							Tenancy:   &tenancy,
							GroupName: aws.String("placement-group1"),
						},
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-1"),
								Groups:      []*string{aws.String("2"), aws.String("3")},
							},
						},
						TagSpecifications: []*ec2.TagSpecification{
							{
								ResourceType: aws.String("instance"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
										Value: aws.String("default/machine-aws-test1"),
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
								ResourceType: aws.String("network-interface"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
						UserData: aws.String(base64.StdEncoding.EncodeToString(data)),
					})).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
									Tenancy:          &tenancy,
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
		{
			name: "with custom placement group and partition number",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:            "m5.large",
				PlacementGroupName:      "placement-group1",
				PlacementGroupPartition: 2,
				UncompressedUserData:    &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							GroupName:       aws.String("placement-group1"),
							PartitionNumber: aws.Int64(2),
						},
						NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
							{
								DeviceIndex: aws.Int64(0),
								SubnetId:    aws.String("subnet-1"),
								Groups:      aws.StringSlice([]string{"2", "3"}),
							},
						},
						TagSpecifications: []*ec2.TagSpecification{
							{
								ResourceType: aws.String("instance"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
										Value: aws.String("default/machine-aws-test1"),
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
								ResourceType: aws.String("network-interface"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("MachineName"),
										Value: aws.String("default/machine-aws-test1"),
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
					})).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
									GroupName:        aws.String("placement-group1"),
									PartitionNumber:  aws.Int64(2),
								},
							},
						},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
		{
			name: "expect error when placementGroupPartition is set, but placementGroupName is empty",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:            "m5.large",
				PlacementGroupPartition: 2,
				UncompressedUserData:    &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
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
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "placementGroupPartition is set but placementGroupName is empty"
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: `%s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "expect the default SSH key when none is provided",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName == nil {
							t.Fatal("Expected key name not to be nil")
						}
						if *input.KeyName != defaultSSHKeyName {
							t.Fatalf("Expected SSH key name to be '%s', not '%s'", defaultSSHKeyName, *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect to use the cluster level ssh key name when no machine key name is provided",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					SSHKeyName: aws.String("specific-cluster-key-name"),
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName == nil {
							t.Fatal("Expected key name not to be nil")
						}
						if *input.KeyName != "specific-cluster-key-name" {
							t.Fatalf("Expected SSH key name to be '%s', not '%s'", "specific-cluster-key-name", *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect to use the machine level ssh key name when both cluster and machine key names are provided",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
				SSHKeyName:   aws.String("specific-machine-ssh-key-name"),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					SSHKeyName: aws.String("specific-cluster-key-name"),
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName == nil {
							t.Fatal("Expected key name not to be nil")
						}
						if *input.KeyName != "specific-machine-ssh-key-name" {
							t.Fatalf("Expected SSH key name to be '%s', not '%s'", "specific-machine-ssh-key-name", *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect ssh key to be unset when cluster key name is empty string and machine key name is nil",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
				SSHKeyName:   nil,
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					SSHKeyName: aws.String(""),
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName != nil {
							t.Fatalf("Expected key name to be nil/unspecified, not '%s'", *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect ssh key to be unset when cluster key name is empty string and machine key name is empty string",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
				SSHKeyName:   aws.String(""),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					SSHKeyName: aws.String(""),
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName != nil {
							t.Fatalf("Expected key name to be nil/unspecified, not '%s'", *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect ssh key to be unset when cluster key name is nil and machine key name is empty string",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
					Version: ptr.To[string]("v1.16.1"),
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				InstanceType: "m5.large",
				SSHKeyName:   aws.String(""),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
						},
					},
					SSHKeyName: nil,
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
							aws.String("m5.large"),
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
				m.
					DescribeImagesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, input *ec2.RunInstancesInput, requestOptions ...request.Option) (*ec2.Reservation, error) {
						if input.KeyName != nil {
							t.Fatalf("Expected key name to be nil/unspecified, not '%s'", *input.KeyName)
						}
						return &ec2.Reservation{
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
									SubnetId:       aws.String("subnet-1"),
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
						}, nil
					})
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
		{
			name: "expect instace PrivateDNSName to be different when DHCP Option has domain name is set in the VPC",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-exists",
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
								NetworkInterfaces: []*ec2.InstanceNetworkInterface{
									{
										NetworkInterfaceId: aws.String("eni-1"),
										PrivateIpAddress:   aws.String("192.168.1.10"),
										PrivateDnsName:     aws.String("ip-192-168-1-10.ec2.internal"),
									},
								},
								VpcId: aws.String("vpc-exists"),
							},
						},
					}, nil)
				m.
					DescribeInstanceTypesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstanceTypesInput{
						InstanceTypes: []*string{
							aws.String("m5.large"),
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
				m.
					DescribeNetworkInterfacesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeNetworkInterfacesOutput{
						NetworkInterfaces: []*ec2.NetworkInterface{},
					}, nil)
				m.
					DescribeVpcs(&ec2.DescribeVpcsInput{
						VpcIds: []*string{aws.String("vpc-exists")},
					}).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							VpcId:         aws.String("vpc-exists"),
							CidrBlock:     aws.String("192.168.1.0/24"),
							IsDefault:     aws.Bool(false),
							State:         aws.String("available"),
							DhcpOptionsId: aws.String("dopt-12345678"),
						},
					},
				}, nil)
				m.
					DescribeDhcpOptions(&ec2.DescribeDhcpOptionsInput{
						DhcpOptionsIds: []*string{aws.String("dopt-12345678")},
					}).Return(&ec2.DescribeDhcpOptionsOutput{
					DhcpOptions: []*ec2.DhcpOptions{
						{
							DhcpConfigurations: []*ec2.DhcpConfiguration{
								{
									Key: aws.String("domain-name"),
									Values: []*ec2.AttributeValue{
										{
											Value: aws.String("example.com"),
										},
									},
								},
							},
						},
					},
				}, nil)
			},
			check: func(instance *infrav1.Instance, err error) {
				g := NewWithT(t)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(len(instance.Addresses)).To(Equal(3))

				for _, address := range instance.Addresses {
					if address.Type == clusterv1.MachineInternalIP {
						g.Expect(address.Address).To(Equal("192.168.1.10"))
					}

					if address.Type == clusterv1.MachineInternalDNS {
						g.Expect(address.Address).To(Or(Equal("ip-192-168-1-10.ec2.internal"), Equal("ip-192-168-1-10.example.com")))
					}
				}
			},
		},
		{
			name: "Simple, setting MarketType to MarketTypeCapacityBlock and providing CapacityReservationID",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:          "m5.large",
				MarketType:            infrav1.MarketTypeCapacityBlock,
				CapacityReservationID: aws.String("cr-12345678901234567"),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{

					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
								InstanceLifecycle:     aws.String(ec2.MarketTypeCapacityBlock),
								CapacityReservationId: aws.String("cr-12345678901234567"),
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
		{
			name: "expect error when MarketType to MarketTypeCapacityBlock set but not providing CapacityReservationID",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				MarketType:              infrav1.MarketTypeCapacityBlock,
				InstanceType:            "m5.large",
				PlacementGroupPartition: 2,
				UncompressedUserData:    &isUncompressedFalse,
			},
			awsCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
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
							aws.String("m5.large"),
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
			},
			check: func(instance *infrav1.Instance, err error) {
				expectedErrMsg := "capacityReservationID is required when CapacityBlock is enabled"
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Fatalf("Expected error: %s\nInstead got: `%s", expectedErrMsg, err.Error())
				}
			},
		},
		{
			name: "Simple, setting not setting MarketType and proving CapacityReservationID",
			machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:          "m5.large",
				CapacityReservationID: aws.String("cr-12345678901234567"),
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{

					NetworkSpec: infrav1.NetworkSpec{
						Subnets: infrav1.Subnets{
							infrav1.SubnetSpec{
								ID:       "subnet-1",
								IsPublic: false,
							},
							infrav1.SubnetSpec{
								IsPublic: false,
							},
						},
						VPC: infrav1.VPCSpec{
							ID: "vpc-test",
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
							aws.String("m5.large"),
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
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstancesWithContext(context.TODO(), gomock.Any()).
					Return(&ec2.Reservation{
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
								SubnetId:       aws.String("subnet-1"),
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
								CapacityReservationId: aws.String("cr-12345678901234567"),
								InstanceLifecycle:     aws.String("scheduled"),
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

func TestGetInstanceMarketOptionsRequest(t *testing.T) {
	mockCapacityReservationID := ptr.To[string]("cr-123")
	testCases := []struct {
		name            string
		instance        *infrav1.Instance
		expectedRequest *ec2.InstanceMarketOptionsRequest
		expectedError   error
	}{
		{
			name:            "with no Spot options specified",
			expectedRequest: nil,
			instance: &infrav1.Instance{
				SpotMarketOptions: nil,
			},
			expectedError: nil,
		},
		{
			name:            "with no MarketType",
			expectedRequest: nil,
			instance:        &infrav1.Instance{},
			expectedError:   nil,
		},
		{
			name:            "invalid MarketType specified",
			expectedRequest: nil,
			instance: &infrav1.Instance{
				MarketType: infrav1.MarketType("inValid"),
			},
			expectedError: errors.New("invalid MarketType \"inValid\""),
		},
		{
			name: "with an empty Spot options specified",
			instance: &infrav1.Instance{
				SpotMarketOptions: &infrav1.SpotMarketOptions{},
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
				},
			},
			expectedError: nil,
		},
		{
			name: "with marketType Spot specified",
			instance: &infrav1.Instance{
				MarketType: infrav1.MarketTypeSpot,
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
				},
			},
		},
		{
			name: "with marketType Spot and capacityRerservationID specified",
			instance: &infrav1.Instance{
				MarketType:            infrav1.MarketTypeSpot,
				CapacityReservationID: mockCapacityReservationID,
			},
			expectedError: errors.Errorf("unable to generate marketOptions for spot instance, capacityReservationID is incompatible with marketType spot and spotMarketOptions"),
		},
		{
			name: "with spotMarketOptions and capacityRerservationID specified",
			instance: &infrav1.Instance{
				SpotMarketOptions:     &infrav1.SpotMarketOptions{},
				CapacityReservationID: mockCapacityReservationID,
			},
			expectedError: errors.Errorf("unable to generate marketOptions for spot instance, capacityReservationID is incompatible with marketType spot and spotMarketOptions"),
		},
		{
			name: "with an empty MaxPrice specified",
			instance: &infrav1.Instance{
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String(""),
				},
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
				},
			},
			expectedError: nil,
		},
		{
			name: "with a valid MaxPrice specified",
			instance: &infrav1.Instance{
				SpotMarketOptions: &infrav1.SpotMarketOptions{
					MaxPrice: aws.String("0.01"),
				},
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
					MaxPrice:                     aws.String("0.01"),
				},
			},
			expectedError: nil,
		},
		{
			name:            "with no MarketTypeCapacityBlock options specified",
			instance:        &infrav1.Instance{},
			expectedRequest: nil,
			expectedError:   nil,
		},
		{
			name: "with a MarketType to MarketTypeCapacityBlock specified with capacityReservationID set to nil",
			instance: &infrav1.Instance{
				MarketType:            infrav1.MarketTypeCapacityBlock,
				CapacityReservationID: nil,
			},
			expectedRequest: nil,
			expectedError:   errors.Errorf("capacityReservationID is required when CapacityBlock is enabled"),
		},
		{
			name: "with a MarketType to MarketTypeCapacityBlock with capacityReservationID set to nil",
			instance: &infrav1.Instance{
				MarketType:            infrav1.MarketTypeCapacityBlock,
				CapacityReservationID: mockCapacityReservationID,
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeCapacityBlock),
			},
			expectedError: nil,
		},
		{
			name: "with a MarketType to MarketTypeCapacityBlock set with capacityReservationID set and empty Spot options specified",
			instance: &infrav1.Instance{
				MarketType:            infrav1.MarketTypeCapacityBlock,
				SpotMarketOptions:     &infrav1.SpotMarketOptions{},
				CapacityReservationID: mockCapacityReservationID,
			},
			expectedRequest: nil,
			expectedError:   errors.New("can't create spot capacity-blocks, remove spot market request"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := getInstanceMarketOptionsRequest(tc.instance)
			g := NewWithT(t)
			if tc.expectedError != nil {
				g.Expect(err.Error()).To(Equal(tc.expectedError.Error()))
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(request).To(Equal(tc.expectedRequest))
		})
	}
}

func TestGetFilteredSecurityGroupID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	securityGroupFilterName := "sg1"
	securityGroupFilterValues := []string{"test"}
	securityGroupID := "1"

	testCases := []struct {
		name          string
		securityGroup infrav1.AWSResourceReference
		expect        func(m *mocks.MockEC2APIMockRecorder)
		check         func(ids []string, err error)
	}{
		{
			name: "successfully return single security group id",
			securityGroup: infrav1.AWSResourceReference{
				Filters: []infrav1.Filter{
					{
						Name: securityGroupFilterName, Values: securityGroupFilterValues,
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String(securityGroupFilterName),
							Values: aws.StringSlice(securityGroupFilterValues),
						},
					},
				})).Return(
					&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{
								GroupId: aws.String(securityGroupID),
							},
						},
					}, nil)
			},
			check: func(ids []string, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if ids[0] != securityGroupID {
					t.Fatalf("expected security group id %v but got: %v", securityGroupID, ids[0])
				}
			},
		},
		{
			name: "allow returning multiple security groups",
			securityGroup: infrav1.AWSResourceReference{
				Filters: []infrav1.Filter{
					{
						Name: securityGroupFilterName, Values: securityGroupFilterValues,
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String(securityGroupFilterName),
							Values: aws.StringSlice(securityGroupFilterValues),
						},
					},
				})).Return(
					&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{
							{
								GroupId: aws.String(securityGroupID),
							},
							{
								GroupId: aws.String(securityGroupID),
							},
							{
								GroupId: aws.String(securityGroupID),
							},
						},
					}, nil)
			},
			check: func(ids []string, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				for _, id := range ids {
					if id != securityGroupID {
						t.Fatalf("expected security group id %v but got: %v", securityGroupID, id)
					}
				}
			},
		},
		{
			name:          "return early when filters are missing",
			securityGroup: infrav1.AWSResourceReference{},
			expect:        func(m *mocks.MockEC2APIMockRecorder) {},
			check: func(ids []string, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if len(ids) > 0 {
					t.Fatalf("didn't expect security group ids %v", ids)
				}
			},
		},
		{
			name: "error describing security group",
			securityGroup: infrav1.AWSResourceReference{
				Filters: []infrav1.Filter{
					{
						Name: securityGroupFilterName, Values: securityGroupFilterValues,
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String(securityGroupFilterName),
							Values: aws.StringSlice(securityGroupFilterValues),
						},
					},
				})).Return(nil, errors.New("some error"))
			},
			check: func(_ []string, err error) {
				if err == nil {
					t.Fatalf("expected error but got none.")
				}
			},
		},
		{
			name: "no error when no security groups found",
			securityGroup: infrav1.AWSResourceReference{
				Filters: []infrav1.Filter{
					{
						Name: securityGroupFilterName, Values: securityGroupFilterValues,
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSecurityGroupsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String(securityGroupFilterName),
							Values: aws.StringSlice(securityGroupFilterValues),
						},
					},
				})).Return(
					&ec2.DescribeSecurityGroupsOutput{
						SecurityGroups: []*ec2.SecurityGroup{},
					}, nil)
			},
			check: func(ids []string, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if len(ids) > 0 {
					t.Fatalf("didn't expect security group ids %v", ids)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())

			s := Service{
				EC2Client: ec2Mock,
			}

			ids, err := s.getFilteredSecurityGroupIDs(tc.securityGroup)
			tc.check(ids, err)
		})
	}
}

func TestGetDHCPOptionSetDomainName(t *testing.T) {
	testsCases := []struct {
		name                   string
		vpcID                  string
		dhcpOpt                *ec2.DhcpOptions
		expectedPrivateDNSName *string
		mockCalls              func(m *mocks.MockEC2APIMockRecorder)
	}{
		{
			name:  "dhcpOptions with domain-name",
			vpcID: "vpc-exists",
			dhcpOpt: &ec2.DhcpOptions{
				DhcpConfigurations: []*ec2.DhcpConfiguration{
					{
						Key: aws.String("domain-name"),
						Values: []*ec2.AttributeValue{
							{
								Value: aws.String("example.com"),
							},
						},
					},
				},
			},
			expectedPrivateDNSName: aws.String("example.com"),
			mockCalls:              mockedGetPrivateDNSDomainNameFromDHCPOptionsCalls,
		},
		{
			name:  "dhcpOptions without domain-name",
			vpcID: "vpc-empty-domain-name",
			dhcpOpt: &ec2.DhcpOptions{
				DhcpConfigurations: []*ec2.DhcpConfiguration{
					{
						Key:    aws.String("domain-name"),
						Values: []*ec2.AttributeValue{},
					},
				},
			},
			expectedPrivateDNSName: nil,
			mockCalls:              mockedGetPrivateDNSDomainNameFromDHCPOptionsEmptyCalls,
		},
	}
	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			scheme, err := setupScheme()
			g.Expect(err).ToNot(HaveOccurred())
			expect := func(m *mocks.MockEC2APIMockRecorder) {
				tc.mockCalls(m)
			}
			expect(ec2Mock.EXPECT())

			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := scope.NewClusterScope(
				scope.ClusterScopeParams{
					Client:  client,
					Cluster: &clusterv1.Cluster{},
					AWSCluster: &infrav1.AWSCluster{
						ObjectMeta: metav1.ObjectMeta{Name: "test"},
						Spec: infrav1.AWSClusterSpec{
							NetworkSpec: infrav1.NetworkSpec{
								VPC: infrav1.VPCSpec{
									ID: tc.vpcID,
								},
							},
						},
					},
				})
			g.Expect(err).ToNot(HaveOccurred())

			ec2Svc := NewService(cs)
			ec2Svc.EC2Client = ec2Mock
			dhcpOptsDomainName := ec2Svc.GetDHCPOptionSetDomainName(ec2Svc.EC2Client, &cs.VPC().ID)
			g.Expect(dhcpOptsDomainName).To(Equal(tc.expectedPrivateDNSName))
		})
	}
}

func mockedGetPrivateDNSDomainNameFromDHCPOptionsCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String("vpc-exists")},
	}).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{
				VpcId:         aws.String("vpc-exists"),
				CidrBlock:     aws.String("10.0.0.0/16"),
				IsDefault:     aws.Bool(false),
				State:         aws.String("available"),
				DhcpOptionsId: aws.String("dopt-12345678"),
			},
		},
	}, nil)
	m.DescribeDhcpOptions(&ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{aws.String("dopt-12345678")},
	}).Return(&ec2.DescribeDhcpOptionsOutput{
		DhcpOptions: []*ec2.DhcpOptions{
			{
				DhcpConfigurations: []*ec2.DhcpConfiguration{
					{
						Key: aws.String("domain-name"),
						Values: []*ec2.AttributeValue{
							{
								Value: aws.String("example.com"),
							},
						},
					},
				},
			},
		},
	}, nil)
}

func mockedGetPrivateDNSDomainNameFromDHCPOptionsEmptyCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String("vpc-empty-domain-name")},
	}).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{
				VpcId:         aws.String("vpc-exists"),
				CidrBlock:     aws.String("10.0.0.0/16"),
				IsDefault:     aws.Bool(false),
				State:         aws.String("available"),
				DhcpOptionsId: aws.String("dopt-empty"),
			},
		},
	}, nil)
	m.DescribeDhcpOptions(&ec2.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{aws.String("dopt-empty")},
	}).Return(&ec2.DescribeDhcpOptionsOutput{
		DhcpOptions: []*ec2.DhcpOptions{
			{
				DhcpConfigurations: []*ec2.DhcpConfiguration{
					{
						Key:    aws.String("domain-name"),
						Values: []*ec2.AttributeValue{},
					},
				},
			},
		},
	}, nil)
}

func TestGetCapacityReservationSpecification(t *testing.T) {
	mockCapacityReservationID := "cr-123"
	mockCapacityReservationIDPtr := &mockCapacityReservationID
	testCases := []struct {
		name                  string
		capacityReservationID *string
		expectedRequest       *ec2.CapacityReservationSpecification
	}{
		{
			name:                  "with no CapacityReservationID options specified",
			capacityReservationID: nil,
			expectedRequest:       nil,
		},
		{
			name:                  "with a valid CapacityReservationID specified",
			capacityReservationID: mockCapacityReservationIDPtr,
			expectedRequest: &ec2.CapacityReservationSpecification{
				CapacityReservationTarget: &ec2.CapacityReservationTarget{
					CapacityReservationId: aws.String(mockCapacityReservationID),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := getCapacityReservationSpecification(tc.capacityReservationID)
			if !cmp.Equal(request, tc.expectedRequest) {
				t.Errorf("Case: %s. Got: %v, expected: %v", tc.name, request, tc.expectedRequest)
			}
		})
	}
}
