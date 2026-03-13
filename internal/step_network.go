package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// networkZoneCreateStep implements step.okta_network_zone_create
type networkZoneCreateStep struct {
	name       string
	moduleName string
}

func newNetworkZoneCreateStep(name string, config map[string]any) (*networkZoneCreateStep, error) {
	return &networkZoneCreateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneCreateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	zoneName := resolveValue("name", current, config)
	zoneType := resolveValue("type", current, config)
	if zoneName == "" || zoneType == "" {
		return &sdk.StepResult{Output: errResult("name and type are required")}, nil
	}
	body := map[string]any{
		"type":    zoneType,
		"name":    zoneName,
		"gateways": resolveStringSlice("gateways", current, config),
		"proxies":  resolveStringSlice("proxies", current, config),
	}
	result, err := oktaPost(client, "/api/v1/zones", body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// networkZoneGetStep implements step.okta_network_zone_get
type networkZoneGetStep struct {
	name       string
	moduleName string
}

func newNetworkZoneGetStep(name string, config map[string]any) (*networkZoneGetStep, error) {
	return &networkZoneGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	zoneID := resolveValue("zoneId", current, config)
	if zoneID == "" {
		return &sdk.StepResult{Output: errResult("zoneId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/zones/"+zoneID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// networkZoneListStep implements step.okta_network_zone_list
type networkZoneListStep struct {
	name       string
	moduleName string
}

func newNetworkZoneListStep(name string, config map[string]any) (*networkZoneListStep, error) {
	return &networkZoneListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/zones", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("zones", items)}, nil
}

// networkZoneDeleteStep implements step.okta_network_zone_delete
type networkZoneDeleteStep struct {
	name       string
	moduleName string
}

func newNetworkZoneDeleteStep(name string, config map[string]any) (*networkZoneDeleteStep, error) {
	return &networkZoneDeleteStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneDeleteStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	zoneID := resolveValue("zoneId", current, config)
	if zoneID == "" {
		return &sdk.StepResult{Output: errResult("zoneId is required")}, nil
	}
	if err := oktaDelete(client, "/api/v1/zones/"+zoneID); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true}}, nil
}

// networkZoneActivateStep implements step.okta_network_zone_activate
type networkZoneActivateStep struct {
	name       string
	moduleName string
}

func newNetworkZoneActivateStep(name string, config map[string]any) (*networkZoneActivateStep, error) {
	return &networkZoneActivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneActivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	zoneID := resolveValue("zoneId", current, config)
	if zoneID == "" {
		return &sdk.StepResult{Output: errResult("zoneId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/zones/"+zoneID+"/lifecycle/activate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"activated": true}}, nil
}

// networkZoneDeactivateStep implements step.okta_network_zone_deactivate
type networkZoneDeactivateStep struct {
	name       string
	moduleName string
}

func newNetworkZoneDeactivateStep(name string, config map[string]any) (*networkZoneDeactivateStep, error) {
	return &networkZoneDeactivateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *networkZoneDeactivateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	zoneID := resolveValue("zoneId", current, config)
	if zoneID == "" {
		return &sdk.StepResult{Output: errResult("zoneId is required")}, nil
	}
	if _, err := oktaPostEmpty(client, "/api/v1/zones/"+zoneID+"/lifecycle/deactivate"); err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deactivated": true}}, nil
}
