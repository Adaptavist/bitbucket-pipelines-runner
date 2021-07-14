package spec

import (
	"fmt"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"
)

type Variables map[string]string

// Merge two Variable lists
func (v Variables) Merge(merge Variables) (merged Variables) {
	if len(v) > 0 {
		merged = v
		for key, val := range merge {
			v[key] = val
		}
	} else {
		merged = merge
	}
	return
}

// PipelineTarget
type PipelineTarget struct {
	Workspace    string
	Repo         string
	RefType      string
	RefName      string
	CustomTarget string
}

// Pipeline represent a spec in our YAML config
type Pipeline struct {
	Pipeline  string    `yaml:"pipeline"`
	Variables Variables `yaml:"variables"`
}

// Pipelines maps our pipelines in our YAML config
type Pipelines map[string]Pipeline

// Spec of a spec which will be mapped to BitBuckets API later
type Spec struct {
	Pipelines Pipelines `yaml:"pipelines"`
	Variables Variables `yaml:"variables"`
}

// MakePipelineOpts from a Pipeline found in the Spec file
func (s Spec) MakePipelineOpts(name string) (opts client.PipelineOpts, err error) {
	spec, ok := s.Pipelines[name]

	if !ok {
		return opts, fmt.Errorf("%s spec spec not found", name)
	}

	targetSpec, err := spec.GetTarget()
	if err != nil {
		return
	}

	target := model.NewTarget(targetSpec.RefType, targetSpec.RefName)

	if targetSpec.CustomTarget != "" {
		target.WithCustomTarget(targetSpec.CustomTarget)
	}

	opts = bitbucket.NewPipelineOpts().
		WithTarget(target).
		WithRepo(client.NewRepo(targetSpec.Workspace, targetSpec.Repo)).
		WithVariables(s.ToBitbucketVariables(spec.Variables.Merge(s.Variables)))

	return
}
