package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/http"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/cmd/utils"
)

// func hasFailedSteps(steps []model.Step) bool {
// 	fails := model.FilterSteps(steps, func(s model.Step) bool {
// 		return s.State.Result.HasError()
// 	})
// 	return len(fails) > 0
// }

// printStepLogs with a pretty lazy implementation
func printStepLogs(logs map[string]string) {
	for step, logStr := range logs {
		fmt.Printf("%s output >\n", step)
		fmt.Println(logStr)
	}
}

func DoDryRun(opts client.PipelineOpts) map[string]string {
	jsonStr := utils.MarshalFormatted(opts)
	for _, pipelineVariable := range opts.Variables {
		if pipelineVariable.Secured {
			jsonStr = strings.ReplaceAll(jsonStr, pipelineVariable.Value, "!SECURED!")
		}
	}

	return map[string]string{"dry": jsonStr}
}

func DoRun(http http.Client, opts client.PipelineOpts) (logs map[string]string, err error) {
	if opts.Dry {
		logs = DoDryRun(opts)
		return
	}

	bitbucket := client.NewClient(http).WithSleep(2 * time.Second)

	// if the target is a tag, we need to look it up, get the commit hash and add it to the request
	if opts.Target.RefType == model.RefTypeTag {
		var tag model.TagResponse
		tag, err = bitbucket.GetTag(opts)

		if err != nil {
			return
		}

		opts.Target.Commit = &model.Commit{
			Type: "commit",
			Hash: tag.Target.Hash,
		}
	}

	pipeline, err := bitbucket.PostPipelineAndWait(opts)
	pipelineFailed := err != nil || pipeline.State.Result.Name == "FAILED"

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
