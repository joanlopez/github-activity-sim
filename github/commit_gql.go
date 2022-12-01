package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func (c Client) CreateCommit(ctx context.Context, branch, msg string) error {
	if _, ok := c.lastOid[branch]; !ok {
		if _, err := c.GetLastBranchReference(ctx, branch); err != nil {
			return err
		}
	}

	var m struct {
		CreateCommitOnBranch struct {
			Commit struct {
				Oid string
			}
		} `graphql:"createCommitOnBranch(input: $input)"`
	}

	input := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(
				githubv4.String(fmt.Sprintf("%s/%s", c.repoOwner, c.repoName)),
			),
			BranchName: githubv4.NewString(
				githubv4.String(branchRef(branch)),
			),
		},
		Message:         githubv4.CommitMessage{Headline: githubv4.String(msg)},
		ExpectedHeadOid: githubv4.GitObjectID(c.lastOid[branch]),
		FileChanges: &githubv4.FileChanges{
			Additions: nil,
			Deletions: nil,
		},
	}

	err := c.gql.Mutate(ctx, &m, input, nil)
	if err != nil {
		return err
	}

	lastOid := m.CreateCommitOnBranch.Commit.Oid
	c.lastOid[branch] = lastOid
	fmt.Printf("Commit successfully created; oid: %s\n", lastOid)

	return nil
}
