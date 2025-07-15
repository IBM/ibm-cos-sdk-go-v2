package ibmiam

import (
	"context"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam/token"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
)

const (
	// Constants
	// Default IBM IAM Authentication Server Endpoint
	defaultAuthEndPoint = `https://iam.cloud.ibm.com/identity/token`

	// Logger constants
	// Debug Log constant
	debugLog = "DEBUG"
	// IBM IAM Provider Log constant
	ibmIamProviderLog = "IBM IAM PROVIDER"
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
	provider.providerType = "oauth"

	if apiKey == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "IbmApiKeyIdNotFound",
			Message: "IBM API Key Id not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(debugLog, "<IBM IAM PROVIDER BUILD>", "IBM API Key Id not found", provider.ErrorStatus)
		return
	}

	provider.serviceInstanceID = serviceInstanceID

	if authEndPoint == "" {
		authEndPoint = defaultAuthEndPoint
		provider.logger.Logf(debugLog, "<IBM IAM PROVIDER BUILD>", "using default auth endpoint", authEndPoint)
	}

	// New code to create a new authenticator using the API Key and auth endpoint
	authenticator, err := core.NewIamAuthenticatorBuilder().
		SetApiKey(apiKey).
		SetURL(authEndPoint).
		Build()

	if err != nil {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "IbmAuthenticatorError",
			Message: "error creating authenticator",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(debugLog, ibmIamProviderLog, provider.providerName, provider.ErrorStatus)
		return
	}
	provider.authenticator = authenticator
	return provider

	// End of new code
}

func (p Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {

	if p.ErrorStatus != nil {
		// if p.logLevel.Matches(aws.LogDebug) {
		// 	p.logger.Log(debugLog, ibmIamProviderLog, p.providerName, p.ErrorStatus)
		// }
		return aws.Credentials{Source: p.providerName}, p.ErrorStatus
	}
	tokenValue, err := p.authenticator.(*core.IamAuthenticator).GetToken()
	if err != nil {
		var returnErr error
		returnErr = &smithy.GenericAPIError{
			Code:    "TokenManagerRetrieveError",
			Message: "error retrieving the token",
			Fault:   smithy.FaultClient,
		}
		p.logger.Logf(debugLog, ibmIamProviderLog, p.providerName, returnErr)
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
func (p *Provider) IsValid() bool {
	return nil == p.ErrorStatus
}

// IsExpired ...
// Returns: bool
//
//	Provider expired or not - boolean
func (p *Provider) IsExpired() bool {
	return true
}
