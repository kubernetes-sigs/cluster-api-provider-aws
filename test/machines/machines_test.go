package machines

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/util/uuid"
	// "k8s.io/apimachinery/pkg/util/uuid"

	"sigs.k8s.io/cluster-api-provider-aws/test/framework"
	"sigs.k8s.io/cluster-api-provider-aws/test/utils"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	machineutils "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/client"
)

const (
	clusterAPIStackDefaultNamespace   = "kube-system"
	poolTimeout                       = 20 * time.Second
	pollInterval                      = 1 * time.Second
	poolClusterAPIDeploymentTimeout   = 10 * time.Minute
	timeoutPoolMachineRunningInterval = 10 * time.Minute
	region                            = "us-east-1"
	awsCredentialsSecretName          = "aws-credentials-secret"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Machine Suite")
}

var _ = framework.SigKubeDescribe("Machines", func() {
	f := framework.NewFramework()
	var err error
	var testNamespace *apiv1.Namespace
	var testMachine *clusterv1alpha1.Machine

	BeforeEach(func() {
		f.BeforeEach()
		By("Deploy cluster API stack components")

		certsSecret := utils.ClusterAPIServerCertsSecret()
		_, err = f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Create(certsSecret)
		Expect(err).NotTo(HaveOccurred())

		err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
			if _, err := f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Get(certsSecret.Name, metav1.GetOptions{}); err != nil {
				return false, nil
			}
			return true, nil
		})
		Expect(err).NotTo(HaveOccurred())

		apiAPIService := utils.ClusterAPIAPIService()
		_, err = f.APIRegistrationClient.Apiregistration().APIServices().Create(apiAPIService)
		Expect(err).NotTo(HaveOccurred())

		apiService := utils.ClusterAPIService()
		_, err = f.KubeClient.CoreV1().Services(apiService.Namespace).Create(apiService)
		Expect(err).NotTo(HaveOccurred())

		clusterAPIDeployment := utils.ClusterAPIDeployment()
		_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIDeployment)
		Expect(err).NotTo(HaveOccurred())

		clusterAPIControllersDeployment := utils.ClusterAPIControllersDeployment()
		_, err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Create(clusterAPIControllersDeployment)
		Expect(err).NotTo(HaveOccurred())

		clusterAPIRoleBinding := utils.ClusterAPIRoleBinding()
		_, err = f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Create(clusterAPIRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		clusterAPIEtcdCluster := utils.ClusterAPIEtcdCluster()
		_, err = f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Create(clusterAPIEtcdCluster)
		Expect(err).NotTo(HaveOccurred())

		etcdService := utils.ClusterAPIEtcdService()
		_, err = f.KubeClient.CoreV1().Services(etcdService.Namespace).Create(etcdService)
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for cluster API stack to come up")
		err = wait.Poll(pollInterval, poolClusterAPIDeploymentTimeout, func() (bool, error) {
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

	})

	AfterEach(func() {
		var err error
		var foregroundDeletepolicy metav1.DeletionPropagation = "Foreground"

		// Make sure the test machine is deleted before deleting its namespace
		if testMachine != nil {
			log.Infof(testMachine.Name+": %#v", testMachine)
			By("Deleting testing machine")
			err = f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Delete(testMachine.Name, &metav1.DeleteOptions{})
			framework.IgnoreNotFoundErr(err)

			// Verify the testing machine has been destroyed
			err = wait.Poll(pollInterval, timeoutPoolMachineRunningInterval, func() (bool, error) {
				_, err := f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Get(testMachine.Name, metav1.GetOptions{})
				if err == nil {
					log.Info("Waiting for machine to be deleted")
					return false, nil
				}
				if strings.Contains(err.Error(), "not found") {
					return true, nil
				}
				return false, nil
			})
			framework.IgnoreNotFoundErr(err)
		}

		if testNamespace != nil {
			log.Infof(testNamespace.Name+": %#v", testNamespace)
			By(fmt.Sprintf("Destroying %q namespace", testNamespace.Name))
			err = f.KubeClient.CoreV1().Namespaces().Delete(testNamespace.Name, &metav1.DeleteOptions{})
			framework.IgnoreNotFoundErr(err)
		}

		etcdService := utils.ClusterAPIEtcdService()
		err = f.KubeClient.CoreV1().Services(etcdService.Namespace).Delete(etcdService.Name, &metav1.DeleteOptions{})
		framework.IgnoreNotFoundErr(err)

		clusterAPIEtcdCluster := utils.ClusterAPIEtcdCluster()
		err = f.KubeClient.AppsV1beta2().StatefulSets(clusterAPIEtcdCluster.Namespace).Delete(clusterAPIEtcdCluster.Name, &metav1.DeleteOptions{PropagationPolicy: &foregroundDeletepolicy})
		framework.IgnoreNotFoundErr(err)

		clusterAPIRoleBinding := utils.ClusterAPIRoleBinding()
		err = f.KubeClient.RbacV1().RoleBindings(clusterAPIRoleBinding.Namespace).Delete(clusterAPIRoleBinding.Name, &metav1.DeleteOptions{})
		framework.IgnoreNotFoundErr(err)

		clusterAPIControllersDeployment := utils.ClusterAPIControllersDeployment()
		err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIControllersDeployment.Namespace).Delete(clusterAPIControllersDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &foregroundDeletepolicy})
		framework.IgnoreNotFoundErr(err)

		clusterAPIDeployment := utils.ClusterAPIDeployment()
		err = f.KubeClient.AppsV1beta2().Deployments(clusterAPIDeployment.Namespace).Delete(clusterAPIDeployment.Name, &metav1.DeleteOptions{PropagationPolicy: &foregroundDeletepolicy})
		framework.IgnoreNotFoundErr(err)

		apiService := utils.ClusterAPIService()
		err = f.KubeClient.CoreV1().Services(apiService.Namespace).Delete(apiService.Name, &metav1.DeleteOptions{})
		framework.IgnoreNotFoundErr(err)

		apiAPIService := utils.ClusterAPIAPIService()
		err = f.APIRegistrationClient.Apiregistration().APIServices().Delete(apiAPIService.Name, &metav1.DeleteOptions{})
		framework.IgnoreNotFoundErr(err)

		certsSecret := utils.ClusterAPIServerCertsSecret()
		err = f.KubeClient.CoreV1().Secrets(certsSecret.Namespace).Delete(certsSecret.Name, &metav1.DeleteOptions{})
		framework.IgnoreNotFoundErr(err)

	})

	// Any of the tests run assumes the cluster-api stack is already deployed.
	// So all the machine, resp. machineset related tests must be run on top
	// of the same cluster-api stack. Once the machine, resp. machineset objects
	// are defined through CRD, we can relax the restriction.
	Context("AWS actuator", func() {

		It("Can create AWS instances", func() {
			testNamespace = &apiv1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "namespace-" + string(uuid.NewUUID()),
				},
			}

			By(fmt.Sprintf("Creating %q namespace", testNamespace.Name))
			_, err = f.KubeClient.CoreV1().Namespaces().Create(testNamespace)
			Expect(err).NotTo(HaveOccurred())

			awsCredSecret := utils.GenerateAwsCredentialsSecretFromEnv(awsCredentialsSecretName, testNamespace.Name)
			_, err = f.KubeClient.CoreV1().Secrets(awsCredSecret.Namespace).Create(awsCredSecret)
			Expect(err).NotTo(HaveOccurred())

			err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
				if _, err := f.KubeClient.CoreV1().Secrets(awsCredSecret.Namespace).Get(awsCredSecret.Name, metav1.GetOptions{}); err != nil {
					return false, nil
				}
				return true, nil
			})
			Expect(err).NotTo(HaveOccurred())

			clusterID := framework.ClusterID
			if clusterID == "" {
				clusterID = "cluster-" + string(uuid.NewUUID())
			}

			cluster := &clusterv1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      clusterID,
					Namespace: testNamespace.Name,
				},
				Spec: clusterv1alpha1.ClusterSpec{
					ClusterNetwork: clusterv1alpha1.ClusterNetworkingConfig{
						Services: clusterv1alpha1.NetworkRanges{
							CIDRBlocks: []string{"10.0.0.1/24"},
						},
						Pods: clusterv1alpha1.NetworkRanges{
							CIDRBlocks: []string{"10.0.0.1/24"},
						},
						ServiceDomain: "example.com",
					},
				},
			}

			By(fmt.Sprintf("Creating %q cluster", cluster.Name))
			err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
				_, err = f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Create(cluster)
				if err != nil {
					return false, nil
				}
				return true, nil
			})

			Expect(err).NotTo(HaveOccurred())
			err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
				_, err := f.CAPIClient.ClusterV1alpha1().Clusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
				if err != nil {
					return false, nil
				}
				return true, nil
			})
			Expect(err).NotTo(HaveOccurred())

			testMachine = utils.TestingMachine(awsCredSecret.Name, cluster.Name, cluster.Namespace)
			By(fmt.Sprintf("Creating %q machine", testMachine.Name))
			_, err = f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Create(testMachine)
			Expect(err).NotTo(HaveOccurred())

			// Verify cluster and machine have been deployed
			err = wait.Poll(pollInterval, timeoutPoolMachineRunningInterval, func() (bool, error) {
				if _, err := f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Get(testMachine.Name, metav1.GetOptions{}); err != nil {
					log.Info("Waiting for cluster and machine to be created")
					return false, nil
				}
				return true, nil
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verify AWS instance is running")
			awsClient, err := awsclient.NewClient(f.KubeClient, awsCredSecret.Name, awsCredSecret.Namespace, region)
			Expect(err).NotTo(HaveOccurred())

			err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
				log.Info("Waiting for aws instances to come up")
				runningInstances, err := machineutils.GetRunningInstances(testMachine, awsClient)
				if err != nil {
					return false, fmt.Errorf("unable to get running instances from aws: %v", err)
				}
				runningInstancesLen := len(runningInstances)
				if runningInstancesLen == 1 {
					log.Info("Machine is running on aws")
					return true, nil
				}
				if runningInstancesLen > 1 {
					return false, fmt.Errorf("Found %q instances instead of one", runningInstancesLen)
				}
				return false, nil
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verify AWS instance is terminated")
			err = f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Delete(testMachine.Name, &metav1.DeleteOptions{})
			framework.IgnoreNotFoundErr(err)

			// Verify the testing machine has been destroyed
			err = wait.Poll(pollInterval, poolTimeout, func() (bool, error) {
				_, err := f.CAPIClient.ClusterV1alpha1().Machines(testMachine.Namespace).Get(testMachine.Name, metav1.GetOptions{})
				if err == nil {
					log.Info("Waiting for machine to be deleted")
					return false, nil
				}
				if strings.Contains(err.Error(), "not found") {
					return true, nil
				}
				return false, nil
			})
			framework.IgnoreNotFoundErr(err)

		})
	})

})
