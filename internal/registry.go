package internal

import (
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"

	oktasdk "github.com/okta/okta-sdk-golang/v6/okta"
)

// OktaClient wraps the official Okta SDK APIClient for REST API access.
type OktaClient struct {
	SdkClient  *oktasdk.APIClient
	OrgURL     string
	APIToken   string
	AuthMode   string // "token" (default) or "private_key"
	tokenCache *goCache.Cache
}

// NewOktaClient creates an OktaClient with the appropriate auth configuration.
func NewOktaClient(sdkClient *oktasdk.APIClient, orgURL, apiToken, authMode string) *OktaClient {
	c := &OktaClient{
		SdkClient: sdkClient,
		OrgURL:    orgURL,
		APIToken:  apiToken,
		AuthMode:  authMode,
	}
	if authMode == "private_key" {
		c.tokenCache = goCache.New(5*time.Minute, 10*time.Minute)
	}
	return c
}

var (
	clientMu       sync.RWMutex
	clientRegistry = make(map[string]*OktaClient)
)

// RegisterClient adds an Okta client to the global registry under the given name.
func RegisterClient(name string, c *OktaClient) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clientRegistry[name] = c
}

// GetClient looks up an Okta client by name.
func GetClient(name string) (*OktaClient, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := clientRegistry[name]
	return c, ok
}

// UnregisterClient removes a client from the registry.
func UnregisterClient(name string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clientRegistry, name)
}
