package main

import (
	"fmt"

	"os"

	"github.com/revanite-io/pvtr-github-repo/baseline"
	"github.com/revanite-io/pvtr-github-repo/data"
	"github.com/revanite-io/pvtr-github-repo/evaluation_plans"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/pluginkit"
)

var (
	// Version is to be replaced at build time by the associated tag
	Version = "0.0.0"
	// VersionPostfix is a marker for the version such as "dev", "beta", "rc", etc.
	VersionPostfix = "dev"
	// GitCommitHash is the commit at build time
	GitCommitHash = ""
	// BuiltAt is the actual build datetime
	BuiltAt = ""

	PluginName   = "github-repo"
	RequiredVars = []string{
		"owner",
		"repo",
		"token",
	}
)

func main() {
	if VersionPostfix != "" {
		Version = fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}

	// NewVessel may take a payload for all suites to reference
	pvtrVessel := pluginkit.NewEvaluationOrchestrator(PluginName, data.Loader, RequiredVars)

	// Evaluation Suite may optionally take a payload to selectively override the data specified in NewVessel
	pvtrVessel.AddEvaluationSuite("OSPS_B", data.Loader, evaluation_plans.OSPS_B)

	//TODO create a map of key value pairs key -> assessment.requireentId value -> recommendation from the baseline reader package

	br := baseline.NewReader()

	data, err := br.ReadAllYAMLFiles()
	if err != nil {
		//do something
	}

	recommendation := data.Catalog.ControlFamilies[0].Controls[0].AssessmentRequirements[0].Recommendation
	
	for _, family := range data.Catalog.ControlFamilies {
		for control := range family.Controls{
			for requirement := range control.
		}
	}



	runCmd := command.NewPluginCommands(
		PluginName,
		Version,
		VersionPostfix,
		GitCommitHash,
		pvtrVessel,
	)

	err := runCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
