package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// idpCreateStep implements step.okta_idp_create
type idpCreateStep struct {
	name       string
	moduleName string
}

func newIdpCreateStep(name string, config map[string]any) (*idpCreateStep, error) {
	return &idpCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	idpName := resolveValue("name", current, config)
	idpType := resolveValue("type", current, config)
	if idpName == "" || idpType == "" {
		return &sdk.StepResult{Output: errResult("name and type are required")}, nil
	}
	body := map[string]any{
		"type":     idpType,
		"name":     idpName,
		"protocol": resolveMap("protocol", current, config),
		"policy":   resolveMap("policy", current, config),
	}
	result, err := oktaPost(client, "/api/v1/idps", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// idpGetStep implements step.okta_idp_get
type idpGetStep struct {
	name       string
	moduleName string
}

func newIdpGetStep(name string, config map[string]any) (*idpGetStep, error) {
	return &idpGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	idpID := resolveValue("idpId", current, config)
	if idpID == "" {
		return &sdk.StepResult{Output: errResult("idpId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/idps/"+idpID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// idpListStep implements step.okta_idp_list
type idpListStep struct {
	name       string
	moduleName string
}

func newIdpListStep(name string, config map[string]any) (*idpListStep, error) {
	return &idpListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/idps", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("idps", items)}, nil
}

// idpDeleteStep implements step.okta_idp_delete
type idpDeleteStep struct {
	name       string
	moduleName string
}

func newIdpDeleteStep(name string, config map[string]any) (*idpDeleteStep, error) {
	return &idpDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	idpID := resolveValue("idpId", current, config)
	if idpID == "" {
		return &sdk.StepResult{Output: errResult("idpId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/idps/"+idpID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// idpActivateStep implements step.okta_idp_activate
type idpActivateStep struct {
	name       string
	moduleName string
}

func newIdpActivateStep(name string, config map[string]any) (*idpActivateStep, error) {
	return &idpActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	idpID := resolveValue("idpId", current, config)
	if idpID == "" {
		return &sdk.StepResult{Output: errResult("idpId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/idps/"+idpID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// idpDeactivateStep implements step.okta_idp_deactivate
type idpDeactivateStep struct {
	name       string
	moduleName string
}

func newIdpDeactivateStep(name string, config map[string]any) (*idpDeactivateStep, error) {
	return &idpDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *idpDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	idpID := resolveValue("idpId", current, config)
	if idpID == "" {
		return &sdk.StepResult{Output: errResult("idpId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/idps/"+idpID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}
