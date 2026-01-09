package ibmiam

import (
	"context"
	"os"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam/token"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
)

// Provider Struct
type Provider struct {
	// Name of Provider
	providerName string

	// Type of Provider - SharedCred, SharedConfig, etc.
	providerType string

	// Authenticator instance will be assigned dynamically
	authenticator core.Authenticator

	// Service Instance ID passes in a provider
	serviceInstanceID string

	// Error
	ErrorStatus error

	//Logger attributes
	logger logging.Logger
}

func NewProvider(providerName string, apiKey string, authEndPoint string, serviceInstanceID string) (provider Provider) { //linter complain about (provider *Provider) {
	provider = *new(Provider)
	provider.providerName = providerName
	provider.providerType = ProviderTypeOauth
	provider.logger = logging.NewStandardLogger(os.Stderr)

	if apiKey == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "IbmApiKeyIdNotFound",
			Message: "IBM API Key Id not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] %s error: %v", "<IBM IAM PROVIDER BUILD>", "IBM API Key Id not found", provider.ErrorStatus)
		return
	}

	provider.serviceInstanceID = serviceInstanceID

	if authEndPoint == "" {
		authEndPoint = defaultAuthEndPoint
		provider.logger.Logf(logging.Debug, "[%s] %s: %v", "<IBM IAM PROVIDER BUILD>", "using default auth endpoint", authEndPoint)
	}

	// New code to create a new authenticator using the API Key and auth endpoint
	authenticator, err := core.NewIamAuthenticatorBuilder().
		SetApiKey(apiKey).
		SetURL(authEndPoint).
		SetDisableSSLVerification(true).
		Build()

	if err != nil {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "IbmAuthenticatorError",
			Message: "error creating authenticator",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s], %s error: %v", ibmIamProviderLog, provider.providerName, provider.ErrorStatus)
		return
	}
	provider.authenticator = authenticator
	return provider

	// End of new code
}

func (p Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {

	logger := middleware.GetLogger(ctx)

	if p.ErrorStatus != nil {
		logger.Logf(logging.Debug, "[%s] Provider %s error: %v", ibmIamProviderLog, p.providerName, p.ErrorStatus)
		return aws.Credentials{Source: p.providerName}, p.ErrorStatus
	}
	tokenValue, err := p.authenticator.(*core.IamAuthenticator).GetToken()

	if err != nil {
		logger.Logf(logging.Warn, "Token retrieval failed for provider %s: %v", p.providerName, err)
		var returnErr error
		returnErr = &smithy.GenericAPIError{
			Code:    "TokenManagerRetrieveError",
			Message: "error retrieving the token",
			Fault:   smithy.FaultClient,
		}
		return aws.Credentials{}, returnErr
	}

	return aws.Credentials{
		Token: token.Token{
			AccessToken: tokenValue,
			TokenType:   "Bearer",
		},
		Source:            p.providerName,
		ServiceInstanceID: p.serviceInstanceID,
		SessionToken:      tokenValue,
	}, nil
}

// IsValid ...
// Returns: bool
//
//	Provider validation - boolean
func (p Provider) IsValid() bool {
	return nil == p.ErrorStatus
}

// IsExpired ...
// Returns: bool
//
//	Provider expired or not - boolean
func (p Provider) IsExpired() bool {
	return true
}
