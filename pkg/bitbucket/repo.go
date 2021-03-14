package bitbucket

import (
	"fmt"
)

// Repo for BitBucket
type Repo struct {
	Workspace string
	Slug      string
}

// String of the repo
func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Workspace, r.Slug)
}

// NewRepo constructs a Repo variable
func NewRepo(owner, repoSlug string) Repo {
	return Repo{
		Workspace: owner,
		Slug:      repoSlug,
	}
}
