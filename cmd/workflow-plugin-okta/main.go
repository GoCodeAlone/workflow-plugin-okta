package main

import (
	"github.com/GoCodeAlone/workflow-plugin-okta/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

var version = "dev"

func main() {
	sdk.Serve(internal.NewOktaPlugin(), sdk.WithBuildVersion(sdk.ResolveBuildVersion(internal.Version)))
}
