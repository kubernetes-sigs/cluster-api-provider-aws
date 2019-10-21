/*
Copyright 2019 The Kubernetes Authors.

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

package controllers

import (
	"bytes"
	"context"
	"flag"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope" //nolint
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/mock_services" //nolint
)

var _ = Describe("AWSMachineReconciler", func() {
	var (
		reconciler AWSMachineReconciler
		cs         *scope.ClusterScope
		ms         *scope.MachineScope
		mockCtrl   *gomock.Controller
		ec2Svc     *mock_services.MockEC2MachineInterface
	)

	BeforeEach(func() {
		// https://github.com/kubernetes/klog/issues/87#issuecomment-540153080
		// TODO: Replace with LogToOutput when https://github.com/kubernetes/klog/pull/99 merges
		flag.Set("logtostderr", "false")
		klog.SetOutput(GinkgoWriter)
		var err error

		ms, err = scope.NewMachineScope(
			scope.MachineScopeParams{
				Client: fake.NewFakeClient(),
				Cluster: &clusterv1.Cluster{
					Status: clusterv1.ClusterStatus{
						InfrastructureReady: true,
					},
				},
				Machine: &clusterv1.Machine{
					Spec: clusterv1.MachineSpec{
						Bootstrap: clusterv1.Bootstrap{
							Data: pointer.StringPtr("best pony: all of them"),
						},
					},
				},
				AWSCluster: &infrav1.AWSCluster{},
				AWSMachine: &infrav1.AWSMachine{},
			},
		)
		Expect(err).To(BeNil())

		cs, err = scope.NewClusterScope(
			scope.ClusterScopeParams{
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			},
		)
		Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(GinkgoT())
		ec2Svc = mock_services.NewMockEC2MachineInterface(mockCtrl)
		reconciler = AWSMachineReconciler{
			serviceFactory: func(*scope.ClusterScope) services.EC2MachineInterface {
				return ec2Svc
			},
		}

	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Reconciling an AWSMachine", func() {
		When("we can't reach amazon", func() {
			var expectedErr = errors.New("no connection available ")

			BeforeEach(func() {
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, expectedErr).AnyTimes()
			})

			It("should exit immediately on an error state", func() {
				er := capierrors.CreateMachineError
				ms.AWSMachine.Status.ErrorReason = &er
				ms.AWSMachine.Status.ErrorMessage = pointer.StringPtr("Couldn't create machine")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(err).To(BeNil())
				Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
			})

			It("should add our finalizer to the machine", func() {
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)

				Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
			})

			It("should exit immediately if cluster infra isn't ready", func() {
				ms.Cluster.Status.InfrastructureReady = false

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(err).To(BeNil())
				Expect(buf.String()).To(ContainSubstring("Cluster infrastructure is not ready yet"))
			})

			It("should exit immediately if bootstrap data isn't available", func() {
				ms.Machine.Spec.Bootstrap.Data = nil

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(err).To(BeNil())
				Expect(buf.String()).To(ContainSubstring("Bootstrap data is not yet available"))
			})

			It("should return an error when we can't list instances by tags", func() {

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})

		When("there's a provider ID", func() {
			var id = "aws:////myMachine"
			BeforeEach(func() {
				_, err := noderefutil.NewProviderID(id)
				Expect(err).To(BeNil())

				ms.AWSMachine.Spec.ProviderID = &id
			})

			It("it should look up by provider ID when one exists", func() {
				expectedErr := errors.New("no connection available ")
				ec2Svc.EXPECT().InstanceIfExists(PointsTo("myMachine")).Return(nil, expectedErr)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			It("should try to create a new machine if none exists", func() {
				expectedErr := errors.New("Invalid instance")
				ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
				ec2Svc.EXPECT().CreateInstance(gomock.Any()).Return(nil, expectedErr)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})

		When("instance creation succeeds", func() {
			var instance *infrav1.Instance
			BeforeEach(func() {
				instance = &infrav1.Instance{
					ID: "myMachine",
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				ec2Svc.EXPECT().CreateInstance(gomock.Any()).Return(instance, nil)
			})

			Context("instance security group errors", func() {
				BeforeEach(func() {
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(nil, errors.New("stop here"))
				})

				It("should set attributes after creating an instance", func() {
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)
					Expect(ms.AWSMachine.Spec.ProviderID).To(PointTo(Equal("aws:////myMachine")))
					Expect(ms.AWSMachine.Annotations).To(Equal(map[string]string{"cluster-api-provider-aws": "true"}))
				})

				Context("with captured logging", func() {
					var buf *bytes.Buffer

					BeforeEach(func() {
						buf = new(bytes.Buffer)
						klog.SetOutput(buf)
					})

					It("should set instance to running", func() {
						instance.State = infrav1.InstanceStateRunning
						_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)
						Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
						Expect(buf.String()).To(ContainSubstring(("Machine instance is running")))
					})

					It("should set instance to pending", func() {
						instance.State = infrav1.InstanceStatePending
						_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)
						Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStatePending)))
						Expect(buf.String()).To(ContainSubstring(("Machine instance is pending")))
					})
				})

				It("should set error message when instance status unknown", func() {
					instance.State = infrav1.InstanceStateStopping
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)
					Expect(ms.AWSMachine.Status.ErrorReason).To(PointTo(Equal(capierrors.UpdateMachineError)))
					Expect(ms.AWSMachine.Status.ErrorMessage).To(PointTo(Equal("EC2 instance state \"stopping\" is unexpected")))
				})

				// TODO: Mock out ELB API as well
				XIt("should call security ELB attachment when node is control plane", func() {
				})
			})

			Context("Security Groups succeed", func() {
				BeforeEach(func() {
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
						Return(map[string][]string{"eid": {}}, nil)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
				})

				It("should reconcile security groups", func() {
					ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1.AWSResourceReference{
						{
							ID: pointer.StringPtr("sg-2345"),
						},
					}
					// ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(instance.ID, []string{"sg-2345"})

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs)
				})

				It("should not tag anything if there's not tags", func() {
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Times(0)
					reconciler.reconcileNormal(context.Background(), ms, cs)
				})

				It("should tag instances from machine and cluster tags", func() {

					ms.AWSMachine.Spec.AdditionalTags = infrav1.Tags{"kind": "alicorn"}
					ms.AWSCluster.Spec.AdditionalTags = infrav1.Tags{"colour": "lavender"}

					ec2Svc.EXPECT().UpdateResourceTags(
						PointsTo("myMachine"),
						map[string]string{
							"kind":   "alicorn",
							"colour": "lavender",
						},
						map[string]string{},
					).Return(nil)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs)
					Expect(err).To(BeNil())
				})
			})
		})
	})
})

func PointsTo(s string) gomock.Matcher {
	return &pointsTo{
		val: s,
	}
}

type pointsTo struct {
	val string
}

func (p *pointsTo) Matches(x interface{}) bool {
	ptr, ok := x.(*string)
	if !ok {
		return false
	}

	if ptr == nil {
		return false
	}

	return *ptr == p.val
}

func (p *pointsTo) String() string {
	return fmt.Sprintf("Pointer to string %q", p.val)
}
