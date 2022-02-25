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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/instancestate/mock_eventbridgeiface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/instancestate/mock_sqsiface"
)

func TestReconcileRules(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ruleName := "test-cluster-ec2-rule"

	testCases := []struct {
		name                        string
		eventBridgeExpect           func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder)
		postCreateEventBridgeExpect func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder)
		sqsExpect                   func(m *mock_sqsiface.MockSQSAPIMockRecorder)
		expectErr                   bool
	}{
		{
			name: "successfully creates missing rule and target",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.DescribeRule(gomock.Eq(&eventbridge.DescribeRuleInput{
					Name: aws.String(ruleName),
				})).Return(nil, awserr.New(eventbridge.ErrCodeResourceNotFoundException, "", nil))
				e := &eventPattern{
					Source:     []string{"aws.autoscaling"},
				}
				data, _ := json.Marshal(e)
				m.PutRule(gomock.Eq(&eventbridge.PutRuleInput{
					Name:         aws.String(ruleName),
					State:        aws.String(eventbridge.RuleStateDisabled),
					EventPattern: aws.String(string(data)),
				}))
			},
			postCreateEventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.DescribeRule(gomock.Eq(&eventbridge.DescribeRuleInput{
					Name: aws.String(ruleName),
				})).Return(&eventbridge.DescribeRuleOutput{Name: aws.String(ruleName), Arn: aws.String("rule-arn")}, nil)
				m.ListTargetsByRule(&eventbridge.ListTargetsByRuleInput{
					Rule: aws.String(ruleName),
				}).Return(&eventbridge.ListTargetsByRuleOutput{
					Targets: []*eventbridge.Target{{
						Id:  aws.String("another-queue"),
						Arn: aws.String("another-queue-arn"),
					}},
				}, nil)
				m.PutTargets(gomock.Eq(&eventbridge.PutTargetsInput{
					Rule: aws.String(ruleName),
					Targets: []*eventbridge.Target{{
						Arn: aws.String("test-cluster-queue-arn"),
						Id:  aws.String("test-cluster-queue"),
					}},
				}))
			},
			sqsExpect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(gomock.Eq(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				})).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("test-cluster-queue-url")}, nil)
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNameQueueArn] = "test-cluster-queue-arn"
				m.GetQueueAttributes(gomock.Eq(&sqs.GetQueueAttributesInput{
					AttributeNames: aws.StringSlice([]string{sqs.QueueAttributeNameQueueArn, sqs.QueueAttributeNamePolicy}),
					QueueUrl:       aws.String("test-cluster-queue-url"),
				})).Return(&sqs.GetQueueAttributesOutput{Attributes: aws.StringMap(attrs)}, nil)
				m.SetQueueAttributes(gomock.AssignableToTypeOf(&sqs.SetQueueAttributesInput{})).Return(nil, nil)
			},
			expectErr: false,
		},
		{
			name: "skips creating target and queue policy if they already exist",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.DescribeRule(gomock.Eq(&eventbridge.DescribeRuleInput{
					Name: aws.String(ruleName),
				})).Return(&eventbridge.DescribeRuleOutput{Name: aws.String(ruleName), Arn: aws.String("rule-arn")}, nil)
				m.ListTargetsByRule(gomock.AssignableToTypeOf(&eventbridge.ListTargetsByRuleInput{})).Return(&eventbridge.ListTargetsByRuleOutput{
					Targets: []*eventbridge.Target{{
						Id:  aws.String("test-cluster-queue"),
						Arn: aws.String("test-cluster-queue-arn"),
					}},
				}, nil)
			},
			postCreateEventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {},
			sqsExpect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(gomock.AssignableToTypeOf(&sqs.GetQueueUrlInput{})).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("test-cluster-queue-url")}, nil)
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNameQueueArn] = "test-cluster-queue-arn"
				attrs[sqs.QueueAttributeNamePolicy] = "some policy"
				m.GetQueueAttributes(gomock.AssignableToTypeOf(&sqs.GetQueueAttributesInput{})).Return(&sqs.GetQueueAttributesOutput{Attributes: aws.StringMap(attrs)}, nil)
			},
		},
		{
			name: "returns error if DescribeRule runs into unexpected error",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.DescribeRule(gomock.Eq(&eventbridge.DescribeRuleInput{
					Name: aws.String(ruleName),
				})).Return(nil, errors.New("some error"))
			},
			postCreateEventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {},
			sqsExpect:                   func(m *mock_sqsiface.MockSQSAPIMockRecorder) {},
			expectErr:                   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			eventbridgeMock := mock_eventbridgeiface.NewMockEventBridgeAPI(mockCtrl)
			sqsMock := mock_sqsiface.NewMockSQSAPI(mockCtrl)
			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))
			tc.sqsExpect(sqsMock.EXPECT())
			tc.eventBridgeExpect(eventbridgeMock.EXPECT())
			tc.postCreateEventBridgeExpect(eventbridgeMock.EXPECT())

			s := NewService(clusterScope)
			s.EventBridgeClient = eventbridgeMock
			s.SQSClient = sqsMock

			err = s.reconcileRules()
			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestDeleteRules(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name              string
		eventBridgeExpect func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder)
		expectErr         bool
	}{
		{
			name: "removes target and ec2 rule successfully when they both exist",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.RemoveTargets(gomock.Eq(&eventbridge.RemoveTargetsInput{
					Rule: aws.String("test-cluster-ec2-rule"),
					Ids:  aws.StringSlice([]string{"test-cluster-queue"}),
				})).Return(nil, nil)
				m.DeleteRule(gomock.Eq(&eventbridge.DeleteRuleInput{
					Name: aws.String("test-cluster-ec2-rule"),
				})).Return(nil, nil)
			},
			expectErr: false,
		},
		{
			name: "continues to remove rule when target doesn't exist",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.RemoveTargets(gomock.AssignableToTypeOf(&eventbridge.RemoveTargetsInput{})).
					Return(nil, awserr.New(eventbridge.ErrCodeResourceNotFoundException, "", nil))
				m.DeleteRule(gomock.Eq(&eventbridge.DeleteRuleInput{
					Name: aws.String("test-cluster-ec2-rule"),
				})).Return(nil, nil)
			},
			expectErr: false,
		},
		{
			name: "returns error when remove target fails unexpectedly",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.RemoveTargets(gomock.AssignableToTypeOf(&eventbridge.RemoveTargetsInput{})).Return(nil, errors.New("some error"))
			},
			expectErr: true,
		},
		{
			name: "returns error when delete rule fails unexpectedly",
			eventBridgeExpect: func(m *mock_eventbridgeiface.MockEventBridgeAPIMockRecorder) {
				m.RemoveTargets(gomock.AssignableToTypeOf(&eventbridge.RemoveTargetsInput{})).Return(nil, nil)
				m.DeleteRule(gomock.AssignableToTypeOf(&eventbridge.DeleteRuleInput{})).Return(nil, errors.New("some error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			eventbridgeMock := mock_eventbridgeiface.NewMockEventBridgeAPI(mockCtrl)
			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))
			tc.eventBridgeExpect(eventbridgeMock.EXPECT())

			s := NewService(clusterScope)
			s.EventBridgeClient = eventbridgeMock

			err = s.deleteRules()
			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
