package github

import (
	"context"
)

type Client interface {
	GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error)
	CommentOnPullRequest(ctx context.Context, owner, repo string, number int, comment string) error
	HasExistingComments(ctx context.Context, owner, repo string, number int) (bool, error)
}
