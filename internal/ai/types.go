package ai

import "context"

// GenerateInput is provider-agnostic chat-style completion input.
type GenerateInput struct {
	SystemPrompt string
	UserPrompt   string
}

// GenerateOutput is normalized model text (raw JSON or prose).
type GenerateOutput struct {
	Text         string
	ProviderName string
}

// Provider abstracts any LLM or hosted inference (OpenAI-compatible HTTP, on-prem, etc.).
type Provider interface {
	// Name identifies the implementation for logging and API responses.
	Name() string
	// Generate returns model output for the given prompts.
	Generate(ctx context.Context, in GenerateInput) (GenerateOutput, error)
}
