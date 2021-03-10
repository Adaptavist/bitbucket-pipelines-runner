package bitbucket

import "encoding/json"

// Creds for authenticating with BitBucket (BasicAuth)
type Creds struct {
	Username string
	Password string
}

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
	Type     string                 `json:"type"`
	RefType  string                 `json:"ref_type"`
	RefName  string                 `json:"ref_name"`
	Selector map[string]interface{} `json:"selector"`
}

// NewBranchTarget constructs a Target
func NewBranchTarget(branch string, target string) Target {
	return Target{
		Type:    "pipeline_ref_target",
		RefType: "branch",
		RefName: branch,
		Selector: map[string]interface{}{
			"type":    "custom",
			"pattern": target,
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

// PipelineStep on a running/complete pipeline
type PipelineStep struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// PipelineSteps response
type PipelineSteps struct {
	Page   int            `json:"page"`
	Values []PipelineStep `json:"values"`
	Next   string         `json:"next"`
}
