/*
Copyright YEAR The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package poseidon

import (
	"log"

	extensionsv1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	olympusv1beta1 "github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/olympus/v1beta1"
	listers "github.com/kubernetes-incubator/apiserver-builder/example/pkg/client/listers_generated/olympus/v1beta1"
	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/controller/sharedinformers"
	"k8s.io/api/extensions/v1beta1"
)

// +controller:group=olympus,version=v1beta1,kind=Poseidon,resource=poseidons
type PoseidonControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Poseidon
	lister listers.PoseidonLister

	// lister indexes properties about Deployments
	deploymentLister extensionsv1beta1listers.DeploymentLister
}

// Init initializes the controller and is called by the generated code
func (c *PoseidonControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	i := arguments.GetSharedInformers().Factory.Olympus().V1beta1().Poseidons()
	c.lister = i.Lister()

	di := arguments.GetSharedInformers().KubernetesFactory.Extensions().V1beta1().Deployments()
	c.deploymentLister = di.Lister()

	// For watching Deployments
	arguments.Watch("PoseidonPod", di.Informer(), c.DeploymentToPoseidon)
}

func (c *PoseidonControllerImpl) DeploymentToPoseidon(i interface{}) (string, error) {
	d, _ := i.(*v1beta1.Deployment)
	log.Printf("Deployment update: %v", d.Name)
	if len(d.OwnerReferences) == 1 && d.OwnerReferences[0].Kind == "Poseidon" {
		return d.Namespace + "/" + d.OwnerReferences[0].Name, nil
	} else {
		// Not owned
		return "", nil
	}
}

// Reconcile handles enqueued messages
func (c *PoseidonControllerImpl) Reconcile(u *olympusv1beta1.Poseidon) error {
	// Implement controller logic here
	log.Printf("Running reconcile Poseidon for %s\n", u.Name)
	return nil
}

func (c *PoseidonControllerImpl) Get(namespace, name string) (*olympusv1beta1.Poseidon, error) {
	return c.lister.Poseidons(namespace).Get(name)
}
