package internal

import (
	"context"
	"fmt"
	"net/http"
)

// oktaModule creates an Okta REST client and registers it.
type oktaModule struct {
	name   string
	config map[string]any
}

func newOktaModule(name string, config map[string]any) (*oktaModule, error) {
	return &oktaModule{name: name, config: config}, nil
}

// Init creates the Okta client and registers it in the global registry.
func (m *oktaModule) Init() error {
	orgURL, _ := m.config["orgUrl"].(string)
	if orgURL == "" {
		return fmt.Errorf("okta.provider %q: orgUrl is required", m.name)
	}

	apiToken, _ := m.config["apiToken"].(string)
	clientID, _ := m.config["clientId"].(string)
	privateKey, _ := m.config["privateKey"].(string)

	if apiToken == "" && (clientID == "" || privateKey == "") {
		return fmt.Errorf("okta.provider %q: either apiToken or clientId+privateKey are required", m.name)
	}

	httpClient := &http.Client{}

	client := &OktaClient{
		HTTPClient: httpClient,
		OrgURL:     orgURL,
		APIToken:   apiToken,
	}

	RegisterClient(m.name, client)
	return nil
}

// Start is a no-op for this module.
func (m *oktaModule) Start(_ context.Context) error { return nil }

// Stop unregisters the Okta client.
func (m *oktaModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
