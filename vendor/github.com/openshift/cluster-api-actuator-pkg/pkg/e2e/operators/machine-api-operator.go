package operators

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	osconfigv1 "github.com/openshift/api/config/v1"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	cvoresourcemerge "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	kappsapi "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = g.Describe("[Feature:Operators] Machine API operator should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		key := types.NamespacedName{
			Namespace: e2e.TestContext.MachineApiNamespace,
			Name:      "machine-api-operator",
		}
		d := &kappsapi.Deployment{}

		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Get(context.TODO(), key, d); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			glog.Infof("Have %v ready replicas, expecting 1", d.Status.ReadyReplicas)
			if d.Status.ReadyReplicas < 1 {
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("reconcile controllers deployment", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		key := types.NamespacedName{
			Namespace: e2e.TestContext.MachineApiNamespace,
			Name:      "clusterapi-manager-controllers",
		}
		d := &kappsapi.Deployment{}

		glog.Info("Get deployment")
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Get(context.TODO(), key, d); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("Getting pod owned by the deployment")
		var deploymentPod *corev1.Pod
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			pods := corev1.PodList{}
			listOpt := &runtimeclient.ListOptions{}
			listOpt.MatchingLabels(d.Spec.Selector.MatchLabels)
			if err := client.List(context.TODO(), listOpt, &pods); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			if len(pods.Items) != 1 {
				glog.Errorf("expected exactly one pod running, have %v", len(pods.Items))
				return false, nil
			}
			deploymentPod = &pods.Items[0]
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Delete deployment")
		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Delete(context.TODO(), d); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		glog.Info("Verify deployment is recreated")
		err = wait.PollImmediate(1*time.Second, e2e.WaitLong, func() (bool, error) {
			if err := client.Get(context.TODO(), key, d); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			if d.Status.ReadyReplicas < 1 || !d.DeletionTimestamp.IsZero() {
				return false, nil
			}

			pods := corev1.PodList{}
			listOpt := &runtimeclient.ListOptions{}
			listOpt.MatchingLabels(d.Spec.Selector.MatchLabels)
			if err := client.List(context.TODO(), listOpt, &pods); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}

			for _, pod := range pods.Items {
				fmt.Printf("pod(%v): %#v\n", pod.Name, pod.Status.Phase)
				if pod.Name == deploymentPod.Name {
					glog.Infof("Ignoring old deployment pod %v", pod.Name)
					continue
				}
				if pod.Status.Phase != corev1.PodRunning {
					glog.Errorf("Deployment pod %v not yet running, retrying...", pod.Name)
					return false, nil
				}
			}

			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

})

var _ = g.Describe("[Feature:Operators] Machine API cluster operator should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		key := types.NamespacedName{
			Namespace: e2e.TestContext.MachineApiNamespace,
			Name:      "machine-api",
		}
		clusterOperator := &osconfigv1.ClusterOperator{}

		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Get(context.TODO(), key, clusterOperator); err != nil {
				glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
				return false, nil
			}
			if available := cvoresourcemerge.FindOperatorStatusCondition(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable); available != nil {
				if available.Status == osconfigv1.ConditionTrue {
					return true, nil
				}
			}
			return false, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})

})
