package main

import (
	"context"
	"log"
	"os"

	"github.com/ziadalikhalifa/code-review-papa/config"
	"github.com/ziadalikhalifa/code-review-papa/internal/ai"
	"github.com/ziadalikhalifa/code-review-papa/internal/analyzer"
	"github.com/ziadalikhalifa/code-review-papa/internal/github"
)

func main() {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Code Review Papa")

	// Load configuration
	cfg := config.LoadConfig()
	if !cfg.Validate() {
		log.Fatal("Invalid configuration: missing required environment variables")
	}

	// Initialize services
	var githubClient github.Client
	var err error

	// Choose authentication method based on available credentials
	if cfg.GithubAppID != 0 && cfg.GithubAppPrivateKey != "" && cfg.GithubAppInstallationID != 0 {
		// Use GitHub App authentication
		log.Println("Using GitHub App authentication")
		githubClient, err = github.NewGithubAppClient(
			cfg.GithubAppID,
			cfg.GithubAppPrivateKey,
			cfg.GithubAppInstallationID,
		)
		if err != nil {
			log.Fatalf("Failed to create GitHub App client: %v", err)
		}
	} else {
		// Fall back to token-based authentication
		log.Println("Using token-based authentication")
		githubClient = github.NewGithubClient(cfg.GithubToken)
	}

	aiService := ai.NewDeepSeekService(cfg.DeepSeekKey)
	prAnalyzer := analyzer.NewPRAnalyzer(githubClient, aiService)

	// For testing, you can hardcode a PR to analyze
	// Replace these values with a real PR you want to test with
	owner := "ziadalikhalifa"
	repo := "Fyyur"
	prNumber := 1 // PR number to analyze

	// Create a context
	ctx := context.Background()

	// Analyze the PR
	err = prAnalyzer.AnalyzePR(ctx, owner, repo, prNumber)
	if err != nil {
		log.Fatalf("Failed to analyze PR: %v", err)
	}

	log.Println("Finished analyzing PR")
}
