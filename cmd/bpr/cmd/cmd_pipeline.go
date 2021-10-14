package cmd

import (
	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/adaptavist/bitbucket_pipelines_runner/cmd/bpr/spec"
	"github.com/spf13/cobra"
	"log"
)

func makePipelineCmdTarget (refType string, refName string, target string) model.PipelineTarget {
	t := builders.Target().Ref(refType, refName)
	if target != "" {
		t.Pattern(target)
	}
	return t.Build()
}

// pipelineCmd represents the run command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Runs a single pipeline",
	Long: `Run pipeline with no variablesFlag
> bpr pipeline workspace/repo_slug/branch[/pipeline_name]
Run pipeline with variables
> bpr pipeline workspace/repo_slug/branch[/pipeline_name] --var "key=value"
Run pipeline with secured variables
> bpr pipeline workspace/repo_slug/branch[/pipeline_name] --secret "key=value"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("pipeline expected")
		}

		client := makeClient()
		variables, err := stringsToVars(variablesFlag, false)
		fatalIfNotNil(err)
		securedVariables, err := stringsToVars(secureVarsFlag, true)
		fatalIfNotNil(err)
		targetSpec, err := spec.StringToTarget(args[0])
		fatalIfNotNil(err)
		target := makePipelineCmdTarget(targetSpec.RefType, targetSpec.RefName, targetSpec.CustomTarget)

		request := model.PostPipelineRequest{
			Workspace:  &targetSpec.Workspace,
			Repository: &targetSpec.Repo,
			Pipeline:   builders.Pipeline().Variables(append(variables, securedVariables...)).Target(target).Build(),
		}

		logs, err := DoRun(client, request, dryRun)
		printStepLogs(logs)
		fatalIfNotNil(err)
	},
}

func init() {
	pipelineCmd.Flags().StringSliceVar(&variablesFlag, "var", []string{}, "-var key=value")
	pipelineCmd.Flags().StringSliceVar(&secureVarsFlag, "secret", []string{}, "-s-var key=value")
	rootCmd.AddCommand(pipelineCmd)
}
