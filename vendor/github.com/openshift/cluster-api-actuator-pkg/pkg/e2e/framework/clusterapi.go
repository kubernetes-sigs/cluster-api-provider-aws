package framework

import (
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/prometheus/common/log"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/manifests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) DeployClusterAPIStack(clusterAPINamespace, actuatorImage, actuatorPrivateKey string) {

	f.By("Deploying cluster API stack components")
	f.By("Deploying API server API service")
	certsSecret, apiAPIService, err := manifests.ClusterAPIServerAPIServiceObjects(clusterAPINamespace)
	f.ErrNotExpected(err)
	_, err = f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Create(certsSecret)
	f.ErrNotExpected(err)

	err = wait.Poll(PollInterval, PoolTimeout, func() (bool, error) {
		if _, err := f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Get(certsSecret.Name, metav1.GetOptions{}); err != nil {
			return false, nil
		}
		return true, nil
	})
	f.ErrNotExpected(err)

	f.By("Deploying API service")
	_, err = f.APIRegistrationClient.Apiregistration().APIServices().Create(apiAPIService)
	f.ErrNotExpected(err)

	apiService := manifests.ClusterAPIService(clusterAPINamespace)
	_, err = f.KubeClient.CoreV1().Services(apiService.Namespace).Create(apiService)
	f.ErrNotExpected(err)

	f.By("Deploying apiserver")
	clusterAPIDeployment := manifests.ClusterAPIDeployment(clusterAPINamespace)
	_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIDeployment)
	f.ErrNotExpected(err)

	f.By("Deploying controllers ")
	clusterAPIControllersDeployment := manifests.ClusterAPIControllersDeployment(clusterAPINamespace, actuatorImage, actuatorPrivateKey)
	_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIControllersDeployment)
	f.ErrNotExpected(err)

	f.By("Deploying role binding")
	clusterAPIRoleBinding := manifests.ClusterAPIRoleBinding(clusterAPINamespace)
	_, err = f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Create(clusterAPIRoleBinding)
	f.ErrNotExpected(err)

	f.By("Deploying etcd cluster")
	clusterAPIEtcdCluster := manifests.ClusterAPIEtcdCluster(clusterAPINamespace)
	_, err = f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Create(clusterAPIEtcdCluster)
	f.ErrNotExpected(err)

	f.By("Deploying etcd service")
	etcdService := manifests.ClusterAPIEtcdService(clusterAPINamespace)
	_, err = f.KubeClient.CoreV1().Services(etcdService.Namespace).Create(etcdService)
	f.ErrNotExpected(err)

	f.By("Waiting for cluster API stack to come up")
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
	f.ErrNotExpected(err)

	f.By("Cluster API stack deployed")
}

func (f *Framework) DestroyClusterAPIStack(clusterAPINamespace, actuatorImage, actuatorPrivateKey string) {
	var orphanDeletepolicy metav1.DeletionPropagation = "Orphan"
	var zero int64 = 0

	f.By("Deleting etcd service")
	etcdService := manifests.ClusterAPIEtcdService(clusterAPINamespace)
	err := WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Services(etcdService.Namespace).Delete(etcdService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Services(etcdService.Namespace).Get(etcdService.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Scaling down etcd cluster")
	clusterAPIEtcdCluster := manifests.ClusterAPIEtcdCluster(clusterAPINamespace)
	f.ScaleSatefulSetDownToZero(clusterAPIEtcdCluster)
	f.ErrNotExpected(err)

	f.By("Deleting etcd cluster")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Delete(clusterAPIEtcdCluster.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		_, err := f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Get(clusterAPIEtcdCluster.Name, metav1.GetOptions{})
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	f.By("Deleting role binding")
	clusterAPIRoleBinding := manifests.ClusterAPIRoleBinding(clusterAPINamespace)
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Delete(clusterAPIRoleBinding.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Get(clusterAPIRoleBinding.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	clusterAPIControllersDeployment := manifests.ClusterAPIControllersDeployment(clusterAPINamespace, actuatorImage, actuatorPrivateKey)
	f.By("Scaling down controllers deployment")
	err = f.ScaleDeploymentDownToZero(clusterAPIControllersDeployment)
	f.ErrNotExpected(err)

	f.By("Deleting controllers deployment")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().Deployments(clusterAPIControllersDeployment.Namespace).Delete(clusterAPIControllersDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		_, err := f.KubeClient.AppsV1beta2().Deployments(clusterAPIControllersDeployment.Namespace).Get(clusterAPIControllersDeployment.Name, metav1.GetOptions{})
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	clusterAPIDeployment := manifests.ClusterAPIDeployment(clusterAPINamespace)
	f.By("Scaling down apiserver deployment")
	err = f.ScaleDeploymentDownToZero(clusterAPIDeployment)
	f.ErrNotExpected(err)

	f.By("Deleting apiserver deployment")
	WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Delete(clusterAPIDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &orphanDeletepolicy, GracePeriodSeconds: &zero})
	}, func() error {
		_, err := f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Get(clusterAPIDeployment.Name, metav1.GetOptions{})
		return err
	})
	// Ignore the error, the deployment has 0 replicas.
	// No longer affecting future deployments since it lives in a different namespace.

	f.By("Deleting cluster api service")
	apiService := manifests.ClusterAPIService(clusterAPINamespace)
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Services(apiService.Namespace).Delete(apiService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Services(apiService.Namespace).Get(apiService.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	// Even though the certs are different, only the secret name(space) and apiservice name(space) are actually used
	certsSecret, apiAPIService, err := manifests.ClusterAPIServerAPIServiceObjects(clusterAPINamespace)
	f.ErrNotExpected(err)

	f.By("Deleting cluster api api service")
	err = WaitUntilDeleted(func() error {
		return f.APIRegistrationClient.Apiregistration().APIServices().Delete(apiAPIService.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.APIRegistrationClient.Apiregistration().APIServices().Get(apiAPIService.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deleting api server certs")
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Delete(certsSecret.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Get(certsSecret.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)
}
