package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// groupCreateStep implements step.okta_group_create
type groupCreateStep struct {
	name       string
	moduleName string
}

func newGroupCreateStep(name string, config map[string]any) (*groupCreateStep, error) {
	return &groupCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupName := resolveValue("name", current, config)
	if groupName == "" {
		return &sdk.StepResult{Output: errResult("name is required")}, nil
	}
	body := map[string]any{
		"profile": map[string]any{
			"name":        groupName,
			"description": resolveValue("description", current, config),
		},
	}
	result, err := oktaPost(client, "/api/v1/groups", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// groupGetStep implements step.okta_group_get
type groupGetStep struct {
	name       string
	moduleName string
}

func newGroupGetStep(name string, config map[string]any) (*groupGetStep, error) {
	return &groupGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupID := resolveValue("groupId", current, config)
	if groupID == "" {
		return &sdk.StepResult{Output: errResult("groupId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/groups/"+groupID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// groupListStep implements step.okta_group_list
type groupListStep struct {
	name       string
	moduleName string
}

func newGroupListStep(name string, config map[string]any) (*groupListStep, error) {
	return &groupListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	qp := url.Values{}
	if q := resolveValue("q", current, config); q != "" {
		qp.Set("q", q)
	}
	if limit := resolveInt("limit", current, config); limit > 0 {
		qp.Set("limit", fmt.Sprintf("%d", limit))
	}
	result, err := oktaGet(client, "/api/v1/groups", qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("groups", items)}, nil
}

// groupDeleteStep implements step.okta_group_delete
type groupDeleteStep struct {
	name       string
	moduleName string
}

func newGroupDeleteStep(name string, config map[string]any) (*groupDeleteStep, error) {
	return &groupDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupID := resolveValue("groupId", current, config)
	if groupID == "" {
		return &sdk.StepResult{Output: errResult("groupId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/groups/"+groupID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// groupAddUserStep implements step.okta_group_add_user
type groupAddUserStep struct {
	name       string
	moduleName string
}

func newGroupAddUserStep(name string, config map[string]any) (*groupAddUserStep, error) {
	return &groupAddUserStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupAddUserStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupID := resolveValue("groupId", current, config)
	userID := resolveValue("userId", current, config)
	if groupID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("groupId and userId are required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/groups/"+groupID+"/users/"+userID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"added": true}}, nil
}

// groupRemoveUserStep implements step.okta_group_remove_user
type groupRemoveUserStep struct {
	name       string
	moduleName string
}

func newGroupRemoveUserStep(name string, config map[string]any) (*groupRemoveUserStep, error) {
	return &groupRemoveUserStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRemoveUserStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupID := resolveValue("groupId", current, config)
	userID := resolveValue("userId", current, config)
	if groupID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("groupId and userId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/groups/"+groupID+"/users/"+userID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"removed": true}}, nil
}

// groupListUsersStep implements step.okta_group_list_users
type groupListUsersStep struct {
	name       string
	moduleName string
}

func newGroupListUsersStep(name string, config map[string]any) (*groupListUsersStep, error) {
	return &groupListUsersStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupListUsersStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	groupID := resolveValue("groupId", current, config)
	if groupID == "" {
		return &sdk.StepResult{Output: errResult("groupId is required")}, nil
	}
	qp := url.Values{}
	if limit := resolveInt("limit", current, config); limit > 0 {
		qp.Set("limit", fmt.Sprintf("%d", limit))
	}
	result, err := oktaGet(client, "/api/v1/groups/"+groupID+"/users", qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("users", items)}, nil
}

// groupRuleCreateStep implements step.okta_group_rule_create
type groupRuleCreateStep struct {
	name       string
	moduleName string
}

func newGroupRuleCreateStep(name string, config map[string]any) (*groupRuleCreateStep, error) {
	return &groupRuleCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	ruleName := resolveValue("name", current, config)
	if ruleName == "" {
		return &sdk.StepResult{Output: errResult("name is required")}, nil
	}
	body := map[string]any{
		"type": "group_rule",
		"name": ruleName,
	}
	if conditions := resolveMap("conditions", current, config); conditions != nil {
		body["conditions"] = conditions
	}
	if actions := resolveMap("actions", current, config); actions != nil {
		body["actions"] = actions
	}
	result, err := oktaPost(client, "/api/v1/groups/rules", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// groupRuleGetStep implements step.okta_group_rule_get
type groupRuleGetStep struct {
	name       string
	moduleName string
}

func newGroupRuleGetStep(name string, config map[string]any) (*groupRuleGetStep, error) {
	return &groupRuleGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	ruleID := resolveValue("ruleId", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: errResult("ruleId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/groups/rules/"+ruleID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// groupRuleListStep implements step.okta_group_rule_list
type groupRuleListStep struct {
	name       string
	moduleName string
}

func newGroupRuleListStep(name string, config map[string]any) (*groupRuleListStep, error) {
	return &groupRuleListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/groups/rules", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("rules", items)}, nil
}

// groupRuleDeleteStep implements step.okta_group_rule_delete
type groupRuleDeleteStep struct {
	name       string
	moduleName string
}

func newGroupRuleDeleteStep(name string, config map[string]any) (*groupRuleDeleteStep, error) {
	return &groupRuleDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	ruleID := resolveValue("ruleId", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: errResult("ruleId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/groups/rules/"+ruleID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// groupRuleActivateStep implements step.okta_group_rule_activate
type groupRuleActivateStep struct {
	name       string
	moduleName string
}

func newGroupRuleActivateStep(name string, config map[string]any) (*groupRuleActivateStep, error) {
	return &groupRuleActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	ruleID := resolveValue("ruleId", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: errResult("ruleId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/groups/rules/"+ruleID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// groupRuleDeactivateStep implements step.okta_group_rule_deactivate
type groupRuleDeactivateStep struct {
	name       string
	moduleName string
}

func newGroupRuleDeactivateStep(name string, config map[string]any) (*groupRuleDeactivateStep, error) {
	return &groupRuleDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *groupRuleDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	ruleID := resolveValue("ruleId", current, config)
	if ruleID == "" {
		return &sdk.StepResult{Output: errResult("ruleId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/groups/rules/"+ruleID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}
