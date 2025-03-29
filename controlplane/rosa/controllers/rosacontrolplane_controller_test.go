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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	sts "github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3/mock_stsiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		mockCluster, _ := cmv1.NewCluster().
			AWS(cmv1.NewAWS().
				AuditLog(cmv1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/AuditLogRole"))).
			RegistryConfig(cmv1.NewClusterRegistryConfig().
				AdditionalTrustedCa(map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
				AllowedRegistriesForImport(cmv1.NewRegistryLocation().
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

		mockCluster, _ := cmv1.NewCluster().
			AWS(cmv1.NewAWS().
				AuditLog(cmv1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/OldAuditLogRole"))).Build()

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

		mockCluster, _ := cmv1.NewCluster().
			RegistryConfig(cmv1.NewClusterRegistryConfig().
				AdditionalTrustedCa(map[string]string{"old-trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
				AllowedRegistriesForImport(cmv1.NewRegistryLocation().
					DomainName("old-registry.com").
					Insecure(false)).RegistrySources(cmv1.NewRegistrySources().BlockedRegistries([]string{"blocked.io", "blocked.org"}...))).
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

		mockCluster, _ := cmv1.NewCluster().
			RegistryConfig(cmv1.NewClusterRegistryConfig().
				AllowedRegistriesForImport(cmv1.NewRegistryLocation().
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
}

func TestRosaControlPlaneReconcile(t *testing.T) {
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
	identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterStaticIdentity"))

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
			Region:           "us-east-1",
			Version:          "4.15.20",
			ChannelGroup:     "stable",
			RolesRef:         rosacontrolplanev1.AWSRolesRef{},
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
			Ready:   true,
			ID:      "rosa-control-plane-1",
			Version: "4.15.20",
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

	rosaControlPlane.OwnerReferences = []metav1.OwnerReference{
		{
			Name:       ownerCluster.Name,
			UID:        ownerCluster.UID,
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
	}

	mockCtrl := gomock.NewController(t)
	// recorder := record.NewFakeRecorder(10)
	ctx := context.TODO()
	ocmMock := mocks.NewMockOCMClient(mockCtrl)
	stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)

	getCallerIdentityResult := &sts.GetCallerIdentityOutput{Account: aws.String("foo"), Arn: aws.String("arn:aws:iam::123456789012:rosa/foo")}
	stsMock.EXPECT().GetCallerIdentity(gomock.Any()).Return(getCallerIdentityResult, nil).Times(1)

	expect := func(m *mocks.MockOCMClientMockRecorder) {
		m.ValidateHypershiftVersion(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterId string, nodePoolID string) (bool, error) {
			return true, nil
		}).Times(1)
		m.GetCluster(gomock.Any(), gomock.Any()).DoAndReturn(func(clusterKey string, creator *rosaaws.Creator) (*v1.Cluster, error) {
			// Mock Cluster input

			sts := (&v1.STSBuilder{}).OIDCEndpointURL("oidc.com/oidc1")
			aws := cmv1.NewAWS().AuditLog(cmv1.NewAuditLog().RoleArn("arn:aws:iam::123456789012:role/AuditLogRole")).STS(sts)
			console := (&v1.ClusterConsoleBuilder{}).URL("console.redhat.com/cluster-123")
			status := (&cmv1.ClusterStatusBuilder{}).State(cmv1.ClusterStateError)
			version := (&cmv1.VersionBuilder{}).RawID(rosaControlPlane.Spec.Version)
			mockCluster, _ := cmv1.NewCluster().AWS(aws).ID("cluster-1").Version(version).Status(status).Console(console).
				RegistryConfig(cmv1.NewClusterRegistryConfig().
					AdditionalTrustedCa(map[string]string{"trusted-ca": "-----BEGIN CERTIFICATE----- testcert -----END CERTIFICATE-----"}).
					AllowedRegistriesForImport(cmv1.NewRegistryLocation().
						DomainName("registry1.com").
						Insecure(false))).
				Build()
			return mockCluster, nil
		}).Times(1)
	}

	expect(ocmMock.EXPECT())

	r := ROSAControlPlaneReconciler{
		// Recorder:         recorder,
		WatchFilterValue: "",
		Endpoints:        []scope.ServiceEndpoint{},
		Client:           testEnv,
		NewStsClient:     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI { return stsMock },
		NewOCMClient: func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error) {
			return ocmMock, nil
		},
	}

	objects := []client.Object{ownerCluster, rosaControlPlane, secret, identity}

	for _, obj := range objects {
		createObject(g, obj, ns.Name)
	}

	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace}
	_, errReconcile := r.Reconcile(ctx, req)
	g.Expect(errReconcile).ToNot(HaveOccurred())
	// g.Expect(result).To(Equal("aa"))
	time.Sleep(50 * time.Millisecond)

	m := &rosacontrolplanev1.ROSAControlPlane{}
	key := client.ObjectKey{Name: rosaControlPlane.Name, Namespace: rosaControlPlane.Namespace}
	errGet := testEnv.Get(ctx, key, m)
	g.Expect(errGet).NotTo(HaveOccurred())
	g.Expect(m.Status.Version).To(Equal("ahoj"))

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
