package ibmiam

// "github.com/IBM/ibm-cos-sdk-go-v2/aws"

// NewStaticProvider constructor of the IBM IAM provider that uses IAM details passed directly
// Returns: New Provider (AWS type)
func NewStaticProvider(authEndPoint, apiKey, serviceInstanceID string) Provider {
	return NewProvider(IBMProvider.StaticProviderName, apiKey, authEndPoint, serviceInstanceID)
}

// NewStaticCredentials constructor for IBM IAM that uses IAM credentials passed in
// Returns: Provider which implements the aws credentialsProvider Interface
func NewStaticCredentials(authEndPoint, apiKey, serviceInstanceID string) Provider {
	return NewStaticProvider(authEndPoint, apiKey, serviceInstanceID)
}

// NewStaticTrustedProfileProviderCR -> constructor for IBM Trusted Profile that uses details passed in
// Returns: New TrustedProfileProvider which implements aws credentialProvider Interface
func NewTrustedProfileProviderCR(authEndPoint string, trustedProfileID string, crTokenFilePath string, serviceInstanceID string) TrustedProfileProvider {
	return NewTrustedProfileProvider(IBMProvider.TrustedProfileProviderName, authEndPoint, trustedProfileID, crTokenFilePath, serviceInstanceID, ResourceComputeResource)
}
