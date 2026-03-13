package internal

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// logListStep implements step.okta_log_list
type logListStep struct {
	name       string
	moduleName string
}

func newLogListStep(name string, config map[string]any) (*logListStep, error) {
	return &logListStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *logListStep) Execute(_ context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: errResult("okta client not found: " + s.moduleName)}, nil
	}
	qp := url.Values{}
	if since := resolveValue("since", current, config); since != "" {
		qp.Set("since", since)
	}
	if until := resolveValue("until", current, config); until != "" {
		qp.Set("until", until)
	}
	if filter := resolveValue("filter", current, config); filter != "" {
		qp.Set("filter", filter)
	}
	if q := resolveValue("q", current, config); q != "" {
		qp.Set("q", q)
	}
	if limit := resolveInt("limit", current, config); limit > 0 {
		qp.Set("limit", fmt.Sprintf("%d", limit))
	}
	result, err := oktaGet(client, "/api/v1/logs", qp)
	if err != nil {
		return &sdk.StepResult{Output: errResult(err.Error())}, nil
	}
	items := toSlice(result)
	if items == nil {
		items = []any{}
	}
	return &sdk.StepResult{Output: listResult("events", items)}, nil
}
