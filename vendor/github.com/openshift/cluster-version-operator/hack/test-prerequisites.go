package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/ghodss/yaml"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
)

// main installs the CV CRD to a cluster for integration testing.
func main() {
	log.SetFlags(0)
	kcfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{})
	cfg, err := kcfg.ClientConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	client := apiext.NewForConfigOrDie(cfg)
	for _, path := range []string{
		"install/0000_00_cluster-version-operator_01_clusterversion.crd.yaml",
		"install/0000_00_cluster-version-operator_01_clusteroperator.crd.yaml",
	} {
		var name string
		err := wait.PollImmediate(time.Second, 30*time.Second, func() (bool, error) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Unable to read %s: %v", path, err)
			}
			var crd v1beta1.CustomResourceDefinition
			if err := yaml.Unmarshal(data, &crd); err != nil {
				log.Fatalf("Unable to parse CRD %s: %v", path, err)
			}
			name = crd.Name
			_, err = client.Apiextensions().CustomResourceDefinitions().Create(&crd)
			if errors.IsAlreadyExists(err) {
				return true, nil
			}
			if err != nil {
				log.Printf("error: failed creating CRD %s: %v", name, err)
				return false, nil
			}
			log.Printf("Installed %s CRD", crd.Name)
			return true, nil
		})
		if err != nil {
			log.Fatalf("Could not install %s CRD: %v", name, err)
		}
	}
}
