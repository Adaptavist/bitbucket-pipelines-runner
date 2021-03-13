package bitbucket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
)

const sleep = 10 * time.Second
const pipelinesURLTemplate = "https://api.bitbucket.org/2.0/repositories/%s/%s/pipelines/"

// GetState of the pipeline
func GetState(r Repo, a http.Auth, uuid string) (resp PipelineResponse, err error) {
	url := strings.ToLower(fmt.Sprintf(pipelinesURLTemplate, r.Workspace, r.Slug)) + uuid
	res, err := http.Get(a, url)

	if err != nil {
		return
	}

	err = json.Unmarshal(res, &resp)
	return
}

// Run a bit bucket pipeline to till its finished
func Run(r Repo, a http.Auth, p Pipeline) (result PipelineResponse, err error) {
	url := strings.ToLower(fmt.Sprintf(pipelinesURLTemplate, r.Workspace, r.Slug))
	res, err := http.Post(a, url, p)

	if err != nil {
		return
	}

	err = json.Unmarshal(res, &result)

	if err != nil {
		return
	}

	// We are going to keep looking at the pipeline till it finishes or fails
	for {
		time.Sleep(sleep)

		log.Printf("%s %s:%s ", strings.ToLower(result.State.Name), r.String(), p.Target.GetTargetDescriptor())

		if result.State.Name == "COMPLETED" || result.State.Name == "FAILED" {
			break
		}

		result, err = GetState(r, a, result.UUID)

		if err != nil {
			return
		}
	}

	return
}

// GetSteps - gets a list of steps in a given pipeline
func GetSteps(r Repo, a http.Auth, pipelineUUID string) (steps []PipelineStep, err error) {
	url := strings.ToLower(fmt.Sprintf(pipelinesURLTemplate, r.Workspace, r.Slug)) + fmt.Sprintf("%s/steps/", pipelineUUID)
	res, err := http.Get(a, url)

	if err != nil {
		return
	}

	var resp PipelineSteps
	err = json.Unmarshal(res, &resp)

	if err != nil {
		return
	}

	steps = append(steps, resp.Values...)

	for {

		if utils.Empty(resp.Next) {
			break
		}

		res, err = http.Get(a, resp.Next)

		if err != nil {
			return
		}

		err = json.Unmarshal(res, &resp)

		if err != nil {
			return
		}

		steps = append(steps, resp.Values...)
	}

	err = json.Unmarshal(res, &resp)

	return
}

// StepLogs from a given pipeline
func StepLogs(r Repo, a http.Auth, pipelineUUID string, stepUUID string) (log string, err error) {
	url := strings.ToLower(fmt.Sprintf(pipelinesURLTemplate, r.Workspace, r.Slug)) +
		fmt.Sprintf("%s/steps/%s/log", pipelineUUID, stepUUID)
	res, err := http.Get(a, url)

	if err != nil {
		return
	}

	log = string(res)

	return
}
