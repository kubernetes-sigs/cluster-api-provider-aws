package framework

import (
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (f *Framework) CreateClusterAndWait(cluster *clusterv1alpha1.Cluster) {
	f.By(fmt.Sprintf("Creating %q cluster", cluster.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Create(cluster)
		if err != nil {
			glog.V(2).Infof("error creating cluster: %v", err)
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)
}
