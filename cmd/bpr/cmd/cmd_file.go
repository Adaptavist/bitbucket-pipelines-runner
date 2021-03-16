package cmd

import (
	"log"
	"path/filepath"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/spec"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/utils"
	"github.com/spf13/cobra"
)

var targetPipeline string
var targetBranch string
var onlyRun string

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
	Run: func(cmd *cobra.Command, args []string) {
		// If true, it will skip till the end
		var hasFailures = false
		// Map of files to spec.Spec objects
		var specFiles = make(map[string]spec.Spec)
		// Map of Pipeline key to spec file, kinda of like specFiles bit inverts for lookups
		var specMap = make(map[string]string)
		// Get our client ready early
		var httpClient = getHTTPClient()

		// Flags
		variables, err := stringsToVars(variablesFlag, false)
		fatalIfNotNil(err)
		securedVariables, err := stringsToVars(secureVarsFlag, true)
		fatalIfNotNil(err)

		// Get a list of spec files
		matches, err := filepath.Glob("*.bpr.yml")
		utils.PanicIfNotNil(err)
		for _, match := range matches {
			specFile, err := spec.UnmarshalSpecsFile(match)
			utils.PanicIfNotNil(err)
			log.Printf("loaded %s\n", match)
			specFiles[match] = specFile
			// Add pipeline keys to specMap and error if a conflicted key is found
			for pKey := range specFile.Pipelines {
				dVal, found := specMap[pKey]
				if found {
					log.Printf("duplicate pipeline key (%s) found in %s and %s", pKey, dVal, match)
					hasFailures = true
				} else {
					specMap[pKey] = match
				}
			}
		}

		// Exit now if errors have already been found
		if hasFailures {
			log.Fatal("exiting")
		}

		// OK now we can start doing some real work
		for file, pipelineSpecs := range specFiles {
			for key, pipelineSpec := range pipelineSpecs.Pipelines {
				log.Printf("%s/%s", file, key)

				// Prepare the pipeline
				opts, err := pipelineSpecs.MakePipelineOpts(key)
				if err != nil {
					hasFailures = true
					log.Printf("failed to make pipelines opts: %s\n", err)
					continue
				}

				// Overrides
				opts = opts.WithDry(dryRun)

				if len(variables) > 0 {
					opts.Variables = model.AppendVariables(opts.Variables, variables)
				}

				if len(securedVariables) > 0 {
					opts.Variables = model.AppendVariables(opts.Variables, securedVariables)
				}

				if targetPipeline != "" {
					opts.Target.Selector.Pattern = targetPipeline
				}

				if targetBranch != "" {
					opts.Target.RefName = targetBranch
				}

				// Checks
				if hasFailures {
					log.Printf("skipped %s %s:%s", file, key, pipelineSpec)
					continue
				}

				if onlyRun != "" && onlyRun != key {
					continue
				}

				// Go!
				logs, err := DoRun(httpClient, opts)
				printStepLogs(logs)

				if err != nil {
					log.Print(err.Error())
					hasFailures = true
				}
			}
		}

		if hasFailures {
			log.Fatal("finished with errors")
		}
	},
}

func init() {
	specCmd.Flags().StringVar(&targetPipeline, "target-pipeline", "", "--target-pipeline example")
	specCmd.Flags().StringVar(&targetBranch, "target-branch", "", "--target-branch master")
	specCmd.Flags().StringVar(&onlyRun, "only", "", "--only my_pipeline")
	rootCmd.AddCommand(specCmd)
}
