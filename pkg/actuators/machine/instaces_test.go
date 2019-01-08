package machine

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"

	"github.com/golang/mock/gomock"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
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

func TestBuildEC2Filters(t *testing.T) {
	filter1 := "filter1"
	filter2 := "filter2"
	value1 := "A"
	value2 := "B"
	value3 := "C"

	inputFilters := []providerconfigv1.Filter{
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
							stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c"),
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
							stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c"),
							stubInstance("ami-a9acbbd7", "i-02fcb933c5da7085d"),
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

func TestRunningInstance(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}

	mockCtrl := gomock.NewController(t)
	mockAWSClient := mockaws.NewMockClient(mockCtrl)

	// Error describing instances
	mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						ImageId:    aws.String("ami-a9acbbd6"),
						InstanceId: aws.String("i-02fcb933c5da7085c"),
						State: &ec2.InstanceState{
							Name: aws.String("Running"),
						},
						LaunchTime: aws.Time(time.Now()),
					},
				},
			},
		},
	}, nil).AnyTimes()
	getRunningInstance(machine, mockAWSClient)
}

func TestLaunchInstance(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}

	cases := []struct {
		name                string
		providerConfig      *providerconfigv1.AWSMachineProviderConfig
		securityGroupOutput *ec2.DescribeSecurityGroupsOutput
		securityGroupErr    error
		subnetOutput        *ec2.DescribeSubnetsOutput
		subnetErr           error
		imageOutput         *ec2.DescribeImagesOutput
		imageErr            error
	}{
		{
			name: "Security groups with filters",
			providerConfig: stubPCSecurityGroups(
				[]providerconfigv1.AWSResourceReference{
					{
						Filters: []providerconfigv1.Filter{},
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
		},
		{
			name: "Security groups with filters with error",
			providerConfig: stubPCSecurityGroups(
				[]providerconfigv1.AWSResourceReference{
					{
						Filters: []providerconfigv1.Filter{},
					},
				},
			),
			securityGroupErr: fmt.Errorf("error"),
		},
		{
			name: "No security group",
			providerConfig: stubPCSecurityGroups(
				[]providerconfigv1.AWSResourceReference{
					{
						Filters: []providerconfigv1.Filter{},
					},
				},
			),
			securityGroupOutput: &ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []*ec2.SecurityGroup{},
			},
		},
		{
			name: "Subnet with filters",
			providerConfig: stubPCSubnet(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{},
			}),
			subnetOutput: &ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId: aws.String("subnetID"),
					},
				},
			},
		},
		{
			name: "Subnet with filters with error",
			providerConfig: stubPCSubnet(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{},
			}),
			subnetErr: fmt.Errorf("error"),
		},
		{
			name: "AMI with filters",
			providerConfig: stubPCAMI(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{},
			}),
			imageOutput: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{
					{
						CreationDate: aws.String(time.RFC3339),
						ImageId:      aws.String("ami-1111"),
					},
				},
			},
		},
		{
			name: "AMI with filters with error",
			providerConfig: stubPCAMI(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{},
			}),
			imageErr: fmt.Errorf("error"),
		},
		{
			name: "AMI with filters with no image",
			providerConfig: stubPCAMI(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{
					{
						Name:   "image_stage",
						Values: []string{"base"},
					},
				},
			}),
			imageOutput: &ec2.DescribeImagesOutput{
				Images: []*ec2.Image{},
			},
		},
		{
			name: "AMI with filters with two images",
			providerConfig: stubPCAMI(providerconfigv1.AWSResourceReference{
				Filters: []providerconfigv1.Filter{
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
		},
		{
			name:           "AMI not specified",
			providerConfig: stubPCAMI(providerconfigv1.AWSResourceReference{}),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			mockAWSClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(tc.securityGroupOutput, tc.securityGroupErr).AnyTimes()
			mockAWSClient.EXPECT().DescribeSubnets(gomock.Any()).Return(tc.subnetOutput, tc.subnetErr).AnyTimes()
			mockAWSClient.EXPECT().DescribeImages(gomock.Any()).Return(tc.imageOutput, tc.imageErr).AnyTimes()
			mockAWSClient.EXPECT().RunInstances(gomock.Any())

			launchInstance(machine, tc.providerConfig, nil, mockAWSClient)
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
