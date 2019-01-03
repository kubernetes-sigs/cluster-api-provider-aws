package framework

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/tools/clientcmd"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	// apiregistrationclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"

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

var kubeconfig string

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

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
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
	KubeClient         *kubernetes.Clientset
	CAPIClient         *clientset.Clientset
	APIExtensionClient *apiextensionsclientset.Clientset
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
	if kubeconfig == "" {
		return nil, fmt.Errorf("-kubeconfig not set")
	}
	f := &Framework{
		Kubeconfig: kubeconfig,
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
