package spec

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultsPipeline(t *testing.T) {
	spec, err := UnmarshalSpec(`
pipelines:
  test:
    pipeline: owner/repo/branch`)

	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, spec.Pipelines["test"], "test should not be nil")

	opts, err := spec.MakePipelineOpts("test")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "owner", opts.Repo.Workspace, "expected owner")
	assert.Equal(t, "repo", opts.Repo.Slug, "expected repo")
	assert.Equal(t, "branch", opts.Target.RefName, "expected ref")
	assert.Equal(t, "default", opts.Target.Selector.Pattern, "expected default")
}

func TestSpecifyPipelineName(t *testing.T) {
	spec, err := UnmarshalSpec(`
pipelines:
  test:
    pipeline: owner/repo/branch/pipeline`)

	assert.Nil(t, err, "err should be nil")
	opts, _ := spec.MakePipelineOpts("test")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, "pipeline", opts.Target.Selector.Pattern, "expected pipeline")
}

func TestGlobalVars(t *testing.T) {
	spec, err := UnmarshalSpec(`
variables:
  var_1: value_1
pipelines:
  test:
    pipeline: owner/repo/branch`)

	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, spec, "spec should not be nil")
	assert.NotNil(t, spec.Pipelines, "pipelines should not be nil")

	opts, err := spec.MakePipelineOpts("test")
	assert.Nil(t, err, "err should be nil")
	fmt.Print(opts.Variables)
	assert.Len(t, opts.Variables, 1, "pipeline should have one variable, go %s")

}

func TestGlobalVarsGetMerged(t *testing.T) {
	spec, err := UnmarshalSpec(`
variables:
  var_1: value_1
pipelines:
  test:
    pipeline: owner/repo/branch
    variables:
      var_2: value_2`)

	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, spec, "spec should not be nil")
	assert.NotNil(t, spec.Pipelines, "pipelines should not be nil")

	opts, err := spec.MakePipelineOpts("test")
	assert.Nil(t, err, "err should be nil")
	assert.Len(t, opts.Variables, 2, "pipeline should have two variables")

}
