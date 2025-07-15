package config

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/internal/awstesting"
	"github.com/aws/smithy-go/middleware"
)

func swapECSContainerURI(path string) func() {
	o := ecsContainerEndpoint
	ecsContainerEndpoint = path
	return func() {
		ecsContainerEndpoint = o
	}
}

const ecsFullPathResponse = `{
  "Code": "Success",
  "Type": "AWS-HMAC",
  "AccessKeyId": "ecs-full-path-access-key",
  "SecretAccessKey": "ecs-full-path-ecs-secret-key",
  "Token": "ecs-full-path-token",
  "Expiration": "2100-01-01T00:00:00Z",
  "LastUpdated": "2009-11-23T00:00:00Z"
}`

const assumeRoleRespEcsFullPathMsg = `
<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
    <AssumeRoleResult>
        <AssumedRoleUser>
            <Arn>arn:aws:sts::account_id:assumed-role/role/session_name</Arn>
            <AssumedRoleId>AKID:session_name</AssumedRoleId>
        </AssumedRoleUser>
        <Credentials>
            <AccessKeyId>AKID-Full-Path</AccessKeyId>
            <SecretAccessKey>SECRET-Full-Path</SecretAccessKey>
            <SessionToken>SESSION_TOKEN-Full-Path</SessionToken>
            <Expiration>%s</Expiration>
        </Credentials>
    </AssumeRoleResult>
    <ResponseMetadata>
        <RequestId>request-id</RequestId>
    </ResponseMetadata>
</AssumeRoleResponse>
`

var ecsMetadataServerURL string

func TestSharedConfigCredentialSource(t *testing.T) {
	var configFileForWindows = filepath.Join("testdata", "config_source_shared_for_windows")
	var configFile = filepath.Join("testdata", "config_source_shared")

	var credFileForWindows = filepath.Join("testdata", "credentials_source_shared_for_windows")
	var credFile = filepath.Join("testdata", "credentials_source_shared")

	cases := map[string]struct {
		name                 string
		envProfile           string
		configProfile        string
		expectedError        string
		expectedAccessKey    string
		expectedSecretKey    string
		expectedSessionToken string
		expectedChain        []string
		init                 func() (func(), error)
		dependentOnOS        bool
	}{
		"credential source and source profile": {
			envProfile:    "invalid_source_and_credential_source",
			expectedError: "only one credential type may be specified per profile",
			init: func() (func(), error) {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
				return func() {}, nil
			},
		},
		"env var credential source": {
			configProfile:        "env_var_credential_source",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_env",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
				return func() {}, nil
			},
		},
		"ec2metadata credential source": {
			envProfile: "ec2metadata",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ec2",
			},
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
		},
		"ecs container credential source": {
			envProfile:           "ecscontainer",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ecs",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"chained assume role with env creds": {
			envProfile:           "chained_assume_role",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_chain",
				"assume_role_w_creds_role_arn_ec2",
			},
		},
		"credential process with no ARN set": {
			envProfile:        "cred_proc_no_arn_set",
			dependentOnOS:     true,
			expectedAccessKey: "cred_proc_akid",
			expectedSecretKey: "cred_proc_secret",
		},
		"credential process with ARN set": {
			envProfile:           "cred_proc_arn_set",
			dependentOnOS:        true,
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_proc_role_arn",
			},
		},
		"chained assume role with credential process": {
			envProfile:           "chained_cred_proc",
			dependentOnOS:        true,
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_proc_source_prof",
			},
		},
		"credential source overrides config source": {
			envProfile:           "credentials_overide",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ec2",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"only credential source": {
			envProfile:           "only_credentials_source",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ecs",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"web identity": {
			envProfile:           "webident",
			expectedAccessKey:    "WEB_IDENTITY_AKID",
			expectedSecretKey:    "WEB_IDENTITY_SECRET",
			expectedSessionToken: "WEB_IDENTITY_SESSION_TOKEN",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			restoreEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(restoreEnv)

			if c.dependentOnOS && runtime.GOOS == "windows" {
				os.Setenv("AWS_CONFIG_FILE", configFileForWindows)
				os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credFileForWindows)
			} else {
				os.Setenv("AWS_CONFIG_FILE", configFile)
				os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credFile)
			}

			os.Setenv("AWS_REGION", "us-east-1")
			if len(c.envProfile) != 0 {
				os.Setenv("AWS_PROFILE", c.envProfile)
			}

			var cleanup func()
			if c.init != nil {
				var err error
				cleanup, err = c.init()
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
				defer cleanup()
			}

			var credChain []string

			loadOptions := []func(*LoadOptions) error{
				WithAPIOptions([]func(*middleware.Stack) error{
					func(stack *middleware.Stack) error {
						return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("GetRoleArns",
							func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
							) (
								out middleware.InitializeOutput, metadata middleware.Metadata, err error,
							) {
								return next.HandleInitialize(ctx, in)
							}), middleware.After)
					},
				}),
			}

			if len(c.configProfile) != 0 {
				loadOptions = append(loadOptions, WithSharedConfigProfile(c.configProfile))
			}

			config, err := LoadDefaultConfig(context.Background(), loadOptions...)
			if err != nil {
				if len(c.expectedError) > 0 {
					if e, a := c.expectedError, err.Error(); !strings.Contains(a, e) {
						t.Fatalf("expect %v, but got %v", e, a)
					}
					return
				}
				t.Fatalf("expect no error, got %v", err)
			} else if len(c.expectedError) > 0 {
				t.Fatalf("expect error, got none")
			}

			creds, err := config.Credentials.Retrieve(context.Background())
			if err != nil {
				t.Fatalf("expected no error, but received %v", err)
			}

			if e, a := c.expectedChain, credChain; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}

			if e, a := c.expectedAccessKey, creds.AccessKeyID; e != a {
				t.Errorf("expected %v, but received %v", e, a)
			}

			if e, a := c.expectedSecretKey, creds.SecretAccessKey; e != a {
				t.Errorf("expect %v, but received %v", e, a)
			}

			if e, a := c.expectedSessionToken, creds.SessionToken; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestResolveCredentialsCacheOptions(t *testing.T) {
	var cfg aws.Config
	var optionsFnCalled bool

	err := resolveCredentials(context.Background(), &cfg, configs{LoadOptions{
		CredentialsCacheOptions: func(o *aws.CredentialsCacheOptions) {
			optionsFnCalled = true
			o.ExpiryWindow = time.Minute * 5
		},
	}})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if !optionsFnCalled {
		t.Errorf("expect options to be called")
	}
}

