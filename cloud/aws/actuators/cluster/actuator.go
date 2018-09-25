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
	"fmt"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

type ec2Svc interface {
	ReconcileNetwork(*providerconfigv1.Network) error
}

type codec interface {
	DecodeFromProviderConfig(clusterv1.ProviderConfig, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	codec         codec
	clusterClient client.ClusterInterface
	ec2           ec2Svc
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Codec         codec
	ClusterClient client.ClusterInterface
	EC2Service    ec2Svc
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	return &Actuator{
		codec:         params.Codec,
		clusterClient: params.ClusterClient,
		ec2:           params.EC2Service,
	}, nil
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	glog.Infof("Reconciling cluster %v.", cluster.Name)

	status, err := a.loadProviderStatus(cluster)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	if err := a.ec2.ReconcileNetwork(&status.Network); err != nil {
		return errors.Errorf("unable to reconcile network: %v", err)
	}

	if err := a.storeProviderStatus(cluster, status); err != nil {
		return errors.Errorf("unable to store cluster provider status: %v", err)
	}

	if _, err := a.clusterClient.UpdateStatus(cluster); err != nil {
		return errors.Errorf("failed to update cluster status: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	glog.Infof("Deleting cluster %v.", cluster.Name)
	return fmt.Errorf("TODO: Not yet implemented")
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

func (a *Actuator) storeProviderStatus(cluster *clusterv1.Cluster, status *providerconfigv1.AWSClusterProviderStatus) error {
	raw, err := a.codec.EncodeProviderStatus(status)
	if err != nil {
		return errors.Errorf("failed to encode provider status: %v", err)
	}

	cluster.Status.ProviderStatus = raw
	return nil
}
