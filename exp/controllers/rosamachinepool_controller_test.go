package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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

func TestNodePoolToRosaMachinePoolSpec(t *testing.T) {
	g := NewWithT(t)

	rosaMachinePoolSpec := expinfrav1.RosaMachinePoolSpec{
		NodePoolName:  "test-nodepool",
		Version:       "4.14.5",
		Subnet:        "subnet-id",
		AutoRepair:    true,
		InstanceType:  "m5.large",
		TuningConfigs: []string{"config1"},
		VolumeSize:    199,
		NodeDrainGracePeriod: &metav1.Duration{
			Duration: time.Minute * 10,
		},
		UpdateConfig: &expinfrav1.RosaUpdateConfig{
			RollingUpdate: &expinfrav1.RollingUpdate{
				MaxSurge:       ptr.To(intstr.FromInt32(3)),
				MaxUnavailable: ptr.To(intstr.FromInt32(5)),
			},
		},
		AdditionalSecurityGroups: []string{
			"id-1",
			"id-2",
		},
		Labels: map[string]string{
			"label1": "value1",
			"label2": "value2",
		},
		Taints: []expinfrav1.RosaTaint{
			{
				Key:    "myKey",
				Value:  "myValue",
				Effect: corev1.TaintEffectNoExecute,
			},
		},
		CapacityReservationID: "capacity-reservation-id",
	}

	machinePoolSpec := expclusterv1.MachinePoolSpec{
		Replicas: ptr.To[int32](2),
	}

	nodePoolBuilder := nodePoolBuilder(rosaMachinePoolSpec, machinePoolSpec, rosacontrolplanev1.Stable)
	nodePoolSpec, err := nodePoolBuilder.Build()
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(computeSpecDiff(rosaMachinePoolSpec, nodePoolSpec)).To(BeEmpty())
}

