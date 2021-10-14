package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/client"

	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/adaptavist/bitbucket_pipelines_runner/cmd/bpr/spec"
	"github.com/adaptavist/bitbucket_pipelines_runner/cmd/bpr/utils"
	"github.com/spf13/cobra"
)

var targetPipeline string
var targetType string
var targetName string
var onlyRun string
var chdir string

// specCmd represents the run command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Control numerous pipelines via yaml files",
	Long: `Run pipelines all .bpr.yml files
> bpr spec
Run pipelines with variables and secured variables
> bpr spec --var "key=value" --secrete "key2=value"
Run pipelines with overridden target pipeline
> bpr spec --target-pipeline pipeline_name --target-branch master
Run a specific pipeline from the files by its YAML key
> bpr spec --only my_pipeline`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cli = makeClient()

		if chdir != "" {
			err := os.Chdir(chdir)
			if err != nil {
				return err
			}
		}

		// Flags
		variables, err := stringsToVars(variablesFlag, false)

		if err != nil {
			return err
		}

		securedVariables, err := stringsToVars(secureVarsFlag, true)

		if err != nil {
			return err
		}

		// Spec files
		pipelines, err := loadPipelinesFromSpecFiles()

		if err != nil {
			return err
		}

		if onlyRun != "" {
			request, ok := pipelines[onlyRun]

			if !ok {
				log.Fatalf("pipeline not found: %s", onlyRun)
			}

			e := run(cli, request, variables, securedVariables)

			if e != nil {
				return e
			}
		} else {
			for _, request := range pipelines {
				e := run(cli, request, variables, securedVariables)

				if e != nil {
					return e
				}
			}
		}
		return nil
	},
}

func init() {
	specCmd.Flags().StringVar(&targetPipeline, "target-pipeline", "", "--target-pipeline example")
	specCmd.Flags().StringVar(&targetType, "target-type", "", "--target-type branch")
	specCmd.Flags().StringVar(&targetName, "target-name", "", "--target-name feat")
	specCmd.Flags().StringVar(&onlyRun, "only", "", "--only my_pipeline")
	specCmd.Flags().StringVar(&chdir, "chdir", "", "--cwd $CWD")
	rootCmd.AddCommand(specCmd)
}

// loadPipelinesFromSpecFiles loads all .bpr.yml files but also checks for duplicate pipeline keys and error if duplicates are found.
func loadPipelinesFromSpecFiles() (pipelines map[string]model.PostPipelineRequest, err error) {
	matches, err := filepath.Glob("*.bpr.yml")

	if err != nil {
		return
	}

	if len(matches) < 1 {
		err = errors.New("cannot find any .bpr.yml files")
		return
	}

	// map of filename to unmarshalled spec
	specFiles := make(map[string]spec.Spec)
	// links a pipeline key to a file, for duplicate checks
	specMap := make(map[string]string)

	// pass 1 - check for duplicates
	for _, match := range matches {
		specFile, err := spec.UnmarshalSpecsFile(match)
		utils.PanicIfNotNil(err)
		fmt.Printf("loaded %s\n", match)
		specFiles[match] = specFile
		// Add pipeline keys to specMap and error if a conflicted key is found
		for pKey := range specFile.Pipelines {
			dVal, found := specMap[pKey]
			if found {
				return nil, fmt.Errorf("duplicate pipeline key (%s) found in %s and %s", pKey, dVal, match)
			} else {
				specMap[pKey] = match
			}
		}
	}

	// pass 2 - reduce to just a map of pipelines
	pipelines = make(map[string]model.PostPipelineRequest)

	for _, fileSpec := range specFiles {
		for key := range fileSpec.Pipelines {
			pipeline, inErr := fileSpec.MakePostPipelineRequests(key)
			if inErr != nil {
				return nil, inErr
			}
			pipelines[key] = pipeline
		}
	}

	return
}

// run will run a pipeline with its variables with the ability to override
func run(cli client.Client, request model.PostPipelineRequest, variables model.PipelineVariables, securedVariables model.PipelineVariables) (err error) {
	if len(variables) > 0 {
		vars := append(*request.Variables, variables...)
		request.Variables = &vars
	}

	if len(securedVariables) > 0 {
		secVars := append(*request.Variables, securedVariables...)
		request.Variables = &secVars
	}

	// Use a builder to help override the target using args
	targetBuilder := builders.Target()
	targetBuilder.PipelineTarget = *request.Target

	if targetPipeline != "" {
		targetBuilder.Pattern(targetPipeline)
	}

	if targetType != "" {
		targetBuilder.RefType = targetType
	}

	if targetName != "" {
		targetBuilder.RefName = targetName
	}

	target := targetBuilder.Build()
	request.Target = &target

	fmt.Println("============================================================")
	fmt.Printf("running %s/%s\n", *request.Workspace, *request.Repository)
	logs, err := DoRun(cli, request, dryRun)
	printStepLogs(logs)
	return err
}
