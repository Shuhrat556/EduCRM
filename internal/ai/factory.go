package ai

import (
	"fmt"
	"strings"

	"github.com/educrm/educrm-backend/internal/config"
)

// NewProviderFromConfig returns the configured Provider implementation.
func NewProviderFromConfig(cfg config.AIConfig) (Provider, error) {
	p := strings.ToLower(strings.TrimSpace(cfg.Provider))
	switch p {
	case "noop", "":
		return NoopProvider{}, nil
	case "http":
		return NewHTTPChatProvider(cfg.HTTPBaseURL, cfg.HTTPAPIKey, cfg.HTTPModel, cfg.HTTPTimeout), nil
	default:
		return nil, fmt.Errorf("unknown AI_PROVIDER %q (use noop or http)", cfg.Provider)
	}
}
