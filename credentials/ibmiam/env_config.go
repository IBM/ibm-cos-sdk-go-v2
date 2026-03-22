package ibmiam

import (
	"context"
	"os"

	"github.com/aws/smithy-go"
)

// Environment variables that will be read for config
const (
	ibmAuthEndpointEnv      = "IBM_AUTH_ENDPOINT"
	ibmServiceInstanceIdEnv = "IBM_SERVICE_INSTANCE_ID"
	ibmApiKeyEnv            = "IBM_API_KEY_ID"

	// Trusted Profile Env Variables
	ibmTrustedProfileIdEnv        = "TRUSTED_PROFILE_ID"
	ibmCrTokenFilePathEnv         = "CR_TOKEN_FILE_PATH"
	ibmTrustedProfileResourceType = "TRUSTED_PROFILE_RESOURCE_TYPE"

	ibmConfigFileEnv            = "AWS_CONFIG_FILE"
	ibmProfileEnv               = "AWS_PROFILE"
	ibmSharedCredentialsFileEnv = "AWS_SHARED_CREDENTIALS_FILE"
)

type IBMEnvConfig struct {
	AuthEndpoint          string `json:"authEndpoint"`
	ServiceInstanceId     string `json:"serviceInstanceId"`
	ApiKey                string `json:"apiKey"`
	TrustedProfileId      string `json:"trustedProfileId"`
	CrTokenFilePath       string `json:"crTokenFilePath"`
	ResourceType          string `json:"resourceType"`
	SharedConfigFile      string `json:"configFile"`
	Profile               string `json:"profile"`
	SharedCredentialsFile string `json:"sharedCredentialsFile"`
	ConfigSource          string
}

func GetIBMCredentials(ctx context.Context) (ExtendedCredentialsProvider, error) {
	cfg := loadIBMEnvConfig()
	provider, err := resolveEnvConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func loadIBMEnvConfig() IBMEnvConfig {
	return NewEnvConfig()
}

func NewEnvConfig() IBMEnvConfig {
	var cfg IBMEnvConfig

	//IAM
	cfg.ServiceInstanceId = os.Getenv(ibmServiceInstanceIdEnv)
	cfg.AuthEndpoint = os.Getenv(ibmAuthEndpointEnv)
	cfg.ApiKey = os.Getenv(ibmApiKeyEnv)

	// Trusted Profile
	cfg.TrustedProfileId = os.Getenv(ibmTrustedProfileIdEnv)
	cfg.CrTokenFilePath = os.Getenv(ibmCrTokenFilePathEnv)
	cfg.ResourceType = os.Getenv(ibmTrustedProfileResourceType)

	// Shared Config
	cfg.SharedConfigFile = os.Getenv(ibmConfigFileEnv)
	cfg.Profile = os.Getenv(ibmProfileEnv)
	cfg.SharedCredentialsFile = os.Getenv(ibmSharedCredentialsFileEnv)

	return cfg
}

func resolveEnvConfig(ctx context.Context, cfg IBMEnvConfig) (ExtendedCredentialsProvider, error) {
	if cfg.checkForConflictingCredentials() {
		return nil, &smithy.GenericAPIError{
			Code:    "InvalidCredentials",
			Message: "only one of Api Key or Trusted Profile ID should be set, not both",
			Fault:   smithy.FaultClient,
		}
	}

	if provider := NewEnvProvider(cfg); provider.IsValid() {
		return provider, nil
	}

	if provider := NewEnvProviderTrustedProfile(cfg); provider.IsValid() {
		return provider, nil
	}

	if provider := NewSharedCredentials(cfg, "", ""); provider.IsValid() {
		return provider, nil
	}

	if provider := NewConfigCredentials(cfg, "", ""); provider.IsValid() {
		return provider, nil
	}
	return nil, &smithy.GenericAPIError{
		Code:    "InvalidCredentials",
		Message: "no valid credentials found",
		Fault:   smithy.FaultClient,
	}
}

func (cfg IBMEnvConfig) HasCredentials() bool {
	return (cfg.HasIAM() || cfg.HasTrustedProfile()) && !cfg.checkForConflictingCredentials()
}

func (cfg IBMEnvConfig) checkForConflictingCredentials() bool {
	return cfg.ApiKey != "" && cfg.TrustedProfileId != ""
}

func (cfg IBMEnvConfig) HasIAM() bool {
	return cfg.ServiceInstanceId != "" && cfg.AuthEndpoint != "" && cfg.ApiKey != ""
}

func (cfg IBMEnvConfig) HasTrustedProfile() bool {
	return cfg.TrustedProfileId != "" &&
		cfg.CrTokenFilePath != "" &&
		cfg.ResourceType != "" &&
		cfg.ServiceInstanceId == "" &&
		cfg.AuthEndpoint == ""
}
