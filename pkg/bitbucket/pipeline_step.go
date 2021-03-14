package bitbucket

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

// String representation of a PipelineStep
func (s PipelineStep) String() string {
	return s.Name
}

// PipelineStepStateResultError so we can extract a bitbucket specific error
type PipelineStepStateResultError struct {
	Message string `json:"message"`
}

// PipelineStepStateResult shows if the step completed successful
type PipelineStepStateResult struct {
	Name  string                       `json:"name"`
	Error PipelineStepStateResultError `json:"error"`
}

// String representation of a PipelineStepStateResult
func (s PipelineStepStateResult) String() string {
	return s.Name
}

// HasError does what is said
func (r PipelineStepStateResult) HasError() bool {
	return r.Error != PipelineStepStateResultError{}
}

// OK does what it says
func (s PipelineStepStateResult) OK() bool {
	return s.Name == "SUCCESSFUL"
}

// PipelineStepsResponse response
type PipelineStepsResponse struct {
	Page   int            `json:"page"`
	Values []PipelineStep `json:"values"`
	Next   string         `json:"next"`
}

// PipelineSteps is a list of steps
type PipelineSteps []PipelineStep

// FilterSteps filters are list of steps
func FilterSteps(steps PipelineSteps, test func(PipelineStep) bool) (ret PipelineSteps) {
	for _, s := range steps {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}