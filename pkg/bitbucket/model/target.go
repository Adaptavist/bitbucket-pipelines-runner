package model

import (
	"fmt"
)

// Target spec for a Pipeline to run
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

// TargetSelector description what branch/tag and spec we are to execute
type TargetSelector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

// NewTarget constructs a Target
func NewTarget(branch string, pipeline string) Target {
	return Target{
		Type:    "pipeline_ref_target",
		RefType: "branch",
		RefName: branch,
		Selector: TargetSelector{
			Type:    "custom",
			Pattern: pipeline,
		},
	}
}
