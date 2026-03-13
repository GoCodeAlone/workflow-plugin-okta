package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// authzServerCreateStep implements step.okta_authz_server_create
type authzServerCreateStep struct {
	name       string
	moduleName string
}

func newAuthzServerCreateStep(name string, config map[string]any) (*authzServerCreateStep, error) {
	return &authzServerCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverName := resolveValue("name", current, config)
	if serverName == "" {
		return &sdk.StepResult{Output: errResult("name is required")}, nil
	}
	body := map[string]any{
		"name":        serverName,
		"description": resolveValue("description", current, config),
		"audiences":   resolveStringSlice("audiences", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzServerGetStep implements step.okta_authz_server_get
type authzServerGetStep struct {
	name       string
	moduleName string
}

func newAuthzServerGetStep(name string, config map[string]any) (*authzServerGetStep, error) {
	return &authzServerGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzServerListStep implements step.okta_authz_server_list
type authzServerListStep struct {
	name       string
	moduleName string
}

func newAuthzServerListStep(name string, config map[string]any) (*authzServerListStep, error) {
	return &authzServerListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("authorizationServers", items)}, nil
}

// authzServerUpdateStep implements step.okta_authz_server_update
type authzServerUpdateStep struct {
	name       string
	moduleName string
}

func newAuthzServerUpdateStep(name string, config map[string]any) (*authzServerUpdateStep, error) {
	return &authzServerUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	body := map[string]any{}
	if n := resolveValue("name", current, config); n != "" {
		body["name"] = n
	}
	if d := resolveValue("description", current, config); d != "" {
		body["description"] = d
	}
	result, err := oktaPut(client, "/api/v1/authorizationServers/"+serverID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzServerDeleteStep implements step.okta_authz_server_delete
type authzServerDeleteStep struct {
	name       string
	moduleName string
}

func newAuthzServerDeleteStep(name string, config map[string]any) (*authzServerDeleteStep, error) {
	return &authzServerDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/authorizationServers/"+serverID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// authzServerActivateStep implements step.okta_authz_server_activate
type authzServerActivateStep struct {
	name       string
	moduleName string
}

func newAuthzServerActivateStep(name string, config map[string]any) (*authzServerActivateStep, error) {
	return &authzServerActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/authorizationServers/"+serverID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// authzServerDeactivateStep implements step.okta_authz_server_deactivate
type authzServerDeactivateStep struct {
	name       string
	moduleName string
}

func newAuthzServerDeactivateStep(name string, config map[string]any) (*authzServerDeactivateStep, error) {
	return &authzServerDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzServerDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/authorizationServers/"+serverID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// authzClaimCreateStep implements step.okta_authz_claim_create
type authzClaimCreateStep struct {
	name       string
	moduleName string
}

func newAuthzClaimCreateStep(name string, config map[string]any) (*authzClaimCreateStep, error) {
	return &authzClaimCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzClaimCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	claimName := resolveValue("name", current, config)
	if serverID == "" || claimName == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and name are required")}, nil
	}
	body := map[string]any{
		"name":      claimName,
		"valueType": resolveValue("valueType", current, config),
		"value":     resolveValue("value", current, config),
		"claimType": resolveValue("claimType", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers/"+serverID+"/claims", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzClaimListStep implements step.okta_authz_claim_list
type authzClaimListStep struct {
	name       string
	moduleName string
}

func newAuthzClaimListStep(name string, config map[string]any) (*authzClaimListStep, error) {
	return &authzClaimListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzClaimListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID+"/claims", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("claims", items)}, nil
}

// authzClaimDeleteStep implements step.okta_authz_claim_delete
type authzClaimDeleteStep struct {
	name       string
	moduleName string
}

func newAuthzClaimDeleteStep(name string, config map[string]any) (*authzClaimDeleteStep, error) {
	return &authzClaimDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzClaimDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	claimID := resolveValue("claimId", current, config)
	if serverID == "" || claimID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and claimId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/authorizationServers/"+serverID+"/claims/"+claimID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// authzScopeCreateStep implements step.okta_authz_scope_create
type authzScopeCreateStep struct {
	name       string
	moduleName string
}

func newAuthzScopeCreateStep(name string, config map[string]any) (*authzScopeCreateStep, error) {
	return &authzScopeCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzScopeCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	scopeName := resolveValue("name", current, config)
	if serverID == "" || scopeName == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and name are required")}, nil
	}
	body := map[string]any{
		"name":        scopeName,
		"description": resolveValue("description", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers/"+serverID+"/scopes", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzScopeListStep implements step.okta_authz_scope_list
type authzScopeListStep struct {
	name       string
	moduleName string
}

func newAuthzScopeListStep(name string, config map[string]any) (*authzScopeListStep, error) {
	return &authzScopeListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzScopeListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID+"/scopes", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("scopes", items)}, nil
}

// authzScopeDeleteStep implements step.okta_authz_scope_delete
type authzScopeDeleteStep struct {
	name       string
	moduleName string
}

func newAuthzScopeDeleteStep(name string, config map[string]any) (*authzScopeDeleteStep, error) {
	return &authzScopeDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzScopeDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	scopeID := resolveValue("scopeId", current, config)
	if serverID == "" || scopeID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and scopeId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/authorizationServers/"+serverID+"/scopes/"+scopeID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// authzPolicyCreateStep implements step.okta_authz_policy_create
type authzPolicyCreateStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyCreateStep(name string, config map[string]any) (*authzPolicyCreateStep, error) {
	return &authzPolicyCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	policyName := resolveValue("name", current, config)
	if serverID == "" || policyName == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and name are required")}, nil
	}
	body := map[string]any{
		"type":        "OAUTH_AUTHORIZATION_POLICY",
		"name":        policyName,
		"description": resolveValue("description", current, config),
		"priority":    1,
		"conditions":  resolveMap("conditions", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers/"+serverID+"/policies", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzPolicyListStep implements step.okta_authz_policy_list
type authzPolicyListStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyListStep(name string, config map[string]any) (*authzPolicyListStep, error) {
	return &authzPolicyListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID+"/policies", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("policies", items)}, nil
}

// authzPolicyDeleteStep implements step.okta_authz_policy_delete
type authzPolicyDeleteStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyDeleteStep(name string, config map[string]any) (*authzPolicyDeleteStep, error) {
	return &authzPolicyDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	policyID := resolveValue("policyId", current, config)
	if serverID == "" || policyID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and policyId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/authorizationServers/"+serverID+"/policies/"+policyID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// authzPolicyRuleCreateStep implements step.okta_authz_policy_rule_create
type authzPolicyRuleCreateStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyRuleCreateStep(name string, config map[string]any) (*authzPolicyRuleCreateStep, error) {
	return &authzPolicyRuleCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyRuleCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	policyID := resolveValue("policyId", current, config)
	ruleName := resolveValue("name", current, config)
	if serverID == "" || policyID == "" || ruleName == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId, policyId and name are required")}, nil
	}
	body := map[string]any{
		"type":       "RESOURCE_ACCESS",
		"name":       ruleName,
		"priority":   1,
		"conditions": resolveMap("conditions", current, config),
		"actions":    resolveMap("actions", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers/"+serverID+"/policies/"+policyID+"/rules", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authzPolicyRuleListStep implements step.okta_authz_policy_rule_list
type authzPolicyRuleListStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyRuleListStep(name string, config map[string]any) (*authzPolicyRuleListStep, error) {
	return &authzPolicyRuleListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyRuleListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	policyID := resolveValue("policyId", current, config)
	if serverID == "" || policyID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId and policyId are required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID+"/policies/"+policyID+"/rules", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("rules", items)}, nil
}

// authzPolicyRuleDeleteStep implements step.okta_authz_policy_rule_delete
type authzPolicyRuleDeleteStep struct {
	name       string
	moduleName string
}

func newAuthzPolicyRuleDeleteStep(name string, config map[string]any) (*authzPolicyRuleDeleteStep, error) {
	return &authzPolicyRuleDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzPolicyRuleDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	policyID := resolveValue("policyId", current, config)
	ruleID := resolveValue("ruleId", current, config)
	if serverID == "" || policyID == "" || ruleID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId, policyId and ruleId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/authorizationServers/"+serverID+"/policies/"+policyID+"/rules/"+ruleID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// authzKeyListStep implements step.okta_authz_key_list
type authzKeyListStep struct {
	name       string
	moduleName string
}

func newAuthzKeyListStep(name string, config map[string]any) (*authzKeyListStep, error) {
	return &authzKeyListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzKeyListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authorizationServers/"+serverID+"/credentials/keys", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("keys", items)}, nil
}

// authzKeyRotateStep implements step.okta_authz_key_rotate
type authzKeyRotateStep struct {
	name       string
	moduleName string
}

func newAuthzKeyRotateStep(name string, config map[string]any) (*authzKeyRotateStep, error) {
	return &authzKeyRotateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authzKeyRotateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	serverID := resolveValue("authorizationServerId", current, config)
	if serverID == "" {
		return &sdk.StepResult{Output: errResult("authorizationServerId is required")}, nil
	}
	body := map[string]any{
		"use": resolveValue("use", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authorizationServers/"+serverID+"/credentials/lifecycle/keyRotate", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("keys", items)}, nil
}
