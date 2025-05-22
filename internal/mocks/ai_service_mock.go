package mocks

import (
	"context"
	"fmt"

	"github.com/ziadalikhalifa/code-review-papa/internal/ai"
)

var _ ai.Service = &MockAIService{}

// MockAIService is a mock implementation of the ai.Service interface.
type MockAIService struct {
	AnalyzeCodeFunc func(ctx context.Context, diff string) (string, error)
}

// AnalyzeCode calls the AnalyzeCodeFunc field.
func (m *MockAIService) AnalyzeCode(ctx context.Context, diff string) (string, error) {
	if m.AnalyzeCodeFunc == nil {
		return "", fmt.Errorf("MockAIService.AnalyzeCodeFunc is not set")
	}
	return m.AnalyzeCodeFunc(ctx, diff)
}
