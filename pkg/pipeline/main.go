package pipeline

import (
	"bytes"
	"io/ioutil"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"gopkg.in/yaml.v3"
)

// Spec of a pipeline which will be mapped to BitBuckets API later
type Spec struct {
	Owner     string              `yaml:"owner"`
	Repo      string              `yaml:"repo"`
	Ref       string              `yaml:"ref"`
	Pipeline  string              `yaml:"pipeline"`
	Variables bitbucket.Variables `yaml:"variables"`
}

// Specs list
type Specs []Spec

// UnmarshalSpecsFile into specs
func UnmarshalSpecsFile(file string) (specs Specs, err error) {
	data, readErr := ioutil.ReadFile(file)

	if readErr != nil {
		err = readErr
		return
	}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&specs)
	return
}

// GetWorkspaceRepo from spec
func (s Spec) GetWorkspaceRepo() bitbucket.Repo {
	return bitbucket.NewRepo(s.Owner, s.Repo)
}

// GetPipeline from spec
func (s Spec) GetPipeline() bitbucket.Pipeline {
	return bitbucket.NewPipeline(bitbucket.NewBranchTarget(s.Ref, s.Pipeline), s.Variables)
}
