package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
)

type PipelineRunner struct {
	Auth     http.Auth
	Repo     bitbucket.Repo
	Pipeline bitbucket.Pipeline
}

func NewPipelineRunner(auth http.Auth, repo bitbucket.Repo, pipeline bitbucket.Pipeline) PipelineRunner {
	return PipelineRunner{
		Auth:     auth,
		Repo:     repo,
		Pipeline: pipeline,
	}
}

func (r PipelineRunner) hasFailedSteps(steps []bitbucket.PipelineStep) bool {
	fails := bitbucket.FilterSteps(steps, func(s bitbucket.PipelineStep) bool {
		return s.State.Result.HasError()
	})
	return len(fails) > 0
}

func (r PipelineRunner) String() string {
	return fmt.Sprintf("%s:%s", r.Repo.String(), r.Pipeline.Target.GetTargetDescriptor())
}

func (r PipelineRunner) run() bool {
	pipelineRun, err := bitbucket.Run(r.Repo, r.Auth, r.Pipeline)
	utils.PanicIfNotNil(err)

	pipelineSteps, err := bitbucket.GetSteps(r.Repo, r.Auth, pipelineRun.UUID)
	utils.PanicIfNotNil(err)

	// Get step logs
	for _, pipelineStep := range pipelineSteps {
		if pipelineStep.State.Result.HasError() {
			log.Printf("%s %s:%s - %s",
				strings.ToLower(pipelineStep.State.Result.Name),
				r,
				pipelineStep,
				pipelineStep.State.Result.Error.Message)
		} else {
			stepLog, err := bitbucket.StepLogs(r.Repo, r.Auth, pipelineRun.UUID, pipelineStep.UUID)
			if err != nil {
				log.Print(err.Error())
			} else {
				log.Printf("log %s:%s\n%s\n", r, pipelineStep.Name, stepLog)
			}
		}
	}

	return !r.hasFailedSteps(pipelineSteps)
}

// RunPipeline on Bitbucket pipelines
func RunPipeline(auth http.Auth, repo bitbucket.Repo, pipeline bitbucket.Pipeline) bool {
	return NewPipelineRunner(auth, repo, pipeline).run()
}
