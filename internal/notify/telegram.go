package notify

import "context"

// TelegramProvider sends outbound messages via a Telegram Bot API integration.
// Wire a real implementation (HTTP client to api.telegram.org) when the bot is ready.
type TelegramProvider interface {
	SendText(ctx context.Context, chatID, text string) error
}

// NoopTelegramProvider satisfies TelegramProvider without sending (default until configured).
type NoopTelegramProvider struct{}

// SendText implements TelegramProvider as a no-op.
func (NoopTelegramProvider) SendText(ctx context.Context, chatID, text string) error {
	_ = ctx
	_ = chatID
	_ = text
	return nil
}
