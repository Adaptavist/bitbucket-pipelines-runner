package cmd

import (
	"errors"
	"fmt"
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/adaptavist/bitbucket_pipelines_runner/cmd/bpr/utils"
	"strings"
)

// printStepLogs with a pretty lazy implementation
func printStepLogs(logs map[string]string) {
	for step, logStr := range logs {
		fmt.Printf("%s output >\n", step)
		fmt.Println(logStr)
	}
}

func DoDryRun(request model.PostPipelineRequest) map[string]string {
	jsonStr := utils.MarshalFormatted(request)
	for _, pipelineVariable := range *request.Variables {
		if pipelineVariable.Secured {
			jsonStr = strings.ReplaceAll(jsonStr, pipelineVariable.Value, "!SECURED!")
		}
	}

	return map[string]string{"dry": jsonStr}
}

func DoRun(bitbucket client.Client, request model.PostPipelineRequest, dry bool) (logs map[string]string, err error) {
	if dry {
		logs = DoDryRun(request)
		return
	}

	// if the target is a tag, we need to look it up, get the commit hash and add it to the request
	if request.Target.RefType == model.RefTypeTag {
		var tag *model.TagResponse

		tag, err = bitbucket.GetTag(model.GetTagRequest{
			Workspace:  request.Workspace,
			Repository: request.Repository,
			Tag:        request.Target.RefName,
		})

		if err != nil {
			return
		}

		request.Target.Commit = &model.PipelineTargetCommit{
			Type: "commit",
			Hash: tag.Target.Hash,
		}
	}

	pipeline, err := bitbucket.RunPipeline(request)
	pipelineFailed := err != nil || pipeline.State.Result.Name == "FAILED"


	logs = make(map[string]string)

	if err != nil {
		return
	}

	steps, err := bitbucket.GetPipelineSteps(model.GetPipelineRequest{
		Workspace:  request.Workspace,
		Repository: request.Repository,
		Pipeline:   pipeline,
	})

	if err != nil {
		return
	}

	// get step logs
	for _, step := range steps {
		if step.State.Result.HasError() {
			logs[step.Name] = step.State.Result.Error.Message
		} else {
			log, e := bitbucket.GetPipelineStepLog(model.GetPipelineStepRequest{
				Workspace: request.Workspace,
				Repository: request.Repository,
				Pipeline: request.Pipeline,
				PipelineStep: &step,
			})
			if e != nil {
				err = e
				return
			} else {
				logs[step.Name] = string(log[:])
			}
		}
	}

	if pipelineFailed {
		err = errors.New("pipeline failed")
	}

	return
}
