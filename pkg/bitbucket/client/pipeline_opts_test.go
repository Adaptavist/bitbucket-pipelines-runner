package client

import (
	"testing"

	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/bitbucket/model"
	"github.com/stretchr/testify/assert"
)

/*
TestOptsAPI
Crappy test to ensure the interface for PipelineOpts works as expected and that it constructs the expected objects.
*/
func TestOptsAPI(t *testing.T) {
	trgt := model.NewTarget("branch", "example")
	trgt.WithCustomTarget("spec")
	o := PipelineOpts{}.
		WithRepo(NewRepo("user", "slug")).
		WithTarget(trgt).
		WithVariables(model.Variables{{Key: "KEY", Value: "VALUE", Secured: false}})

	p := o.NewPipelineRequest()

	assert.NotNil(t, p.Target, "Target should not be nil")
	assert.Equal(t, p.Target.RefName, "example", "Target Ref does not match example")
	assert.Equal(t, p.Target.Selector.Pattern, "spec", "Select doesn't match spec")
	assert.NotNil(t, p.Variables, "Variables should not be nil")
	assert.Len(t, p.Variables, 1, "Variables expected")
}
