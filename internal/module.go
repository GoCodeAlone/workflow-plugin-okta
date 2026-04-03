package internal

import (
	"context"
	"fmt"

	oktaprovider "github.com/GoCodeAlone/workflow-plugin-okta/okta"
)

// oktaModule creates an Okta SDK client and registers it.
type oktaModule struct {
	name   string
	config map[string]any
}

func newOktaModule(name string, config map[string]any) (*oktaModule, error) {
	return &oktaModule{name: name, config: config}, nil
}

// Init creates the Okta SDK client via the provider and registers it.
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

	cfg := oktaprovider.Config{
		OrgURL:     orgURL,
		APIToken:   apiToken,
		ClientID:   clientID,
		PrivateKey: privateKey,
	}

	authMode := "token"
	if apiToken == "" {
		authMode = "private_key"
		cfg.AuthMode = authMode
	}

	scopes := resolveStringSlice("scopes", m.config, nil)
	if authMode == "private_key" && len(scopes) == 0 {
		scopes = []string{"okta.users.manage", "okta.groups.manage", "okta.apps.manage"}
	}
	cfg.Scopes = scopes

	provider, err := oktaprovider.NewProvider(cfg)
	if err != nil {
		return fmt.Errorf("okta.provider %q: %w", m.name, err)
	}

	RegisterClient(m.name, NewOktaClient(provider.Client, orgURL, apiToken, authMode))
	return nil
}

// Start is a no-op for this module.
func (m *oktaModule) Start(_ context.Context) error { return nil }

// Stop unregisters the Okta client.
func (m *oktaModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
