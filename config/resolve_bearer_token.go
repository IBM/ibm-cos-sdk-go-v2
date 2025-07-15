package config

import (
	"context"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	smithybearer "github.com/aws/smithy-go/auth/bearer"
)

// resolveBearerAuthTokenProvider extracts the first instance of
// BearerAuthTokenProvider from the config sources.
//
// The resolved BearerAuthTokenProvider will be wrapped in a cache to ensure
// the Token is only refreshed when needed. This also protects the
// TokenProvider so it can be used concurrently.
//
// Config providers used:
// * bearerAuthTokenProviderProvider
func resolveBearerAuthTokenProvider(ctx context.Context, cfg *aws.Config, configs configs) (bool, error) {
	tokenProvider, found, err := getBearerAuthTokenProvider(ctx, configs)
	if !found || err != nil {
		return false, err
	}

	cfg.BearerAuthTokenProvider, err = wrapWithBearerAuthTokenCache(
		ctx, configs, tokenProvider)
	if err != nil {
		return false, err
	}

	return true, nil
}

// wrapWithBearerAuthTokenCache will wrap provider with an smithy-go
// bearer/auth#TokenCache with the provided options if the provider is not
// already a TokenCache.
func wrapWithBearerAuthTokenCache(
	ctx context.Context,
	cfgs configs,
	provider smithybearer.TokenProvider,
	optFns ...func(*smithybearer.TokenCacheOptions),
) (smithybearer.TokenProvider, error) {
	_, ok := provider.(*smithybearer.TokenCache)
	if ok {
		return provider, nil
	}

	tokenCacheConfigOptions, optionsFound, err := getBearerAuthTokenCacheOptions(ctx, cfgs)
	if err != nil {
		return nil, err
	}

	opts := make([]func(*smithybearer.TokenCacheOptions), 0, 2+len(optFns))
	opts = append(opts, func(o *smithybearer.TokenCacheOptions) {
		o.RefreshBeforeExpires = 5 * time.Minute
		o.RetrieveBearerTokenTimeout = 30 * time.Second
	})
	opts = append(opts, optFns...)
	if optionsFound {
		opts = append(opts, tokenCacheConfigOptions)
	}

	return smithybearer.NewTokenCache(provider, opts...), nil
}
