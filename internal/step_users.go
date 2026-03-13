package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// userCreateStep implements step.okta_user_create
type userCreateStep struct {
	name       string
	moduleName string
}

func newUserCreateStep(name string, config map[string]any) (*userCreateStep, error) {
	return &userCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	body := map[string]any{}
	if profile := resolveMap("profile", current, config); profile != nil {
		body["profile"] = profile
	} else {
		firstName := resolveValue("firstName", current, config)
		lastName := resolveValue("lastName", current, config)
		email := resolveValue("email", current, config)
		login := resolveValue("login", current, config)
		if login == "" {
			login = email
		}
		if email == "" {
			return &sdk.StepResult{Output: errResult("email is required")}, nil
		}
		body["profile"] = map[string]any{
			"firstName": firstName,
			"lastName":  lastName,
			"email":     email,
			"login":     login,
		}
	}
	if creds := resolveMap("credentials", current, config); creds != nil {
		body["credentials"] = creds
	}
	activate := resolveBool("activate", current, config)
	qp := url.Values{}
	if activate {
		qp.Set("activate", "true")
	}
	result, _, err := oktaRequest(client, "POST", "/api/v1/users", body, qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// userGetStep implements step.okta_user_get
type userGetStep struct {
	name       string
	moduleName string
}

func newUserGetStep(name string, config map[string]any) (*userGetStep, error) {
	return &userGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/users/"+userID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// userListStep implements step.okta_user_list
type userListStep struct {
	name       string
	moduleName string
}

func newUserListStep(name string, config map[string]any) (*userListStep, error) {
	return &userListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	qp := url.Values{}
	if q := resolveValue("q", current, config); q != "" {
		qp.Set("q", q)
	}
	if filter := resolveValue("filter", current, config); filter != "" {
		qp.Set("filter", filter)
	}
	if limit := resolveInt("limit", current, config); limit > 0 {
		qp.Set("limit", fmt.Sprintf("%d", limit))
	}
	result, err := oktaGet(client, "/api/v1/users", qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("users", items)}, nil
}

// userUpdateStep implements step.okta_user_update
type userUpdateStep struct {
	name       string
	moduleName string
}

func newUserUpdateStep(name string, config map[string]any) (*userUpdateStep, error) {
	return &userUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	body := map[string]any{}
	if profile := resolveMap("profile", current, config); profile != nil {
		body["profile"] = profile
	}
	if creds := resolveMap("credentials", current, config); creds != nil {
		body["credentials"] = creds
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// userDeleteStep implements step.okta_user_delete
type userDeleteStep struct {
	name       string
	moduleName string
}

func newUserDeleteStep(name string, config map[string]any) (*userDeleteStep, error) {
	return &userDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	// First deactivate, then delete
	if err := oktaDelete(client, "/api/v1/users/"+userID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "userId": userID}}, nil
}

// userActivateStep implements step.okta_user_activate
type userActivateStep struct {
	name       string
	moduleName string
}

func newUserActivateStep(name string, config map[string]any) (*userActivateStep, error) {
	return &userActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	sendEmail := resolveBool("sendEmail", current, config)
	qp := url.Values{}
	qp.Set("sendEmail", fmt.Sprintf("%v", sendEmail))
	result, _, err := oktaRequest(client, "POST", "/api/v1/users/"+userID+"/lifecycle/activate", map[string]any{}, qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// userDeactivateStep implements step.okta_user_deactivate
type userDeactivateStep struct {
	name       string
	moduleName string
}

func newUserDeactivateStep(name string, config map[string]any) (*userDeactivateStep, error) {
	return &userDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/users/"+userID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}

// userReactivateStep implements step.okta_user_reactivate
type userReactivateStep struct {
	name       string
	moduleName string
}

func newUserReactivateStep(name string, config map[string]any) (*userReactivateStep, error) {
	return &userReactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userReactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	result, _, err := oktaRequest(client, "POST", "/api/v1/users/"+userID+"/lifecycle/reactivate", map[string]any{}, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"reactivated": true}}, nil
}

// userSuspendStep implements step.okta_user_suspend
type userSuspendStep struct {
	name       string
	moduleName string
}

func newUserSuspendStep(name string, config map[string]any) (*userSuspendStep, error) {
	return &userSuspendStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userSuspendStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/users/"+userID+"/lifecycle/suspend"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"suspended": true}}, nil
}

// userUnsuspendStep implements step.okta_user_unsuspend
type userUnsuspendStep struct {
	name       string
	moduleName string
}

func newUserUnsuspendStep(name string, config map[string]any) (*userUnsuspendStep, error) {
	return &userUnsuspendStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userUnsuspendStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/users/"+userID+"/lifecycle/unsuspend"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"unsuspended": true}}, nil
}

// userUnlockStep implements step.okta_user_unlock
type userUnlockStep struct {
	name       string
	moduleName string
}

func newUserUnlockStep(name string, config map[string]any) (*userUnlockStep, error) {
	return &userUnlockStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userUnlockStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/users/"+userID+"/lifecycle/unlock"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"unlocked": true}}, nil
}

// userResetFactorsStep implements step.okta_user_reset_factors
type userResetFactorsStep struct {
	name       string
	moduleName string
}

func newUserResetFactorsStep(name string, config map[string]any) (*userResetFactorsStep, error) {
	return &userResetFactorsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userResetFactorsStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/users/"+userID+"/lifecycle/reset_factors"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"reset": true}}, nil
}

// userChangePasswordStep implements step.okta_user_change_password
type userChangePasswordStep struct {
	name       string
	moduleName string
}

func newUserChangePasswordStep(name string, config map[string]any) (*userChangePasswordStep, error) {
	return &userChangePasswordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userChangePasswordStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	oldPassword := resolveValue("oldPassword", current, config)
	newPassword := resolveValue("newPassword", current, config)
	if newPassword == "" {
		return &sdk.StepResult{Output: errResult("newPassword is required")}, nil
	}
	body := map[string]any{
		"oldPassword": map[string]any{"value": oldPassword},
		"newPassword": map[string]any{"value": newPassword},
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/credentials/change_password", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"changed": true}}, nil
}

// userResetPasswordStep implements step.okta_user_reset_password
type userResetPasswordStep struct {
	name       string
	moduleName string
}

func newUserResetPasswordStep(name string, config map[string]any) (*userResetPasswordStep, error) {
	return &userResetPasswordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userResetPasswordStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	sendEmail := resolveBool("sendEmail", current, config)
	qp := url.Values{}
	qp.Set("sendEmail", fmt.Sprintf("%v", sendEmail))
	result, _, err := oktaRequest(client, "POST", "/api/v1/users/"+userID+"/lifecycle/reset_password", map[string]any{}, qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"reset": true}}, nil
}

// userExpirePasswordStep implements step.okta_user_expire_password
type userExpirePasswordStep struct {
	name       string
	moduleName string
}

func newUserExpirePasswordStep(name string, config map[string]any) (*userExpirePasswordStep, error) {
	return &userExpirePasswordStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userExpirePasswordStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/lifecycle/expire_password", map[string]any{})
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"expired": true}}, nil
}

// userSetRecoveryQuestionStep implements step.okta_user_set_recovery_question
type userSetRecoveryQuestionStep struct {
	name       string
	moduleName string
}

func newUserSetRecoveryQuestionStep(name string, config map[string]any) (*userSetRecoveryQuestionStep, error) {
	return &userSetRecoveryQuestionStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *userSetRecoveryQuestionStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	userID := resolveValue("userId", current, config)
	if userID == "" {
		return &sdk.StepResult{Output: errResult("userId is required")}, nil
	}
	question := resolveValue("question", current, config)
	answer := resolveValue("answer", current, config)
	password := resolveValue("password", current, config)
	body := map[string]any{
		"password": map[string]any{"value": password},
		"recovery_question": map[string]any{
			"question": question,
			"answer":   answer,
		},
	}
	result, err := oktaPost(client, "/api/v1/users/"+userID+"/credentials/change_recovery_question", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true}}, nil
}
