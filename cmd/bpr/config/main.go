package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adaptavist/bitbucket_pipelines_runner/cmd/bpr/utils"
	"github.com/spf13/viper"
)

// Config for the CLI
type Config struct {
	BitbucketUsername string `mapstructure:"bitbucket_username"`
	BitbucketPassword string `mapstructure:"bitbucket_password"`
}

func findLoadConfigFile() {
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
	if homeErr == nil {
		viper.AddConfigPath(fmt.Sprintf("%s/.bpr", home))
		viper.SetConfigType("env")
		err := viper.ReadInConfig()
		if err != nil {
			log.Print(err.Error())
		}
		return
	} else {
		log.Print(homeErr.Error())
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
		findLoadConfigFile()
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatal(err)
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