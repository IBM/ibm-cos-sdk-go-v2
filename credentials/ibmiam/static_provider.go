package ibmiam

// "github.com/IBM/ibm-cos-sdk-go-v2/aws"

// NewStaticProvider constructor of the IBM IAM provider that uses IAM details passed directly
// Returns: New Provider (AWS type)
func NewStaticProvider(authEndPoint, apiKey, serviceInstanceID string) Provider {
	return NewProvider(IBMProvider.StaticProviderName, apiKey, authEndPoint, serviceInstanceID)
}

// NewStaticCredentials constructor for IBM IAM that uses IAM credentials passed in
// Returns: credentials.NewCredentials(newStaticProvider()) (AWS type)
func NewStaticCredentials(authEndPoint, apiKey, serviceInstanceID string) Provider {
	return NewStaticProvider(authEndPoint, apiKey, serviceInstanceID)
}
