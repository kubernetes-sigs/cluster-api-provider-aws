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

package machine

import (
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	apicorev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// ControlPlaneInitLocker provides a locking mechanism for cluster initialization.
type ControlPlaneInitLocker interface {
	// GetOrCreateInitLock will attempt to create a lock specifying the given machine is the one
	// that will init the cluster. If no such lock exists, or the existing lock matches the
	// specified machine, return "true" and the idempotency token.
	// otherwise, return false
	AcquireWithToken(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, string)
}

// controlPlaneInitLocker uses a ConfigMap to synchronize cluster initialization.
type controlPlaneInitLocker struct {
	log             logr.Logger
	configMapClient corev1.ConfigMapsGetter
	random          func() string
}

var _ ControlPlaneInitLocker = &controlPlaneInitLocker{}

func newControlPlaneInitLocker(log logr.Logger, configMapClient corev1.ConfigMapsGetter) *controlPlaneInitLocker {
	return &controlPlaneInitLocker{
		log:             log,
		configMapClient: configMapClient,
		random: func() string {
			return uuid.New().String()
		},
	}
}

func (l *controlPlaneInitLocker) AcquireWithToken(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, string) {
	configMapName := actuators.ControlPlaneConfigMapName(cluster)
	log := l.log.WithValues("namespace", cluster.Namespace, "cluster-name", cluster.Name, "configmap-name", configMapName, "machine-name", machine.Name)

	lockInfo, err := l.getLock(cluster.Namespace, configMapName)
	if err != nil {
		log.Error(err, "Error checking for control plane configmap lock existence")
		return false, ""
	}

	if lockInfo != nil {
		if lockInfo.MachineName == machine.Name {
			return true, lockInfo.IdempotencyToken
		}

		log.Info("waiting on on another machine to initialize", "init-machine", lockInfo.MachineName)
		return false, ""
	}

	token := l.random()
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

	lockInfo = &LockInformation{
		MachineName:      machine.Name,
		IdempotencyToken: token,
	}

	if err := WriteLockInfo(controlPlaneConfigMap, lockInfo); err != nil {
		log.Error(err, "Failed to add lock information to control plane config map")
		return false, ""
	}

	log.Info("Attempting to create control plane configmap lock")

	if _, err := l.configMapClient.ConfigMaps(cluster.Namespace).Create(controlPlaneConfigMap); err != nil {
		if apierrors.IsAlreadyExists(err) {
			// Someone else beat us to it
			log.Info("Control plane configmap lock already exists")
		} else {
			log.Error(err, "Error creating control plane configmap lock")
		}

		// Unable to acquire
		return false, ""
	}

	// Successfully acquired
	return true, token
}

func (l *controlPlaneInitLocker) getLock(namespace, name string) (*LockInformation, error) {
	cm, err := l.configMapClient.ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ReadLockInfo(cm)
}

func (l *controlPlaneInitLocker) configMapExists(namespace, name string) (bool, error) {
	_, err := l.configMapClient.ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	return err == nil, err
}
