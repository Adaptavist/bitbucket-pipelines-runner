package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/config"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/spec"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/utils"
	"log"
)

// directCommandOKOrPanic - check any of the vars are set or unset, if there is no consistancy... panic!
func directCommandOKOrPanic(owner string, repoSlug string, ref string, pipeline string) bool {
	if !utils.Empty(owner) || !utils.Empty(repoSlug) || !utils.Empty(ref) || !utils.Empty(pipeline) {
		// Error if any of them are not set at this point
		if utils.Empty(owner) || utils.Empty(repoSlug) || utils.Empty(ref) || utils.Empty(pipeline) {
			panic(errors.New("-owner, -repo, -ref, -pipeline must all be set to run a pipeline via flags"))
		}
	}
	return true
}

// printStepLogs with a pretty lazy implementation
func printStepLogs(logs map[string]string) {
	for step, logStr := range logs {
		log.Printf("step (%s) output >\n", step)
		fmt.Println(logStr)
	}
}

func main() {
	// Flags - will replace with cobra at some point
	ownerPtr := flag.String("owner", "", "-owner Username")
	repoSlugPtr := flag.String("repo", "", "-repo slug")
	refPtr := flag.String("ref", "", "-ref master")
	pipelinePtr := flag.String("pipeline", "", "-pipeline deploy")
	variablesPtr := flag.String("vars", "", `-vars '[{"key":"VAR_NAME", "value": "VAR_VALUE"}]'`)
	dryPtr := flag.Bool("dry", false, "-dry")
	// Any additional arguments are treated as filename and each one will be a configuration for a pipeline
	flag.Parse()
	files := flag.Args()

	// Initialise configuration
	configuration := config.LoadConfigOrPanic(true)
	http := configuration.GetHttp()

	// If any of the direct run params are set
	if !utils.Empty(*ownerPtr) || !utils.Empty(*repoSlugPtr) || !utils.Empty(*refPtr) || !utils.Empty(*pipelinePtr) {
		directCommandOKOrPanic(*ownerPtr, *repoSlugPtr, *refPtr, *pipelinePtr)

		if len(files) > 0 {
			panic(errors.New("you cannot run a newPipeline directly and reference newPipeline run specs at the same time"))
		}

		var variables model.Variables

		// Construct variables
		err := json.Unmarshal([]byte(utils.DefaultWhenEmpty(*variablesPtr, "[]")), &variables)
		utils.PanicIfNotNil(err)

		opts := client.PipelineOpts{
			Dry:       *dryPtr,
			Repo:      client.NewRepo(*ownerPtr, *repoSlugPtr),
			Target:    model.NewTarget(*refPtr, *pipelinePtr),
			Variables: variables,
		}

		logs, err := run(http, opts)

		printStepLogs(logs)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		var variables model.Variables
		err := json.Unmarshal([]byte(utils.DefaultWhenEmpty(*variablesPtr, "[]")), &variables)
		utils.PanicIfNotNil(err)
		hasFailures := false

		var specFiles = make(map[string]spec.Spec)
		for _, file := range files {
			specFile, err := spec.UnmarshalSpecsFile(file)
			utils.PanicIfNotNil(err)
			log.Printf("loaded %s\n", file)
			specFiles[file] = specFile
		}

		// OK now we can start doing some real work
		for file, pipelineSpecs := range specFiles {
			for key, pipelineSpec := range pipelineSpecs.Pipelines {
				log.Printf("%s/%s", file, key)
				opts, err := pipelineSpecs.MakePipelineOpts(key)
				// Are we good or are we skipping
				if err != nil {
					hasFailures = true
					log.Printf("failed to make pipelines opts: %s\n", err)
				}

				opts = opts.WithDry(*dryPtr)

				if hasFailures {
					log.Printf("skipped %s %s:%s", file, key, pipelineSpec)
				} else {
					// Append variables if we have any from the flags
					opts.Variables = model.AppendVariables(opts.Variables, variables)
					logs, err := run(http, opts)
					printStepLogs(logs)

					if err != nil {
						hasFailures = true
						log.Print(err)
					}
				}
			}
		}
	}
}
