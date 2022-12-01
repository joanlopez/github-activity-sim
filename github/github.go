package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
)

type Client struct {
	repoId    string
	repoName  string
	repoOwner string

	defBranch string
	lastOid   map[string]string

	gql  *githubv4.Client
	http *github.Client
}

func NewClient(ctx context.Context, name, owner, defBranch string, oauth2 *http.Client) (*Client, error) {
	client := &Client{
		repoName:  name,
		repoOwner: owner,
		defBranch: defBranch,
		lastOid:   make(map[string]string),
		gql:       githubv4.NewClient(oauth2),
		http:      github.NewClient(oauth2),
	}

	if err := client.init(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) init(ctx context.Context) error {
	var q struct {
		Repository struct {
			Id               string
			DefaultBranchRef struct {
				Name   string
				Target struct {
					Commit struct {
						History struct {
							Nodes []struct {
								Oid string
							}
						} `graphql:"history(first: 1)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":  githubv4.String(c.repoName),
		"owner": githubv4.String(c.repoOwner),
	}

	err := c.gql.Query(ctx, &q, variables)
	if err != nil {
		fmt.Printf("Client initialization failed; err: %s\n", err.Error())
		return err
	}

	c.repoId = q.Repository.Id
	c.defBranch = q.Repository.DefaultBranchRef.Name
	c.lastOid[c.defBranch] = q.Repository.DefaultBranchRef.Target.Commit.History.Nodes[0].Oid

	fmt.Printf(
		"Client successfully initialized; repo_id: %s; branch: %s, main oid: %s\n",
		c.repoId, c.defBranch, c.lastOid[c.defBranch],
	)

	return nil
}
