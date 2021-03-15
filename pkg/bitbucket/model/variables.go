package model

// Variable for a BitBucket spec
type Variable struct {
	Key     string `json:"key" yaml:"key"`
	Value   string `json:"value" yaml:"value"`
	Secured bool   `json:"secured" yaml:"secured"`
}

// Variables for a BitBucket spec
type Variables []Variable

// Append variables
func AppendVariables(vars Variables, add Variables) Variables {
	return append(vars, add...)
}
