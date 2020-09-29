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
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/instancestate/mock_sqsiface"
)

func TestReconcileSQSQueue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name      string
		expect    func(m *mock_sqsiface.MockSQSAPIMockRecorder)
		expectErr bool
	}{
		{
			name: "successfully creates an SQS queue",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNameReceiveMessageWaitTimeSeconds] = "20"
				m.CreateQueue(&sqs.CreateQueueInput{
					QueueName:  aws.String("test-cluster-queue"),
					Attributes: aws.StringMap(attrs),
				}).Return(nil, nil)
			},
			expectErr: false,
		},
		{
			name: "does not error if queue already exists",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNameReceiveMessageWaitTimeSeconds] = "20"
				m.CreateQueue(&sqs.CreateQueueInput{
					QueueName:  aws.String("test-cluster-queue"),
					Attributes: aws.StringMap(attrs),
				}).Return(nil, awserr.New(sqs.ErrCodeQueueNameExists, "", nil))
			},
			expectErr: false,
		},
		{
			name: "errors when unexpected error occurs",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNameReceiveMessageWaitTimeSeconds] = "20"
				m.CreateQueue(&sqs.CreateQueueInput{
					QueueName:  aws.String("test-cluster-queue"),
					Attributes: aws.StringMap(attrs),
				}).Return(nil, errors.New("some error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			sqsMock := mock_sqsiface.NewMockSQSAPI(mockCtrl)
			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))

			tc.expect(sqsMock.EXPECT())
			s := NewService(clusterScope)
			s.SQSClient = sqsMock

			err = s.reconcileSQSQueue()

			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestDeleteSQSQueue(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name      string
		expect    func(m *mock_sqsiface.MockSQSAPIMockRecorder)
		expectErr bool
	}{
		{
			name: "deletes queue successfully",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				}).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("test-cluster-queue-url")}, nil)
				m.DeleteQueue(&sqs.DeleteQueueInput{
					QueueUrl: aws.String("test-cluster-queue-url"),
				}).Return(nil, nil)
			},
			expectErr: false,
		},
		{
			name: "doesn't return error if queue not found when calling GetQueueUrl",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				}).Return(nil, awserr.New(sqs.ErrCodeQueueDoesNotExist, "", nil))
			},
			expectErr: false,
		},
		{
			name: "returns error if Describe Queue failed for unexpected reason",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				}).Return(nil, errors.New("some error"))
			},
			expectErr: true,
		},
		{
			name: "doesn't return error if queue not found when attempting delete",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				}).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("test-cluster-queue-url")}, nil)
				m.DeleteQueue(&sqs.DeleteQueueInput{
					QueueUrl: aws.String("test-cluster-queue-url"),
				}).Return(nil, awserr.New(sqs.ErrCodeQueueDoesNotExist, "", nil))
			},
			expectErr: false,
		},
		{
			name: "returns error if delete queue failed for unexpected reason",
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				m.GetQueueUrl(&sqs.GetQueueUrlInput{
					QueueName: aws.String("test-cluster-queue"),
				}).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("test-cluster-queue-url")}, nil)
				m.DeleteQueue(&sqs.DeleteQueueInput{
					QueueUrl: aws.String("test-cluster-queue-url"),
				}).Return(nil, errors.New("some error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			sqsMock := mock_sqsiface.NewMockSQSAPI(mockCtrl)
			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))

			tc.expect(sqsMock.EXPECT())
			s := NewService(clusterScope)
			s.SQSClient = sqsMock

			err = s.deleteSQSQueue()

			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestCreatePolicyForRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name      string
		input     *createPolicyForRuleInput
		expect    func(m *mock_sqsiface.MockSQSAPIMockRecorder)
		expectErr bool
	}{
		{
			name: "creates a policy for a given rule",
			input: &createPolicyForRuleInput{
				QueueArn: "test-cluster-queue-arn",
				QueueURL: "test-cluster-queue-url",
				RuleArn:  "test-cluster-rule-arn",
			},
			expect: func(m *mock_sqsiface.MockSQSAPIMockRecorder) {
				buffer := new(bytes.Buffer)
				_ = json.Compact(buffer, []byte(expectedPolicyJSON))
				attrs := make(map[string]string)
				attrs[sqs.QueueAttributeNamePolicy] = buffer.String()
				m.SetQueueAttributes(&sqs.SetQueueAttributesInput{
					QueueUrl:   aws.String("test-cluster-queue-url"),
					Attributes: aws.StringMap(attrs),
				}).Return(nil, nil)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			sqsMock := mock_sqsiface.NewMockSQSAPI(mockCtrl)
			clusterScope, err := setupCluster("test-cluster")
			g.Expect(err).To(Not(HaveOccurred()))

			tc.expect(sqsMock.EXPECT())
			s := NewService(clusterScope)
			s.SQSClient = sqsMock

			err = s.createPolicyForRule(tc.input)

			if tc.expectErr {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestGenerateQueueName(t *testing.T) {
	testCases := []struct {
		name              string
		clusterName       string
		expectedQueueName string
	}{
		{
			name:              "unchanged when cluster name doesn't include a .",
			clusterName:       "test-cluster",
			expectedQueueName: "test-cluster-queue",
		},
		{
			name:              "replaces . with - in cluster name",
			clusterName:       "some.cluster.name",
			expectedQueueName: "some-cluster-name-queue",
		},
	}

	for _, tc := range testCases {
		g := NewWithT(t)
		g.Expect(GenerateQueueName(tc.clusterName)).To(Equal(tc.expectedQueueName))
	}
}

const expectedPolicyJSON = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CAPAEvents_test-cluster-ec2-rule_test-cluster-queue",
      "Principal": {
        "Service": [
          "events.amazonaws.com"
        ]
      },
      "Effect": "Allow",
      "Action": [
        "sqs:SendMessage"
      ],
      "Resource": [
        "test-cluster-queue-arn"
      ],
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "test-cluster-rule-arn"
        }
      }
    }
  ],
  "Id": "test-cluster-queue-arn"
}`
