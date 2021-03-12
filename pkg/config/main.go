package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/http"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/utils"
	"github.com/spf13/viper"
)

// Config for the CLI to work with BitBucket
type Config struct {
	BitbucketUsername string `mapstructure:"bitbucket_username"`
	BitbucketPassword string `mapstructure:"bitbucket_password"`
}

func findLoadConfigFile(config *Config) {
	// try from the environment first
	envPath, envOK := os.LookupEnv("BPR_CONFIG_PATH")
	if envOK {
		viper.SetConfigFile(envPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
		return
	}

	// then try the home directory
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		viper.AddConfigPath(fmt.Sprintf("%s/.bpr", home))
		viper.SetConfigType("env")
		err := viper.ReadInConfig()
		if err != nil {
			// TODO: replace with better log levels
			log.Println(err.Error())
		}
		return
	}
}

// LoadConfig using ~/.bpr/config.env or environment
func LoadConfig(loadFile bool) (config Config, err error) {
	viper.SetEnvPrefix("BPR")
	viper.AutomaticEnv()
	_ = viper.BindEnv("bitbucket_username", "BITBUCKET_USERNAME")
	_ = viper.BindEnv("bitbucket_password", "BITBUCKET_PASSWORD", "BITBUCKET_APP_PASSWORD")

	if loadFile {
		findLoadConfigFile(&config)
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Println(err)
		return
	}

	if utils.Empty(config.BitbucketUsername) {
		err = errors.New("BITBUCKET_USERNAME required")
		return
	}

	if utils.Empty(config.BitbucketPassword) {
		err = errors.New("BITBUCKET_PASSWORD required")
		return
	}

	return
}

// LoadConfigOrPanic does exactly what it says
func LoadConfigOrPanic(loadFile bool) Config {
	config, err := LoadConfig(loadFile)

	if err != nil {
		panic(err)
	}

	return config
}

// GetAuth from the config
func (c Config) GetAuth() http.Auth {
	return http.Auth{
		Username: c.BitbucketUsername,
		Password: c.BitbucketPassword,
	}
}
