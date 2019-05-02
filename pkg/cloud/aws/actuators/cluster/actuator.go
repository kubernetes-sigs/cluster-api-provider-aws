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
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
)

//+kubebuilder:rbac:groups=awsprovider.k8s.io,resources=awsclusterproviderconfigs;awsclusterproviderstatuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.k8s.io,resources=clusters;clusters/status,verbs=get;list;watch;create;update;patch;delete

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	*deployer.Deployer

	client client.ClusterV1alpha1Interface
	log    logr.Logger
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client         client.ClusterV1alpha1Interface
	LoggingContext string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		Deployer: deployer.New(deployer.Params{ScopeGetter: actuators.DefaultScopeGetter}),
		client:   params.Client,
		log:      klogr.New().WithName(params.LoggingContext),
	}
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	a.log.Info("Reconciling Cluster", "cluster-name", cluster.Name, "cluster-namespace", cluster.Namespace)

	scope, err := actuators.NewScope(actuators.ScopeParams{Cluster: cluster, Client: a.client, Logger: a.log})
	if err != nil {
		return errors.Errorf("failed to create scope: %+v", err)
	}

	defer scope.Close()

	ec2svc := ec2.NewService(scope)
	elbsvc := elb.NewService(scope)
	certSvc := certificates.NewService(scope)

	// Store cert material in spec.
	if err := certSvc.ReconcileCertificates(); err != nil {
		return errors.Wrapf(err, "failed to reconcile certificates for cluster %q", cluster.Name)
	}

	if err := ec2svc.ReconcileNetwork(); err != nil {
		return errors.Wrapf(err, "failed to reconcile network for cluster %q", cluster.Name)
	}

	if err := ec2svc.ReconcileBastion(); err != nil {
		return errors.Wrapf(err, "failed to reconcile bastion host for cluster %q", cluster.Name)
	}

	if err := elbsvc.ReconcileLoadbalancers(); err != nil {
		return errors.Wrapf(err, "failed to reconcile load balancers for cluster %q", cluster.Name)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	a.log.Info("Deleting cluster", "cluster-name", cluster.Name, "cluster-namespace", cluster.Namespace)

	scope, err := actuators.NewScope(actuators.ScopeParams{
		Cluster: cluster,
		Client:  a.client,
		Logger:  a.log,
	})
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
		a.log.Error(err, "Error deleting cluster", "cluster-name", cluster.Name, "cluster-namespace", cluster.Namespace)
		return &controllerError.RequeueAfterError{
			RequeueAfter: 5 * 1000 * 1000 * 1000,
		}
	}

	return nil
}
