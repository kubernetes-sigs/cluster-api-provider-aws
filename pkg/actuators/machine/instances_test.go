package machine

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRemoveDuplicatedTags(t *testing.T) {
	cases := []struct {
		tagList  []*ec2.Tag
		expected []*ec2.Tag
	}{
		{
			// empty tags
			tagList:  []*ec2.Tag{},
			expected: []*ec2.Tag{},
		},
		{
			// no duplicate tags
			tagList: []*ec2.Tag{
				{Key: aws.String("clusterID"), Value: aws.String("test-ClusterIDValue")},
			},
			expected: []*ec2.Tag{
				{Key: aws.String("clusterID"), Value: aws.String("test-ClusterIDValue")},
			},
		},
		{
			// multiple duplicate tags
			tagList: []*ec2.Tag{
				{Key: aws.String("clusterID"), Value: aws.String("test-ClusterIDValue")},
				{Key: aws.String("clusterSize"), Value: aws.String("test-ClusterSizeValue")},
				{Key: aws.String("clusterSize"), Value: aws.String("test-ClusterSizeDuplicatedValue")},
			},
			expected: []*ec2.Tag{
				{Key: aws.String("clusterID"), Value: aws.String("test-ClusterIDValue")},
				{Key: aws.String("clusterSize"), Value: aws.String("test-ClusterSizeValue")},
			},
		},
	}

	for i, c := range cases {
		actual := removeDuplicatedTags(c.tagList)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Errorf("test #%d: expected %+v, got %+v", i, c.expected, actual)
		}
	}
}

