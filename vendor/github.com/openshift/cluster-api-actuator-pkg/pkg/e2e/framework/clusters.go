package framework

import (
	"fmt"

	"github.com/golang/glog"

	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (f *Framework) CreateClusterAndWait(cluster *machinev1beta1.Cluster) {
	f.By(fmt.Sprintf("Creating %q cluster", cluster.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().Clusters(cluster.Namespace).Create(cluster)
		if err != nil {
			glog.V(2).Infof("error creating cluster: %v", err)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.MachineV1beta1().Clusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)
}
