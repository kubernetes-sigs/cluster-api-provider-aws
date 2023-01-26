/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestOIDCReconcile(t *testing.T) {
	tests := []struct {
		name    string
		expect  func(m *mock_iamauth.MockIAMAPIMockRecorder, url string)
		cluster func(url string) eks.Cluster
	}{
		{
			name: "cluster create with no OIDC provider present yet should create one",
			cluster: func(url string) eks.Cluster {
				return eks.Cluster{
					Name:    aws.String("cluster-test"),
					Arn:     aws.String("arn:arn"),
					RoleArn: aws.String("arn:role"),
					Identity: &eks.Identity{
						Oidc: &eks.OIDC{
							Issuer: aws.String(url),
						},
					},
				}
			},
			expect: func(m *mock_iamauth.MockIAMAPIMockRecorder, url string) {
				m.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{}).Return(&iam.ListOpenIDConnectProvidersOutput{
					OpenIDConnectProviderList: []*iam.OpenIDConnectProviderListEntry{},
				}, nil)
				m.CreateOpenIDConnectProvider(&iam.CreateOpenIDConnectProviderInput{
					ClientIDList:   aws.StringSlice([]string{"sts.amazonaws.com"}),
					ThumbprintList: aws.StringSlice([]string{"15dbd260c7465ecca6de2c0b2181187f66ee0d1a"}),
					Url:            &url,
				}).Return(&iam.CreateOpenIDConnectProviderOutput{
					OpenIDConnectProviderArn: aws.String("arn::oidc"),
				}, nil)
				m.TagOpenIDConnectProvider(&iam.TagOpenIDConnectProviderInput{
					OpenIDConnectProviderArn: aws.String("arn::oidc"),
					Tags:                     []*iam.Tag{},
				}).Return(&iam.TagOpenIDConnectProviderOutput{}, nil)
			},
		},
		{
			name: "cluster create with existing OIDC provider which is retrieved",
			cluster: func(url string) eks.Cluster {
				return eks.Cluster{
					Name:    aws.String("cluster-test"),
					Arn:     aws.String("arn:arn"),
					RoleArn: aws.String("arn:role"),
					Identity: &eks.Identity{
						Oidc: &eks.OIDC{
							Issuer: aws.String(url),
						},
					},
				}
			},
			expect: func(m *mock_iamauth.MockIAMAPIMockRecorder, url string) {
				m.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{}).Return(&iam.ListOpenIDConnectProvidersOutput{
					OpenIDConnectProviderList: []*iam.OpenIDConnectProviderListEntry{
						{
							Arn: aws.String("arn::oidc"),
						},
					},
				}, nil)
				// This should equal with what we provide.
				m.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{
					OpenIDConnectProviderArn: aws.String("arn::oidc"),
				}).Return(&iam.GetOpenIDConnectProviderOutput{
					ClientIDList:   aws.StringSlice([]string{"sts.amazonaws.com"}),
					ThumbprintList: aws.StringSlice([]string{"15dbd260c7465ecca6de2c0b2181187f66ee0d1a"}),
					Url:            &url,
				}, nil)
				m.TagOpenIDConnectProvider(&iam.TagOpenIDConnectProviderInput{
					OpenIDConnectProviderArn: aws.String("arn::oidc"),
					Tags:                     []*iam.Tag{},
				}).Return(&iam.TagOpenIDConnectProviderOutput{}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			_ = corev1.AddToScheme(scheme)

			ts := httptest.NewTLSServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Send response to be tested
				rw.WriteHeader(http.StatusOK)
				rw.Write([]byte(`OK`))
			}))
			defer ts.Close()

			controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source",
					Namespace: "ns",
				},
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					Version:               aws.String("1.25"),
					AssociateOIDCProvider: true,
				},
			}
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "capi-name-kubeconfig",
					Namespace: "ns",
				},
				Data: map[string][]byte{
					"value": kubeConfig,
				},
			}
			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(controlPlane, secret).Build()
			scope, _ := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      "capi-name",
					},
				},
				ControlPlane: controlPlane,
				EnableIAM:    true,
			})

			iamMock := mock_iamauth.NewMockIAMAPI(mockControl)
			tc.expect(iamMock.EXPECT(), ts.URL)
			s := NewService(scope, WithIAMClient(ts.Client()))
			s.IAMClient = iamMock

			cluster := tc.cluster(ts.URL)
			err := s.reconcileOIDCProvider(&cluster)
			// We reached the trusted policy reconcile which will fail because it tries to connect to the server.
			// But at this point, we already know that the critical area has been covered.
			g.Expect(err).To(MatchError(ContainSubstring("dial tcp: lookup test-cluster-api.nodomain.example.com")))
		})
	}
}

