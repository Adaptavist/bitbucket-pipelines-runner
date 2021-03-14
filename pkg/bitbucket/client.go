package bitbucket

import (
	"fmt"
	"log"
	"time"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
)

const sleep = 10 * time.Second
const pipelinesURLTemplate = "https://api.bitbucket.org/2.0/repositories/%s/%s/pipelines/"

// Client for bitbucket
type Client struct {
	Http http.Client
}

// NewClient for bitbucket
func NewClient(http http.Client) Client {
	return Client{
		Http: http,
	}
}

// PostPipeline started a bitbucket pipeline
func (c Client) PostPipeline(workspace, repoSlug, branch, pipeline string, variables PipelineVariables) (p PipelineResponse, err error) {
	u := fmt.Sprintf(pipelinesURLTemplate, workspace, repoSlug)
	r := NewPipeline(branch, pipeline, variables)
	err = c.Http.PostUnmarshalled(u, r, &p)
	return
}

// GetPipeline so we can check the status
func (c Client) GetPipeline(workspace, repoSlug, uuid string) (r PipelineResponse, err error) {
	u := fmt.Sprintf(pipelinesURLTemplate, workspace, repoSlug) + uuid
	err = c.Http.GetUnmarshalled(u, &r)
	return
}

// PostPipelineAndWait Run a pipeline till its complete
func (c Client) PostPipelineAndWait(workspace, repo, branch, pipeline string, variables PipelineVariables) (p PipelineResponse, err error) {
	p, err = c.PostPipeline(workspace, repo, branch, pipeline, variables)

	if err != nil {
		return
	}

	log.Printf("https://bitbucket.org/%s/%s/addon/pipelines/home#!/results/%d", workspace, repo, p.BuildNumber)

	for {
		log.Print(p.State)
		time.Sleep(sleep)

		if p.State.Name == "COMPLETED" {
			break
		}

		p, err = c.GetPipeline(workspace, repo, p.UUID)

		if err != nil {
			return
		}
	}

	return
}

// GetSteps - gets a list of steps in a given pipeline
func (c Client) GetSteps(workspace, repo, UUID string) (steps []PipelineStep, err error) {
	u := fmt.Sprintf(pipelinesURLTemplate, workspace, repo) + UUID + "/steps/"

	var resp PipelineStepsResponse
	err = c.Http.GetUnmarshalled(u, &resp)

	if err != nil {
		return
	}

	steps = append(steps, resp.Values...)

	for {
		if utils.Empty(resp.Next) {
			break
		}

		err = c.Http.GetUnmarshalled(resp.Next, &resp)

		if err != nil {
			return
		}

		steps = append(steps, resp.Values...)
	}
	return
}

// StepLogs from a given pipeline
func (c Client) GetStepLogs(workspace, slug, UUID, stepUUID string) (log string, err error) {
	url := fmt.Sprintf(pipelinesURLTemplate, workspace, slug) + UUID + "/steps/" + stepUUID + "/log"
	res, err := c.Http.Get(url)

	if err != nil {
		return
	}

	log = string(res)

	return
}
