package machinehealthcheck

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/glog"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	"github.com/openshift/machine-api-operator/pkg/util/conditions"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("[Feature:MachineHealthCheck] MachineHealthCheck controller", func() {
	var client runtimeclient.Client
	var numberOfReadyWorkers int
	var workerNode *corev1.Node
	var workerMachine *mapiv1beta1.Machine

	stopKubeletAndValidateMachineDeletion := func(workerNodeName *corev1.Node, workerMachine *mapiv1beta1.Machine, timeout time.Duration) {
		By(fmt.Sprintf("Stopping kubelet service on the node %s", workerNode.Name))
		err := e2e.StopKubelet(workerNode.Name)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("Validating that node %s has 'NotReady' condition", workerNode.Name))
		waitForNodeUnhealthyCondition(workerNode.Name)

		By(fmt.Sprintf("Validating that machine %s is deleted", workerMachine.Name))
		machine := &mapiv1beta1.Machine{}
		key := types.NamespacedName{
			Namespace: workerMachine.Namespace,
			Name:      workerMachine.Name,
		}
		Eventually(func() bool {
			err := client.Get(context.TODO(), key, machine)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, timeout, 5*time.Second).Should(BeTrue())
	}

	BeforeEach(func() {
		var err error
		client, err = e2e.LoadClient()
		Expect(err).ToNot(HaveOccurred())

		// TODO: enable once https://github.com/openshift/cluster-api-actuator-pkg/pull/61 is fixed
		glog.V(2).Info("Skipping machine health checking test")
		Skip("Skipping machine health checking test")

		workerNodes, err := e2e.GetWorkerNodes(client)
		Expect(err).ToNot(HaveOccurred())

		readyWorkerNodes := e2e.FilterReadyNodes(workerNodes)
		Expect(readyWorkerNodes).ToNot(BeEmpty())

		numberOfReadyWorkers = len(readyWorkerNodes)
		workerNode = &readyWorkerNodes[0]
		glog.V(2).Infof("Worker node %s", workerNode.Name)

		workerMachine, err = e2e.GetMachineFromNode(client, workerNode)
		Expect(err).ToNot(HaveOccurred())
		glog.V(2).Infof("Worker machine %s", workerMachine.Name)

		glog.V(2).Infof("Create machine health check with label selector: %s", workerMachine.Labels)
		err = e2e.CreateMachineHealthCheck(workerMachine.Labels)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with node-unhealthy-conditions configmap", func() {
		BeforeEach(func() {
			unhealthyConditions := &conditions.UnhealthyConditions{
				Items: []conditions.UnhealthyCondition{
					{
						Name:    "Ready",
						Status:  "Unknown",
						Timeout: "60s",
					},
				},
			}
			glog.V(2).Infof("Create node-unhealthy-conditions configmap")
			err := e2e.CreateUnhealthyConditionsConfigMap(unhealthyConditions)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should delete unhealthy machine", func() {
			stopKubeletAndValidateMachineDeletion(workerNode, workerMachine, 2*time.Minute)
		})

		AfterEach(func() {
			glog.V(2).Infof("Delete node-unhealthy-conditions configmap")
			err := e2e.DeleteUnhealthyConditionsConfigMap()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("should delete unhealthy machine", func() {
		stopKubeletAndValidateMachineDeletion(workerNode, workerMachine, 6*time.Minute)
	})

	AfterEach(func() {
		// TODO: enable once https://github.com/openshift/cluster-api-actuator-pkg/pull/61 is fixed
		glog.V(2).Info("Skipping machine health checking test")
		Skip("Skipping machine health checking test")

		waitForWorkersToGetReady(numberOfReadyWorkers)
		deleteMachineHealthCheck(e2e.MachineHealthCheckName)
		deleteKubeletKillerPods()
	})
})

func waitForNodeUnhealthyCondition(workerNodeName string) {
	client, err := e2e.LoadClient()
	Expect(err).ToNot(HaveOccurred())

	key := types.NamespacedName{
		Name:      workerNodeName,
		Namespace: e2e.TestContext.MachineApiNamespace,
	}
	node := &corev1.Node{}
	glog.Infof("Wait until node %s will have 'Ready' condition with the status %s", node.Name, corev1.ConditionUnknown)
	Eventually(func() bool {
		err := client.Get(context.TODO(), key, node)
		if err != nil {
			return false
		}
		readyCond := conditions.GetNodeCondition(node, corev1.NodeReady)
		glog.V(2).Infof("Node %s has 'Ready' condition with the status %s", node.Name, readyCond.Status)
		return readyCond.Status == corev1.ConditionUnknown
	}, e2e.WaitLong, 10*time.Second).Should(BeTrue())
}

func waitForWorkersToGetReady(numberOfReadyWorkers int) {
	client, err := e2e.LoadClient()
	Expect(err).ToNot(HaveOccurred())

	glog.V(2).Infof("Wait until the environment will have %d ready workers", numberOfReadyWorkers)
	Eventually(func() bool {
		workerNodes, err := e2e.GetWorkerNodes(client)
		if err != nil {
			return false
		}

		readyWorkerNodes := e2e.FilterReadyNodes(workerNodes)
		glog.V(2).Infof("Number of ready workers %d", len(readyWorkerNodes))
		return len(readyWorkerNodes) == numberOfReadyWorkers
	}, 15*time.Minute, 10*time.Second).Should(BeTrue())
}

func deleteMachineHealthCheck(healthcheckName string) {
	client, err := e2e.LoadClient()
	Expect(err).ToNot(HaveOccurred())

	key := types.NamespacedName{
		Name:      healthcheckName,
		Namespace: e2e.TestContext.MachineApiNamespace,
	}
	healthcheck := &healthcheckingv1alpha1.MachineHealthCheck{}
	err = client.Get(context.TODO(), key, healthcheck)
	Expect(err).ToNot(HaveOccurred())

	glog.V(2).Infof("Delete machine health check %s", healthcheck.Name)
	err = client.Delete(context.TODO(), healthcheck)
	Expect(err).ToNot(HaveOccurred())
}

func deleteKubeletKillerPods() {
	client, err := e2e.LoadClient()
	Expect(err).ToNot(HaveOccurred())

	listOptions := runtimeclient.ListOptions{
		Namespace: e2e.TestContext.MachineApiNamespace,
	}
	listOptions.MatchingLabels(map[string]string{e2e.KubeletKillerPodName: ""})
	podList := &corev1.PodList{}
	err = client.List(context.TODO(), &listOptions, podList)
	Expect(err).ToNot(HaveOccurred())

	for _, pod := range podList.Items {
		glog.V(2).Infof("Delete kubelet killer pod %s", pod.Name)
		err = client.Delete(context.TODO(), &pod)
		Expect(err).ToNot(HaveOccurred())
	}
}
