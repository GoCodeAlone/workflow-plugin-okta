package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoCodeAlone/workflow-plugin-okta/internal/contracts"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func typedOktaModuleConfig(cfg *contracts.ProviderConfig) map[string]any {
	if cfg == nil {
		return map[string]any{}
	}
	values := map[string]any{}
	if cfg.GetOrgUrl() != "" {
		values["orgUrl"] = cfg.GetOrgUrl()
	}
	if cfg.GetApiToken() != "" {
		values["apiToken"] = cfg.GetApiToken()
	}
	if cfg.GetClientId() != "" {
		values["clientId"] = cfg.GetClientId()
	}
	if cfg.GetPrivateKey() != "" {
		values["privateKey"] = cfg.GetPrivateKey()
	}
	if len(cfg.GetScopes()) > 0 {
		values["scopes"] = append([]string(nil), cfg.GetScopes()...)
	}
	return values
}

func typedOktaStepHandler(typeName string) sdk.TypedStepHandler[*contracts.OktaStepConfig, *contracts.OktaStepInput, *contracts.OktaStepOutput] {
	return func(ctx context.Context, req sdk.TypedStepRequest[*contracts.OktaStepConfig, *contracts.OktaStepInput]) (*sdk.TypedStepResult[*contracts.OktaStepOutput], error) {
		config, err := structToPlainMap(req.Config.GetValues())
		if err != nil {
			return nil, err
		}
		if req.Config.GetModule() != "" {
			config["module"] = req.Config.GetModule()
		}
		input, err := structToPlainMap(req.Input.GetValues())
		if err != nil {
			return nil, err
		}
		step, err := createStep(typeName, "typed", config)
		if err != nil {
			return nil, err
		}
		result, err := step.Execute(ctx, req.TriggerData, req.StepOutputs, mergePlainMaps(req.Current, input), req.Metadata, config)
		if err != nil {
			return nil, err
		}
		output, err := mapToStruct(result.Output)
		if err != nil {
			return nil, err
		}
		return &sdk.TypedStepResult[*contracts.OktaStepOutput]{
			Output:       &contracts.OktaStepOutput{Values: output, StopPipeline: result.StopPipeline},
			StopPipeline: result.StopPipeline,
		}, nil
	}
}

func typedAuthProviderDescribe(ctx context.Context, req sdk.TypedStepRequest[*contracts.AuthProviderDescribeConfig, *contracts.AuthProviderDescribeInput]) (*sdk.TypedStepResult[*contracts.AuthProviderDescribeOutput], error) {
	config := map[string]any{}
	if req.Config.GetProviderId() != "" {
		config["provider_id"] = req.Config.GetProviderId()
	}
	if req.Config.GetOrgUrl() != "" {
		config["org_url"] = req.Config.GetOrgUrl()
	}
	if len(req.Config.GetScopes()) > 0 {
		config["scopes"] = append([]string(nil), req.Config.GetScopes()...)
	}
	current := map[string]any{}
	if req.Input.GetProviderId() != "" {
		current["provider_id"] = req.Input.GetProviderId()
	}
	if req.Input.GetOrgUrl() != "" {
		current["org_url"] = req.Input.GetOrgUrl()
	}
	if len(req.Input.GetScopes()) > 0 {
		current["scopes"] = append([]string(nil), req.Input.GetScopes()...)
	}
	step := newAuthProviderDescribeStep("typed", config)
	result, err := step.Execute(ctx, req.TriggerData, req.StepOutputs, current, req.Metadata, nil)
	if err != nil {
		return nil, err
	}
	output, err := mapToProtoMessage(result.Output, &contracts.AuthProviderDescribeOutput{})
	if err != nil {
		return nil, err
	}
	return &sdk.TypedStepResult[*contracts.AuthProviderDescribeOutput]{Output: output}, nil
}

func structToPlainMap(value *structpb.Struct) (map[string]any, error) {
	if value == nil {
		return map[string]any{}, nil
	}
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal typed struct: %w", err)
	}
	values := map[string]any{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("decode typed struct: %w", err)
	}
	return values, nil
}

func mapToStruct(values map[string]any) (*structpb.Struct, error) {
	if values == nil {
		values = map[string]any{}
	}
	data, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("marshal step output map: %w", err)
	}
	decoded := map[string]any{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, fmt.Errorf("decode step output map: %w", err)
	}
	out, err := structpb.NewStruct(decoded)
	if err != nil {
		return nil, fmt.Errorf("encode typed output struct: %w", err)
	}
	return out, nil
}

func mergePlainMaps(sources ...map[string]any) map[string]any {
	merged := map[string]any{}
	for _, source := range sources {
		for key, value := range source {
			merged[key] = value
		}
	}
	return merged
}

func mapToProtoMessage[O proto.Message](values map[string]any, target O) (O, error) {
	typed := proto.Clone(target).(O)
	data, err := json.Marshal(values)
	if err != nil {
		return typed, fmt.Errorf("marshal output map: %w", err)
	}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: false}).Unmarshal(data, typed); err != nil {
		return typed, fmt.Errorf("decode typed protobuf output: %w", err)
	}
	return typed, nil
}