func TestRosaMachinePoolReconcile(t *testing.T) {
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

	rosaControlPlane := func(i int) *rosacontrolplanev1.ROSAControlPlane {
		return &rosacontrolplanev1.ROSAControlPlane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("rosa-control-plane-%v", i),
				Namespace: ns.Name},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSAControlPlane",
				APIVersion: rosacontrolplanev1.GroupVersion.String(),
			},
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				RosaClusterName:   fmt.Sprintf("rosa-control-plane-%v", i),
				Subnets:           []string{"subnet-0ac99a6230b408813", "subnet-1ac99a6230b408811"},
				AvailabilityZones: []string{"az-1", "az-2"},
				Network: &rosacontrolplanev1.NetworkSpec{
					MachineCIDR: "10.0.0.0/16",
					PodCIDR:     "10.128.0.0/14",
					ServiceCIDR: "172.30.0.0/16",
				},
				Region:       "us-east-1",
				Version:      "4.15.20",
				ChannelGroup: "stable",
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
				OIDCID:           "iodcid1",
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
			Status: rosacontrolplanev1.RosaControlPlaneStatus{
				Ready:   true,
				ID:      fmt.Sprintf("rosa-control-plane-%v", i),
				Version: "4.15.20",
			},
		}
	}

	ownerCluster := func(i int) *clusterv1.Cluster {
		return &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("owner-cluster-%v", i),
				Namespace: ns.Name,
			},
			Spec: clusterv1.ClusterSpec{
				ControlPlaneRef: &corev1.ObjectReference{
					Name:       rosaControlPlane(i).Name,
					Kind:       "ROSAControlPlane",
					APIVersion: rosacontrolplanev1.GroupVersion.String(),
				},
			},
		}
	}

	rosaMachinePool := func(i int) *expinfrav1.ROSAMachinePool {
		return &expinfrav1.ROSAMachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("rosa-machinepool-%v", i),
				Namespace: ns.Name,
				UID:       types.UID(fmt.Sprintf("rosa-machinepool-%v", i)),
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSAMachinePool",
				APIVersion: expinfrav1.GroupVersion.String(),
			},
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-nodepool",
				Version:      "4.14.5",
				Subnet:       "subnet-id",
				InstanceType: "m5.large",
			},
		}
	}

	ownerMachinePool := func(i int) *expclusterv1.MachinePool {
		return &expclusterv1.MachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("machinepool-%v", i),
				Namespace: ns.Name,
				Labels:    map[string]string{clusterv1.ClusterNameLabel: ownerCluster(i).Name},
				UID:       types.UID(fmt.Sprintf("owner-mp-uid--%v", i)),
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "MachinePool",
				APIVersion: clusterv1.GroupVersion.String(),
			},
			Spec: expclusterv1.MachinePoolSpec{
				ClusterName: fmt.Sprintf("owner-cluster-%v", i),
				Template: clusterv1.MachineTemplateSpec{
					Spec: clusterv1.MachineSpec{
						ClusterName: fmt.Sprintf("owner-cluster-%v", i),
						InfrastructureRef: corev1.ObjectReference{
							UID:        rosaMachinePool(i).UID,
							Name:       rosaMachinePool(i).Name,
							Namespace:  ns.Namespace,
							Kind:       "ROSAMachinePool",
							APIVersion: expclusterv1.GroupVersion.String(),
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name               string
		newROSAMachinePool *expinfrav1.ROSAMachinePool
		oldROSAMachinePool *expinfrav1.ROSAMachinePool
		machinePool        *expclusterv1.MachinePool
		expect             func(m *mocks.MockOCMClientMockRecorder)
		result             reconcile.Result
	}{
		{
			name:               "create node pool, nodepool doesn't exist",
			machinePool:        ownerMachinePool(0),
			oldROSAMachinePool: rosaMachinePool(0),
			newROSAMachinePool: &expinfrav1.ROSAMachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rosa-machinepool",
					Namespace: ns.Name,
					UID:       "rosa-machinepool",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ROSAMachinePool",
					APIVersion: expinfrav1.GroupVersion.String(),
				},
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName: "test-nodepool",
					Version:      "4.14.5",
					Subnet:       "subnet-id",
					InstanceType: "m5.large",
				},
				Status: expinfrav1.RosaMachinePoolStatus{
					Ready: false,
					ID:    rosaMachinePool(0).Spec.NodePoolName,
				},
			},
			result: ctrl.Result{},
			expect: func(m *mocks.MockOCMClientMockRecorder) {
				m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
					return nil, false, nil
				}).Times(1)
				m.CreateNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
					return nodePool, nil
				}).Times(1)
			},
		},
		{
			name:               "Nodepool exist, but is not ready",
			machinePool:        ownerMachinePool(1),
			oldROSAMachinePool: rosaMachinePool(1),
			newROSAMachinePool: &expinfrav1.ROSAMachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rosa-machinepool",
					Namespace: ns.Name,
					UID:       "rosa-machinepool",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ROSAMachinePool",
					APIVersion: expinfrav1.GroupVersion.String(),
				},
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName: "test-nodepool",
					Version:      "4.14.5",
					Subnet:       "subnet-id",
					InstanceType: "m5.large",
				},
				Status: expinfrav1.RosaMachinePoolStatus{
					Ready:    false,
					Replicas: 0,
				},
			},
			result: ctrl.Result{RequeueAfter: time.Second * 60},
			expect: func(m *mocks.MockOCMClientMockRecorder) {
				m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
					nodePoolBuilder := nodePoolBuilder(rosaMachinePool(1).Spec, ownerMachinePool(1).Spec, rosacontrolplanev1.Stable)
					nodePool, err := nodePoolBuilder.ID("node-pool-1").Build()
					g.Expect(err).To(BeNil())
					return nodePool, true, nil
				}).Times(1)
				m.UpdateNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
					return nodePool, nil
				}).Times(1)
				m.CreateNodePool(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:               "Nodepool is ready",
			machinePool:        ownerMachinePool(2),
			oldROSAMachinePool: rosaMachinePool(2),
			newROSAMachinePool: &expinfrav1.ROSAMachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rosa-machinepool",
					Namespace: ns.Name,
					UID:       "rosa-machinepool",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ROSAMachinePool",
					APIVersion: expinfrav1.GroupVersion.String(),
				},
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName: "test-nodepool",
					Version:      "4.14.5",
					Subnet:       "subnet-id",
					InstanceType: "m5.large",
				},
				Status: expinfrav1.RosaMachinePoolStatus{
					Ready:    true,
					Replicas: 1,
				},
			},
			result: ctrl.Result{},
			expect: func(m *mocks.MockOCMClientMockRecorder) {
				m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
					nodePoolBuilder := nodePoolBuilder(rosaMachinePool(2).Spec, ownerMachinePool(2).Spec, rosacontrolplanev1.Stable)
					statusBuilder := (&cmv1.NodePoolStatusBuilder{}).CurrentReplicas(1)
					autoscalingBuilder := (&cmv1.NodePoolAutoscalingBuilder{}).MinReplica(1).MaxReplica(1)
					nodePool, err := nodePoolBuilder.ID("node-pool-1").Autoscaling(autoscalingBuilder).Replicas(1).Status(statusBuilder).Build()
					g.Expect(err).NotTo(HaveOccurred())

					return nodePool, true, nil
				}).Times(1)
				m.UpdateNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
					statusBuilder := (&cmv1.NodePoolStatusBuilder{}).CurrentReplicas(1)
					version := (&cmv1.VersionBuilder{}).RawID("4.14.5")
					npBuilder := cmv1.NodePoolBuilder{}
					updatedNodePool, err := npBuilder.Copy(nodePool).Status(statusBuilder).Version(version).Build()
					g.Expect(err).NotTo(HaveOccurred())

					return updatedNodePool, nil
				}).Times(1)
				m.CreateNodePool(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "Create nodepool, replicas are set in MachinePool",
			machinePool: &expclusterv1.MachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ownerMachinePool(3).Name,
					Namespace: ns.Name,
					Labels:    map[string]string{clusterv1.ClusterNameLabel: ownerCluster(3).Name},
					UID:       types.UID("owner-mp-uid--3"),
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachinePool",
					APIVersion: clusterv1.GroupVersion.String(),
				},
				Spec: expclusterv1.MachinePoolSpec{
					ClusterName: ownerCluster(3).Name,
					Replicas:    ptr.To[int32](2),
					Template: clusterv1.MachineTemplateSpec{
						Spec: clusterv1.MachineSpec{
							ClusterName: ownerCluster(3).Name,
							InfrastructureRef: corev1.ObjectReference{
								UID:        rosaMachinePool(3).UID,
								Name:       rosaMachinePool(3).Name,
								Namespace:  ns.Namespace,
								Kind:       "ROSAMachinePool",
								APIVersion: expclusterv1.GroupVersion.String(),
							},
						},
					},
				},
			},
			oldROSAMachinePool: rosaMachinePool(3),
			newROSAMachinePool: &expinfrav1.ROSAMachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rosa-machinepool",
					Namespace: ns.Name,
					UID:       "rosa-machinepool",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ROSAMachinePool",
					APIVersion: expinfrav1.GroupVersion.String(),
				},
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName: "test-nodepool",
					Version:      "4.14.5",
					Subnet:       "subnet-id",
					InstanceType: "m5.large",
				},
				Status: expinfrav1.RosaMachinePoolStatus{
					Ready: false,
					ID:    rosaMachinePool(3).Spec.NodePoolName,
				},
			},
			result: ctrl.Result{},
			expect: func(m *mocks.MockOCMClientMockRecorder) {
				m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
					return nil, false, nil
				}).Times(1)
				m.CreateNodePool(gomock.Any(), matchesReplicas(2)).DoAndReturn(func(clusterId string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
					return nodePool, nil
				}).Times(1)
			},
		},
		{
			name: "Update nodepool, replicas are updated from MachinePool",
			machinePool: &expclusterv1.MachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ownerMachinePool(4).Name,
					Namespace: ns.Name,
					Labels:    map[string]string{clusterv1.ClusterNameLabel: ownerCluster(4).Name},
					UID:       types.UID("owner-mp-uid--4"),
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachinePool",
					APIVersion: clusterv1.GroupVersion.String(),
				},
				Spec: expclusterv1.MachinePoolSpec{
					ClusterName: ownerCluster(4).Name,
					Replicas:    ptr.To[int32](2),
					Template: clusterv1.MachineTemplateSpec{
						Spec: clusterv1.MachineSpec{
							ClusterName: ownerCluster(4).Name,
							InfrastructureRef: corev1.ObjectReference{
								UID:        rosaMachinePool(4).UID,
								Name:       rosaMachinePool(4).Name,
								Namespace:  ns.Namespace,
								Kind:       "ROSAMachinePool",
								APIVersion: expclusterv1.GroupVersion.String(),
							},
						},
					},
				},
			},
			oldROSAMachinePool: rosaMachinePool(4),
			newROSAMachinePool: &expinfrav1.ROSAMachinePool{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rosa-machinepool",
					Namespace: ns.Name,
					UID:       "rosa-machinepool",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ROSAMachinePool",
					APIVersion: expinfrav1.GroupVersion.String(),
				},
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName: "test-nodepool",
					Version:      "4.14.5",
					Subnet:       "subnet-id",
					InstanceType: "m5.large",
				},
				Status: expinfrav1.RosaMachinePoolStatus{
					// Ready:    false,
					Ready:    true,
					Replicas: 2,
				},
			},
			result: ctrl.Result{},
			expect: func(m *mocks.MockOCMClientMockRecorder) {
				m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
					rosaSpec := rosaMachinePool(4).Spec
					rosaSpec.UpdateConfig = &expinfrav1.RosaUpdateConfig{
						RollingUpdate: &expinfrav1.RollingUpdate{
							MaxSurge:       ptr.To(intstr.FromInt32(1)),
							MaxUnavailable: ptr.To(intstr.FromInt32(0)),
						},
					}
					nodePoolBuilder := nodePoolBuilder(rosaSpec, ownerMachinePool(4).Spec, rosacontrolplanev1.Stable)
					statusBuilder := (&cmv1.NodePoolStatusBuilder{}).CurrentReplicas(1)
					nodeGrace := (&cmv1.ValueBuilder{}).Unit("s").Value(0)
					nodePool, err := nodePoolBuilder.ID("test-nodepool").Replicas(1).Status(statusBuilder).AutoRepair(true).NodeDrainGracePeriod(nodeGrace).Build()
					g.Expect(err).NotTo(HaveOccurred())

					return nodePool, true, nil
				}).Times(1)
				m.UpdateNodePool(gomock.Any(), matchesReplicas(2)).DoAndReturn(func(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
					statusBuilder := (&cmv1.NodePoolStatusBuilder{}).CurrentReplicas(2)
					version := (&cmv1.VersionBuilder{}).RawID("4.14.5")
					npBuilder := cmv1.NodePoolBuilder{}
					updatedNodePool, err := npBuilder.Copy(nodePool).Status(statusBuilder).Version(version).Build()
					g.Expect(err).NotTo(HaveOccurred())

					return updatedNodePool, nil
				}).Times(1)
				m.CreateNodePool(gomock.Any(), gomock.Any()).Times(0)
			},
		},
	}

	createObject(g, secret, ns.Name)
	createObject(g, identity, ns.Name)
	defer cleanupObject(g, secret)
	defer cleanupObject(g, identity)

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// This is set by CAPI MachinePool reconcile
			test.oldROSAMachinePool.OwnerReferences = []metav1.OwnerReference{
				{
					Name:       ownerMachinePool(i).Name,
					UID:        ownerMachinePool(i).UID,
					Kind:       "MachinePool",
					APIVersion: clusterv1.GroupVersion.String(),
				},
			}
			cp := rosaControlPlane(i)
			objects := []client.Object{ownerCluster(i), test.machinePool, cp, test.oldROSAMachinePool}

			for _, obj := range objects {
				createObject(g, obj, ns.Name)
			}
			// make Control Plane ready, can't do this duirng creation
			mpPh, err := patch.NewHelper(cp, testEnv)
			cp.Status.Ready = true
			cp.Status.Version = cp.Spec.Version
			g.Expect(mpPh.Patch(ctx, cp)).To(Succeed())
			g.Expect(err).ShouldNot(HaveOccurred())

			// patch status conditions
			rmpPh, err := patch.NewHelper(test.oldROSAMachinePool, testEnv)
			test.oldROSAMachinePool.Status.Conditions = clusterv1.Conditions{
				{
					Type:   "Paused",
					Status: corev1.ConditionFalse,
					Reason: "NotPaused",
				},
			}

			g.Expect(rmpPh.Patch(ctx, test.oldROSAMachinePool)).To(Succeed())
			g.Expect(err).ShouldNot(HaveOccurred())

			// patching is not reliably synchronous
			time.Sleep(50 * time.Millisecond)

			mockCtrl := gomock.NewController(t)
			recorder := record.NewFakeRecorder(10)
			ctx := context.TODO()
			ocmMock := mocks.NewMockOCMClient(mockCtrl)
			test.expect(ocmMock.EXPECT())

			stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
			stsMock.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Times(1)

			r := ROSAMachinePoolReconciler{
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

			req := ctrl.Request{}
			req.NamespacedName = types.NamespacedName{Name: test.oldROSAMachinePool.Name, Namespace: ns.Name}

			result, errReconcile := r.Reconcile(ctx, req)
			g.Expect(errReconcile).ToNot(HaveOccurred())
			g.Expect(result).To(Equal(test.result))
			time.Sleep(50 * time.Millisecond)

			m := &expinfrav1.ROSAMachinePool{}
			key := client.ObjectKey{Name: test.oldROSAMachinePool.Name, Namespace: ns.Name}
			errGet := testEnv.Get(ctx, key, m)
			g.Expect(errGet).NotTo(HaveOccurred())
			g.Expect(m.Status.Ready).To(Equal(test.newROSAMachinePool.Status.Ready))
			g.Expect(m.Status.Replicas).To(Equal(test.newROSAMachinePool.Status.Replicas))
			g.Expect(m.Status.ID).To(Equal(test.newROSAMachinePool.Status.ID))

			// cleanup
			for _, obj := range objects {
				cleanupObject(g, obj)
			}
			mockCtrl.Finish()
		})
	}

	t.Run("Reconcile delete", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl := gomock.NewController(t)
		recorder := record.NewFakeRecorder(10)
		ctx := context.TODO()
		controlPlaneName := "rosa-control-plane-9"
		mp := &expinfrav1.ROSAMachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:       controlPlaneName,
				Namespace:  ns.Name,
				UID:        types.UID(controlPlaneName),
				Finalizers: []string{expinfrav1.RosaMachinePoolFinalizer},
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSAMachinePool",
				APIVersion: expinfrav1.GroupVersion.String(),
			},
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-nodepool-1",
				Version:      "4.14.5",

				Subnet:       "subnet-id",
				InstanceType: "m5.large",
			},
		}
		oc := ownerCluster(9)
		omp := ownerMachinePool(9)
		cp := rosaControlPlane(9)
		objects := []client.Object{oc, omp, cp, mp}

		for _, obj := range objects {
			createObject(g, obj, ns.Name)
		}

		cpPh, err := patch.NewHelper(cp, testEnv)
		cp.Status.Ready = true
		cp.Status.ID = controlPlaneName
		cp.Status.Version = cp.Spec.Version
		g.Expect(cpPh.Patch(ctx, cp)).To(Succeed())
		g.Expect(err).ShouldNot(HaveOccurred())

		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		nodePoolName := "node-pool-1"
		expect := func(m *mocks.MockOCMClientMockRecorder) {
			m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
				nodePoolBuilder := nodePoolBuilder(mp.Spec, omp.Spec, rosacontrolplanev1.Stable)
				nodePool, err := nodePoolBuilder.ID(nodePoolName).Build()
				g.Expect(err).NotTo(HaveOccurred())
				return nodePool, true, nil
			}).Times(1)
			m.DeleteNodePool(controlPlaneName, nodePoolName).DoAndReturn(func(clusterId string, nodePoolID string) error {
				testEnv.Delete(ctx, mp)
				return nil
			}).Times(1)
		}
		expect(ocmMock.EXPECT())

		stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)
		stsMock.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Times(1)

		r := ROSAMachinePoolReconciler{
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

		log := logger.FromContext(ctx)
		machinePoolScope, err1 := scope.NewRosaMachinePoolScope(scope.RosaMachinePoolScopeParams{
			Client:          r.Client,
			ControllerName:  "rosamachinepool",
			Cluster:         oc,
			ControlPlane:    cp,
			MachinePool:     omp,
			RosaMachinePool: mp,
			Logger:          log,
		})
		g.Expect(err1).ToNot(HaveOccurred())

		rosaControlPlaneScope, err2 := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
			Client:         r.Client,
			Cluster:        oc,
			ControlPlane:   cp,
			ControllerName: "rosaControlPlane",
			NewStsClient:   r.NewStsClient,
		})
		g.Expect(err2).ToNot(HaveOccurred())

		err3 := r.reconcileDelete(ctx, machinePoolScope, rosaControlPlaneScope)
		g.Expect(err3).ToNot(HaveOccurred())

		machinePoolScope.Close()
		time.Sleep(50 * time.Millisecond)
		rosaMachinePool := &expinfrav1.ROSAMachinePool{}
		key := client.ObjectKey{Name: mp.Name, Namespace: ns.Name}
		err4 := testEnv.Get(ctx, key, rosaMachinePool)
		g.Expect(err4).To(HaveOccurred())
		g.Expect(rosaMachinePool.Finalizers).To(BeNil())

		for _, obj := range objects {
			cleanupObject(g, obj)
		}
		mockCtrl.Finish()
	})
}

func createObject(g *WithT, obj client.Object, namespace string) {
	if obj.DeepCopyObject() != nil {
		obj.SetNamespace(namespace)
		g.Expect(testEnv.Create(ctx, obj)).To(Succeed())
	}
}

func cleanupObject(g *WithT, obj client.Object) {
	if obj.DeepCopyObject() != nil {
		g.Expect(testEnv.Cleanup(ctx, obj)).To(Succeed())
	}
}

type replicasMatcher struct {
	replicas int
}

func (m replicasMatcher) Matches(arg interface{}) bool {
	sarg := arg.(*cmv1.NodePool)
	if sarg != nil && sarg.Replicas() == m.replicas {
		return true
	}
	return false
}

// Not used here, but satisfies the Matcher interface.
func (m replicasMatcher) String() string {
	return fmt.Sprintf("Replicas %v", m.replicas)
}

func matchesReplicas(replicas int) gomock.Matcher {
	return replicasMatcher{replicas: replicas}
}
