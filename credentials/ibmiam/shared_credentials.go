package ibmiam

import (
	"os"

	"github.com/IBM/ibm-cos-sdk-go-v2/internal/shareddefaults"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/logging"
)

func NewSharedCredentialsProvider(cfg IBMEnvConfig, filename, profileName string) Provider {

	provider := *new(Provider)
	provider.providerName = IBMProvider.SharedCredentialsProviderName
	provider.logger = logging.NewStandardLogger(os.Stderr)

	// Sets the file name from possible locations
	//	- AWS_SHARED_CREDENTIALS_FILE environment variable
	// Error if the filename is missing
	if filename == "" {
		filename = cfg.SharedCredentialsFile // os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
		if filename == "" {
			// BUG where will we use home?
			home := shareddefaults.UserHomeDir()
			if home == "" {
				provider.ErrorStatus = &smithy.GenericAPIError{
					Code:    "SharedCredentialsHomeNotFound",
					Message: "Shared Credentials Home folder not found",
					Fault:   smithy.FaultClient,
				}
				provider.logger.Logf(logging.Debug, "[%s] error: %v", "<IBM IAM PROVIDER BUILD>", provider.ErrorStatus)
				return provider
				//e := awserr.New("SharedCredentialsHomeNotFound", "Shared Credentials Home folder not found", nil)
				//logFromConfigHelper(config, "<DEBUG>", "<IBM IAM PROVIDER BUILD>", SharedCredsProviderName, e)
				//return &Provider{
				//	providerName: SharedCredsProviderName,
				//	ErrorStatus:  e,
				//}
			}
			filename = shareddefaults.SharedCredentialsFilename()
		}
	}

	// Sets the profile name from AWS_PROFILE environment variable
	// Otherwise sets the profile name with defaultProfile passed in
	if profileName == "" {
		profileName = cfg.Profile
		if profileName == "" {
			profileName = defaultProfile
		}
	}

	return commonIniProvider(provider.providerName, cfg, filename, profileName)
}
