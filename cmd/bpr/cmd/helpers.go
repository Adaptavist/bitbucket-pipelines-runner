package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/adaptavist/bitbucket-pipelines-runner/v1/pkg/bitbucket/model"
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
			Key:     strings.TrimSpace(parts[0]),
			Value:   strings.TrimSpace(parts[1]),
			Secured: secured,
		})
	}
	return
}
