package framework

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/log"
)

func (f *Framework) CreateClusterAndWait(cluster *clusterv1alpha1.Cluster) {
	By(fmt.Sprintf("Creating %q cluster", cluster.Name))
	err := wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Create(cluster)
		if err != nil {
			log.Infof("error creating cluster: %v", err)
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		_, err := f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())
}
