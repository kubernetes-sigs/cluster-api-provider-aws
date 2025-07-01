package rosa

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	sdk "github.com/openshift-online/ocm-sdk-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

func createROSAControlPlaneScopeWithSecrets(cp *rosacontrolplanev1.ROSAControlPlane, secrets ...*corev1.Secret) *scope.ROSAControlPlaneScope {
	// k8s mock (fake) client
	fakeClientBuilder := fake.NewClientBuilder()
	for _, sec := range secrets {
		fakeClientBuilder.WithObjects(sec)
	}

	fakeClient := fakeClientBuilder.Build()

	// ROSA Control Plane Scope
	rcpScope := &scope.ROSAControlPlaneScope{
		Client:       fakeClient,
		ControlPlane: cp,
		Logger:       *logger.NewLogger(klog.Background()),
	}

	return rcpScope
}

func createSecret(name, namespace, token, url, clientID, clientSecret string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"ocmToken":        []byte(token),
			"ocmApiUrl":       []byte(url),
			"ocmClientID":     []byte(clientID),
			"ocmClientSecret": []byte(clientSecret),
		},
	}
}

func createCP(name string, namespace string, credSecretName string) *rosacontrolplanev1.ROSAControlPlane {
	return &rosacontrolplanev1.ROSAControlPlane{
		Spec: rosacontrolplanev1.RosaControlPlaneSpec{
			CredentialsSecretRef: &corev1.LocalObjectReference{
				Name: credSecretName,
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func TestNewOCMRawConnection(t *testing.T) {
	g := NewWithT(t)
	wlSecret := createSecret("rosa-hcp-creds-secret", "default", "fake-token", "https://api.stage.openshift.com", "", "")
	cp := createCP("rosa-hcp-cp", "default", "rosa-hcp-creds-secret")
	rcpScope := createROSAControlPlaneScopeWithSecrets(cp, wlSecret)

	conn, _ := newOCMRawConnection(context.Background(), rcpScope)
	g.Expect(conn.Agent()).To(Equal(capaAgentName + "/" + version.Get().GitVersion + " " + sdk.DefaultAgent))
}
func TestOcmCredentials(t *testing.T) {
	g := NewWithT(t)

	wlSecret := createSecret("rosa-creds-secret", "default", "", "url", "client-id", "client-secret")
	mgrSecret := createSecret("rosa-creds-secret", system.GetManagerNamespace(), "", "url", "global-client-id", "global-client-secret")

	cp := createCP("rosa-cp", "default", "rosa-creds-secret")

	// Test that ocmCredentials() prefers workload secret to global and environment secrets
	os.Setenv("OCM_API_URL", "env-url")
	os.Setenv("OCM_TOKEN", "env-token")
	rcpScope := createROSAControlPlaneScopeWithSecrets(cp, wlSecret, mgrSecret)
	token, url, clientID, clientSecret, err := ocmCredentials(context.Background(), rcpScope)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(token).To(Equal(string(wlSecret.Data["ocmToken"])))
	g.Expect(url).To(Equal(string(wlSecret.Data["ocmApiUrl"])))
	g.Expect(clientID).To(Equal(string(wlSecret.Data["ocmClientID"])))
	g.Expect(clientSecret).To(Equal(string(wlSecret.Data["ocmClientSecret"])))

	// Test that ocmCredentials() prefers global manager secret to environment secret in case workload secret is not specified
	cp.Spec = rosacontrolplanev1.RosaControlPlaneSpec{}
	rcpScope = createROSAControlPlaneScopeWithSecrets(cp, mgrSecret)
	token, url, clientID, clientSecret, err = ocmCredentials(context.Background(), rcpScope)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(token).To(Equal(string(mgrSecret.Data["ocmToken"])))
	g.Expect(url).To(Equal(string(mgrSecret.Data["ocmApiUrl"])))
	g.Expect(clientID).To(Equal(string(mgrSecret.Data["ocmClientID"])))
	g.Expect(clientSecret).To(Equal(string(mgrSecret.Data["ocmClientSecret"])))

	// Test that ocmCredentials() returns environment secret in case workload and manager secret are not specified
	cp.Spec = rosacontrolplanev1.RosaControlPlaneSpec{}
	rcpScope = createROSAControlPlaneScopeWithSecrets(cp)
	token, url, clientID, clientSecret, err = ocmCredentials(context.Background(), rcpScope)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(token).To(Equal(os.Getenv("OCM_TOKEN")))
	g.Expect(url).To(Equal(os.Getenv("OCM_API_URL")))
	g.Expect(clientID).To(Equal(""))
	g.Expect(clientSecret).To(Equal(""))

	// Test that ocmCredentials() returns error in case none of the secrets has been provided
	os.Unsetenv("OCM_API_URL")
	os.Unsetenv("OCM_TOKEN")
	token, url, clientID, clientSecret, err = ocmCredentials(context.Background(), rcpScope)

	g.Expect(err).To(HaveOccurred())
	g.Expect(token).To(Equal(""))
	g.Expect(url).To(Equal(""))
	g.Expect(clientID).To(Equal(""))
	g.Expect(clientSecret).To(Equal(""))
}
