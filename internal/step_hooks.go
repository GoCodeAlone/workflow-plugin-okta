package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// eventHookCreateStep implements step.okta_event_hook_create
type eventHookCreateStep struct {
	name       string
	moduleName string
}

func newEventHookCreateStep(name string, config map[string]any) (*eventHookCreateStep, error) {
	return &eventHookCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookName := resolveValue("name", current, config)
	if hookName == "" {
		return &sdk.StepResult{Output: errResult("name is required")}, nil
	}
	body := map[string]any{
		"name":    hookName,
		"events":  resolveMap("events", current, config),
		"channel": resolveMap("channel", current, config),
	}
	result, err := oktaPost(client, "/api/v1/eventHooks", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// eventHookGetStep implements step.okta_event_hook_get
type eventHookGetStep struct {
	name       string
	moduleName string
}

func newEventHookGetStep(name string, config map[string]any) (*eventHookGetStep, error) {
	return &eventHookGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("eventHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("eventHookId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/eventHooks/"+hookID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// eventHookListStep implements step.okta_event_hook_list
type eventHookListStep struct {
	name       string
	moduleName string
}

func newEventHookListStep(name string, config map[string]any) (*eventHookListStep, error) {
	return &eventHookListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/eventHooks", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("eventHooks", items)}, nil
}

// eventHookDeleteStep implements step.okta_event_hook_delete
type eventHookDeleteStep struct {
	name       string
	moduleName string
}

func newEventHookDeleteStep(name string, config map[string]any) (*eventHookDeleteStep, error) {
	return &eventHookDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("eventHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("eventHookId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/eventHooks/"+hookID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// eventHookActivateStep implements step.okta_event_hook_activate
type eventHookActivateStep struct {
	name       string
	moduleName string
}

func newEventHookActivateStep(name string, config map[string]any) (*eventHookActivateStep, error) {
	return &eventHookActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("eventHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("eventHookId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/eventHooks/"+hookID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// eventHookDeactivateStep implements step.okta_event_hook_deactivate
type eventHookDeactivateStep struct {
	name       string
	moduleName string
}

func newEventHookDeactivateStep(name string, config map[string]any) (*eventHookDeactivateStep, error) {
	return &eventHookDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("eventHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("eventHookId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/eventHooks/"+hookID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// eventHookVerifyStep implements step.okta_event_hook_verify
type eventHookVerifyStep struct {
	name       string
	moduleName string
}

func newEventHookVerifyStep(name string, config map[string]any) (*eventHookVerifyStep, error) {
	return &eventHookVerifyStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *eventHookVerifyStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("eventHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("eventHookId is required")}, nil
	}
	result, err := oktaPost(client, "/api/v1/eventHooks/"+hookID+"/lifecycle/verify", map[string]any{})
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"verified": true}}, nil
}

// inlineHookCreateStep implements step.okta_inline_hook_create
type inlineHookCreateStep struct {
	name       string
	moduleName string
}

func newInlineHookCreateStep(name string, config map[string]any) (*inlineHookCreateStep, error) {
	return &inlineHookCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookName := resolveValue("name", current, config)
	hookType := resolveValue("type", current, config)
	if hookName == "" || hookType == "" {
		return &sdk.StepResult{Output: errResult("name and type are required")}, nil
	}
	body := map[string]any{
		"name":    hookName,
		"type":    hookType,
		"version": resolveValue("version", current, config),
		"channel": resolveMap("channel", current, config),
	}
	result, err := oktaPost(client, "/api/v1/inlineHooks", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// inlineHookGetStep implements step.okta_inline_hook_get
type inlineHookGetStep struct {
	name       string
	moduleName string
}

func newInlineHookGetStep(name string, config map[string]any) (*inlineHookGetStep, error) {
	return &inlineHookGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("inlineHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("inlineHookId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/inlineHooks/"+hookID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// inlineHookListStep implements step.okta_inline_hook_list
type inlineHookListStep struct {
	name       string
	moduleName string
}

func newInlineHookListStep(name string, config map[string]any) (*inlineHookListStep, error) {
	return &inlineHookListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/inlineHooks", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("inlineHooks", items)}, nil
}

// inlineHookDeleteStep implements step.okta_inline_hook_delete
type inlineHookDeleteStep struct {
	name       string
	moduleName string
}

func newInlineHookDeleteStep(name string, config map[string]any) (*inlineHookDeleteStep, error) {
	return &inlineHookDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("inlineHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("inlineHookId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/inlineHooks/"+hookID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// inlineHookActivateStep implements step.okta_inline_hook_activate
type inlineHookActivateStep struct {
	name       string
	moduleName string
}

func newInlineHookActivateStep(name string, config map[string]any) (*inlineHookActivateStep, error) {
	return &inlineHookActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("inlineHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("inlineHookId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/inlineHooks/"+hookID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// inlineHookDeactivateStep implements step.okta_inline_hook_deactivate
type inlineHookDeactivateStep struct {
	name       string
	moduleName string
}

func newInlineHookDeactivateStep(name string, config map[string]any) (*inlineHookDeactivateStep, error) {
	return &inlineHookDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("inlineHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("inlineHookId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/inlineHooks/"+hookID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// inlineHookExecuteStep implements step.okta_inline_hook_execute
type inlineHookExecuteStep struct {
	name       string
	moduleName string
}

func newInlineHookExecuteStep(name string, config map[string]any) (*inlineHookExecuteStep, error) {
	return &inlineHookExecuteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *inlineHookExecuteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	hookID := resolveValue("inlineHookId", current, config)
	if hookID == "" {
		return &sdk.StepResult{Output: errResult("inlineHookId is required")}, nil
	}
	body := map[string]any{}
	if payload := resolveMap("payload", current, config); payload != nil {
		body = payload
	}
	result, err := oktaPost(client, "/api/v1/inlineHooks/"+hookID+"/execute", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}
