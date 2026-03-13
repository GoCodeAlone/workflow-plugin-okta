package internal

import (
	"context"
	"testing"
)

func TestUserCreateStep_MissingClient(t *testing.T) {
	step, err := newUserCreateStep("test", map[string]any{"module": "nonexistent-okta"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"email": "test@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestUserCreateStep_MissingEmail(t *testing.T) {
	// Register a dummy client
	RegisterClient("test-create", &OktaClient{OrgURL: "https://test.okta.com", APIToken: "tok"})
	defer UnregisterClient("test-create")

	step, _ := newUserCreateStep("test", map[string]any{"module": "test-create"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing email")
	}
}

func TestUserGetStep_MissingUserId(t *testing.T) {
	step, _ := newUserGetStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing userId")
	}
}

func TestUserListStep_MissingClient(t *testing.T) {
	step, _ := newUserListStep("test", map[string]any{"module": "nonexistent-list"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestUserUpdateStep_MissingUserId(t *testing.T) {
	step, _ := newUserUpdateStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing userId")
	}
}

func TestUserDeleteStep_MissingUserId(t *testing.T) {
	step, _ := newUserDeleteStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing userId")
	}
}

func TestUserActivateStep_MissingUserId(t *testing.T) {
	step, _ := newUserActivateStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing userId")
	}
}

func TestUserChangePasswordStep_MissingNewPassword(t *testing.T) {
	// Register a dummy client
	RegisterClient("test-chpw", &OktaClient{OrgURL: "https://test.okta.com", APIToken: "tok"})
	defer UnregisterClient("test-chpw")

	step, _ := newUserChangePasswordStep("test", map[string]any{"module": "test-chpw"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{"userId": "user123"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing newPassword")
	}
}

func TestGroupCreateStep_MissingName(t *testing.T) {
	step, _ := newGroupCreateStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing name")
	}
}

func TestAppGetStep_MissingAppId(t *testing.T) {
	step, _ := newAppGetStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing appId")
	}
}

func TestFactorEnrollStep_MissingParams(t *testing.T) {
	step, _ := newFactorEnrollStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing userId and factorType")
	}
}

func TestSessionGetStep_MissingSessionId(t *testing.T) {
	step, _ := newSessionGetStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing sessionId")
	}
}

func TestOrgGetStep_MissingClient(t *testing.T) {
	step, _ := newOrgGetStep("test", map[string]any{"module": "nonexistent-okta"})
	result, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Output["error"] == nil {
		t.Error("expected error for missing client")
	}
}

func TestAllStepTypes_Count(t *testing.T) {
	types := allStepTypes()
	if len(types) < 100 {
		t.Errorf("expected at least 100 step types, got %d", len(types))
	}
}

func TestCreateStep_UnknownType(t *testing.T) {
	_, err := createStep("step.okta_nonexistent", "test", map[string]any{})
	if err == nil {
		t.Error("expected error for unknown step type")
	}
}

func TestResolveValue_PrefersCurrentOverConfig(t *testing.T) {
	current := map[string]any{"key": "from-current"}
	config := map[string]any{"key": "from-config"}
	if v := resolveValue("key", current, config); v != "from-current" {
		t.Errorf("expected 'from-current', got %q", v)
	}
}

func TestResolveValue_FallsBackToConfig(t *testing.T) {
	current := map[string]any{}
	config := map[string]any{"key": "from-config"}
	if v := resolveValue("key", current, config); v != "from-config" {
		t.Errorf("expected 'from-config', got %q", v)
	}
}

func TestGetModuleName_Default(t *testing.T) {
	if n := getModuleName(map[string]any{}); n != "okta" {
		t.Errorf("expected default module name 'okta', got %q", n)
	}
}

func TestGetModuleName_Custom(t *testing.T) {
	if n := getModuleName(map[string]any{"module": "my-okta"}); n != "my-okta" {
		t.Errorf("expected 'my-okta', got %q", n)
	}
}
