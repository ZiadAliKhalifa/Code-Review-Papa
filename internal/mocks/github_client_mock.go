package mocks

import (
	"context"
	"fmt"

	"github.com/ziadalikhalifa/code-review-papa/internal/github"
)

var _ github.Client = &MockGithubClient{}

// MockGithubClient is a mock implementation of the github.Client interface.
type MockGithubClient struct {
	GetPullRequestDiffFunc  func(ctx context.Context, owner, repo string, number int) (string, error)
	CommentOnPullRequestFunc func(ctx context.Context, owner, repo string, number int, comment string) error
	HasExistingCommentsFunc  func(ctx context.Context, owner, repo string, number int) (bool, error)
}

// GetPullRequestDiff calls the GetPullRequestDiffFunc field.
func (m *MockGithubClient) GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	if m.GetPullRequestDiffFunc == nil {
		return "", fmt.Errorf("MockGithubClient.GetPullRequestDiffFunc is not set")
	}
	return m.GetPullRequestDiffFunc(ctx, owner, repo, number)
}

// CommentOnPullRequest calls the CommentOnPullRequestFunc field.
func (m *MockGithubClient) CommentOnPullRequest(ctx context.Context, owner, repo string, number int, comment string) error {
	if m.CommentOnPullRequestFunc == nil {
		return fmt.Errorf("MockGithubClient.CommentOnPullRequestFunc is not set")
	}
	return m.CommentOnPullRequestFunc(ctx, owner, repo, number, comment)
}

// HasExistingComments calls the HasExistingCommentsFunc field.
func (m *MockGithubClient) HasExistingComments(ctx context.Context, owner, repo string, number int) (bool, error) {
	if m.HasExistingCommentsFunc == nil {
		return false, fmt.Errorf("MockGithubClient.HasExistingCommentsFunc is not set")
	}
	return m.HasExistingCommentsFunc(ctx, owner, repo, number)
}
