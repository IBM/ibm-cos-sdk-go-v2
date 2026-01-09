package config

import (
	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
)

const execEnvVar = "AWS_EXECUTION_ENV"

// DefaultsModeOptions is the set of options that are used to configure
type DefaultsModeOptions struct {
	// The SDK configuration defaults mode. Defaults to legacy if not specified.
	//
	// Supported modes are: auto, cross-region, in-region, legacy, mobile, standard
	Mode aws.DefaultsMode

	// The EC2 Instance Metadata Client that should be used when performing environment
	// discovery when aws.DefaultsModeAuto is set.
	//
	// If not specified the SDK will construct a client if the instance metadata service has not been disabled by
	// the AWS_EC2_METADATA_DISABLED environment variable.
	//IMDSClient *imds.Client
}
