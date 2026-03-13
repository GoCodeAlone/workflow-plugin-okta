package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// domainCreateStep implements step.okta_domain_create
type domainCreateStep struct {
	name       string
	moduleName string
}

func newDomainCreateStep(name string, config map[string]any) (*domainCreateStep, error) {
	return &domainCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *domainCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	domain := resolveValue("domain", current, config)
	if domain == "" {
		return &sdk.StepResult{Output: errResult("domain is required")}, nil
	}
	body := map[string]any{
		"domain": domain,
	}
	result, err := oktaPost(client, "/api/v1/domains", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// domainGetStep implements step.okta_domain_get
type domainGetStep struct {
	name       string
	moduleName string
}

func newDomainGetStep(name string, config map[string]any) (*domainGetStep, error) {
	return &domainGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *domainGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	domainID := resolveValue("domainId", current, config)
	if domainID == "" {
		return &sdk.StepResult{Output: errResult("domainId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/domains/"+domainID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// domainListStep implements step.okta_domain_list
type domainListStep struct {
	name       string
	moduleName string
}

func newDomainListStep(name string, config map[string]any) (*domainListStep, error) {
	return &domainListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *domainListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/domains", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		if domains, ok := m["domains"].([]any); ok {
			return &sdk.StepResult{Output: listResult("domains", domains)}, nil
		}
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("domains", items)}, nil
}

// domainDeleteStep implements step.okta_domain_delete
type domainDeleteStep struct {
	name       string
	moduleName string
}

func newDomainDeleteStep(name string, config map[string]any) (*domainDeleteStep, error) {
	return &domainDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *domainDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	domainID := resolveValue("domainId", current, config)
	if domainID == "" {
		return &sdk.StepResult{Output: errResult("domainId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/domains/"+domainID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// domainVerifyStep implements step.okta_domain_verify
type domainVerifyStep struct {
	name       string
	moduleName string
}

func newDomainVerifyStep(name string, config map[string]any) (*domainVerifyStep, error) {
	return &domainVerifyStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *domainVerifyStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	domainID := resolveValue("domainId", current, config)
	if domainID == "" {
		return &sdk.StepResult{Output: errResult("domainId is required")}, nil
	}
	result, err := oktaPost(client, "/api/v1/domains/"+domainID+"/verify", map[string]any{})
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"verified": true}}, nil
}
