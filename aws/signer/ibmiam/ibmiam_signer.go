package ibmiam

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/aws/smithy-go/logging"
)

// IBM COS SDK Code -- START

type IBMCOSSigner struct {
	logger logging.Logger
}

// NewIBMCOSSigner creates a new IBM COS signer
func NewIBMCOSSigner(options ...func(*IBMCOSSigner)) *IBMCOSSigner {
	signer := &IBMCOSSigner{}
	for _, option := range options {
		option(signer)
	}
	return signer
}

// WithLogger sets the logger for the signer
func WithLogger(logger logging.Logger) func(*IBMCOSSigner) {
	return func(s *IBMCOSSigner) {
		s.logger = logger
	}
}

// SignHTTP implements the aws.HTTPSigner interface
// This is the main method that gets called by the signing middleware
func (s *IBMCOSSigner) SignHTTP(
	ctx context.Context,
	credentials aws.Credentials,
	r *http.Request,
	payloadHash string,
	service string,
	region string,
	signingTime time.Time,
) error {
	token := credentials.Token.AccessToken
	tokenType := credentials.Token.TokenType
	serviceInstanceID := r.Header.Get("ibm-service-instance-id")

	if serviceInstanceID == "" && credentials.ServiceInstanceID != "" {
		// Log the Service Instance ID
		//if s.logger != nil {
		//	s.logger.Logf(logging.Debug, "Setting Service Instance ID: %s", credentials.ServiceInstanceID)
		//}
		r.Header.Set("ibm-service-instance-id", credentials.ServiceInstanceID)
	}

	if token == "" {
		return fmt.Errorf("no bearer token found in credentials")
	}

	// Remove any existing AWS authorization headers
	r.Header.Del("Authorization")
	r.Header.Del("X-Amz-Date")
	r.Header.Del("X-Amz-Security-Token")
	r.Header.Del("X-Amz-Content-Sha256")

	// Add IBM COS bearer token authorization
	authString := tokenType + " " + token
	// r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.Header.Set("Authorization", authString)

	// Log the signing operation if logger is available
	//if s.logger != nil {
	//	s.logger.Logf(logging.Debug, "IBM COS signing request for service %s in region %s", service, region)
	//}

	return nil
}

// IBM COS SDK Code -- END
