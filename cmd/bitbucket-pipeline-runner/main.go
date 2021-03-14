package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/config"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/pipeline"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
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
	ownerPtr := flag.String("owner", "", "owner of the repo")
	repoSlugPtr := flag.String("repo", "", "repo of the pipeline")
	refPtr := flag.String("ref", "", "git ref where the pipeline exists, could be a tag or branch name")
	pipelinePtr := flag.String("pipeline", "", "which pipeline in the ref do you want to run")
	variablesPtr := flag.String("vars", "", "JSON encoded string of variables to pass towards the pipeline")
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

		var variables bitbucket.PipelineVariables

		// Construct variables
		err := json.Unmarshal([]byte(utils.DefaultWhenEmpty(*variablesPtr, "[]")), &variables)
		utils.PanicIfNotNil(err)

		logs, err := run(http, *ownerPtr, *repoSlugPtr, *refPtr, *pipelinePtr, variables)

		printStepLogs(logs)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		var variables bitbucket.PipelineVariables
		err := json.Unmarshal([]byte(utils.DefaultWhenEmpty(*variablesPtr, "[]")), &variables)
		utils.PanicIfNotNil(err)
		hasFailures := false

		for _, file := range files {
			specs, err := pipeline.UnmarshalSpecsFile(file)
			utils.PanicIfNotNil(err)
			log.Printf("load %s", file)

			for i, spec := range specs {
				if hasFailures {
					log.Printf("skipped %s#%d:%s", file, i, spec)
				} else {
					logs, err := run(http, spec.Owner, spec.Repo, spec.Ref, spec.Pipeline, append(spec.Variables, variables...))
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
