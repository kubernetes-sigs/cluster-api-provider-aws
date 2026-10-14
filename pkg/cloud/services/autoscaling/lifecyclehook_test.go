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

package asg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
)

func TestLifecycleHookNeedsUpdate(t *testing.T) {
	defaultResultAbandon := expinfrav1.LifecycleHookDefaultResultAbandon
	defaultResultContinue := expinfrav1.LifecycleHookDefaultResultContinue

	tests := []struct {
		name       string
		existing   expinfrav1.AWSLifecycleHook
		expected   expinfrav1.AWSLifecycleHook
		wantUpdate bool
	}{
		{
			name: "exactly equal",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			wantUpdate: false,
		},

		{
			name: "exactly equal (all fields filled)",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			wantUpdate: false,
		},

		{
			name: "heartbeatTimeout and defaultResult not set in manifest, but set to defaults by AWS",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				// Describing with AWS SDK always fills these fields with the defaults
				HeartbeatTimeout: &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:    &defaultResultAbandon,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
			},
			wantUpdate: false,
		},

		{
			name: "transition differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_LAUNCHING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			wantUpdate: true,
		},

		{
			name: "heartbeat timeout differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3601 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			wantUpdate: true,
		},

		{
			name: "default result differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultAbandon,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                "andreas-test",
				LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:    &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:       &defaultResultContinue,
			},
			wantUpdate: true,
		},

		{
			name: "role ARN differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification2"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			wantUpdate: true,
		},

		{
			name: "notification target ARN differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                  "test",
				NotificationTargetARN: ptr.To("arn:aws:sqs:eu-west-1:123456789012:mycluster-nth2"),
				RoleARN:               ptr.To("arn:aws:iam::123456789012:role/mycluster-nth-notification"),
				LifecycleTransition:   "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:      &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:         &defaultResultAbandon,
				NotificationMetadata:  nil,
			},
			wantUpdate: true,
		},

		{
			name: "notification metadata both empty",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                 "test",
				LifecycleTransition:  "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:     &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:        &defaultResultAbandon,
				NotificationMetadata: nil,
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                 "test",
				LifecycleTransition:  "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:     &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:        &defaultResultAbandon,
				NotificationMetadata: ptr.To(""),
			},
			wantUpdate: false,
		},

		{
			name: "notification metadata differs",
			existing: expinfrav1.AWSLifecycleHook{
				Name:                 "test",
				LifecycleTransition:  "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:     &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:        &defaultResultAbandon,
				NotificationMetadata: ptr.To("abc"),
			},
			expected: expinfrav1.AWSLifecycleHook{
				Name:                 "test",
				LifecycleTransition:  "autoscaling:EC2_INSTANCE_TERMINATING",
				HeartbeatTimeout:     &metav1.Duration{Duration: 3600 * time.Second},
				DefaultResult:        &defaultResultAbandon,
				NotificationMetadata: ptr.To("xyz"),
			},
			wantUpdate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := tt.existing
			expected := tt.expected
			g := NewWithT(t)
			g.Expect(lifecycleHookNeedsUpdate(&existing, &expected)).To(Equal(tt.wantUpdate))
		})
	}
}

func TestLifecycleHookDualWrite(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	service := mock_services.NewMockASGInterface(mockCtrl)
	pool := &clusterv1.MachinePool{ObjectMeta: metav1.ObjectMeta{Generation: 3}}
	hook := &expinfrav1.AWSLifecycleHook{Name: "hook"}
	failure := errors.New("creation failed")
	service.EXPECT().CreateLifecycleHook(ctx, "asg", hook).Return(failure)
	g.Expect(reconcileLifecycleHook(ctx, service, "asg", hook, nil, pool, logger.FromContext(ctx))).To(MatchError(failure))
	legacy := deprecatedv1beta1conditions.Get(pool, clusterv1.ConditionType(expinfrav1.LifecycleHookReadyCondition))
	g.Expect(legacy).ToNot(BeNil())
	standard := conditions.Get(pool, expinfrav1.AWSMachinePoolLifecycleHookReadyV1Beta2Condition)
	g.Expect(standard).ToNot(BeNil())
	g.Expect(standard.Status).To(Equal(metav1.ConditionFalse))
	g.Expect(standard.Reason).To(Equal(expinfrav1.AWSMachinePoolLifecycleHookCreationFailedV1Beta2Reason))
	g.Expect(standard.Message).To(Equal(legacy.Message))
	g.Expect(standard.ObservedGeneration).To(Equal(int64(3)))

	// Observing the desired hook clears both failure conditions.
	g.Expect(reconcileLifecycleHook(ctx, service, "asg", hook, []*expinfrav1.AWSLifecycleHook{hook}, pool, logger.FromContext(ctx))).To(Succeed())
	g.Expect(deprecatedv1beta1conditions.IsTrue(pool, clusterv1.ConditionType(expinfrav1.LifecycleHookReadyCondition))).To(BeTrue())
	g.Expect(conditions.IsTrue(pool, expinfrav1.AWSMachinePoolLifecycleHookReadyV1Beta2Condition)).To(BeTrue())
	g.Expect(conditions.Get(pool, expinfrav1.AWSMachinePoolLifecycleHookReadyV1Beta2Condition).Message).To(BeEmpty())
}
