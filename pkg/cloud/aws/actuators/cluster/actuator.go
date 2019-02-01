/*
Copyright 2018 The Kubernetes Authors.

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

package cluster

import (
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
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
	return &Actuator{
		Deployer: deployer.New(deployer.Params{ScopeGetter: actuators.DefaultScopeGetter}),
		client:   params.Client,
	}
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	klog.Infof("Reconciling cluster %v", cluster.Name)

	scope, err := actuators.NewScope(actuators.ScopeParams{Cluster: cluster, Client: a.client})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)
	elbsvc := elb.NewService(scope)

	// Store some config parameters in the status.
	if !scope.ClusterConfig.CAKeyPair.HasCertAndKey() {
		caCert, caKey, err := actuators.GetOrGenerateKeyPair(&scope.ClusterConfig.CAKeyPair, actuators.ClusterCA)
		if err != nil {
			return errors.Wrap(err, "Failed to generate a CA for the control plane")
		}
		scope.ClusterConfig.CAKeyPair = v1alpha1.KeyPair{
			Cert: caCert,
			Key:  caKey,
		}
	}

	if !scope.ClusterConfig.EtcdCAKeyPair.HasCertAndKey() {
		etcdCACert, etcdCAKey, err := actuators.GetOrGenerateKeyPair(&scope.ClusterConfig.EtcdCAKeyPair, actuators.EtcdCA)
		if err != nil {
			return errors.Wrapf(err, "failed to get or generate etcd CA cert")
		}
		scope.ClusterConfig.EtcdCAKeyPair = v1alpha1.KeyPair{
			Cert: etcdCACert,
			Key:  etcdCAKey,
		}
	}

	if !scope.ClusterConfig.FrontProxyCAKeyPair.HasCertAndKey() {
		fpCACert, fpCAKey, err := actuators.GetOrGenerateKeyPair(&scope.ClusterConfig.FrontProxyCAKeyPair, actuators.FrontProxyCA)
		if err != nil {
			return errors.Wrapf(err, "failed to get or generate front-proxy CA cert")
		}
		scope.ClusterConfig.FrontProxyCAKeyPair = v1alpha1.KeyPair{
			Cert: fpCACert,
			Key:  fpCAKey,
		}
	}

	if !scope.ClusterConfig.SAKeyPair.HasCertAndKey() {
		saPub, saKey, err := actuators.GetOrGenerateKeyPair(&scope.ClusterConfig.SAKeyPair, actuators.ServiceAccount)
		if err != nil {
			return errors.Wrapf(err, "failed to get or generate service-account certificates")
		}
		scope.ClusterConfig.SAKeyPair = v1alpha1.KeyPair{
			Cert: saPub,
			Key:  saKey,
		}
	}

	if err := ec2svc.ReconcileNetwork(); err != nil {
		return errors.Errorf("unable to reconcile network: %+v", err)
	}

	if err := ec2svc.ReconcileBastion(); err != nil {
		return errors.Errorf("unable to reconcile network: %+v", err)
	}

	if err := elbsvc.ReconcileLoadbalancers(); err != nil {
		return errors.Errorf("unable to reconcile load balancers: %+v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	klog.Infof("Deleting cluster %v.", cluster.Name)

	scope, err := actuators.NewScope(actuators.ScopeParams{Cluster: cluster, Client: a.client})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)
	elbsvc := elb.NewService(scope)

	if err := elbsvc.DeleteLoadbalancers(); err != nil {
		return errors.Errorf("unable to delete load balancers: %+v", err)
	}

	if err := ec2svc.DeleteBastion(); err != nil {
		return errors.Errorf("unable to delete bastion: %+v", err)
	}

	if err := ec2svc.DeleteNetwork(); err != nil {
		klog.Errorf("Error deleting cluster %v: %v.", cluster.Name, err)
		return &controllerError.RequeueAfterError{
			RequeueAfter: 5 * 1000 * 1000 * 1000,
		}
	}

	return nil
}
