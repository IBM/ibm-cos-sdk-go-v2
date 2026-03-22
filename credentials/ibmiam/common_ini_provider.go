package ibmiam

import (
	"os"

	"github.com/IBM/ibm-cos-sdk-go-v2/internal/ini"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
)

func commonIniProvider(providerName string, config IBMEnvConfig, filename, profileName string) Provider {

	provider := *new(Provider)
	provider.providerName = providerName
	provider.logger = logging.NewStandardLogger(os.Stderr)
	// Opens an ini file with the filename passed in for shared credentials
	// If fails, returns error
	in, err := ini.OpenFile(filename)
	if err != nil {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "SharedCredentialsOpenError",
			Message: "Shared Credentials Open Error",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", "<IBM IAM PROVIDER BUILD>", provider.ErrorStatus)
		return provider
	}

	// Gets section of the shared credentials ini file
	// If fails, returns error
	iniProfile, ok := in.GetSection(profileName)
	if !ok {
		provider.ErrorStatus = &smithy.GenericAPIError{
			Code:    "SharedCredentialsProfileNotFound",
			Message: "Shared Credentials Section '" + profileName + "' not Found in file '" + filename + "'",
			Fault:   smithy.FaultClient,
		}
		provider.logger.Logf(logging.Debug, "[%s] error: %v", "<IBM IAM PROVIDER BUILD>", provider.ErrorStatus)
		return provider
	}

	// Populate the IBM IAM Credential values
	apiKey := iniProfile.String("ibm_api_key_id")
	serviceInstanceID := iniProfile.String("ibm_service_instance_id")
	authEndPoint := iniProfile.String("ibm_auth_endpoint")

	return NewProvider(providerName, apiKey, authEndPoint, serviceInstanceID)
}
