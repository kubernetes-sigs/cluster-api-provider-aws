package rosa

import (
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
)

type rosaClient struct {
	ocm *sdk.Connection
}

// NewRosaClientWithConnection creates a client with a preexisting connection for testing purpose
func NewRosaClientWithConnection(connection *sdk.Connection) *rosaClient {
	return &rosaClient{
		ocm: connection,
	}
}

func NewRosaClient(token string) (*rosaClient, error) {
	ocmAPIUrl := os.Getenv("OCM_API_URL")
	if ocmAPIUrl == "" {
		ocmAPIUrl = "https://api.openshift.com"
	}

	// Create a logger that has the debug level enabled:
	logger, err := sdk.NewGoLoggerBuilder().
		Debug(true).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(token).
		URL(ocmAPIUrl).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ocm connection: %w", err)
	}

	return &rosaClient{
		ocm: connection,
	}, nil
}

func (c *rosaClient) Close() error {
	return c.ocm.Close()
}

func (c *rosaClient) GetConnectionURL() string {
	return c.ocm.URL()
}

func (c *rosaClient) GetConnectionTokens() (string, string, error) {
	return c.ocm.Tokens()
}
