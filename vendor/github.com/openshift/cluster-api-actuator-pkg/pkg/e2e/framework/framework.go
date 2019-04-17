package framework

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	healthcheckingclient "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned"
	kappsapi "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// Default timeout for pools
	PoolTimeout = 5 * time.Minute
	// Default waiting interval for pools
	PollInterval = 5 * time.Second
	// Node waiting internal
	PollNodeInterval = 5 * time.Second
	// Pool timeout for cluster API deployment
	PoolClusterAPIDeploymentTimeout = 10 * time.Minute
	PoolDeletionTimeout             = 1 * time.Minute
	// Pool timeout for kubeconfig
	PoolKubeConfigTimeout = 10 * time.Minute
	PoolNodesReadyTimeout = 10 * time.Minute
	// Instances are running timeout
	TimeoutPoolMachineRunningInterval = 10 * time.Minute
)

// ClusterID set by -cluster-id flag
var ClusterID string

// Path to private ssh to connect to instances (e.g. to download kubeconfig or copy docker images)
var sshkey string

// Default ssh user
var sshuser string

var machineControllerImage string
var machineManagerImage string
var nodelinkControllerImage string

var libvirtURI string
var libvirtPK string

type TestContextType struct {
	KubeConfig          string
	MachineApiNamespace string
	Host                string
}

var TestContext TestContextType

func init() {
	flag.StringVar(&TestContext.KubeConfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&TestContext.MachineApiNamespace, "machine-api-namespace", "openshift-machine-api", "Default machine API namespace")
	flag.StringVar(&ClusterID, "cluster-id", "", "cluster ID")
	flag.StringVar(&sshkey, "ssh-key", "", "Path to private ssh to connect to instances (e.g. to download kubeconfig or copy docker images)")
	flag.StringVar(&sshuser, "ssh-user", "ec2-user", "Ssh user to connect to instances")
	flag.StringVar(&machineControllerImage, "machine-controller-image", "gcr.io/k8s-cluster-api/machine-controller:0.0.1", "Machine controller (actuator) image to run")
	flag.StringVar(&machineManagerImage, "machine-manager-image", "gcr.io/k8s-cluster-api/machine-controller:0.0.1", "Machine manager image to run")
	flag.StringVar(&nodelinkControllerImage, "nodelink-controller-image", "gcr.io/k8s-cluster-api/machine-controller:0.0.1", "Nodelink controller image to run")

	// libvirt specific flags
	flag.StringVar(&libvirtURI, "libvirt-uri", "", "Libvirt URI to connect to libvirt from within machine controller container")
	flag.StringVar(&libvirtPK, "libvirt-pk", "", "Private key to connect to qemu+ssh libvirt uri")

	flag.Parse()
}

type ErrNotExpectedFnc func(error)
type ByFnc func(string)

type SSHConfig struct {
	Key  string
	User string
	Host string
}

// Framework supports common operations used by tests
type Framework struct {
	KubeClient           *kubernetes.Clientset
	CAPIClient           *clientset.Clientset
	APIExtensionClient   *apiextensionsclientset.Clientset
	HealthCheckingClient *healthcheckingclient.Clientset

	// APIRegistrationClient *apiregistrationclientset.Clientset
	Kubeconfig string
	RestConfig *rest.Config

	SSH *SSHConfig

	LibvirtURI string
	LibvirtPK  string

	MachineControllerImage  string
	MachineManagerImage     string
	NodelinkControllerImage string

	ErrNotExpected ErrNotExpectedFnc
	By             ByFnc
}

// NewFramework setups a new framework
func NewFramework() (*Framework, error) {
	if sshkey == "" {
		return nil, fmt.Errorf("-sshkey not set")
	}
	if TestContext.KubeConfig == "" {
		return nil, fmt.Errorf("-kubeconfig not set")
	}
	f := &Framework{
		Kubeconfig: TestContext.KubeConfig,
		SSH: &SSHConfig{
			Key:  sshkey,
			User: sshuser,
		},

		LibvirtURI: libvirtURI,
		LibvirtPK:  libvirtPK,

		MachineControllerImage:  machineControllerImage,
		MachineManagerImage:     machineManagerImage,
		NodelinkControllerImage: nodelinkControllerImage,
	}

	f.ErrNotExpected = f.DefaultErrNotExpected
	f.By = f.DefaultBy

	BeforeEach(f.BeforeEach)
	return f, nil
}

func DefaultSSHConfig() (*SSHConfig, error) {
	if sshkey == "" {
		return nil, fmt.Errorf("-sshkey not set")
	}

	return &SSHConfig{
		Key:  sshkey,
		User: sshuser,
	}, nil
}

func NewFrameworkFromConfig(config *rest.Config, sshConfig *SSHConfig) (*Framework, error) {
	f := &Framework{
		RestConfig: config,
		SSH:        sshConfig,
		MachineControllerImage:  machineControllerImage,
		MachineManagerImage:     machineManagerImage,
		NodelinkControllerImage: nodelinkControllerImage,
		LibvirtURI:              libvirtURI,
		LibvirtPK:               libvirtPK,
	}

	f.ErrNotExpected = f.DefaultErrNotExpected
	f.By = f.DefaultBy

	err := f.buildClientsets()
	return f, err
}

