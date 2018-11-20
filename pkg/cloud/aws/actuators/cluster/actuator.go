// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
)

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	*deployer.Deployer

	client client.ClusterV1alpha1Interface
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client client.ClusterV1alpha1Interface
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) *Actuator {
	res := &Actuator{
		client: params.Client,
	}

	res.Deployer = deployer.New(services.NewSDKGetter())
	return res
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	klog.Infof("Reconciling cluster %v", cluster.Name)

	scope, err := actuators.NewScope(actuators.ScopeParams{Cluster: cluster, Client: a.client})
	if err != nil {
		return err
	}

	defer scope.Close()

	// Store some config parameters in the status.
	if len(scope.ClusterConfig.CACertificate) == 0 {
		caCert, caKey, err := certificates.NewCertificateAuthority()
		if err != nil {
			return errors.Wrap(err, "Failed to generate a CA for the control plane")
		}

		scope.ClusterConfig.CACertificate = certificates.EncodeCertPEM(caCert)
		scope.ClusterConfig.CAPrivateKey = certificates.EncodePrivateKeyPEM(caKey)
	}

	if err := scope.EC2.ReconcileNetwork(cluster.Name, &scope.ClusterStatus.Network); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	if err := scope.EC2.ReconcileBastion(cluster.Name, scope.ClusterConfig.SSHKeyName, scope.ClusterStatus); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	if err := scope.ELB.ReconcileLoadbalancers(cluster.Name, &scope.ClusterStatus.Network); err != nil {
		return errors.Errorf("unable to reconcile load balancers: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	klog.Infof("Deleting cluster %v.", cluster.Name)

	scope, err := actuators.NewScope(actuators.ScopeParams{Cluster: cluster, Client: a.client})
	if err != nil {
		return err
	}

	defer scope.Close()

	if err := scope.ELB.DeleteLoadbalancers(cluster.Name, &scope.ClusterStatus.Network); err != nil {
		return errors.Errorf("unable to delete load balancers: %v", err)
	}

	if err := scope.EC2.DeleteBastion(cluster.Name, scope.ClusterStatus); err != nil {
		return errors.Errorf("unable to delete bastion: %v", err)
	}

	if err := scope.EC2.DeleteNetwork(cluster.Name, &scope.ClusterStatus.Network); err != nil {
		klog.Errorf("Error deleting cluster %v: %v.", cluster.Name, err)
		return &controllerError.RequeueAfterError{
			RequeueAfter: 5 * 1000 * 1000 * 1000,
		}
	}

	return nil
}
