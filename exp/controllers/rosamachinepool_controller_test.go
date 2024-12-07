package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
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
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
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
	}

	machinePoolSpec := expclusterv1.MachinePoolSpec{
		Replicas: ptr.To[int32](2),
	}

	nodePoolBuilder := nodePoolBuilder(rosaMachinePoolSpec, machinePoolSpec)
	nodePoolSpec, err := nodePoolBuilder.Build()
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(computeSpecDiff(rosaMachinePoolSpec, nodePoolSpec)).To(BeEmpty())
}

func TestRosaMachinePoolReconcile(t *testing.T) {
	g := NewWithT(t)
	var (
		recorder         *record.FakeRecorder
		mockCtrl         *gomock.Controller
		ctx              context.Context
		scheme           *runtime.Scheme
		ns               *corev1.Namespace
		secret           *corev1.Secret
		rosaControlPlane *rosacontrolplanev1.ROSAControlPlane
		ownerCluster     *clusterv1.Cluster
		ownerMachinePool *expclusterv1.MachinePool
		rosaMachinePool  *expinfrav1.ROSAMachinePool
		ocmMock          *mocks.MockOCMClient
		objects          []client.Object
		err              error
	)

	setup := func(t *testing.T) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		recorder = record.NewFakeRecorder(10)
		ctx = context.TODO()
		scheme = runtime.NewScheme()
		// scheme = testEnv.Scheme()
		ns, err = testEnv.CreateNamespace(ctx, "test-namespace")
		g.Expect(err).To(BeNil())

		g.Expect(expinfrav1.AddToScheme(scheme)).To(Succeed())
		g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
		g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())
		g.Expect(expclusterv1.AddToScheme(scheme)).To(Succeed())
		g.Expect(rosacontrolplanev1.AddToScheme(scheme)).To(Succeed())
		g.Expect(corev1.AddToScheme(scheme)).To(Succeed())

		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rosa-secret",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"ocmToken": []byte("secret-ocm-token-string"),
			},
		}

		rosaControlPlane = &rosacontrolplanev1.ROSAControlPlane{
			ObjectMeta: metav1.ObjectMeta{Name: "rosa-control-plane", Namespace: ns.Name},
			TypeMeta: metav1.TypeMeta{
				Kind:       "ROSAControlPlane",
				APIVersion: rosacontrolplanev1.GroupVersion.String(),
			},
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				RosaClusterName:   "rosa-control-plane",
				Subnets:           []string{"subnet-0ac99a6230b408813", "subnet-1ac99a6230b408811"},
				AvailabilityZones: []string{"az-1", "az-2"},
				Network: &rosacontrolplanev1.NetworkSpec{
					MachineCIDR: "10.0.0.0/16",
					PodCIDR:     "10.128.0.0/14",
					ServiceCIDR: "172.30.0.0/16",
				},
				Region:           "us-east-1",
				Version:          "4.15.20",
				RolesRef:         rosacontrolplanev1.AWSRolesRef{},
				OIDCID:           "iodcid1",
				InstallerRoleARN: "arn1",
				WorkerRoleARN:    "arn2",
				SupportRoleARN:   "arn3",
				CredentialsSecretRef: &corev1.LocalObjectReference{
					Name: secret.Name,
				},
			},
			Status: rosacontrolplanev1.RosaControlPlaneStatus{
				Ready: true,
				ID:    "rosa-control-plane1",
			},
		}

		ownerCluster = &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "capi-test-6",
				Namespace: ns.Name,
			},
			Spec: clusterv1.ClusterSpec{
				ControlPlaneRef: &corev1.ObjectReference{
					Name:       rosaControlPlane.Name,
					Kind:       "ROSAControlPlane",
					APIVersion: rosacontrolplanev1.GroupVersion.String(),
				},
			},
		}

		rosaMachinePool = &expinfrav1.ROSAMachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rosa-machinepool",
				Namespace: ns.Name,
				UID:       "rosa-machinepool-1",
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

		ownerMachinePool = &expclusterv1.MachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "machinepool-1",
				Namespace: ns.Name,
				Labels:    map[string]string{clusterv1.ClusterNameLabel: ownerCluster.Name},
				UID:       "owner-mp-uid-1",
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "MachinePool",
				APIVersion: clusterv1.GroupVersion.String(),
			},
			Spec: expclusterv1.MachinePoolSpec{
				ClusterName: ownerCluster.Name,
				Template: clusterv1.MachineTemplateSpec{
					Spec: clusterv1.MachineSpec{
						ClusterName: ownerCluster.Name,
						InfrastructureRef: corev1.ObjectReference{
							UID:        rosaMachinePool.UID,
							Name:       rosaMachinePool.Name,
							Namespace:  ns.Namespace,
							Kind:       "ROSAMachinePool",
							APIVersion: expclusterv1.GroupVersion.String(),
						},
					},
				},
			},
		}

		// This is set by CAPI MachinePool reconcile
		rosaMachinePool.OwnerReferences = []metav1.OwnerReference{
			{
				Name:       ownerMachinePool.Name,
				UID:        ownerMachinePool.UID,
				Kind:       "MachinePool",
				APIVersion: clusterv1.GroupVersion.String(),
			},
		}

		objects = []client.Object{secret, ownerCluster, ownerMachinePool, rosaMachinePool}

		for _, obj := range objects {
			createObject(g, obj, ns.Name)
		}
	}

	teardown := func() {
		mockCtrl.Finish()
		for _, obj := range objects {
			cleanupObject(g, obj)
		}
	}

	t.Run("Reconcile create node pool", func(t *testing.T) {
		fmt.Println("START test")
		setup(t)
		defer teardown()
		ocmMock = mocks.NewMockOCMClient(mockCtrl)
		expect := func(m *mocks.MockOCMClientMockRecorder) {
			m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
				return nil, false, nil
			}).Times(1)
			m.CreateNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
				fmt.Println("NODE POOL", nodePool.ID())
				return nodePool, nil
			}).Times(1)
		}
		expect(ocmMock.EXPECT())

		fmt.Println("REC", scheme.Recognizes(rosaMachinePool.GroupVersionKind()))
		fmt.Println("REC", scheme.Recognizes(ownerCluster.GroupVersionKind()))

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(rosaMachinePool, ownerCluster, ownerMachinePool, rosaControlPlane, secret).Build()

		// failureMessage := "fail blabla"
		// rosaMachinePool.Status.FailureMessage = &failureMessage

		// err = c.Status().Patch(ctx, rosaMachinePool, client.MergeFrom(rosaMachinePool))
		// err = c.Status().Patch(ctx, ownerMachinePool, client.MergeFrom(ownerMachinePool))
		g.Expect(err).NotTo(HaveOccurred())
		mpPh, err := patch.NewHelper(rosaMachinePool, testEnv)

		rosaMachinePool.Status.Ready = true
		l := map[string]string{"key": "value"}
		rosaMachinePool.SetLabels(l)

		err = mpPh.Patch(ctx, rosaMachinePool, patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.RosaMachinePoolReadyCondition,
		}})

		m := &expinfrav1.ROSAMachinePool{}
		key := client.ObjectKey{Name: rosaMachinePool.Name, Namespace: ns.Name}
		c.Get(ctx, key, m)
		fmt.Println("ROSAMACHINEPOOL:", m.Name, m.Status, m.Namespace, m.Labels)

		g.Expect(err).NotTo(HaveOccurred())

		stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
		stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Times(1)

		r := ROSAMachinePoolReconciler{
			Recorder:         recorder,
			WatchFilterValue: "",
			Endpoints:        []scope.ServiceEndpoint{},
			Client:           c,
			NewStsClient:     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI { return stsMock },
			NewOCMClient: func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error) {
				return ocmMock, nil
			},
		}

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: "rosa-machinepool", Namespace: ns.Name}

		g.Expect(r.Endpoints).To(Equal("a"))
		// result, err := r.Reconcile(ctx, req)
		// key := client.ObjectKey{Name: ownerMachinePool.Name, Namespace: ns.Name}

		g.Expect(m.Status).To(Equal(map[string]string{"cluster.x-k8s.io/replicas-managed-by": "rosa"}))

		g.Expect(err).ToNot(HaveOccurred())

		// g.Expect(result).To(Equal(ctrl.Result{}))
	})

	// t.Run("Reconcile delete", func(t *testing.T) {
	// 	setup(t)
	// 	defer teardown()

	// 	deleteTime := metav1.NewTime(time.Now().Add(5 * time.Second))
	// 	rosaMachinePool.ObjectMeta.Finalizers = []string{"finalizer-rosa"}
	// 	rosaMachinePool.ObjectMeta.DeletionTimestamp = &deleteTime

	// 	ocmMock := mocks.NewMockOCMClient(mockCtrl)
	// 	expect := func(m *mocks.MockOCMClientMockRecorder) {
	// 		m.GetNodePool(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (*cmv1.NodePool, bool, error) {
	// 			nodePoolBuilder := nodePoolBuilder(rosaMachinePool.Spec, ownerMachinePool.Spec)
	// 			nodePool, err := nodePoolBuilder.ID("node-pool-1").Build()
	// 			g.Expect(err).To(BeNil())
	// 			return nodePool, true, nil
	// 		}).Times(1)
	// 		m.DeleteNodePool("rosa-control-plane1", "node-pool-1").DoAndReturn(func(clusterId string, nodePoolID string) error {
	// 			return nil
	// 		}).Times(1)
	// 	}
	// 	expect(ocmMock.EXPECT())

	// 	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(rosaMachinePool, ownerCluster, ownerMachinePool, rosaControlPlane, secret).Build()
	// 	stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
	// 	stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Times(1)

	// 	r := ROSAMachinePoolReconciler{
	// 		Recorder:         recorder,
	// 		WatchFilterValue: "",
	// 		Endpoints:        []scope.ServiceEndpoint{},
	// 		Client:           client,
	// 		NewStsClient:     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI { return stsMock },
	// 		NewOCMClient: func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error) {
	// 			return ocmMock, nil
	// 		},
	// 	}

	// 	req := ctrl.Request{}
	// 	req.NamespacedName = types.NamespacedName{Name: "rosa-machinepool", Namespace: ns.Name}

	// 	result, err := r.Reconcile(ctx, req)
	// 	g.Expect(err).ToNot(HaveOccurred())
	// 	g.Expect(result).To(Equal(ctrl.Result{}))
	// })
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
