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

package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

const testAutoscalerClusterID = "cluster-123"

func newAutoscalerTestScope(spec *rosacontrolplanev1.ClusterAutoscaler) *scope.ROSAControlPlaneScope {
	return &scope.ROSAControlPlaneScope{
		ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{ClusterAutoscaler: spec},
		},
		Logger: *logger.FromContext(context.TODO()),
	}
}

func newAutoscalerTestReconciler(client rosa.ClusterAutoscalerClient) *ROSAControlPlaneReconciler {
	return &ROSAControlPlaneReconciler{
		NewClusterAutoscalerClient: func(_ context.Context, _ *scope.ROSAControlPlaneScope) (rosa.ClusterAutoscalerClient, error) {
			return client, nil
		},
	}
}

func newTestOCMAutoscaler(g *WithT, maxNodesTotal int) *v1.ClusterAutoscaler {
	autoscaler, err := v1.NewClusterAutoscaler().
		ResourceLimits(v1.NewAutoscalerResourceLimits().MaxNodesTotal(maxNodesTotal)).
		MaxPodGracePeriod(600).
		MaxNodeProvisionTime("15m").
		PodPriorityThreshold(-10).
		Build()
	g.Expect(err).ToNot(HaveOccurred())
	return autoscaler
}

func newTestCluster(g *WithT) *v1.Cluster {
	cluster, err := v1.NewCluster().ID(testAutoscalerClusterID).Build()
	g.Expect(err).ToNot(HaveOccurred())
	return cluster
}

func TestReconcileClusterAutoscaler(t *testing.T) {
	t.Run("not managed: populates status and never patches", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(newTestOCMAutoscaler(g, 50), nil)
		autoscalerMock.EXPECT().UpdateClusterAutoscaler(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(nil)
		v1beta1conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rosaScope.ControlPlane.Status.ClusterAutoscaler.ResourceLimits.MaxNodesTotal).To(HaveValue(Equal(int32(50))))
		g.Expect(v1beta1conditions.Get(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).To(BeNil())
	})

	t.Run("not managed: a read failure is not returned and keeps the previous status", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(nil, &rosa.OCMAPIError{Status: http.StatusServiceUnavailable})
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(nil)
		previous := &rosacontrolplanev1.ClusterAutoscalerStatus{MaxPodGracePeriod: ptr.To[int32](600)}
		rosaScope.ControlPlane.Status.ClusterAutoscaler = previous

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rosaScope.ControlPlane.Status.ClusterAutoscaler).To(BeIdenticalTo(previous))
	})

	t.Run("managed: no PATCH when OCM already matches", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(newTestOCMAutoscaler(g, 50), nil)
		autoscalerMock.EXPECT().UpdateClusterAutoscaler(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](50)},
		})

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v1beta1conditions.IsTrue(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).To(BeTrue())
	})

	t.Run("managed: PATCHes only the differing field and reports the updated values", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(newTestOCMAutoscaler(g, 0), nil)
		autoscalerMock.EXPECT().UpdateClusterAutoscaler(gomock.Any(), testAutoscalerClusterID, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, patch *v1.ClusterAutoscaler) (*v1.ClusterAutoscaler, error) {
				g.Expect(patch.ResourceLimits().MaxNodesTotal()).To(Equal(50))
				_, sent := patch.GetMaxPodGracePeriod()
				g.Expect(sent).To(BeFalse())
				return newTestOCMAutoscaler(g, 50), nil
			})
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](50)},
		})

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v1beta1conditions.IsTrue(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).To(BeTrue())
		g.Expect(rosaScope.ControlPlane.Status.ClusterAutoscaler.ResourceLimits.MaxNodesTotal).To(HaveValue(Equal(int32(50))))
	})

	t.Run("managed: a 400 from OCM is reported in the condition and not retried", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(newTestOCMAutoscaler(g, 0), nil)
		autoscalerMock.EXPECT().UpdateClusterAutoscaler(gomock.Any(), testAutoscalerClusterID, gomock.Any()).
			Return(nil, &rosa.OCMAPIError{Status: http.StatusBadRequest, Reason: "cluster size exceeds the limit", OperationID: "op-1"})
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](1000)},
		})

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		condition := v1beta1conditions.Get(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)
		g.Expect(condition).ToNot(BeNil())
		g.Expect(condition.Status).To(BeEquivalentTo("False"))
		g.Expect(condition.Reason).To(Equal(rosacontrolplanev1.ClusterAutoscalerValidationFailedReason))
		g.Expect(condition.Severity).To(Equal(clusterv1beta1.ConditionSeverityError))
		g.Expect(condition.Message).To(ContainSubstring("cluster size exceeds the limit"))
		g.Expect(condition.Message).To(ContainSubstring("op-1"))
		g.Expect(rosaScope.ControlPlane.Status.ClusterAutoscaler.ResourceLimits.MaxNodesTotal).To(HaveValue(Equal(int32(0))))
	})

	t.Run("managed: a 404 from OCM is reported as not available and not retried", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).
			Return(nil, &rosa.OCMAPIError{Status: http.StatusNotFound, Reason: "Autoscaler for cluster ID 'cluster-123' is not found"})
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{})

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v1beta1conditions.GetReason(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).
			To(Equal(rosacontrolplanev1.ClusterAutoscalerNotAvailableReason))
	})

	t.Run("managed: a 5xx from OCM is returned to be retried", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		autoscalerMock := mocks.NewMockClusterAutoscalerClient(mockCtrl)
		autoscalerMock.EXPECT().GetClusterAutoscaler(gomock.Any(), testAutoscalerClusterID).Return(newTestOCMAutoscaler(g, 0), nil)
		autoscalerMock.EXPECT().UpdateClusterAutoscaler(gomock.Any(), testAutoscalerClusterID, gomock.Any()).
			Return(nil, &rosa.OCMAPIError{Status: http.StatusServiceUnavailable})
		autoscalerMock.EXPECT().Close().Return(nil)

		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To[int32](50)},
		})

		err := newAutoscalerTestReconciler(autoscalerMock).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).To(HaveOccurred())
		g.Expect(v1beta1conditions.GetReason(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).
			To(Equal(rosacontrolplanev1.ReconciliationFailedReason))
	})

	t.Run("managed: failing to create the client is returned to be retried", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &ROSAControlPlaneReconciler{
			NewClusterAutoscalerClient: func(_ context.Context, _ *scope.ROSAControlPlaneScope) (rosa.ClusterAutoscalerClient, error) {
				return nil, errors.New("missing OCM credentials")
			},
		}
		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{})

		err := reconciler.reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).To(MatchError(ContainSubstring("missing OCM credentials")))
	})

	t.Run("managed: a nil client factory is returned as an error", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{})

		err := (&ROSAControlPlaneReconciler{}).reconcileClusterAutoscaler(context.TODO(), rosaScope, newTestCluster(g))
		g.Expect(err).To(HaveOccurred())
	})
}

