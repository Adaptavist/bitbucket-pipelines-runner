package bitbucket

import (
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/client"
	"github.com/adaptavist/bitbucket-pipeline-runner/v1/pkg/bitbucket/http"
)

// NewClient helper for BitBucket
func NewClient(http http.Client) client.Client {
	return client.NewClient(http)
}

// NewPipelineOpts helper for BitBucket
func NewPipelineOpts() client.PipelineOpts {
	return client.NewPipelineOpts()
}
