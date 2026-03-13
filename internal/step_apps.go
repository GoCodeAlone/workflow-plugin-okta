package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// appCreateStep implements step.okta_app_create
type appCreateStep struct {
	name       string
	moduleName string
}

func newAppCreateStep(name string, config map[string]any) (*appCreateStep, error) {
	return &appCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appName := resolveValue("name", current, config)
	appLabel := resolveValue("label", current, config)
	if appLabel == "" {
		appLabel = appName
	}
	signOnMode := resolveValue("signOnMode", current, config)
	if signOnMode == "" {
		signOnMode = "AUTO"
	}
	body := map[string]any{
		"name":       appName,
		"label":      appLabel,
		"signOnMode": signOnMode,
	}
	if settings := resolveMap("settings", current, config); settings != nil {
		body["settings"] = settings
	}
	if visibility := resolveMap("visibility", current, config); visibility != nil {
		body["visibility"] = visibility
	}
	result, err := oktaPost(client, "/api/v1/apps", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appGetStep implements step.okta_app_get
type appGetStep struct {
	name       string
	moduleName string
}

func newAppGetStep(name string, config map[string]any) (*appGetStep, error) {
	return &appGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/apps/"+appID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appListStep implements step.okta_app_list
type appListStep struct {
	name       string
	moduleName string
}

func newAppListStep(name string, config map[string]any) (*appListStep, error) {
	return &appListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
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
	result, err := oktaGet(client, "/api/v1/apps", qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("apps", items)}, nil
}

// appUpdateStep implements step.okta_app_update
type appUpdateStep struct {
	name       string
	moduleName string
}

func newAppUpdateStep(name string, config map[string]any) (*appUpdateStep, error) {
	return &appUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	body := map[string]any{}
	if label := resolveValue("label", current, config); label != "" {
		body["label"] = label
	}
	if settings := resolveMap("settings", current, config); settings != nil {
		body["settings"] = settings
	}
	result, err := oktaPut(client, "/api/v1/apps/"+appID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appDeleteStep implements step.okta_app_delete
type appDeleteStep struct {
	name       string
	moduleName string
}

func newAppDeleteStep(name string, config map[string]any) (*appDeleteStep, error) {
	return &appDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/apps/"+appID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// appActivateStep implements step.okta_app_activate
type appActivateStep struct {
	name       string
	moduleName string
}

func newAppActivateStep(name string, config map[string]any) (*appActivateStep, error) {
	return &appActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/apps/"+appID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// appDeactivateStep implements step.okta_app_deactivate
type appDeactivateStep struct {
	name       string
	moduleName string
}

func newAppDeactivateStep(name string, config map[string]any) (*appDeactivateStep, error) {
	return &appDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/apps/"+appID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// appUserAssignStep implements step.okta_app_user_assign
type appUserAssignStep struct {
	name       string
	moduleName string
}

func newAppUserAssignStep(name string, config map[string]any) (*appUserAssignStep, error) {
	return &appUserAssignStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUserAssignStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	userID := resolveValue("userId", current, config)
	if appID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("appId and userId are required")}, nil
	}
	body := map[string]any{"id": userID}
	if profile := resolveMap("profile", current, config); profile != nil {
		body["profile"] = profile
	}
	result, err := oktaPost(client, "/api/v1/apps/"+appID+"/users", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"assigned": true}}, nil
}

// appUserGetStep implements step.okta_app_user_get
type appUserGetStep struct {
	name       string
	moduleName string
}

func newAppUserGetStep(name string, config map[string]any) (*appUserGetStep, error) {
	return &appUserGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUserGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	userID := resolveValue("userId", current, config)
	if appID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("appId and userId are required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/apps/"+appID+"/users/"+userID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appUserListStep implements step.okta_app_user_list
type appUserListStep struct {
	name       string
	moduleName string
}

func newAppUserListStep(name string, config map[string]any) (*appUserListStep, error) {
	return &appUserListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUserListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/apps/"+appID+"/users", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("users", items)}, nil
}

// appUserUpdateStep implements step.okta_app_user_update
type appUserUpdateStep struct {
	name       string
	moduleName string
}

func newAppUserUpdateStep(name string, config map[string]any) (*appUserUpdateStep, error) {
	return &appUserUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUserUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	userID := resolveValue("userId", current, config)
	if appID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("appId and userId are required")}, nil
	}
	body := map[string]any{}
	if profile := resolveMap("profile", current, config); profile != nil {
		body["profile"] = profile
	}
	result, err := oktaPost(client, "/api/v1/apps/"+appID+"/users/"+userID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appUserUnassignStep implements step.okta_app_user_unassign
type appUserUnassignStep struct {
	name       string
	moduleName string
}

func newAppUserUnassignStep(name string, config map[string]any) (*appUserUnassignStep, error) {
	return &appUserUnassignStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appUserUnassignStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	userID := resolveValue("userId", current, config)
	if appID == "" || userID == "" {
		return &sdk.StepResult{Output: errResult("appId and userId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/apps/"+appID+"/users/"+userID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"unassigned": true}}, nil
}

// appGroupAssignStep implements step.okta_app_group_assign
type appGroupAssignStep struct {
	name       string
	moduleName string
}

func newAppGroupAssignStep(name string, config map[string]any) (*appGroupAssignStep, error) {
	return &appGroupAssignStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGroupAssignStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	groupID := resolveValue("groupId", current, config)
	if appID == "" || groupID == "" {
		return &sdk.StepResult{Output: errResult("appId and groupId are required")}, nil
	}
	body := map[string]any{}
	if priority := resolveInt("priority", current, config); priority > 0 {
		body["priority"] = priority
	}
	result, err := oktaPut(client, "/api/v1/apps/"+appID+"/groups/"+groupID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"assigned": true}}, nil
}

// appGroupGetStep implements step.okta_app_group_get
type appGroupGetStep struct {
	name       string
	moduleName string
}

func newAppGroupGetStep(name string, config map[string]any) (*appGroupGetStep, error) {
	return &appGroupGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGroupGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	groupID := resolveValue("groupId", current, config)
	if appID == "" || groupID == "" {
		return &sdk.StepResult{Output: errResult("appId and groupId are required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/apps/"+appID+"/groups/"+groupID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// appGroupListStep implements step.okta_app_group_list
type appGroupListStep struct {
	name       string
	moduleName string
}

func newAppGroupListStep(name string, config map[string]any) (*appGroupListStep, error) {
	return &appGroupListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGroupListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	if appID == "" {
		return &sdk.StepResult{Output: errResult("appId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/apps/"+appID+"/groups", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("groups", items)}, nil
}

// appGroupUpdateStep implements step.okta_app_group_update
type appGroupUpdateStep struct {
	name       string
	moduleName string
}

func newAppGroupUpdateStep(name string, config map[string]any) (*appGroupUpdateStep, error) {
	return &appGroupUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGroupUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	groupID := resolveValue("groupId", current, config)
	if appID == "" || groupID == "" {
		return &sdk.StepResult{Output: errResult("appId and groupId are required")}, nil
	}
	body := map[string]any{}
	if priority := resolveInt("priority", current, config); priority > 0 {
		body["priority"] = priority
	}
	result, err := oktaPut(client, "/api/v1/apps/"+appID+"/groups/"+groupID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true}}, nil
}

// appGroupUnassignStep implements step.okta_app_group_unassign
type appGroupUnassignStep struct {
	name       string
	moduleName string
}

func newAppGroupUnassignStep(name string, config map[string]any) (*appGroupUnassignStep, error) {
	return &appGroupUnassignStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *appGroupUnassignStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	appID := resolveValue("appId", current, config)
	groupID := resolveValue("groupId", current, config)
	if appID == "" || groupID == "" {
		return &sdk.StepResult{Output: errResult("appId and groupId are required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/apps/"+appID+"/groups/"+groupID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"unassigned": true}}, nil
}