func TestMarkClusterAutoscalerWaitingForCluster(t *testing.T) {
	t.Run("sets the condition when the autoscaler is managed", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := newAutoscalerTestScope(&rosacontrolplanev1.ClusterAutoscaler{})

		markClusterAutoscalerWaitingForCluster(rosaScope, v1.ClusterStateInstalling)
		condition := v1beta1conditions.Get(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)
		g.Expect(condition).ToNot(BeNil())
		g.Expect(condition.Reason).To(Equal(rosacontrolplanev1.ClusterAutoscalerWaitingForClusterReason))
		g.Expect(condition.Severity).To(Equal(clusterv1beta1.ConditionSeverityInfo))
		g.Expect(condition.Message).To(ContainSubstring("installing"))
	})

	t.Run("does nothing when the autoscaler is not managed", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := newAutoscalerTestScope(nil)

		markClusterAutoscalerWaitingForCluster(rosaScope, v1.ClusterStateInstalling)
		g.Expect(v1beta1conditions.Get(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)).To(BeNil())
	})
}

// TestClusterAutoscalerCRDValidation checks the CEL and OpenAPI validation of spec.clusterAutoscaler
// against the generated CRD loaded in envtest.
func TestClusterAutoscalerCRDValidation(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, "test-autoscaler-validation")
	g.Expect(err).ToNot(HaveOccurred())

	newControlPlane := func(autoscaler *rosacontrolplanev1.ClusterAutoscaler) *rosacontrolplanev1.ROSAControlPlane {
		return &rosacontrolplanev1.ROSAControlPlane{
			ObjectMeta: metav1.ObjectMeta{Name: "rosa-autoscaler-validation", Namespace: ns.Name},
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				RosaClusterName:   "rosa-autoscaler-validation",
				Subnets:           []string{"subnet-0ac99a6230b408813", "subnet-1ac99a6230b408811"},
				AvailabilityZones: []string{"az-1", "az-2"},
				Region:            "us-east-1",
				Version:           "4.19.10",
				ChannelGroup:      "stable",
				RolesRef: rosacontrolplanev1.AWSRolesRef{
					IngressARN:              "op-arn1",
					ImageRegistryARN:        "op-arn2",
					StorageARN:              "op-arn3",
					NetworkARN:              "op-arn4",
					KubeCloudControllerARN:  "op-arn5",
					NodePoolManagementARN:   "op-arn6",
					ControlPlaneOperatorARN: "op-arn7",
					KMSProviderARN:          "op-arn8",
				},
				OIDCID:            "oidcid1",
				InstallerRoleARN:  "arn:aws:iam::123456789012:role/installer",
				WorkerRoleARN:     "arn:aws:iam::123456789012:role/worker",
				SupportRoleARN:    "arn:aws:iam::123456789012:role/support",
				VersionGate:       "Acknowledge",
				ClusterAutoscaler: autoscaler,
			},
		}
	}
	minutes := func(m int) *metav1.Duration { return &metav1.Duration{Duration: time.Duration(m) * time.Minute} }
	maxNodesTotal := func(n int32) *rosacontrolplanev1.ClusterAutoscaler {
		return &rosacontrolplanev1.ClusterAutoscaler{
			ResourceLimits: &rosacontrolplanev1.ClusterAutoscalerResourceLimits{MaxNodesTotal: ptr.To(n)},
		}
	}

	tests := []struct {
		name       string
		autoscaler *rosacontrolplanev1.ClusterAutoscaler
		wantErr    string
	}{
		{name: "unset", autoscaler: nil},
		{name: "empty object", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{}},
		{name: "maxNodesTotal 0 (no limit)", autoscaler: maxNodesTotal(0)},
		{name: "maxNodesTotal 2", autoscaler: maxNodesTotal(2)},
		{name: "maxNodesTotal 1", autoscaler: maxNodesTotal(1), wantErr: "maxNodesTotal must be 0 or at least 2"},
		{name: "maxNodesTotal negative", autoscaler: maxNodesTotal(-1), wantErr: "should be greater than or equal to 0"},
		{name: "maxPodGracePeriod 0", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxPodGracePeriod: ptr.To[int32](0)}},
		{name: "maxPodGracePeriod 600", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxPodGracePeriod: ptr.To[int32](600)}},
		{name: "maxPodGracePeriod 300", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxPodGracePeriod: ptr.To[int32](300)}, wantErr: "maxPodGracePeriod must be 0 or at least 600"},
		{name: "maxNodeProvisionTime 15m", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxNodeProvisionTime: minutes(15)}},
		{name: "maxNodeProvisionTime 60m", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxNodeProvisionTime: minutes(60)}},
		{name: "maxNodeProvisionTime 10m", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxNodeProvisionTime: minutes(10)}, wantErr: "maxNodeProvisionTime must be between 15m and 60m"},
		{name: "maxNodeProvisionTime 61m", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{MaxNodeProvisionTime: minutes(61)}, wantErr: "maxNodeProvisionTime must be between 15m and 60m"},
		{name: "podPriorityThreshold -10", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{PodPriorityThreshold: ptr.To[int32](-10)}},
		{name: "podPriorityThreshold max", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{PodPriorityThreshold: ptr.To[int32](1000000000)}},
		{name: "podPriorityThreshold above max", autoscaler: &rosacontrolplanev1.ClusterAutoscaler{PodPriorityThreshold: ptr.To[int32](1000000001)}, wantErr: "should be less than or equal to 1000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			err := testEnv.Create(ctx, newControlPlane(tt.autoscaler), client.DryRunAll)
			if tt.wantErr == "" {
				g.Expect(err).ToNot(HaveOccurred())
				return
			}
			g.Expect(err).To(MatchError(ContainSubstring(tt.wantErr)))
		})
	}
}
