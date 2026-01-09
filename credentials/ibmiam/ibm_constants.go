package ibmiam

import (
	"reflect"
)

type ProviderEnum struct {
	StaticProviderName         string
	TrustedProfileProviderName string
	IBMIAMProviderLog          string
	SharedConfProviderName     string
}

const (
	// Default IBM IAM Authentication Server Endpoint
	defaultAuthEndPoint = `https://iam.cloud.ibm.com/identity/token`
	// Debug Log constant
	debugLog = "DEBUG"
	// IBM IAM Provider Log constant
	ibmIamProviderLog       = "IBM IAM PROVIDER"
	ProviderTypeOauth       = "oauth"
	ResourceComputeResource = "CR"
	profilePrefix           = "profile "
)

// IBMProvider -> enum instance with values
var IBMProvider = ProviderEnum{
	StaticProviderName:         "StaticProviderIBM",
	TrustedProfileProviderName: "TrustedProfileProviderIBM",
	IBMIAMProviderLog:          "IBM IAM PROVIDER", // New enum - only add here!
	SharedConfProviderName:     "SharedConfigProviderIBM",
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
