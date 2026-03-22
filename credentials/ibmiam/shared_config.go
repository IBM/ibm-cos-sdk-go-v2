package ibmiam

import (
	"os"

	"github.com/IBM/ibm-cos-sdk-go-v2/internal/shareddefaults"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
)

func NewSharedConfigProvider(cfg IBMEnvConfig, filename string, profileName string) Provider {
	provider := *new(Provider)
	provider.providerName = IBMProvider.SharedConfigProviderName
	provider.logger = logging.NewStandardLogger(os.Stderr)

	if filename == "" {
		filename = cfg.SharedConfigFile
		if filename == "" {
			// BUG?
			home := shareddefaults.UserHomeDir()
			if home == "" {
				provider.ErrorStatus = &smithy.GenericAPIError{
					Code:    "SharedCredentialsHomeNotFound",
					Message: "Shared Credentials Home folder not found",
					Fault:   smithy.FaultClient,
				}
				provider.logger.Logf(logging.Debug, "[%s] error: %v", "<IBM IAM PROVIDER BUILD>", provider.ErrorStatus)
				return provider
			}
			filename = shareddefaults.SharedConfigFilename()
		}
		if profileName == "" {
			profileName = cfg.Profile
			if profileName == "" {
				profileName = defaultProfile
			} else {
				profileName = profilePrefix + profileName
			}
		} else {
			profileName = profilePrefix + profileName
		}
	}
	return commonIniProvider(provider.providerName, cfg, filename, profileName)
}
