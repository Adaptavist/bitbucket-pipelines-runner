package client

import (
	"fmt"
	"time"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/http"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/urls"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/cmd/utils"
)

// Client for BitBucket
type Client struct {
	Http  http.Client
	Sleep time.Duration
}

// NewClient for BitBucket
func NewClient(http http.Client) Client {
	return Client{
		Http: http,
	}
}

// WithSleep as a sleep time some of the larger methods
func (c Client) WithSleep(seconds time.Duration) Client {
	c.Sleep = seconds
	return c
}

func (c Client) sleep() {
	if c.Sleep != 0 {
		time.Sleep(c.Sleep)
	}
}

// PostPipeline using PipelineOpts to start a pipeline
func (c Client) PostPipeline(opts PipelineOpts) (p model.Pipeline, err error) {
	u := urls.Pipelines(opts.Repo.Workspace, opts.Repo.Slug)
	r := opts.NewPipelineRequest()
	err = c.Http.PostUnmarshalled(u, r, &p)
	return
}

// GetPipeline so we can check the status
func (c Client) GetPipeline(opts PipelineOpts, p model.Pipeline) (r model.Pipeline, err error) {
	u := urls.Pipeline(opts.Repo.Workspace, opts.Repo.Slug, p.UUID)
	err = c.Http.GetUnmarshalled(u, &r)
	return
}

// PostPipelineAndWait runs a spec till its complete
func (c Client) PostPipelineAndWait(opts PipelineOpts) (p model.Pipeline, err error) {
	p, err = c.PostPipeline(opts)

	if err != nil {
		return
	}

	// Would be handy for the end user to have a link
	fmt.Println(urls.PipelineWeb(opts.Repo.Workspace, opts.Repo.Slug, p.BuildNumber))

	for {
		fmt.Print(".")

		c.sleep()

		if p.State.Name == "COMPLETED" {
			fmt.Print("\n")
			break
		}

		p, err = c.GetPipeline(opts, p)

		if err != nil {
			return
		}
	}

	return
}

// GetSteps gets model.Steps in a given model.Pipeline
func (c Client) GetSteps(opts PipelineOpts, p model.Pipeline) (steps model.Steps, err error) {
	var resp model.StepsResponse
	err = c.Http.GetUnmarshalled(urls.PipelineSteps(opts.Repo.Workspace, opts.Repo.Slug, p.UUID), &resp)

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
func (c Client) GetStepLogs(opts PipelineOpts, p model.Pipeline, s model.Step) (logContent string, err error) {
	res, err := c.Http.Get(urls.PipelineStepLogs(opts.Repo.Workspace, opts.Repo.Slug, p.UUID, s.UUID))

	if err != nil {
		return
	}

	logContent = string(res)

	return
}

// GetTag returns a tag if it exists from the bitbucket API
func (c Client) GetTag(opts PipelineOpts) (tag model.TagResponse, err error) {
	url := urls.Tag(opts.Repo.Workspace, opts.Repo.Slug, opts.Target.RefName)
	err = c.Http.GetUnmarshalled(url, &tag)
	return
}
