package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func (c *Client) CreateIssue(ctx context.Context, title, body string) (Issue, error) {
	var m struct {
		CreateIssue struct {
			Issue Issue
		} `graphql:"createIssue(input: $input)"`
	}

	input := githubv4.CreateIssueInput{
		RepositoryID: githubv4.ID(c.repoId),
		Title:        githubv4.String(title),
		Body:         githubv4.NewString(githubv4.String(body)),
	}

	err := c.gql.Mutate(ctx, &m, input, nil)
	if err != nil {
		fmt.Printf("Issue creation failed; err: %s\n", err.Error())
		return Issue{}, err
	}

	issueId := m.CreateIssue.Issue.Id
	fmt.Printf("Issue successfully created; id: %s\n", issueId)
	return m.CreateIssue.Issue, nil
}
