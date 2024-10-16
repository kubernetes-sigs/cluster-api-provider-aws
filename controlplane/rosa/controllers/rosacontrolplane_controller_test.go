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
	"testing"

	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
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
