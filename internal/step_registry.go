package internal

import (
	"fmt"
	"sort"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// stepConstructor is a function that creates a StepInstance.
type stepConstructor func(name string, config map[string]any) (sdk.StepInstance, error)

// stepRegistry maps step type strings to constructor functions.
var stepRegistry = map[string]stepConstructor{
	// Provider descriptors
	"step.okta_auth_provider_describe": func(n string, c map[string]any) (sdk.StepInstance, error) {
		return newAuthProviderDescribeStep(n, c), nil
	},

	// Users — CRUD
	"step.okta_user_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserCreateStep(n, c) },
	"step.okta_user_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserGetStep(n, c) },
	"step.okta_user_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserListStep(n, c) },
	"step.okta_user_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserUpdateStep(n, c) },
	"step.okta_user_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserDeleteStep(n, c) },

	// Users — Lifecycle
	"step.okta_user_activate":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserActivateStep(n, c) },
	"step.okta_user_deactivate":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserDeactivateStep(n, c) },
	"step.okta_user_reactivate":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserReactivateStep(n, c) },
	"step.okta_user_suspend":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserSuspendStep(n, c) },
	"step.okta_user_unsuspend":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserUnsuspendStep(n, c) },
	"step.okta_user_unlock":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserUnlockStep(n, c) },
	"step.okta_user_reset_factors": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserResetFactorsStep(n, c) },

	// Users — Credentials
	"step.okta_user_change_password": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserChangePasswordStep(n, c) },
	"step.okta_user_reset_password":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserResetPasswordStep(n, c) },
	"step.okta_user_expire_password": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserExpirePasswordStep(n, c) },
	"step.okta_user_set_recovery_question": func(n string, c map[string]any) (sdk.StepInstance, error) {
		return newUserSetRecoveryQuestionStep(n, c)
	},

	// Groups — CRUD
	"step.okta_group_create":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupCreateStep(n, c) },
	"step.okta_group_get":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupGetStep(n, c) },
	"step.okta_group_list":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupListStep(n, c) },
	"step.okta_group_delete":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupDeleteStep(n, c) },
	"step.okta_group_add_user":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupAddUserStep(n, c) },
	"step.okta_group_remove_user": func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRemoveUserStep(n, c) },
	"step.okta_group_list_users":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupListUsersStep(n, c) },

	// Group Rules
	"step.okta_group_rule_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleCreateStep(n, c) },
	"step.okta_group_rule_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleGetStep(n, c) },
	"step.okta_group_rule_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleListStep(n, c) },
	"step.okta_group_rule_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleDeleteStep(n, c) },
	"step.okta_group_rule_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleActivateStep(n, c) },
	"step.okta_group_rule_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newGroupRuleDeactivateStep(n, c) },

	// Applications — Core
	"step.okta_app_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppCreateStep(n, c) },
	"step.okta_app_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGetStep(n, c) },
	"step.okta_app_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppListStep(n, c) },
	"step.okta_app_update":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUpdateStep(n, c) },
	"step.okta_app_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppDeleteStep(n, c) },
	"step.okta_app_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppActivateStep(n, c) },
	"step.okta_app_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppDeactivateStep(n, c) },

	// Applications — Users
	"step.okta_app_user_assign":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUserAssignStep(n, c) },
	"step.okta_app_user_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUserGetStep(n, c) },
	"step.okta_app_user_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUserListStep(n, c) },
	"step.okta_app_user_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUserUpdateStep(n, c) },
	"step.okta_app_user_unassign": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppUserUnassignStep(n, c) },

	// Applications — Groups
	"step.okta_app_group_assign":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGroupAssignStep(n, c) },
	"step.okta_app_group_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGroupGetStep(n, c) },
	"step.okta_app_group_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGroupListStep(n, c) },
	"step.okta_app_group_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGroupUpdateStep(n, c) },
	"step.okta_app_group_unassign": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAppGroupUnassignStep(n, c) },

	// Authorization Servers
	"step.okta_authz_server_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerCreateStep(n, c) },
	"step.okta_authz_server_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerGetStep(n, c) },
	"step.okta_authz_server_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerListStep(n, c) },
	"step.okta_authz_server_update":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerUpdateStep(n, c) },
	"step.okta_authz_server_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerDeleteStep(n, c) },
	"step.okta_authz_server_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerActivateStep(n, c) },
	"step.okta_authz_server_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzServerDeactivateStep(n, c) },

	// Auth Server — Claims
	"step.okta_authz_claim_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzClaimCreateStep(n, c) },
	"step.okta_authz_claim_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzClaimListStep(n, c) },
	"step.okta_authz_claim_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzClaimDeleteStep(n, c) },

	// Auth Server — Scopes
	"step.okta_authz_scope_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzScopeCreateStep(n, c) },
	"step.okta_authz_scope_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzScopeListStep(n, c) },
	"step.okta_authz_scope_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzScopeDeleteStep(n, c) },

	// Auth Server — Policies
	"step.okta_authz_policy_create":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyCreateStep(n, c) },
	"step.okta_authz_policy_list":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyListStep(n, c) },
	"step.okta_authz_policy_delete":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyDeleteStep(n, c) },
	"step.okta_authz_policy_rule_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyRuleCreateStep(n, c) },
	"step.okta_authz_policy_rule_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyRuleListStep(n, c) },
	"step.okta_authz_policy_rule_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzPolicyRuleDeleteStep(n, c) },
	"step.okta_authz_key_list":           func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzKeyListStep(n, c) },
	"step.okta_authz_key_rotate":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthzKeyRotateStep(n, c) },

	// Policies
	"step.okta_policy_create":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyCreateStep(n, c) },
	"step.okta_policy_get":             func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyGetStep(n, c) },
	"step.okta_policy_list":            func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyListStep(n, c) },
	"step.okta_policy_delete":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyDeleteStep(n, c) },
	"step.okta_policy_activate":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyActivateStep(n, c) },
	"step.okta_policy_deactivate":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyDeactivateStep(n, c) },
	"step.okta_policy_rule_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyRuleCreateStep(n, c) },
	"step.okta_policy_rule_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyRuleListStep(n, c) },
	"step.okta_policy_rule_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyRuleDeleteStep(n, c) },
	"step.okta_policy_rule_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyRuleActivateStep(n, c) },
	"step.okta_policy_rule_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newPolicyRuleDeactivateStep(n, c) },

	// Authenticators (MFA)
	"step.okta_authenticator_create":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthenticatorCreateStep(n, c) },
	"step.okta_authenticator_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthenticatorGetStep(n, c) },
	"step.okta_authenticator_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthenticatorListStep(n, c) },
	"step.okta_authenticator_activate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newAuthenticatorActivateStep(n, c) },
	"step.okta_authenticator_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) {
		return newAuthenticatorDeactivateStep(n, c)
	},

	// User Factors
	"step.okta_factor_enroll":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newFactorEnrollStep(n, c) },
	"step.okta_factor_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newFactorListStep(n, c) },
	"step.okta_factor_verify":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newFactorVerifyStep(n, c) },
	"step.okta_factor_unenroll": func(n string, c map[string]any) (sdk.StepInstance, error) { return newFactorUnenrollStep(n, c) },
	"step.okta_factor_activate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newFactorActivateStep(n, c) },

	// Identity Providers
	"step.okta_idp_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpCreateStep(n, c) },
	"step.okta_idp_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpGetStep(n, c) },
	"step.okta_idp_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpListStep(n, c) },
	"step.okta_idp_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpDeleteStep(n, c) },
	"step.okta_idp_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpActivateStep(n, c) },
	"step.okta_idp_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdpDeactivateStep(n, c) },

	// Sessions
	"step.okta_session_get":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newSessionGetStep(n, c) },
	"step.okta_session_refresh": func(n string, c map[string]any) (sdk.StepInstance, error) { return newSessionRefreshStep(n, c) },
	"step.okta_session_revoke":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newSessionRevokeStep(n, c) },

	// Network Zones
	"step.okta_network_zone_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneCreateStep(n, c) },
	"step.okta_network_zone_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneGetStep(n, c) },
	"step.okta_network_zone_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneListStep(n, c) },
	"step.okta_network_zone_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneDeleteStep(n, c) },
	"step.okta_network_zone_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneActivateStep(n, c) },
	"step.okta_network_zone_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newNetworkZoneDeactivateStep(n, c) },

	// System Log
	"step.okta_log_list": func(n string, c map[string]any) (sdk.StepInstance, error) { return newLogListStep(n, c) },

	// Event Hooks
	"step.okta_event_hook_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookCreateStep(n, c) },
	"step.okta_event_hook_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookGetStep(n, c) },
	"step.okta_event_hook_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookListStep(n, c) },
	"step.okta_event_hook_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookDeleteStep(n, c) },
	"step.okta_event_hook_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookActivateStep(n, c) },
	"step.okta_event_hook_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookDeactivateStep(n, c) },
	"step.okta_event_hook_verify":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventHookVerifyStep(n, c) },

	// Inline Hooks
	"step.okta_inline_hook_create":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookCreateStep(n, c) },
	"step.okta_inline_hook_get":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookGetStep(n, c) },
	"step.okta_inline_hook_list":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookListStep(n, c) },
	"step.okta_inline_hook_delete":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookDeleteStep(n, c) },
	"step.okta_inline_hook_activate":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookActivateStep(n, c) },
	"step.okta_inline_hook_deactivate": func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookDeactivateStep(n, c) },
	"step.okta_inline_hook_execute":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newInlineHookExecuteStep(n, c) },

	// Domains
	"step.okta_domain_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDomainCreateStep(n, c) },
	"step.okta_domain_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newDomainGetStep(n, c) },
	"step.okta_domain_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newDomainListStep(n, c) },
	"step.okta_domain_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDomainDeleteStep(n, c) },
	"step.okta_domain_verify": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDomainVerifyStep(n, c) },

	// Brands & Themes
	"step.okta_brand_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newBrandGetStep(n, c) },
	"step.okta_brand_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newBrandListStep(n, c) },
	"step.okta_brand_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newBrandUpdateStep(n, c) },
	"step.okta_theme_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newThemeGetStep(n, c) },
	"step.okta_theme_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newThemeListStep(n, c) },
	"step.okta_theme_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newThemeUpdateStep(n, c) },

	// Org Settings
	"step.okta_org_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newOrgGetStep(n, c) },
	"step.okta_org_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newOrgUpdateStep(n, c) },
}

// createStep dispatches to the appropriate step constructor.
func createStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	constructor, ok := stepRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("okta plugin: unknown step type %q", typeName)
	}
	return constructor(name, config)
}

// allStepTypes returns all registered step type strings.
func allStepTypes() []string {
	types := make([]string, 0, len(stepRegistry))
	for k := range stepRegistry {
		types = append(types, k)
	}
	sort.Strings(types)
	return types
}
