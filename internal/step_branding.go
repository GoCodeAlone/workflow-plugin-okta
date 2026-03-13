package internal

import (
	"context"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// brandGetStep implements step.okta_brand_get
type brandGetStep struct {
	name       string
	moduleName string
}

func newBrandGetStep(name string, config map[string]any) (*brandGetStep, error) {
	return &brandGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *brandGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	brandID := resolveValue("brandId", current, config)
	if brandID == "" {
		return &sdk.StepResult{Output: errResult("brandId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/brands/"+brandID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// brandListStep implements step.okta_brand_list
type brandListStep struct {
	name       string
	moduleName string
}

func newBrandListStep(name string, config map[string]any) (*brandListStep, error) {
	return &brandListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *brandListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	result, err := oktaGet(client, "/api/v1/brands", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("brands", items)}, nil
}

// brandUpdateStep implements step.okta_brand_update
type brandUpdateStep struct {
	name       string
	moduleName string
}

func newBrandUpdateStep(name string, config map[string]any) (*brandUpdateStep, error) {
	return &brandUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *brandUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	brandID := resolveValue("brandId", current, config)
	if brandID == "" {
		return &sdk.StepResult{Output: errResult("brandId is required")}, nil
	}
	body := map[string]any{}
	if agreeToCustomPrivacyPolicy := resolveBool("agreeToCustomPrivacyPolicy", current, config); agreeToCustomPrivacyPolicy {
		body["agreeToCustomPrivacyPolicy"] = true
	}
	if customPrivacyPolicyURL := resolveValue("customPrivacyPolicyUrl", current, config); customPrivacyPolicyURL != "" {
		body["customPrivacyPolicyUrl"] = customPrivacyPolicyURL
	}
	if removePoweredByOkta := resolveBool("removePoweredByOkta", current, config); removePoweredByOkta {
		body["removePoweredByOkta"] = true
	}
	result, err := oktaPut(client, "/api/v1/brands/"+brandID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true}}, nil
}

// themeGetStep implements step.okta_theme_get
type themeGetStep struct {
	name       string
	moduleName string
}

func newThemeGetStep(name string, config map[string]any) (*themeGetStep, error) {
	return &themeGetStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *themeGetStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	brandID := resolveValue("brandId", current, config)
	themeID := resolveValue("themeId", current, config)
	if brandID == "" || themeID == "" {
		return &sdk.StepResult{Output: errResult("brandId and themeId are required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/brands/"+brandID+"/themes/"+themeID, nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": result}}, nil
}

// themeListStep implements step.okta_theme_list
type themeListStep struct {
	name       string
	moduleName string
}

func newThemeListStep(name string, config map[string]any) (*themeListStep, error) {
	return &themeListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *themeListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	brandID := resolveValue("brandId", current, config)
	if brandID == "" {
		return &sdk.StepResult{Output: errResult("brandId is required")}, nil
	}
	result, err := oktaGet(client, "/api/v1/brands/"+brandID+"/themes", nil)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("themes", items)}, nil
}

// themeUpdateStep implements step.okta_theme_update
type themeUpdateStep struct {
	name       string
	moduleName string
}

func newThemeUpdateStep(name string, config map[string]any) (*themeUpdateStep, error) {
	return &themeUpdateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *themeUpdateStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	brandID := resolveValue("brandId", current, config)
	themeID := resolveValue("themeId", current, config)
	if brandID == "" || themeID == "" {
		return &sdk.StepResult{Output: errResult("brandId and themeId are required")}, nil
	}
	body := map[string]any{}
	if primaryColorHex := resolveValue("primaryColorHex", current, config); primaryColorHex != "" {
		body["primaryColorHex"] = primaryColorHex
	}
	if secondaryColorHex := resolveValue("secondaryColorHex", current, config); secondaryColorHex != "" {
		body["secondaryColorHex"] = secondaryColorHex
	}
	if signInPageTouchPointVariant := resolveValue("signInPageTouchPointVariant", current, config); signInPageTouchPointVariant != "" {
		body["signInPageTouchPointVariant"] = signInPageTouchPointVariant
	}
	result, err := oktaPut(client, "/api/v1/brands/"+brandID+"/themes/"+themeID, body)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	if m := toMap(result); m != nil {
		return &sdk.StepResult{Output: mapResult(m)}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"updated": true}}, nil
}
