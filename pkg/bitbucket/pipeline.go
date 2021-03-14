package bitbucket

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PipelineStateResponse in BitBucket
type PipelineStateResponse struct {
	Name string `json:"name"`
}

func (s PipelineStateResponse) String() string {
	return strings.ToLower(s.Name)
}

// PipelineResponse containing details we need to lookup steps and status
type PipelineResponse struct {
	UUID        string                `json:"uuid"`
	BuildNumber int                   `json:"build_number"`
	CompletedOn string                `json:"completed_on"`
	State       PipelineStateResponse `json:"state"`
}

func (p PipelineResponse) String() string {
	return strings.ToLower(fmt.Sprintf("%s %s", p.UUID, p.State))
}

// PipelineVariable for a BitBucket pipeline
type PipelineVariable struct {
	Key     string `json:"key" yaml:"key"`
	Value   string `json:"value" yaml:"value"`
	Secured bool   `json:"secured" yaml:"secured"`
}

// PipelineVariables for a BitBucket pipeline
type PipelineVariables []PipelineVariable

// PipelineTarget pipeline for a BitBucket to run
type PipelineTarget struct {
	Type     string                 `json:"type"`
	RefType  string                 `json:"ref_type"`
	RefName  string                 `json:"ref_name"`
	Selector PipelineTargetSelector `json:"selector"`
}

func (t PipelineTarget) String() string {
	return fmt.Sprintf("%s:%s", t.RefName, t.Selector.Pattern)
}

func (t PipelineTarget) GetTargetDescriptor() string {
	return t.String()
}

// PipelineTargetSelector description what branch/tag and pipeline we are to execute
type PipelineTargetSelector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

// NewPipelineTarget constructs a PipelineTarget
func NewPipelineTarget(branch string, pipeline string) PipelineTarget {
	return PipelineTarget{
		Type:    "pipeline_ref_target",
		RefType: "branch",
		RefName: branch,
		Selector: PipelineTargetSelector{
			Type:    "custom",
			Pattern: pipeline,
		},
	}
}

// Pipeline spec to for the BitBucket Pipeline API
type Pipeline struct {
	Target    PipelineTarget    `json:"target"`
	Variables PipelineVariables `json:"variables"`
}

func (p Pipeline) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

// NewPipeline constructor
func NewPipeline(branch, pipeline string, variables PipelineVariables) Pipeline {
	return Pipeline{
		Target:    NewPipelineTarget(branch, pipeline),
		Variables: variables,
	}
}
