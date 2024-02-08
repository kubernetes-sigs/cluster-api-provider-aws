/*
Copyright 2022 The Kubernetes Authors.

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

package secretsmanager

import (
	"crypto/rand"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/secretsmanager/mock_secretsmanageriface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestServiceCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	generateBytes := func(g *WithT, count int64) []byte {
		token := make([]byte, count)
		_, err := rand.Read(token)
		g.Expect(err).NotTo(HaveOccurred())
		return token
	}

	sortTagsByKey := func(tags []*secretsmanager.Tag) {
		sort.Slice(tags, func(i, j int) bool {
			return *(tags[i].Key) < *(tags[j].Key)
		})
	}

	expectedTags := []*secretsmanager.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("infra-cluster"),
		},
		{
			Key:   aws.String("kubernetes.io/cluster/test"),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test"),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
			Value: aws.String("node"),
		},
	}

	check := func(g *WithT, actualPrefix string, expectedPrefix string, err error, wantErr bool) {
		g.Expect(actualPrefix).Should(HavePrefix(expectedPrefix))
		if wantErr {
			g.Expect(err).To(HaveOccurred())
			return
		}
		g.Expect(err).NotTo(HaveOccurred())
	}

	tests := []struct {
		name           string
		bytesCount     int64
		secretPrefix   string
		expectedPrefix string
		wantErr        bool
		expect         func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder)
	}{
		{
			name:           "Should not store data in secret manager if data is having zero bytes",
			bytesCount:     0,
			secretPrefix:   "/awsprefix",
			expectedPrefix: "/awsprefix",
			wantErr:        false,
			expect: func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.CreateSecret(gomock.Any()).Times(0)
			},
		},
		{
			name:           "Should store data in secret manager if data is having non-zero bytes",
			bytesCount:     10000,
			secretPrefix:   "prefix",
			expectedPrefix: "prefix",
			wantErr:        false,
			expect: func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).MinTimes(1).Return(&secretsmanager.CreateSecretOutput{}, nil).Do(
					func(createSecretInput *secretsmanager.CreateSecretInput) {
						g.Expect(*(createSecretInput.Name)).To(HavePrefix("prefix-"))
						sortTagsByKey(createSecretInput.Tags)
						g.Expect(createSecretInput.Tags).To(Equal(expectedTags))
					},
				)
			},
		},
		{
			name:           "Should not retry if non-retryable error occurred while storing data in secret manager",
			bytesCount:     10,
			secretPrefix:   "/prefix",
			expectedPrefix: "/prefix",
			wantErr:        true,
			expect: func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(nil, &secretsmanager.InternalServiceError{}).Do(
					func(createSecretInput *secretsmanager.CreateSecretInput) {
						g.Expect(*(createSecretInput.Name)).To(HavePrefix("/prefix-"))
						sortTagsByKey(createSecretInput.Tags)
						g.Expect(createSecretInput.Tags).To(Equal(expectedTags))
					},
				)
			},
		},
		{
			name:           "Should retry if retryable error occurred while storing data in secret manager",
			bytesCount:     10,
			secretPrefix:   "",
			expectedPrefix: "aws.cluster.x-k8s.io",
			wantErr:        false,
			expect: func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(nil, &secretsmanager.InvalidRequestException{})
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(nil, &secretsmanager.ResourceNotFoundException{})
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(&secretsmanager.CreateSecretOutput{}, nil)
			},
		},
		{
			name:           "Should delete and retry creation if resource already exists while storing data in secret manager",
			bytesCount:     10,
			secretPrefix:   "",
			expectedPrefix: "aws.cluster.x-k8s.io",
			wantErr:        false,
			expect: func(g *WithT, m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(nil, &secretsmanager.ResourceExistsException{})
				m.DeleteSecret(gomock.AssignableToTypeOf(&secretsmanager.DeleteSecretInput{})).Return(&secretsmanager.DeleteSecretOutput{}, nil)
				m.CreateSecret(gomock.AssignableToTypeOf(&secretsmanager.CreateSecretInput{})).Return(&secretsmanager.CreateSecretOutput{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := getClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			secretManagerClientMock := mock_secretsmanageriface.NewMockSecretsManagerAPI(mockCtrl)
			tt.expect(g, secretManagerClientMock.EXPECT())
			s := NewService(clusterScope)
			s.SecretsManagerClient = secretManagerClientMock
			ms, err := getMachineScope(client, clusterScope)
			g.Expect(err).NotTo(HaveOccurred())
			ms.SetSecretPrefix(tt.secretPrefix)
			data := generateBytes(g, tt.bytesCount)

			prefix, _, err := s.Create(ms, data)
			check(g, prefix, tt.expectedPrefix, err, tt.wantErr)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		secretCount int32
		expect      func(m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder)
		check       func(*WithT, error)
	}{
		{
			name:        "Should not call AWS when secret count has zero value",
			secretCount: 0,
			check: func(g *WithT, err error) {
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name:        "Should not return error when delete is successful",
			secretCount: 1,
			expect: func(m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.DeleteSecret(gomock.Eq(&secretsmanager.DeleteSecretInput{
					SecretId:                   aws.String("prefix-0"),
					ForceDeleteWithoutRecovery: aws.Bool(true),
				})).Return(&secretsmanager.DeleteSecretOutput{}, nil)
			},
			check: func(g *WithT, err error) {
				g.Expect(err).NotTo(HaveOccurred())
			},
		},
		{
			name:        "Should return all errors except not found errors",
			secretCount: 3,
			expect: func(m *mock_secretsmanageriface.MockSecretsManagerAPIMockRecorder) {
				m.DeleteSecret(gomock.Eq(&secretsmanager.DeleteSecretInput{
					SecretId:                   aws.String("prefix-0"),
					ForceDeleteWithoutRecovery: aws.Bool(true),
				})).Return(nil, awserrors.NewFailedDependency("failed dependency"))
				m.DeleteSecret(gomock.Eq(&secretsmanager.DeleteSecretInput{
					SecretId:                   aws.String("prefix-1"),
					ForceDeleteWithoutRecovery: aws.Bool(true),
				})).Return(nil, awserrors.NewNotFound("not found"))
				m.DeleteSecret(gomock.Eq(&secretsmanager.DeleteSecretInput{
					SecretId:                   aws.String("prefix-2"),
					ForceDeleteWithoutRecovery: aws.Bool(true),
				})).Return(nil, awserrors.NewConflict("new conflict"))
			},
			check: func(g *WithT, err error) {
				g.Expect(err).ToNot(BeNil())
				g.Expect((err.Error())).To(Equal("[failed dependency, new conflict]"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			client := setupClient()
			clusterScope, err := getClusterScope(client)
			g.Expect(err).NotTo(HaveOccurred())

			secretManagerClientMock := mock_secretsmanageriface.NewMockSecretsManagerAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(secretManagerClientMock.EXPECT())
			}
			s := NewService(clusterScope)
			s.SecretsManagerClient = secretManagerClientMock
			ms, err := getMachineScope(client, clusterScope)
			g.Expect(err).NotTo(HaveOccurred())

			ms.SetSecretPrefix("prefix")
			ms.SetSecretCount(tt.secretCount)
			err = s.Delete(ms)
			tt.check(g, err)
		})
	}
}

func setupClient() client.Client {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	return fake.NewClientBuilder().WithScheme(scheme).Build()
}

func getClusterScope(client client.Client) (*scope.ClusterScope, error) {
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	return scope.NewClusterScope(scope.ClusterScopeParams{
		Client:     client,
		Cluster:    cluster,
		AWSCluster: &infrav1.AWSCluster{},
	})
}

func getMachineScope(client client.Client, clusterScope *scope.ClusterScope) (*scope.MachineScope, error) {
	return scope.NewMachineScope(scope.MachineScopeParams{
		Client:       client,
		ControlPlane: &unstructured.Unstructured{},
		Cluster:      clusterScope.Cluster,
		Machine:      &clusterv1.Machine{},
		InfraCluster: clusterScope,
		AWSMachine: &infrav1.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "infra-cluster",
			},
		},
	})
}
