package model

// TagResponse gives us enough information to do what we need to run a tag pipeline using the hash
type TagResponse struct {
	Name string `json:"name"`
	Target struct {
		Hash string `json:"hash"`
	} `json:"target"`
}