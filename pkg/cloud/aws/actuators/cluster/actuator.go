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
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
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

	clustersGetter client.ClustersGetter
	servicesGetter services.Getter
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClustersGetter client.ClustersGetter
	ServicesGetter services.Getter
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) *Actuator {
	res := &Actuator{
		clustersGetter: params.ClustersGetter,
		servicesGetter: params.ServicesGetter,
	}

	if res.servicesGetter == nil {
		res.servicesGetter = services.NewSDKGetter()
	}

	res.Deployer = deployer.New(res.servicesGetter)
	return res
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) (reterr error) {
	klog.Infof("Reconciling cluster %v.", cluster.Name)

	// Load provider config.
	config, err := providerv1.ClusterConfigFromProviderConfig(cluster.Spec.ProviderConfig)
	if err != nil {
		return errors.Errorf("failed to load cluster provider config: %v", err)
	}

	// Load provider status.
	status, err := providerv1.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	defer func() {
		if err := a.storeClusterConfig(cluster, config); err != nil {
			klog.Errorf("failed to store provider config for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
		}

		if err := a.storeClusterStatus(cluster, status); err != nil {
			klog.Errorf("failed to store provider status for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
		}
	}()

	// Store some config parameters in the status.
	status.Region = config.Region

	if len(config.CACertificate) == 0 {
		caCert, caKey, err := certificates.NewCertificateAuthority()
		if err != nil {
			return errors.Wrap(err, "Failed to generate a CA for the control plane")
		}

		config.CACertificate = certificates.EncodeCertPEM(caCert)
		config.CAPrivateKey = certificates.EncodePrivateKeyPEM(caKey)
	}

	// Create new aws session.
	sess := a.servicesGetter.Session(config)

	// Load ec2 client.
	ec2 := a.servicesGetter.EC2(sess)

	if err := ec2.ReconcileNetwork(cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	if err := ec2.ReconcileBastion(cluster.Name, config.SSHKeyName, status); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	// Load elb client.
	elb := a.servicesGetter.ELB(sess)

	if err := elb.ReconcileLoadbalancers(cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to reconcile load balancers: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	klog.Infof("Deleting cluster %v.", cluster.Name)

	// Load provider config.
	config, err := providerv1.ClusterConfigFromProviderConfig(cluster.Spec.ProviderConfig)
	if err != nil {
		return errors.Errorf("failed to load cluster provider config: %v", err)
	}

	// Load provider status.
	status, err := providerv1.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	// Store some config parameters in the status.
	status.Region = config.Region

	// Create new aws session.
	sess := a.servicesGetter.Session(config)

	// Load ec2 client.
	ec2 := a.servicesGetter.EC2(sess)

	// Load elb client.
	elb := a.servicesGetter.ELB(sess)

	if err := elb.DeleteLoadbalancers(cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to delete load balancers: %v", err)
	}

	if err := ec2.DeleteBastion(cluster.Name, status); err != nil {
		return errors.Errorf("unable to delete bastion: %v", err)
	}

	if err := ec2.DeleteNetwork(cluster.Name, &status.Network); err != nil {
		klog.Errorf("Error deleting cluster %v: %v.", cluster.Name, err)
		return &controllerError.RequeueAfterError{
			RequeueAfter: 5 * 1000 * 1000 * 1000,
		}
	}

	return nil
}

func (a *Actuator) storeClusterConfig(cluster *clusterv1.Cluster, config *providerv1.AWSClusterProviderConfig) error {
	clusterClient := a.clustersGetter.Clusters(cluster.Namespace)

	ext, err := providerv1.EncodeClusterConfig(config)
	if err != nil {
		return err
	}

	cluster.Spec.ProviderConfig.Value = ext

	if _, err := clusterClient.Update(cluster); err != nil {
		return err
	}

	return nil
}

func (a *Actuator) storeClusterStatus(cluster *clusterv1.Cluster, status *providerv1.AWSClusterProviderStatus) error {
	clusterClient := a.clustersGetter.Clusters(cluster.Namespace)

	ext, err := providerv1.EncodeClusterStatus(status)
	if err != nil {
		return err
	}

	cluster.Status.ProviderStatus = ext

	if _, err := clusterClient.UpdateStatus(cluster); err != nil {
		return err
	}

	return nil
}
