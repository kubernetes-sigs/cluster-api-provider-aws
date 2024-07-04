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

package asg

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func (s *Service) LifecycleHookNeedsUpdate(scope scope.LifecycleHookScope, existing *expinfrav1.AWSLifecycleHook, expected *expinfrav1.AWSLifecycleHook) bool {
	return existing.DefaultResult != expected.DefaultResult ||
		existing.HeartbeatTimeout != expected.HeartbeatTimeout ||
		existing.LifecycleTransition != expected.LifecycleTransition ||
		existing.NotificationTargetARN != expected.NotificationTargetARN ||
		existing.NotificationMetadata != expected.NotificationMetadata
}

func (s *Service) GetLifecycleHooks(scope scope.LifecycleHookScope) ([]*expinfrav1.AWSLifecycleHook, error) {
	asgName := scope.GetASGName()
	input := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: aws.String(asgName),
	}

	out, err := s.ASGClient.DescribeLifecycleHooksWithContext(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe lifecycle hooks for AutoScalingGroup: %q", scope.GetASGName())
	}

	hooks := make([]*expinfrav1.AWSLifecycleHook, len(out.LifecycleHooks))
	for i, hook := range out.LifecycleHooks {
		hooks[i] = s.SDKToLifecycleHook(hook)
	}

	return hooks, nil
}

func (s *Service) GetLifecycleHook(scope scope.LifecycleHookScope, hook *expinfrav1.AWSLifecycleHook) (*expinfrav1.AWSLifecycleHook, error) {
	asgName := scope.GetASGName()
	input := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: aws.String(asgName),
		LifecycleHookNames:   []*string{aws.String(hook.Name)},
	}

	out, err := s.ASGClient.DescribeLifecycleHooksWithContext(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe lifecycle hook %q for AutoScalingGroup: %q", hook.Name, scope.GetASGName())
	}

	if len(out.LifecycleHooks) == 0 {
		return nil, nil
	}

	return s.SDKToLifecycleHook(out.LifecycleHooks[0]), nil
}

func (s *Service) CreateLifecycleHook(scope scope.LifecycleHookScope, hook *expinfrav1.AWSLifecycleHook) error {
	asgName := scope.GetASGName()
	input := &autoscaling.PutLifecycleHookInput{
		AutoScalingGroupName: aws.String(asgName),
		LifecycleHookName:    aws.String(hook.Name),
		LifecycleTransition:  aws.String(hook.LifecycleTransition.String()),
	}

	// Optional parameters
	if hook.DefaultResult != nil {
		input.DefaultResult = aws.String(hook.DefaultResult.String())
	}

	if hook.HeartbeatTimeout != nil {
		timeoutSeconds := hook.HeartbeatTimeout.Duration.Seconds()
		input.HeartbeatTimeout = aws.Int64(int64(timeoutSeconds))
	}

	if hook.NotificationTargetARN != nil {
		input.NotificationTargetARN = hook.NotificationTargetARN
	}

	if hook.RoleARN != nil {
		input.RoleARN = hook.RoleARN
	}

	if hook.NotificationMetadata != nil {
		input.NotificationMetadata = hook.NotificationMetadata
	}

	if _, err := s.ASGClient.PutLifecycleHookWithContext(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to create lifecycle hook %q for AutoScalingGroup: %q", hook.Name, scope.GetASGName())
	}

	return nil
}

func (s *Service) UpdateLifecycleHook(scope scope.LifecycleHookScope, hook *expinfrav1.AWSLifecycleHook) error {
	asgName := scope.GetASGName()
	input := &autoscaling.PutLifecycleHookInput{
		AutoScalingGroupName: aws.String(asgName),
		LifecycleHookName:    aws.String(hook.Name),
		LifecycleTransition:  aws.String(hook.LifecycleTransition.String()),
	}

	// Optional parameters
	if hook.DefaultResult != nil {
		input.DefaultResult = aws.String(hook.DefaultResult.String())
	}

	if hook.HeartbeatTimeout != nil {
		timeoutSeconds := hook.HeartbeatTimeout.Duration.Seconds()
		input.HeartbeatTimeout = aws.Int64(int64(timeoutSeconds))
	}

	if hook.NotificationTargetARN != nil {
		input.NotificationTargetARN = hook.NotificationTargetARN
	}

	if hook.RoleARN != nil {
		input.RoleARN = hook.RoleARN
	}

	if hook.NotificationMetadata != nil {
		input.NotificationMetadata = hook.NotificationMetadata
	}

	if _, err := s.ASGClient.PutLifecycleHookWithContext(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to update lifecycle hook %q for AutoScalingGroup: %q", hook.Name, scope.GetASGName())
	}

	return nil
}

func (s *Service) DeleteLifecycleHook(
	scope scope.LifecycleHookScope,
	hook *expinfrav1.AWSLifecycleHook,
) error {
	input := &autoscaling.DeleteLifecycleHookInput{
		AutoScalingGroupName: aws.String(scope.GetASGName()),
		LifecycleHookName:    aws.String(hook.Name),
	}

	if _, err := s.ASGClient.DeleteLifecycleHookWithContext(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to delete lifecycle hook %q for AutoScalingGroup: %q", hook.Name, scope.GetASGName())
	}

	return nil
}

func (s *Service) SDKToLifecycleHook(hook *autoscaling.LifecycleHook) *expinfrav1.AWSLifecycleHook {
	timeoutDuration := time.Duration(*hook.HeartbeatTimeout) * time.Second
	metav1Duration := metav1.Duration{Duration: timeoutDuration}
	defaultResult := expinfrav1.DefaultResult(*hook.DefaultResult)
	lifecycleTransition := expinfrav1.LifecycleTransition(*hook.LifecycleTransition)

	return &expinfrav1.AWSLifecycleHook{
		Name:                  *hook.LifecycleHookName,
		DefaultResult:         &defaultResult,
		HeartbeatTimeout:      &metav1Duration,
		LifecycleTransition:   lifecycleTransition,
		NotificationTargetARN: hook.NotificationTargetARN,
		RoleARN:               hook.RoleARN,
		NotificationMetadata:  hook.NotificationMetadata,
	}
}
