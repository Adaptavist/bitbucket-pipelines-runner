package model

import "encoding/json"

// Step on a running/complete spec
type Step struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	State State  `json:"state"`
}

// String representation of a Step
func (s Step) String() string {
	return s.Name
}

// ToJSON marshals the Step to JSON
func (s Step) ToJSON() ([]byte, error) {
	return json.Marshal(s)
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
