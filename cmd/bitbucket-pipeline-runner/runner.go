package main

import (
	"errors"
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
)

func hasFailedSteps(steps []bitbucket.PipelineStep) bool {
	fails := bitbucket.FilterSteps(steps, func(s bitbucket.PipelineStep) bool {
		return s.State.Result.HasError()
	})
	return len(fails) > 0
}

func run(http http.Client, workspace, repo, branch, pipeline string, variables bitbucket.PipelineVariables) (logs map[string]string, err error) {
	c := bitbucket.NewClient(http)
	p, err := c.PostPipelineAndWait(workspace, repo, branch, pipeline, variables)

	logs = make(map[string]string)

	if err != nil {
		return logs, fmt.Errorf("failed to run pipeline: %s", err)
	}

	steps, err := c.GetSteps(workspace, repo, p.UUID)

	if err != nil {
		return logs, fmt.Errorf("failed to get steps: %s", err)
	}

	// Get step logs
	for _, step := range steps {
		if step.State.Result.HasError() {
			logs[step.Name] = step.State.Result.Error.Message
		} else {
			stepLog, err := c.GetStepLogs(workspace, repo, p.UUID, step.UUID)
			if err != nil {
				return logs, fmt.Errorf("unable to get all logs: %s", err)
			} else {
				logs[step.Name] = stepLog
			}
		}
	}

	if hasFailedSteps(steps) {
		return logs, errors.New("pipeline has failed steps")
	}

	return
}
