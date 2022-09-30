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

			ts, url, err := testServer()
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

func testServer() (*httptest.Server, *url.URL, error) {
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
MIIC3DCCAcQCCQDKSKIAwGGsezANBgkqhkiG9w0BAQUFADAwMQswCQYDVQQGEwJV
UzENMAsGA1UECAwEVXRhaDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIyMDkzMDA1
MDk0NVoXDTMyMDkyNzA1MDk0NVowMDELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFV0
YWgxEjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAJet5JFSOFRh4McYs/V7ZWE4NfGwOYrjwExocwuxf+3rZp2LsAHefN01
rps8fDBk57PbolC5WAutGZBS06asT6/j7XGi1SSIr+C1Sr5X5lnrQWlqimYyVK+k
cPqRkgEVVmYdgESIi0UV1ulEIqfeqgo49S/2u46lt1S/Cvb3dV9oX+aP/CBihPal
z00QtqPdgM3ebG0K/V+JKF5VGkduHfCwIR710pbSvrscPhuQBW+FtGkVGgGsT53w
+m+bpUo8w6FIqp6oQ1gqXIZTDWNtqF7RmuzgohSuo0xfuqkazWMKOsucKJirS0Z2
6wbFG1O/e/GrQ/T1Yp3u8dvSG0KPZy8CAwEAATANBgkqhkiG9w0BAQUFAAOCAQEA
JJFduoko6CTlcgJo4bUJUb/nBDvGZ4e52YJlbsKnqT9bWCgEtaiekw08PBFwjIWK
GXNVUHVhCyFk+hFwCx9TkFfpDNTeiv/xoRBi24Nl60x0kv9SfOyaPaeC3g9cN4HU
JEg72P4A47Owj94RVkqZmwmRcZQ/fh8qTuvSmgoJaMfqXLRGFJbWyPUa3wYzHyjY
CGzMRQTnwJ8Ky4xHoVClbcBTXXTm2tdmojzJP1hwt1zBraq/3tBRBYrKvV4Eqsg4
j3OcbtBxfcVpm/tHlS1JkPznTryVNZhoxf/a4LXSwBGsAHTb14FfNbguuoyl0vZZ
qE56RZINYB2h11pH/1sC0g==
-----END CERTIFICATE-----`)
var serverKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCXreSRUjhUYeDH
GLP1e2VhODXxsDmK48BMaHMLsX/t62adi7AB3nzdNa6bPHwwZOez26JQuVgLrRmQ
UtOmrE+v4+1xotUkiK/gtUq+V+ZZ60FpaopmMlSvpHD6kZIBFVZmHYBEiItFFdbp
RCKn3qoKOPUv9ruOpbdUvwr293VfaF/mj/wgYoT2pc9NELaj3YDN3mxtCv1fiShe
VRpHbh3wsCEe9dKW0r67HD4bkAVvhbRpFRoBrE+d8Ppvm6VKPMOhSKqeqENYKlyG
Uw1jbahe0Zrs4KIUrqNMX7qpGs1jCjrLnCiYq0tGdusGxRtTv3vxq0P09WKd7vHb
0htCj2cvAgMBAAECggEADv66/P7i4Ly4axZvHBKx6BWVh6pDVg7EAQnGbd6DZjMC
dwrLQLQNJhVbiK9HG8Wt/mL1PgPEx4q6X0FA+VZJnnrrC3PsnGsC8DUcCYtJE5Sl
Z9WHjyjkpGSeYrcndwH0A65g8uWI1zCciX0Z6/ygVNhirPY4fpa1dCRa4iV+rgrN
jyZ7+drI3yiiQYds1mF9qBlvDEvomoWlPQ74BMzjQs7BlgRTATwdIztRCiKNUK4G
4PvnNGjDVZsGUeSAaR/+FOE/mCElDE8QR+5eD9iLYUyO0aTXKoMKFZCJMjzqaUlW
XjOE3d6jNzP4qM1vZVc0ozloXsRVuakJduyh/tmL6QKBgQDJ2gmOgTMfj1jKlIOo
T2fa4iAW9PMMErn/x4wTgF+BItuwf8rCVAZYtzMt022CajwHODDpIAiQQZsOOCyE
nd0poYH+IU0PVNvUJRcH7SYSvVlzD9nkM64BLnaX1Fnf4We83GmUhweVqUXTszDn
U0bJcRAxM9kovABkkrHeDQr4mwKBgQDAXlBUtzfSslqsXXLEDfTqaW2u+Z3QaKfQ
VI3z3D1GKmgKzr6Z+oaa6JboKLL9if757GK9xRLTnay775F9x2xlNzz+iq6lZSyb
n1yUaLP1LXa6hmOunP5lo6KduSFP4kqreIHWWuoyKbrB3tEBMqr1Lp2ydL6kP/Dg
0Oz+rnuC/QKBgA55BrRkCSFbKtejnGkGAIFOM1TSDVcxRIrVaPLBApgEwtG95/DV
C3ty70V64mA2c8VkvwUIGfUV7yMu3epIU2I3xVVOV/Mgd36XhjY4R8GSOAaq/UmC
dxh4l2I9hJAr3j9JYnyWzfFqKKqML5Z2fx3UcH/Gouxrxm9voTc1ojK/AoGATPWu
d6XxLFb0VZ7xKiRXRmy1V9o/W8By2rLpM5V54hdXFnPN5zZGIbVJokmeCjbqDjyW
6ErulECxeWKHt2VQJVIrEb6TzlGivgPMewdEb6Mnq8nWGWZvlGQZy7Xj8NyceOs2
Lnai2Ty+nY8x2KPXp01mA54XIwj9qkOLfPx7J1UCgYA7zJ3auVehaKIphx/gLpzL
mWdrQHrrvlS04jy6IfQfKcRo9lGFXgWiPSQomWKbvA2WJik0EP9CQO28oAYKZWP/
jhckhSOsc4+cMSi5b3OqlNiFiL164COTy8I5OLLG1nhIqWAUVYOjQDLlWRCG69xS
VSuwY/kUqjW6vpvWP5j1kg==
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
