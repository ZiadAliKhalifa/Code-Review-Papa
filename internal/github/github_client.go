package github

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

// NewGithubAppClient creates a new GitHub client using GitHub App authentication
func NewGithubAppClient(appID int64, privateKey string, installationID int64) (*GitHubClient, error) {
	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create a JWT for GitHub App authentication
	jwtToken, err := createJWT(appID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT: %w", err)
	}

	// Create a temporary client to get an installation token
	httpClient := &http.Client{}

	// Create a request to get an installation token
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "code-review-papa")

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get installation token, status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse installation token response: %w", err)
	}

	// Create a new client with the installation token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenResp.Token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GitHubClient{
		client: client,
	}, nil
}

// createJWT creates a JWT for GitHub App authentication
func createJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	// JWT expiration time: 10 minutes from now
	now := time.Now()
	expirationTime := now.Add(10 * time.Minute)

	// Create the JWT claims
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    fmt.Sprintf("%d", appID),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
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
