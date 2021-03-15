package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const username = "uname"
const password = "passwd"

func env() (err error) {
	err = os.Setenv("BPR_BITBUCKET_USERNAME", username)
	if err != nil {
		return err
	}

	err = os.Setenv("BPR_BITBUCKET_PASSWORD", password)
	if err != nil {
		return err
	}

	return
}

func clear() {
	os.Clearenv()
}

func TestConfigBitbucketEnv(t *testing.T) {
	env()
	config, err := LoadConfig(false)
	assert.NoError(t, err)
	assert.Equal(t, username, config.BitbucketUsername)
	assert.Equal(t, password, config.BitbucketPassword)
	clear()
}

func TestConfigOkIfFileDoesNotExist(t *testing.T) {
	env()
	config, err := LoadConfig(true)
	assert.NoError(t, err)
	assert.Equal(t, username, config.BitbucketUsername)
	assert.Equal(t, password, config.BitbucketPassword)
	clear()
}

func TestEnvCanSetConfigPath(t *testing.T) {
	os.Setenv("BPR_CONFIG_PATH", "test/fixtures/config.env")
	config, err := LoadConfig(true)
	assert.NoError(t, err)
	assert.Equal(t, username, config.BitbucketUsername)
	assert.Equal(t, password, config.BitbucketPassword)
	clear()
}
