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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/controllers"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWSInstanceStateController", func() {
	It("should maintain list of cluster queue URLs and reconcile failing machines", func() {
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

		Expect(k8sManager.GetFieldIndexer().IndexField(&infrav1.AWSMachine{},
			controllers.InstanceIDIndex,
			func(o runtime.Object) []string {
				m := o.(*infrav1.AWSMachine)
				if m.Spec.InstanceID != nil {
					return []string{*m.Spec.InstanceID}
				}
				return nil
			},
		)).ToNot(HaveOccurred())
		err := instanceStateReconciler.SetupWithManager(k8sManager, controller.Options{})
		Expect(err).ToNot(HaveOccurred())
		stop := make(chan struct{})
		defer close(stop)
		go func() {
			defer GinkgoRecover()
			err := k8sManager.Start(stop)
			Expect(err).ToNot(HaveOccurred())
		}()

		k8sClient = k8sManager.GetClient()
		Expect(k8sClient).ToNot(BeNil())

		persistObject(createAWSCluster("aws-cluster-1"))
		persistObject(createAWSCluster("aws-cluster-2"))

		machine1 := &infrav1.AWSMachine{
			Spec: infrav1.AWSMachineSpec{
				InstanceID: pointer.StringPtr("i-failing-instance-1"),
			},
			ObjectMeta: failingMachineMeta,
		}
		persistObject(machine1)

		By("Ensuring queue URLs are up-to-date")
		Eventually(func() bool {
			exist := true
			for _, cluster := range []string{"aws-cluster-1", "aws-cluster-2"} {
				_, ok := instanceStateReconciler.queueURLs.Load(cluster)
				exist = exist && ok
			}
			return exist
		}, 10*time.Second).Should(Equal(true))

		deleteAWSCluster("aws-cluster-2")
		By("Ensuring we stop tracking deleted queue")
		Eventually(func() bool {
			_, ok := instanceStateReconciler.queueURLs.Load("aws-cluster-2")
			return ok
		}, 10*time.Second).Should(Equal(false))

		persistObject(createAWSCluster("aws-cluster-3"))
		By("Ensuring newly created cluster is added to tracked clusters")
		Eventually(func() bool {
			exist := true
			for _, cluster := range []string{"aws-cluster-1", "aws-cluster-3"} {
				_, ok := instanceStateReconciler.queueURLs.Load(cluster)
				exist = exist && ok
			}
			return exist
		}, 10*time.Second).Should(Equal(true))

		By("Ensuring machine is labelled with correct instance state")
		Eventually(func() bool {
			m := &infrav1.AWSMachine{}
			key := types.NamespacedName{
				Namespace: failingMachineMeta.Namespace,
				Name:      failingMachineMeta.Name,
			}
			Expect(k8sClient.Get(context.TODO(), key, m)).NotTo(HaveOccurred())
			labels := m.GetLabels()
			val := labels[Ec2InstanceStateLabelKey]
			return val == "shutting-down"
		}, 10*time.Second).Should(Equal(true))
	})
})

const messageBodyJSON = `{
	"source": "aws.ec2",
	"detail-type": "EC2 Instance State-change Notification",
	"detail": {
		"instance-id": "i-failing-instance-1",
		"state": "shutting-down"
	}
}`
