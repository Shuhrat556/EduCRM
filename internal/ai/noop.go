package ai

import (
	"context"
	"fmt"
	"strings"
)

// NoopProvider returns deterministic sample text for local dev and tests.
type NoopProvider struct{}

// Name implements Provider.
func (NoopProvider) Name() string { return "noop" }

// Generate implements Provider.
func (NoopProvider) Generate(ctx context.Context, in GenerateInput) (GenerateOutput, error) {
	_ = ctx
	preview := strings.TrimSpace(in.UserPrompt)
	if len(preview) > 400 {
		preview = preview[:400] + "…"
	}
	text := fmt.Sprintf("[noop AI sample] System instructions were applied. User context preview:\n%s\n\nReplace AI_PROVIDER=http and set AI_HTTP_* to call a real model.", preview)
	return GenerateOutput{Text: text, ProviderName: "noop"}, nil
}
