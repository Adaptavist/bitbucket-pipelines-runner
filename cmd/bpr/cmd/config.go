package cmd

import (
	"github.com/adaptavist/bitbucket-pipelines-runner/v2/pkg/bitbucket/http"
	"github.com/spf13/viper"
)

func getHTTPClient() http.Client {
	bitbucketUsername = viper.GetString("BITBUCKET_USERNAME")
	fatalIfEmpty(bitbucketUsername, "BitBucket username required")
	bitbucketPassword = viper.GetString("BITBUCKET_PASSWORD")
	fatalIfEmpty(bitbucketPassword, "BitBucket password required")
	return http.Client{Auth: http.Auth{
		Username: bitbucketUsername,
		Password: bitbucketPassword,
	}}
}
