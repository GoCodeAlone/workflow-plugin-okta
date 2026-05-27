// Package internal implements the workflow-plugin-okta plugin.
package internal

import (
	"fmt"

	"github.com/GoCodeAlone/workflow-plugin-okta/internal/contracts"
	pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
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

// TypedModuleTypes returns protobuf-typed module types.
func (p *oktaPlugin) TypedModuleTypes() []string {
	return p.ModuleTypes()
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

// CreateTypedModule creates a module instance from protobuf-typed config.
func (p *oktaPlugin) CreateTypedModule(typeName, name string, config *anypb.Any) (sdk.ModuleInstance, error) {
	if typeName != "okta.provider" {
		return nil, fmt.Errorf("okta plugin: unknown typed module type %q", typeName)
	}
	factory := sdk.NewTypedModuleFactory(typeName, &contracts.ProviderConfig{}, func(name string, cfg *contracts.ProviderConfig) (sdk.ModuleInstance, error) {
		return newOktaModule(name, typedOktaModuleConfig(cfg))
	})
	return factory.CreateTypedModule(typeName, name, config)
}

// StepTypes returns the step type names this plugin provides.
func (p *oktaPlugin) StepTypes() []string {
	return allStepTypes()
}

// TypedStepTypes returns protobuf-typed step types.
func (p *oktaPlugin) TypedStepTypes() []string {
	return p.StepTypes()
}

// CreateStep creates a step instance of the given type.
func (p *oktaPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	return createStep(typeName, name, config)
}

// CreateTypedStep creates a step instance from protobuf-typed config.
func (p *oktaPlugin) CreateTypedStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	if _, ok := stepRegistry[typeName]; !ok {
		return nil, fmt.Errorf("%w: step type %q", sdk.ErrTypedContractNotHandled, typeName)
	}
	if typeName == "step.okta_auth_provider_describe" {
		return sdk.NewTypedStepFactory(typeName, &contracts.AuthProviderDescribeConfig{}, &contracts.AuthProviderDescribeInput{}, typedAuthProviderDescribe).CreateTypedStep(typeName, name, config)
	}
	return sdk.NewTypedStepFactory(typeName, &contracts.OktaStepConfig{}, &contracts.OktaStepInput{}, typedOktaStepHandler(typeName)).CreateTypedStep(typeName, name, config)
}

// ContractRegistry exposes strict protobuf contracts for modules and steps.
func (p *oktaPlugin) ContractRegistry() *pb.ContractRegistry {
	return oktaContractRegistry
}

var oktaContractRegistry = &pb.ContractRegistry{
	FileDescriptorSet: &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			protodesc.ToFileDescriptorProto(structpb.File_google_protobuf_struct_proto),
			protodesc.ToFileDescriptorProto(contracts.File_internal_contracts_okta_proto),
		},
	},
	Contracts: oktaContractDescriptors(),
}

func oktaContractDescriptors() []*pb.ContractDescriptor {
	contracts := []*pb.ContractDescriptor{
		oktaModuleContract("okta.provider", "ProviderConfig"),
	}
	for _, stepType := range allStepTypes() {
		if stepType == "step.okta_auth_provider_describe" {
			contracts = append(contracts, oktaStepContract(stepType, "AuthProviderDescribeConfig", "AuthProviderDescribeInput", "AuthProviderDescribeOutput"))
			continue
		}
		contracts = append(contracts, oktaStepContract(stepType, "OktaStepConfig", "OktaStepInput", "OktaStepOutput"))
	}
	return contracts
}

func oktaModuleContract(moduleType, configMessage string) *pb.ContractDescriptor {
	const pkg = "workflow.plugins.okta.v1."
	return &pb.ContractDescriptor{
		Kind:          pb.ContractKind_CONTRACT_KIND_MODULE,
		ModuleType:    moduleType,
		ConfigMessage: pkg + configMessage,
		Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
	}
}

func oktaStepContract(stepType, configMessage, inputMessage, outputMessage string) *pb.ContractDescriptor {
	const pkg = "workflow.plugins.okta.v1."
	return &pb.ContractDescriptor{
		Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
		StepType:      stepType,
		ConfigMessage: pkg + configMessage,
		InputMessage:  pkg + inputMessage,
		OutputMessage: pkg + outputMessage,
		Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
	}
}
