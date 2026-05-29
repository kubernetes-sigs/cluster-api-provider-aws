/*
Copyright 2023 The Kubernetes Authors.

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

// Package controllers provides a way to reconcile ROSA resources.
package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	restclient "k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// fakeStsAPIClient is a minimal implementation of rosaaws.StsApiClient for use in tests.
// It returns a fixed caller identity without needing a mock framework.
type fakeStsAPIClient struct {
	account string
	arn     string
	userID  string
}

func (f *fakeStsAPIClient) GetCallerIdentity(_ context.Context, _ *stsv2.GetCallerIdentityInput, _ ...func(*stsv2.Options)) (*stsv2.GetCallerIdentityOutput, error) {
	return &stsv2.GetCallerIdentityOutput{
		Account: aws.String(f.account),
		Arn:     aws.String(f.arn),
		UserId:  aws.String(f.userID),
	}, nil
}

func (f *fakeStsAPIClient) AssumeRole(_ context.Context, _ *stsv2.AssumeRoleInput, _ ...func(*stsv2.Options)) (*stsv2.AssumeRoleOutput, error) {
	return nil, nil
}

func (f *fakeStsAPIClient) AssumeRoleWithWebIdentity(_ context.Context, _ *stsv2.AssumeRoleWithWebIdentityInput, _ ...func(*stsv2.Options)) (*stsv2.AssumeRoleWithWebIdentityOutput, error) {
	return nil, nil
}

// newFakeAWSClientFactory returns an awsClientFactory that injects a rosaaws.Client backed by
// the given fakeStsAPIClient. GetCreator() on the returned client will use that STS stub.
func newFakeAWSClientFactory(fakeSts *fakeStsAPIClient) func(*scope.ROSAControlPlaneScope) (rosaaws.Client, error) {
	return func(_ *scope.ROSAControlPlaneScope) (rosaaws.Client, error) {
		return rosaaws.New(
			aws.Config{},
			rosaaws.NewLoggerWrapper(logrus.New(), nil),
			nil, nil, nil, nil, nil,
			fakeSts,
			nil, nil, nil,
			&rosaaws.AccessKey{},
			false,
		), nil
	}
}

func TestUpdateOCMClusterSpec(t *testing.T) {
	g := NewWithT(t)
	// Test case 1: No updates, everything matches
	t.Run("No Updates When Specs Are Same", func(t *testing.T) {
		// Mock ROSAControlPlane input
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AuditLogRoleARN: "arn:aws:iam::123456789012:role/AuditLogRole",
				ClusterRegistryConfig: &rosacontrolplanev1.RegistryConfig{
					AdditionalTrustedCAs: map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"},
					AllowedRegistriesForImport: []rosacontrolplanev1.RegistryLocation{
						{DomainName: "registry1.com", Insecure: false},
					},
				},
			},
		}

		// Mock Cluster input
		mockCluster, _ := v1.NewCluster().
			AWS(v1.NewAWS().
				AuditLog(v1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/AuditLogRole"))).
			RegistryConfig(v1.NewClusterRegistryConfig().
				AdditionalTrustedCa(map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
				AllowedRegistriesForImport(v1.NewRegistryLocation().
					DomainName("registry1.com").
					Insecure(false))).Build()

		expectedOCMSpec := ocm.Spec{}

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeFalse())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 2: Update when AuditLogRoleARN is different
	t.Run("Update AuditLogRoleARN", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AuditLogRoleARN: "arn:aws:iam::123456789012:role/NewAuditLogRole",
			},
		}

		mockCluster, _ := v1.NewCluster().
			AWS(v1.NewAWS().
				AuditLog(v1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/OldAuditLogRole"))).Build()

		expectedOCMSpec := ocm.Spec{
			AuditLogRoleARN: &rosaControlPlane.Spec.AuditLogRoleARN,
		}

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 3: Update when RegistryConfig is different
	t.Run("Update RegistryConfig", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ClusterRegistryConfig: &rosacontrolplanev1.RegistryConfig{
					AdditionalTrustedCAs: map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"},
					AllowedRegistriesForImport: []rosacontrolplanev1.RegistryLocation{
						{DomainName: "new-registry.com", Insecure: true},
					},
					RegistrySources: &rosacontrolplanev1.RegistrySources{
						AllowedRegistries: []string{"quay.io", "reg1.org"},
					},
				},
			},
		}

		mockCluster, _ := v1.NewCluster().
			RegistryConfig(v1.NewClusterRegistryConfig().
				AdditionalTrustedCa(map[string]string{"old-trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
				AllowedRegistriesForImport(v1.NewRegistryLocation().
					DomainName("old-registry.com").
					Insecure(false)).RegistrySources(v1.NewRegistrySources().BlockedRegistries([]string{"blocked.io", "blocked.org"}...))).
			Build()

		expectedOCMSpec := ocm.Spec{
			AdditionalTrustedCa:        rosaControlPlane.Spec.ClusterRegistryConfig.AdditionalTrustedCAs,
			AllowedRegistriesForImport: "new-registry.com:true",
			AllowedRegistries:          rosaControlPlane.Spec.ClusterRegistryConfig.RegistrySources.AllowedRegistries,
		}
		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 4: AllowedRegistriesForImport mismatch
	t.Run("Update AllowedRegistriesForImport", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ClusterRegistryConfig: &rosacontrolplanev1.RegistryConfig{
					AllowedRegistriesForImport: []rosacontrolplanev1.RegistryLocation{},
				},
			},
		}

		mockCluster, _ := v1.NewCluster().
			RegistryConfig(v1.NewClusterRegistryConfig().
				AllowedRegistriesForImport(v1.NewRegistryLocation().
					DomainName("old-registry.com").
					Insecure(false))).
			Build()

		expectedOCMSpec := ocm.Spec{
			AllowedRegistriesForImport: "",
		}

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 6: channel group update
	t.Run("Update channel group", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ChannelGroup: rosacontrolplanev1.Candidate,
			},
		}

		mockCluster, _ := v1.NewCluster().
			Version(v1.NewVersion().
				ChannelGroup("stable")).
			Build()

		expectedOCMSpec := ocm.Spec{
			ChannelGroup: "candidate",
		}

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 7: AutoNode update
	t.Run("Update Auto Node", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AutoNode: &rosacontrolplanev1.AutoNode{
					Mode:    rosacontrolplanev1.AutoNodeModeEnabled,
					RoleARN: "autoNodeARN",
				},
			},
		}

		mockCluster, _ := v1.NewCluster().
			AutoNode(v1.NewClusterAutoNode().Mode("disabled")).
			AWS(v1.NewAWS().AutoNode(v1.NewAwsAutoNode().RoleArn("anyARN"))).
			Build()

		expectedOCMSpec := ocm.Spec{
			AutoNodeMode:    "enabled",
			AutoNodeRoleARN: "autoNodeARN",
		}

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec).To(Equal(expectedOCMSpec))
	})

	// Test case 8: Channel explicitly set - use it directly
	t.Run("Channel explicitly set", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel:      "eus-4.18",
				ChannelGroup: rosacontrolplanev1.Stable, // Different from channel, but channel takes precedence
				Version:      "4.16.5",
			},
		}

		mockCluster, _ := v1.NewCluster().
			Version(v1.NewVersion().
				ID("4.16.5").
				ChannelGroup("stable")).
			Build()

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec.Channel).To(Equal("eus-4.18"))
		g.Expect(ocmSpec.ChannelGroup).To(BeEmpty())
	})

	// Test case 9: ChannelGroup matches current - no update
	t.Run("ChannelGroup matches current channel group - no update", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ChannelGroup: rosacontrolplanev1.Stable,
				Version:      "4.16.5",
			},
		}

		mockCluster, _ := v1.NewCluster().
			Version(v1.NewVersion().
				ID("4.18.3").
				ChannelGroup("stable")).
			Build()

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeFalse())
		g.Expect(ocmSpec.Channel).To(BeEmpty())
		g.Expect(ocmSpec.ChannelGroup).To(BeEmpty())
	})

	// Test case 10: ChannelGroup changes - update channelGroup
	t.Run("ChannelGroup changes from stable to eus", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ChannelGroup: rosacontrolplanev1.Eus, // Changing from stable to eus
				Version:      "4.16.5",
			},
		}

		mockCluster, _ := v1.NewCluster().
			Version(v1.NewVersion().
				ID("4.18.3").
				ChannelGroup("stable")).
			Build()

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeTrue())
		g.Expect(ocmSpec.ChannelGroup).To(Equal("eus"))
		g.Expect(ocmSpec.Channel).To(BeEmpty())
	})

	// Test case 11: ChannelGroup set, no current version info
	t.Run("ChannelGroup set with no current version info", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ChannelGroup: rosacontrolplanev1.Eus,
				Version:      "4.16.5",
			},
		}

		// Cluster has no version info (edge case)
		// When version is nil, we skip the channelGroup update
		mockCluster, _ := v1.NewCluster().Build()

		reconciler := &ROSAControlPlaneReconciler{}
		ocmSpec, updated := reconciler.updateOCMClusterSpec(rosaControlPlane, mockCluster)

		g.Expect(updated).To(BeFalse())
		g.Expect(ocmSpec.ChannelGroup).To(BeEmpty())
		g.Expect(ocmSpec.Channel).To(BeEmpty())
	})
}

func TestValidateControlPlaneSpec(t *testing.T) {
	g := NewWithT(t)

	mockCtrl := gomock.NewController(t)
	ocmMock := mocks.NewMockOCMClient(mockCtrl)
	expect := func(m *mocks.MockOCMClientMockRecorder) {
		m.ValidateHypershiftVersion(gomock.Any(), gomock.Any()).DoAndReturn(func(versionRawID string, channelGroup string) (bool, error) {
			return true, nil
		}).AnyTimes()
	}
	expect(ocmMock.EXPECT())

	// Test case 1: AutoNode and Version are set valid
	t.Run("AutoNode is valid.", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AutoNode: &rosacontrolplanev1.AutoNode{
					Mode:    rosacontrolplanev1.AutoNodeModeEnabled,
					RoleARN: "autoNodeARN",
				},
				Version:      "4.19.0",
				ChannelGroup: rosacontrolplanev1.Stable,
			},
		}
		str, err := validateControlPlaneSpec(ocmMock, rosaControlPlane)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(str).To(Equal(""))
	})

	// Test case 2: AutoNode is enabled and AutoNode ARN is empty.
	t.Run("AutoNode is enabled and AutoNode ARN is empty", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AutoNode: &rosacontrolplanev1.AutoNode{
					Mode:    rosacontrolplanev1.AutoNodeModeEnabled,
					RoleARN: "",
				},
				Version:      "4.19.0",
				ChannelGroup: rosacontrolplanev1.Stable,
			},
		}
		str, err := validateControlPlaneSpec(ocmMock, rosaControlPlane)
		g.Expect(err).To(HaveOccurred())
		g.Expect(strings.Contains(err.Error(), "autoNode.roleARN, must be set when autoNode mode is enabled")).To(BeTrue())
		g.Expect(str).To(Equal(""))
	})

	// Test case 3: AutoNode is disabled and AutoNode ARN is empty.
	t.Run("AutoNode is disabled and AutoNode ARN is empty.", func(t *testing.T) {
		rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				AutoNode: &rosacontrolplanev1.AutoNode{
					Mode:    rosacontrolplanev1.AutoNodeModeDisabled,
					RoleARN: "",
				},
				Version:      "4.19.0",
				ChannelGroup: rosacontrolplanev1.Stable,
			},
		}
		str, err := validateControlPlaneSpec(ocmMock, rosaControlPlane)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(str).To(Equal(""))
	})
}

func TestRosaControlPlaneReconcileStatusVersion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, "test-namespace")
	g.Expect(err).ToNot(HaveOccurred())

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rosa-secret",
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"ocmToken": []byte("secret-ocm-token-string"),
		},
	}

	identity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterControllerIdentity"))

	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rosa-control-plane-1",
			Namespace: ns.Name,
			UID:       types.UID("rosa-control-plane-1"),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ROSAControlPlane",
			APIVersion: rosacontrolplanev1.GroupVersion.String(),
		},
		Spec: rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   "rosa-control-plane-1",
			Subnets:           []string{"subnet-0ac99a6230b408813", "subnet-1ac99a6230b408811"},
			AvailabilityZones: []string{"az-1", "az-2"},
			Network: &rosacontrolplanev1.NetworkSpec{
				MachineCIDR: "10.0.0.0/16",
				PodCIDR:     "10.128.0.0/14",
				ServiceCIDR: "172.30.0.0/16",
			},
			Region:       "us-east-1",
			Version:      "4.15.20",
			ChannelGroup: "stable",
			RolesRef: rosacontrolplanev1.AWSRolesRef{
				IngressARN:              "op-arn1",
				ImageRegistryARN:        "op-arn2",
				StorageARN:              "op-arn3",
				NetworkARN:              "op-arn4",
				KubeCloudControllerARN:  "op-arn5",
				NodePoolManagementARN:   "op-arn6",
				ControlPlaneOperatorARN: "op-arn7",
				KMSProviderARN:          "op-arn8",
			},
			OIDCID:           "iodcid1",
			InstallerRoleARN: "arn:aws:iam::123456789012:role/installer",
			WorkerRoleARN:    "arn:aws:iam::123456789012:role/worker",
			SupportRoleARN:   "arn:aws:iam::123456789012:role/support",
			CredentialsSecretRef: &corev1.LocalObjectReference{
				Name: secret.Name,
			},
			VersionGate: "Acknowledge",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: identity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
		Status: rosacontrolplanev1.RosaControlPlaneStatus{
			ID: "rosa-control-plane-1",
			Conditions: clusterv1beta1.Conditions{clusterv1beta1.Condition{
				Type:               "Paused",
				Status:             "False",
				Severity:           "",
				Reason:             "NotPaused",
				Message:            "",
				LastTransitionTime: metav1.NewTime(time.Now()),
			}},
		},
	}

	ownerCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "owner-cluster-1",
			Namespace: ns.Name,
			UID:       types.UID("owner-cluster-1"),
		},
		Spec: clusterv1.ClusterSpec{
			ControlPlaneRef: clusterv1.ContractVersionedObjectReference{
				Name:     rosaControlPlane.Name,
				Kind:     "ROSAControlPlane",
				APIGroup: rosacontrolplanev1.GroupVersion.Group,
			},
		},
	}

	kubeconfigSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-kubeconfig", ownerCluster.Name),
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"kubeconfig": []byte("secret-kubeconfig-string"),
		},
	}

	rosaControlPlane.OwnerReferences = []metav1.OwnerReference{
		{
			Name:       ownerCluster.Name,
			UID:        ownerCluster.UID,
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
	}

	mockCtrl := gomock.NewController(t)
	ctx := context.TODO()
	ocmMock := mocks.NewMockOCMClient(mockCtrl)

	fakeSts := &fakeStsAPIClient{
		account: "foo",
		arn:     "arn:aws:iam::123456789012:rosa/foo",
		userID:  "user-id",
	}

	expect := func(m *mocks.MockOCMClientMockRecorder) {
		m.ValidateHypershiftVersion(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (bool, error) {
			return true, nil
		}).Times(1)
		m.GetCluster(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterKey string, creator *rosaaws.Creator) (*v1.Cluster, error) {
			sts := (&v1.STSBuilder{}).OIDCEndpointURL("oidc.com/oidc1")
			aws := v1.NewAWS().AuditLog(v1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/AuditLogRole")).STS(sts)
			console := (&v1.ClusterConsoleBuilder{}).URL("https://console.redhat.com/cluster-123")
			status := (&v1.ClusterStatusBuilder{}).State(v1.ClusterStateReady)
			api := (&v1.ClusterAPIBuilder{}).URL("https://url.com:5000")
			version := (&v1.VersionBuilder{}).RawID(rosaControlPlane.Spec.Version)
			mockCluster, _ := v1.NewCluster().AWS(aws).ID("cluster-1").Version(version).Status(status).Console(console).API(api).
				RegistryConfig(v1.NewClusterRegistryConfig().
					AdditionalTrustedCa(map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
					AllowedRegistriesForImport(v1.NewRegistryLocation().
						DomainName("registry1.com").
						Insecure(false))).
				Build()
			return mockCluster, nil
		}).Times(1)
		m.GetLogForwarders(gomock.Any()).DoAndReturn(func(clusterID string) ([]*v1.LogForwarder, error) {
			logs := []*v1.LogForwarder{}
			return logs, nil
		}).Times(1)
		m.UpdateCluster(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(clusterKey string, creator *rosaaws.Creator, config ocm.Spec) error {
			return nil
		}).Times(1)
		m.GetAvailableChannels(gomock.Any()).DoAndReturn(func(versionID string) ([]string, error) {
			return []string{"stable-4.15"}, nil
		}).Times(1)
		m.GetIdentityProviders(gomock.Any()).DoAndReturn(func(cclusterID string) ([]*v1.IdentityProvider, error) {
			ip := []*v1.IdentityProvider{}
			return ip, nil
		}).Times(1)
		m.GetUser(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(clusterID string, group string, username string) (*v1.User, error) {
			user, _ := (&v1.UserBuilder{}).ID("userid-1").Build()
			return user, nil
		}).Times(1)
		m.CreateIdentityProvider(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterID string, idp *v1.IdentityProvider) (*v1.IdentityProvider, error) {
			return idp, nil
		}).Times(1)
	}

	expect(ocmMock.EXPECT())

	// We need to mock http GET on this url
	redirectURL := "https://url.com/oauth/authorize?response_type=token&client_id=openshift-challenging-client#access_token=mocktoken&expires_in=1234"
	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate 302 redirect
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}))
	defer mockServer.Close()

	fakeAPIURL := "https://url.com"
	mockURL, _ := url.Parse(mockServer.URL)
	dialer := &net.Dialer{}
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Force redirect all traffic to mockServer
			return dialer.DialContext(ctx, "tcp", mockURL.Host)
		},
	}

	cfg := &restclient.Config{
		Host:      fakeAPIURL,
		Transport: customTransport,
	}

	objects := []client.Object{ownerCluster, rosaControlPlane, secret, identity, kubeconfigSecret}
	for _, obj := range objects {
		createObject(g, obj, ns.Name)
	}

	// Add conditions, can't do this duirng creation
	cpPh, err := patch.NewHelper(rosaControlPlane, testEnv)
	g.Expect(err).ShouldNot(HaveOccurred())
	rosaControlPlane.Status = rosacontrolplanev1.RosaControlPlaneStatus{
		ID: "rosa-control-plane-1",
		Conditions: clusterv1beta1.Conditions{clusterv1beta1.Condition{
			Type:               "Paused",
			Status:             "False",
			Severity:           "",
			Reason:             "NotPaused",
			Message:            "",
			LastTransitionTime: metav1.NewTime(time.Now()),
		}},
	}

	g.Expect(cpPh.Patch(ctx, rosaControlPlane)).To(Succeed())

	time.Sleep(50 * time.Millisecond)

	cp := &rosacontrolplanev1.ROSAControlPlane{}
	key := client.ObjectKey{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace}
	errGet := testEnv.Get(ctx, key, cp)
	g.Expect(errGet).NotTo(HaveOccurred())
	oldCondition := v1beta1conditions.Get(cp, clusterv1beta1.PausedV1Beta2Condition)
	g.Expect(oldCondition).NotTo(BeNil())

	r := ROSAControlPlaneReconciler{
		WatchFilterValue: "",
		Client:           testEnv,
		restClientConfig: cfg,
		awsClientFactory: newFakeAWSClientFactory(fakeSts),
		NewOCMClient: func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error) {
			return ocmMock, nil
		},
	}

	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace}
	_, errReconcile := r.Reconcile(ctx, req)
	g.Expect(errReconcile).ToNot(HaveOccurred())
	time.Sleep(50 * time.Millisecond)

	errGet2 := testEnv.Get(ctx, key, cp)
	g.Expect(errGet2).NotTo(HaveOccurred())
	g.Expect(cp.Status.Version).To(Equal(rosaControlPlane.Spec.Version))
	g.Expect(cp.Status.Ready).To(BeTrue())

	// cleanup
	for _, obj := range objects {
		cleanupObject(g, obj)
	}
}

func TestBuildLogForwarders(t *testing.T) {
	tests := []struct {
		name                string
		cwConfig            *rosacontrolplanev1.CloudWatchLogForwarderConfig
		s3Config            *rosacontrolplanev1.S3LogForwarderConfig
		expectCW            bool
		expectS3            bool
		expectCWGroupsCount int
		expectS3GroupsCount int
		expectCWApps        []string
		expectS3Apps        []string
	}{
		{
			name:     "both configs nil",
			expectCW: false,
			expectS3: false,
		},
		{
			name: "cloudwatch only",
			cwConfig: &rosacontrolplanev1.CloudWatchLogForwarderConfig{
				CloudWatchLogRoleArn:   "arn:aws:iam::123456789012:role/test",
				CloudWatchLogGroupName: "cw-log-group",
				Applications:           []string{"kube-apiserver", "openshift-apiserver"},
				GroupLogIDs:            []string{"audit", "infrastructure"},
			},
			expectCW:            true,
			expectCWGroupsCount: 2,
			expectCWApps:        []string{"kube-apiserver", "openshift-apiserver"},
		},
		{
			name: "s3 only",
			s3Config: &rosacontrolplanev1.S3LogForwarderConfig{
				S3ConfigBucketName:   "my-bucket",
				S3ConfigBucketPrefix: "logs/",
				Applications:         []string{"machine-config-daemon"},
				GroupLogIDs:          []string{"audit"},
			},
			expectS3:            true,
			expectS3GroupsCount: 1,
			expectS3Apps:        []string{"machine-config-daemon"},
		},
		{
			name: "cloudwatch and s3",
			cwConfig: &rosacontrolplanev1.CloudWatchLogForwarderConfig{
				CloudWatchLogRoleArn:   "arn:aws:iam::123456789012:role/test",
				CloudWatchLogGroupName: "cw-log-group",
				Applications:           []string{"kube-apiserver"},
			},
			s3Config: &rosacontrolplanev1.S3LogForwarderConfig{
				S3ConfigBucketName:   "my-bucket",
				S3ConfigBucketPrefix: "logs/",
				Applications:         []string{"machine-config-daemon"},
			},
			expectCW: true,
			expectS3: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			cwForwarder, s3Forwarder, err :=
				buildlogForwarders(tt.cwConfig, tt.s3Config)

			g.Expect(err).ToNot(HaveOccurred())

			// CloudWatch assertions
			if tt.expectCW {
				g.Expect(cwForwarder).ToNot(BeNil())
				g.Expect(cwForwarder.Cloudwatch()).ToNot(BeNil())

				if len(tt.expectCWApps) > 0 {
					g.Expect(cwForwarder.Applications()).
						To(ConsistOf(tt.expectCWApps))
				}

				if tt.expectCWGroupsCount > 0 {
					g.Expect(cwForwarder.Groups()).
						To(HaveLen(tt.expectCWGroupsCount))
				}
			} else {
				g.Expect(cwForwarder).To(BeNil())
			}

			// S3 assertions
			if tt.expectS3 {
				g.Expect(s3Forwarder).ToNot(BeNil())
				g.Expect(s3Forwarder.S3()).ToNot(BeNil())

				if len(tt.expectS3Apps) > 0 {
					g.Expect(s3Forwarder.Applications()).
						To(ConsistOf(tt.expectS3Apps))
				}

				if tt.expectS3GroupsCount > 0 {
					g.Expect(s3Forwarder.Groups()).
						To(HaveLen(tt.expectS3GroupsCount))
				}
			} else {
				g.Expect(s3Forwarder).To(BeNil())
			}
		})
	}
}

func TestReconcileLogForwarders_CreateCloudWatchAndS3(t *testing.T) {
	g := NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOCM := mocks.NewMockOCMClient(ctrl)

	clusterID := "cluster-123"
	cluster, _ := v1.NewCluster().ID(clusterID).Name("test-cluster").Build()

	rosaScope := &scope.ROSAControlPlaneScope{
		ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				CloudWatchLogForwarder: &rosacontrolplanev1.CloudWatchLogForwarderConfig{
					CloudWatchLogRoleArn:   "arn:aws:iam::123:role/test",
					CloudWatchLogGroupName: "cw-group",
					Applications:           []string{"kube-apiserver"},
				},
				S3LogForwarder: &rosacontrolplanev1.S3LogForwarderConfig{
					S3ConfigBucketName:   "bucket",
					S3ConfigBucketPrefix: "logs/",
					Applications:         []string{"machine-config-daemon"},
				},
			},
		},
		Logger: *logger.FromContext(context.TODO()),
	}

	// No existing log forwarders
	mockOCM.
		EXPECT().
		GetLogForwarders(clusterID).
		Return([]*v1.LogForwarder{}, nil)

	// Expect creation of both
	mockOCM.
		EXPECT().
		SetLogForwarder(clusterID, gomock.Any()).
		Times(2)

	reconciler := &ROSAControlPlaneReconciler{}

	err := reconciler.reconcileLogForwarders(rosaScope, mockOCM, cluster)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestReconcileLogForwarders_UpdateCW_DeleteS3(t *testing.T) {
	g := NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOCM := mocks.NewMockOCMClient(ctrl)

	clusterID := "cluster-123"
	cluster, _ := v1.NewCluster().ID(clusterID).Name("test-cluster").Build()

	existingCW, _ := v1.NewLogForwarder().
		ID("cw-id").
		Cloudwatch(v1.NewLogForwarderCloudWatchConfig()).Build()

	existingS3, _ := v1.NewLogForwarder().
		ID("s3-id").
		S3(v1.NewLogForwarderS3Config()).Build()

	rosaScope := &scope.ROSAControlPlaneScope{
		ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				CloudWatchLogForwarder: &rosacontrolplanev1.CloudWatchLogForwarderConfig{
					CloudWatchLogRoleArn:   "arn:aws:iam::123:role/test",
					CloudWatchLogGroupName: "cw-group",
				},
				S3LogForwarder: nil, // removed
			},
		},
		Logger: *logger.FromContext(context.TODO()),
	}

	mockOCM.
		EXPECT().
		GetLogForwarders(clusterID).
		Return([]*v1.LogForwarder{existingCW, existingS3}, nil)

	mockOCM.
		EXPECT().
		UpdateLogForwarder(gomock.Any(), "cw-id", clusterID).
		Return(nil)

	mockOCM.
		EXPECT().
		DeleteLogForwarder(clusterID, "s3-id").
		Return(nil)

	reconciler := &ROSAControlPlaneReconciler{}

	err := reconciler.reconcileLogForwarders(rosaScope, mockOCM, cluster)
	g.Expect(err).ToNot(HaveOccurred())
}

func createObject(g *WithT, obj client.Object, namespace string) {
	if obj.DeepCopyObject() != nil {
		obj.SetNamespace(namespace)
		g.Expect(testEnv.Create(ctx, obj)).To(Succeed())
	}
}

func cleanupObject(g *WithT, obj client.Object) {
	if obj.DeepCopyObject() != nil {
		g.Expect(testEnv.Cleanup(ctx, obj)).To(Succeed())
	}
}

// generateTestID returns a unique suffix for test object names.
func generateTestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// fakeStsHTTPResponse writes a valid STS AssumeRole XML response for use in httptest servers.
func fakeStsHTTPResponse(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials>
      <AccessKeyId>ASIAIOSFODNN7EXAMPLE</AccessKeyId>
      <SecretAccessKey>wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY</SecretAccessKey>
      <SessionToken>AQoXnyc4MCrrlandlJKwBQ==</SessionToken>
      <Expiration>2030-01-01T00:00:00Z</Expiration>
    </Credentials>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::123456789012:assumed-role/fake-rosa-role/test</Arn>
      <AssumedRoleId>ARO123EXAMPLE123:test</AssumedRoleId>
    </AssumedRoleUser>
  </AssumeRoleResult>
  <ResponseMetadata>
    <RequestId>12345678-1234-1234-1234-123456789012</RequestId>
  </ResponseMetadata>
</AssumeRoleResponse>`)
}

// TestROSAControlPlaneReconcilerWithRoleIdentity verifies that when an AWSClusterRoleIdentity
// (with AllowedNamespaces permitting all namespaces) is referenced as the identity, the
// ROSAControlPlaneReconciler successfully resolves credentials and reaches the AWS client
// factory (our canary), rather than failing at scope creation.
func TestROSAControlPlaneReconcilerWithRoleIdentity(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

	stsServer := httptest.NewServer(http.HandlerFunc(fakeStsHTTPResponse))
	defer stsServer.Close()

	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key-id")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")
	t.Setenv("AWS_REGION", "us-east-1")

	testID := generateTestID()

	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-ns-cp-ri-%s", testID[:12]))
	g.Expect(err).ToNot(HaveOccurred())

	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	controllerIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterControllerIdentity"))
	createObject(g, controllerIdentity, ns.Name)
	defer cleanupObject(g, controllerIdentity)

	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("fake-role-%s", testID[:12]),
		},
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:     fmt.Sprintf("arn:aws:iam::123456789012:role/fake-cp-role-%s", testID[:12]),
				SessionName: "test-session",
			},
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
			SourceIdentityRef: &infrav1.AWSIdentityReference{
				Name: controllerIdentity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}
	roleIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
	createObject(g, roleIdentity, ns.Name)
	defer cleanupObject(g, roleIdentity)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rosa-secret-%s", testID[:12]),
			Namespace: ns.Name,
		},
		Data: map[string][]byte{"ocmToken": []byte("fake-token")},
	}
	createObject(g, secret, ns.Name)
	defer cleanupObject(g, secret)

	cpName := fmt.Sprintf("rosa-cp-%s", testID[:12])
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cpName,
			Namespace: ns.Name,
			UID:       types.UID(cpName),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ROSAControlPlane",
			APIVersion: rosacontrolplanev1.GroupVersion.String(),
		},
		Spec: rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   cpName,
			Subnets:           []string{"subnet-0ac99a6230b408813"},
			AvailabilityZones: []string{"us-east-1a"},
			Network: &rosacontrolplanev1.NetworkSpec{
				MachineCIDR: "10.0.0.0/16",
				PodCIDR:     "10.128.0.0/14",
				ServiceCIDR: "172.30.0.0/16",
			},
			Region:           "us-east-1",
			Version:          "4.15.0",
			ChannelGroup:     "stable",
			OIDCID:           "oidcid-test",
			InstallerRoleARN: "arn:aws:iam::123456789012:role/installer",
			WorkerRoleARN:    "arn:aws:iam::123456789012:role/worker",
			SupportRoleARN:   "arn:aws:iam::123456789012:role/support",
			RolesRef: rosacontrolplanev1.AWSRolesRef{
				IngressARN:              "op-arn1",
				ImageRegistryARN:        "op-arn2",
				StorageARN:              "op-arn3",
				NetworkARN:              "op-arn4",
				KubeCloudControllerARN:  "op-arn5",
				NodePoolManagementARN:   "op-arn6",
				ControlPlaneOperatorARN: "op-arn7",
				KMSProviderARN:          "op-arn8",
			},
			CredentialsSecretRef: &corev1.LocalObjectReference{Name: secret.Name},
			VersionGate:          "Acknowledge",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: roleIdentity.Name,
				Kind: infrav1.ClusterRoleIdentityKind,
			},
		},
	}

	ownerCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("owner-cluster-%s", testID[:12]),
			Namespace: ns.Name,
			UID:       types.UID(fmt.Sprintf("owner-cluster-%s", testID[:12])),
		},
		Spec: clusterv1.ClusterSpec{
			ControlPlaneRef: clusterv1.ContractVersionedObjectReference{
				Name:     rosaControlPlane.Name,
				Kind:     "ROSAControlPlane",
				APIGroup: rosacontrolplanev1.GroupVersion.Group,
			},
		},
	}

	rosaControlPlane.OwnerReferences = []metav1.OwnerReference{
		{
			Name:       ownerCluster.Name,
			UID:        ownerCluster.UID,
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
	}

	for _, obj := range []client.Object{ownerCluster, rosaControlPlane} {
		createObject(g, obj, ns.Name)
	}
	defer cleanupObject(g, rosaControlPlane)
	defer cleanupObject(g, ownerCluster)

	awsClientCalled := false
	reconciler := &ROSAControlPlaneReconciler{
		Client: testEnv,
		// awsClientFactory is the canary: it is called only when scope creation
		// (including AWSClusterRoleIdentity credential resolution) succeeds.
		awsClientFactory: func(_ *scope.ROSAControlPlaneScope) (rosaaws.Client, error) {
			awsClientCalled = true
			return nil, fmt.Errorf("sentinel: awsClientFactory reached after role-identity scope creation")
		},
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace},
	}

	// EnsurePausedCondition on first call returns early (conditionChanged=true).
	// Retry until scope creation succeeds and awsClientFactory is reached.
	g.Eventually(func(g Gomega) {
		_, errReconcile := reconciler.Reconcile(ctx, req)
		g.Expect(errReconcile).To(HaveOccurred())
		g.Expect(errReconcile.Error()).To(ContainSubstring("failed to create AWS client"))
		g.Expect(errReconcile.Error()).NotTo(ContainSubstring("failed to create scope"))
		g.Expect(awsClientCalled).To(BeTrue())
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

// TestROSAControlPlaneReconcilerWithRoleIdentityNamespaceNotAllowed verifies that when an
// AWSClusterRoleIdentity restricts usage to a namespace that does not contain the
// ROSAControlPlane, scope creation is rejected before the AWS client factory is reached.
func TestROSAControlPlaneReconcilerWithRoleIdentityNamespaceNotAllowed(t *testing.T) {
	RegisterTestingT(t)
	g := NewWithT(t)
	ctx := context.TODO()

	stsServer := httptest.NewServer(http.HandlerFunc(fakeStsHTTPResponse))
	defer stsServer.Close()

	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key-id")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")
	t.Setenv("AWS_REGION", "us-east-1")

	testID := generateTestID()

	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("test-ns-cp-ri-denied-%s", testID[:12]))
	g.Expect(err).ToNot(HaveOccurred())

	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	controllerIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterControllerIdentity"))
	createObject(g, controllerIdentity, ns.Name)
	defer cleanupObject(g, controllerIdentity)

	// AllowedNamespaces permits only "other-namespace" — the control plane's namespace is excluded.
	roleIdentity := &infrav1.AWSClusterRoleIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("restricted-role-%s", testID[:12]),
		},
		Spec: infrav1.AWSClusterRoleIdentitySpec{
			AWSRoleSpec: infrav1.AWSRoleSpec{
				RoleArn:     fmt.Sprintf("arn:aws:iam::123456789012:role/restricted-cp-role-%s", testID[:12]),
				SessionName: "test-session",
			},
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{
					NamespaceList: []string{"other-namespace"},
				},
			},
			SourceIdentityRef: &infrav1.AWSIdentityReference{
				Name: controllerIdentity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}
	roleIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
	createObject(g, roleIdentity, ns.Name)
	defer cleanupObject(g, roleIdentity)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rosa-secret-denied-%s", testID[:12]),
			Namespace: ns.Name,
		},
		Data: map[string][]byte{"ocmToken": []byte("fake-token")},
	}
	createObject(g, secret, ns.Name)
	defer cleanupObject(g, secret)

	cpName := fmt.Sprintf("rosa-cp-denied-%s", testID[:12])
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cpName,
			Namespace: ns.Name,
			UID:       types.UID(cpName),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ROSAControlPlane",
			APIVersion: rosacontrolplanev1.GroupVersion.String(),
		},
		Spec: rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   cpName,
			Subnets:           []string{"subnet-0ac99a6230b408813"},
			AvailabilityZones: []string{"us-east-1a"},
			Network: &rosacontrolplanev1.NetworkSpec{
				MachineCIDR: "10.0.0.0/16",
				PodCIDR:     "10.128.0.0/14",
				ServiceCIDR: "172.30.0.0/16",
			},
			Region:           "us-east-1",
			Version:          "4.15.0",
			ChannelGroup:     "stable",
			OIDCID:           "oidcid-test",
			InstallerRoleARN: "arn:aws:iam::123456789012:role/installer",
			WorkerRoleARN:    "arn:aws:iam::123456789012:role/worker",
			SupportRoleARN:   "arn:aws:iam::123456789012:role/support",
			RolesRef: rosacontrolplanev1.AWSRolesRef{
				IngressARN:              "op-arn1",
				ImageRegistryARN:        "op-arn2",
				StorageARN:              "op-arn3",
				NetworkARN:              "op-arn4",
				KubeCloudControllerARN:  "op-arn5",
				NodePoolManagementARN:   "op-arn6",
				ControlPlaneOperatorARN: "op-arn7",
				KMSProviderARN:          "op-arn8",
			},
			CredentialsSecretRef: &corev1.LocalObjectReference{Name: secret.Name},
			VersionGate:          "Acknowledge",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: roleIdentity.Name,
				Kind: infrav1.ClusterRoleIdentityKind,
			},
		},
	}

	ownerCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("owner-cluster-denied-%s", testID[:12]),
			Namespace: ns.Name,
			UID:       types.UID(fmt.Sprintf("owner-cluster-denied-%s", testID[:12])),
		},
		Spec: clusterv1.ClusterSpec{
			ControlPlaneRef: clusterv1.ContractVersionedObjectReference{
				Name:     rosaControlPlane.Name,
				Kind:     "ROSAControlPlane",
				APIGroup: rosacontrolplanev1.GroupVersion.Group,
			},
		},
	}

	rosaControlPlane.OwnerReferences = []metav1.OwnerReference{
		{
			Name:       ownerCluster.Name,
			UID:        ownerCluster.UID,
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
	}

	for _, obj := range []client.Object{ownerCluster, rosaControlPlane} {
		createObject(g, obj, ns.Name)
	}
	defer cleanupObject(g, rosaControlPlane)
	defer cleanupObject(g, ownerCluster)

	awsClientCalled := false
	reconciler := &ROSAControlPlaneReconciler{
		Client: testEnv,
		awsClientFactory: func(_ *scope.ROSAControlPlaneScope) (rosaaws.Client, error) {
			awsClientCalled = true
			return nil, nil
		},
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace},
	}

	// Retry until the informer cache has indexed the control plane and the namespace
	// restriction is enforced, causing scope creation to fail.
	g.Eventually(func(g Gomega) {
		_, errReconcile := reconciler.Reconcile(ctx, req)
		g.Expect(errReconcile).To(HaveOccurred())
		g.Expect(errReconcile.Error()).To(ContainSubstring("failed to create scope"))
		g.Expect(awsClientCalled).To(BeFalse())
	}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(Succeed())
}

// fakeMultiActionStsHandler handles both AssumeRole and GetCallerIdentity STS requests.
// It dispatches based on the Action form value, returning canned XML responses so that
// cross-account creator resolution can be exercised without real AWS credentials.
func fakeMultiActionStsHandler(targetAccount, targetARN string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		action := r.FormValue("Action")
		w.Header().Set("Content-Type", "text/xml")
		switch action {
		case "AssumeRole":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials>
      <AccessKeyId>ASIAIOSFODNN7EXAMPLE</AccessKeyId>
      <SecretAccessKey>wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY</SecretAccessKey>
      <SessionToken>AQoXnyc4MCrrlandlJKwBQ==</SessionToken>
      <Expiration>2030-01-01T00:00:00Z</Expiration>
    </Credentials>
    <AssumedRoleUser>
      <Arn>%s</Arn>
      <AssumedRoleId>ARO123EXAMPLE123:capa-session</AssumedRoleId>
    </AssumedRoleUser>
  </AssumeRoleResult>
  <ResponseMetadata><RequestId>12345678-1234-1234-1234-123456789012</RequestId></ResponseMetadata>
</AssumeRoleResponse>`, targetARN)
		case "GetCallerIdentity":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Account>%s</Account>
    <Arn>%s</Arn>
    <UserId>ARO123EXAMPLE123:capa-session</UserId>
  </GetCallerIdentityResult>
  <ResponseMetadata><RequestId>12345678-1234-1234-1234-123456789012</RequestId></ResponseMetadata>
</GetCallerIdentityResponse>`, targetAccount, targetARN)
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
		}
	}
}

func TestResolveCreatorForTargetAccount(t *testing.T) {
	const (
		sourceAccount   = "111111111111"
		targetAccount   = "999999999999"
		rosaClusterName = "test-cluster"
		// targetARN reflects the actual RoleSessionName generated by resolveCreatorForTargetAccount:
		// "capa-session-<RosaClusterName>"
		targetARN          = "arn:aws:sts::999999999999:assumed-role/installer/capa-session-" + rosaClusterName
		expectedCreatorARN = "arn:aws:iam::999999999999:role/installer" // ROSA normalizes assumed-role ARNs to the base IAM role ARN
	)

	sourceCreator := &rosaaws.Creator{
		AccountID: sourceAccount,
		ARN:       "arn:aws:iam::111111111111:user/management-user",
	}

	t.Run("EmptyInstallerRoleARN_returnsCreatorUnchanged", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					InstallerRoleARN: "",
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		reconciler := &ROSAControlPlaneReconciler{}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result).To(Equal(sourceCreator))
	})

	t.Run("SameAccountInstallerRoleARN_returnsCreatorUnchanged", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					// Same account as sourceCreator
					InstallerRoleARN: "arn:aws:iam::111111111111:role/installer",
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		reconciler := &ROSAControlPlaneReconciler{}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result).To(Equal(sourceCreator))
	})

	t.Run("InvalidInstallerRoleARN_returnsError", func(t *testing.T) {
		g := NewWithT(t)
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					InstallerRoleARN: "not-a-valid-arn",
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		reconciler := &ROSAControlPlaneReconciler{}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("failed to parse account ID from InstallerRoleARN"))
		g.Expect(result).To(BeNil())
	})

	t.Run("RosaRoleConfigRef_NotFound_returnsError", func(t *testing.T) {
		g := NewWithT(t)
		ns, err := testEnv.CreateNamespace(ctx, "test-resolve-creator-notfound")
		g.Expect(err).ToNot(HaveOccurred())

		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
				},
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					RosaRoleConfigRef: &corev1.LocalObjectReference{
						Name: "non-existent-role-config",
					},
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		reconciler := &ROSAControlPlaneReconciler{Client: testEnv}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).To(HaveOccurred())
		g.Expect(result).To(BeNil())
	})

	t.Run("RosaRoleConfigRef_SameAccount_returnsCreatorUnchanged", func(t *testing.T) {
		g := NewWithT(t)
		ns, err := testEnv.CreateNamespace(ctx, "test-resolve-creator-sameaccount")
		g.Expect(err).ToNot(HaveOccurred())

		roleConfig := &expinfrav1.ROSARoleConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "role-config-same-account",
				Namespace: ns.Name,
			},
			Spec: expinfrav1.ROSARoleConfigSpec{
				OidcProviderType: "Managed",
				AccountRoleConfig: expinfrav1.AccountRoleConfig{
					Prefix:  "test",
					Version: "4.19.0",
				},
				OperatorRoleConfig: expinfrav1.OperatorRoleConfig{
					Prefix: "test",
				},
			},
		}
		createObject(g, roleConfig, ns.Name)
		defer cleanupObject(g, roleConfig)

		// Patch status to set InstallerRoleARN to the same account as sourceCreator
		roleConfig.Status = expinfrav1.ROSARoleConfigStatus{
			AccountRolesRef: expinfrav1.AccountRolesRef{
				InstallerRoleARN: "arn:aws:iam::111111111111:role/installer",
			},
		}
		g.Expect(testEnv.Status().Update(ctx, roleConfig)).To(Succeed())

		// The cached client is updated asynchronously after a status write. Wait for
		// the cache to reflect the new InstallerRoleARN before calling into the function
		// under test, otherwise it may see an empty status and return "not ready".
		g.Eventually(func(g Gomega) {
			fresh := &expinfrav1.ROSARoleConfig{}
			g.Expect(testEnv.Get(ctx, client.ObjectKeyFromObject(roleConfig), fresh)).To(Succeed())
			g.Expect(fresh.Status.AccountRolesRef.InstallerRoleARN).ToNot(BeEmpty())
		}).WithTimeout(10 * time.Second).WithPolling(100 * time.Millisecond).Should(Succeed())

		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
				},
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					RosaRoleConfigRef: &corev1.LocalObjectReference{
						Name: roleConfig.Name,
					},
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		reconciler := &ROSAControlPlaneReconciler{Client: testEnv}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result).To(Equal(sourceCreator))
	})

	t.Run("CrossAccountInstallerRoleARN_assumesRoleAndReturnsTargetCreator", func(t *testing.T) {
		g := NewWithT(t)

		stsServer := httptest.NewServer(fakeMultiActionStsHandler(targetAccount, targetARN))
		defer stsServer.Close()

		t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
		t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key-id")
		t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")

		fakeSession, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("us-east-1"))
		g.Expect(err).ToNot(HaveOccurred())

		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					RosaClusterName:  rosaClusterName,
					Region:           "us-east-1",
					InstallerRoleARN: "arn:aws:iam::999999999999:role/installer",
				},
			},
			Logger: *logger.FromContext(context.TODO()),
		}
		rosaScope.SetSession(fakeSession)

		reconciler := &ROSAControlPlaneReconciler{}
		result, err := reconciler.resolveCreatorForTargetAccount(context.TODO(), rosaScope, sourceCreator)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result).ToNot(BeNil())
		g.Expect(result.AccountID).To(Equal(targetAccount))
		g.Expect(result.ARN).To(Equal(expectedCreatorARN))
	})
}

func TestBuildOCMClusterSpec(t *testing.T) {
	// Mock AWS Creator
	mockCreator := &rosaaws.Creator{
		AccountID: "123456789012",
		ARN:       "arn:aws:iam::123456789012:user/test-user",
	}

	// Mock ROSARoleConfig
	mockRoleConfig := &expinfrav1.ROSARoleConfig{
		Status: expinfrav1.ROSARoleConfigStatus{
			OIDCID: "test-oidc-id",
			AccountRolesRef: expinfrav1.AccountRolesRef{
				InstallerRoleARN: "arn:aws:iam::123456789012:role/installer",
				SupportRoleARN:   "arn:aws:iam::123456789012:role/support",
				WorkerRoleARN:    "arn:aws:iam::123456789012:role/worker",
			},
			OperatorRolesRef: rosacontrolplanev1.AWSRolesRef{
				IngressARN:              "arn:aws:iam::123456789012:role/ingress",
				ImageRegistryARN:        "arn:aws:iam::123456789012:role/image-registry",
				StorageARN:              "arn:aws:iam::123456789012:role/storage",
				NetworkARN:              "arn:aws:iam::123456789012:role/network",
				KubeCloudControllerARN:  "arn:aws:iam::123456789012:role/kube-cloud-controller",
				KMSProviderARN:          "arn:aws:iam::123456789012:role/kms-provider",
				ControlPlaneOperatorARN: "arn:aws:iam::123456789012:role/control-plane-operator",
				NodePoolManagementARN:   "arn:aws:iam::123456789012:role/nodepool-management",
			},
		},
	}

	// Test case 1: FIPS enabled
	t.Run("FIPS Enabled", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   "test-cluster",
			Region:            "us-west-2",
			Version:           "4.14.5",
			FIPS:              rosacontrolplanev1.FIPSEnabled, // FIPS enabled
			Subnets:           []string{"subnet-1", "subnet-2"},
			AvailabilityZones: []string{"us-west-2a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, mockRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.FIPS).To(BeTrue())
		g.Expect(ocmSpec.Name).To(Equal("test-cluster"))
		g.Expect(ocmSpec.Region).To(Equal("us-west-2"))
	})

	// Test case 2: FIPS explicitly disabled
	t.Run("FIPS Disabled (Explicit)", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   "test-cluster-no-fips",
			Region:            "us-east-1",
			Version:           "4.14.5",
			FIPS:              rosacontrolplanev1.FIPSDisabled, // FIPS explicitly disabled
			Subnets:           []string{"subnet-1", "subnet-2"},
			AvailabilityZones: []string{"us-east-1a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, mockRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.FIPS).To(BeFalse())
		g.Expect(ocmSpec.Name).To(Equal("test-cluster-no-fips"))
		g.Expect(ocmSpec.Region).To(Equal("us-east-1"))
	})

	// Test case 3: Zero value FIPS (should be false)
	t.Run("FIPS Zero Value", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName: "test-cluster-zero-fips",
			Region:          "us-west-1",
			Version:         "4.14.5",
			// FIPS field not explicitly set (zero value)
			Subnets:           []string{"subnet-1", "subnet-2"},
			AvailabilityZones: []string{"us-west-1a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, mockRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.FIPS).To(BeFalse())
		g.Expect(ocmSpec.Name).To(Equal("test-cluster-zero-fips"))
	})

	// Test case 4: TrustPolicyExternalID set on ROSARoleConfig (takes precedence)
	t.Run("TrustPolicyExternalID Set", func(t *testing.T) {
		g := NewWithT(t)
		roleConfigWithExtID := mockRoleConfig.DeepCopy()
		roleConfigWithExtID.Spec.AccountRoleConfig.TrustPolicyExternalID = "223B9588-36A5-ECA4-BE8D-7C673B77CEC1"

		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   "test-cluster-extid",
			Region:            "us-west-2",
			Version:           "4.14.5",
			Subnets:           []string{"subnet-1", "subnet-2"},
			AvailabilityZones: []string{"us-west-2a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, roleConfigWithExtID, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.ExternalID).To(Equal("223B9588-36A5-ECA4-BE8D-7C673B77CEC1"))
		g.Expect(ocmSpec.Name).To(Equal("test-cluster-extid"))
	})

	// Test case 5: TrustPolicyExternalID empty (should remain empty)
	t.Run("TrustPolicyExternalID Empty", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:   "test-cluster-no-extid",
			Region:            "us-west-2",
			Version:           "4.14.5",
			Subnets:           []string{"subnet-1", "subnet-2"},
			AvailabilityZones: []string{"us-west-2a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, mockRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.ExternalID).To(BeEmpty())
	})

	// Test case 6: TrustPolicyExternalID on controlPlane is ignored when ROSARoleConfig is present
	t.Run("TrustPolicyExternalID ControlPlane Ignored With RoleConfig", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:       "test-cluster-cp-extid",
			Region:                "us-west-2",
			Version:               "4.14.5",
			TrustPolicyExternalID: "should-be-ignored",
			Subnets:               []string{"subnet-1", "subnet-2"},
			AvailabilityZones:     []string{"us-west-2a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, mockRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.ExternalID).To(BeEmpty())
	})

	// Test case 7: TrustPolicyExternalID from controlPlane propagated via synthesized roleConfig
	// (simulates direct-role-ARN flow where reconcileRosaRoleConfig copies the value)
	t.Run("TrustPolicyExternalID ControlPlane Without RoleConfig", func(t *testing.T) {
		g := NewWithT(t)
		controlPlaneSpec := rosacontrolplanev1.RosaControlPlaneSpec{
			RosaClusterName:       "test-cluster-direct",
			Region:                "us-west-2",
			Version:               "4.14.5",
			TrustPolicyExternalID: "direct-external-id",
			Subnets:               []string{"subnet-1", "subnet-2"},
			AvailabilityZones:     []string{"us-west-2a"},
			DefaultMachinePoolSpec: rosacontrolplanev1.DefaultMachinePoolSpec{
				InstanceType: "m5.xlarge",
			},
		}

		// Simulate what reconcileRosaRoleConfig does in the else branch:
		// it copies TrustPolicyExternalID from the control plane spec onto the roleConfig
		synthesizedRoleConfig := mockRoleConfig.DeepCopy()
		synthesizedRoleConfig.Spec.AccountRoleConfig.TrustPolicyExternalID = controlPlaneSpec.TrustPolicyExternalID

		ocmSpec, err := buildOCMClusterSpec(controlPlaneSpec, synthesizedRoleConfig, nil, mockCreator)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ocmSpec.ExternalID).To(Equal("direct-external-id"))
	})
}
