package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ziadalikhalifa/code-review-papa/config"
	"github.com/ziadalikhalifa/code-review-papa/internal/ai"
	"github.com/ziadalikhalifa/code-review-papa/internal/analyzer"
	"github.com/ziadalikhalifa/code-review-papa/internal/github"
)

// PREvent represents the structure of a GitHub PR webhook event
type PREvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		URL string `json:"url"`
	} `json:"pull_request"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Code Review Papa Lambda")

	log.Println("Request body:", request.Body)

	// Parse the webhook payload
	var prEvent PREvent
	if err := json.Unmarshal([]byte(request.Body), &prEvent); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid webhook payload",
		}, nil
	}

	// Only process opened or synchronize (new commits) events
	if prEvent.Action != "opened" && prEvent.Action != "synchronize" && prEvent.Action != "reopened" && prEvent.Action != "edited" {
		log.Printf("Ignoring PR event with action: %s", prEvent.Action)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Event ignored - not an opened or synchronize action",
		}, nil
	}

	// Load configuration
	cfg := config.LoadConfig()
	if !cfg.Validate() {
		log.Fatal("Invalid configuration: missing required environment variables")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Server configuration error",
		}, nil
	}

	// Initialize services
	var githubClient github.Client
	var err error

	// Choose authentication method based on available credentials
	if cfg.GithubAppPrivateKey != "" {

		githubClient, err = github.NewGithubAppClient(
			cfg.GithubAppID,
			cfg.GithubAppPrivateKey,
			prEvent.Installation.ID,
		)
		if err != nil {
			log.Printf("Failed to create GitHub App client: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Failed to initialize GitHub client",
			}, nil
		}
	} else {
		// Fall back to token-based authentication
		log.Println("Using token-based authentication")
		githubClient = github.NewGithubClient(cfg.GithubToken)
	}

	aiService := ai.NewDeepSeekService(cfg.DeepSeekKey)
	prAnalyzer := analyzer.NewPRAnalyzer(githubClient, aiService)

	// Extract PR details from the event
	owner := prEvent.Repository.Owner.Login
	repo := prEvent.Repository.Name
	prNumber := prEvent.Number

	log.Printf("Analyzing PR #%d in %s/%s", prNumber, owner, repo)

	// Analyze the PR
	err = prAnalyzer.AnalyzePR(ctx, owner, repo, prNumber)
	if err != nil {
		log.Printf("Failed to analyze PR: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to analyze PR",
		}, nil
	}

	log.Println("Finished analyzing PR")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "PR analysis completed successfully",
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
