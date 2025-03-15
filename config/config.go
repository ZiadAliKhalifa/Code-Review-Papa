package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	GithubToken         string
	DeepSeekKey         string
	GithubAppID         int64
	GithubAppPrivateKey string
}

func LoadConfig() *Config {
	privateKey := getEnv("GITHUB_APP_PRIVATE_KEY", "")

	if privateKey != "" && !strings.Contains(privateKey, "-----BEGIN RSA PRIVATE KEY-----\n") {
		privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")
	}

	// Check if we should load from a file instead
	keyPath := getEnv("GITHUB_APP_PRIVATE_KEY_PATH", "")
	if privateKey == "" && keyPath != "" {
		keyData, err := os.ReadFile(keyPath)
		if err == nil {
			privateKey = string(keyData)
		}
	}

	return &Config{
		GithubToken:         getEnv("GITHUB_TOKEN", ""),
		DeepSeekKey:         getEnv("DEEPSEEK_KEY", ""),
		GithubAppID:         getEnvAsInt64("GITHUB_APP_ID", 0),
		GithubAppPrivateKey: privateKey,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() bool {
	// Check if we have either token-based auth or app-based auth
	hasTokenAuth := c.GithubToken != ""
	hasAppAuth := c.GithubAppPrivateKey != ""

	return (hasTokenAuth || hasAppAuth) && c.DeepSeekKey != ""
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt64 gets an environment variable as int64 or returns a default value
func getEnvAsInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return intValue
}
