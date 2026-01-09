package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/endpointcreds"
)

const (
	// valid credential source values
	credSourceEnvironment      = "Environment"
	httpProviderAuthFileEnvVar = "AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE"
)

// resolveCredentials extracts a credential provider from slice of config
// sources.
//
// If an explicit credential provider is not found the resolver will fallback
// to resolving credentials by extracting a credential provider from EnvConfig
// and SharedConfig.
func resolveCredentials(ctx context.Context, cfg *aws.Config, configs configs) error {
	found, err := resolveCredentialProvider(ctx, cfg, configs)
	if found || err != nil {
		return err
	}

	return resolveCredentialChain(ctx, cfg, configs)
}

// resolveCredentialProvider extracts the first instance of Credentials from the
// config slices.
//
// The resolved CredentialProvider will be wrapped in a cache to ensure the
// credentials are only refreshed when needed. This also protects the
// credential provider to be used concurrently.
//
// Config providers used:
// * credentialsProviderProvider
func resolveCredentialProvider(ctx context.Context, cfg *aws.Config, configs configs) (bool, error) {
	credProvider, found, err := getCredentialsProvider(ctx, configs)
	if !found || err != nil {
		return false, err
	}

	cfg.Credentials, err = wrapWithCredentialsCache(ctx, configs, credProvider)
	if err != nil {
		return false, err
	}

	return true, nil
}

// resolveCredentialChain resolves a credential provider chain using EnvConfig
// and SharedConfig if present in the slice of provided configs.
//
// The resolved CredentialProvider will be wrapped in a cache to ensure the
// credentials are only refreshed when needed. This also protects the
// credential provider to be used concurrently.
//func resolveCredentialChain(ctx context.Context, cfg *aws.Config, configs configs) (err error) {
//	envConfig, _, _ := getAWSConfigSources(configs)
//
//	// When checking if a profile was specified programmatically we should only consider the "other"
//	// configuration sources that have been provided. This ensures we correctly honor the expected credential
//	// hierarchy.
//
//	switch {
//	case envConfig.Credentials.HasKeys():
//		ctx = addCredentialSource(ctx, aws.CredentialSourceEnvVars)
//		cfg.Credentials = credentials.StaticCredentialsProvider{Value: envConfig.Credentials, Source: getCredentialSources(ctx)}
//	case len(envConfig.WebIdentityTokenFilePath) > 0:
//		ctx = addCredentialSource(ctx, aws.CredentialSourceEnvVarsSTSWebIDToken)
//		//err = assumeWebIdentity(ctx, cfg, envConfig.WebIdentityTokenFilePath, envConfig.RoleARN, envConfig.RoleSessionName, configs)
//	default:
//		//ctx, err = resolveCredsFromProfile(ctx, cfg, envConfig, sharedConfig, other)
//	}
//	if err != nil {
//		return err
//	}
//
//	// Wrap the resolved provider in a cache so the SDK will cache credentials.
//	cfg.Credentials, err = wrapWithCredentialsCache(ctx, configs, cfg.Credentials)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// resolveCredentialProvider extracts the first instance of Credentials from the
// config slices.
//
// The resolved CredentialProvider will be wrapped in a cache to ensure the
// credentials are only refreshed when needed. This also protects the
// credential provider to be used concurrently.
//
// Config providers used:
// * credentialsProviderProvider

// resolveCredentialChain resolves a credential provider chain using EnvConfig
// and SharedConfig if present in the slice of provided configs.
//
// The resolved CredentialProvider will be wrapped in a cache to ensure the
// credentials are only refreshed when needed. This also protects the
// credential provider to be used concurrently.
func resolveCredentialChain(ctx context.Context, cfg *aws.Config, configs configs) (err error) {
	envConfig, sharedConfig, other := getAWSConfigSources(configs)

	// When checking if a profile was specified programmatically we should only consider the "other"
	// configuration sources that have been provided. This ensures we correctly honor the expected credential
	// hierarchy.
	_, sharedProfileSet, err := getSharedConfigProfile(ctx, other)
	if err != nil {
		return err
	}

	switch {
	case sharedProfileSet:
		ctx, err = resolveCredsFromProfile(ctx, cfg, envConfig, sharedConfig, other)

	//case
	case envConfig.Credentials.HasKeys():
		ctx = addCredentialSource(ctx, aws.CredentialSourceEnvVars)
		cfg.Credentials = credentials.StaticCredentialsProvider{Value: envConfig.Credentials, Source: getCredentialSources(ctx)}
	//case len(envConfig.WebIdentityTokenFilePath) > 0:
	//	ctx = addCredentialSource(ctx, aws.CredentialSourceEnvVarsSTSWebIDToken)
	//	err = assumeWebIdentity(ctx, cfg, envConfig.WebIdentityTokenFilePath, envConfig.RoleARN, envConfig.RoleSessionName, configs)
	default:
		ctx, err = resolveCredsFromProfile(ctx, cfg, envConfig, sharedConfig, other)
	}
	if err != nil {
		return err
	}

	// Wrap the resolved provider in a cache so the SDK will cache credentials.
	cfg.Credentials, err = wrapWithCredentialsCache(ctx, configs, cfg.Credentials)
	if err != nil {
		return err
	}

	return nil
}

