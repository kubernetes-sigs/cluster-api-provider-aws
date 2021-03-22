package machine

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	awsproviderv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
	configv1.AddToScheme(scheme.Scheme)
}

func TestAvailabilityZone(t *testing.T) {
	cases := []struct {
		name             string
		availabilityZone string
		subnet           string
		expectedError    error
	}{
		{
			name:             "availability zone only",
			availabilityZone: "us-east-1a",
			expectedError:    errors.New("failed to launch instance: error getting subnet IDs: no subnet IDs were found"),
		},
		{
			name:   "subnet only",
			subnet: "subnet-b46032ec",
		},
		{
			name:             "availability zone and subnet",
			availabilityZone: "us-east-1a",
			subnet:           "subnet-b46032ec",
		},
	}

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			machine, err := stubMachine()
			if err != nil {
				t.Fatal(err)
			}

			machinePc, err := awsproviderv1.ProviderSpecFromRawExtension(machine.Spec.ProviderSpec.Value)
			if err != nil {
				t.Fatal(err)
			}

			// no load balancers tested
			machinePc.LoadBalancers = nil

			machinePc.Placement.AvailabilityZone = tc.availabilityZone
			if tc.subnet == "" {
				machinePc.Subnet.ID = nil
			} else {
				machinePc.Subnet.ID = aws.String(tc.subnet)
			}

			config, err := awsproviderv1.RawExtensionFromProviderSpec(machinePc)
			if err != nil {
				t.Fatal(err)
			}

			machine.Spec.ProviderSpec = machinev1.ProviderSpec{Value: config}

			fakeClient := fake.NewFakeClient(machine, awsCredentialsSecret, userDataSecret)

			err = fakeClient.Create(context.Background(), &configv1.Infrastructure{ObjectMeta: metav1.ObjectMeta{Name: awsclient.GlobalInfrastuctureName}})
			if err != nil {
				t.Fatal(err)
			}

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			machineScope, err := newMachineScope(machineScopeParams{
				client:  fakeClient,
				machine: machine,
				awsClientBuilder: func(client runtimeclient.Client, secretName, namespace, region string, configManagedClient runtimeclient.Client) (awsclient.Client, error) {
					return mockAWSClient, nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			reconciler := newReconciler(machineScope)

			var placement *ec2.Placement
			if tc.availabilityZone != "" && tc.subnet == "" {
				placement = &ec2.Placement{AvailabilityZone: aws.String(tc.availabilityZone)}
			}

			mockAWSClient.EXPECT().RunInstances(placementMatcher{placement}).Return(
				&ec2.Reservation{
					Instances: []*ec2.Instance{
						{
							ImageId:    aws.String("ami-a9acbbd6"),
							InstanceId: aws.String("i-02fcb933c5da7085c"),
							State: &ec2.InstanceState{
								Name: aws.String(ec2.InstanceStateNameRunning),
							},
							LaunchTime: aws.Time(time.Now()),
							Placement: &ec2.Placement{
								AvailabilityZone: aws.String("us-east-1a"),
							},
						},
					},
				}, nil)

			mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(
				&ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								{
									ImageId:    aws.String("ami-a9acbbd6"),
									InstanceId: aws.String("i-02fcb933c5da7085c"),
									State: &ec2.InstanceState{
										Name: aws.String(ec2.InstanceStateNameRunning),
										Code: aws.Int64(16),
									},
									LaunchTime: aws.Time(time.Now()),
								},
							},
						},
					},
				}, nil).AnyTimes()

			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil)
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).AnyTimes()
			mockAWSClient.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().DescribeSubnets(gomock.Any()).Return(&ec2.DescribeSubnetsOutput{}, nil)
			mockAWSClient.EXPECT().DescribeVpcs(gomock.Any()).Return(StubDescribeVPCs()).AnyTimes()
			mockAWSClient.EXPECT().DescribeDHCPOptions(gomock.Any()).Return(StubDescribeDHCPOptions()).AnyTimes()

			err = reconciler.create()
			if tc.expectedError != nil {
				if err == nil {
					t.Error("reconciler was expected to return error")
				}
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected: %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("reconciler was not expected to return error: %v", err)
				}
			}
		})
	}
}

type placementMatcher struct {
	placement *ec2.Placement
}

func (m placementMatcher) Matches(input interface{}) bool {
	runInstancesInput, ok := input.(*ec2.RunInstancesInput)
	if !ok {
		return false
	}
	if runInstancesInput.Placement == m.placement {
		return true
	}
	return false
}

func (m placementMatcher) String() string {
	return fmt.Sprintf("is placement: %#v", m.placement)
}

