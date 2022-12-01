package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v48/github"
)

func (c *Client) CreateBranch(ctx context.Context, branch string) error {
	_, _, err := c.http.Git.CreateRef(
		ctx,
		c.repoOwner,
		c.repoName,
		&github.Reference{
			Ref: github.String(branchRef(branch)),
			Object: &github.GitObject{
				SHA: github.String(c.lastOid[c.defBranch]),
			},
		})

	if err != nil {
		fmt.Printf("Branch creation failed; err: %s\n", err.Error())
		if isBranchAlreadyExistsErr(err) {
			if _, err := c.GetLastBranchReference(ctx, branch); err != nil {
				return err
			}
			return ErrBranchAlreadyExists
		}

		return err
	}

	c.lastOid[branch] = c.lastOid[c.defBranch]
	fmt.Printf("Branch successfully created; name: %s\n", branch)
	return nil
}

func (c *Client) DeleteBranch(ctx context.Context, branch string) error {
	_, err := c.http.Git.DeleteRef(
		ctx,
		c.repoOwner,
		c.repoName,
		branchRef(branch),
	)

	if err != nil {
		fmt.Printf("Branch deletion failed; err: %s\n", err.Error())
	}

	fmt.Printf("Branch successfully deleted; name: %s\n", branch)
	return nil
}

func isBranchAlreadyExistsErr(err error) bool {
	var errRes *github.ErrorResponse
	return errors.As(err, &errRes) &&
		errRes.Response.StatusCode == http.StatusUnprocessableEntity &&
		errRes.Message == "Reference already exists"
}
