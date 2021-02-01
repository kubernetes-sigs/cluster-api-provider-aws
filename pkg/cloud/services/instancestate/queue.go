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

package instancestate

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/iam/v1alpha1"
	"strings"
)

func (s *Service) reconcileSQSQueue() error {
	attrs := make(map[string]string)
	attrs[sqs.QueueAttributeNameReceiveMessageWaitTimeSeconds] = "20"

	_, err := s.SQSClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName:  aws.String(GenerateQueueName(s.scope.Name())),
		Attributes: aws.StringMap(attrs),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == sqs.ErrCodeQueueNameExists {
				return nil
			}
		}
	}
	return errors.Wrap(err, "unable to create new queue")
}

func (s *Service) deleteSQSQueue() error {
	resp, err := s.SQSClient.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(GenerateQueueName(s.scope.Name()))})
	if err != nil {
		if queueNotFoundError(err) {
			return nil
		}
		return errors.Wrap(err, "unable to get queue URL")
	}
	_, err = s.SQSClient.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
	if err != nil && queueNotFoundError(err) {
		return nil
	}

	return errors.Wrap(err, "unable to delete queue")
}

func (s *Service) createPolicyForRule(input *createPolicyForRuleInput) error {
	attrs := make(map[string]string)
	policy := v1alpha1.PolicyDocument{
		Version: v1alpha1.CurrentVersion,
		ID:      input.QueueArn,
		Statement: v1alpha1.Statements{
			v1alpha1.StatementEntry{
				Sid:       fmt.Sprintf("CAPAEvents_%s_%s", s.getEC2RuleName(), GenerateQueueName(s.scope.Name())),
				Effect:    v1alpha1.EffectAllow,
				Principal: v1alpha1.Principals{v1alpha1.PrincipalService: v1alpha1.PrincipalID{"events.amazonaws.com"}},
				Action:    v1alpha1.Actions{"sqs:SendMessage"},
				Resource:  v1alpha1.Resources{input.QueueArn},
				Condition: v1alpha1.Conditions{
					"ArnEquals": map[string]string{"aws:SourceArn": input.RuleArn},
				},
			},
		},
	}
	policyData, err := json.Marshal(policy)
	if err != nil {
		return errors.Wrap(err, "unable to JSON marshal policy")
	}
	attrs[sqs.QueueAttributeNamePolicy] = string(policyData)

	_, err = s.SQSClient.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		QueueUrl:   aws.String(input.QueueURL),
		Attributes: aws.StringMap(attrs),
	})

	return errors.Wrap(err, "unable to update queue attributes")
}

func GenerateQueueName(clusterName string) string {
	adjusted := strings.Replace(clusterName, ".", "-", -1)
	return fmt.Sprintf("%s-queue", adjusted)
}

func queueNotFoundError(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			return true
		}
	}
	return false
}

type createPolicyForRuleInput struct {
	QueueArn string
	QueueURL string
	RuleArn  string
}
