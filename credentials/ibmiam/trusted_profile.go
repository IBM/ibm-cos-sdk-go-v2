package ibmiam

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam/token"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	"golang.org/x/net/context"
	"os"
)

// Provider Struct
type TrustedProfileProvider struct {
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

func NewTrustedProfileProvider(providerName string, authEndPoint string, trustedProfileID string, crTokenFilePath string, serviceInstanceID string, resourceType string) (provider TrustedProfileProvider) {

	provider = *new(TrustedProfileProvider)
	provider.providerName = providerName
	provider.providerType = ProviderTypeOauth
	provider.logger = logging.NewStandardLogger(os.Stderr)

	if authEndPoint == "" {
		authEndPoint = defaultAuthEndPoint
		provider.logger.Logf(logging.Debug, "[%s] %s error: %v", ibmIamProviderLog, "using default auth endpoint", authEndPoint)
	}

	if trustedProfileID == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "trustedProfileIDNotFound",
			Message: "Trusted Profile ID not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", ibmIamProviderLog, provider.ErrorStatus)
		return
	}

	if crTokenFilePath == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "crTokenFilePathNotFound",
			Message: "CR Token file path not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", ibmIamProviderLog, provider.ErrorStatus)
		return
	}

	provider.serviceInstanceID = serviceInstanceID

	authenticator, err := core.NewContainerAuthenticatorBuilder().
		SetCRTokenFilename(crTokenFilePath).
		SetIAMProfileID(trustedProfileID).
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
}

func (p TrustedProfileProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {

	// SDK's middleware logger from context
	logger := middleware.GetLogger(ctx)

	if p.ErrorStatus != nil {
		logger.Logf(logging.Debug, "Provider %s error: %v", p.providerName, p.ErrorStatus)
		return aws.Credentials{Source: p.providerName}, p.ErrorStatus
	}

	tokenValue, err := p.authenticator.(*core.ContainerAuthenticator).GetToken()

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
