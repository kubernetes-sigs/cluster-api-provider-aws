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
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	restclient "k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsiface "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

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
			UID:       types.UID("rosa-control-plane-1")},
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
			InstallerRoleARN: "arn1",
			WorkerRoleARN:    "arn2",
			SupportRoleARN:   "arn3",
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
			Conditions: clusterv1.Conditions{clusterv1.Condition{
				Type:     "Paused",
				Status:   "False",
				Severity: "",
				Reason:   "NotPaused",
				Message:  "",
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
			ControlPlaneRef: &corev1.ObjectReference{
				Name:       rosaControlPlane.Name,
				Kind:       "ROSAControlPlane",
				APIVersion: rosacontrolplanev1.GroupVersion.String(),
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
	stsMock := mock_stsiface.NewMockSTSClient(mockCtrl)

	getCallerIdentityResult := &stsv2.GetCallerIdentityOutput{Account: aws.String("foo"), Arn: aws.String("arn:aws:iam::123456789012:rosa/foo")}
	stsMock.EXPECT().GetCallerIdentity(gomock.Any(), gomock.Any()).Return(getCallerIdentityResult, nil).Times(1)

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
		m.UpdateCluster(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(clusterKey string, creator *rosaaws.Creator, config ocm.Spec) error {
			return nil
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
		Conditions: clusterv1.Conditions{clusterv1.Condition{
			Type:     "Paused",
			Status:   "False",
			Severity: "",
			Reason:   "NotPaused",
			Message:  "",
		}},
	}

	g.Expect(cpPh.Patch(ctx, rosaControlPlane)).To(Succeed())

	time.Sleep(50 * time.Millisecond)

	cp := &rosacontrolplanev1.ROSAControlPlane{}
	key := client.ObjectKey{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace}
	errGet := testEnv.Get(ctx, key, cp)
	g.Expect(errGet).NotTo(HaveOccurred())
	oldCondition := conditions.Get(cp, clusterv1.PausedV1Beta2Condition)
	g.Expect(oldCondition).NotTo(BeNil())

	r := ROSAControlPlaneReconciler{
		WatchFilterValue: "",
		Client:           testEnv,
		restClientConfig: cfg,
		NewStsClient: func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSClient {
			return stsMock
		},
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
