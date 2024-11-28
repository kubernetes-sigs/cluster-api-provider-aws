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
	"testing"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(lifecycleHookNeedsUpdate(&tt.existing, &tt.expected)).To(Equal(tt.wantUpdate))
		})
	}
}
