package eks

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/iamauth/mock_iamauth"
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
					ThumbprintList: aws.StringSlice([]string{"c7a33e1de97f8bf5413ef9a833da98507c95416c"}),
					Url:            &url,
				}).Return(&iam.CreateOpenIDConnectProviderOutput{
					OpenIDConnectProviderArn: aws.String("arn::oidc"),
				}, nil)
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
					ThumbprintList: aws.StringSlice([]string{"c7a33e1de97f8bf5413ef9a833da98507c95416c"}),
					Url:            &url,
				}, nil)
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

			ts, url, err := testServer(serverCert)
			g.Expect(err).To(Succeed())
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
			tc.expect(iamMock.EXPECT(), url.String())
			s := NewService(scope)
			s.IAMClient = iamMock

			cluster := tc.cluster(url.String())
			err = s.reconcileOIDCProvider(&cluster)
			// We reached the trusted policy reconcile which will fail because it tries to connect to the server.
			// But at this point, we already know that the critical area has been covered.
			g.Expect(err).To(MatchError(ContainSubstring("dial tcp: lookup test-cluster-api.nodomain.example.com: no such host")))
		})
	}
}

func testServer(serverCert []byte) (*httptest.Server, *url.URL, error) {
	rootCAs := x509.NewCertPool()

	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init x509 cert/key pair: %w", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	tlsServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	tlsServer.TLS = tlsConfig
	tlsServer.StartTLS()

	serverURL, err := url.Parse(tlsServer.URL)
	if err != nil {
		tlsServer.Close()
		return nil, nil, fmt.Errorf("failed to parse the testserver URL: %w", err)
	}
	serverURL.Host = net.JoinHostPort("localhost", serverURL.Port())

	return tlsServer, serverURL, nil
}

