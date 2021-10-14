package cmd

import (
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/spf13/viper"
)

func makeClient() (c client.Client) {
	config := client.Config{}

	config.Username = viper.GetString("BITBUCKET_USERNAME")
	fatalIfEmpty(config.Username, "BitBucket username required")

	config.Password = viper.GetString("BITBUCKET_PASSWORD")
	fatalIfEmpty(config.Password, "BitBucket password required")

	baseURL := viper.GetString("BITBUCKET_BASE_URL")
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	c = client.Client{
		Config: &config,
	}
	return
}
