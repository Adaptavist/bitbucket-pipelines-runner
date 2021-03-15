package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// State of a Pipeline
type State struct {
	Name string `json:"name"`
}

// String representation of State
func (s State) String() string {
	return strings.ToLower(s.Name)
}

// PipelineRequest containing details we need to start a pipeline
type PipelineRequest struct {
	Target    Target    `json:"target"`
	Variables Variables `json:"variables"`
}

// Pipeline containing details we need to lookup steps and status
type Pipeline struct {
	UUID        string    `json:"uuid"`
	BuildNumber int       `json:"build_number"`
	CompletedOn string    `json:"completed_on"`
	State       State     `json:"state"`
	Target      Target    `json:"target"`
	Variables   Variables `json:"variables"`
}

// String representation of the Pipeline (uuid state)
func (p Pipeline) String() string {
	return strings.ToLower(fmt.Sprintf("%s %s", p.UUID, p.State))
}

// ToJSON marshals the Pipeline to JSON
func (p Pipeline) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}
