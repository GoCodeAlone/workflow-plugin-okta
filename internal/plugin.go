// Package internal implements the workflow-plugin-okta plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// Version is set at build time via -ldflags
// "-X github.com/GoCodeAlone/workflow-plugin-okta/internal.Version=X.Y.Z"
var Version = "0.0.0"

// oktaPlugin implements sdk.PluginProvider, sdk.ModuleProvider, and sdk.StepProvider.
type oktaPlugin struct{}

// NewOktaPlugin returns a new oktaPlugin instance.
func NewOktaPlugin() sdk.PluginProvider {
	return &oktaPlugin{}
}

// Manifest returns plugin metadata.
func (p *oktaPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-okta",
		Version:     Version,
		Author:      "GoCodeAlone",
		Description: "Okta identity platform plugin (~130 step types across all Okta APIs)",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *oktaPlugin) ModuleTypes() []string {
	return []string{"okta.provider"}
}

// CreateModule creates a module instance of the given type.
func (p *oktaPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "okta.provider":
		m, err := newOktaModule(name, config)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("okta plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns the step type names this plugin provides.
func (p *oktaPlugin) StepTypes() []string {
	return allStepTypes()
}

// CreateStep creates a step instance of the given type.
func (p *oktaPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	return createStep(typeName, name, config)
}