func TestBuildTagList(t *testing.T) {
	cases := []struct {
		name            string
		machineSpecTags []machinev1.TagSpecification
		infra           *configv1.Infrastructure
		expected        []*ec2.Tag
	}{
		{
			name:            "with empty infra and provider spec should return default tags",
			machineSpecTags: []machinev1.TagSpecification{},
			infra: &configv1.Infrastructure{
				Status: configv1.InfrastructureStatus{
					PlatformStatus: &configv1.PlatformStatus{
						AWS: &configv1.AWSPlatformStatus{
							ResourceTags: []configv1.AWSResourceTag{},
						},
					},
				},
			},
			expected: []*ec2.Tag{
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
		{
			name:            "with empty infra should return default tags",
			machineSpecTags: []machinev1.TagSpecification{},
			infra:           &configv1.Infrastructure{}, // should work with empty infra object
			expected: []*ec2.Tag{
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
		{
			name:            "with nil infra should  return default tags",
			machineSpecTags: []machinev1.TagSpecification{},
			infra:           nil, // should work with nil infra object
			expected: []*ec2.Tag{
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
		{
			name: "should filter out bad tags from provider spec",
			machineSpecTags: []machinev1.TagSpecification{
				{Name: "Name", Value: "badname"},
				{Name: "kubernetes.io/cluster/badid", Value: "badvalue"},
				{Name: "good", Value: "goodvalue"},
			},
			infra: nil,
			// Invalid tags get dropped and the valid clusterID and Name get applied last.
			expected: []*ec2.Tag{
				{Key: aws.String("good"), Value: aws.String("goodvalue")},
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
		{
			name:            "should filter out bad tags from infra object",
			machineSpecTags: []machinev1.TagSpecification{},
			infra: &configv1.Infrastructure{
				Status: configv1.InfrastructureStatus{
					PlatformStatus: &configv1.PlatformStatus{
						AWS: &configv1.AWSPlatformStatus{
							ResourceTags: []configv1.AWSResourceTag{
								{
									Key:   "kubernetes.io/cluster/badid",
									Value: "badvalue",
								},
								{
									Key:   "Name",
									Value: "badname",
								},
								{
									Key:   "good",
									Value: "goodvalue",
								},
							},
						},
					},
				},
			},
			// Invalid tags get dropped and the valid clusterID and Name get applied last.
			expected: []*ec2.Tag{
				{Key: aws.String("good"), Value: aws.String("goodvalue")},
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
		{
			name: "tags from machine object should have precedence",
			machineSpecTags: []machinev1.TagSpecification{
				{Name: "Name", Value: "badname"},
				{Name: "kubernetes.io/cluster/badid", Value: "badvalue"},
				{Name: "good", Value: "goodvalue"},
			},
			infra: &configv1.Infrastructure{
				Status: configv1.InfrastructureStatus{
					PlatformStatus: &configv1.PlatformStatus{
						AWS: &configv1.AWSPlatformStatus{
							ResourceTags: []configv1.AWSResourceTag{
								{
									Key:   "good",
									Value: "should-be-overwritten",
								},
							},
						},
					},
				},
			},
			expected: []*ec2.Tag{
				{Key: aws.String("good"), Value: aws.String("goodvalue")},
				{Key: aws.String("kubernetes.io/cluster/clusterID"), Value: aws.String("owned")},
				{Key: aws.String("Name"), Value: aws.String("machineName")},
			},
		},
	}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := buildTagList("machineName", "clusterID", c.machineSpecTags, c.infra)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("test #%d: expected %+v, got %+v", i, c.expected, actual)
			}
		})
	}
}

func TestBuildEC2Filters(t *testing.T) {
	filter1 := "filter1"
	filter2 := "filter2"
	value1 := "A"
	value2 := "B"
	value3 := "C"

	inputFilters := []machinev1.Filter{
		{
			Name:   filter1,
			Values: []string{value1, value2},
		},
		{
			Name:   filter2,
			Values: []string{value3},
		},
	}

	expected := []*ec2.Filter{
		{
			Name:   &filter1,
			Values: []*string{&value1, &value2},
		},
		{
			Name:   &filter2,
			Values: []*string{&value3},
		},
	}

	got := buildEC2Filters(inputFilters)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("failed to buildEC2Filters. Expected: %+v, got: %+v", expected, got)
	}
}

func TestGetBlockDeviceMappings(t *testing.T) {
	rootDeviceName := "/dev/sda1"
	volumeSize := int64(16384)
	deviceName2 := "/dev/sda2"
	volumeSize2 := int64(16385)
	deleteOnTermination := true
	volumeType := "ssd"

	mockCtrl := gomock.NewController(t)
	mockAWSClient := mockaws.NewMockClient(mockCtrl)
	mockAWSClient.EXPECT().DescribeImages(gomock.Any()).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{
			{
				CreationDate:   aws.String(time.RFC3339),
				ImageId:        aws.String("ami-1111"),
				RootDeviceName: &rootDeviceName,
			},
		},
	}, nil).AnyTimes()

	oneBlockDevice := []machinev1.BlockDeviceMappingSpec{
		{
			DeviceName: &rootDeviceName,
			EBS: &machinev1.EBSBlockDeviceSpec{
				VolumeSize: &volumeSize,
				VolumeType: &volumeType,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
	}

	oneExpectedBlockDevice := []*ec2.BlockDeviceMapping{
		{
			DeviceName: &rootDeviceName,
			Ebs: &ec2.EbsBlockDevice{
				VolumeSize:          &volumeSize,
				VolumeType:          &volumeType,
				DeleteOnTermination: &deleteOnTermination,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
	}

	blockDevices := []machinev1.BlockDeviceMappingSpec{
		{
			DeviceName: &rootDeviceName,
			EBS: &machinev1.EBSBlockDeviceSpec{
				VolumeSize: &volumeSize,
				VolumeType: &volumeType,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
		{
			DeviceName: &deviceName2,
			EBS: &machinev1.EBSBlockDeviceSpec{
				VolumeSize: &volumeSize2,
				VolumeType: &volumeType,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
	}

	twoExpectedDevices := []*ec2.BlockDeviceMapping{
		{
			DeviceName: &rootDeviceName,
			Ebs: &ec2.EbsBlockDevice{
				VolumeSize:          &volumeSize,
				VolumeType:          &volumeType,
				DeleteOnTermination: &deleteOnTermination,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
		{
			DeviceName: &deviceName2,
			Ebs: &ec2.EbsBlockDevice{
				VolumeSize:          &volumeSize2,
				VolumeType:          &volumeType,
				DeleteOnTermination: &deleteOnTermination,
			},
			NoDevice:    nil,
			VirtualName: nil,
		},
	}

	blockDevicesOneEmptyName := make([]machinev1.BlockDeviceMappingSpec, len(blockDevices))
	copy(blockDevicesOneEmptyName, blockDevices)
	blockDevicesOneEmptyName[0].DeviceName = nil

	blockDevicesTwoEmptyNames := make([]machinev1.BlockDeviceMappingSpec, len(blockDevicesOneEmptyName))
	copy(blockDevicesTwoEmptyNames, blockDevicesOneEmptyName)
	blockDevicesTwoEmptyNames[1].DeviceName = nil

	testCases := []struct {
		description  string
		blockDevices []machinev1.BlockDeviceMappingSpec
		expected     []*ec2.BlockDeviceMapping
		expectedErr  bool
	}{
		{
			description:  "When it gets an empty blockDevices list",
			blockDevices: []machinev1.BlockDeviceMappingSpec{},
			expected:     []*ec2.BlockDeviceMapping{},
		},
		{
			description:  "When it gets one blockDevice",
			blockDevices: oneBlockDevice,
			expected:     oneExpectedBlockDevice,
		},
		{
			description:  "When it gets two blockDevices",
			blockDevices: blockDevices,
			expected:     twoExpectedDevices,
		},
		{
			description:  "When it gets two blockDevices and one with empty device name",
			blockDevices: blockDevicesOneEmptyName,
			expected:     twoExpectedDevices,
		},
		{
			description:  "Fail when it gets two blockDevices and two with empty device name",
			blockDevices: blockDevicesTwoEmptyNames,
			expectedErr:  true,
		},
	}

	fakeMachineKey := client.ObjectKey{
		Name:      "fake",
		Namespace: "fake",
	}
	for _, tc := range testCases {
		got, err := getBlockDeviceMappings(fakeMachineKey, tc.blockDevices, "existing-AMI", mockAWSClient)
		if tc.expectedErr {
			if err == nil {
				t.Error("Expected error")
			}
		} else {
			if err != nil {
				t.Errorf("error when calling getBlockDeviceMappings: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Got: %v, expected: %v", got, tc.expected)
			}
		}
	}
}

func TestRemoveStoppedMachine(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}

	cases := []struct {
		name   string
		output *ec2.DescribeInstancesOutput
		err    error
	}{
		{
			name:   "DescribeInstances with error",
			output: &ec2.DescribeInstancesOutput{},
			// any non-nil error will do
			err: fmt.Errorf("error describing instances"),
		},
		{
			name: "No instances to stop",
			output: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
		},
		{
			name: "One instance to stop",
			output: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							stubInstance(stubAMIID, stubInstanceID, true),
						},
					},
				},
			},
		},
		{
			name: "Two instances to stop",
			output: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							stubInstance(stubAMIID, stubInstanceID, true),
							stubInstance("ami-a9acbbd7", "i-02fcb933c5da7085d", true),
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			// Not here to check how many times all the mocked methods get called.
			// Rather to provide fake outputs to get through all possible execution paths.
			mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(tc.output, tc.err).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).AnyTimes()
			removeStoppedMachine(machine, mockAWSClient)
		})
	}
}

func TestLaunchInstance(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}

	providerConfig := stubProviderConfig()
	stubTagList := buildTagList(machine.Name, stubClusterID, providerConfig.Tags, nil)

	infra := &configv1.Infrastructure{
		Status: configv1.InfrastructureStatus{
			PlatformStatus: &configv1.PlatformStatus{
				AWS: &configv1.AWSPlatformStatus{
					ResourceTags: []configv1.AWSResourceTag{
						{
							Key:   "infra-tag-key",
							Value: "infra-tag-value",
						},
					},
				},
			},
		},
	}

	stubTagListWithInfraObject := buildTagList(machine.Name, stubClusterID, providerConfig.Tags, infra)

	cases := []struct {
		name                string
		providerConfig      *machinev1.AWSMachineProviderConfig
		securityGroupOutput *ec2.DescribeSecurityGroupsOutput
		securityGroupErr    error
		subnetOutput        *ec2.DescribeSubnetsOutput
		subnetErr           error
		azErr               error
		imageOutput         *ec2.DescribeImagesOutput
		imageErr            error
		instancesOutput     *ec2.Reservation
		instancesErr        error
		succeeds            bool
		runInstancesInput   *ec2.RunInstancesInput
		infra               *configv1.Infrastructure
	}{
		{
			name: "Security groups with filters",
			providerConfig: stubPCSecurityGroups(
				[]machinev1.AWSResourceReference{
					{
						Filters: []machinev1.Filter{},
					},
				},
			),
			securityGroupOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{
					{
						GroupId: aws.String("groupID"),
					},
				},
			},
			instancesOutput: stubReservation(stubAMIID, stubInstanceID, "192.168.0.10"),
			succeeds:        true,
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
						Groups:                   []*string{aws.String("groupID")},
					},
				},
				UserData: aws.String(""),
			},
		},
		{
			name: "Security groups with filters with error",
			providerConfig: stubPCSecurityGroups(
				[]machinev1.AWSResourceReference{
					{
						Filters: []machinev1.Filter{},
					},
				},
			),
			securityGroupErr: fmt.Errorf("error"),
		},
		{
			name: "No security group",
			providerConfig: stubPCSecurityGroups(
				[]machinev1.AWSResourceReference{
					{
						Filters: []machinev1.Filter{},
					},
				},
			),
			securityGroupOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{},
			},
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
					},
				},
				UserData: aws.String(""),
			},
		},
		{
			name: "Subnet with filters",
			providerConfig: stubPCSubnet(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{},
			}),
			subnetOutput: &ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId: aws.String("subnetID"),
					},
				},
			},
			instancesOutput: stubReservation(stubAMIID, stubInstanceID, "192.168.0.10"),
			succeeds:        true,
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 aws.String("subnetID"),
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
				Placement: &ec2.Placement{
					AvailabilityZone: aws.String("us-east-1a"),
				},
			},
		},
		{
			name: "Subnet with filters with error",
			providerConfig: stubPCSubnet(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{},
			}),
			subnetErr: fmt.Errorf("error"),
		},
		{
			name: "Subnet with availability zone with error",
			providerConfig: stubPCSubnet(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{},
			}),
			azErr: fmt.Errorf("error"),
		},
		{
			name: "AMI with filters",
			providerConfig: stubPCAMI(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{
					{
						Name:   "foo",
						Values: []string{"bar"},
					},
				},
			}),
			imageOutput: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{
					{
						CreationDate: aws.String("2006-01-02T15:04:05Z"),
						ImageId:      aws.String("ami-1111"),
					},
				},
			},
			instancesOutput: stubReservation(stubAMIID, stubInstanceID, "192.168.0.10"),
			succeeds:        true,
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String("ami-1111"),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
			},
		},
		{
			name: "AMI with filters with error",
			providerConfig: stubPCAMI(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{},
			}),
			imageErr: fmt.Errorf("error"),
		},
		{
			name: "AMI with filters with no image",
			providerConfig: stubPCAMI(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{
					{
						Name:   "image_stage",
						Values: []string{"base"},
					},
				},
			}),
			imageOutput: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{},
			},
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 aws.String("subnetID"),
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
				Placement: &ec2.Placement{
					AvailabilityZone: aws.String("us-east-1a"),
				},
			},
		},
		{
			name: "AMI with filters with two images",
			providerConfig: stubPCAMI(machinev1.AWSResourceReference{
				Filters: []machinev1.Filter{
					{
						Name:   "image_stage",
						Values: []string{"base"},
					},
				},
			}),
			imageOutput: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{
					{
						CreationDate: aws.String("2006-01-02T15:04:05Z"),
						ImageId:      aws.String("ami-1111"),
					},
					{
						CreationDate: aws.String("2006-01-02T15:04:05Z"),
						ImageId:      aws.String("ami-2222"),
					},
				},
			},
			instancesOutput: stubReservation(stubAMIID, stubInstanceID, "192.168.0.10"),
			succeeds:        true,
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String("ami-1111"),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
			},
		},
		{
			name:           "AMI not specified",
			providerConfig: stubPCAMI(machinev1.AWSResourceReference{}),
		},
		{
			name:           "Dedicated instance tenancy",
			providerConfig: stubDedicatedInstanceTenancy(),
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagList,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagList,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
				Placement: &ec2.Placement{
					Tenancy: aws.String("dedicated"),
				},
			},
		},
		{
			name:           "Dedicated instance tenancy",
			providerConfig: stubInvalidInstanceTenancy(),
		},
		{
			name:           "Attach infrastructure object tags",
			providerConfig: providerConfig,
			infra:          infra,
			runInstancesInput: &ec2.RunInstancesInput{
				IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
					Name: aws.String(*providerConfig.IAMInstanceProfile.ID),
				},
				ImageId:      aws.String(*providerConfig.AMI.ID),
				InstanceType: &providerConfig.InstanceType,
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				KeyName:      providerConfig.KeyName,
				TagSpecifications: []*ec2.TagSpecification{{
					ResourceType: aws.String("instance"),
					Tags:         stubTagListWithInfraObject,
				}, {
					ResourceType: aws.String("volume"),
					Tags:         stubTagListWithInfraObject,
				}},
				NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
					{
						DeviceIndex:              aws.Int64(providerConfig.DeviceIndex),
						AssociatePublicIpAddress: providerConfig.PublicIP,
						SubnetId:                 providerConfig.Subnet.ID,
						Groups: []*string{
							aws.String("sg-00868b02fbe29de17"),
							aws.String("sg-0a4658991dc5eb40a"),
							aws.String("sg-009a70e28fa4ba84e"),
							aws.String("sg-07323d56fb932c84c"),
							aws.String("sg-08b1ffd32874d59a2"),
						},
					},
				},
				UserData: aws.String(""),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			mockAWSClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(tc.securityGroupOutput, tc.securityGroupErr).AnyTimes()
			mockAWSClient.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(nil, tc.azErr).AnyTimes()
			mockAWSClient.EXPECT().DescribeSubnets(gomock.Any()).Return(tc.subnetOutput, tc.subnetErr).AnyTimes()
			mockAWSClient.EXPECT().DescribeImages(gomock.Any()).Return(tc.imageOutput, tc.imageErr).AnyTimes()
			mockAWSClient.EXPECT().RunInstances(tc.runInstancesInput).Return(tc.instancesOutput, tc.instancesErr).AnyTimes()

			_, launchErr := launchInstance(machine, tc.providerConfig, nil, mockAWSClient, tc.infra)
			t.Log(launchErr)
			if launchErr == nil {
				if !tc.succeeds {
					t.Errorf("Call to launchInstance did not fail as expected")
				}
			} else {
				if tc.succeeds {
					t.Errorf("Call to launchInstance did not succeed as expected")
				}
			}
		})
	}
}

