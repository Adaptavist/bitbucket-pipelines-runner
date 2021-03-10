package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/config"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
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

func run(repo bitbucket.Repo, auth http.Auth, pipeline bitbucket.Pipeline) {
	pipelineRun, err := bitbucket.Run(repo, auth, pipeline)
	utils.PanicIfNotNil(err)
	pipelineSteps, err := bitbucket.Steps(repo, auth, pipelineRun.UUID)
	utils.PanicIfNotNil(err)

	// Get step logs
	for _, pipelineStep := range pipelineSteps {
		log, err := bitbucket.StepLogs(repo, auth, pipelineRun.UUID, pipelineStep.UUID)
		utils.PanicIfNotNil(err)
		fmt.Println(pipelineStep.Name)
		fmt.Println("================================================================================")
		fmt.Println(log)
	}
}

func main() {
	// Flags - will replace with cobra at some point
	ownerPtr := flag.String("owner", "", "owner of the repo")
	repoSlugPtr := flag.String("repo", "", "repo of the pipeline")
	refPtr := flag.String("ref", "", "git ref where the pipeline exists, could be a tag or branch name")
	pipelinePtr := flag.String("pipeline", "", "which pipeling in the ref do you want to run")
	variablesPtr := flag.String("vars", "", "JSON encoded string of variables to pass towards the pipeline")
	// Any additional arguments are treated as filename and each one will be a config for a pipeline
	flag.Parse()
	files := flag.Args()

	// Initialise config
	config := config.LoadConfigOrPanic()
	auth := config.GetAuth()

	// If any of the direct run params are set
	if !utils.Empty(*ownerPtr) || !utils.Empty(*repoSlugPtr) || !utils.Empty(*refPtr) || !utils.Empty(*pipelinePtr) {

		directIssetOrPanic(*ownerPtr, *repoSlugPtr, *refPtr, *pipelinePtr)

		if len(files) > 0 {
			panic(errors.New("You cannot run a pipeline directly and reference pipeline run specs at the same time"))
		}

		log.Printf("Running %s/%s/%s/%s", *ownerPtr, *repoSlugPtr, *refPtr, *pipelinePtr)

		// Construct variales
		repo := bitbucket.NewRepo(*ownerPtr, *repoSlugPtr)
		target := bitbucket.NewBranchTarget(*refPtr, *pipelinePtr)
		variables, err := bitbucket.UnmarshalVariables(utils.DefaultWhenEmpty(*variablesPtr, "[]"))
		utils.PanicIfNotNil(err)
		pipeline := bitbucket.NewPipeline(target, variables)

		// Run the pipeline
		run(repo, auth, pipeline)
	} else {
		variables, err := bitbucket.UnmarshalVariables(utils.DefaultWhenEmpty(*variablesPtr, "[]"))
		utils.PanicIfNotNil(err)
		for _, file := range files {
			specs, err := pipeline.UnmarshalSpecsFile(file)
			utils.PanicIfNotNil(err)
			for i, spec := range specs {
				fmt.Printf("Running: %s:%s\n", file, strconv.Itoa(i))
				repo := spec.GetWorkspaceRepo()
				pipeline := spec.GetPipeline()
				pipeline.Variables = append(pipeline.Variables, variables...)
				run(repo, auth, pipeline)
			}
		}
	}
}
