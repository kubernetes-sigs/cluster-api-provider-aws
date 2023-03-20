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
	"encoding/base64"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
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
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
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
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
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
				m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
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
				m.DescribeInstances(&ec2.DescribeInstancesInput{
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

			instance, err := s.InstanceIfExists(&tc.instanceID)
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
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
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
				m.TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
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
		machine       clusterv1.Machine
		machineConfig *infrav1.AWSMachineSpec
		awsCluster    *infrav1.AWSCluster
		expect        func(m *mocks.MockEC2APIMockRecorder)
		check         func(instance *infrav1.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			name: "with ImageLookupOrg specified at the machine level",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-18.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-18.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
				amiName, err := GenerateAmiName("capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*", "ubuntu-18.04", "1.16.1")
				if err != nil {
					t.Fatalf("Failed to process ami format: %v", err)
				}
				m.
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Eq(&ec2.DescribeImagesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					RunInstances(gomock.Any()).
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
							{Name: aws.String("subnet-id"), Values: aws.StringSlice([]string{"matching-subnet"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId: aws.String("matching-subnet"),
						}},
					}, nil)
				m.
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeSubnets(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
							filter.EC2.VPC("vpc-id"),
							{Name: aws.String("tag:some-tag"), Values: aws.StringSlice([]string{"some-value"})},
						},
					}).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{{
							SubnetId:            aws.String("filtered-subnet-1"),
							MapPublicIpOnLaunch: aws.Bool(true),
						}},
					}, nil)
				m.
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Any()).
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
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
					RunInstances(gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							Tenancy: &tenancy,
						},
						SecurityGroupIds: []*string{aws.String("2"), aws.String("3")},
						SubnetId:         aws.String("subnet-1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			name: "with dedicated tenancy ignition",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"set": "node"},
					Namespace: "default",
					Name:      "machine-aws-test1",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
				},
			},
			machineConfig: &infrav1.AWSMachineSpec{
				AMI: infrav1.AMIReference{
					ID: aws.String("abc"),
				},
				InstanceType:         "m5.large",
				Tenancy:              "dedicated",
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					RunInstances(gomock.Eq(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("m5.large"),
						KeyName:      aws.String("default"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						Placement: &ec2.Placement{
							Tenancy: &tenancy,
						},
						SecurityGroupIds: []*string{aws.String("2"), aws.String("3")},
						SubnetId:         aws.String("subnet-1"),
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			name: "expect the default SSH key when none is provided",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-data"),
					},
					Version: pointer.String("v1.16.1"),
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
					DescribeInstanceTypes(gomock.Eq(&ec2.DescribeInstanceTypesInput{
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
					DescribeImages(gomock.Any()).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								Name:         aws.String("ami-1"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
				m. // TODO: Restore these parameters, but with the tags as well
					RunInstances(gomock.Any()).
					DoAndReturn(func(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
					DescribeNetworkInterfaces(gomock.Any()).
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

			machine := &tc.machine

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
	testCases := []struct {
		name              string
		spotMarketOptions *infrav1.SpotMarketOptions
		expectedRequest   *ec2.InstanceMarketOptionsRequest
	}{
		{
			name:              "with no Spot options specified",
			spotMarketOptions: nil,
			expectedRequest:   nil,
		},
		{
			name:              "with an empty Spot options specified",
			spotMarketOptions: &infrav1.SpotMarketOptions{},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
				},
			},
		},
		{
			name: "with an empty MaxPrice specified",
			spotMarketOptions: &infrav1.SpotMarketOptions{
				MaxPrice: aws.String(""),
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
			name: "with a valid MaxPrice specified",
			spotMarketOptions: &infrav1.SpotMarketOptions{
				MaxPrice: aws.String("0.01"),
			},
			expectedRequest: &ec2.InstanceMarketOptionsRequest{
				MarketType: aws.String(ec2.MarketTypeSpot),
				SpotOptions: &ec2.SpotMarketOptions{
					InstanceInterruptionBehavior: aws.String(ec2.InstanceInterruptionBehaviorTerminate),
					SpotInstanceType:             aws.String(ec2.SpotInstanceTypeOneTime),
					MaxPrice:                     aws.String("0.01"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := getInstanceMarketOptionsRequest(tc.spotMarketOptions)
			if !cmp.Equal(request, tc.expectedRequest) {
				t.Errorf("Case: %s. Got: %v, expected: %v", tc.name, request, tc.expectedRequest)
			}
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
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{
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
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{
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
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{
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
				m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{
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
