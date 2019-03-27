package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"bytes"
	"reflect"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// ClusterConfigNamespace is the namespace containing the cluster config
	ClusterConfigNamespace = "kube-system"
	// ClusterConfigName is the name of the cluster config configmap
	ClusterConfigName = "cluster-config-v1"
	// InstallConfigKey is the key in the cluster config configmap containing yaml installConfig data
	InstallConfigKey = "install-config"
	// AWSPlatformType is used to install on AWS
	AWSProvider = Provider("aws")
	// LibvirtPlatformType is used to install of libvirt
	LibvirtProvider = Provider("libvirt")
	// OpenStackPlatformType is used to install on OpenStack
	OpenStackProvider = Provider("openstack")
	// KubemarkPlatformType is used to install on Kubemark
	KubemarkProvider = Provider("kubemark")
)

type Provider string

// OperatorConfig contains configuration for MAO
type OperatorConfig struct {
	TargetNamespace string `json:"targetNamespace"`
	Controllers     Controllers
}

type Controllers struct {
	Provider           string
	NodeLink           string
	MachineHealthCheck string
}

// Images allows build systems to inject images for MAO components
type Images struct {
	MachineAPIOperator            string `json:"machineAPIOperator"`
	ClusterAPIControllerAWS       string `json:"clusterAPIControllerAWS"`
	ClusterAPIControllerOpenStack string `json:"clusterAPIControllerOpenStack"`
	ClusterAPIControllerLibvirt   string `json:"clusterAPIControllerLibvirt"`
	ClusterAPIControllerKubemark  string `json:"clusterAPIControllerKubemark"`
}

// InstallConfig contains the mao relevant config coming from the install config, i.e provider
type InstallConfig struct {
	InstallPlatform `json:"platform"`
}

// InstallPlatform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set
type InstallPlatform struct {
	// AWS is the configuration used when running on AWS
	AWS interface{} `json:"aws,omitempty"`

	// Libvirt is the configuration used when running on libvirt
	Libvirt interface{} `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when running on OpenStack
	OpenStack interface{} `json:"openstack,omitempty"`

	// Kubemark is the configuration used when running with Kubemark
	Kubemark interface{} `json:"kubemark,omitempty"`
}

func getInstallConfig(client kubernetes.Interface) (*InstallConfig, error) {
	cm, err := client.CoreV1().ConfigMaps(ClusterConfigNamespace).Get(ClusterConfigName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting clusterconfig %s/%s: %v", ClusterConfigNamespace, ClusterConfigName, err)
	}

	return getInstallConfigFromClusterConfig(cm)
}

// getInstallConfigFromClusterConfig builds an install config from the cluster config.
func getInstallConfigFromClusterConfig(clusterConfig *corev1.ConfigMap) (*InstallConfig, error) {
	icYaml, ok := clusterConfig.Data[InstallConfigKey]
	if !ok {
		return nil, fmt.Errorf("missing %q in configmap", InstallConfigKey)
	}
	var ic InstallConfig
	if err := yaml.Unmarshal([]byte(icYaml), &ic); err != nil {
		return nil, fmt.Errorf("invalid InstallConfig: %v yaml: %s", err, icYaml)
	}
	return &ic, nil
}

func getProviderFromInstallConfig(installConfig *InstallConfig) (Provider, error) {
	v := reflect.ValueOf(installConfig.InstallPlatform)
	var nonNilFields int

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != nil {
			nonNilFields = nonNilFields + 1
		}
		if nonNilFields > 1 {
			return "", fmt.Errorf("more than one platform provider given")
		}
	}

	if installConfig.AWS != nil {
		return AWSProvider, nil
	}
	if installConfig.Libvirt != nil {
		return LibvirtProvider, nil
	}
	if installConfig.OpenStack != nil {
		return OpenStackProvider, nil
	}
	if installConfig.Kubemark != nil {
		return KubemarkProvider, nil
	}
	return "", fmt.Errorf("no platform provider found on install config")
}

func getImagesFromJSONFile(filePath string) (*Images, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var i Images
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return &i, nil
}

func getProviderControllerFromImages(provider Provider, images Images) (string, error) {
	switch provider {
	case AWSProvider:
		return images.ClusterAPIControllerAWS, nil
	case LibvirtProvider:
		return images.ClusterAPIControllerLibvirt, nil
	case OpenStackProvider:
		return images.ClusterAPIControllerOpenStack, nil
	case KubemarkProvider:
		return images.ClusterAPIControllerKubemark, nil
	}
	return "", fmt.Errorf("not known platform provider given %s", provider)
}

func getMachineAPIOperatorFromImages(images Images) (string, error) {
	if images.MachineAPIOperator == "" {
		return "", fmt.Errorf("failed gettingMachineAPIOperator image. It is empty")
	}
	return images.MachineAPIOperator, nil
}

// PopulateTemplate receives a template file path and renders its content populated with the config
func PopulateTemplate(config *OperatorConfig, path string) ([]byte, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading file, %v", err)
	}

	buf := &bytes.Buffer{}
	tmpl, err := template.New("").Option("missingkey=error").Parse(string(data))
	if err != nil {
		return nil, err
	}

	tmplData := struct {
		OperatorConfig
	}{
		OperatorConfig: *config,
	}

	if err := tmpl.Execute(buf, tmplData); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
