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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/controllers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/instancestate/mock_sqsiface"
)

func TestAWSInstanceStateController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	sqsSvs = mock_sqsiface.NewMockSQSAPI(mockCtrl)
	instanceStateReconciler = &AwsInstanceStateReconciler{
		Client: testEnv.Client,
		Log:    ctrl.Log.WithName("controllers").WithName("AWSInstanceState"),
		sqsServiceFactory: func() sqsiface.SQSAPI {
			return sqsSvs
		},
	}
	defer mockCtrl.Finish()

	t.Run("should maintain list of cluster queue URLs and reconcile failing machines", func(t *testing.T) {
		g := NewWithT(t)

		failingMachineMeta := metav1.ObjectMeta{
			Name:      "aws-cluster-1-instance-1",
			Namespace: "default",
		}
		sqsSvs.EXPECT().GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String("aws-cluster-1-queue")}).AnyTimes().
			Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("aws-cluster-1-url")}, nil)
		sqsSvs.EXPECT().GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String("aws-cluster-2-queue")}).AnyTimes().
			Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("aws-cluster-2-url")}, nil)
		sqsSvs.EXPECT().GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String("aws-cluster-3-queue")}).AnyTimes().
			Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String("aws-cluster-3-url")}, nil)
		sqsSvs.EXPECT().ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String("aws-cluster-1-url")}).AnyTimes().
			DoAndReturn(func(arg *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
				m := &infrav1.AWSMachine{}
				lookupKey := types.NamespacedName{
					Namespace: failingMachineMeta.Namespace,
					Name:      failingMachineMeta.Name,
				}
				err := k8sClient.Get(context.TODO(), lookupKey, m)
				// start returning a message once the AWSMachine is available
				if err == nil {
					return &sqs.ReceiveMessageOutput{
						Messages: []*sqs.Message{{
							ReceiptHandle: aws.String("message-receipt-handle"),
							Body:          aws.String(messageBodyJSON),
						}},
					}, nil
				}

				return &sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil
			})

		sqsSvs.EXPECT().ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String("aws-cluster-2-url")}).AnyTimes().
			Return(&sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil)
		sqsSvs.EXPECT().ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String("aws-cluster-3-url")}).AnyTimes().
			Return(&sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil)
		sqsSvs.EXPECT().DeleteMessage(&sqs.DeleteMessageInput{QueueUrl: aws.String("aws-cluster-1-url"), ReceiptHandle: aws.String("message-receipt-handle")}).AnyTimes().
			Return(nil, nil)

		g.Expect(testEnv.Manager.GetFieldIndexer().IndexField(context.Background(), &infrav1.AWSMachine{},
			controllers.InstanceIDIndex,
			func(o client.Object) []string {
				m := o.(*infrav1.AWSMachine)
				if m.Spec.InstanceID != nil {
					return []string{*m.Spec.InstanceID}
				}
				return nil
			},
		)).ToNot(HaveOccurred())

		err := instanceStateReconciler.SetupWithManager(context.Background(), testEnv.Manager, controller.Options{})
		g.Expect(err).ToNot(HaveOccurred())
		go func() {
			fmt.Println("Starting the manager")
			if err := testEnv.StartManager(ctx); err != nil {
				panic(fmt.Sprintf("Failed to start the envtest manager: %v", err))
			}
		}()
		testEnv.WaitForWebhooks()

		k8sClient = testEnv.GetClient()

		persistObject(g, createAWSCluster("aws-cluster-1"))
		persistObject(g, createAWSCluster("aws-cluster-2"))

		machine1 := &infrav1.AWSMachine{
			Spec: infrav1.AWSMachineSpec{
				InstanceID:   pointer.String("i-failing-instance-1"),
				InstanceType: "test",
			},
			ObjectMeta: failingMachineMeta,
		}
		persistObject(g, machine1)

		t.Log("Ensuring queue URLs are up-to-date")
		g.Eventually(func() bool {
			exist := true
			for _, cluster := range []string{"aws-cluster-1", "aws-cluster-2"} {
				_, ok := instanceStateReconciler.queueURLs.Load(cluster)
				exist = exist && ok
			}
			return exist
		}, 10*time.Second).Should(Equal(true))

		deleteAWSCluster(g, "aws-cluster-2")
		t.Log("Ensuring we stop tracking deleted queue")
		g.Eventually(func() bool {
			_, ok := instanceStateReconciler.queueURLs.Load("aws-cluster-2")
			return ok
		}, 10*time.Second).Should(BeFalse())

		persistObject(g, createAWSCluster("aws-cluster-3"))
		t.Log("Ensuring newly created cluster is added to tracked clusters")
		g.Eventually(func() bool {
			exist := true
			for _, cluster := range []string{"aws-cluster-1", "aws-cluster-3"} {
				_, ok := instanceStateReconciler.queueURLs.Load(cluster)
				exist = exist && ok
			}
			return exist
		}, 10*time.Second).Should(Equal(true))

		t.Log("Ensuring machine is labelled with correct instance state")
		g.Eventually(func() bool {
			m := &infrav1.AWSMachine{}
			key := types.NamespacedName{
				Namespace: failingMachineMeta.Namespace,
				Name:      failingMachineMeta.Name,
			}
			g.Expect(k8sClient.Get(context.TODO(), key, m)).NotTo(HaveOccurred())
			labels := m.GetLabels()
			val := labels[Ec2InstanceStateLabelKey]
			return val == "shutting-down"
		}, 10*time.Second).Should(Equal(true))
	})
}

const messageBodyJSON = `{
	"source": "aws.ec2",
	"detail-type": "EC2 Instance State-change Notification",
	"detail": {
		"instance-id": "i-failing-instance-1",
		"state": "shutting-down"
	}
}`
