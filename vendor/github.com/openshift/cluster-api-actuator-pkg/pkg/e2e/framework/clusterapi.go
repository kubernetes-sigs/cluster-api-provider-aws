package framework

import (
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/manifests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) DeployClusterAPIStack(clusterAPINamespace, actuatorPrivateKey string) {

	f.By("Deploying cluster API stack components")

	f.By("Deploying cluster CRD manifest")
	clusterCRDManifest := manifests.ClusterCRDManifest()
	err := WaitUntilCreated(func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(clusterCRDManifest)
		return err
	}, func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(clusterCRDManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deploying machine CRD manifest")
	machineCRDManifest := manifests.MachineCRDManifest()
	err = WaitUntilCreated(func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(machineCRDManifest)
		return err
	}, func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineCRDManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deploying machineset CRD manifest")
	machineSetCRDManifest := manifests.MachineSetCRDManifest()
	err = WaitUntilCreated(func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(machineSetCRDManifest)
		return err
	}, func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineSetCRDManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deploying machinedeployment CRD manifest")
	machineDeploymentCRDManifest := manifests.MachineDeploymentCRDManifest()
	err = WaitUntilCreated(func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(machineDeploymentCRDManifest)
		return err
	}, func() error {
		_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineDeploymentCRDManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deploying cluster role")
	clusterRoleManifest := manifests.ClusterRoleManifest()
	_, err = f.KubeClient.RbacV1().ClusterRoles().Create(clusterRoleManifest)
	f.ErrNotExpected(err)

	clusterRoleBinding := manifests.ClusterRoleBinding(clusterAPINamespace)
	_, err = f.KubeClient.RbacV1().ClusterRoleBindings().Create(clusterRoleBinding)
	f.ErrNotExpected(err)

	f.By("Deploying machine API controllers")
	deploymentManifest := manifests.ClusterAPIControllersDeployment(clusterAPINamespace, f.MachineControllerImage, f.MachineManagerImage, f.NodelinkControllerImage, actuatorPrivateKey)
	_, err = f.KubeClient.AppsV1().Deployments(deploymentManifest.Namespace).Create(deploymentManifest)
	f.ErrNotExpected(err)

	f.By("Waiting until cluster objects can be listed")
	err = wait.Poll(PollInterval, PoolClusterAPIDeploymentTimeout, func() (bool, error) {
		// Any namespace will do just to that one can list cluster objects
		_, err := f.CAPIClient.ClusterV1alpha1().Clusters("default").List(metav1.ListOptions{})
		if err != nil {
			glog.V(2).Infof("unable to list clusters: %v", err)
			return false, nil
		}

		return true, nil
	})

	f.By("Cluster API stack deployed")
}

func (f *Framework) DestroyClusterAPIStack(clusterAPINamespace, actuatorPrivateKey string) {

	f.By("Deleting machine API controllers")
	deploymentManifest := manifests.ClusterAPIControllersDeployment(clusterAPINamespace, f.MachineControllerImage, f.MachineManagerImage, f.NodelinkControllerImage, actuatorPrivateKey)
	err := WaitUntilDeleted(func() error {
		return f.KubeClient.AppsV1().Deployments(deploymentManifest.Namespace).Delete(deploymentManifest.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.AppsV1().Deployments(deploymentManifest.Namespace).Get(deploymentManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	f.By("Deleting cluster role")
	clusterRoleBinding := manifests.ClusterRoleBinding(clusterAPINamespace)
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.RbacV1().ClusterRoleBindings().Delete(clusterRoleBinding.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.RbacV1().ClusterRoleBindings().Get(clusterRoleBinding.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	clusterRoleManifest := manifests.ClusterRoleManifest()
	err = WaitUntilDeleted(func() error {
		return f.KubeClient.RbacV1().ClusterRoles().Delete(clusterRoleManifest.Name, &metav1.DeleteOptions{})
	}, func() error {
		_, err := f.KubeClient.RbacV1().ClusterRoles().Get(clusterRoleManifest.Name, metav1.GetOptions{})
		return err
	})
	f.ErrNotExpected(err)

	// Do not delete the CRDs as they are always in the default namespace.
	// Deleting and creating the same object in a row puts the object into terminating state for a long time at some point.
	// f.By("Deleting machinedeployment CRD manifest")
	// machineDeploymentCRDManifest := manifests.MachineDeploymentCRDManifest()
	// err = WaitUntilDeleted(func() error {
	// 	return f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(machineDeploymentCRDManifest.Name, &metav1.DeleteOptions{})
	// }, func() error {
	// 	_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineDeploymentCRDManifest.Name, metav1.GetOptions{})
	// 	return err
	// })
	// f.ErrNotExpected(err)
	//
	// f.By("Deleting machineset CRD manifest")
	// machineSetCRDManifest := manifests.MachineSetCRDManifest()
	// err = WaitUntilDeleted(func() error {
	// 	return f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(machineSetCRDManifest.Name, &metav1.DeleteOptions{})
	// }, func() error {
	// 	_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineSetCRDManifest.Name, metav1.GetOptions{})
	// 	return err
	// })
	// f.ErrNotExpected(err)
	//
	// f.By("Deleting machine CRD manifest")
	// machineCRDManifest := manifests.MachineCRDManifest()
	// err = WaitUntilDeleted(func() error {
	// 	return f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(machineCRDManifest.Name, &metav1.DeleteOptions{})
	// }, func() error {
	// 	_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(machineCRDManifest.Name, metav1.GetOptions{})
	// 	return err
	// })
	// f.ErrNotExpected(err)
	//
	// f.By("Deleting cluster CRD manifest")
	// clusterCRDManifest := manifests.ClusterCRDManifest()
	// err = WaitUntilDeleted(func() error {
	// 	return f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(clusterCRDManifest.Name, &metav1.DeleteOptions{})
	// }, func() error {
	// 	_, err := f.APIExtensionClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(clusterCRDManifest.Name, metav1.GetOptions{})
	// 	return err
	// })
	// f.ErrNotExpected(err)
}
