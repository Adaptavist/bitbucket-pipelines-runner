package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PanicIfNotNil - Does what it says
func PanicIfNotNil(v error) {
	if v != nil {
		fmt.Println("PANIC!")
		panic(v)
	}
}

// Trim as string
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// Empty check
func Empty(str string) bool {
	return len(Trim(str)) == 0
}

// DefaultWhenEmpty checks if a value is empty otherwise returns a default
func DefaultWhenEmpty(expected string, def string) string {
	if Empty(expected) {
		return def
	}
	return expected
}

// Marshall object into a JSON string
func Marshall(v interface{}) string {
	str, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(str)
}

// MarshalFormatted marshalls an object into a formatted JSON string
func MarshalFormatted(v interface{}) string {
	str, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(str)
}
