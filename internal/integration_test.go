package internal_test

import (
	"testing"

	"github.com/GoCodeAlone/workflow/wftest"
)

// TestIntegration_GetUser verifies a pipeline that looks up a user and sets a flag.
func TestIntegration_GetUser(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  lookup:
    steps:
      - name: get
        type: step.okta_user_get
        config:
          user_id: "user-123"
      - name: confirm
        type: step.set
        config:
          values:
            found: true
`),
		wftest.MockStep("step.okta_user_get", wftest.Returns(map[string]any{
			"id":     "user-123",
			"status": "ACTIVE",
		})),
	)
	result := h.ExecutePipeline("lookup", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["found"] != true {
		t.Errorf("expected found=true, got %v", result.Output["found"])
	}
	step := result.StepOutput("get")
	if step == nil {
		t.Fatal("expected step output for 'get'")
	}
	if step["id"] != "user-123" {
		t.Errorf("expected id=user-123, got %v", step["id"])
	}
}

// TestIntegration_CreateUser verifies a pipeline that creates a user and records the call.
func TestIntegration_CreateUser(t *testing.T) {
	rec := wftest.RecordStep("step.okta_user_create")
	rec.WithOutput(map[string]any{
		"id":     "new-user-456",
		"status": "STAGED",
	})
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  provision:
    steps:
      - name: create
        type: step.okta_user_create
        config:
          email: "alice@example.com"
          firstName: "Alice"
          lastName: "Smith"
`),
		rec,
	)
	result := h.ExecutePipeline("provision", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if rec.CallCount() != 1 {
		t.Errorf("expected 1 call to okta_user_create, got %d", rec.CallCount())
	}
	step := result.StepOutput("create")
	if step == nil {
		t.Fatal("expected step output for 'create'")
	}
	if step["id"] != "new-user-456" {
		t.Errorf("expected id=new-user-456, got %v", step["id"])
	}
}

// TestIntegration_ListGroupsAndAddUser verifies a pipeline that lists groups then adds a user.
func TestIntegration_ListGroupsAndAddUser(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  onboard:
    steps:
      - name: list_groups
        type: step.okta_group_list
        config: {}
      - name: add_user
        type: step.okta_group_add_user
        config:
          group_id: "grp-001"
          user_id: "user-123"
      - name: done
        type: step.set
        config:
          values:
            enrolled: true
`),
		wftest.MockStep("step.okta_group_list", wftest.Returns(map[string]any{
			"groups": []any{
				map[string]any{"id": "grp-001", "profile": map[string]any{"name": "Engineers"}},
			},
		})),
		wftest.MockStep("step.okta_group_add_user", wftest.Returns(map[string]any{
			"success": true,
		})),
	)
	result := h.ExecutePipeline("onboard", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["enrolled"] != true {
		t.Errorf("expected enrolled=true, got %v", result.Output["enrolled"])
	}
	if !result.StepExecuted("list_groups") {
		t.Error("expected list_groups step to be executed")
	}
	if !result.StepExecuted("add_user") {
		t.Error("expected add_user step to be executed")
	}
}
