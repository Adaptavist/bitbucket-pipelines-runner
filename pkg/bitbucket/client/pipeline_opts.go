package client

import (
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
)

// PipelineOpts is in internal model for requesting a spec to be executed
type PipelineOpts struct {
	Dry       bool
	Repo      Repo
	Target    model.Target
	Variables model.Variables
}

func (o PipelineOpts) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", o.Repo.Workspace, o.Repo.Slug, o.Target.RefName, o.Target.Selector.Pattern)
}

// NewPipelineOpts constructs empty PipelineOpts object
func NewPipelineOpts() PipelineOpts {
	return PipelineOpts{}
}

func (o PipelineOpts) WithDry(dry bool) PipelineOpts {
	o.Dry = dry
	return o
}

// WithRepo adds Repo to the current PipelineOpts
func (o PipelineOpts) WithRepo(repo Repo) PipelineOpts {
	o.Repo = repo
	return o
}

// WithTarget adds Target to the current PipelineOpts
func (o PipelineOpts) WithTarget(target model.Target) PipelineOpts {
	o.Target = target
	return o
}

// WithVariables adds PipelineVariables to the current PipelineOpts
func (o PipelineOpts) WithVariables(variables model.Variables) PipelineOpts {
	o.Variables = variables
	return o
}

// NewPipelineRequest creates model.PipelineRequest from PipelineOpts
func (o PipelineOpts) NewPipelineRequest() model.PipelineRequest {
	return model.PipelineRequest{
		Target:    o.Target,
		Variables: o.Variables,
	}
}
