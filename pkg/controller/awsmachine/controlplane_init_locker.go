/*
Copyright 2019 The Kubernetes Authors.

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

package awsmachine

import (
	"github.com/go-logr/logr"
	apicorev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
)

// ControlPlaneInitLocker provides a locking mechanism for cluster initialization.
type ControlPlaneInitLocker interface {
	// Acquire returns true if it acquires the lock for the cluster.
	Acquire(cluster *clusterv1.Cluster) bool
}

// controlPlaneInitLocker uses a ConfigMap to synchronize cluster initialization.
type controlPlaneInitLocker struct {
	log             logr.Logger
	configMapClient corev1.ConfigMapsGetter
}

var _ ControlPlaneInitLocker = &controlPlaneInitLocker{}

func newControlPlaneInitLocker(log logr.Logger, configMapClient corev1.ConfigMapsGetter) *controlPlaneInitLocker {
	return &controlPlaneInitLocker{
		log:             log,
		configMapClient: configMapClient,
	}
}

func (l *controlPlaneInitLocker) Acquire(cluster *clusterv1.Cluster) bool {
	configMapName := actuators.ControlPlaneConfigMapName(cluster)
	log := l.log.WithValues("namespace", cluster.Namespace, "cluster-name", cluster.Name, "configmap-name", configMapName)

	exists, err := l.configMapExists(cluster.Namespace, configMapName)
	if err != nil {
		log.Error(err, "Error checking for control plane configmap lock existence")
		return false
	}
	if exists {
		return false
	}

	controlPlaneConfigMap := &apicorev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      configMapName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: cluster.APIVersion,
					Kind:       cluster.Kind,
					Name:       cluster.Name,
					UID:        cluster.UID,
				},
			},
		},
	}

	log.Info("Attempting to create control plane configmap lock")
	_, err = l.configMapClient.ConfigMaps(cluster.Namespace).Create(controlPlaneConfigMap)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			// Someone else beat us to it
			log.Info("Control plane configmap lock already exists")
		} else {
			log.Error(err, "Error creating control plane configmap lock")
		}

		// Unable to acquire
		return false
	}

	// Successfully acquired
	return true
}

func (l *controlPlaneInitLocker) configMapExists(namespace, name string) (bool, error) {
	_, err := l.configMapClient.ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	return err == nil, err
}
