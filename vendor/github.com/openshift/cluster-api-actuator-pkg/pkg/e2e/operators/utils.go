package operators

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	kappsapi "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func getDeployment(client runtimeclient.Client, name string) (*kappsapi.Deployment, error) {
	key := types.NamespacedName{
		Namespace: e2e.TestContext.MachineApiNamespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	if err := wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("Error querying api for Deployment object %q: %v, retrying...", name, err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("error getting deployment %q: %v", name, err)
	}
	return d, nil
}

func deleteDeployment(client runtimeclient.Client, deployment *kappsapi.Deployment) error {
	return wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.Delete(context.TODO(), deployment); err != nil {
			glog.Errorf("error querying api for deployment object %q: %v, retrying...", deployment.Name, err)
			return false, nil
		}
		return true, nil
	})
}

func isDeploymentAvailable(client runtimeclient.Client, name string) bool {
	if err := wait.PollImmediate(1*time.Second, e2e.WaitLong, func() (bool, error) {
		d, err := getDeployment(client, name)
		if err != nil {
			glog.Errorf("Error getting deployment: %v", err)
			return false, nil
		}
		if d.Status.AvailableReplicas < 1 {
			glog.Errorf("Deployment %q is not available. Status: (replicas: %d, updated: %d, ready: %d, available: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
			return false, nil
		}
		glog.Infof("Deployment %q is available. Status: (replicas: %d, updated: %d, ready: %d, available: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
		return true, nil
	}); err != nil {
		glog.Errorf("Error checking isDeploymentAvailable: %v", err)
		return false
	}
	return true
}

func isStatusAvailable(client runtimeclient.Client, name string) bool {
	key := types.NamespacedName{
		Namespace: e2e.TestContext.MachineApiNamespace,
		Name:      name,
	}
	clusterOperator := &osconfigv1.ClusterOperator{}

	if err := wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.Get(context.TODO(), key, clusterOperator); err != nil {
			glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
			return false, nil
		}
		if cov1helpers.IsStatusConditionFalse(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable) {
			glog.Errorf("Condition: %q is false", osconfigv1.OperatorAvailable)
			return false, nil
		}
		if cov1helpers.IsStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorProgressing) {
			glog.Errorf("Condition: %q is true", osconfigv1.OperatorProgressing)
			return false, nil
		}
		if cov1helpers.IsStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorDegraded) {
			glog.Errorf("Condition: %q is true", osconfigv1.OperatorDegraded)
			return false, nil
		}
		return true, nil
	}); err != nil {
		glog.Errorf("Error checking isStatusAvailable: %v", err)
		return false
	}
	return true

}
