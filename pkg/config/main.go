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

// LoadConfig using ~/.bpr/config.env or environment
func LoadConfig() (config Config, err error) {
	home, homeErr := os.UserHomeDir()
	if homeErr == nil {
		viper.AddConfigPath(fmt.Sprintf("%s/.bpr", home))
		viper.SetConfigType("env")
	}
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		log.Println(err)
	}

	err = viper.Unmarshal(&config)

	if utils.Empty(config.BitbucketUsername) {
		err = errors.New("BITBUCKET_USERNAME required")
	}

	if utils.Empty(config.BitbucketPassword) {
		err = errors.New("BITBUCKET_PASSWORD required")
	}

	return
}

// LoadConfigOrPanic does exactly what it says
func LoadConfigOrPanic() Config {
	config, err := LoadConfig()

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
