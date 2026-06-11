package ibmiam

import (
	"os"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam/token"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	"golang.org/x/net/context"
)

// TrustedProfileProvider Provider Struct
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

// TrustedProfileConfig has all the authentication parameters for trusted profile.
type TrustedProfileConfig struct {
	TrustedProfileID       string
	CrTokenFilePath        string
	DisableSSLVerification bool
}

// NewTrustedProfileProvider allows the creation of a custom IBM IAM Trusted Profile Provider
// Parameters:
//
//	Provider Name
//	IBM IAM Authentication Server Endpoint
//	Trusted Profile ID
//	CR token file path
//	Service Instance ID
//	Resource type
//
// Returns:
//
//	TrustedProfileProvider
//
// Deprecated: Use NewTrustedProfileProviderWithConfig instead.
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

	// Note: SSL verification is enabled by default for security. Only disable for testing.
	authenticator, err := core.NewContainerAuthenticatorBuilder().
		SetCRTokenFilename(crTokenFilePath).
		SetIAMProfileID(trustedProfileID).
		SetURL(authEndPoint).
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

// NewTrustedProfileProviderWithConfig allows the creation of a custom IBM IAM Trusted Profile Provider with trusted profile config
// Parameters:
//
//	Provider Name
//	IBM IAM Authentication Server Endpoint
//	Trusted Profile Config
//	Service Instance ID
//	Resource type
//
// Returns:
//
//	TrustedProfileProvider
func NewTrustedProfileProviderWithConfig(providerName string, authEndPoint string, trustedProfileConfig *TrustedProfileConfig, serviceInstanceID string, resourceType string) (provider TrustedProfileProvider) {

	provider = *new(TrustedProfileProvider)
	provider.providerName = providerName
	provider.providerType = ProviderTypeOauth
	provider.logger = logging.NewStandardLogger(os.Stderr)

	if trustedProfileConfig == nil {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "trustedProfileConfigNil",
			Message: "TrustedProfileConfig cannot be nil",
			Fault:   smithy.FaultClient,
		}
		return
	}

	if authEndPoint == "" {
		authEndPoint = defaultAuthEndPoint
		provider.logger.Logf(logging.Debug, "[%s] %s error: %v", ibmIamProviderLog, "using default auth endpoint", authEndPoint)
	}

	if trustedProfileConfig.TrustedProfileID == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "trustedProfileIDNotFound",
			Message: "Trusted Profile ID not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", ibmIamProviderLog, provider.ErrorStatus)
		return
	}

	if trustedProfileConfig.CrTokenFilePath == "" {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "crTokenFilePathNotFound",
			Message: "CR Token file path not found",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", ibmIamProviderLog, provider.ErrorStatus)
		return
	}

	provider.serviceInstanceID = serviceInstanceID

	// Build the authenticator with configurable SSL verification
	builder := core.NewContainerAuthenticatorBuilder().
		SetCRTokenFilename(trustedProfileConfig.CrTokenFilePath).
		SetIAMProfileID(trustedProfileConfig.TrustedProfileID).
		SetURL(authEndPoint)

	// Check if SSL verification should be disabled
	// Priority: 1. Config field, 2. Environment variable
	disableSSL := trustedProfileConfig.DisableSSLVerification
	if !disableSSL {
		// Check environment variable if not explicitly set in config
		envValue := os.Getenv("TRUSTED_PROFILE_CR_DISABLE_SSL_VERIFICATION")
		if strings.ToLower(envValue) == "true" {
			disableSSL = true
		}
	}

	// Only disable SSL verification if explicitly requested (not recommended for production)
	if disableSSL {
		builder.SetDisableSSLVerification(true)
	}

	authenticator, err := builder.Build()

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

func (p TrustedProfileProvider) IsValid() bool {
	return nil == p.ErrorStatus
}

func (p TrustedProfileProvider) GetProviderName() string {
	return p.providerName
}
