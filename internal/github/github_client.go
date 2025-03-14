package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
}

// NewGithubClient creates a new GitHub client with the provided token
func NewGithubClient(token string) *GitHubClient {
	// Create an OAuth2 client with the token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	// Create a new GitHub client using the OAuth2 client
	client := github.NewClient(tc)

	return &GitHubClient{
		client: client,
	}
}

func (g *GitHubClient) GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	// Get the raw diff for the PR
	diff, _, err := g.client.PullRequests.GetRaw(
		ctx,
		owner,
		repo,
		number,
		github.RawOptions{Type: github.Diff},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get PR diff: %w", err)
	}
	return diff, nil
}

// CommentOnPullRequest adds a comment to a PR
func (g *GitHubClient) CommentOnPullRequest(ctx context.Context, owner, repo string, number int, comment string) error {
	// Create a comment object with the provided text
	commentBody := &github.IssueComment{
		Body: github.String(comment),
	}

	// Post the comment to the PR
	_, _, err := g.client.Issues.CreateComment(ctx, owner, repo, number, commentBody)
	if err != nil {
		return fmt.Errorf("failed to comment on PR: %w", err)
	}
	return nil
}

// HasExistingComments checks if the bot has already commented on the PR
func (g *GitHubClient) HasExistingComments(ctx context.Context, owner, repo string, number int) (bool, error) {
	// Get all comments on the PR
	comments, _, err := g.client.Issues.ListComments(ctx, owner, repo, number, nil)
	if err != nil {
		return false, fmt.Errorf("failed to list PR comments: %w", err)
	}

	// Look for comments that might be from our bot
	botSignature := "Code Review Papa"
	for _, comment := range comments {
		if comment.Body != nil && strings.Contains(*comment.Body, botSignature) {
			return true, nil
		}
	}
	return false, nil
}
