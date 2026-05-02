package internal

import (
	"context"
	"testing"
)

func TestModuleInit_RegistersClient(t *testing.T) {
	m, err := newOktaModule("test-init", map[string]any{
		"orgUrl":   "https://dev-test.okta.com",
		"apiToken": "test-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	c, ok := GetClient("test-init")
	if !ok || c == nil {
		t.Error("expected client to be registered")
	}
	// cleanup
	UnregisterClient("test-init")
}

func TestModuleStop_UnregistersClient(t *testing.T) {
	m, _ := newOktaModule("test-stop", map[string]any{
		"orgUrl":   "https://dev-test.okta.com",
		"apiToken": "test-token",
	})
	_ = m.Init()
	_ = m.Stop(context.Background())
	_, ok := GetClient("test-stop")
	if ok {
		t.Error("expected client to be unregistered after stop")
	}
}

func TestModuleInit_MissingOrgUrl(t *testing.T) {
	m, err := newOktaModule("test-missing-url", map[string]any{
		"apiToken": "test-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing orgUrl")
		UnregisterClient("test-missing-url")
	}
}

func TestModuleInit_MissingCredentials(t *testing.T) {
	m, err := newOktaModule("test-missing-creds", map[string]any{
		"orgUrl": "https://dev-test.okta.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for missing credentials")
		UnregisterClient("test-missing-creds")
	}
}

func TestModuleInit_RejectsMixedCredentialModes(t *testing.T) {
	m, err := newOktaModule("test-mixed-creds", map[string]any{
		"orgUrl":     "https://dev-test.okta.com",
		"apiToken":   "test-token",
		"clientId":   "0oa123",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err == nil {
		t.Error("expected error for mixed credential modes")
		UnregisterClient("test-mixed-creds")
	}
}

func TestModuleInit_WithClientIdAndPrivateKey(t *testing.T) {
	m, err := newOktaModule("test-oauth", map[string]any{
		"orgUrl":     "https://dev-test.okta.com",
		"clientId":   "0oa123",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	c, ok := GetClient("test-oauth")
	if !ok || c == nil {
		t.Error("expected client to be registered with OAuth credentials")
	}
	UnregisterClient("test-oauth")
}
