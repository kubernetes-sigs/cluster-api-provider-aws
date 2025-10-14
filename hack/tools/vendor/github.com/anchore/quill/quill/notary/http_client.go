package notary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anchore/quill/internal/log"
)

type httpClient struct {
	client *http.Client
	token  string
}

func newHTTPClient(token string, httpTimeout time.Duration) *httpClient {
	if httpTimeout == 0 {
		httpTimeout = time.Second * 30
	}

	return &httpClient{
		client: &http.Client{
			Timeout: httpTimeout,
		},
		token: token,
	}
}

func (s httpClient) get(ctx context.Context, endpoint string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, body)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	return s.do(request)
}

func (s httpClient) post(ctx context.Context, endpoint string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return s.do(request)
}

func (s httpClient) do(request *http.Request) (*http.Response, error) {
	log.Tracef("http %s %s", request.Method, request.URL)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	return s.client.Do(request)
}