var kubeConfig = []byte(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1Ea3lPVEl3TWpnek1Gb1hEVE15TURreU5qSXdNamd6TUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTFBqCldzdlNEM1kxR1daNGpPSDdqdm4zNUNKUUJvdm8vVnljN3BQdHB6OEVWaFJNNnpRSTMrU2EvdDZyMWdSeHcwM1QKalhHTlRvamNOU0dUVGhHSnN6K28vRjc0Tml5enN2bk5zaThHem9rRU42QmpVU1NmeDg2RVZrM3J4ekVkeFhEaQpoZmNmcDFrQkJBa3lyMGltUGlSZDBaWGFSTnA1dEhldDI3eXp4TTBLZDRjRUxPcHJQc1QzRlp0bGNQTU01YVhzCmhzcGR6dkpmMFNUeWtCNWRtUmU4WHVEc0VDeVgvSTBVbVdXNVkvaWRDMmN0WUE5bEExMjdlcEFrMFUwb29lU2IKUWdMZ0tScjJ6UUl0UzlVNEdEaUttdGVsZVI2dWg2bjdTRkxRZVl1Y25hYXZWb0lBZ1A0UDRRejVMVi9OelF3bApVZENFR3lPODkxQVVDVkIrbENNQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZOTmc4L2ZueStjUW40YmZhdzVGRE1ld3kvcm9NQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBRkFMSGZTcFNmTVQwQjlpTGFPVgpqN3d0V2I5MUY1Y25uZm15QWhVbm5qZXk2a1YyOFlNSWl2enZ3ZlQ5ZmVPb1llT3p4Tldla3YxUUVEaGw3a2pZCmJ2L2V5SDR3dDZGSHVCNlduQ1lESVpzK1doaXNURmd6bE5QeDJ0UVZLYjhzdy9WdGI0UU1WWVp3QjdpT2p6V1QKWjA4MW5MTmJpalQ3eEdIeWRWMWQ0SDR5eS91ajNJdWU3bkxYNHFPZk9udi8wQ2Vvb2Evd0VQeG1HMjJYb09WZgpzSlRWZnhrK1Zpak1Fc1kzRmZidWR1d3llNHc0cmxmUXhCNFZtbE1INEFrRmFvT1hLTGdGS3FrQkFLNVgwekhKClQvWWJkTm9jOThlcnJRNXZkRXhDZkV4RjFCWWtnbUVwcGZOV2UwK01xekgwZ2RTTTBzNEFBUmhrME4xNWRwVXoKeTBnPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://test-cluster-api.nodomain.example.com:6443
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJRUlJK3dnRVhXWVF3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpBNU1qa3lNREk0TXpCYUZ3MHlNekE1TWpreU1ESTRNekZhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXFNcHJqYU1NaXdTN2dTOXgKcnd5VGFJRDZncUxOcklpMEh3SnRzTVpMRERTWTJLZWd5VCtwaDcxNjc4bHB3SHk2Njg3dVJ0WXpQcU90cXVNRQpEcTQzdmpxMzNMTng1ZVI4SnhSTk53d2Q3VXJZNmt4R2U1UUF3MXdWRW9OcmZTZk1BdTBOMEtIb1FKRjhEZDNlCjNFTVl5YmxySEYzMlN5MnluNHpWMmZRMDdpV2RUa2x3WDNZbkpTcFlFRTFDM3k2NjFHVVdGSXZCZm03b2NuOFUKeGdzQ01XNkxrbzVaMXh4OGVzZm5SSU5oZHZnS1BuN3dQSEtMUEQzRDNNUUdwM2V2QVVIWVExclpnTXRJNDQ3WQpVeVlkSFo5NDVLcVpRZ1ppS2FCdE1lblpmcVJndzArckNkeG5qMTRFYSt4RmtPMlNYQ0wrRWtNaUdvaFU3T29RCnVnay80d0lEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JUVFlQUDM1OHZuRUorRzMyc09SUXpIc012Ngo2REFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBTklMSVp3TDhZdlh0QTRqZm9VNlBlYWc4bFBSaWQ2TDlEdTcxCm03NDZMRWQ0cVMvQ2VFb1Z3Q0JqUnplQytLRkcvMERFa3JvYXRTbzhQMXZQQVkxYm5yN3FJdmo3S0RIS25lWFcKSS9saGo4M3ZyYmhoRjN1TXVnTXRiaUI0cnB0eUxjMjc5cGpnWDJqMkFxN09OUDNnVVJoVmJBZG1JTmQwNlVhYQpnaTR1dFBGV1Z2cENsTlpKWXhqUnJVZzJCR0JSQ0RQVU9JWkVkeHBVRnQ5cWsrWWxva0RQb29lR1QzVGlKNnE3ClJwS01UQ04yOWo3cS96cEwxSlNGNXFEVWprdXV5eWd3aUNUcXR0SVBwajAvaU5kak9TVGJlcG5sMUdPOTVNTUEKbGN6NzQ4NEt1dTlGSEtTcjhvcHVEK1hWYXBRbWpuZVdIYmtQUVo4elMwTGExdHc1V2c9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBcU1wcmphTU1pd1M3Z1M5eHJ3eVRhSUQ2Z3FMTnJJaTBId0p0c01aTEREU1kyS2VnCnlUK3BoNzE2NzhscHdIeTY2ODd1UnRZelBxT3RxdU1FRHE0M3ZqcTMzTE54NWVSOEp4Uk5Od3dkN1VyWTZreEcKZTVRQXcxd1ZFb05yZlNmTUF1ME4wS0hvUUpGOERkM2UzRU1ZeWJsckhGMzJTeTJ5bjR6VjJmUTA3aVdkVGtsdwpYM1luSlNwWUVFMUMzeTY2MUdVV0ZJdkJmbTdvY244VXhnc0NNVzZMa281WjF4eDhlc2ZuUklOaGR2Z0tQbjd3ClBIS0xQRDNEM01RR3AzZXZBVUhZUTFyWmdNdEk0NDdZVXlZZEhaOTQ1S3FaUWdaaUthQnRNZW5aZnFSZ3cwK3IKQ2R4bmoxNEVhK3hGa08yU1hDTCtFa01pR29oVTdPb1F1Z2svNHdJREFRQUJBb0lCQUNBTUJxMm1wbXdDb3dNZApHZTJOYXJOdHdhSnAvTGprWDZaL2xJbjZyQ2NPR1hNUktKTHRObWZpVHVRV0RyRVFQWUVtRWRGN085R0p6Q0JrCjU5Rk52S0d1amxnbDdkc2pMWHRSL0hNV0p0eDEySWRyb2ZvMm1JcC9BalU0cElEbnZIRlZ4c2kwNU43VmdJTTEKZStuQUI0WE5ZWXZLUDBmNHpqQkMwaHVHcFVJTnJTWEF5NEJUL0RQajF2bWkzQVZ5UGUwazNmV1RhY3RxRUN4dwpPUmRRMDhIeCtnRlNzNlpsYldZUU8xWnRlZ1AySlBKUVR3R0k5MGV3Q0JweCtNWC9Fdk5nRDdqbnhFS0ZRYUIzCko3RkpVVFIrcU5qZEs4c2wxeDhBUSs4R3lxVFo3SkNRRHI5WlRUamxBeW1UY0xNcHFSQXB3Z2hoYVZMNXlCejQKanBNODdIRUNnWUVBeENJdUVGMktFbFoxN29leE5RTUpyMGtBM2prbVVoMlpZbG5sOVdmdU1HNU1paldLNzZxNgpUWnVpVjB1c0dDandldDgvd1lvdEdHOFNqSFBsM2VXR0RzMTFJb1doRmdrK29keTNIbU5Gbm9wVWVSbmVxVnNvCnJLT0I0VGpuVjJpdkl2M1FuclFpV2NETjNGd0JoczlYNlQ0ZUwrSUZrVkk4LzdJUDRlbVlmNHNDZ1lFQTNFK3UKSkxZVHZKYm9YU2l5cHJlVVlnZGs1UjdlNVMvK3FNczJPOTh4d3hQRVgvbmxER0FWZlBXMEJ4ODhTMjk2c1dtTQpqYy8xdW95cDJrTWRBZm9DdVVYSDZncGFZVU1qSlpwQ01Vd1dyUDVuTGFQejhJMjZMQzBtY3M0T3JJNjI3MnFXCm5wQ3d1T1VMbzYyYVZrVWJDVGlmMkg0NkNHNkhUY1JGaHB0ZXpBa0NnWUVBaWtkQ3pMejJCRm02eVpJWFNMVzgKbFQxV0JGYXNnc1psaHFhMDd5RDRHR01iU1hIWVk0S3QyTnQ2U0N1TXlIZk1uQVJiMGRyV1VseTA2aHNvSEJxZgpPajUyY0FGZ2djWEF4Nk54NDFYQUZyZVdPTThaWWJOb2FOYmFVZXlwaGNIRGdGc01RMmZpcy82djVNVmxPaU5pCjZvbW1CTUpJaEoxRGJrNmV6ZnJBVG1NQ2dZRUFnSGl5bTFQV0ZJNkh0L09Jb25IQlJKejlPQ01WWmQ3a0NQaGYKaXZCdnEwdDJvMlV0TFZkR2tKVVRRMmZ5bUNiTkRISDVkYVVFcmFGalZ4VDE4SFlqYW5rSHlESDdYR1p6TTNWTwpEa05Kb2QzRXV6ZTFnOXlSNlRyM0JkR2xldmpLTXJrY1ZpRVgvT29NTEltS3k2NEd3d3pUSWNNU0FtSzU0aDZIClVLUi8xa2tDZ1lCQTN2R1lDTlJDS2hIazdKLzZBVnpESVN2VVMvQk9ma0pMTUplcDZ2cWt2SHl1bU9ISkVVcjAKa25KNVJHY3NqY3VsUE1EM2F1TjJlWWovV1k3dkJIclBiSk9sRkFlUVpCc2dKTEg5ZXlzV29tY1haNzRNQ0tUegpUTXhMWDhhZG9Sa3Y0NnhCdlB0YzR2WWVJUWErVWFxRDhVTDY3S3NaWnJVekdDdVRNdnIwWEE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=

`)
