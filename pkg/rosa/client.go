// Package rosa provides a way to interact with the Red Hat OpenShift Service on AWS (ROSA) API.
package rosa

import (
	"context"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	ocmcfg "github.com/openshift/rosa/pkg/config"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

const (
	ocmTokenKey        = "ocmToken"
	ocmAPIURLKey       = "ocmApiUrl"
	ocmClientIdKey     = "ocmClientId"
	ocmClientSecretKey = "ocmClientSecret"
)

// NewOCMClient creates a new OCM client.
func NewOCMClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*ocm.Client, error) {
	token, url, clientId, clientSecret, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}

	ocmConfig := ocmcfg.Config{
		URL: url,
	}

	if clientId != "" && clientSecret != "" {
		ocmConfig.ClientID = clientId
		ocmConfig.ClientSecret = clientSecret
	} else if token != "" {
		ocmConfig.AccessToken = token
	}

	return ocm.NewClient().Logger(logrus.New()).Config(&ocmConfig).Build()
}

func newOCMRawConnection(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*sdk.Connection, error) {
	ocmSdkLogger, err := sdk.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	token, url, clientId, clientSecret, err := ocmCredentials(ctx, rosaScope)
	if err != nil {
		return nil, err
	}

	connBuilder := sdk.NewConnectionBuilder().
		Logger(ocmSdkLogger).
		URL(url)

	if clientId != "" && clientSecret != "" {
		connBuilder.Client(clientId, clientSecret)
	} else if token != "" {
		connBuilder.Tokens(token)
	}

	connection, err := connBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create ocm connection: %w", err)
	}

	return connection, nil
}

func ocmCredentials(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, string, string, string, error) {
	var token string           // Offline SSO token
	var ocmClientId string     // Service account client id
	var ocmClientSecret string // Service account client secret
	var ocmAPIUrl string       // https://api.openshift.com by default
	var secret *corev1.Secret

	secret = rosaScope.CredentialsSecret() // We'll retrieve the OCM credentials from the ROSA control plane
	if secret != nil {
		if err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret); err != nil {
			return "", "", "", "", fmt.Errorf("failed to get credentials secret: %w", err)
		}
	} else { // If the reference to OCM secret wasn't specified in the ROSA control plane, we'll default to a predefined secret name
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default-rosa-creds-secret",
				Namespace: rosaScope.Namespace(),
			},
		}
	}

	if err := rosaScope.Client.Get(ctx, client.ObjectKeyFromObject(secret), secret); err != nil {
		return "", "", "", "", fmt.Errorf("failed to get credentials secret: %w", err)
	}

	token = string(secret.Data[ocmTokenKey])
	ocmAPIUrl = string(secret.Data[ocmAPIURLKey])
	ocmClientId = string(secret.Data[ocmClientIdKey])
	ocmClientSecret = string(secret.Data[ocmClientSecretKey])

	// Deprecation warning in case SSO offline token was used
	if token != "" {
		rosaScope.Info(fmt.Sprintf("Using SSO offline token (%s) is deprecated, use service account credentials instead (%s and %s)",
			ocmTokenKey, ocmClientIdKey, ocmClientSecretKey))
	}

	if token == "" && (ocmClientId == "" || ocmClientSecret == "") {
		// Last fall-back is to use OCM_TOKEN & OCM_API_URL environment variables (soon to be deprecated)
		token = os.Getenv("OCM_TOKEN")
		ocmAPIUrl = os.Getenv("OCM_API_URL")

		if token != "" {
			rosaScope.Info(fmt.Sprintf("Defining OCM credentials in environment variable is deprecated, use secret with service account credentials instead (%s and %s)",
				ocmTokenKey, ocmClientIdKey, ocmClientSecretKey))
		} else {
			return "", "", "", "",
				fmt.Errorf("OCM credentials have not been provided. Make sure to set the secret with service account credentials (%s and %s)",
					ocmClientIdKey, ocmClientSecretKey)
		}
	}

	if ocmAPIUrl == "" {
		ocmAPIUrl = "https://api.openshift.com" // Defaults to production URL
	}

	return token, ocmAPIUrl, ocmClientId, ocmClientSecret, nil
}
