package spec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"

	"gopkg.in/yaml.v3"
)

func (s PipelineTarget) String() string {
	str := s.Workspace + "/" + s.Repo + "/" + s.RefType + "/" + s.RefName
	if s.CustomTarget != "" {
		str = str + "/" + s.CustomTarget
	}
	return str
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
func (p Pipeline) GetTarget() (PipelineTarget, error) {
	return StringToTarget(p.Pipeline)
}

// StringToTarget takes a string and builds a PipelineTarget object
func StringToTarget(str string) (target PipelineTarget, err error) {
	parts := strings.Split(str, "/")

	if len(parts) < 4 || len(parts) > 5 {
		err = fmt.Errorf("spec identifier must consists 4-5 parts (workspace/repo/ref-type/ref-name[/spec]), but got %d (%s)", len(parts), str)
		return
	}

	target = PipelineTarget{
		Workspace: parts[0],
		Repo:      parts[1],
		RefType:   parts[2],
		RefName:   parts[3],
	}

	if len(parts) == 5 {
		target.CustomTarget = parts[4]
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
