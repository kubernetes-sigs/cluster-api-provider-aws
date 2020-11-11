package client

import (
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestUseCustomCABundle(t *testing.T) {
	cases := []struct {
		name             string
		cm               *corev1.ConfigMap
		expectedCABundle string
	}{
		{
			name: "no configmap",
		},
		{
			name: "no CA bundle in configmap",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "openshift-config-managed",
					Name:      "kube-cloud-config",
				},
				Data: map[string]string{
					"other-key": "other-data",
				},
			},
		},
		{
			name: "custom CA bundle",
			cm: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "openshift-config-managed",
					Name:      "kube-cloud-config",
				},
				Data: map[string]string{
					"ca-bundle.pem": "a custom bundle",
				},
			},
			expectedCABundle: "a custom bundle",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			corev1.AddToScheme(scheme)
			resources := []runtime.Object{}
			if tc.cm != nil {
				resources = append(resources, tc.cm)
			}
			ctrlRuntimeClient := fake.NewFakeClientWithScheme(scheme, resources...)
			awsOptions := &session.Options{}
			err := useCustomCABundle(awsOptions, ctrlRuntimeClient)
			if err != nil {
				t.Fatalf("unexpected error from useCustomCABundle: %v", err)
			}
			actualCABundle := ""
			if awsOptions.CustomCABundle != nil {
				bundleBytes, err := ioutil.ReadAll(awsOptions.CustomCABundle)
				if err != nil {
					t.Fatalf("unexpected error reading bundle: %v", err)
				}
				actualCABundle = string(bundleBytes)
			}
			if a, e := actualCABundle, tc.expectedCABundle; a != e {
				t.Errorf("unexpected CA bundle: expected=%s; got %s", e, a)
			}
		})
	}
}
