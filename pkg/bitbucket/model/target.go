package model

import (
	"fmt"
)

const (
	RefTypeBranch = "branch"
	RefTypeTag    = "tag"
)

// Target spec for a Pipeline to run
type Target struct {
	Type     string          `json:"type,omitempty"`
	RefType  string          `json:"ref_type,omitempty"`
	RefName  string          `json:"ref_name,omitempty"`
	Selector *TargetSelector `json:"selector,omitempty"`
	Commit   *Commit         `json:"commit,omitempty"`
}

func (t Target) String() string {
	return fmt.Sprintf("%s/%s", t.RefName, t.Selector.Pattern)
}

func (t Target) GetTargetDescriptor() string {
	return t.String()
}

// WithCustomTarget adds a custom target
func (t *Target) WithCustomTarget(target string) *Target {
	t.Selector = &TargetSelector{
		Type:    "custom",
		Pattern: target,
	}
	return t
}

type Commit struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
}

// TargetSelector description what branch/tag and spec we are to execute
type TargetSelector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

// NewTarget constructs a Target
func NewTarget(refType, refName string) Target {
	return Target{
		Type:    "pipeline_ref_target",
		RefType: refType,
		RefName: refName,
	}
}
