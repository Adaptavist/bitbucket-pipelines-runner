package cmd

import (
	"errors"
	"fmt"
	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/bitbucket/http"
	"log"
	"os"
	"path/filepath"

	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/cmd/spec"
	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/cmd/utils"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var httpClient = getHTTPClient()

		if chdir != "" {
			err = os.Chdir(chdir)
			if err != nil {
				return
			}
		}

		// Flags
		variables, err := stringsToVars(variablesFlag, false)

		if err != nil {
			return
		}

		securedVariables, err := stringsToVars(secureVarsFlag, true)

		if err != nil {
			return
		}

		// Spec files
		pipelines, err := loadPipelinesFromSpecFiles()

		if err != nil {
			return
		}

		if onlyRun != "" {
			opts, ok := pipelines[onlyRun]
			if !ok {
				log.Fatalf("pipeline not found: %s", onlyRun)
			}
			err = run(httpClient, opts, variables, securedVariables)
			if err != nil {
				return
			}
		} else {
			for _, opts := range pipelines {
				err = run(httpClient, opts, variables, securedVariables)
				if err != nil {
					return
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
func loadPipelinesFromSpecFiles() (pipelines map[string]client.PipelineOpts, err error) {
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
	pipelines = make(map[string]client.PipelineOpts)

	for _, fileSpec := range specFiles {
		for key := range fileSpec.Pipelines {
			pipeline, inErr := fileSpec.MakePipelineOpts(key)
			if inErr != nil {
				return nil, inErr
			}
			pipelines[key] = pipeline
		}
	}

	return
}

// run will run a pipeline with its variables with the ability to override
func run(httpClient http.Client, opts client.PipelineOpts, variables model.Variables, securedVariables model.Variables) (err error) {
	opts = opts.WithDry(dryRun)

	if len(variables) > 0 {
		opts.Variables = model.AppendVariables(opts.Variables, variables)
	}

	if len(securedVariables) > 0 {
		opts.Variables = model.AppendVariables(opts.Variables, securedVariables)
	}

	if targetPipeline != "" {
		opts.Target.WithCustomTarget(targetPipeline)
	}

	if targetType != "" {
		opts.Target.RefType = targetType
	}

	if targetName != "" {
		opts.Target.RefName = targetName
	}

	fmt.Println("============================================================")
	fmt.Printf("running %s\n", opts.String())
	logs, err := DoRun(httpClient, opts)
	printStepLogs(logs)
	return err
}