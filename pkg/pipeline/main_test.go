package pipeline

import (
	"testing"
)

// TestUnmarshalSpecsFromFile - does what it says
func TestUnmarshalSpecsFromFile(t *testing.T) {
	specs, err := UnmarshalSpecsFile("test/fixtures/pipelines.yml")

	if err != nil {
		t.Error(err.Error())
	}

	if len(specs) == 0 {
		t.Error("Expected at least one spec")
	}
}
