package analyzer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ziadalikhalifa/code-review-papa/internal/ai"
	"github.com/ziadalikhalifa/code-review-papa/internal/github"
)

// PRAnalyzer handles the analysis of pull requests
type PRAnalyzer struct {
	githubClient github.Client
	aiService    ai.Service
}

// NewPRAnalyzer creates a new PR analyzer
func NewPRAnalyzer(githubClient github.Client, aiService ai.Service) *PRAnalyzer {
	return &PRAnalyzer{
		githubClient: githubClient,
		aiService:    aiService,
	}
}

// AnalyzePR analyzes a pull request and posts a comment with the analysis
func (a *PRAnalyzer) AnalyzePR(ctx context.Context, owner, repo string, prNumber int) error {
	log.Printf("Analyzing PR #%d in %s/%s", prNumber, owner, repo)

	// Check if we've already commented on this PR
	hasComments, err := a.githubClient.HasExistingComments(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to check for existing comments: %w", err)
	}

	if hasComments {
		log.Printf("Already commented on PR #%d, skipping", prNumber)
		return nil
	}

	// Get the PR diff
	diff, err := a.githubClient.GetPullRequestDiff(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get PR diff: %w", err)
	}

	// Skip if diff is too large or empty
	if len(diff) == 0 {
		log.Printf("Empty diff for PR #%d, skipping", prNumber)
		return nil
	}

	if len(diff) > 100000 {
		log.Printf("Diff too large for PR #%d (%d bytes), skipping", prNumber, len(diff))
		comment := "⚠️ **Code Review Papa**: This PR is too large for automated review. Consider breaking it into smaller PRs for better feedback."
		return a.githubClient.CommentOnPullRequest(ctx, owner, repo, prNumber, comment)
	}

	// Analyze the diff
	analysis, err := a.aiService.AnalyzeCode(ctx, diff)
	if err != nil {
		return fmt.Errorf("failed to analyze code: %w", err)
	}

	// Format the comment
	comment := formatComment(analysis)

	// Post the comment
	// Post the comment
	err = a.githubClient.CommentOnPullRequest(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to comment on PR: %w", err)
	}

	log.Printf("Successfully analyzed and commented on PR #%d", prNumber)
	return nil
}

// formatComment formats the AI analysis into a nice GitHub comment
func formatComment(analysis string) string {
	header := `# 🧙‍♂️ Code Review Papa

I've analyzed this pull request and have some feedback for you!

`
	footer := `

---
*This automated review was generated by [Code Review Papa](https://github.com/ziadalikhalifa/code-review-papa). If you find this helpful, please give it a ⭐!*
`

	// Ensure the analysis doesn't have duplicate headers
	if strings.Contains(strings.ToLower(analysis), "code review") {
		// The AI might have added its own header, so we'll skip our header
		return analysis + footer
	}

	return header + analysis + footer
}
