package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestChDir (t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err, "couldn't get cwd")
	rootCmd.SetArgs([]string{"spec", "--chdir", "../../../test/default", "--dry"})
	err = rootCmd.Execute()
	assert.Nil(t, err, "execute shouldn't error")
	err = os.Chdir(dir)
	assert.Nil(t, err, "couldn't reset dir")
}

func TestDuplicateSpecs (t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err, "couldn't get cwd")
	rootCmd.SetArgs([]string{"spec", "--chdir", "../../../test/duplicates", "--dry"})
	err = rootCmd.Execute()
	assert.NotNil(t, err, "execute should error")
	err = os.Chdir(dir)
	assert.Nil(t, err, "couldn't reset dir")
}