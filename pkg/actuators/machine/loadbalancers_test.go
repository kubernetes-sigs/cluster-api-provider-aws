package machine

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
)

func TestRegisterWithNetworkLoadBalancers(t *testing.T) {
	cases := []struct {
		name              string
		lbErr             error
		targetGroupErr    error
		registerTargetErr error
		err               error
	}{
		{
			name: "No error",
		},
		{
			name:  "With describe lb error",
			lbErr: fmt.Errorf("error"),
		},
		{
			name:           "With target group error",
			targetGroupErr: fmt.Errorf("error"),
		},
		{
			name:              "With register target error",
			registerTargetErr: fmt.Errorf("error"),
		},
	}

	instance := stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr)
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), tc.targetGroupErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, tc.registerTargetErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeTargetHealth(gomock.Any()).Return(&elbv2.DescribeTargetHealthOutput{}, nil).AnyTimes()
			registerWithNetworkLoadBalancers(mockAWSClient, []string{"name1", "name2"}, instance)
		})
	}
}

func TestDeregisterNetworkLoadBalancers(t *testing.T) {
	cases := []struct {
		name                           string
		instance                       *ec2.Instance
		lbErr                          error
		describeLoadBalancersCallTimes int
		targetGroupErr                 error
		describeTargetGroupsCallTimes  int
		unregisterTargetErr            error
		deregisterCallTimes            int
		expectErr                      error
	}{
		{
			name:     "No action if ip is unset",
			instance: stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", false),
		},
		{
			name:                           "No error",
			instance:                       stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true),
			describeLoadBalancersCallTimes: 1,
			describeTargetGroupsCallTimes:  1,
			deregisterCallTimes:            1,
		},
		{
			name:                           "With describe lb error",
			instance:                       stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true),
			lbErr:                          fmt.Errorf("error"),
			describeLoadBalancersCallTimes: 1,
			describeTargetGroupsCallTimes:  0,
			deregisterCallTimes:            0,
			expectErr:                      fmt.Errorf("error"),
		},
		{
			name:                           "With target group error",
			instance:                       stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true),
			targetGroupErr:                 fmt.Errorf("error"),
			describeLoadBalancersCallTimes: 1,
			describeTargetGroupsCallTimes:  1,
			deregisterCallTimes:            0,
			expectErr:                      fmt.Errorf("error"),
		},
		{
			name:                           "With target already unregistered error",
			instance:                       stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true),
			unregisterTargetErr:            awserr.New(elbv2.ErrCodeTargetGroupNotFoundException, "error", nil),
			describeLoadBalancersCallTimes: 1,
			describeTargetGroupsCallTimes:  1,
			deregisterCallTimes:            1,
			expectErr:                      nil,
		},
		{
			name:                           "With register target unknown error",
			instance:                       stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c", true),
			unregisterTargetErr:            fmt.Errorf("error"),
			describeLoadBalancersCallTimes: 1,
			describeTargetGroupsCallTimes:  1,
			deregisterCallTimes:            1,
			expectErr:                      fmt.Errorf("arn2: error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr).Times(tc.describeLoadBalancersCallTimes)
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), tc.targetGroupErr).Times(tc.describeTargetGroupsCallTimes)
			mockAWSClient.EXPECT().ELBv2DeregisterTargets(gomock.Any()).Return(nil, tc.unregisterTargetErr).Times(tc.deregisterCallTimes)
			err := deregisterNetworkLoadBalancers(mockAWSClient, []string{"name1", "name2"}, tc.instance)
			mockCtrl.Finish()

			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tc.expectErr) {
				t.Errorf("Unexpeted error output: expected '%s', got '%s'", tc.expectErr, err)
			}
		})
	}
}
