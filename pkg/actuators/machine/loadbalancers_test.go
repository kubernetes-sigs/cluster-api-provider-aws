package machine

import (
	"fmt"
	"testing"

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

	instance := stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr)
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), tc.targetGroupErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, tc.registerTargetErr).AnyTimes()
			registerWithNetworkLoadBalancers(mockAWSClient, []string{"name1", "name2"}, instance)
		})
	}
}
