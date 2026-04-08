package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPChatProvider calls an OpenAI-compatible POST /chat/completions endpoint.
// BaseURL should be the API root including /v1 when required (e.g. https://api.openai.com/v1).
type HTTPChatProvider struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewHTTPChatProvider builds a chat-completions client. Trailing slashes on baseURL are trimmed.
func NewHTTPChatProvider(baseURL, apiKey, model string, timeout time.Duration) *HTTPChatProvider {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &HTTPChatProvider{
		baseURL: strings.TrimSuffix(strings.TrimSpace(baseURL), "/"),
		apiKey:  strings.TrimSpace(apiKey),
		model:   strings.TrimSpace(model),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Name implements Provider.
func (p *HTTPChatProvider) Name() string { return "http_chat" }

// Generate implements Provider.
func (p *HTTPChatProvider) Generate(ctx context.Context, in GenerateInput) (GenerateOutput, error) {
	if p.baseURL == "" || p.apiKey == "" {
		return GenerateOutput{}, fmt.Errorf("AI_HTTP_BASE_URL and AI_HTTP_API_KEY must be set for http provider")
	}
	if p.model == "" {
		return GenerateOutput{}, fmt.Errorf("AI_HTTP_MODEL is empty")
	}
	url := p.baseURL + "/chat/completions"
	body := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": in.SystemPrompt},
			{"role": "user", "content": in.UserPrompt},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return GenerateOutput{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return GenerateOutput{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return GenerateOutput{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return GenerateOutput{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return GenerateOutput{}, fmt.Errorf("ai http %d: %s", resp.StatusCode, string(b))
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return GenerateOutput{}, fmt.Errorf("decode ai response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return GenerateOutput{}, fmt.Errorf("empty choices from ai response")
	}
	return GenerateOutput{
		Text:         parsed.Choices[0].Message.Content,
		ProviderName: p.Name(),
	}, nil
}
