package config

import (
	"context"
	smithybearer "github.com/aws/smithy-go/auth/bearer"
	"testing"
)

func TestWrapWithBearerAuthTokenProvider(t *testing.T) {
	cases := map[string]struct {
		configs         configs
		provider        smithybearer.TokenProvider
		optFns          []func(*smithybearer.TokenCacheOptions)
		compareProvider bool
		expectToken     smithybearer.Token
	}{
		"already wrapped": {
			provider: smithybearer.NewTokenCache(smithybearer.StaticTokenProvider{
				Token: smithybearer.Token{Value: "abc123"},
			}),
			compareProvider: true,
			expectToken:     smithybearer.Token{Value: "abc123"},
		},
		"to be wrapped": {
			provider: smithybearer.StaticTokenProvider{
				Token: smithybearer.Token{Value: "abc123"},
			},
			expectToken: smithybearer.Token{Value: "abc123"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			provider, err := wrapWithBearerAuthTokenCache(context.Background(),
				c.configs, c.provider, c.optFns...)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if p, ok := provider.(*smithybearer.TokenCache); !ok {
				t.Fatalf("expect provider wrapped in %T, got %T", p, provider)
			}

			if c.compareProvider && provider != c.provider {
				t.Errorf("expect same provider, was not")
			}

			token, err := provider.RetrieveBearerToken(context.Background())
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmpDiff(c.expectToken, token); diff != "" {
				t.Errorf("expect token match\n%s", diff)
			}
		})
	}
}
