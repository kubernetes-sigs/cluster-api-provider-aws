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
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
	"sigs.k8s.io/cluster-api/pkg/controller/remote"
)

//+kubebuilder:rbac:groups=awsprovider.k8s.io,resources=awsclusterproviderconfigs;awsclusterproviderstatuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.k8s.io,resources=clusters;clusters/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=,resources=secrets,verbs=create;get;watch;list

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	*deployer.Deployer

	coreClient corev1.CoreV1Interface
	client     client.ClusterV1alpha1Interface
	log        logr.Logger
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	CoreClient     corev1.CoreV1Interface
	Client         client.ClusterV1alpha1Interface
	LoggingContext string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		client:     params.Client,
		coreClient: params.CoreClient,
		log:        klogr.New().WithName(params.LoggingContext),
		Deployer:   deployer.New(deployer.Params{ScopeGetter: actuators.DefaultScopeGetter}),
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

	// Store KubeConfig for Cluster API NodeRef controller to use.
	kubeConfigSecretName := remote.KubeConfigSecretName(cluster.Name)
	secretClient := a.coreClient.Secrets(cluster.Namespace)
	if _, err := secretClient.Get(kubeConfigSecretName, metav1.GetOptions{}); err != nil && apierrors.IsNotFound(err) {
		kubeConfig, err := a.Deployer.GetKubeConfig(cluster, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to get kubeconfig for cluster %q", cluster.Name)
		}

		kubeConfigSecret := &apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: kubeConfigSecretName,
			},
			StringData: map[string]string{
				"value": kubeConfig,
			},
		}

		if _, err := secretClient.Create(kubeConfigSecret); err != nil {
			return errors.Wrapf(err, "failed to create kubeconfig secret for cluster %q", cluster.Name)
		}
	} else if err != nil {
		return errors.Wrapf(err, "failed to get kubeconfig secret for cluster %q", cluster.Name)
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
