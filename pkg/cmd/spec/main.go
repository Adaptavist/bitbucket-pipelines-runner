package spec

import (
	"bytes"
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

func (s PipelineTarget) String() string {
	return s.Workspace + "/" + s.Repo + "/" + s.Ref + "/" + s.Pipeline
}

func UnmarshalSpec(specStr string) (spec Spec, err error) {
	dec := yaml.NewDecoder(strings.NewReader(specStr))
	err = dec.Decode(&spec)
	return
}

// UnmarshalSpecsFile into specs
func UnmarshalSpecsFile(file string) (spec Spec, err error) {
	data, readErr := ioutil.ReadFile(file)

	if readErr != nil {
		err = readErr
		return
	}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&spec)
	return
}

// GetTarget from the spec string in the YAML config
func (p Pipeline) GetTarget() (t PipelineTarget, err error) {
	parts := strings.Split(p.Pipeline, "/")

	if len(parts) < 3 || len(parts) > 4 {
		err = fmt.Errorf("spec identifier must consists 3-4 parts (workspace/repo/branch[/spec]), but got %d (%s)", len(parts), p.Pipeline)
		return
	}

	t = PipelineTarget{
		Workspace: parts[0],
		Repo:      parts[1],
		Ref:       parts[2],
		Pipeline:  "default",
	}

	if len(parts) == 4 {
		t.Pipeline = parts[3]
		return
	}

	return
}

func (s Spec) ToBitbucketVariables(v Variables) model.Variables {
	vars := model.Variables{}
	for key, value := range v {
		vars = append(vars, model.Variable{
			Key:   key,
			Value: value,
		})
	}
	return vars
}
