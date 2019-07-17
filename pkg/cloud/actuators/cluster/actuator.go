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
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/certificates"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	clientv1 "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/controller/remote"
	controllerError "sigs.k8s.io/cluster-api/pkg/errors"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const waitForControlPlaneMachineDuration = 15 * time.Second //nolint

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	*deployer.Deployer
	client.Client

	coreClient    corev1.CoreV1Interface
	clusterClient clientv1.ClusterV1alpha2Interface
	log           logr.Logger
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client         client.Client
	CoreClient     corev1.CoreV1Interface
	ClusterClient  clientv1.ClusterV1alpha2Interface
	LoggingContext string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		clusterClient: params.ClusterClient,
		coreClient:    params.CoreClient,
		log:           klogr.New().WithName(params.LoggingContext),
		Deployer:      deployer.New(deployer.Params{ClusterScopeGetter: scope.DefaultClusterScopeGetter}),
	}
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	log := a.log.WithValues("cluster-name", cluster.Name, "cluster-namespace", cluster.Namespace)
	log.Info("Reconciling Cluster")

	scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: cluster,
		Logger:  a.log,
	})
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

	if cluster.Annotations == nil {
		cluster.Annotations = make(map[string]string)
	}
	cluster.Annotations[v1alpha2.AnnotationClusterInfrastructureReady] = v1alpha2.ValueReady

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

	// If the control plane is ready, try to delete the control plane configmap lock, if it exists, and return.
	if cluster.Annotations[v1alpha2.AnnotationControlPlaneReady] == v1alpha2.ValueReady {
		configMapName := scope.ControlPlaneConfigMapName()
		log.Info("Checking for existence of control plane configmap lock", "configmap-name", configMapName)

		_, err := a.coreClient.ConfigMaps(cluster.Namespace).Get(configMapName, metav1.GetOptions{})
		switch {
		case apierrors.IsNotFound(err):
			// It doesn't exist - no-op
		case err != nil:
			return errors.Wrapf(err, "Error retrieving control plane configmap lock %q", configMapName)
		default:
			if err := a.coreClient.ConfigMaps(cluster.Namespace).Delete(configMapName, nil); err != nil {
				return errors.Wrapf(err, "Error deleting control plane configmap lock %q", configMapName)
			}
		}

		// Nothing more to reconcile - return early.
		return nil
	}

	log.Info("Cluster does not have ready annotation - checking for ready control plane machines")

	machineList := &clusterv1.MachineList{}
	if err := a.List(context.Background(), machineList, scope.ListOptionsLabelSelector()); err != nil {
		return errors.Wrapf(err, "failed to retrieve machines in cluster %q", cluster.Name)
	}

	controlPlaneMachines := util.GetControlPlaneMachinesFromList(machineList)

	machineReady := false
	for _, machine := range controlPlaneMachines {
		if machine.Status.NodeRef != nil {
			machineReady = true
			break
		}
	}

	if !machineReady {
		log.Info("No control plane machines are ready - requeuing cluster")
		return &controllerError.RequeueAfterError{RequeueAfter: waitForControlPlaneMachineDuration}
	}

	log.Info("Setting cluster ready annotation")
	cluster.Annotations[v1alpha2.AnnotationControlPlaneReady] = v1alpha2.ValueReady

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	a.log.Info("Deleting cluster", "cluster-name", cluster.Name, "cluster-namespace", cluster.Namespace)

	scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: cluster,
		Client:  a.Client,
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
			RequeueAfter: 5 * time.Second,
		}
	}

	return nil
}