func TestSortInstances(t *testing.T) {
	instances := []*ec2.Instance{
		{
			LaunchTime: aws.Time(time.Now()),
		},
		{
			LaunchTime: nil,
		},
		{
			LaunchTime: nil,
		},
		{
			LaunchTime: aws.Time(time.Now()),
		},
	}
	sortInstances(instances)
}

func TestGetInstanceMarketOptionsRequest(t *testing.T) {
	testCases := []struct {
		name              string
		spotMarketOptions *machinev1.SpotMarketOptions
		expectedRequest   *ec2.InstanceMarketOptionsRequest
	}{
		{
			name:              "with no Spot options specified",
			spotMarketOptions: nil,
			expectedRequest:   nil,
		},
		{
			name:              "with an empty Spot options specified",
			spotMarketOptions: &machinev1.SpotMarketOptions{},
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
			spotMarketOptions: &machinev1.SpotMarketOptions{
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
			spotMarketOptions: &machinev1.SpotMarketOptions{
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
			providerConfig := &machinev1.AWSMachineProviderConfig{
				SpotMarketOptions: tc.spotMarketOptions,
			}

			request := getInstanceMarketOptionsRequest(providerConfig)
			if !reflect.DeepEqual(request, tc.expectedRequest) {
				t.Errorf("Case: %s. Got: %v, expected: %v", tc.name, request, tc.expectedRequest)
			}
		})
	}
}

func TestCorrectExistingTags(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}
	clusterID, _ := getClusterID(machine)
	instance := ec2.Instance{
		InstanceId: aws.String(stubInstanceID),
	}
	testCases := []struct {
		name               string
		tags               []*ec2.Tag
		expectedCreateTags bool
	}{
		{
			name: "Valid Tags",
			tags: []*ec2.Tag{
				{
					Key:   aws.String("kubernetes.io/cluster/" + clusterID),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(machine.Name),
				},
			},
			expectedCreateTags: false,
		},
		{
			name: "Invalid Name Tag Correct Cluster",
			tags: []*ec2.Tag{
				{
					Key:   aws.String("kubernetes.io/cluster/" + clusterID),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String("badname"),
				},
			},
			expectedCreateTags: true,
		},
		{
			name: "Invalid Cluster Tag Correct Name",
			tags: []*ec2.Tag{
				{
					Key:   aws.String("kubernetes.io/cluster/" + "badcluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String(machine.Name),
				},
			},
			expectedCreateTags: true,
		},
		{
			name: "Both Tags Wrong",
			tags: []*ec2.Tag{
				{
					Key:   aws.String("kubernetes.io/cluster/" + clusterID),
					Value: aws.String("bad value"),
				},
				{
					Key:   aws.String("Name"),
					Value: aws.String("bad name"),
				},
			},
			expectedCreateTags: true,
		},
		{
			name:               "No Tags",
			tags:               nil,
			expectedCreateTags: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			// if Finish is not called, MinTimes is never enforced
			defer mockCtrl.Finish()
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			instance.Tags = tc.tags

			if tc.expectedCreateTags {
				mockAWSClient.EXPECT().CreateTags(gomock.Any()).Return(&ec2.CreateTagsOutput{}, nil).MinTimes(1)
			}

			err := correctExistingTags(machine, &instance, mockAWSClient)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
