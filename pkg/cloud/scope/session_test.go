package scope

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestPrincipalParsing(t *testing.T) {
	testCases := []struct {
		name string
		awsCluster infrav1.AWSCluster
		principalRef *corev1.ObjectReference
		principal    runtime.Object
		setup func(client.Client, *testing.T)
		expect func(*session.Session) bool
	}{
		{
			name: "Can get a session for a static principal",
			setup: func(c client.Client, t *testing.T) {
				awsCluster := &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta {
						Name: "cluster1",
						Namespace: "default",
					},
					Spec: infrav1.AWSClusterSpec {
						PrincipalRef: corev1.ObjectReference{
							Name: "test1",
							Kind: "AWSClusterStaticPrincipal",
						},
					},
				}
				awsCluster.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSCluster"))
				err := c.Create(context.Background(), awsCluster)
				if err != nil {
					t.Fatal(err)
				}

				principal := &infrav1.AWSClusterStaticPrincipal {
					TypeMeta: metav1.TypeMeta {
						APIVersion:  infrav1.GroupVersion.String(),
						Kind: "AWSClusterStaticPrincipal",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "principal",
					},
					Spec: infrav1.AWSClusterStaticPrincipalSpec {
						SecretRef: corev1.ObjectReference{
							Name: "static-credentials-secret",
							Namespace: "default",

						},
					},
				}
				//principal.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterStaticPrincipal"))
				err = c.Create(context.Background(), principal)
				if err != nil {
					t.Fatal(err)
				}

				credentialsSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta {
						Name: "static-credentials-secret",
						Namespace: "default",
					},
					Data: map[string][]byte {
						"accessKeyID": []byte("1234567890"),
						"SecretAccessKey": []byte("abcdefghijklmnop"),
						"SessionToken": []byte("asdfasdfasdf"),
					},
				}
				credentialsSecret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Kind: "Secret", Version: "v1"})
				err = c.Create(context.Background(), credentialsSecret)
				if err != nil {
					t.Fatal(err)
				}
			},
			expect: func(s *session.Session) bool {
				// session is a static principal
				v, err := s.Config.Credentials.Get()
				if err != nil {
					t.Fatal(err)
				}
				v.AccessKeyID="1234567890"
				v.SecretAccessKey="abcdefghijklmnop"
				v.SessionToken="asdfasdfasdf"
				return true
			},
		},
	}

	scheme := runtime.NewScheme()
	err := infrav1.AddToScheme(scheme)
	if err != nil {
		t.Fatal(err)
	}
	k8sClient := fake.NewFakeClientWithScheme(scheme)
	for _, tc := range testCases {
		tc.setup(k8sClient, t)
		s, err := sessionForRegion(&tc.awsCluster, "us-east-1")
		if err != nil {
			t.Fatal(err)
		}
		if !tc.expect(s) {
			t.Fatal("Failed: " + tc.name)
		}
	}
}