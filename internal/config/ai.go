package config

import "time"

// AIConfig selects an inference backend and prompt template directory.
// Swap AI_PROVIDER and HTTP settings without changing application code.
type AIConfig struct {
	Provider    string        `env:"PROVIDER" envDefault:"noop"` // noop | http (OpenAI-compatible chat completions)
	PromptsDir  string        `env:"PROMPTS_DIR" envDefault:"configs/ai/prompts"`
	HTTPBaseURL string        `env:"HTTP_BASE_URL"` // e.g. https://api.openai.com/v1 or a LiteLLM proxy
	HTTPAPIKey  string        `env:"HTTP_API_KEY"`
	HTTPModel   string        `env:"HTTP_MODEL" envDefault:"gpt-4o-mini"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"60s"`
}
