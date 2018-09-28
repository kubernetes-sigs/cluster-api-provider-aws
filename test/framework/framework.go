package framework

import (
	"flag"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"k8s.io/client-go/kubernetes"
	apiregistrationclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var kubeconfig string

// ClusterID set by -cluster-id flag
var ClusterID string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&ClusterID, "cluster-id", "", "cluster ID")
	flag.Parse()
}

// Framework supports common operations used by tests
type Framework struct {
	KubeClient            *kubernetes.Clientset
	CAPIClient            *clientset.Clientset
	APIRegistrationClient *apiregistrationclientset.Clientset
	Kubeconfig            string
}

// NewFramework setups a new framework
func NewFramework() *Framework {
	f := &Framework{
		Kubeconfig: kubeconfig,
	}

	BeforeEach(f.BeforeEach)

	return f
}

// BeforeEach to be run before each spec responsible for building various clientsets
func (f *Framework) BeforeEach() {
	By("Creating a kubernetes client")

	if f.KubeClient == nil {
		config, err := controller.GetConfig(f.Kubeconfig)
		Expect(err).NotTo(HaveOccurred())
		f.KubeClient, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	}

	if f.CAPIClient == nil {
		config, err := controller.GetConfig(f.Kubeconfig)
		Expect(err).NotTo(HaveOccurred())
		f.CAPIClient, err = clientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	}

	if f.APIRegistrationClient == nil {
		config, err := controller.GetConfig(f.Kubeconfig)
		Expect(err).NotTo(HaveOccurred())
		f.APIRegistrationClient, err = apiregistrationclientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	}

}

// IgnoreNotFoundErr ignores not found errors in case resource
// that does not exist is to be deleted
func IgnoreNotFoundErr(err error) {
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

// SigKubeDescribe is a wrapper function for ginkgo describe.  Adds namespacing.
func SigKubeDescribe(text string, body func()) bool {
	return Describe("[sigs.k8s.io] "+text, body)
}
