package operators

import (
	"context"
	"time"

	"github.com/golang/glog"
	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	osconfigv1 "github.com/openshift/api/config/v1"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	cvoresourcemerge "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	kappsapi "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = g.Describe("[Feature:Operators] Cluster autoscaler operator should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		key := types.NamespacedName{
			Namespace: e2e.TestContext.MachineApiNamespace,
			Name:      "cluster-autoscaler-operator",
		}
		d := &kappsapi.Deployment{}

		err = wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
			if err := client.Get(context.TODO(), key, d); err != nil {
				glog.Errorf("error querying api for Deployment object: %v, retrying...", err)
				return false, nil
			}
			if d.Status.ReadyReplicas < 1 {
				return false, nil
			}
			return true, nil
		})
		o.Expect(err).NotTo(o.HaveOccurred())
	})
})

var _ = g.Describe("[Feature:Operators] Cluster autoscaler cluster operator should", func() {
	defer g.GinkgoRecover()

	g.It("be available", func() {
		var err error
		client, err := e2e.LoadClient()
		o.Expect(err).NotTo(o.HaveOccurred())

		key := types.NamespacedName{
			Namespace: e2e.TestContext.MachineApiNamespace,
			Name:      "cluster-autoscaler",
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
