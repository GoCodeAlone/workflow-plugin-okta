package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// policyCreateStep implements step.okta_policy_create
type policyCreateStep struct {
	name       string
	moduleName string
}

func newPolicyCreateStep(name string, config map[string]any) (*policyCreateStep, error) {
	return &policyCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyName := resolveValue("name", current, config)
	policyType := resolveValue("type", current, config)
	if policyName == "" || policyType == "" {
		return &sdk.StepResult{Output: errResult("name and type are required")}, nil
	}
	body := map[string]any{
		"type":        policyType,
		"name":        policyName,
		"description": resolveValue("description", current, config),
		"priority":    1,
		"status":      "ACTIVE",
	}
	if conditions := resolveMap("conditions", current, config); conditions != nil {
		body["conditions"] = conditions
	}
	result, err := oktaPost(client, "/api/v1/policies", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// policyGetStep implements step.okta_policy_get
type policyGetStep struct {
	name       string
	moduleName string
}

func newPolicyGetStep(name string, config map[string]any) (*policyGetStep, error) {
	return &policyGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	if policyID == "" {
		return &sdk.StepResult{Output: errResult("policyId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/policies/"+policyID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// policyListStep implements step.okta_policy_list
type policyListStep struct {
	name       string
	moduleName string
}

func newPolicyListStep(name string, config map[string]any) (*policyListStep, error) {
	return &policyListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/policies", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("policies", items)}, nil
}

// policyDeleteStep implements step.okta_policy_delete
type policyDeleteStep struct {
	name       string
	moduleName string
}

func newPolicyDeleteStep(name string, config map[string]any) (*policyDeleteStep, error) {
	return &policyDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	if policyID == "" {
		return &sdk.StepResult{Output: errResult("policyId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/policies/"+policyID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// policyActivateStep implements step.okta_policy_activate
type policyActivateStep struct {
	name       string
	moduleName string
}

func newPolicyActivateStep(name string, config map[string]any) (*policyActivateStep, error) {
	return &policyActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	if policyID == "" {
		return &sdk.StepResult{Output: errResult("policyId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/policies/"+policyID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// policyDeactivateStep implements step.okta_policy_deactivate
type policyDeactivateStep struct {
	name       string
	moduleName string
}

func newPolicyDeactivateStep(name string, config map[string]any) (*policyDeactivateStep, error) {
	return &policyDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	if policyID == "" {
		return &sdk.StepResult{Output: errResult("policyId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/policies/"+policyID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// policyRuleCreateStep implements step.okta_policy_rule_create
type policyRuleCreateStep struct {
	name       string
	moduleName string
}

func newPolicyRuleCreateStep(name string, config map[string]any) (*policyRuleCreateStep, error) {
	return &policyRuleCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyRuleCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	ruleName := resolveValue("name", current, config)
	if policyID == "" || ruleName == "" {
		return &sdk.StepResult{Output: errResult("policyId and name are required")}, nil
	}
	body := map[string]any{
		"name":       ruleName,
		"type":       resolveValue("type", current, config),
		"conditions": resolveMap("conditions", current, config),
		"actions":    resolveMap("actions", current, config),
	}
	result, err := oktaPost(client, "/api/v1/policies/"+policyID+"/rules", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// policyRuleListStep implements step.okta_policy_rule_list
type policyRuleListStep struct {
	name       string
	moduleName string
}

func newPolicyRuleListStep(name string, config map[string]any) (*policyRuleListStep, error) {
	return &policyRuleListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyRuleListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	if policyID == "" {
		return &sdk.StepResult{Output: errResult("policyId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/policies/"+policyID+"/rules", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("rules", items)}, nil
}

// policyRuleDeleteStep implements step.okta_policy_rule_delete
type policyRuleDeleteStep struct {
	name       string
	moduleName string
}

func newPolicyRuleDeleteStep(name string, config map[string]any) (*policyRuleDeleteStep, error) {
	return &policyRuleDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyRuleDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	ruleID := resolveValue("ruleId", current, config)
	if policyID == "" || ruleID == "" {
		return &sdk.StepResult{Output: errResult("policyId and ruleId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/policies/"+policyID+"/rules/"+ruleID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// policyRuleActivateStep implements step.okta_policy_rule_activate
type policyRuleActivateStep struct {
	name       string
	moduleName string
}

func newPolicyRuleActivateStep(name string, config map[string]any) (*policyRuleActivateStep, error) {
	return &policyRuleActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyRuleActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	ruleID := resolveValue("ruleId", current, config)
	if policyID == "" || ruleID == "" {
		return &sdk.StepResult{Output: errResult("policyId and ruleId are required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/policies/"+policyID+"/rules/"+ruleID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// policyRuleDeactivateStep implements step.okta_policy_rule_deactivate
type policyRuleDeactivateStep struct {
	name       string
	moduleName string
}

func newPolicyRuleDeactivateStep(name string, config map[string]any) (*policyRuleDeactivateStep, error) {
	return &policyRuleDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *policyRuleDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	policyID := resolveValue("policyId", current, config)
	ruleID := resolveValue("ruleId", current, config)
	if policyID == "" || ruleID == "" {
		return &sdk.StepResult{Output: errResult("policyId and ruleId are required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/policies/"+policyID+"/rules/"+ruleID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}
