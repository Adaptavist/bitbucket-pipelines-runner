package pipeline

import (
	"bytes"
	"io/ioutil"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket"
	"gopkg.in/yaml.v3"
)

// Spec of a pipeline which will be mapped to BitBuckets API later
type Spec struct {
	Owner     string                      `yaml:"owner"`
	Repo      string                      `yaml:"repo"`
	Ref       string                      `yaml:"ref"`
	Pipeline  string                      `yaml:"pipeline"`
	Variables bitbucket.PipelineVariables `yaml:"variables"`
}

func (s Spec) String() string {
	return s.Owner + "/" + s.Repo + ":" + s.Ref + ":" + s.Pipeline
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
