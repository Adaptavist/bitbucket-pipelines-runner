package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// State of a Pipeline
type State struct {
	Name   string `json:"name"`
	Result Result `json:"result"`
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

// ResultError so we can extract a client specific error
type ResultError struct {
	Message string `json:"message"`
}

// Result shows if the step completed successful
type Result struct {
	Name  string      `json:"name"`
	Error ResultError `json:"error"`
}

// String representation of a Result
func (s Result) String() string {
	return s.Name
}

// HasError does what is said
func (s Result) HasError() bool {
	return s.Error != ResultError{}
}

// OK does what it says
func (s Result) OK() bool {
	return s.Name == "SUCCESSFUL"
}
