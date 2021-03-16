package cmd

import (
	"fmt"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/model"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/cmd/spec"
	"log"
	"strings"
)

func fatalIfNotNil(v error) {
	if v != nil {
		log.Fatal(v)
	}
}

func fatalIfEmpty(str, err string) {
	if strings.TrimSpace(str) == "" {
		log.Fatal(err)
	}
}

func stringsToVars(varString []string, secured bool) (vars model.Variables, err error) {
	for _, str := range varString {
		parts := strings.Split(str, "=")

		if len(parts) != 2 {
			err = fmt.Errorf("variable must comprise of a key and equals and a value")
			return
		}

		vars = append(vars, model.Variable{
			Key: strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
			Secured: secured,
		})
	}
	return
}

// NewTarget to run
func NewTarget(target string) (t spec.PipelineTarget, err error) {
	parts := strings.Split(target, "/")

	if len(parts) < 3 || len(parts) > 4 {
		err = fmt.Errorf("identifier must consists 3-4 parts (workspace/repo/branch[/spec]), but got %d (%s)", len(parts), target)
		return
	}

	t = spec.PipelineTarget{
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