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
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/instrumentation"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	service "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services"
	certificates "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/certificates"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	elbsvc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/elb"
	ssmsvc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ssm"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	controllerError "sigs.k8s.io/cluster-api/pkg/controller/error"
)

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	codec          codec
	clustersGetter client.ClustersGetter
	servicesGetter service.Getter
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Codec          codec
	ClustersGetter client.ClustersGetter
	ServicesGetter service.Getter
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	res := &Actuator{
		codec:          params.Codec,
		clustersGetter: params.ClustersGetter,
		servicesGetter: params.ServicesGetter,
	}

	if res.servicesGetter == nil {
		res.servicesGetter = new(defaultServicesGetter)
	}

	return res, nil
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) (reterr error) {
	ctx := context.Background()

	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "cluster", "Reconcile"),
	)
	defer span.End()

	glog.V(2).Infof("Reconciling cluster %v.", cluster.Name)

	// Get a cluster api client for the namespace of the cluster.
	clusterClient := a.clustersGetter.Clusters(cluster.Namespace)

	// Load provider config.
	config, err := a.loadProviderConfig(cluster)
	if err != nil {
		return errors.Errorf("failed to load cluster provider config: %v", err)
	}

	// Load provider status.
	status, err := a.loadProviderStatus(cluster)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	// Store some config parameters in the status.
	status.Region = config.Region

	// Always defer storing the cluster status. In case any of the calls below fails or returns an error
	// the cluster state might have partial changes that should be stored.
	defer func() {
		// TODO(vincepri): remove this after moving to tag-discovery based approach.
		if err := a.storeProviderStatus(clusterClient, cluster, status); err != nil {
			glog.Errorf("failed to store provider status for cluster %q: %v", cluster.Name, err)
		}
	}()

	// Create new aws session.
	sess := a.servicesGetter.Session(config)

	// Load ec2 client.
	ec2 := a.servicesGetter.EC2(sess, instrumentation.AWSInstrumentedConfig())

	// Load ssm client.

	ssm := a.servicesGetter.SSM(sess, instrumentation.AWSInstrumentedConfig())

	caCert, caKey, err := certificates.NewCertificateAuthority()
	if err != nil {
		return errors.Wrap(err, "Failed to generate a CA for the control plane")
	}

	err = ssm.ReconcileParameter(ctx, cluster.Name,
		certificates.SSMCACertificatePath,
		string(certificates.EncodeCertPEM(caCert)))

	if err != nil {
		return errors.Wrap(err, "failed to put CA certificate in SSM Parameter store")
	}

	err = ssm.ReconcileParameter(ctx, cluster.Name,
		certificates.SSMCAPrivateKeyPath,
		string(certificates.EncodePrivateKeyPEM(caKey)))

	if err != nil {
		return errors.Wrap(err, "failed to put CA private key in SSM Parameter store")
	}

	if err := ec2.ReconcileNetwork(ctx, cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	if err := ec2.ReconcileBastion(ctx, cluster.Name, config.SSHKeyName, status); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	// Load elb client.
	elb := a.servicesGetter.ELB(sess, instrumentation.AWSInstrumentedConfig())

	if err := elb.ReconcileLoadbalancers(ctx, cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to reconcile load balancers: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	ctx := context.Background()

	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("actuators", "cluster", "Delete"),
	)
	defer span.End()
	glog.Infof("Deleting cluster %v.", cluster.Name)

	// Get a cluster api client for the namespace of the cluster.
	clusterClient := a.clustersGetter.Clusters(cluster.Namespace)

	// Load provider config.
	config, err := a.loadProviderConfig(cluster)
	if err != nil {
		return errors.Errorf("failed to load cluster provider config: %v", err)
	}

	// Load provider status.
	status, err := a.loadProviderStatus(cluster)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	// Store some config parameters in the status.
	status.Region = config.Region

	// Always defer storing the cluster status. In case any of the calls below fails or returns an error
	// the cluster state might have partial changes that should be stored.
	defer func() {
		// TODO(vincepri): remove this after moving to tag-discovery based approach.
		if err := a.storeProviderStatus(clusterClient, cluster, status); err != nil {
			glog.Errorf("failed to store provider status for cluster %q: %v", cluster.Name, err)
		}
	}()

	// Create new aws session.
	sess := a.servicesGetter.Session(config)

	// Load ec2 client.
	ec2 := a.servicesGetter.EC2(sess, instrumentation.AWSInstrumentedConfig())

	// Load elb client.
	elb := a.servicesGetter.ELB(sess, instrumentation.AWSInstrumentedConfig())

	ssm := a.servicesGetter.SSM(sess, instrumentation.AWSInstrumentedConfig())

	err = ssm.DeleteParameter(ctx, cluster.Name,
		certificates.SSMCACertificatePath)

	if err != nil {
		return errors.Wrap(err, "failed to delete CA certificate from SSM Parameter store")
	}

	err = ssm.DeleteParameter(ctx, cluster.Name,
		certificates.SSMCAPrivateKeyPath)

	if err != nil {
		return errors.Wrap(err, "failed to delete CA private key from SSM Parameter store")
	}

	if err := elb.DeleteLoadbalancers(ctx, cluster.Name, &status.Network); err != nil {
		return errors.Errorf("unable to delete load balancers: %v", err)
	}

	if err := ec2.DeleteBastion(ctx, cluster.Name, status); err != nil {
		return errors.Errorf("unable to delete bastion: %v", err)
	}

	if err := ec2.DeleteNetwork(ctx, cluster.Name, &status.Network); err != nil {
		glog.Errorf("Error deleting cluster %v: %v.", cluster.Name, err)
		return &controllerError.RequeueAfterError{
			RequeueAfter: 5 * 1000 * 1000 * 1000,
		}
	}

	return nil
}

