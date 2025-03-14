package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	GithubToken string
	DeepSeekKey string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		GithubToken: getEnv("GITHUB_TOKEN", ""),
		DeepSeekKey: getEnv("DEEPSEEK_KEY", ""),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() bool {
	return c.GithubToken != "" && c.DeepSeekKey != ""
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
