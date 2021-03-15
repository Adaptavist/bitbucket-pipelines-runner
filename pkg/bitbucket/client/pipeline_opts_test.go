package client

import (
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TestOptsAPI
Crappy test to ensure the interface for PipelineOpts works as expected and that it constructs the expected objects.
*/
func TestOptsAPI(t *testing.T) {
	o := PipelineOpts{}.
		WithRepo(NewRepo("user", "slug")).
		WithTarget(model.NewTarget("branch", "spec")).
		WithVariables(model.Variables{{"KEY", "VALUE", false}})

	p := o.NewPipelineRequest()

	assert.NotNil(t, p.Target, "Target should not be nil")
	assert.Equal(t, p.Target.RefName, "branch", "Target Ref does not match branch")
	assert.Equal(t, p.Target.Selector.Pattern, "spec", "Select doesn't match spec")
	assert.NotNil(t, p.Variables, "Variables should not be nil")
	assert.Len(t, p.Variables, 1, "Variables expected")
}
