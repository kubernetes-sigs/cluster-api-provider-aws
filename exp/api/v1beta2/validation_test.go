/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func TestValidateManagedMachinePoolScaling(t *testing.T) {
	tests := []struct {
		name    string
		scaling *ManagedMachinePoolScaling
		wantErr bool
	}{
		{
			name:    "nil scaling is valid",
			scaling: nil,
			wantErr: false,
		},
		{
			name: "valid scaling",
			scaling: &ManagedMachinePoolScaling{
				MinSize: ptr.To[int32](1),
				MaxSize: ptr.To[int32](10),
			},
			wantErr: false,
		},
		{
			name: "minSize 0 is valid",
			scaling: &ManagedMachinePoolScaling{
				MinSize: ptr.To[int32](0),
				MaxSize: ptr.To[int32](10),
			},
			wantErr: false,
		},
		{
			name: "minSize negative is invalid",
			scaling: &ManagedMachinePoolScaling{
				MinSize: ptr.To[int32](-1),
			},
			wantErr: true,
		},
		{
			name: "maxSize negative is invalid",
			scaling: &ManagedMachinePoolScaling{
				MaxSize: ptr.To[int32](-1),
			},
			wantErr: true,
		},
		{
			name: "minSize > maxSize is invalid",
			scaling: &ManagedMachinePoolScaling{
				MinSize: ptr.To[int32](10),
				MaxSize: ptr.To[int32](5),
			},
			wantErr: true,
		},
		{
			name: "minSize == maxSize is valid",
			scaling: &ManagedMachinePoolScaling{
				MinSize: ptr.To[int32](5),
				MaxSize: ptr.To[int32](5),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			errs := validateManagedMachinePoolScaling(tt.scaling, field.NewPath("spec", "scaling"))
			if tt.wantErr {
				g.Expect(errs).NotTo(BeEmpty())
			} else {
				g.Expect(errs).To(BeEmpty())
			}
		})
	}
}

func TestValidateManagedMachinePoolUpdateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *UpdateConfig
		wantErr bool
	}{
		{
			name:    "nil config is valid",
			config:  nil,
			wantErr: false,
		},
		{
			name: "valid with MaxUnavailable",
			config: &UpdateConfig{
				MaxUnavailable: aws.Int(1),
			},
			wantErr: false,
		},
		{
			name: "valid with MaxUnavailablePercentage",
			config: &UpdateConfig{
				MaxUnavailablePercentage: aws.Int(10),
			},
			wantErr: false,
		},
		{
			name:    "empty config is invalid",
			config:  &UpdateConfig{},
			wantErr: true,
		},
		{
			name: "both values is invalid",
			config: &UpdateConfig{
				MaxUnavailable:           aws.Int(1),
				MaxUnavailablePercentage: aws.Int(10),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			errs := validateManagedMachinePoolUpdateConfig(tt.config, field.NewPath("spec", "updateConfig"))
			if tt.wantErr {
				g.Expect(errs).NotTo(BeEmpty())
			} else {
				g.Expect(errs).To(BeEmpty())
			}
		})
	}
}