func (a *Actuator) loadProviderConfig(cluster *clusterv1.Cluster) (*providerconfigv1.AWSClusterProviderConfig, error) {
	providerConfig := &providerconfigv1.AWSClusterProviderConfig{}
	err := a.codec.DecodeFromProviderConfig(cluster.Spec.ProviderConfig, providerConfig)
	return providerConfig, err
}

func (a *Actuator) loadProviderStatus(cluster *clusterv1.Cluster) (*providerconfigv1.AWSClusterProviderStatus, error) {
	providerStatus := &providerconfigv1.AWSClusterProviderStatus{}
	err := a.codec.DecodeProviderStatus(cluster.Status.ProviderStatus, providerStatus)
	return providerStatus, err
}

func (a *Actuator) storeProviderStatus(clusterClient client.ClusterInterface, cluster *clusterv1.Cluster, status *providerconfigv1.AWSClusterProviderStatus) error {
	raw, err := a.codec.EncodeProviderStatus(status)
	if err != nil {
		return errors.Errorf("failed to encode provider status: %v", err)
	}

	cluster.Status.ProviderStatus = raw
	if _, err := clusterClient.UpdateStatus(cluster); err != nil {
		return err
	}

	return nil
}

type defaultServicesGetter struct{}

func (d *defaultServicesGetter) Session(clusterConfig *providerconfigv1.AWSClusterProviderConfig) *session.Session {
	return session.Must(session.NewSession(instrumentation.AWSInstrumentedConfig().WithRegion(clusterConfig.Region)))
}

func (d *defaultServicesGetter) EC2(session *session.Session, c *aws.Config) service.EC2Interface {
	return ec2svc.NewService(ec2.New(session, instrumentation.AWSInstrumentedConfig()))
}

func (d *defaultServicesGetter) ELB(session *session.Session, c *aws.Config) service.ELBInterface {
	return elbsvc.NewService(elb.New(session, instrumentation.AWSInstrumentedConfig()))
}

func (d *defaultServicesGetter) SSM(session *session.Session, c *aws.Config) service.SSMInterface {
	return ssmsvc.NewService(ssm.New(session, instrumentation.AWSInstrumentedConfig()))
}
