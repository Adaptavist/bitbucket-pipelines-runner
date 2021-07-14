package cmd

import (
	"log"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/cmd/spec"
	"github.com/spf13/cobra"
)

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
		httpClient := getHTTPClient()
		variables, err := stringsToVars(variablesFlag, false)
		fatalIfNotNil(err)
		securedVariables, err := stringsToVars(secureVarsFlag, true)
		fatalIfNotNil(err)
		targetSpec, err := spec.StringToTarget(args[0])
		fatalIfNotNil(err)
		target := model.NewTarget(targetSpec.RefType, targetSpec.RefName)
		if targetSpec.CustomTarget != "" {
			target.WithCustomTarget(targetSpec.CustomTarget)
		}
		opts := client.NewPipelineOpts().
			WithDry(dryRun).
			WithVariables(model.AppendVariables(variables, securedVariables)).
			WithRepo(client.NewRepo(targetSpec.Workspace, targetSpec.Repo)).
			WithTarget(target)
		logs, err := DoRun(httpClient, opts)
		printStepLogs(logs)
		fatalIfNotNil(err)
	},
}

func init() {
	pipelineCmd.Flags().StringSliceVar(&variablesFlag, "var", []string{}, "-var key=value")
	pipelineCmd.Flags().StringSliceVar(&secureVarsFlag, "secret", []string{}, "-s-var key=value")
	rootCmd.AddCommand(pipelineCmd)
}
