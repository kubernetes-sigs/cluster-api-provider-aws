/*
Copyright 2021 The Kubernetes Authors.

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

package s3_test

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	s3svc "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3/mock_s3iface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3/mock_stsiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	testClusterName      = "test-cluster"
	testClusterNamespace = "test-namespace"
)

func TestReconcileBucket(t *testing.T) {
	t.Parallel()

	t.Run("does_nothing_when_bucket_management_is_disabled", func(t *testing.T) {
		t.Parallel()

		svc, _ := testService(t, nil)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("creates_bucket_with_configured_name", func(t *testing.T) {
		t.Parallel()

		expectedBucketName := "baz"

		svc, s3Mock := testService(t, &infrav1.S3Bucket{
			Name: expectedBucketName,
		})

		input := &s3svc.CreateBucketInput{
			Bucket: aws.String(expectedBucketName),
		}

		s3Mock.EXPECT().CreateBucket(gomock.Eq(input)).Return(nil, nil).Times(1)
		s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Return(nil, nil).Times(1)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("hashes_default_bucket_name_if_name_exceeds_maximum_length", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		s3Mock := mock_s3iface.NewMockS3API(mockCtrl)
		stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)

		getCallerIdentityResult := &sts.GetCallerIdentityOutput{Account: aws.String("foo")}
		stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Return(getCallerIdentityResult, nil).AnyTimes()

		scheme := runtime.NewScheme()
		_ = infrav1.AddToScheme(scheme)
		client := fake.NewClientBuilder().WithScheme(scheme).Build()

		longName := strings.Repeat("a", 40)
		scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
			Client: client,
			Cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      longName,
					Namespace: longName,
				},
			},
			AWSCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{},
				},
			},
		})
		if err != nil {
			t.Fatalf("Failed to create test context: %v", err)
		}

		svc := s3.NewService(scope)
		svc.S3Client = s3Mock
		svc.STSClient = stsMock

		s3Mock.EXPECT().CreateBucket(gomock.Any()).Do(func(input *s3svc.CreateBucketInput) {
			if input.Bucket == nil {
				t.Fatalf("CreateBucket request must have Bucket specified")
			}

			if strings.Contains(*input.Bucket, longName) {
				t.Fatalf("Default bucket name be hashed when it's very long, got: %q", *input.Bucket)
			}
		}).Return(nil, nil).Times(1)

		s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Return(nil, nil).Times(1)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("creates_bucket_with_policy_allowing_controlplane_and_worker_nodes_to_read_their_secrets", func(t *testing.T) {
		t.Parallel()

		bucketName := "bar"

		svc, s3Mock := testService(t, &infrav1.S3Bucket{
			Name:                           bucketName,
			ControlPlaneIAMInstanceProfile: fmt.Sprintf("control-plane%s", iamv1.DefaultNameSuffix),
			NodesIAMInstanceProfiles: []string{
				fmt.Sprintf("nodes%s", iamv1.DefaultNameSuffix),
			},
		})

		s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, nil).Times(1)
		s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Do(func(input *s3svc.PutBucketPolicyInput) {
			if input.Policy == nil {
				t.Fatalf("Policy must be defined")
			}

			policy := *input.Policy

			if !strings.Contains(policy, "role/control-plane") {
				t.Errorf("At least one policy should include a reference to control-plane role, got: %v", policy)
			}

			if !strings.Contains(policy, "role/node") {
				t.Errorf("At least one policy should include a reference to node role, got: %v", policy)
			}

			if !strings.Contains(policy, fmt.Sprintf("%s/control-plane/*", bucketName)) {
				t.Errorf("At least one policy should apply for all objects with %q prefix, got: %v", "control-plane", policy)
			}

			if !strings.Contains(policy, fmt.Sprintf("%s/node/*", bucketName)) {
				t.Errorf("At least one policy should apply for all objects with %q prefix, got: %v", "node", policy)
			}

			if !strings.Contains(policy, "arn:aws:iam::foo:role/control-plane.cluster-api-provider-aws.sigs.k8s.io") {
				t.Errorf("Expected arn to contain the right principal; got: %v", policy)
			}
		}).Return(nil, nil).Times(1)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("is_idempotent", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, nil).Times(2)
		s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Return(nil, nil).Times(2)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("ignores_when_bucket_already_exists_but_its_owned_by_the_same_account", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		err := awserr.New(s3svc.ErrCodeBucketAlreadyOwnedByYou, "err", errors.New("err"))

		s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, err).Times(1)
		s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Return(nil, nil).Times(1)

		if err := svc.ReconcileBucket(); err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}
	})

	t.Run("returns_error_when", func(t *testing.T) {
		t.Parallel()

		t.Run("bucket_creation_fails", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, errors.New("error")).Times(1)

			if err := svc.ReconcileBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("bucket_creation_returns_unexpected_AWS_error", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, awserr.New("foo", "", nil)).Times(1)

			if err := svc.ReconcileBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("generating_bucket_policy_fails", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, nil).Times(1)

			mockCtrl := gomock.NewController(t)
			stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
			stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Return(nil, fmt.Errorf(t.Name())).AnyTimes()
			svc.STSClient = stsMock

			if err := svc.ReconcileBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("creating_bucket_policy_fails", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().CreateBucket(gomock.Any()).Return(nil, nil).Times(1)
			s3Mock.EXPECT().PutBucketPolicy(gomock.Any()).Return(nil, errors.New("error")).Times(1)

			if err := svc.ReconcileBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})
	})
}

func TestDeleteBucket(t *testing.T) {
	t.Parallel()

	const bucketName = "foo"

	t.Run("does_nothing_when_bucket_management_is_disabled", func(t *testing.T) {
		t.Parallel()

		svc, _ := testService(t, nil)

		if err := svc.DeleteBucket(); err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}
	})

	t.Run("deletes_bucket_with_configured_name", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{
			Name: bucketName,
		})

		input := &s3svc.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		}

		s3Mock.EXPECT().DeleteBucket(input).Return(nil, nil).Times(1)

		if err := svc.DeleteBucket(); err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}
	})

	t.Run("returns_error_when_bucket_removal_returns", func(t *testing.T) {
		t.Parallel()
		t.Run("unexpected_error", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().DeleteBucket(gomock.Any()).Return(nil, errors.New("err")).Times(1)

			if err := svc.DeleteBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("unexpected_AWS_error", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			s3Mock.EXPECT().DeleteBucket(gomock.Any()).Return(nil, awserr.New("foo", "", nil)).Times(1)

			if err := svc.DeleteBucket(); err == nil {
				t.Fatalf("Expected error")
			}
		})
	})

	t.Run("ignores_when_bucket_has_already_been_removed", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		s3Mock.EXPECT().DeleteBucket(gomock.Any()).Return(nil, awserr.New(s3svc.ErrCodeNoSuchBucket, "", nil)).Times(1)

		if err := svc.DeleteBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("skips_bucket_removal_when_bucket_is_not_empty", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		s3Mock.EXPECT().DeleteBucket(gomock.Any()).Return(nil, awserr.New("BucketNotEmpty", "", nil)).Times(1)

		if err := svc.DeleteBucket(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}

func TestCreateObject(t *testing.T) {
	t.Parallel()

	const (
		bucketName = "foo"
		nodeName   = "aws-test1"
	)

	t.Run("for_machine", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{
			Name: bucketName,
		})

		machineScope := &scope.MachineScope{
			Machine: &clusterv1.Machine{},
			AWSMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			},
		}

		bootstrapData := []byte("foobar")

		s3Mock.EXPECT().PutObject(gomock.Any()).Do(func(putObjectInput *s3svc.PutObjectInput) {
			t.Run("use_configured_bucket_name_on_cluster_level", func(t *testing.T) {
				t.Parallel()

				if *putObjectInput.Bucket != bucketName {
					t.Fatalf("Expected object to be created in bucket %q, got %q", bucketName, *putObjectInput.Bucket)
				}
			})

			t.Run("use_machine_role_and_machine_name_as_key", func(t *testing.T) {
				t.Parallel()

				if !strings.HasPrefix(*putObjectInput.Key, "node") {
					t.Errorf("Expected key to start with node role, got: %q", *putObjectInput.Key)
				}

				if !strings.HasSuffix(*putObjectInput.Key, nodeName) {
					t.Errorf("Expected key to end with node name, got: %q", *putObjectInput.Key)
				}
			})

			t.Run("puts_given_bootstrap_data_untouched", func(t *testing.T) {
				t.Parallel()

				data, err := io.ReadAll(putObjectInput.Body)
				if err != nil {
					t.Fatalf("Reading put object body: %v", err)
				}

				if !reflect.DeepEqual(data, bootstrapData) {
					t.Fatalf("Unexpected request body %q, expected %q", string(data), string(bootstrapData))
				}
			})
		}).Return(nil, nil).Times(1)

		bootstrapDataURL, err := svc.Create(machineScope, bootstrapData)
		if err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}

		t.Run("returns_s3_url_for_created_object", func(t *testing.T) {
			t.Parallel()

			parsedURL, err := url.Parse(bootstrapDataURL)
			if err != nil {
				t.Fatalf("Parsing URL %q: %v", bootstrapDataURL, err)
			}

			expectedScheme := "s3"
			if parsedURL.Scheme != expectedScheme {
				t.Errorf("Unexpected URL scheme, expected %q, got %q", expectedScheme, parsedURL.Scheme)
			}

			if !strings.Contains(parsedURL.Host, bucketName) {
				t.Errorf("URL Host should include bucket %q reference, got %q", bucketName, parsedURL.Host)
			}

			if !strings.HasSuffix(parsedURL.Path, nodeName) {
				t.Errorf("URL Path should end with node name %q, got: %q", nodeName, parsedURL.Path)
			}
		})
	})

	t.Run("is_idempotent", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		machineScope := &scope.MachineScope{
			Machine: &clusterv1.Machine{},
			AWSMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			},
		}

		s3Mock.EXPECT().PutObject(gomock.Any()).Return(nil, nil).Times(2)

		boostrapData := []byte("foo")

		if _, err := svc.Create(machineScope, boostrapData); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if _, err := svc.Create(machineScope, boostrapData); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("returns_error_when", func(t *testing.T) {
		t.Parallel()

		t.Run("object_creation_fails", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			machineScope := &scope.MachineScope{
				Machine: &clusterv1.Machine{},
				AWSMachine: &infrav1.AWSMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: nodeName,
					},
				},
			}

			s3Mock.EXPECT().PutObject(gomock.Any()).Return(nil, errors.New("foo")).Times(1)

			bootstrapDataURL, err := svc.Create(machineScope, []byte("foo"))
			if err == nil {
				t.Fatalf("Expected error")
			}

			if bootstrapDataURL != "" {
				t.Fatalf("Expected empty bootstrap data URL when creation error occurs")
			}
		})

		t.Run("given_empty_machine_scope", func(t *testing.T) {
			t.Parallel()

			svc, _ := testService(t, &infrav1.S3Bucket{})

			bootstrapDataURL, err := svc.Create(nil, []byte("foo"))
			if err == nil {
				t.Fatalf("Expected error")
			}

			if bootstrapDataURL != "" {
				t.Fatalf("Expected empty bootstrap data URL when creation error occurs")
			}
		})

		// If one tries to put empty bootstrap data into S3, most likely something is wrong.
		t.Run("given_empty_bootstrap_data", func(t *testing.T) {
			t.Parallel()

			svc, _ := testService(t, &infrav1.S3Bucket{})

			machineScope := &scope.MachineScope{
				Machine: &clusterv1.Machine{},
				AWSMachine: &infrav1.AWSMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: nodeName,
					},
				},
			}

			bootstrapDataURL, err := svc.Create(machineScope, []byte{})
			if err == nil {
				t.Fatalf("Expected error")
			}

			if bootstrapDataURL != "" {
				t.Fatalf("Expected empty bootstrap data URL when creation error occurs")
			}
		})

		t.Run("bucket_management_is_disabled_clusterwide", func(t *testing.T) {
			t.Parallel()

			svc, _ := testService(t, nil)

			machineScope := &scope.MachineScope{
				Machine: &clusterv1.Machine{},
				AWSMachine: &infrav1.AWSMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: nodeName,
					},
				},
			}

			bootstrapDataURL, err := svc.Create(machineScope, []byte("foo"))
			if err == nil {
				t.Fatalf("Expected error")
			}

			if bootstrapDataURL != "" {
				t.Fatalf("Expected empty bootstrap data URL when creation error occurs")
			}
		})
	})
}

func TestDeleteObject(t *testing.T) {
	t.Parallel()

	const nodeName = "aws-test1"

	t.Run("for_machine", func(t *testing.T) {
		t.Parallel()

		expectedBucketName := "foo"

		svc, s3Mock := testService(t, &infrav1.S3Bucket{
			Name: expectedBucketName,
		})

		machineScope := &scope.MachineScope{
			Machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.MachineControlPlaneLabel: "",
					},
				},
			},
			AWSMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			},
		}

		s3Mock.EXPECT().DeleteObject(gomock.Any()).Do(func(deleteObjectInput *s3svc.DeleteObjectInput) {
			t.Run("use_configured_bucket_name_on_cluster_level", func(t *testing.T) {
				t.Parallel()

				if *deleteObjectInput.Bucket != expectedBucketName {
					t.Fatalf("Expected object to be deleted from bucket %q, got %q", expectedBucketName, *deleteObjectInput.Bucket)
				}
			})

			t.Run("use_machine_role_and_machine_name_as_key", func(t *testing.T) {
				t.Parallel()

				if !strings.HasPrefix(*deleteObjectInput.Key, "control-plane") {
					t.Errorf("Expected key to start with control-plane role, got: %q", *deleteObjectInput.Key)
				}

				if !strings.HasSuffix(*deleteObjectInput.Key, nodeName) {
					t.Errorf("Expected key to end with node name, got: %q", *deleteObjectInput.Key)
				}
			})
		}).Return(nil, nil).Times(1)

		if err := svc.Delete(machineScope); err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}
	})

	t.Run("succeeds_when_bucket_has_already_been_removed", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		machineScope := &scope.MachineScope{
			Machine: &clusterv1.Machine{},
			AWSMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			},
		}

		s3Mock.EXPECT().DeleteObject(gomock.Any()).Return(nil, awserr.New(s3svc.ErrCodeNoSuchBucket, "", nil)).Times(1)

		if err := svc.Delete(machineScope); err != nil {
			t.Fatalf("Unexpected error, got: %v", err)
		}
	})

	t.Run("returns_error_when", func(t *testing.T) {
		t.Parallel()

		t.Run("object_deletion_fails", func(t *testing.T) {
			t.Parallel()

			svc, s3Mock := testService(t, &infrav1.S3Bucket{})

			machineScope := &scope.MachineScope{
				Machine: &clusterv1.Machine{},
				AWSMachine: &infrav1.AWSMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: nodeName,
					},
				},
			}

			s3Mock.EXPECT().DeleteObject(gomock.Any()).Return(nil, errors.New("foo")).Times(1)

			if err := svc.Delete(machineScope); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("given_empty_machine_scope", func(t *testing.T) {
			t.Parallel()

			svc, _ := testService(t, &infrav1.S3Bucket{})

			if err := svc.Delete(nil); err == nil {
				t.Fatalf("Expected error")
			}
		})

		t.Run("bucket_management_is_disabled_clusterwide", func(t *testing.T) {
			t.Parallel()

			svc, _ := testService(t, nil)

			machineScope := &scope.MachineScope{
				Machine: &clusterv1.Machine{},
				AWSMachine: &infrav1.AWSMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: nodeName,
					},
				},
			}

			if err := svc.Delete(machineScope); err == nil {
				t.Fatalf("Expected error")
			}
		})
	})

	t.Run("is_idempotent", func(t *testing.T) {
		t.Parallel()

		svc, s3Mock := testService(t, &infrav1.S3Bucket{})

		machineScope := &scope.MachineScope{
			Machine: &clusterv1.Machine{},
			AWSMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			},
		}

		s3Mock.EXPECT().DeleteObject(gomock.Any()).Return(nil, nil).Times(2)

		if err := svc.Delete(machineScope); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if err := svc.Delete(machineScope); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}

func testService(t *testing.T, bucket *infrav1.S3Bucket) (*s3.Service, *mock_s3iface.MockS3API) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	s3Mock := mock_s3iface.NewMockS3API(mockCtrl)
	stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)

	getCallerIdentityResult := &sts.GetCallerIdentityOutput{Account: aws.String("foo")}
	stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Return(getCallerIdentityResult, nil).AnyTimes()

	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testClusterName,
				Namespace: testClusterNamespace,
			},
		},
		AWSCluster: &infrav1.AWSCluster{
			Spec: infrav1.AWSClusterSpec{
				S3Bucket: bucket,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}

	svc := s3.NewService(scope)
	svc.S3Client = s3Mock
	svc.STSClient = stsMock

	return svc, s3Mock
}
