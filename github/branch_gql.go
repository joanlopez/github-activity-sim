package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func (c *Client) GetLastBranchReference(ctx context.Context, branch string) (string, error) {
	var q struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Target struct {
							Oid string
						}
					}
				}
			} `graphql:"refs(first: 1, query: $branch, refPrefix: $prefix)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":   githubv4.String(c.repoName),
		"owner":  githubv4.String(c.repoOwner),
		"branch": githubv4.String(branch),
		"prefix": githubv4.String(branchRefPrefix),
	}

	err := c.gql.Query(ctx, &q, variables)
	if err != nil {
		fmt.Printf("Branch fetching failed; err: %s\n", err.Error())
		return "", err
	}

	lastOid := q.Repository.Refs.Edges[0].Node.Target.Oid
	c.lastOid[branch] = lastOid
	fmt.Printf("Branch successfully fetched; oid: %s\n", lastOid)
	return lastOid, nil
}
