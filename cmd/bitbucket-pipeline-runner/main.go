package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/config"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/pipeline"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
)

// directIssetOrPanic - check any of the vars are set or unset, if there is no consistancy... panic!
func directIssetOrPanic(owner string, repoSlug string, ref string, pipeline string) bool {
	if !utils.Empty(owner) || !utils.Empty(repoSlug) || !utils.Empty(ref) || !utils.Empty(pipeline) {
		// Error if any of them are not set at this point
		if utils.Empty(owner) || utils.Empty(repoSlug) || utils.Empty(ref) || utils.Empty(pipeline) {
			panic(errors.New("-owner, -repo, -ref, -pipeline must all be set to run a pipeline via flags"))
		}
	}
	return true
}

func exit(ok bool) {
	if !ok {
		os.Exit(1)
	} else {
		os.Exit(0)
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
	auth := configuration.GetAuth()

	// If any of the direct run params are set
	if !utils.Empty(*ownerPtr) || !utils.Empty(*repoSlugPtr) || !utils.Empty(*refPtr) || !utils.Empty(*pipelinePtr) {

		directIssetOrPanic(*ownerPtr, *repoSlugPtr, *refPtr, *pipelinePtr)

		if len(files) > 0 {
			panic(errors.New("You cannot run a newPipeline directly and reference newPipeline run specs at the same time"))
		}

		// Construct variables
		repo := bitbucket.NewRepo(*ownerPtr, *repoSlugPtr)
		target := bitbucket.NewBranchTarget(*refPtr, *pipelinePtr)
		variables, err := bitbucket.UnmarshalVariables(utils.DefaultWhenEmpty(*variablesPtr, "[]"))
		utils.PanicIfNotNil(err)
		newPipeline := bitbucket.NewPipeline(target, variables)
		ok := RunPipeline(auth, repo, newPipeline)
		exit(ok)
	} else {
		variables, err := bitbucket.UnmarshalVariables(utils.DefaultWhenEmpty(*variablesPtr, "[]"))
		utils.PanicIfNotNil(err)
		hasFailures := false

		for _, file := range files {
			specs, err := pipeline.UnmarshalSpecsFile(file)
			utils.PanicIfNotNil(err)
			log.Printf("load %s", file)

			for i, spec := range specs {
				repo := spec.GetWorkspaceRepo()
				targetPipeline := spec.GetPipeline()
				targetPipeline.Variables = append(targetPipeline.Variables, variables...)
				runner := NewPipelineRunner(auth, repo, targetPipeline)

				if hasFailures {
					log.Printf("skipped %s#%d:%s", file, i, targetPipeline.Target)
				} else {
					ok := runner.run()
					if !ok {
						hasFailures = true
					}
					fmt.Print(strconv.FormatBool(ok))
				}
			}
		}

		exit(!hasFailures)
	}
}