func resolveCredsFromProfile(ctx context.Context, cfg *aws.Config, envConfig *EnvConfig, sharedConfig *SharedConfig, configs configs) (ctx2 context.Context, err error) {
	switch {
	case sharedConfig.Source != nil:
		ctx = addCredentialSource(ctx, aws.CredentialSourceProfileSourceProfile)
		// Assume IAM role with credentials source from a different profile.
		ctx, err = resolveCredsFromProfile(ctx, cfg, envConfig, sharedConfig.Source, configs)

	case sharedConfig.Credentials.HasKeys():
		// Static Credentials from Shared Config/Credentials file.
		ctx = addCredentialSource(ctx, aws.CredentialSourceProfile)
		cfg.Credentials = credentials.StaticCredentialsProvider{
			Value:  sharedConfig.Credentials,
			Source: getCredentialSources(ctx),
		}

	case len(sharedConfig.CredentialSource) != 0:
		ctx = addCredentialSource(ctx, aws.CredentialSourceProfileNamedProvider)
		ctx, err = resolveCredsFromSource(ctx, cfg, envConfig, sharedConfig, configs)

	//case sharedConfig.hasSSOConfiguration():
	//	if sharedConfig.hasLegacySSOConfiguration() {
	//		ctx = addCredentialSource(ctx, aws.CredentialSourceProfileSSOLegacy)
	//		ctx = addCredentialSource(ctx, aws.CredentialSourceSSOLegacy)
	//	} else {
	//		ctx = addCredentialSource(ctx, aws.CredentialSourceSSO)
	//	}
	//	if sharedConfig.SSOSession != nil {
	//		ctx = addCredentialSource(ctx, aws.CredentialSourceProfileSSO)
	//	}

	case len(sharedConfig.CredentialProcess) != 0:
		// Get credentials from CredentialProcess
		ctx = addCredentialSource(ctx, aws.CredentialSourceProfileProcess)
		ctx = addCredentialSource(ctx, aws.CredentialSourceProcess)

	default:
		ctx = addCredentialSource(ctx, aws.CredentialSourceIMDS)
	}
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

// isAllowedHost allows host to be loopback or known ECS/EKS container IPs
//
// host can either be an IP address OR an unresolved hostname - resolution will
// be automatically performed in the latter case
func isAllowedHost(host string) (bool, error) {
	if ip := net.ParseIP(host); ip != nil {
		return isIPAllowed(ip), nil
	}

	addrs, err := lookupHostFn(host)
	if err != nil {
		return false, err
	}

	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip == nil || !isIPAllowed(ip) {
			return false, nil
		}
	}

	return true, nil
}

func isIPAllowed(ip net.IP) bool {
	return ip.IsLoopback()
}

func resolveLocalHTTPCredProvider(ctx context.Context, cfg *aws.Config, endpointURL, authToken string, configs configs) error {
	var resolveErr error

	parsed, err := url.Parse(endpointURL)
	if err != nil {
		resolveErr = fmt.Errorf("invalid URL, %w", err)
	} else {
		host := parsed.Hostname()
		if len(host) == 0 {
			resolveErr = fmt.Errorf("unable to parse host from local HTTP cred provider URL")
		} else if parsed.Scheme == "http" {
			if isAllowedHost, allowHostErr := isAllowedHost(host); allowHostErr != nil {
				resolveErr = fmt.Errorf("failed to resolve host %q, %v", host, allowHostErr)
			} else if !isAllowedHost {
				resolveErr = fmt.Errorf("invalid endpoint host, %q, only loopback/ecs/eks hosts are allowed", host)
			}
		}
	}

	if resolveErr != nil {
		return resolveErr
	}

	return resolveHTTPCredProvider(ctx, cfg, endpointURL, authToken, configs)
}

