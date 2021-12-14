package apiserver

import (
	"context"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

type ServiceAcountTokenFetcher struct {
	client *kubernetes.Clientset
}

type ServiceAccountTokenFetcherParams struct {
	ServiceAccount    string
	Namespace         string
	ExpirationSeconds int64
	Audiences         []string
}

func NewServiceAccountTokenFetcher() (*ServiceAcountTokenFetcher, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ServiceAcountTokenFetcher{
		client: clientset,
	}, nil
}

func (saTokenFetcher *ServiceAcountTokenFetcher) FetchToken(params *ServiceAccountTokenFetcherParams) ([]byte, error) {
	tokenRequest := &authv1.TokenRequest{
		Spec: authv1.TokenRequestSpec{
			Audiences:         params.Audiences,
			ExpirationSeconds: &params.ExpirationSeconds,
		},
	}
	tokenRequest, err := saTokenFetcher.client.CoreV1().ServiceAccounts(params.Namespace).CreateToken(context.TODO(), params.ServiceAccount, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return []byte(tokenRequest.Status.Token), nil
}
