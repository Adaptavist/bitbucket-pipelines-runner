package spec

import (
	"fmt"

	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
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

// MakePostPipelineRequests from a Pipeline found in the Spec file
func (s Spec) MakePostPipelineRequests(name string) (request model.PostPipelineRequest, err error) {
	spec, ok := s.Pipelines[name]

	if !ok {
		return request, fmt.Errorf("%s spec spec not found", name)
	}

	targetSpec, err := spec.GetTarget()
	if err != nil {
		return
	}

	targetBuilder := builders.Target().Ref(targetSpec.RefType, targetSpec.RefName)

	if targetSpec.CustomTarget != "" {
		targetBuilder.Pattern(targetSpec.CustomTarget)
	}

	target := targetBuilder.Build()
	variables := s.ToBitbucketVariables(spec.Variables.Merge(s.Variables))
	pipeline := builders.Pipeline().Target(target).Variables(variables).Build()

	request = model.PostPipelineRequest{
		Workspace:  &targetSpec.Workspace,
		Repository: &targetSpec.Repo,
		Pipeline:   pipeline,
	}

	return
}
