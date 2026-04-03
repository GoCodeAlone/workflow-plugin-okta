package internal

import (
	"sync"

	oktasdk "github.com/okta/okta-sdk-golang/v6/okta"
)

// OktaClient wraps the official Okta SDK APIClient for REST API access.
type OktaClient struct {
	SdkClient *oktasdk.APIClient
	OrgURL    string
	APIToken  string
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
