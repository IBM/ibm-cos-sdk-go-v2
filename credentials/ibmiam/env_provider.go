package ibmiam

func NewEnvProvider(cfg IBMEnvConfig) Provider {
	apiKey := cfg.ApiKey
	serviceInstanceId := cfg.ServiceInstanceId
	authEndPoint := cfg.AuthEndpoint
	return NewProvider(IBMProvider.EnvProviderName, apiKey, authEndPoint, serviceInstanceId)
}

func NewEnvProviderTrustedProfile(cfg IBMEnvConfig) TrustedProfileProvider {
	trustedProfileID := cfg.TrustedProfileId
	serviceInstanceId := cfg.ServiceInstanceId
	crTokenFilePath := cfg.CrTokenFilePath
	authEndpoint := cfg.AuthEndpoint

	return NewTrustedProfileProvider(IBMProvider.EnvProviderTrustedProfileName, authEndpoint, trustedProfileID, crTokenFilePath, serviceInstanceId, ResourceComputeResource)
}

func NewConfigCredentials(cfg IBMEnvConfig, filename, profileName string) Provider {
	return NewSharedConfigProvider(cfg, filename, profileName)
}

func NewSharedCredentials(cfg IBMEnvConfig, filename, profileName string) Provider {
	return NewSharedCredentialsProvider(cfg, filename, profileName)
}
