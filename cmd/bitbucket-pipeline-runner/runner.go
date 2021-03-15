package main

import (
	"errors"
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"log"
	"strings"
	"time"
)

func hasFailedSteps(steps []model.Step) bool {
	fails := model.FilterSteps(steps, func(s model.Step) bool {
		return s.State.Result.HasError()
	})
	return len(fails) > 0
}

func dryRun(opts client.PipelineOpts) map[string]string {
	var out []string

	log.Print(opts)
	out = append(out, "variables:")
	for _, pipelineVariable := range opts.Variables {
		if pipelineVariable.Secured {
			out = append(out, fmt.Sprintf("%s: !SECURED!", pipelineVariable.Key))
		} else {
			out = append(out, fmt.Sprintf("%s: %s!", pipelineVariable.Key, pipelineVariable.Value))
		}
	}

	return map[string]string{"dry": strings.Join(out, "\n")}
}

func run(http http.Client, opts client.PipelineOpts) (logs map[string]string, err error) {
	if opts.Dry {
		logs = dryRun(opts)
		return
	}

	bitbucket := client.NewClient(http).WithSleep(2 * time.Second)
	pipeline, err := bitbucket.PostPipelineAndWait(opts)

	logs = make(map[string]string)

	if err != nil {
		return
	}

	steps, err := bitbucket.GetSteps(opts, pipeline)

	if err != nil {
		return
	}

	var stepLog string

	// Get step logs
	for _, step := range steps {
		if step.State.Result.HasError() {
			logs[step.Name] = step.State.Result.Error.Message
		} else {
			stepLog, err = bitbucket.GetStepLogs(opts, pipeline, step)
			if err != nil {
				return
			} else {
				logs[step.Name] = stepLog
			}
		}
	}

	if hasFailedSteps(steps) {
		err = errors.New("spec has failed steps")
	}

	return
}