func TestCreate(t *testing.T) {
	// mock aws API calls
	mockCtrl := gomock.NewController(t)
	mockAWSClient := mockaws.NewMockClient(mockCtrl)
	mockAWSClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(nil, fmt.Errorf("describeSecurityGroups error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(nil, fmt.Errorf("describeAvailabilityZones error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeImages(gomock.Any()).Return(nil, fmt.Errorf("describeImages error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), nil).AnyTimes()
	mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil).AnyTimes()
	mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), nil).AnyTimes()
	mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), nil)
	mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
	mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAWSClient.EXPECT().DescribeVpcs(gomock.Any()).Return(StubDescribeVPCs()).AnyTimes()
	mockAWSClient.EXPECT().DescribeDHCPOptions(gomock.Any()).Return(StubDescribeDHCPOptions()).AnyTimes()

	testCases := []struct {
		testcase             string
		providerConfig       *awsproviderv1.AWSMachineProviderConfig
		userDataSecret       *corev1.Secret
		awsCredentialsSecret *corev1.Secret
		expectedError        error
	}{
		{
			testcase: "Create succeed",
			providerConfig: &awsproviderv1.AWSMachineProviderConfig{
				AMI: awsproviderv1.AWSResourceReference{
					ID: aws.String("ami-a9acbbd6"),
				},
				CredentialsSecret: &corev1.LocalObjectReference{
					Name: awsCredentialsSecretName,
				},
				InstanceType: "m4.xlarge",
				Placement: awsproviderv1.Placement{
					Region:           region,
					AvailabilityZone: defaultAvailabilityZone,
				},
				Subnet: awsproviderv1.AWSResourceReference{
					ID: aws.String("subnet-0e56b13a64ff8a941"),
				},
				IAMInstanceProfile: &awsproviderv1.AWSResourceReference{
					ID: aws.String("openshift_master_launch_instances"),
				},
				KeyName: aws.String(keyName),
				UserDataSecret: &corev1.LocalObjectReference{
					Name: userDataSecretName,
				},
				Tags: []awsproviderv1.TagSpecification{
					{Name: "openshift-node-group-config", Value: "node-config-master"},
					{Name: "host-type", Value: "master"},
					{Name: "sub-host-type", Value: "default"},
				},
				SecurityGroups: []awsproviderv1.AWSResourceReference{
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
				LoadBalancers: []awsproviderv1.LoadBalancerReference{
					{
						Name: "cluster-con",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-ext",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-int",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-net-lb",
						Type: awsproviderv1.NetworkLoadBalancerType,
					},
				},
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			expectedError:        nil,
		},
		{
			testcase: "Bad userData",
			providerConfig: &awsproviderv1.AWSMachineProviderConfig{
				AMI: awsproviderv1.AWSResourceReference{
					ID: aws.String("ami-a9acbbd6"),
				},
				CredentialsSecret: &corev1.LocalObjectReference{
					Name: awsCredentialsSecretName,
				},
				InstanceType: "m4.xlarge",
				Placement: awsproviderv1.Placement{
					Region:           region,
					AvailabilityZone: defaultAvailabilityZone,
				},
				Subnet: awsproviderv1.AWSResourceReference{
					ID: aws.String("subnet-0e56b13a64ff8a941"),
				},
				IAMInstanceProfile: &awsproviderv1.AWSResourceReference{
					ID: aws.String("openshift_master_launch_instances"),
				},
				KeyName: aws.String(keyName),
				UserDataSecret: &corev1.LocalObjectReference{
					Name: userDataSecretName,
				},
				Tags: []awsproviderv1.TagSpecification{
					{Name: "openshift-node-group-config", Value: "node-config-master"},
					{Name: "host-type", Value: "master"},
					{Name: "sub-host-type", Value: "default"},
				},
				SecurityGroups: []awsproviderv1.AWSResourceReference{
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
				LoadBalancers: []awsproviderv1.LoadBalancerReference{
					{
						Name: "cluster-con",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-ext",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-int",
						Type: awsproviderv1.ClassicLoadBalancerType,
					},
					{
						Name: "cluster-net-lb",
						Type: awsproviderv1.NetworkLoadBalancerType,
					},
				},
			},
			userDataSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespace,
				},
				Data: map[string][]byte{
					"badKey": []byte(userDataBlob),
				},
			},
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			expectedError:        errors.New("failed to get user data: secret default/aws-actuator-user-data-secret missing userData key"),
		},
		{
			testcase: "Failed security groups return invalid config machine error",
			providerConfig: &awsproviderv1.AWSMachineProviderConfig{
				AMI: awsproviderv1.AWSResourceReference{
					ID: aws.String("ami-a9acbbd6"),
				},
				CredentialsSecret: &corev1.LocalObjectReference{
					Name: awsCredentialsSecretName,
				},
				InstanceType: "m4.xlarge",
				Placement: awsproviderv1.Placement{
					Region:           region,
					AvailabilityZone: defaultAvailabilityZone,
				},
				Subnet: awsproviderv1.AWSResourceReference{
					ID: aws.String("subnet-0e56b13a64ff8a941"),
				},
				IAMInstanceProfile: &awsproviderv1.AWSResourceReference{
					ID: aws.String("openshift_master_launch_instances"),
				},
				KeyName: aws.String(keyName),
				UserDataSecret: &corev1.LocalObjectReference{
					Name: userDataSecretName,
				},
				Tags: []awsproviderv1.TagSpecification{
					{Name: "openshift-node-group-config", Value: "node-config-master"},
					{Name: "host-type", Value: "master"},
					{Name: "sub-host-type", Value: "default"},
				},
				SecurityGroups: []awsproviderv1.AWSResourceReference{{
					Filters: []awsproviderv1.Filter{{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-%s-sg", stubClusterID, "role")},
					}},
				}},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			expectedError:        errors.New("failed to launch instance: error getting security groups IDs: error describing security groups: describeSecurityGroups error"),
		},
		{
			testcase: "Failed Availability zones return invalid config machine error",
			providerConfig: &awsproviderv1.AWSMachineProviderConfig{
				AMI: awsproviderv1.AWSResourceReference{
					ID: aws.String("ami-a9acbbd6"),
				},
				CredentialsSecret: &corev1.LocalObjectReference{
					Name: awsCredentialsSecretName,
				},
				InstanceType: "m4.xlarge",
				Placement: awsproviderv1.Placement{
					Region:           region,
					AvailabilityZone: defaultAvailabilityZone,
				},
				Subnet: awsproviderv1.AWSResourceReference{
					Filters: []awsproviderv1.Filter{{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-private-%s", stubClusterID, "az")},
					}},
				},
				IAMInstanceProfile: &awsproviderv1.AWSResourceReference{
					ID: aws.String("openshift_master_launch_instances"),
				},
				KeyName: aws.String(keyName),
				UserDataSecret: &corev1.LocalObjectReference{
					Name: userDataSecretName,
				},
				Tags: []awsproviderv1.TagSpecification{
					{Name: "openshift-node-group-config", Value: "node-config-master"},
					{Name: "host-type", Value: "master"},
					{Name: "sub-host-type", Value: "default"},
				},
				SecurityGroups: []awsproviderv1.AWSResourceReference{
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			expectedError:        errors.New("failed to launch instance: error getting subnet IDs: error describing availability zones: describeAvailabilityZones error"),
		},
		{
			testcase: "Failed BlockDevices return invalid config machine error",
			providerConfig: &awsproviderv1.AWSMachineProviderConfig{
				AMI: awsproviderv1.AWSResourceReference{
					ID: aws.String("ami-a9acbbd6"),
				},
				CredentialsSecret: &corev1.LocalObjectReference{
					Name: awsCredentialsSecretName,
				},
				InstanceType: "m4.xlarge",
				Placement: awsproviderv1.Placement{
					Region:           region,
					AvailabilityZone: defaultAvailabilityZone,
				},
				BlockDevices: []awsproviderv1.BlockDeviceMappingSpec{
					{
						EBS: &awsproviderv1.EBSBlockDeviceSpec{
							VolumeType: pointer.StringPtr("type"),
							VolumeSize: pointer.Int64Ptr(int64(1)),
							Iops:       pointer.Int64Ptr(int64(1)),
						},
					},
				},
				Subnet: awsproviderv1.AWSResourceReference{
					ID: aws.String("subnet-0e56b13a64ff8a941"),
				},
				IAMInstanceProfile: &awsproviderv1.AWSResourceReference{
					ID: aws.String("openshift_master_launch_instances"),
				},
				KeyName: aws.String(keyName),
				UserDataSecret: &corev1.LocalObjectReference{
					Name: userDataSecretName,
				},
				Tags: []awsproviderv1.TagSpecification{
					{Name: "openshift-node-group-config", Value: "node-config-master"},
					{Name: "host-type", Value: "master"},
					{Name: "sub-host-type", Value: "default"},
				},
				SecurityGroups: []awsproviderv1.AWSResourceReference{
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			expectedError:        errors.New("failed to launch instance: error getting blockDeviceMappings: error describing AMI: describeImages error"),
		},
	}
	for _, tc := range testCases {
		// create fake resources
		t.Logf("testCase: %v", tc.testcase)

		encodedProviderConfig, err := awsproviderv1.RawExtensionFromProviderSpec(tc.providerConfig)
		if err != nil {
			t.Fatalf("Unexpected error")
		}
		machine, err := stubMachine()
		if err != nil {
			t.Fatal(err)
		}
		machine.Spec.ProviderSpec = machinev1.ProviderSpec{Value: encodedProviderConfig}

		fakeClient := fake.NewFakeClientWithScheme(scheme.Scheme, machine, tc.awsCredentialsSecret, tc.userDataSecret)

		err = fakeClient.Create(context.Background(), &configv1.Infrastructure{ObjectMeta: metav1.ObjectMeta{Name: awsclient.GlobalInfrastuctureName}})
		if err != nil {
			t.Fatal(err)
		}

		machineScope, err := newMachineScope(machineScopeParams{
			client:  fakeClient,
			machine: machine,
			awsClientBuilder: func(client runtimeclient.Client, secretName, namespace, region string, configManagedClient runtimeclient.Client) (awsclient.Client, error) {
				return mockAWSClient, nil
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		reconciler := newReconciler(machineScope)

		// test create
		err = reconciler.create()
		if tc.expectedError != nil {
			if err == nil {
				t.Error("reconciler was expected to return error")
			}
			if err.Error() != tc.expectedError.Error() {
				t.Errorf("Expected: %v, got %v", tc.expectedError, err)
			}
		} else {
			if err != nil {
				t.Errorf("reconciler was not expected to return error: %v", err)
			}
		}
	}
}

func TestGetMachineInstances(t *testing.T) {
	clusterID := "aws-actuator-cluster"
	instanceID := "i-02fa4197109214b46"
	imageID := "ami-a9acbbd6"

	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("unable to build stub machine: %v", err)
	}

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	testCases := []struct {
		testcase       string
		providerStatus awsproviderv1.AWSMachineProviderStatus
		awsClientFunc  func(*gomock.Controller) awsclient.Client
		exists         bool
	}{
		{
			testcase:       "empty-status-search-by-tag",
			providerStatus: awsproviderv1.AWSMachineProviderStatus{},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				request := &ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						{
							Name:   awsTagFilter("Name"),
							Values: aws.StringSlice([]string{machine.Name}),
						},

						clusterFilter(clusterID),
					},
				}

				mockAWSClient.EXPECT().DescribeInstances(request).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameRunning),
					nil,
				).Times(1)

				return mockAWSClient
			},
			exists: true,
		},
		{
			testcase: "has-status-search-by-id-running",
			providerStatus: awsproviderv1.AWSMachineProviderStatus{
				InstanceID: aws.String(instanceID),
			},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				request := &ec2.DescribeInstancesInput{
					InstanceIds: aws.StringSlice([]string{instanceID}),
				}

				mockAWSClient.EXPECT().DescribeInstances(request).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameRunning),
					nil,
				).Times(1)

				return mockAWSClient
			},
			exists: true,
		},
		{
			testcase: "has-status-search-by-id-terminated",
			providerStatus: awsproviderv1.AWSMachineProviderStatus{
				InstanceID: aws.String(instanceID),
			},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				first := mockAWSClient.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
					InstanceIds: aws.StringSlice([]string{instanceID}),
				}).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameTerminated),
					nil,
				).Times(1)

				mockAWSClient.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						{
							Name:   awsTagFilter("Name"),
							Values: aws.StringSlice([]string{machine.Name}),
						},

						clusterFilter(clusterID),
					},
				}).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameTerminated),
					nil,
				).Times(1).After(first)

				return mockAWSClient
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			awsStatusRaw, err := awsproviderv1.RawExtensionFromProviderStatus(&tc.providerStatus)
			if err != nil {
				t.Fatal(err)
			}

			machineCopy := machine.DeepCopy()
			machineCopy.Status.ProviderStatus = awsStatusRaw

			fakeClient := fake.NewFakeClientWithScheme(scheme.Scheme, machine, awsCredentialsSecret, userDataSecret)
			mockAWSClient := tc.awsClientFunc(ctrl)

			machineScope, err := newMachineScope(machineScopeParams{
				client:  fakeClient,
				machine: machineCopy,
				awsClientBuilder: func(client runtimeclient.Client, secretName, namespace, region string, configManagedClient runtimeclient.Client) (awsclient.Client, error) {
					return mockAWSClient, nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			reconciler := newReconciler(machineScope)

			instances, err := reconciler.getMachineInstances()
			if err != nil {
				t.Errorf("Unexpected error from getMachineInstances: %v", err)
			}
			if tc.exists != (len(instances) > 0) {
				t.Errorf("Expected instance exists: %t, got instances: %v", tc.exists, instances)
			}
		})
	}
}
