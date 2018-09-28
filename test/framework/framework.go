package framework

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"github.com/prometheus/common/log"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiregistrationclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// Default timeout for pools
	PoolTimeout = 20 * time.Second
	// Default waiting interval for pools
	PollInterval = 1 * time.Second
	// Node waiting internal
	PollNodeInterval = 5 * time.Second
	// Pool timeout for cluster API deployment
	PoolClusterAPIDeploymentTimeout = 10 * time.Minute
	PoolDeletionTimeout             = 1 * time.Minute
	// Pool timeout for kubeconfig
	PoolKubeConfigTimeout = 10 * time.Minute
	PoolNodesReadyTimeout = 5 * time.Minute
	// Instances are running timeout
	TimeoutPoolMachineRunningInterval = 10 * time.Minute
)

var kubeconfig string

// ClusterID set by -cluster-id flag
var ClusterID string

// Path to private ssh to connect to instances (e.g. to download kubeconfig or copy docker images)
var sshkey string

var actuatorImage string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&ClusterID, "cluster-id", "", "cluster ID")
	flag.StringVar(&sshkey, "ssh-key", "", "Path to private ssh to connect to instances (e.g. to download kubeconfig or copy docker images)")
	flag.StringVar(&actuatorImage, "actuator-image", "gcr.io/k8s-cluster-api/aws-machine-controller:0.0.1", "Actuator image to run")

	flag.Parse()
}

// Framework supports common operations used by tests
type Framework struct {
	KubeClient            *kubernetes.Clientset
	CAPIClient            *clientset.Clientset
	APIRegistrationClient *apiregistrationclientset.Clientset
	Kubeconfig            string
	RestConfig            *rest.Config
	SSHKey                string
	ActuatorImage         string
}

// NewFramework setups a new framework
func NewFramework() *Framework {
	if sshkey == "" {
		panic("-sskey not set")
	}
	if kubeconfig == "" {
		panic("-kubeconfig not set")
	}
	f := &Framework{
		Kubeconfig:    kubeconfig,
		SSHKey:        sshkey,
		ActuatorImage: actuatorImage,
	}

	BeforeEach(f.BeforeEach)
	return f
}

func NewFrameworkFromConfig(config *rest.Config) *Framework {
	if sshkey == "" {
		panic("-sskey not set")
	}

	f := &Framework{
		RestConfig:    config,
		SSHKey:        sshkey,
		ActuatorImage: actuatorImage,
	}

	f.buildClientsets()
	return f
}

func (f *Framework) buildClientsets() {
	var err error

	By("Creating a kubernetes client")
	if f.RestConfig == nil {
		f.RestConfig, err = controller.GetConfig(f.Kubeconfig)
		Expect(err).NotTo(HaveOccurred())
	}

	if f.KubeClient == nil {
		f.KubeClient, err = kubernetes.NewForConfig(f.RestConfig)
		Expect(err).NotTo(HaveOccurred())
	}

	if f.CAPIClient == nil {
		f.CAPIClient, err = clientset.NewForConfig(f.RestConfig)
		Expect(err).NotTo(HaveOccurred())
	}

	if f.APIRegistrationClient == nil {
		f.APIRegistrationClient, err = apiregistrationclientset.NewForConfig(f.RestConfig)
		Expect(err).NotTo(HaveOccurred())
	}
}

// BeforeEach to be run before each spec responsible for building various clientsets
func (f *Framework) BeforeEach() {
	f.buildClientsets()
}

func (f *Framework) ScaleSatefulSetDownToZero(statefulset *appsv1beta2.StatefulSet) error {
	var zero int32 = 0
	statefulset.Spec.Replicas = &zero
	err := wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {
		// give it some time
		_, err := f.KubeClient.AppsV1beta2().StatefulSets(statefulset.Namespace).Update(statefulset)
		log.Infof("ScaleSatefulSetDownToZero.err: %v\n", err)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	// Now wait the number of replicas is really zero
	return wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {
		// give it some time
		result, err := f.KubeClient.AppsV1beta2().StatefulSets(statefulset.Namespace).Get(statefulset.Name, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}
		fmt.Printf("result: %v\n", result.Status.CurrentReplicas)
		if result.Status.CurrentReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}

func (f *Framework) ScaleDeploymentDownToZero(deployment *appsv1beta2.Deployment) error {
	var zero int32 = 0
	deployment.Spec.Replicas = &zero
	err := wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {
		// give it some time
		_, err := f.KubeClient.AppsV1beta2().Deployments(deployment.Namespace).Update(deployment)
		log.Infof("ScaleDeploymentDownToZero.err: %v\n", err)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	// Now wait the number of replicas is really zero
	return wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {
		// give it some time
		result, err := f.KubeClient.AppsV1beta2().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}
		fmt.Printf("result: %v\n", result.Status.AvailableReplicas)
		if result.Status.AvailableReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}

func WaitUntilDeleted(delFnc func() error, getFnc func() error) error {
	return wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {

		err := delFnc()
		log.Infof("del.err: %v\n", err)
		if err != nil {
			if strings.Contains(err.Error(), "object is being deleted") {
				return false, nil
			}
			if strings.Contains(err.Error(), "not found") {
				return true, nil
			}
			return false, nil
		}

		err = getFnc()
		log.Infof("get.err: %v\n", err)
		if err != nil && strings.Contains(err.Error(), "not found") {
			return true, nil
		}
		return false, nil
	})
}

// IgnoreNotFoundErr ignores not found errors in case resource
// that does not exist is to be deleted
func IgnoreNotFoundErr(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

// SigKubeDescribe is a wrapper function for ginkgo describe.  Adds namespacing.
func SigKubeDescribe(text string, body func()) bool {
	return Describe("[sigs.k8s.io] "+text, body)
}

func (f *Framework) UploadDockerImageToInstance(image, targetMachine string) error {
	log.Infof("Uploading %q to the master machine under %q", image, targetMachine)
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		"docker save %v | bzip2 | ssh -o StrictHostKeyChecking=no -i %v ec2-user@%v \"bunzip2 > /tmp/tempimage.bz2 && sudo docker load -i /tmp/tempimage.bz2\"",
		image,
		f.SSHKey,
		targetMachine,
	))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(string(out))
		return err
	}
	log.Info(string(out))
	return nil
}
