package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// sessionGetStep implements step.okta_session_get
type sessionGetStep struct {
	name       string
	moduleName string
}

func newSessionGetStep(name string, config map[string]any) (*sessionGetStep, error) {
	return &sessionGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sessionGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	sessionID := resolveValue("sessionId", current, config)
	if sessionID == "" {
		return &sdk.StepResult{Output: errResult("sessionId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/sessions/"+sessionID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// sessionRefreshStep implements step.okta_session_refresh
type sessionRefreshStep struct {
	name       string
	moduleName string
}

func newSessionRefreshStep(name string, config map[string]any) (*sessionRefreshStep, error) {
	return &sessionRefreshStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sessionRefreshStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	sessionID := resolveValue("sessionId", current, config)
	if sessionID == "" {
		return &sdk.StepResult{Output: errResult("sessionId is required")}, nil
	}
	result, err := oktaPost(client, "/api/v1/sessions/"+sessionID+"/lifecycle/refresh", map[string]any{})
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"refreshed": true}}, nil
}

// sessionRevokeStep implements step.okta_session_revoke
type sessionRevokeStep struct {
	name       string
	moduleName string
}

func newSessionRevokeStep(name string, config map[string]any) (*sessionRevokeStep, error) {
	return &sessionRevokeStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sessionRevokeStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	sessionID := resolveValue("sessionId", current, config)
	if sessionID == "" {
		return &sdk.StepResult{Output: errResult("sessionId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/sessions/"+sessionID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"revoked": true}}, nil
}
