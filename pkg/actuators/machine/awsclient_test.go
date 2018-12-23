package machine

import (
	"fmt"
	"testing"
	"time"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
)

func TestAwsClient(t *testing.T) {
	machine, machinePC, err := stubMachine()
	if err != nil {
		t.Fatalf("Unable to build test machine manifest: %v", err)
	}

	machinePC.AMI = providerconfigv1.AWSResourceReference{}

	cases := []struct {
		name   string
		output *ec2.DescribeInstancesOutput
		err    error
	}{
		{
			name: "valid data",
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
			name: "valid data with more nil fields",
			output: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								ImageId:    aws.String("ami-a9acbbd6"),
								InstanceId: aws.String("i-02fcb933c5da7085c"),
								State: &ec2.InstanceState{
									Name: aws.String("Running"),
								},
								LaunchTime:       aws.Time(time.Now()),
								PublicDnsName:    aws.String(""),
								PrivateIpAddress: aws.String(""),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("key"),
										Value: aws.String("value"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "with error",
			err:  fmt.Errorf("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(tc.output, tc.err).AnyTimes()

			acw := NewAwsClientWrapper(mockAWSClient)
			acw.GetRunningInstances(machine)
			acw.GetPublicDNSName(machine)
			acw.GetPrivateIP(machine)
			acw.GetSecurityGroups(machine)
			acw.GetIAMRole(machine)
			acw.GetTags(machine)
			acw.GetSubnet(machine)
			acw.GetAvailabilityZone(machine)
		})
	}
}
