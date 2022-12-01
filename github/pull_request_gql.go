package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
)

func (c *Client) CreatePullRequest(ctx context.Context, title, from string, to ...string) (PullRequest, error) {
	if len(to) > 1 {
		return PullRequest{}, errors.New("single base ref supported")
	}

	var m struct {
		CreatePullRequest struct {
			PullRequest PullRequest
		} `graphql:"createPullRequest(input: $input)"`
	}

	baseRef := c.defBranch
	if len(to) == 1 {
		baseRef = to[0]
	}

	input := githubv4.CreatePullRequestInput{
		RepositoryID: githubv4.ID(c.repoId),
		BaseRefName:  githubv4.String(baseRef),
		HeadRefName:  githubv4.String(from),
		Title:        githubv4.String(title),
	}

	err := c.gql.Mutate(ctx, &m, input, nil)
	if err != nil {
		fmt.Printf("Pull request creation failed; err: %s\n", err.Error())
		if isPullRequestAlreadyExistsErr(err) {
			pr, err := c.GetPullRequest(ctx, from, to...)
			if err != nil {
				return PullRequest{}, err
			}
			return pr, nil
		}
		return PullRequest{}, err
	}

	pullRequestId := m.CreatePullRequest.PullRequest.Id
	fmt.Printf("Pull request successfully created; id: %s\n", pullRequestId)
	return m.CreatePullRequest.PullRequest, nil
}

func (c *Client) GetPullRequest(ctx context.Context, from string, to ...string) (PullRequest, error) {
	if len(to) > 1 {
		return PullRequest{}, errors.New("single base ref supported")
	}

	var q struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node PullRequest
				}
			} `graphql:"pullRequests(first: 1, baseRefName: $baseRef, headRefName: $headRef)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	baseRef := c.defBranch
	if len(to) == 1 {
		baseRef = to[0]
	}

	variables := map[string]interface{}{
		"name":    githubv4.String(c.repoName),
		"owner":   githubv4.String(c.repoOwner),
		"baseRef": githubv4.String(baseRef),
		"headRef": githubv4.String(from),
	}

	err := c.gql.Query(ctx, &q, variables)
	if err != nil {
		fmt.Printf("Pul request fetching failed; err: %s\n", err.Error())
		return PullRequest{}, err
	}

	pullRequestId := q.Repository.PullRequests.Edges[0].Node.Id
	fmt.Printf("Pull request successfully fetched; id: %s\n", pullRequestId)
	return q.Repository.PullRequests.Edges[0].Node, nil
}

func (c *Client) MergePullRequest(ctx context.Context, opts []MergePullRequestOption, from string, to ...string) error {
	pr, err := c.GetPullRequest(ctx, from, to...)
	if err != nil {
		return err
	}

	var m struct {
		MergePullRequest struct {
			ClientMutationId string
		} `graphql:"mergePullRequest(input: $input)"`
	}

	input := githubv4.MergePullRequestInput{
		PullRequestID: githubv4.ID(pr.Id),
	}

	for _, opt := range opts {
		input = opt(input)
	}

	err = c.gql.Mutate(ctx, &m, input, nil)
	if err != nil {
		fmt.Printf("Pull request merge failed; err: %s\n", err.Error())
		return err
	}

	fmt.Printf("Pull request successfully merged; id: %s\n", pr.Id)

	if err := c.DeleteBranch(ctx, pr.HeadRef.Name); err != nil {
		return err
	}

	return nil
}

func isPullRequestAlreadyExistsErr(err error) bool {
	return strings.Contains(
		err.Error(),
		"A pull request already exists for",
	)
}

type MergePullRequestOption func(githubv4.MergePullRequestInput) githubv4.MergePullRequestInput

type PullRequestMergeMethod githubv4.PullRequestMergeMethod

const (
	PullRequestMergeMethodMerge  = PullRequestMergeMethod(githubv4.PullRequestMergeMethodMerge)
	PullRequestMergeMethodSquash = PullRequestMergeMethod(githubv4.PullRequestMergeMethodSquash)
	PullRequestMergeMethodRebase = PullRequestMergeMethod(githubv4.PullRequestMergeMethodRebase)
)

func MergePullRequestMethod(m PullRequestMergeMethod) MergePullRequestOption {
	return func(i githubv4.MergePullRequestInput) githubv4.MergePullRequestInput {
		mergeMethod := githubv4.PullRequestMergeMethod(m)
		i.MergeMethod = &mergeMethod
		return i
	}
}
