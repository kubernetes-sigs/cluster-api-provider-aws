package iam

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"

	v1certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := infrav1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := admissionv1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	return scheme, nil
}

func TestReconcileMutatingWebhook(t *testing.T) {
	nsString := "cluster-test"
	testSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-test",
			Namespace: nsString,
		},
		Data: map[string][]byte{
			"ca.crt": []byte("myCACert"),
		},
	}

	namespacedSecret := fmt.Sprintf("%s/%s", nsString, testSecret.Name)
	mwhMeta := objectMeta(podIdentityWebhookName, nsString)
	mwhMeta.UID = "uid-test"
	mutate := "/mutate"
	fail := admissionv1.Ignore
	none := admissionv1.SideEffectClassNone
	testPodIdentityWebhookConfiguration := &admissionv1.MutatingWebhookConfiguration{
		ObjectMeta: mwhMeta,
		Webhooks: []admissionv1.MutatingWebhook{
			{
				Name:          podIdentityWebhookName + ".amazonaws.com",
				FailurePolicy: &fail,
				ClientConfig: admissionv1.WebhookClientConfig{
					Service: &admissionv1.ServiceReference{
						Name:      podIdentityWebhookName,
						Namespace: nsString,
						Path:      &mutate,
					},
					CABundle: []byte(""),
				},
				Rules: []admissionv1.RuleWithOperations{
					{
						Operations: []admissionv1.OperationType{admissionv1.Create},
						Rule: admissionv1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"pods"},
						},
					},
				},
				SideEffects:             &none,
				AdmissionReviewVersions: []string{"v1beta1"},
			},
		},
	}

	oldAnnotationKey := "annotationsAlreadyExist"
	oldAnnotationValue := "test"
	oldAnnotation := map[string]string{
		oldAnnotationKey: oldAnnotationValue,
	}

	testCases := []struct {
		name                  string
		ctx                   context.Context
		ns                    string
		secret                *corev1.Secret
		setup                 func(*testing.T, client.Client)
		expectedAnnotationVal string
		hasOldAnnotations     bool
	}{
		{
			name:                  "initialize new MutatingWebhookConfiguration, ensure Annotation map is created and cert-manager annotation is set",
			ctx:                   context.TODO(),
			ns:                    nsString,
			secret:                testSecret,
			expectedAnnotationVal: namespacedSecret,
		}, {
			name: "update existing MutatingWebhookConfiguration that has no Annotations, ensure cert-manager annotation is set",
			ctx:  context.TODO(),
			ns:   nsString,
			setup: func(t *testing.T, c client.Client) {
				t.Helper()
				testPodIdentityWebhookConfiguration.Annotations = nil
				testPodIdentityWebhookConfiguration.SetGroupVersionKind(infrav1.GroupVersion.WithKind("MutatingWebhookConfiguration"))
				err := c.Create(context.TODO(), testPodIdentityWebhookConfiguration)
				if err != nil {
					t.Fatal(err)
				}
			},
			secret:                testSecret,
			expectedAnnotationVal: namespacedSecret,
		}, {
			name: "update existing MutatingWebhookConfiguration that contains existing Annotations, ensure cert-manager annotation is set",
			ctx:  context.TODO(),
			ns:   nsString,
			setup: func(t *testing.T, c client.Client) {
				t.Helper()
				testPodIdentityWebhookConfiguration.ResourceVersion = ""
				testPodIdentityWebhookConfiguration.Annotations = oldAnnotation
				err := c.Create(context.TODO(), testPodIdentityWebhookConfiguration)
				if err != nil {
					t.Fatal(err)
				}
			},
			secret:                testSecret,
			expectedAnnotationVal: namespacedSecret,
			hasOldAnnotations:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()
			if tc.setup != nil {
				tc.setup(t, k8sClient)
			}

			// perform reconcile
			err = reconcileMutatingWebHook(context.TODO(), nsString, testSecret, k8sClient)
			g.Expect(err).To(BeNil())
			mwhc := &admissionv1.MutatingWebhookConfiguration{}

			err = k8sClient.Get(context.TODO(), types.NamespacedName{
				Name:      podIdentityWebhookName,
				Namespace: nsString,
			}, mwhc)
			g.Expect(err).To(BeNil())

			if !g.Expect(mwhc.Annotations[v1certmanager.WantInjectAnnotation]).To(Equal(tc.expectedAnnotationVal)) {
				t.Fatalf("Expected %s annotation to equal %s but got '%s'\n", v1certmanager.WantInjectAnnotation, tc.expectedAnnotationVal, mwhc.Annotations[v1certmanager.WantInjectAnnotation])
			}

			// for test cases with non empty Annotations maps, ensure old annotations are preserved
			if tc.hasOldAnnotations {
				g.Expect(mwhc.Annotations[oldAnnotationKey]).To(Equal(oldAnnotationValue))
			}
		})
	}
}
