package internal

import (
	"context"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-okta/internal/contracts"
	pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func TestContractRegistryIncludesStrictProtoDescriptors(t *testing.T) {
	provider, ok := NewOktaPlugin().(interface {
		ContractRegistry() *pb.ContractRegistry
	})
	if !ok {
		t.Fatal("plugin does not expose ContractRegistry")
	}
	registry := provider.ContractRegistry()
	if registry == nil {
		t.Fatal("ContractRegistry is nil")
	}
	if registry.GetFileDescriptorSet() == nil || len(registry.GetFileDescriptorSet().GetFile()) == 0 {
		t.Fatal("ContractRegistry does not include file descriptors")
	}

	contractsByType := map[string]*pb.ContractDescriptor{}
	for _, contract := range registry.GetContracts() {
		switch contract.GetKind() {
		case pb.ContractKind_CONTRACT_KIND_MODULE:
			contractsByType["module:"+contract.GetModuleType()] = contract
		case pb.ContractKind_CONTRACT_KIND_STEP:
			contractsByType["step:"+contract.GetStepType()] = contract
		}
	}

	module := contractsByType["module:okta.provider"]
	if module == nil {
		t.Fatal("missing okta.provider contract")
	}
	if module.GetMode() != pb.ContractMode_CONTRACT_MODE_STRICT_PROTO {
		t.Fatalf("okta.provider mode = %v, want strict proto", module.GetMode())
	}
	if module.GetConfigMessage() != "workflow.plugins.okta.v1.ProviderConfig" {
		t.Fatalf("okta.provider config = %q", module.GetConfigMessage())
	}

	for _, stepType := range allStepTypes() {
		contract := contractsByType["step:"+stepType]
		if contract == nil {
			t.Fatalf("missing contract for %s", stepType)
		}
		if contract.GetMode() != pb.ContractMode_CONTRACT_MODE_STRICT_PROTO {
			t.Fatalf("%s mode = %v, want strict proto", stepType, contract.GetMode())
		}
	}
}

func TestAuthProviderDescribeStepAdvertisesRealCapabilities(t *testing.T) {
	step := newAuthProviderDescribeStep("describe", map[string]any{
		"org_url": "https://dev-123456.okta.com",
		"scopes":  []string{"okta.users.manage", "okta.logs.read"},
	})
	result, err := step.Execute(context.Background(), nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	providers, ok := result.Output["providers"].([]map[string]any)
	if !ok || len(providers) != 1 {
		t.Fatalf("providers = %#v", result.Output["providers"])
	}
	provider := providers[0]
	categories := stringSet(provider["categories"].([]string))
	for _, category := range []string{"identity_management", "oauth2_oidc", "enterprise_sso", "mfa", "audit"} {
		if !categories[category] {
			t.Fatalf("missing category %q in %#v", category, categories)
		}
	}
	if categories["directory_sync"] {
		t.Fatalf("descriptor must not advertise directory_sync without SCIM/directory steps")
	}

	capabilities := provider["capabilities"].([]map[string]any)
	if len(capabilities) != 5 {
		t.Fatalf("capability count = %d, want 5", len(capabilities))
	}
	for _, capability := range capabilities {
		if capability["supported"] != true {
			t.Fatalf("%s supported = %#v, want true", capability["key"], capability["supported"])
		}
		fields := capability["config_fields"].([]map[string]any)
		if len(fields) == 0 {
			t.Fatalf("%s has no config fields", capability["key"])
		}
	}
}

func TestTypedAuthProviderDescribeOutput(t *testing.T) {
	result, err := typedAuthProviderDescribe(context.Background(), sdk.TypedStepRequest[*contracts.AuthProviderDescribeConfig, *contracts.AuthProviderDescribeInput]{
		Config: &contracts.AuthProviderDescribeConfig{ProviderId: "okta-admin"},
		Input:  &contracts.AuthProviderDescribeInput{OrgUrl: "https://dev-123456.okta.com", Scopes: []string{"okta.users.manage"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output == nil || len(result.Output.GetProviders()) != 1 {
		t.Fatalf("providers = %#v", result.Output)
	}
	provider := result.Output.GetProviders()[0]
	if provider.GetId() != "okta-admin" {
		t.Fatalf("provider id = %q", provider.GetId())
	}
	if len(provider.GetCapabilities()) != 5 {
		t.Fatalf("capabilities = %d, want 5", len(provider.GetCapabilities()))
	}
}

func TestTypedModuleConfigMapsProtoNamesToLegacyNames(t *testing.T) {
	cfg := typedOktaModuleConfig(&contracts.ProviderConfig{
		OrgUrl:     "https://dev-123456.okta.com",
		ApiToken:   "tok",
		ClientId:   "client",
		PrivateKey: "key",
		Scopes:     []string{"okta.users.manage"},
	})
	for _, key := range []string{"orgUrl", "apiToken", "clientId", "privateKey", "scopes"} {
		if _, ok := cfg[key]; !ok {
			t.Fatalf("missing legacy key %q in %#v", key, cfg)
		}
	}
}

func stringSet(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}
