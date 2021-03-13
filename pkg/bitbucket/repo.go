package bitbucket

import (
	"encoding/json"
	"fmt"
)

// Variable for a BitBucket pipeline
type Variable struct {
	Key     string `json:"key" yaml:"key"`
	Value   string `json:"value" yaml:"value"`
	Secured bool   `json:"secured" yaml:"secured"`
}

// Variables for a BitBucket pipeline
type Variables []Variable

// UnmarshalVariables from string JSON input
func UnmarshalVariables(v string) (variables Variables, err error) {
	err = json.Unmarshal([]byte(v), &variables)
	return
}

// Target pipeline for a BitBucket to run
type Target struct {
	Type     string         `json:"type"`
	RefType  string         `json:"ref_type"`
	RefName  string         `json:"ref_name"`
	Selector TargetSelector `json:"selector"`
}

func (t Target) String() string {
	return fmt.Sprintf("%s:%s", t.RefName, t.Selector.Pattern)
}

func (t Target) GetTargetDescriptor() string {
	return t.String()
}

// TargetSelector description what branch/tag and pipline we are to execute
type TargetSelector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

// NewBranchTarget constructs a Target
func NewBranchTarget(branch string, target string) Target {
	return Target{
		Type:    "pipeline_ref_target",
		RefType: "branch",
		RefName: branch,
		Selector: TargetSelector{
			Type:    "custom",
			Pattern: target,
		},
	}
}

// Pipeline spec to for the BitBucket Pipeline API
type Pipeline struct {
	Target    Target    `json:"target"`
	Variables Variables `json:"variables"`
}

// NewPipeline constructor
func NewPipeline(t Target, v Variables) Pipeline {
	return Pipeline{
		Target:    t,
		Variables: v,
	}
}

// Repo for BitBucket
type Repo struct {
	Workspace string
	Slug      string
}

// String of the repo
func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Workspace, r.Slug)
}

// NewRepo constructs a Repo variable
func NewRepo(owner, repoSlug string) Repo {
	return Repo{
		Workspace: owner,
		Slug:      repoSlug,
	}
}

// PipelineState in BitBucket
type PipelineState struct {
	Name string `json:"name"`
}

// PipelineResponse containing details we need to lookup steps and status
type PipelineResponse struct {
	UUID        string        `json:"uuid"`
	CompletedOn string        `json:"completed_on"`
	State       PipelineState `json:"state"`
}

type PipelineStepStateResultError struct {
	Message string `json:"message"`
}

// PipelineStepStateResult shows if the step completed successfull
type PipelineStepStateResult struct {
	Name  string                       `json:"name"`
	Error PipelineStepStateResultError `json:"error"`
}

func (s PipelineStepStateResult) String() string {
	return s.Name
}

func (r PipelineStepStateResult) HasError() bool {
	return r.Error != PipelineStepStateResultError{}
}

func (s PipelineStepStateResult) IsSuccessfull() bool {
	return s.Name == "SUCCESSFUL"
}

// PipelineStepState provides the state the current pipeline step is in
type PipelineStepState struct {
	Name   string                  `json:"name"`
	Result PipelineStepStateResult `json:"result"`
}

func (s PipelineStepState) IsCompleted() bool {
	return s.Name == "COMPLETED"
}

// PipelineStep on a running/complete pipeline
type PipelineStep struct {
	UUID  string            `json:"uuid"`
	Name  string            `json:"name"`
	State PipelineStepState `json:"state"`
}

func (s PipelineStep) String() string {
	return s.Name
}

// FilterSteps does what it says
func FilterSteps(steps []PipelineStep, test func(PipelineStep) bool) (ret []PipelineStep) {
	for _, s := range steps {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// PipelineSteps response
type PipelineSteps struct {
	Page   int            `json:"page"`
	Values []PipelineStep `json:"values"`
	Next   string         `json:"next"`
}
