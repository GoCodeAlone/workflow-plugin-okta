package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// orgGetStep implements step.okta_org_get
type orgGetStep struct {
	name       string
	moduleName string
}

func newOrgGetStep(name string, config map[string]any) (*orgGetStep, error) {
	return &orgGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *orgGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/org", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// orgUpdateStep implements step.okta_org_update
type orgUpdateStep struct {
	name       string
	moduleName string
}

func newOrgUpdateStep(name string, config map[string]any) (*orgUpdateStep, error) {
	return &orgUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *orgUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	body := map[string]any{}
	if companyName := resolveValue("companyName", current, config); companyName != "" {
		body["companyName"] = companyName
	}
	if website := resolveValue("website", current, config); website != "" {
		body["website"] = website
	}
	if phoneNumber := resolveValue("phoneNumber", current, config); phoneNumber != "" {
		body["phoneNumber"] = phoneNumber
	}
	if address1 := resolveValue("address1", current, config); address1 != "" {
		body["address1"] = address1
	}
	if city := resolveValue("city", current, config); city != "" {
		body["city"] = city
	}
	if state := resolveValue("state", current, config); state != "" {
		body["state"] = state
	}
	if country := resolveValue("country", current, config); country != "" {
		body["country"] = country
	}
	result, err := oktaPut(client, "/api/v1/org", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true}}, nil
}
