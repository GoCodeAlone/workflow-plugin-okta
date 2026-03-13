package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// authenticatorCreateStep implements step.okta_authenticator_create
type authenticatorCreateStep struct {
	name       string
	moduleName string
}

func newAuthenticatorCreateStep(name string, config map[string]any) (*authenticatorCreateStep, error) {
	return &authenticatorCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authenticatorCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	key := resolveValue("key", current, config)
	if key == "" {
		return &sdk.StepResult{Output: errResult("key is required")}, nil
	}
	body := map[string]any{
		"key":      key,
		"name":     resolveValue("name", current, config),
		"type":     resolveValue("type", current, config),
		"settings": resolveMap("settings", current, config),
	}
	result, err := oktaPost(client, "/api/v1/authenticators", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authenticatorGetStep implements step.okta_authenticator_get
type authenticatorGetStep struct {
	name       string
	moduleName string
}

func newAuthenticatorGetStep(name string, config map[string]any) (*authenticatorGetStep, error) {
	return &authenticatorGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authenticatorGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	authenticatorID := resolveValue("authenticatorId", current, config)
	if authenticatorID == "" {
		return &sdk.StepResult{Output: errResult("authenticatorId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/authenticators/"+authenticatorID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// authenticatorListStep implements step.okta_authenticator_list
type authenticatorListStep struct {
	name       string
	moduleName string
}

func newAuthenticatorListStep(name string, config map[string]any) (*authenticatorListStep, error) {
	return &authenticatorListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authenticatorListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/authenticators", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("authenticators", items)}, nil
}

// authenticatorActivateStep implements step.okta_authenticator_activate
type authenticatorActivateStep struct {
	name       string
	moduleName string
}

func newAuthenticatorActivateStep(name string, config map[string]any) (*authenticatorActivateStep, error) {
	return &authenticatorActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authenticatorActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	authenticatorID := resolveValue("authenticatorId", current, config)
	if authenticatorID == "" {
		return &sdk.StepResult{Output: errResult("authenticatorId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/authenticators/"+authenticatorID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// authenticatorDeactivateStep implements step.okta_authenticator_deactivate
type authenticatorDeactivateStep struct {
	name       string
	moduleName string
}

func newAuthenticatorDeactivateStep(name string, config map[string]any) (*authenticatorDeactivateStep, error) {
	return &authenticatorDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *authenticatorDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	authenticatorID := resolveValue("authenticatorId", current, config)
	if authenticatorID == "" {
		return &sdk.StepResult{Output: errResult("authenticatorId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/authenticators/"+authenticatorID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// factorEnrollStep implements step.okta_factor_enroll
type factorEnrollStep struct {
	name       string
	moduleName string
}

func newFactorEnrollStep(name string, config map[string]any) (*factorEnrollStep, error) {
	return &factorEnrollStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *factorEnrollStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	factorType := resolveValue("factorType", current, config)
	provider := resolveValue("provider", current, config)
	if userID == "" || factorType == "" {
		return &sdk.StepResult{Output: errResult("userId and factorType are required")}, nil
	}
	body := map[string]any{
		"factorType": factorType,
		"provider":   provider,
	}
	if profile := resolveMap("profile", current, config); profile != nil {
		body["profile"] = profile
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/factors", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// factorListStep implements step.okta_factor_list
type factorListStep struct {
	name       string
	moduleName string
}

func newFactorListStep(name string, config map[string]any) (*factorListStep, error) {
	return &factorListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *factorListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/users/"+userID+"/factors", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("factors", items)}, nil
}

// factorVerifyStep implements step.okta_factor_verify
type factorVerifyStep struct {
	name       string
	moduleName string
}

func newFactorVerifyStep(name string, config map[string]any) (*factorVerifyStep, error) {
	return &factorVerifyStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *factorVerifyStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	factorID := resolveValue("factorId", current, config)
	if userID == "" || factorID == "" {
		return &sdk.StepResult{Output: errResult("userId and factorId are required")}, nil
	}
	body := map[string]any{}
	if passCode := resolveValue("passCode", current, config); passCode != "" {
		body["passCode"] = passCode
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/factors/"+factorID+"/verify", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// factorUnenrollStep implements step.okta_factor_unenroll
type factorUnenrollStep struct {
	name       string
	moduleName string
}

func newFactorUnenrollStep(name string, config map[string]any) (*factorUnenrollStep, error) {
	return &factorUnenrollStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *factorUnenrollStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	factorID := resolveValue("factorId", current, config)
	if userID == "" || factorID == "" {
		return &sdk.StepResult{Output: errResult("userId and factorId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/users/"+userID+"/factors/"+factorID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"unenrolled": true}}, nil
}

// factorActivateStep implements step.okta_factor_activate
type factorActivateStep struct {
	name       string
	moduleName string
}

func newFactorActivateStep(name string, config map[string]any) (*factorActivateStep, error) {
	return &factorActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *factorActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	factorID := resolveValue("factorId", current, config)
	if userID == "" || factorID == "" {
		return &sdk.StepResult{Output: errResult("userId and factorId are required")}, nil
	}
	body := map[string]any{}
	if passCode := resolveValue("passCode", current, config); passCode != "" {
		body["passCode"] = passCode
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/factors/"+factorID+"/lifecycle/activate", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}
