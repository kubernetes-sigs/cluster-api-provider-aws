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

package ssm

import (
	"context"
	"crypto/rand"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm/mock_ssmiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type mockAPIError struct {
	Code    string
	Message string
}

// ErrorCode returns the error's code, making it satisfy one part of the interface.
func (e *mockAPIError) ErrorCode() string {
	return e.Code
}

// Error returns the error's message, satisfying the standard 'error' interface.
func (e *mockAPIError) Error() string {
	return e.Message
}

// Error returns the error's message, satisfying the standard 'error' interface.
func (e *mockAPIError) ErrorMessage() string {
	return e.Message
}

// ErrorFault is the missing method needed to fully implement smithy.APIError.
func (e *mockAPIError) ErrorFault() smithy.ErrorFault {
	// smithy.FaultClient is a good default for most simulated errors.
	return smithy.FaultClient
}

var _ smithy.APIError = (*mockAPIError)(nil)

func TestServiceCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	generateBytes := func(count int64) []byte {
		token := make([]byte, count)
		_, err := rand.Read(token)
		if err != nil {
			t.Fatalf("error while creating data: %v", err)
		}
		return token
	}

	check := func(actualPrefix string, expectedPrefix string, err error, IsErrorExpected bool) {
		if !strings.HasPrefix(actualPrefix, expectedPrefix) {
			t.Fatalf("Prefix is not as expected: %v", actualPrefix)
		}
		if (err != nil) != isErrorRetryable(err, retryableErrors) {
			t.Fatalf("Unexpected error value, error = %v, expectedError %v", err, IsErrorExpected)
		}
	}

	sortTagsByKey := func(tags []types.Tag) {
		sort.Slice(tags, func(i, j int) bool {
			return *(tags[i].Key) < *(tags[j].Key)
		})
	}

	expectedTags := []*types.Tag{
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

	expectedTagsAsValues := make([]types.Tag, 0, len(expectedTags))
	for _, tagPtr := range expectedTags {
		if tagPtr != nil {
			expectedTagsAsValues = append(expectedTagsAsValues, *tagPtr)
		}
	}

	tests := []struct {
		name           string
		bytesCount     int64
		secretPrefix   string
		expectedPrefix string
		wantErr        bool
		expect         func(m *mock_ssmiface.MockSSMAPIMockRecorder)
	}{
		{
			name:           "Should not store data in SSM if data is having zero bytes",
			bytesCount:     0,
			secretPrefix:   "/awsprefix",
			expectedPrefix: "/prefix",
		},
		{
			name:           "Should store data in SSM if data is having non-zero bytes",
			bytesCount:     10000,
			secretPrefix:   "prefix",
			expectedPrefix: "/prefix",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.PutParameter(context.TODO(), gomock.AssignableToTypeOf(&ssm.PutParameterInput{})).MinTimes(1).Return(&ssm.PutParameterOutput{}, nil).Do(
					func(ctx context.Context, putParameterInput *ssm.PutParameterInput, optFns ...func(*ssm.Options)) {
						if !strings.HasPrefix(*(putParameterInput.Name), "/prefix/") {
							t.Fatalf("Prefix is not as expected: %v", putParameterInput.Name)
						}

						sortTagsByKey(putParameterInput.Tags)
						sortTagsByKey(expectedTagsAsValues)

						if diff := cmp.Diff(expectedTagsAsValues, putParameterInput.Tags,
							// This Comparer is the key. It tells cmp how to properly compare two Tag structs.
							cmp.Comparer(func(x, y types.Tag) bool {
								// Safely compare the string content of the Key field.
								keyMatch := (x.Key == nil && y.Key == nil) || (x.Key != nil && y.Key != nil && *x.Key == *y.Key)
								// Safely compare the string content of the Value field.
								valueMatch := (x.Value == nil && y.Value == nil) || (x.Value != nil && y.Value != nil && *x.Value == *y.Value)

								// The tags are equal only if both their key and value content match.
								return keyMatch && valueMatch
							}),
						); diff != "" {
							t.Fatalf("Tags mismatch (-expected +actual):\n%s", diff)
						}
					},
				)
			},
		},
		{
			name:           "Should not retry if non-retryable error occurred while storing data in SSM",
			bytesCount:     10,
			secretPrefix:   "/prefix",
			expectedPrefix: "/prefix",
			wantErr:        true,
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.PutParameter(context.TODO(), gomock.AssignableToTypeOf(&ssm.PutParameterInput{})).Return(nil, &mockAPIError{
					"ParameterAlreadyExists",
					"parameter already exists"}).Do(
					func(ctx context.Context, putParameterInput *ssm.PutParameterInput, optFns ...func(*ssm.Options)) {
						if !strings.HasPrefix(*(putParameterInput.Name), "/prefix/") {
							t.Fatalf("Prefix is not as expected: %v", putParameterInput.Name)
						}
						sortTagsByKey(putParameterInput.Tags)
						sortTagsByKey(expectedTagsAsValues)
						if diff := cmp.Diff(expectedTagsAsValues, putParameterInput.Tags,
							// This Comparer is the key. It tells cmp how to properly compare two Tag structs.
							cmp.Comparer(func(x, y types.Tag) bool {
								// Safely compare the string content of the Key field.
								keyMatch := (x.Key == nil && y.Key == nil) || (x.Key != nil && y.Key != nil && *x.Key == *y.Key)
								// Safely compare the string content of the Value field.
								valueMatch := (x.Value == nil && y.Value == nil) || (x.Value != nil && y.Value != nil && *x.Value == *y.Value)

								// The tags are equal only if both their key and value content match.
								return keyMatch && valueMatch
							}),
						); diff != "" {
							t.Fatalf("Tags mismatch (-expected +actual):\n%s", diff)
						}
					},
				)
			},
		},
		{
			name:           "Should retry if retryable error occurred while storing data in SSM",
			bytesCount:     10,
			secretPrefix:   "",
			expectedPrefix: "/cluster.x-k8s.io",
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.PutParameter(context.TODO(), gomock.AssignableToTypeOf(&ssm.PutParameterInput{})).Return(nil, &mockAPIError{
					"ParameterLimitExceeded",
					"parameter limit exceeded"})
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
			ssmClientMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(ssmClientMock.EXPECT())
			}
			s := NewService(clusterScope)
			s.SSMClient = ssmClientMock

			ms, err := getMachineScope(client, clusterScope)
			g.Expect(err).NotTo(HaveOccurred())
			ms.SetSecretPrefix(tt.secretPrefix)
			data := generateBytes(tt.bytesCount)

			prefix, _, err := s.Create(ms, data)
			check(prefix, tt.expectedPrefix, err, tt.wantErr)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		secretCount int32
		expect      func(m *mock_ssmiface.MockSSMAPIMockRecorder)
		wantErr     bool
		check       func(error)
	}{
		{
			name:        "Should not call AWS when secret count has zero value",
			secretCount: 0,
			expect:      func(m *mock_ssmiface.MockSSMAPIMockRecorder) {},
		},
		{
			name:        "Should not return error when delete is successful",
			secretCount: 1,
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.DeleteParameter(context.TODO(), gomock.Eq(&ssm.DeleteParameterInput{
					Name: aws.String("prefix/0"),
				})).Return(&ssm.DeleteParameterOutput{}, nil)
			},
		},
		{
			name:        "Should return all errors except not found errors",
			secretCount: 3,
			expect: func(m *mock_ssmiface.MockSSMAPIMockRecorder) {
				m.DeleteParameter(context.TODO(), gomock.Eq(&ssm.DeleteParameterInput{
					Name: aws.String("prefix/0"),
				})).Return(nil, &smithy.GenericAPIError{Code: "FailedDependency", Message: "failed dependency"})
				m.DeleteParameter(context.TODO(), gomock.Eq(&ssm.DeleteParameterInput{
					Name: aws.String("prefix/1"),
				})).Return(nil, &mockAPIError{
					Code:    "ParameterNotFound",
					Message: "not found",
				})
				m.DeleteParameter(context.TODO(), gomock.Eq(&ssm.DeleteParameterInput{
					Name: aws.String("prefix/2"),
				})).Return(nil, &smithy.GenericAPIError{Code: "Conflict", Message: "new conflict"})
			},
			wantErr: true,
			check: func(err error) {
				if err.Error() != "[api error FailedDependency: failed dependency, api error Conflict: new conflict]" {
					t.Fatalf("Unexpected error: %v", err)
				}
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

			ssmClientMock := mock_ssmiface.NewMockSSMAPI(mockCtrl)
			tt.expect(ssmClientMock.EXPECT())
			s := NewService(clusterScope)
			s.SSMClient = ssmClientMock
			ms, err := getMachineScope(client, clusterScope)
			g.Expect(err).NotTo(HaveOccurred())

			ms.SetSecretPrefix("prefix")
			ms.SetSecretCount(tt.secretCount)

			err = s.Delete(ms)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
				if tt.check != nil {
					tt.check(err)
				}
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
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
