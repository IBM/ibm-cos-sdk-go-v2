package ibmiam

import (
	"reflect"
)

const (
	// Default IBM IAM Authentication Server Endpoint
	defaultAuthEndPoint = `https://iam.cloud.ibm.com/identity/token`
	// Debug Log constant
	debugLog                = "DEBUG"
	ibmIamProviderLog       = "IBM IAM PROVIDER"
	ProviderTypeOauth       = "oauth"
	ResourceComputeResource = "CR"
	profilePrefix           = "profile "
	defaultProfile          = "default"
)

type ProviderEnum struct {
	StaticProviderName            string
	TrustedProfileProviderName    string
	EnvProviderTrustedProfileName string
	SharedConfigProviderName      string
	SharedCredentialsProviderName string
	EnvProviderName               string
}

// IBMProvider -> enum instance with values
var IBMProvider = ProviderEnum{
	StaticProviderName:            "StaticProviderIBM",
	TrustedProfileProviderName:    "TrustedProfileProviderIBM",
	EnvProviderTrustedProfileName: "EnvProviderTrustedProfileIBM",
	SharedConfigProviderName:      "SharedConfigProviderIBM",
	SharedCredentialsProviderName: "SharedCredentialsProviderIBM",
	EnvProviderName:               "EnvProviderIBM",
}

func (p ProviderEnum) IsValid(value string) bool {
	val := reflect.ValueOf(p)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == value {
			return true
		}
	}
	return false
}
