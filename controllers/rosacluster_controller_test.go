/*
Copyright 2025 The Kubernetes Authors.

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

// Package controllers provides a way to reconcile ROSA resources.
package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsiface "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

func TestRosaClusterReconcile(t *testing.T) {
	t.Run("Reconcile Rosa Cluster", func(t *testing.T) {
		g := NewWithT(t)
		ns, err := testEnv.CreateNamespace(ctx, "test-namespace")
		g.Expect(err).ToNot(HaveOccurred())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rosa-secret",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"ocmToken": []byte("secret-ocm-token-string"),
			},
		}

		identity := &infrav1.AWSClusterControllerIdentity{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: infrav1.AWSClusterControllerIdentitySpec{
				AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
					AllowedNamespaces: &infrav1.AllowedNamespaces{},
				},
			},
		}
		identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterStaticIdentity"))

		rosaClusterName := "rosa-controlplane-1"
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rosaClusterName,
				Namespace: ns.Name,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSAControlPlane",
				APIVersion: rosacontrolplanev1.GroupVersion.String(),
			},
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				RosaClusterName:   rosaClusterName,
				Subnets:           []string{"subnet-0ac99a6230b408813", "subnet-1ac99a6230b408811"},
				AvailabilityZones: []string{"az-1", "az-2"},
				Network: &rosacontrolplanev1.NetworkSpec{
					MachineCIDR: "10.0.0.0/16",
					PodCIDR:     "10.128.0.0/14",
					ServiceCIDR: "172.30.0.0/16",
				},
				Region:       "us-east-1",
				Version:      "4.19.20",
				ChannelGroup: "stable",
				RolesRef: rosacontrolplanev1.AWSRolesRef{
					IngressARN:              "ingress-arn",
					ImageRegistryARN:        "image-arn",
					StorageARN:              "storage-arn",
					NetworkARN:              "net-arn",
					KubeCloudControllerARN:  "kube-arn",
					NodePoolManagementARN:   "node-arn",
					ControlPlaneOperatorARN: "control-arn",
					KMSProviderARN:          "kms-arn",
				},
				OIDCID:           "oidcid1",
				InstallerRoleARN: "arn1",
				WorkerRoleARN:    "arn2",
				SupportRoleARN:   "arn3",
				CredentialsSecretRef: &corev1.LocalObjectReference{
					Name: secret.Name,
				},
				VersionGate: "Acknowledge",
				IdentityRef: &infrav1.AWSIdentityReference{
					Name: identity.Name,
					Kind: infrav1.ControllerIdentityKind,
				},
			},
		}

		rosaCluster := &expinfrav1.ROSACluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSACluster",
				APIVersion: expinfrav1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rosa-cluster",
				Namespace: ns.Name,
			},
		}

		capiCluster := &clusterv1.Cluster{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Cluster",
				APIVersion: clusterv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "capi-cluster-1",
				Namespace: ns.Name,
				UID:       types.UID("capi-cluster-1"),
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					Name:       rosaCluster.Name,
					Kind:       "ROSACluster",
					APIVersion: expinfrav1.GroupVersion.String(),
					Namespace:  ns.Name,
				},
				ControlPlaneRef: &corev1.ObjectReference{
					Name:       rosaControlPlane.Name,
					Kind:       "ROSAControlPlane",
					APIVersion: rosacontrolplanev1.GroupVersion.String(),
					Namespace:  ns.Name,
				},
				Paused: false,
			},
		}

		rosaCluster.OwnerReferences = []metav1.OwnerReference{
			{
				Name:       capiCluster.Name,
				Kind:       "Cluster",
				APIVersion: clusterv1.GroupVersion.String(),
				UID:        capiCluster.UID,
			},
		}

		createObject(g, secret, ns.Name)
		createObject(g, identity, ns.Name)
		createObject(g, capiCluster, ns.Name)
		createObject(g, rosaControlPlane, ns.Name)
		createObject(g, rosaCluster, ns.Name)

		// set controlplane status
		rosaCPPatch, err := patch.NewHelper(rosaControlPlane, testEnv)
		rosaControlPlane.Status.Ready = true
		rosaControlPlane.Status.Version = "4.19.20"
		rosaControlPlane.Status.ID = rosaClusterName
		g.Expect(rosaCPPatch.Patch(ctx, rosaControlPlane)).To(Succeed())
		g.Expect(err).ShouldNot(HaveOccurred())

		// set rosaCluster pause conditions
		rosaClsPatch, err := patch.NewHelper(rosaCluster, testEnv)
		rosaCluster.Status.Conditions = clusterv1.Conditions{
			clusterv1.Condition{
				Type:    clusterv1.PausedV1Beta2Condition,
				Status:  corev1.ConditionFalse,
				Reason:  clusterv1.NotPausedV1Beta2Reason,
				Message: "",
			},
		}
		g.Expect(rosaClsPatch.Patch(ctx, rosaCluster)).To(Succeed())
		g.Expect(err).ShouldNot(HaveOccurred())

		// set capiCluster pause condition
		clsPatch, err := patch.NewHelper(capiCluster, testEnv)
		capiCluster.Status.Conditions = clusterv1.Conditions{
			clusterv1.Condition{
				Type:    clusterv1.PausedV1Beta2Condition,
				Status:  corev1.ConditionFalse,
				Reason:  clusterv1.NotPausedV1Beta2Reason,
				Message: "",
			},
		}
		g.Expect(clsPatch.Patch(ctx, capiCluster)).To(Succeed())
		g.Expect(err).ShouldNot(HaveOccurred())

		// patching is not reliably synchronous
		time.Sleep(50 * time.Millisecond)

		mockCtrl := gomock.NewController(t)
		recorder := record.NewFakeRecorder(10)
		ctx := context.TODO()
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
		stsMock.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).AnyTimes()

		nodePoolName := "nodepool-1"
		expect := func(m *mocks.MockOCMClientMockRecorder) {
			m.GetNodePools(gomock.Any()).AnyTimes().DoAndReturn(func(clusterId string) ([]*cmv1.NodePool, error) {
				// Build a NodePool.
				builder := cmv1.NewNodePool().
					ID(nodePoolName).
					Version(cmv1.NewVersion().ID("openshift-v4.15.0")).
					AvailabilityZone("us-east-1a").
					Subnet("subnet-12345").
					Labels(map[string]string{"role": "worker"}).
					AutoRepair(true).
					TuningConfigs("tuning1").
					AWSNodePool(
						cmv1.NewAWSNodePool().
							InstanceType("m5.large").
							AdditionalSecurityGroupIds("sg-123", "sg-456").
							RootVolume(cmv1.NewAWSVolume().Size(120)),
					).
					Taints(
						cmv1.NewTaint().Key("dedicated").Value("gpu").Effect(string(corev1.TaintEffectNoSchedule)),
					).
					NodeDrainGracePeriod(
						cmv1.NewValue().Value(10),
					).
					ManagementUpgrade(
						cmv1.NewNodePoolManagementUpgrade().
							MaxSurge("1").
							MaxUnavailable("2"),
					).
					Replicas(2).
					Status(
						cmv1.NewNodePoolStatus().
							Message("").
							CurrentReplicas(2),
					)

				nodePool, err := builder.Build()
				g.Expect(err).ToNot(HaveOccurred())
				return []*cmv1.NodePool{nodePool}, err
			})
		}
		expect(ocmMock.EXPECT())

		r := ROSAClusterReconciler{
			Recorder:         recorder,
			WatchFilterValue: "",
			Client:           testEnv,
			NewStsClient: func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSClient {
				return stsMock
			},
			NewOCMClient: func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error) {
				return ocmMock, nil
			},
		}

		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: rosaCluster.Namespace,
				Name:      rosaCluster.Name,
			},
		}

		_, err = r.Reconcile(ctx, req)
		g.Expect(err).ToNot(HaveOccurred())

		// Check RosamachinePool & MachinePool are created.
		rosaMachinePool := &expinfrav1.ROSAMachinePool{}
		keyRosaMP := client.ObjectKey{Name: nodePoolName, Namespace: ns.Name}
		errRosaMP := testEnv.Get(ctx, keyRosaMP, rosaMachinePool)
		g.Expect(errRosaMP).ToNot(HaveOccurred())

		machinePool := &expclusterv1.MachinePool{}
		keyMP := client.ObjectKey{Name: nodePoolName, Namespace: ns.Name}
		errMP := testEnv.Get(ctx, keyMP, machinePool)
		g.Expect(errMP).ToNot(HaveOccurred())

		// Test get RosaMachinePoolNames
		rosaMachinePools, err := r.getRosaMachinePoolNames(ctx, capiCluster)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(len(rosaMachinePools)).To(Equal(1))

		cleanupObject(g, rosaMachinePool)
		cleanupObject(g, machinePool)
		cleanupObject(g, rosaCluster)
		cleanupObject(g, rosaControlPlane)
		cleanupObject(g, capiCluster)
		mockCtrl.Finish()
	})
}
