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

package conversion

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cabpkv1a2 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
	"sigs.k8s.io/cluster-api/util/certs"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/yaml"
)

const (
	// EtcdCA is the secret name suffix for the Etcd CA
	etcdCA secret.Purpose = "etcd"

	// ServiceAccount is the secret name suffix for the Service Account keys
	serviceAccount secret.Purpose = "sa"

	// FrontProxyCA is the secret name suffix for Front Proxy CA
	frontProxyCA secret.Purpose = "proxy"
)

type ClusterConverter struct {
	oldCluster    *capiv1a1.Cluster
	oldAWSCluster *capav1a1.AWSClusterProviderSpec
}

func NewClusterConverter(cluster *capiv1a1.Cluster) *ClusterConverter {
	return &ClusterConverter{
		oldCluster:    cluster,
		oldAWSCluster: nil,
	}
}

func (c *ClusterConverter) GetCluster(cluster *capiv1a2.Cluster) error {
	if err := capiv1a2.Convert_v1alpha1_Cluster_To_v1alpha2_Cluster(c.oldCluster, cluster, nil); err != nil {
		return errors.WithStack(err)
	}

	ref := corev1.ObjectReference{
		Name:       c.oldCluster.Name,
		Namespace:  c.oldCluster.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSCluster",
	}

	cluster.Spec.InfrastructureRef = &ref

	return nil
}

func (c *ClusterConverter) getOldAWSCluster() (*capav1a1.AWSClusterProviderSpec, error) {
	if c.oldAWSCluster == nil {
		var oldAWSCluster capav1a1.AWSClusterProviderSpec
		if c.oldCluster.Spec.ProviderSpec.Value == nil {
			return nil, nil
		}

		if err := yaml.Unmarshal(c.oldCluster.Spec.ProviderSpec.Value.Raw, &oldAWSCluster); err != nil {
			return nil, errors.Wrap(err, "couldn't decode ProviderSpec")
		}

		c.oldAWSCluster = &oldAWSCluster
	}

	return c.oldAWSCluster, nil
}

func (c *ClusterConverter) GetAWSCluster(cluster *capav1a2.AWSCluster) error {
	oldCluster, err := c.getOldAWSCluster()
	if err != nil {
		return err
	}

	if oldCluster == nil {
		return nil
	}

	if err := capav1a2.Convert_v1alpha1_AWSClusterProviderSpec_To_v1alpha2_AWSClusterSpec(oldCluster, &cluster.Spec, nil); err != nil {
		return err
	}

	cluster.Name = c.oldCluster.Name
	cluster.Namespace = c.oldCluster.Namespace

	return nil
}

func (c *ClusterConverter) GetSecrets(cluster *capiv1a2.Cluster, cfg *cabpkv1a2.KubeadmConfig) ([]*corev1.Secret, error) {
	oldCluster, err := c.getOldAWSCluster()
	if err != nil {
		return []*corev1.Secret{}, err
	}

	certificates := certificates{
		ClusterCA:      convertKeypair(&oldCluster.CAKeyPair),
		EtcdCA:         convertKeypair(&oldCluster.EtcdCAKeyPair),
		FrontProxyCA:   convertKeypair(&oldCluster.FrontProxyCAKeyPair),
		ServiceAccount: convertKeypair(&oldCluster.SAKeyPair),
	}

	return newSecretsFromCertificates(cluster, cfg, &certificates), nil

}

// newSecretsFromCertificates returns a list of secrets, 1 for each certificate.
func newSecretsFromCertificates(cluster *capiv1a2.Cluster, config *cabpkv1a2.KubeadmConfig, c *certificates) []*corev1.Secret {
	return []*corev1.Secret{
		keyPairToSecret(cluster, config, string(secret.ClusterCA), c.ClusterCA),
		keyPairToSecret(cluster, config, string(etcdCA), c.EtcdCA),
		keyPairToSecret(cluster, config, string(frontProxyCA), c.FrontProxyCA),
		keyPairToSecret(cluster, config, string(serviceAccount), c.ServiceAccount),
	}
}

// certificates hold all the certificates necessary for a Kubernetes cluster
type certificates struct {
	ClusterCA      *certs.KeyPair
	EtcdCA         *certs.KeyPair
	FrontProxyCA   *certs.KeyPair
	ServiceAccount *certs.KeyPair
}

func convertKeypair(in *capav1a1.KeyPair) *certs.KeyPair {
	return &certs.KeyPair{
		Cert: in.Cert,
		Key:  in.Key,
	}
}

// KeyPairToSecret creates a Secret from a KeyPair.
func keyPairToSecret(cluster *capiv1a2.Cluster, config *cabpkv1a2.KubeadmConfig, name string, keyPair *certs.KeyPair) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      fmt.Sprintf("%s-%s", cluster.Name, name),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: cabpkv1a2.GroupVersion.String(),
					Kind:       "KubeadmConfig",
					Name:       config.Name,
					UID:        config.UID,
				},
			},
			Labels: map[string]string{
				capiv1a2.MachineClusterLabelName: cluster.Name,
			},
		},
		Data: map[string][]byte{
			secret.TLSKeyDataName: keyPair.Key,
			secret.TLSCrtDataName: keyPair.Cert,
		},
	}
}