func (f *Framework) buildClientsets() error {
	var err error

	if f.RestConfig == nil {
		config, err := clientcmd.LoadFromFile(f.Kubeconfig)
		if err != nil {
			return err
		}
		f.RestConfig, err = clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return err
		}
	}

	if f.KubeClient == nil {
		f.KubeClient, err = kubernetes.NewForConfig(f.RestConfig)
		if err != nil {
			return err
		}
	}

	if f.CAPIClient == nil {
		f.CAPIClient, err = clientset.NewForConfig(f.RestConfig)
		if err != nil {
			return err
		}
	}

	if f.APIExtensionClient == nil {
		f.APIExtensionClient, err = apiextensionsclientset.NewForConfig(f.RestConfig)
		if err != nil {
			return err
		}
	}

	if f.HealthCheckingClient == nil {
		f.HealthCheckingClient, err = healthcheckingclient.NewForConfig(f.RestConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// BeforeEach to be run before each spec responsible for building various clientsets
func (f *Framework) BeforeEach() {
	err := f.buildClientsets()
	f.ErrNotExpected(err)
}

func (f *Framework) ScaleSatefulSetDownToZero(statefulset *appsv1beta2.StatefulSet) error {
	var zero int32 = 0
	statefulset.Spec.Replicas = &zero
	err := wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {
		// give it some time
		_, err := f.KubeClient.AppsV1beta2().StatefulSets(statefulset.Namespace).Update(statefulset)
		glog.V(2).Infof("ScaleSatefulSetDownToZero.err: %v\n", err)
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
		glog.V(2).Infof("ScaleDeploymentDownToZero.err: %v\n", err)
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

		if result.Status.AvailableReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}

func WaitUntilDeleted(delFnc func() error, getFnc func() error) error {
	return wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {

		err := delFnc()
		glog.V(2).Infof("del.err: %v\n", err)
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
		glog.V(2).Infof("get.err: %v\n", err)
		if err != nil && strings.Contains(err.Error(), "not found") {
			return true, nil
		}
		return false, nil
	})
}

func WaitUntilCreated(createFnc func() error, getFnc func() error) error {
	return wait.Poll(PollInterval, PoolDeletionTimeout, func() (bool, error) {

		err := createFnc()
		glog.V(2).Infof("create.err: %v\n", err)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				return true, nil
			}
			return false, nil
		}

		err = getFnc()
		glog.V(2).Infof("get.err: %v\n", err)
		if err == nil {
			return true, nil
		}
		return false, nil
	})
}

func (f *Framework) DefaultErrNotExpected(err error) {
	Expect(err).NotTo(HaveOccurred())
}

func (f *Framework) DefaultBy(msg string) {
	By(msg)
}

// IgnoreNotFoundErr ignores not found errors in case resource
// that does not exist is to be deleted
func (f *Framework) IgnoreNotFoundErr(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return
		}
		f.ErrNotExpected(err)
	}
}

// SigKubeDescribe is a wrapper function for ginkgo describe.  Adds namespacing.
func SigKubeDescribe(text string, body func()) bool {
	return Describe("[sigs.k8s.io] "+text, body)
}

func (f *Framework) UploadDockerImageToInstance(image, targetMachine string) error {
	glog.V(2).Infof("Uploading %q to the master machine under %q", image, targetMachine)
	cmdStr := fmt.Sprintf(
		"docker save %v | bzip2 | ssh -o StrictHostKeyChecking=no -i %v %v@%v \"bunzip2 > /tmp/tempimage.bz2 && sudo docker load -i /tmp/tempimage.bz2\"",
		image,
		f.SSH.Key,
		f.SSH.User,
		targetMachine,
	)
	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		glog.V(2).Info(string(out))
		return err
	}
	glog.V(2).Info(string(out))
	return nil
}

// RestclientConfig builds a REST client config
func RestclientConfig() (*clientcmdapi.Config, error) {
	glog.Infof(">>> kubeConfig: %s", TestContext.KubeConfig)
	if TestContext.KubeConfig == "" {
		return nil, fmt.Errorf("KubeConfig must be specified to load client config")
	}
	c, err := clientcmd.LoadFromFile(TestContext.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error loading KubeConfig: %v", err.Error())
	}
	return c, nil
}

// LoadConfig builds config from kubernetes config
func LoadConfig() (*rest.Config, error) {
	c, err := RestclientConfig()
	if err != nil {
		if TestContext.KubeConfig == "" {
			return rest.InClusterConfig()
		}
		return nil, err
	}
	return clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: TestContext.Host}}).ClientConfig()
}

// LoadClient builds controller runtime client that accepts any registered type
func LoadClient() (runtimeclient.Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err.Error())
	}
	return runtimeclient.New(config, runtimeclient.Options{})
}

func IsNodeReady(node *corev1.Node) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func WaitUntilAllNodesAreReady(client runtimeclient.Client) error {
	return wait.PollImmediate(1*time.Second, PoolNodesReadyTimeout, func() (bool, error) {
		nodeList := corev1.NodeList{}
		if err := client.List(context.TODO(), &runtimeclient.ListOptions{}, &nodeList); err != nil {
			glog.Errorf("error querying api for nodeList object: %v, retrying...", err)
			return false, nil
		}
		// All nodes needs to be ready
		for _, node := range nodeList.Items {
			if !IsNodeReady(&node) {
				glog.Errorf("Node %q is not ready", node.Name)
				return false, nil
			}
		}
		return true, nil
	})
}

func IsKubemarkProvider(client runtimeclient.Client) (bool, error) {
	key := types.NamespacedName{
		Namespace: TestContext.MachineApiNamespace,
		Name:      "machineapi-kubemark-controllers",
	}
	glog.Infof("Checking if deployment %q exists", key.Name)
	d := &kappsapi.Deployment{}
	if err := client.Get(context.TODO(), key, d); err != nil {
		if strings.Contains(err.Error(), "not found") {
			glog.Infof("Deployment %q does not exists", key.Name)
			return false, nil
		}
		return false, fmt.Errorf("Error querying api for Deployment object %q: %v", key.Name, err)
	}
	glog.Infof("Deployment %q exists", key.Name)
	return true, nil
}
