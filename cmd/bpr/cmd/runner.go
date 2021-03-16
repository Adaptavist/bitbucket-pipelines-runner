package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
)

func hasFailedSteps(steps []model.Step) bool {
	fails := model.FilterSteps(steps, func(s model.Step) bool {
		return s.State.Result.HasError()
	})
	return len(fails) > 0
}

// printStepLogs with a pretty lazy implementation
func printStepLogs(logs map[string]string) {
	for step, logStr := range logs {
		log.Printf("%s output >\n", step)
		fmt.Println(logStr)
	}
}

func DoDryRun(opts client.PipelineOpts) map[string]string {
	var out []string
	out = append(out, "variables:")
	for _, pipelineVariable := range opts.Variables {
		if pipelineVariable.Secured {
			out = append(out, fmt.Sprintf("  %s: !SECURED!", pipelineVariable.Key))
		} else {
			out = append(out, fmt.Sprintf("  %s: %s", pipelineVariable.Key, pipelineVariable.Value))
		}
	}

	return map[string]string{"dry": strings.Join(out, "\n")}
}

func DoRun(http http.Client, opts client.PipelineOpts) (logs map[string]string, err error) {
	if opts.Dry {
		logs = DoDryRun(opts)
		return
	}

	bitbucket := client.NewClient(http).WithSleep(2 * time.Second)
	pipeline, err := bitbucket.PostPipelineAndWait(opts)
	pipelineFailed := pipeline.State.Result.Name == "FAILED"

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

	if pipelineFailed {
		err = errors.New("pipeline failed")
	}

	return
}