func resolveHTTPCredProvider(ctx context.Context, cfg *aws.Config, url, authToken string, configs configs) error {
	optFns := []func(*endpointcreds.Options){
		func(options *endpointcreds.Options) {
			if len(authToken) != 0 {
				options.AuthorizationToken = authToken
			}
			if authFilePath := os.Getenv(httpProviderAuthFileEnvVar); authFilePath != "" {
				options.AuthorizationTokenProvider = endpointcreds.TokenProviderFunc(func() (string, error) {
					var contents []byte
					var err error
					if contents, err = ioutil.ReadFile(authFilePath); err != nil {
						return "", fmt.Errorf("failed to read authorization token from %v: %v", authFilePath, err)
					}
					return string(contents), nil
				})
			}
			options.APIOptions = cfg.APIOptions
			if cfg.Retryer != nil {
				options.Retryer = cfg.Retryer()
			}
			options.CredentialSources = getCredentialSources(ctx)
		},
	}

	optFn, found, err := getEndpointCredentialProviderOptions(ctx, configs)
	if err != nil {
		return err
	}
	if found {
		optFns = append(optFns, optFn)
	}

	provider := endpointcreds.New(url, optFns...)

	cfg.Credentials, err = wrapWithCredentialsCache(ctx, configs, provider, func(options *aws.CredentialsCacheOptions) {
		options.ExpiryWindow = 5 * time.Minute
	})
	if err != nil {
		return err
	}

	return nil
}

func resolveCredsFromSource(ctx context.Context, cfg *aws.Config, envConfig *EnvConfig, sharedCfg *SharedConfig, configs configs) (context.Context, error) {
	switch sharedCfg.CredentialSource {
	case credSourceEnvironment:
		ctx = addCredentialSource(ctx, aws.CredentialSourceHTTP)
		cfg.Credentials = credentials.StaticCredentialsProvider{Value: envConfig.Credentials, Source: getCredentialSources(ctx)}

	default:
		return ctx, fmt.Errorf("credential_source values must be EcsContainer, Ec2InstanceMetadata, or Environment")
	}

	return ctx, nil
}

func getAWSConfigSources(cfgs configs) (*EnvConfig, *SharedConfig, configs) {
	var (
		//ibmEnvConfig *EnvConfig
		envConfig    *EnvConfig
		sharedConfig *SharedConfig
		other        configs
	)

	for i := range cfgs {
		switch c := cfgs[i].(type) {
		case EnvConfig:
			if envConfig == nil {
				envConfig = &c
			}
		case *EnvConfig:
			if envConfig == nil {
				envConfig = c
			}
		case SharedConfig:
			if sharedConfig == nil {
				sharedConfig = &c
			}
		case *SharedConfig:
			if envConfig == nil {
				sharedConfig = c
			}
		default:
			other = append(other, c)
		}
	}

	if envConfig == nil {
		envConfig = &EnvConfig{}
	}

	if sharedConfig == nil {
		sharedConfig = &SharedConfig{}
	}

	return envConfig, sharedConfig, other
}

// wrapWithCredentialsCache will wrap provider with an aws.CredentialsCache
// with the provided options if the provider is not already a
// aws.CredentialsCache.
func wrapWithCredentialsCache(
	ctx context.Context,
	cfgs configs,
	provider aws.CredentialsProvider,
	optFns ...func(options *aws.CredentialsCacheOptions),
) (aws.CredentialsProvider, error) {
	_, ok := provider.(*aws.CredentialsCache)
	if ok {
		return provider, nil
	}

	credCacheOptions, optionsFound, err := getCredentialsCacheOptionsProvider(ctx, cfgs)
	if err != nil {
		return nil, err
	}

	// force allocation of a new slice if the additional options are
	// needed, to prevent overwriting the passed in slice of options.
	optFns = optFns[:len(optFns):len(optFns)]
	if optionsFound {
		optFns = append(optFns, credCacheOptions)
	}

	return aws.NewCredentialsCache(provider, optFns...), nil
}

// credentialSource stores the chain of providers that was used to create an instance of
// a credentials provider on the context
type credentialSource struct{}

func addCredentialSource(ctx context.Context, source aws.CredentialSource) context.Context {
	existing, ok := ctx.Value(credentialSource{}).([]aws.CredentialSource)
	if !ok {
		existing = []aws.CredentialSource{source}
	} else {
		existing = append(existing, source)
	}
	return context.WithValue(ctx, credentialSource{}, existing)
}

func getCredentialSources(ctx context.Context) []aws.CredentialSource {
	return ctx.Value(credentialSource{}).([]aws.CredentialSource)
}
