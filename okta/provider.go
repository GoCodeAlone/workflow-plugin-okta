// Package okta provides an exported Okta SDK client provider for cross-plugin imports.
package okta

import (
	"fmt"

	oktasdk "github.com/okta/okta-sdk-golang/v6/okta"
)

// Provider wraps an Okta SDK APIClient for typed API access.
type Provider struct {
	Client *oktasdk.APIClient
	OrgURL string
}

// Config holds the parameters needed to create an Okta SDK client.
type Config struct {
	OrgURL     string
	APIToken   string
	ClientID   string
	PrivateKey string
	AuthMode   string // "token" (default) or "private_key"
}

// NewProvider creates a Provider backed by the official Okta Go SDK v6.
func NewProvider(cfg Config) (*Provider, error) {
	opts := []oktasdk.ConfigSetter{
		oktasdk.WithOrgUrl(cfg.OrgURL),
		oktasdk.WithRateLimitMaxRetries(3),
	}

	if cfg.AuthMode == "private_key" && cfg.ClientID != "" && cfg.PrivateKey != "" {
		opts = append(opts,
			oktasdk.WithClientId(cfg.ClientID),
			oktasdk.WithPrivateKey(cfg.PrivateKey),
			oktasdk.WithAuthorizationMode("PrivateKey"),
		)
	} else if cfg.APIToken != "" {
		opts = append(opts, oktasdk.WithToken(cfg.APIToken))
	} else {
		return nil, fmt.Errorf("okta: either apiToken or clientId+privateKey are required")
	}

	sdkCfg, err := oktasdk.NewConfiguration(opts...)
	if err != nil {
		return nil, fmt.Errorf("okta: failed to create SDK configuration: %w", err)
	}

	client := oktasdk.NewAPIClient(sdkCfg)
	return &Provider{Client: client, OrgURL: cfg.OrgURL}, nil
}
