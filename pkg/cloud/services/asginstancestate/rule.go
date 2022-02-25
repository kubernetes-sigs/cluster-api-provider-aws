/*
Copyright 2020 The Kubernetes Authors.

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

package asginstancestate

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

// reconcileRules creates rules and attaches the queue as a target.
func (s Service) reconcileRules() error {
	var ruleNotFound bool
	ruleResp, err := s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
		Name: aws.String(s.getASGRuleName()),
	})
	if err != nil {
		if resourceNotFoundError(err) {
			ruleNotFound = true
		} else {
			return errors.Wrapf(err, "unable to describe rule %s", s.getASGRuleName())
		}
	}

	if ruleNotFound {
		err = s.createRule()
		if err != nil {
			return errors.Wrap(err, "unable to create rule")
		}
		// fetch newly created rule
		ruleResp, err = s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
			Name: aws.String(s.getASGRuleName()),
		})

		if err != nil {
			return errors.Wrapf(err, "unable to describe new rule %s", s.getASGRuleName())
		}
	}

	queueURLResp, err := s.SQSClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(GenerateQueueName(s.scope.Name())),
	})

	if err != nil {
		return errors.Wrap(err, "unable to get queue URL")
	}
	queueAttrs, err := s.SQSClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: aws.StringSlice([]string{sqs.QueueAttributeNameQueueArn, sqs.QueueAttributeNamePolicy}),
		QueueUrl:       queueURLResp.QueueUrl,
	})

	if err != nil {
		return errors.Wrap(err, "unable to get queue attributes")
	}

	targetsResp, err := s.EventBridgeClient.ListTargetsByRule(&eventbridge.ListTargetsByRuleInput{
		Rule: aws.String(s.getASGRuleName()),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to list targets for rule %s", s.getASGRuleName())
	}

	targetFound := false
	for _, target := range targetsResp.Targets {
		// check if queue is already added as a target
		if *target.Id == GenerateQueueName(s.scope.Name()) && *target.Arn == *queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn] {
			targetFound = true
		}
	}

	if !targetFound {
		_, err = s.EventBridgeClient.PutTargets(&eventbridge.PutTargetsInput{
			Rule: ruleResp.Name,
			Targets: []*eventbridge.Target{{
				Arn: queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn],
				Id:  aws.String(GenerateQueueName(s.scope.Name())),
			}},
		})

		if err != nil {
			return errors.Wrapf(err, "unable to add SQS target %s to rule %s", GenerateQueueName(s.scope.Name()), s.getASGRuleName())
		}
	}

	if queueAttrs.Attributes[sqs.QueueAttributeNamePolicy] == nil {
		// add a policy for the rule so the rule is authorized to emit messages to the queue
		err = s.createPolicyForRule(&createPolicyForRuleInput{
			QueueArn: *queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn],
			QueueURL: *queueURLResp.QueueUrl,
			RuleArn:  *ruleResp.Arn,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) createRule() error {
	eventPattern := eventPattern{
		Source:     []string{"aws.autoscaling"},
	}
	data, _ := json.Marshal(eventPattern)
	_, err := s.EventBridgeClient.PutRule(&eventbridge.PutRuleInput{
		Name:         aws.String(s.getASGRuleName()),
		EventPattern: aws.String(string(data)),
		State:        aws.String(eventbridge.RuleStateEnabled),
	})

	return err
}

func (s Service) deleteRules() error {
	_, err := s.EventBridgeClient.RemoveTargets(&eventbridge.RemoveTargetsInput{
		Rule: aws.String(s.getASGRuleName()),
		Ids:  aws.StringSlice([]string{GenerateQueueName(s.scope.Name())}),
	})
	if err != nil && !resourceNotFoundError(err) {
		return errors.Wrapf(err, "unable to remove target %s for rule %s", GenerateQueueName(s.scope.Name()), s.getASGRuleName())
	}
	_, err = s.EventBridgeClient.DeleteRule(&eventbridge.DeleteRuleInput{
		Name: aws.String(s.getASGRuleName()),
	})

	if err != nil && resourceNotFoundError(err) {
		return nil
	}

	return err
}

func (s Service) getASGRuleName() string {
	return fmt.Sprintf("%s-asg-rule", s.scope.Name())
}

func resourceNotFoundError(err error) bool {
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == eventbridge.ErrCodeResourceNotFoundException {
		return true
	}
	return false
}

type eventPattern struct {
	Source      []string     `json:"source"`
}
