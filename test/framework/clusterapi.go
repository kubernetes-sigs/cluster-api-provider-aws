package framework

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/wait"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/log"

	"sigs.k8s.io/cluster-api-provider-aws/test/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) DeployClusterAPIStack(clusterAPINamespace string) {

	By("Deploy cluster API stack components")
	certsSecret, apiAPIService, err := utils.ClusterAPIServerAPIServiceObjects(clusterAPINamespace)
	Expect(err).NotTo(HaveOccurred())
	_, err = f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Create(certsSecret)
	Expect(err).NotTo(HaveOccurred())

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		if _, err := f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Get(certsSecret.Name, metav1.GetOptions{}); err != nil {
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = f.APIRegistrationClient.Apiregistration().APIServices().Create(apiAPIService)
	Expect(err).NotTo(HaveOccurred())

	apiService := utils.ClusterAPIService(clusterAPINamespace)
	_, err = f.KubeClient.CoreV1().Services(apiService.Namespace).Create(apiService)
	Expect(err).NotTo(HaveOccurred())

	clusterAPIDeployment := utils.ClusterAPIDeployment(clusterAPINamespace)
	_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIDeployment)
	Expect(err).NotTo(HaveOccurred())

	clusterAPIControllersDeployment := utils.ClusterAPIControllersDeployment(clusterAPINamespace)
	_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIControllersDeployment)
	Expect(err).NotTo(HaveOccurred())

	clusterAPIRoleBinding := utils.ClusterAPIRoleBinding(clusterAPINamespace)
	_, err = f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Create(clusterAPIRoleBinding)
	Expect(err).NotTo(HaveOccurred())

	clusterAPIEtcdCluster := utils.ClusterAPIEtcdCluster(clusterAPINamespace)
	_, err = f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Create(clusterAPIEtcdCluster)
	Expect(err).NotTo(HaveOccurred())

	etcdService := utils.ClusterAPIEtcdService(clusterAPINamespace)
	_, err = f.KubeClient.CoreV1().Services(etcdService.Namespace).Create(etcdService)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting for cluster API stack to come up")
	err = wait.Poll(PollInterval, PoolClusterAPIDeploymentTimeout, func() (bool, error) {
		if deployment, err := f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Get(clusterAPIDeployment.Name, metav1.GetOptions{}); err == nil {
			// Check all the pods are running
			log.Infof("Waiting for all cluster-api deployment pods to be ready, have %v, expecting 1", deployment.Status.ReadyReplicas)
			if deployment.Status.ReadyReplicas < 1 {
				return false, nil
			}
			return true, nil
		}

		return false, nil
	})
	Expect(err).NotTo(HaveOccurred())

	By("Cluster API stack deployed")
}

func (f *Framework) DestroyClusterAPIStack(clusterAPINamespace string) {
	var orphanDeletepolicy metav1.DeletionPropagation = "Orphan"
	var zero int64 = 0

	By("Deleting etcd service")
	etcdService := utils.ClusterAPIEtcdService(clusterAPINamespace)
	err := WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Services(etcdService.Namespace).Delete(etcdService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Services(etcdService.Namespace).Get(etcdService.Name, metav1.GetOptions{})
		return err
	})
	Expect(err).NotTo(HaveOccurred())

	By("Scaling down etcd cluster")
	clusterAPIEtcdCluster := utils.ClusterAPIEtcdCluster(clusterAPINamespace)
	f.ScaleSatefulSetDownToZero(clusterAPIEtcdCluster)
	Expect(err).NotTo(HaveOccurred())

	By("Deleting etcd cluster")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Delete(clusterAPIEtcdCluster.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		obj, err := f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Get(clusterAPIEtcdCluster.Name, metav1.GetOptions{})
		fmt.Printf("obj: %#v\n", obj)
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	By("Deleting role binding")
	clusterAPIRoleBinding := utils.ClusterAPIRoleBinding(clusterAPINamespace)
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Delete(clusterAPIRoleBinding.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Get(clusterAPIRoleBinding.Name, metav1.GetOptions{})
		return err
	})
	Expect(err).NotTo(HaveOccurred())

	clusterAPIControllersDeployment := utils.ClusterAPIControllersDeployment(clusterAPINamespace)
	By("Scaling down controllers deployment")
	err = f.ScaleDeploymentDownToZero(clusterAPIControllersDeployment)
	Expect(err).NotTo(HaveOccurred())

	By("Deleting controllers deployment")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().Deployments(clusterAPIControllersDeployment.Namespace).Delete(clusterAPIControllersDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		_, err := f.KubeClient.AppsV1beta2().Deployments(clusterAPIControllersDeployment.Namespace).Get(clusterAPIControllersDeployment.Name, metav1.GetOptions{})
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	clusterAPIDeployment := utils.ClusterAPIDeployment(clusterAPINamespace)
	By("Scaling down apiserver deployment")
	err = f.ScaleDeploymentDownToZero(clusterAPIDeployment)
	Expect(err).NotTo(HaveOccurred())

	By("Deleting apiserver deployment")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Delete(clusterAPIDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		_, err := f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Get(clusterAPIDeployment.Name, metav1.GetOptions{})
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	By("Deleting cluster api service")
	apiService := utils.ClusterAPIService(clusterAPINamespace)
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Services(apiService.Namespace).Delete(apiService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Services(apiService.Namespace).Get(apiService.Name, metav1.GetOptions{})
		return err
	})
	Expect(err).NotTo(HaveOccurred())

	// Even though the certs are different, only the secret name(space) and apiservice name(space) are actually used
	certsSecret, apiAPIService, err := utils.ClusterAPIServerAPIServiceObjects(clusterAPINamespace)
	Expect(err).NotTo(HaveOccurred())

	By("Deleting cluster api api service")
	err = WaitUntilDeleted(func() error {
		return f.APIRegistrationClient.Apiregistration().APIServices().Delete(apiAPIService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.APIRegistrationClient.Apiregistration().APIServices().Get(apiAPIService.Name, metav1.GetOptions{})
		return err
	})
	Expect(err).NotTo(HaveOccurred())

	By("Deleting api server certs")
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Delete(certsSecret.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Get(certsSecret.Name, metav1.GetOptions{})
		return err
	})
	Expect(err).NotTo(HaveOccurred())
}
