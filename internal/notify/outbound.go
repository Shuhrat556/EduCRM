package notify

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/educrm/educrm-backend/internal/domain"
)

// Outbound fans out to external channels after an in-app notification is stored.
// Errors from individual providers are ignored here; real implementations may log/metrics.
type Outbound struct {
	Telegram TelegramProvider
	Push     PushProvider
}

// NewOutbound returns a dispatcher. Nil telegram or push are replaced with no-op implementations.
func NewOutbound(telegram TelegramProvider, push PushProvider) *Outbound {
	if telegram == nil {
		telegram = NoopTelegramProvider{}
	}
	if push == nil {
		push = NoopPushProvider{}
	}
	return &Outbound{Telegram: telegram, Push: push}
}

// DispatchOnCreated runs best-effort delivery for a newly persisted notification.
func (o *Outbound) DispatchOnCreated(ctx context.Context, n *domain.Notification) {
	if o == nil || n == nil {
		return
	}
	text := strings.TrimSpace(n.Title)
	if b := strings.TrimSpace(n.Body); b != "" {
		if text != "" {
			text += "\n"
		}
		text += b
	}
	if chatID, ok := telegramChatIDFromMetadata(n.Metadata); ok {
		_ = o.Telegram.SendText(ctx, chatID, text)
	}
	_ = o.Push.Send(ctx, PushPayload{
		UserID: n.UserID,
		Title:  n.Title,
		Body:   n.Body,
		Data:   stringMapFromJSONMetadata(n.Metadata),
	})
}

func telegramChatIDFromMetadata(raw json.RawMessage) (string, bool) {
	if len(raw) == 0 {
		return "", false
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return "", false
	}
	v, ok := m["telegram_chat_id"]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}

func stringMapFromJSONMetadata(raw json.RawMessage) map[string]string {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}
