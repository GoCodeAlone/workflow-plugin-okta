package internal

import (
	"context"
	"strings"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

type authProviderDescribeStep struct {
	name   string
	config map[string]any
}

func newAuthProviderDescribeStep(name string, config map[string]any) sdk.StepInstance {
	return &authProviderDescribeStep{name: name, config: config}
}

func (s *authProviderDescribeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current, _, _ map[string]any) (*sdk.StepResult, error) {
	providerID := firstString("provider_id", s.config, current)
	if providerID == "" {
		providerID = "okta"
	}
	orgURL := firstString("org_url", s.config, current)
	if orgURL == "" {
		orgURL = firstString("orgUrl", s.config, current)
	}
	scopes := firstStringSlice("scopes", s.config, current)
	if len(scopes) == 0 {
		scopes = []string{"okta.users.manage", "okta.groups.manage", "okta.apps.manage"}
	}
	return &sdk.StepResult{Output: map[string]any{
		"providers": []map[string]any{oktaAuthProviderDescriptor(providerID, orgURL, scopes)},
	}}, nil
}

func oktaAuthProviderDescriptor(providerID, orgURL string, scopes []string) map[string]any {
	return map[string]any{
		"id":             providerID,
		"label":          "Okta",
		"description":    "Okta workforce identity, SSO, MFA, OAuth/OIDC authorization server, users, groups, applications, and audit administration.",
		"categories":     []string{"identity_management", "oauth2_oidc", "enterprise_sso", "mfa", "audit"},
		"implementation": "workflow-plugin-okta",
		"version":        Version,
		"docs_url":       "https://github.com/GoCodeAlone/workflow-plugin-okta",
		"support_level":  "management",
		"capabilities": []map[string]any{
			oktaCapability("okta_identity_management", "Identity management", "identity_management", "Manage Okta users, groups, group rules, app assignments, and lifecycle actions.", []string{"okta.users.manage", "okta.groups.manage", "okta.apps.manage"}, oktaManagementFields(orgURL, scopes)),
			oktaCapability("okta_oauth_oidc_admin", "OAuth/OIDC authorization server", "oauth2_oidc", "Manage Okta authorization servers, scopes, claims, policies, policy rules, and signing keys.", []string{"okta.authorizationServers.manage"}, oktaManagementFields(orgURL, scopes)),
			oktaCapability("okta_enterprise_sso", "Enterprise SSO", "enterprise_sso", "Manage Okta applications, identity providers, domains, brands, and org settings used for enterprise SSO.", []string{"okta.apps.manage", "okta.idps.manage", "okta.domains.manage"}, oktaManagementFields(orgURL, scopes)),
			oktaCapability("okta_mfa_policy", "MFA and authenticators", "mfa", "Manage Okta authenticators, factors, and sign-on policies.", []string{"okta.authenticators.manage", "okta.policies.manage"}, oktaManagementFields(orgURL, scopes)),
			oktaCapability("okta_audit_hooks", "Audit and hooks", "audit", "Read Okta system logs and manage event hooks and inline hooks.", []string{"okta.logs.read", "okta.eventHooks.manage", "okta.inlineHooks.manage"}, oktaManagementFields(orgURL, scopes)),
		},
	}
}

func oktaCapability(key, label, category, description string, appScopes []string, fields []map[string]any) map[string]any {
	return map[string]any{
		"key":                key,
		"label":              label,
		"category":           category,
		"description":        description,
		"supported":          true,
		"app_scopes":         appScopes,
		"admin_read_scopes":  []string{"admin.auth.providers.read"},
		"admin_write_scopes": []string{"admin.auth.providers.write"},
		"config_fields":      fields,
	}
}

func oktaManagementFields(orgURL string, scopes []string) []map[string]any {
	return []map[string]any{
		oktaField("okta_org_url", "Okta org URL", "url", "Base URL for the Okta org, for example https://dev-123456.okta.com.", "Use the exact org domain configured in Okta.", false, true, optionIfSet(orgURL)),
		oktaField("okta_auth_mode", "Credential mode", "select", "How Workflow authenticates administrative API calls to Okta.", "Prefer OAuth private key for least-privilege service integrations. API token remains supported for existing Okta tenants.", false, true, []map[string]any{
			{"value": "private_key", "label": "OAuth private key", "description": "Use an Okta API service app with client ID, private key, and explicit scopes."},
			{"value": "api_token", "label": "API token", "description": "Use an Okta SSWS API token. Treat as highly privileged secret material."},
		}),
		oktaField("okta_api_token", "API token", "secret", "Okta SSWS API token used when credential mode is API token.", "Write-only secret. Leave blank to keep the existing token.", true, false, nil),
		oktaField("okta_client_id", "OAuth client ID", "text", "Okta API service application client ID used for OAuth private key mode.", "Required when credential mode is OAuth private key.", false, false, nil),
		oktaField("okta_private_key", "OAuth private key", "secret", "Private key for the Okta API service application.", "Write-only secret. Store through the application's secret provider.", true, false, nil),
		oktaField("okta_scopes", "OAuth management scopes", "multiselect", "Okta management scopes granted to the service application.", "Select the least-privilege scopes required by enabled Okta capabilities.", false, false, scopeOptions(scopes)),
	}
}

func oktaField(key, label, inputType, description, helpText string, secret, required bool, options []map[string]any) map[string]any {
	return map[string]any{
		"key":         key,
		"label":       label,
		"input_type":  inputType,
		"description": description,
		"help_text":   helpText,
		"secret":      secret,
		"required":    required,
		"options":     options,
	}
}

func optionIfSet(value string) []map[string]any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return []map[string]any{{"value": value, "label": value}}
}

func scopeOptions(selected []string) []map[string]any {
	known := []string{
		"okta.users.manage",
		"okta.groups.manage",
		"okta.apps.manage",
		"okta.authorizationServers.manage",
		"okta.idps.manage",
		"okta.domains.manage",
		"okta.authenticators.manage",
		"okta.policies.manage",
		"okta.logs.read",
		"okta.eventHooks.manage",
		"okta.inlineHooks.manage",
	}
	seen := map[string]bool{}
	var options []map[string]any
	for _, scope := range append(known, selected...) {
		scope = strings.TrimSpace(scope)
		if scope == "" || seen[scope] {
			continue
		}
		seen[scope] = true
		options = append(options, map[string]any{"value": scope, "label": scope})
	}
	return options
}

func firstString(key string, sources ...map[string]any) string {
	for _, source := range sources {
		if source == nil {
			continue
		}
		if value, ok := source[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstStringSlice(key string, sources ...map[string]any) []string {
	for _, source := range sources {
		if values := resolveStringSlice(key, source, nil); len(values) > 0 {
			return values
		}
	}
	return nil
}
