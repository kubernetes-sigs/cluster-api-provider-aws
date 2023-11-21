package eks

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts/mock_stsiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/secret"
)

func Test_createCAPIKubeconfigSecret(t *testing.T) {
	testCases := []struct {
		name        string
		input       *eks.Cluster
		serviceFunc func() *Service
		wantErr     bool
	}{
		{
			name: "create kubeconfig secret",
			input: &eks.Cluster{
				CertificateAuthority: &eks.Certificate{Data: aws.String("")},
				Endpoint:             aws.String("https://F00BA4.gr4.us-east-2.eks.amazonaws.com"),
			},
			serviceFunc: func() *Service {
				mockCtrl := gomock.NewController(t)
				stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
				op := request.Request{
					Operation: &request.Operation{Name: "GetCallerIdentity",
						HTTPMethod: "POST",
						HTTPPath:   "/",
					},
					HTTPRequest: &http.Request{
						Header: make(http.Header),
						URL: &url.URL{
							Scheme: "https",
							Host:   "F00BA4.gr4.us-east-2.eks.amazonaws.com",
						},
					},
				}
				stsMock.EXPECT().GetCallerIdentityRequest(gomock.Any()).Return(&op, &sts.GetCallerIdentityOutput{})

				scheme := runtime.NewScheme()
				_ = infrav1.AddToScheme(scheme)
				_ = ekscontrolplanev1.AddToScheme(scheme)
				_ = corev1.AddToScheme(scheme)

				client := fake.NewClientBuilder().WithScheme(scheme).Build()
				managedScope, _ := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
					Client: client,
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
						},
					},
					ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
							UID:       types.UID("1"),
						},
						Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
							EKSClusterName: "cluster-foo",
						},
					},
				})

				service := NewService(managedScope)
				service.STSClient = stsMock
				return service
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			service := tc.serviceFunc()
			clusterRef := types.NamespacedName{
				Namespace: service.scope.Namespace(),
				Name:      service.scope.Name(),
			}
			err := service.createCAPIKubeconfigSecret(context.TODO(), tc.input, &clusterRef)
			if tc.wantErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				var kubeconfigSecret corev1.Secret
				g.Expect(service.scope.Client.Get(context.TODO(), types.NamespacedName{Namespace: "ns", Name: "capi-cluster-foo-kubeconfig"}, &kubeconfigSecret)).To(BeNil())
				g.Expect(kubeconfigSecret.Data).ToNot(BeNil())
				g.Expect(len(kubeconfigSecret.Data)).To(BeIdenticalTo(3))
				g.Expect(kubeconfigSecret.Data[secret.KubeconfigDataName]).ToNot(BeEmpty())
				g.Expect(kubeconfigSecret.Data[relativeKubeconfigKey]).ToNot(BeEmpty())
				g.Expect(kubeconfigSecret.Data[relativeTokenFileKey]).ToNot(BeEmpty())
			}
		})
	}
}

func Test_updateCAPIKubeconfigSecret(t *testing.T) {
	type testCase struct {
		name        string
		input       *eks.Cluster
		secret      *corev1.Secret
		serviceFunc func(tc testCase) *Service
		wantErr     bool
	}
	testCases := []testCase{
		{
			name: "update kubeconfig secret",
			input: &eks.Cluster{
				Name:                 aws.String("cluster-foo"),
				CertificateAuthority: &eks.Certificate{Data: aws.String("")},
				Endpoint:             aws.String("https://F00BA4.gr4.us-east-2.eks.amazonaws.com"),
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "capi-cluster-foo-kubeconfig",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "controlplane.cluster.x-k8s.io/v1beta2",
							Kind:       "AWSManagedControlPlane",
							Name:       "capi-cluster-foo",
							UID:        "1",
							Controller: aws.Bool(true),
						},
					},
				},
				Data: make(map[string][]byte),
			},
			serviceFunc: func(tc testCase) *Service {
				mockCtrl := gomock.NewController(t)
				stsMock := mock_stsiface.NewMockSTSAPI(mockCtrl)
				op := request.Request{
					Operation: &request.Operation{Name: "GetCallerIdentity",
						HTTPMethod: "POST",
						HTTPPath:   "/",
					},
					HTTPRequest: &http.Request{
						Header: make(http.Header),
						URL: &url.URL{
							Scheme: "https",
							Host:   "F00BA4.gr4.us-east-2.eks.amazonaws.com",
						},
					},
				}
				stsMock.EXPECT().GetCallerIdentityRequest(gomock.Any()).Return(&op, &sts.GetCallerIdentityOutput{})

				scheme := runtime.NewScheme()
				_ = infrav1.AddToScheme(scheme)
				_ = ekscontrolplanev1.AddToScheme(scheme)
				_ = corev1.AddToScheme(scheme)

				client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(tc.secret).Build()
				managedScope, _ := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
					Client: client,
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
						},
					},
					ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
							UID:       "1",
						},
						Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
							EKSClusterName: "cluster-foo",
						},
					},
				})

				service := NewService(managedScope)
				service.STSClient = stsMock
				return service
			},
		},
		{
			name: "detect incorrect ownership on the kubeconfig secret",
			input: &eks.Cluster{
				Name:                 aws.String("cluster-foo"),
				CertificateAuthority: &eks.Certificate{Data: aws.String("")},
				Endpoint:             aws.String("https://F00BA4.gr4.us-east-2.eks.amazonaws.com"),
			},
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "capi-cluster-foo-kubeconfig",
				},
				Data: make(map[string][]byte),
			},
			serviceFunc: func(tc testCase) *Service {
				scheme := runtime.NewScheme()
				_ = infrav1.AddToScheme(scheme)
				_ = ekscontrolplanev1.AddToScheme(scheme)
				_ = corev1.AddToScheme(scheme)

				client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(tc.secret).Build()
				managedScope, _ := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
					Client: client,
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
						},
					},
					ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns",
							Name:      "capi-cluster-foo",
							UID:       "1",
						},
						Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
							EKSClusterName: "cluster-foo",
						},
					},
				})

				service := NewService(managedScope)
				return service
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			service := tc.serviceFunc(tc)
			err := service.updateCAPIKubeconfigSecret(context.TODO(), tc.secret, tc.input)
			if tc.wantErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
				var kubeconfigSecret corev1.Secret
				g.Expect(service.scope.Client.Get(context.TODO(), types.NamespacedName{Namespace: "ns", Name: "capi-cluster-foo-kubeconfig"}, &kubeconfigSecret)).To(BeNil())
				g.Expect(kubeconfigSecret.Data).ToNot(BeNil())
				g.Expect(len(kubeconfigSecret.Data)).To(BeIdenticalTo(3))
				g.Expect(kubeconfigSecret.Data[secret.KubeconfigDataName]).ToNot(BeEmpty())
				g.Expect(kubeconfigSecret.Data[relativeKubeconfigKey]).ToNot(BeEmpty())
				g.Expect(kubeconfigSecret.Data[relativeTokenFileKey]).ToNot(BeEmpty())
			}
		})
	}
}
