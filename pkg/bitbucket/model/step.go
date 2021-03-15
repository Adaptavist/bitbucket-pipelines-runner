package model

import "encoding/json"

// StepState provides the state the current spec step is in
type StepState struct {
	Name   string          `json:"name"`
	Result StepStateResult `json:"result"`
}

// IsCompleted is true if the Step's StepState is complete
func (s StepState) IsCompleted() bool {
	return s.Name == "COMPLETED"
}

// Step on a running/complete spec
type Step struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	State StepState `json:"state"`
}

// String representation of a Step
func (s Step) String() string {
	return s.Name
}

// ToJSON marshals the Step to JSON
func (s Step) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// StepStateResultError so we can extract a client specific error
type StepStateResultError struct {
	Message string `json:"message"`
}

// StepStateResult shows if the step completed successful
type StepStateResult struct {
	Name  string               `json:"name"`
	Error StepStateResultError `json:"error"`
}

// String representation of a StepStateResult
func (s StepStateResult) String() string {
	return s.Name
}

// HasError does what is said
func (s StepStateResult) HasError() bool {
	return s.Error != StepStateResultError{}
}

// OK does what it says
func (s StepStateResult) OK() bool {
	return s.Name == "SUCCESSFUL"
}

// StepsResponse response
type StepsResponse struct {
	Page   int    `json:"page"`
	Values []Step `json:"values"`
	Next   string `json:"next"`
}

// Steps is a list of steps
type Steps []Step

// Filter Steps using a callable function
func (s Steps) Filter(test func(step Step) bool) Steps {
	return FilterSteps(s, test)
}

// Filter Steps using a callable function
func FilterSteps(steps Steps, test func(Step) bool) (results Steps) {
	for _, s := range steps {
		if test(s) {
			results = append(results, s)
		}
	}
	return
}