// generated with `mkcert example.com "*.example.com" example.test localhost 127.0.0.1 ::1`.
var serverCert = []byte(`-----BEGIN CERTIFICATE-----
MIIEkDCCAvigAwIBAgIQKmk0KCrPT4AkOBJ2bIML1TANBgkqhkiG9w0BAQsFADCB
iTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMS8wLQYDVQQLDCZza2Fy
bHNvQEdCcy1NYWNCb29rLVByby5jaGVsbG8uaHUgKEdCKTE2MDQGA1UEAwwtbWtj
ZXJ0IHNrYXJsc29AR0JzLU1hY0Jvb2stUHJvLmNoZWxsby5odSAoR0IpMB4XDTIy
MDkyOTE5NTAzOVoXDTI0MTIyOTIwNTAzOVowWjEnMCUGA1UEChMebWtjZXJ0IGRl
dmVsb3BtZW50IGNlcnRpZmljYXRlMS8wLQYDVQQLDCZza2FybHNvQEdCcy1NYWNC
b29rLVByby5jaGVsbG8uaHUgKEdCKTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAKo3e8iCx2fFzcohYZgi9Z7OMs31jlqG4E9lxAolQ4hi6ehPoAgz8NB1
X5nwqjN63E9L8nsGIASS6+PkQGS06HQgV9F5EoXl7fdrAwCMkOLj7fwrTUHIrYNV
XCCPdeo6d9yWAREjjCOaE31taFx8/VBkH11/F7nSsNbkGyoFN9Ob0Z5dPxoDp/it
YzyKpOLEXRrfKA44+gCnE2wczri1qdMVa7gI9EnwZKa6aiuOKwR20X6gdBXW7eSR
UaU67NXl5bxzc1BiRdXI/k0tAt0E7k0xz1qcqVLeLFriHDt30y5yV627B1qis1KI
EJVgZ2LrB/9h9JMOlBuyUcMoBCC56P0CAwEAAaOBoTCBnjAOBgNVHQ8BAf8EBAMC
BaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwHwYDVR0jBBgwFoAUe+QKVJzlMSlp9SBA
8ojl54B4VB8wVgYDVR0RBE8wTYILZXhhbXBsZS5jb22CDSouZXhhbXBsZS5jb22C
DGV4YW1wbGUudGVzdIIJbG9jYWxob3N0hwR/AAABhxAAAAAAAAAAAAAAAAAAAAAB
MA0GCSqGSIb3DQEBCwUAA4IBgQAYnY9zDQb6vQwkO7e5n4ZMuH/4GXpwO+CBEfTH
sMk9fM/T88vcNo085scX+5LMUI+jLFWV0kBLwasLDMVUMQVPwDsVcTgGFc3rTfjd
D3LwEaGl67+jo+a9CWfcbhC3KYJFKgofceEVI6D+lqTsmVLs6wCkiQal683KZc06
5qHXpRDwbNDVh8Yj9RSXBAQTPeI4dKReS3xMr8bSzHKd/M+UCuzT3taLSqbB5RRZ
a35SfMCqz2DDfrQMb5uMTaV7MyyCoYUqJ9JDlGKq4JuayUpBpsw2fyLUcp3rfUzI
I7PTHEEvVq3JOMcM1/wngAJ9sRYng92Fdo0SbN3fT+HVTwef45DKQ7dWwkbjU3EH
1CYm+MCwt2kt2SrmhwyJWDIS0mqQ77+rGIfwPnPxaBzm6LSoiL0W+K0WxaYJ9u9g
uTiE4SmmyQbCOMCgsnfmCMdb1cTpt9VfiUyA8tTEx6/HHRitAdPVHE7/HyXxbkXc
E0WWYq8B9rsnoCV5l8LC5bf75aI=
-----END CERTIFICATE-----
`)
var serverKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCqN3vIgsdnxc3K
IWGYIvWezjLN9Y5ahuBPZcQKJUOIYunoT6AIM/DQdV+Z8KozetxPS/J7BiAEkuvj
5EBktOh0IFfReRKF5e33awMAjJDi4+38K01ByK2DVVwgj3XqOnfclgERI4wjmhN9
bWhcfP1QZB9dfxe50rDW5BsqBTfTm9GeXT8aA6f4rWM8iqTixF0a3ygOOPoApxNs
HM64tanTFWu4CPRJ8GSmumorjisEdtF+oHQV1u3kkVGlOuzV5eW8c3NQYkXVyP5N
LQLdBO5NMc9anKlS3ixa4hw7d9MucletuwdaorNSiBCVYGdi6wf/YfSTDpQbslHD
KAQguej9AgMBAAECggEAcvhq6XI8EcGvUDnf24yxboREGI0li8vSMo7ALUIiLTry
VvGBXEkI/SRqYGYH8CGqMP1RWPs4IR7Dqff/7oWrBNTbvCcU9c/qPdXP/0zyh+4A
TvVIh4huemjrgyqjMIQxdZL0QYlLHLjzNSLy/JWH3AbkkJhJhD/lJiCIoIdJv+lD
W238dIrAX9Ed0BC0p7Ebf5r7asGY1cedGazqSj/T86R4YP1cFVKpGyodOY4w0JLS
Nzoh6k40s+E4Ywcy+E+LyePslMxE9L+ZW4FG4AvXQC1gu4K/wtYSiWsy57e/7+m+
xj6hZPRsXkNnZGHWaIjR5uevAoZKjrp7GMk9u8EbGQKBgQDP6I/mq1TjIdwmaJRz
SRb0O6HPay0es7Z2N/3EcEbosN30u1XSbs+bB44Y6s0xRsxj8sFc3wGQpOP6tGm7
rvafGnApjgqv15776bKWjJx4CE6yAKFjDYOM8i19HVtyUAF+RVU5kW0m4TzUT8wl
3ywtM+ghLWC7b4lUIFti+b2V1wKBgQDRlvs67F5C3R+9vZzCBD2IPHBw28kQFPG1
Mk4jj3kX+kCOTuRyAefyl1LycdtSXr34e277bde9uVJhLPEgf87mn4QisfKT3sru
vg0kJLesS7sZ4F0JytrUJ6hNKgJJWD+zEKtvl17VSLVk6uErDOqYrsJBW62b2asF
jTbjKvK1SwKBgQC1WjjbjquXDBwKbMLA5Qpes/1q/iP3We9Yo3J5/S39Hvoc1aQA
0KPKqQZr+bROvWDf9gpwxh2JXCt4rhJkojOBiQA5Xys3Qy/ssWcUJ0b89NIgNqiP
zGPpd/3x2r+/sMX8rOGwO4gol+QFli2PA2J3c4WSGxD7rkjt1uOgLBQRNQKBgQDI
kW+CB8h8vBcwAFAO6vfnc882cV2L4j8cYzObnCUJ6RX2GVFMOL66zE04bfSwcrHh
JF4khg07JinLjLKDo0tgL67HdPrqvv38UitJN0n9u8slDCx8vn+DHyBUF6twfN8Y
gQ9ODtFV0eqk1JD+HbIywqpq2UzeJAMhoO2xntv82QKBgQCY40cf/os4CY5Vx932
1T04bzbva5eS/3zUBcO7NVVcL7f+y/BgkWImuTtWNxQqlyQRVM78yiRaPBcGYoPP
GWLtipfCWovkqbksXvEXR1gpISJCsmzgcwrsojvlyu8Zbb+kHI46j81QH9YGWPGI
6p63faEx5JFQB6VS0ShqXlnZHg==
-----END PRIVATE KEY-----`)

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