func TestResolveCredentialsEcsContainer(t *testing.T) {
	testCases := map[string]struct {
		expectedAccessKey string
		expectedSecretKey string
		envVar            map[string]string
		configFile        string
	}{
		"only relative ECS URI set": {
			expectedAccessKey: "ecs-access-key",
			expectedSecretKey: "ecs-secret-key",
			envVar: map[string]string{
				"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI": "/ECS",
			},
		},
		"only full ECS URI set": {
			expectedAccessKey: "ecs-full-path-access-key",
			expectedSecretKey: "ecs-full-path-ecs-secret-key",
			envVar: map[string]string{
				"AWS_CONTAINER_CREDENTIALS_FULL_URI": "placeholder-replaced-at-runtime",
			},
		},
		"relative ECS URI has precedence over full": {
			expectedAccessKey: "ecs-access-key",
			expectedSecretKey: "ecs-secret-key",
			envVar: map[string]string{
				"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI": "/ECS",
				"AWS_CONTAINER_CREDENTIALS_FULL_URI":     "placeholder-replaced-at-runtime",
			},
		},
		"credential source only relative ECS URI set": {
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			envVar: map[string]string{
				"AWS_PROFILE":                            "ecscontainer",
				"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI": "/ECS",
			},
			configFile: filepath.Join("testdata", "config_source_shared"),
		},
		"credential source only full ECS URI set": {
			expectedAccessKey: "AKID-Full-Path",
			expectedSecretKey: "SECRET-Full-Path",
			envVar: map[string]string{
				"AWS_CONTAINER_CREDENTIALS_FULL_URI": "placeholder-replaced-at-runtime",
				"AWS_PROFILE":                        "ecscontainer",
			},
			configFile: filepath.Join("testdata", "config_source_shared"),
		},
		"credential source relative ECS URI has precedence over full": {
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			envVar: map[string]string{
				"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI": "/ECS",
				"AWS_CONTAINER_CREDENTIALS_FULL_URI":     "placeholder-replaced-at-runtime",
				"AWS_PROFILE":                            "ecscontainer",
			},
			configFile: filepath.Join("testdata", "config_source_shared"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			restoreEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(restoreEnv)
			var sharedConfigFiles []string
			if tc.configFile != "" {
				sharedConfigFiles = append(sharedConfigFiles, tc.configFile)
			}
			opts := []func(*LoadOptions) error{
				WithRetryer(func() aws.Retryer { return aws.NopRetryer{} }),
				WithSharedConfigFiles(sharedConfigFiles),
				WithSharedCredentialsFiles([]string{}),
			}
			for k, v := range tc.envVar {
				// since we don't know the value of this until the server starts
				if k == "AWS_CONTAINER_CREDENTIALS_FULL_URI" {
					v = ecsMetadataServerURL + "/ECSFullPath"
				}
				os.Setenv(k, v)
			}
			cfg, err := LoadDefaultConfig(context.TODO(), opts...)
			if err != nil {
				t.Fatalf("could not load config: %s", err)
			}
			actual, err := cfg.Credentials.Retrieve(context.TODO())
			if err != nil {
				t.Fatalf("could not retrieve credentials: %s", err)
			}
			if actual.AccessKeyID != tc.expectedAccessKey {
				t.Errorf("expected access key to be %s, got %s", tc.expectedAccessKey, actual.AccessKeyID)
			}
			if actual.SecretAccessKey != tc.expectedSecretKey {
				t.Errorf("expected secret key to be %s, got %s", tc.expectedSecretKey, actual.SecretAccessKey)
			}
		})
	}

}

type stubErrorClient struct {
	err error
}

func (c stubErrorClient) Do(*http.Request) (*http.Response, error) { return nil, c.err }
