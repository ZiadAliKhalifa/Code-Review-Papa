package ai

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Service defines the interface for AI operations
type Service interface {
	AnalyzeCode(ctx context.Context, diff string) (string, error)
}

// DeepSeekService implements the Service interface using DeepSeek
type DeepSeekService struct {
	client      *resty.Client
	apiKey      string
	modelName   string
	maxTokens   int
	temperature float64
}

// NewDeepSeekService creates a new DeepSeek service
func NewDeepSeekService(apiKey string) *DeepSeekService {
	return &DeepSeekService{
		client:      resty.New(),
		apiKey:      apiKey,
		modelName:   "deepseek-coder",
		maxTokens:   1000,
		temperature: 0.7,
	}
}

// Define the request and response structures for the DeepSeek API
type deepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []deepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens"`
	Temperature float64           `json:"temperature"`
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// AnalyzeCode sends the diff to DeepSeek for analysis
func (s *DeepSeekService) AnalyzeCode(ctx context.Context, diff string) (string, error) {
	// Create the prompt for the AI
	prompt := fmt.Sprintf(`
				You are a senior software engineer reviewing a pull request. 
				Analyze the following code diff and provide:
				1. A concise summary of the changes
				2. Potential issues or bugs
				3. Suggestions for improvements
				4. Any security concerns
				5. Code quality feedback

				Format your response in markdown with clear sections.

				Here's the diff:
				%s
	`, diff)

	// Create the request to the DeepSeek API
	request := deepSeekRequest{
		Model: s.modelName,
		Messages: []deepSeekMessage{
			{
				Role:    "system",
				Content: "You are a helpful code review assistant that provides concise, actionable feedback.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
	}

	// Send the request to the DeepSeek API
	var response deepSeekResponse
	resp, err := s.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post("https://api.deepseek.com/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("failed to call DeepSeek API: %w", err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("DeepSeek API returned error: %s", resp.String())
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from DeepSeek")
	}

	if response.Error != nil {
		return "", fmt.Errorf("DeepSeek error: %s", response.Error.Message)
	}

	// Return the AI's analysis
	return response.Choices[0].Message.Content, nil
}
