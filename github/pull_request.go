package github

type PullRequest struct {
	Id      string
	BaseRef Ref
	HeadRef Ref
}