func TestValidateManagedMachinePoolRemoteAccess(t *testing.T) {
	tests := []struct {
		name    string
		access  *ManagedRemoteAccess
		wantErr bool
	}{
		{
			name:    "nil access is valid",
			access:  nil,
			wantErr: false,
		},
		{
			name: "private with sourceSecurityGroups is valid",
			access: &ManagedRemoteAccess{
				Public:               false,
				SourceSecurityGroups: []string{"sg-123"},
			},
			wantErr: false,
		},
		{
			name: "public with no sourceSecurityGroups is valid",
			access: &ManagedRemoteAccess{
				Public: true,
			},
			wantErr: false,
		},
		{
			name: "public with sourceSecurityGroups is invalid",
			access: &ManagedRemoteAccess{
				Public:               true,
				SourceSecurityGroups: []string{"sg-123"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			errs := validateManagedMachinePoolRemoteAccess(tt.access, field.NewPath("spec", "remoteAccess"))
			if tt.wantErr {
				g.Expect(errs).NotTo(BeEmpty())
			} else {
				g.Expect(errs).To(BeEmpty())
			}
		})
	}
}

func TestValidateManagedMachinePoolLaunchTemplate(t *testing.T) {
	tests := []struct {
		name         string
		lt           *AWSLaunchTemplate
		instanceType *string
		diskSize     *int32
		wantErr      bool
	}{
		{
			name:    "nil launch template is valid",
			lt:      nil,
			wantErr: false,
		},
		{
			name: "launch template without conflicts is valid",
			lt: &AWSLaunchTemplate{
				Name: "my-lt",
			},
			wantErr: false,
		},
		{
			name: "launch template with instanceType is invalid",
			lt: &AWSLaunchTemplate{
				Name: "my-lt",
			},
			instanceType: ptr.To[string]("t3.medium"),
			wantErr:      true,
		},
		{
			name: "launch template with diskSize is invalid",
			lt: &AWSLaunchTemplate{
				Name: "my-lt",
			},
			diskSize: ptr.To[int32](100),
			wantErr:  true,
		},
		{
			name: "launch template with IAM instance profile is invalid",
			lt: &AWSLaunchTemplate{
				Name:               "my-lt",
				IamInstanceProfile: "my-profile",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			errs := validateManagedMachinePoolLaunchTemplate(tt.lt, tt.instanceType, tt.diskSize, field.NewPath("spec"))
			if tt.wantErr {
				g.Expect(errs).NotTo(BeEmpty())
			} else {
				g.Expect(errs).To(BeEmpty())
			}
		})
	}
}

func TestValidateLifecycleHooks(t *testing.T) {
	tests := []struct {
		name    string
		hooks   []AWSLifecycleHook
		wantErr bool
	}{
		{
			name:    "nil hooks is valid",
			hooks:   nil,
			wantErr: false,
		},
		{
			name:    "empty hooks is valid",
			hooks:   []AWSLifecycleHook{},
			wantErr: false,
		},
		{
			name: "valid hook",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
				},
			},
			wantErr: false,
		},
		{
			name: "missing name is invalid",
			hooks: []AWSLifecycleHook{
				{
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
				},
			},
			wantErr: true,
		},
		{
			name: "notificationTargetARN without roleARN is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                  "my-hook",
					LifecycleTransition:   LifecycleHookTransitionInstanceLaunching,
					NotificationTargetARN: ptr.To[string]("arn:aws:sns:us-east-1:123456789012:my-topic"),
				},
			},
			wantErr: true,
		},
		{
			name: "roleARN without notificationTargetARN is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					RoleARN:             ptr.To[string]("arn:aws:iam::123456789012:role/my-role"),
				},
			},
			wantErr: true,
		},
		{
			name: "both notificationTargetARN and roleARN is valid",
			hooks: []AWSLifecycleHook{
				{
					Name:                  "my-hook",
					LifecycleTransition:   LifecycleHookTransitionInstanceLaunching,
					NotificationTargetARN: ptr.To[string]("arn:aws:sns:us-east-1:123456789012:my-topic"),
					RoleARN:               ptr.To[string]("arn:aws:iam::123456789012:role/my-role"),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid lifecycle transition is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid default result is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					DefaultResult:       ptr.To[LifecycleHookDefaultResult]("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "valid default result Continue",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					DefaultResult:       ptr.To[LifecycleHookDefaultResult](LifecycleHookDefaultResultContinue),
				},
			},
			wantErr: false,
		},
		{
			name: "valid default result Abandon",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					DefaultResult:       ptr.To[LifecycleHookDefaultResult](LifecycleHookDefaultResultAbandon),
				},
			},
			wantErr: false,
		},
		{
			name: "heartbeat timeout too low is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					HeartbeatTimeout:    &metav1.Duration{Duration: 29 * time.Second},
				},
			},
			wantErr: true,
		},
		{
			name: "heartbeat timeout too high is invalid",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					HeartbeatTimeout:    &metav1.Duration{Duration: 172801 * time.Second},
				},
			},
			wantErr: true,
		},
		{
			name: "valid heartbeat timeout",
			hooks: []AWSLifecycleHook{
				{
					Name:                "my-hook",
					LifecycleTransition: LifecycleHookTransitionInstanceLaunching,
					HeartbeatTimeout:    &metav1.Duration{Duration: 300 * time.Second},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			errs := validateLifecycleHooks(tt.hooks)
			if tt.wantErr {
				g.Expect(errs).NotTo(BeEmpty())
			} else {
				g.Expect(errs).To(BeEmpty())
			}
		})
	}
}
