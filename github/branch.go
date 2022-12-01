package github

import (
	"errors"
	"fmt"
)

var ErrBranchAlreadyExists = errors.New("branch already exists")

const branchRefPrefix = "refs/heads/"

func branchRef(branch string) string {
	return fmt.Sprintf("%s%s", branchRefPrefix, branch)
}
