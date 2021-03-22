package urls

import (
	"fmt"
)

const baseURL = "https://api.bitbucket.org/2.0"

func PipelineWeb(workspace, slug string, buildNumber int) string {
	return fmt.Sprintf("https://bitbucket.org/%s/%s/addon/pipelines/home#!/results/%d", workspace, slug, buildNumber)
}

func Pipelines(workspace, slug string) string {
	return fmt.Sprintf("%s/repositories/%s/%s/pipelines/", baseURL, workspace, slug)
}

func Pipeline(workspace, slug, UUID string) string {
	return Pipelines(workspace, slug) + UUID
}

func PipelineSteps(workspace, slug, UUID string) string {
	return fmt.Sprintf("%s/steps/", Pipeline(workspace, slug, UUID))
}

func PipelineStepLogs(workspace, slug, pipelineUUID, stepUUID string) string {
	return fmt.Sprintf("%s%s/log", PipelineSteps(workspace, slug, pipelineUUID), stepUUID)
}

func Tag(workspace, slug, tag string) string {
	return fmt.Sprintf("%s/repositories/%s/%s/refs/tags/%s", baseURL, workspace, slug, tag)
}