package ibmiam

import (
	"reflect"
)

type ProviderEnum struct {
	StaticProviderName         string
	TrustedProfileProviderName string
	IBMIAMProviderLog          string
}

// Create the enum instance with values
var IBMProvider = ProviderEnum{
	StaticProviderName:         "StaticProviderIBM",
	TrustedProfileProviderName: "TrustedProfileProviderIBM",
	IBMIAMProviderLog:          "IBM IAM PROVIDER", // New enum - only add here!
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
