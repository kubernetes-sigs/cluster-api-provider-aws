package controllers

import (
	"bytes"
	"testing"
	"time"

	clustersmgmtv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	corev1 "k8s.io/api/core/v1"
	clocktesting "k8s.io/utils/clock/testing"
	"k8s.io/utils/ptr"

	rosacontrolplanev1beta2 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	testutil "sigs.k8s.io/cluster-api-provider-aws/v2/test/helpers/fixture"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestOCMCluster(t *testing.T) {
	now, err := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z")
	if err != nil {
		t.Fatalf("could not parse time: %v", err)
	}
	clock := clocktesting.NewFakeClock(now)

	for _, testCase := range []struct {
		name string
		in   *rosacontrolplanev1beta2.ROSAControlPlane
	}{
		{
			name: "exhaustive case",
			in: &rosacontrolplanev1beta2.ROSAControlPlane{
				Spec: rosacontrolplanev1beta2.RosaControlPlaneSpec{
					RosaClusterName:   "rosa-cluster-name",
					Subnets:           []string{"subnet-1", "subnet-2"},
					AvailabilityZones: []string{"az-1", "az-2"},
					MachineCIDR:       ptr.To("machine-cidr"),
					Region:            ptr.To("region"),
					Version:           ptr.To("version"),
					ControlPlaneEndpoint: clusterv1.APIEndpoint{
						Host: "example.com",
						Port: 1234,
					},
					RolesRef: rosacontrolplanev1beta2.AWSRolesRef{
						IngressARN:              "ingress-arn",
						ImageRegistryARN:        "image-registry-arn",
						StorageARN:              "storage-arn",
						NetworkARN:              "network-arn",
						KubeCloudControllerARN:  "kube-cloud-controller-arn",
						NodePoolManagementARN:   "node-pool-management-arn",
						ControlPlaneOperatorARN: "control-plane-operator-arn",
						KMSProviderARN:          "kms-provider-arn",
					},
					OIDCID:           ptr.To("oidc-id"),
					AccountID:        ptr.To("account-id"),
					CreatorARN:       ptr.To("creator-arn"),
					InstallerRoleARN: ptr.To("installer-role-arn"),
					SupportRoleARN:   ptr.To("support-role-arn"),
					WorkerRoleARN:    ptr.To("worker-role-arn"),
					CredentialsSecretRef: &corev1.LocalObjectReference{
						Name: "credentials-secret-name",
					},
				},
			},
		},
	} {
		got, err := ocmCluster(&scope.ROSAControlPlaneScope{ControlPlane: testCase.in}, clock.Now)
		if err != nil {
			t.Fatalf("failed to create cluster: %v", err)
		}
		out := bytes.Buffer{}
		if err := clustersmgmtv1.MarshalCluster(got, &out); err != nil {
			t.Fatalf("failed to marshal cluster: %v", err)
		}
		testutil.CompareWithFixture(t, out.Bytes(), testutil.WithExtension(".json"))
	}
}
