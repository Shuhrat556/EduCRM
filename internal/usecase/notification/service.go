package notification

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/notify"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles in-app notifications and triggers outbound dispatch.
type Service struct {
	repo     repository.NotificationRepository
	users    repository.UserRepository
	outbound *notify.Outbound
}

// NewService constructs the notification service.
func NewService(
	repo repository.NotificationRepository,
	users repository.UserRepository,
	outbound *notify.Outbound,
) *Service {
	return &Service{repo: repo, users: users, outbound: outbound}
}

// CreateInput is staff-created notification payload.
type CreateInput struct {
	UserID   uuid.UUID
	Type     domain.NotificationType
	Title    string
	Body     string
	Metadata json.RawMessage
}

// Create persists a notification and runs outbound dispatch.
func (s *Service) Create(ctx context.Context, actorRole domain.Role, in CreateInput) (*domain.Notification, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	u, err := s.users.FindByID(ctx, in.UserID)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.Validation("user_id", "user not found")
	}
	if len(in.Metadata) > 0 && !json.Valid(in.Metadata) {
		return nil, apperror.Validation("metadata", "must be valid JSON")
	}
	if strings.TrimSpace(in.Title) == "" {
		return nil, apperror.Validation("title", "required")
	}
	now := time.Now().UTC()
	var meta json.RawMessage
	if len(in.Metadata) > 0 {
		meta = append(json.RawMessage(nil), in.Metadata...)
	}
	row := &domain.Notification{
		ID:        uuid.New(),
		UserID:    in.UserID,
		Type:      in.Type,
		Title:     strings.TrimSpace(in.Title),
		Body:      strings.TrimSpace(in.Body),
		Metadata:  meta,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, row); err != nil {
		return nil, apperror.Internal("create notification").Wrap(err)
	}
	if s.outbound != nil {
		s.outbound.DispatchOnCreated(ctx, row)
	}
	return row, nil
}

// Get returns one notification if the actor may access it.
func (s *Service) Get(ctx context.Context, actorRole domain.Role, actorUserID, id uuid.UUID) (*domain.Notification, error) {
	row, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load notification").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("notification")
	}
	if err := s.assertAccess(actorRole, actorUserID, row.UserID); err != nil {
		return nil, err
	}
	return row, nil
}

// List returns paginated notifications for the target user (self or staff-selected).
func (s *Service) List(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, targetUserID *uuid.UUID, unreadOnly bool, page, pageSize int) ([]domain.Notification, int64, error) {
	uid, err := s.resolveListUser(actorRole, actorUserID, targetUserID)
	if err != nil {
		return nil, 0, err
	}
	items, total, err := s.repo.List(ctx, repository.NotificationListParams{
		UserID:     uid,
		UnreadOnly: unreadOnly,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return nil, 0, apperror.Internal("list notifications").Wrap(err)
	}
	return items, total, nil
}

// MarkRead sets read_at when the actor may access the notification.
func (s *Service) MarkRead(ctx context.Context, actorRole domain.Role, actorUserID, id uuid.UUID) (*domain.Notification, error) {
	row, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load notification").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("notification")
	}
	if err := s.assertAccess(actorRole, actorUserID, row.UserID); err != nil {
		return nil, err
	}
	if row.ReadAt != nil {
		return row, nil
	}
	now := time.Now().UTC()
	row.ReadAt = &now
	row.UpdatedAt = now
	if err := s.repo.Update(ctx, row); err != nil {
		return nil, apperror.Internal("update notification").Wrap(err)
	}
	return row, nil
}

// Delete removes a notification when the actor may access it.
func (s *Service) Delete(ctx context.Context, actorRole domain.Role, actorUserID, id uuid.UUID) error {
	row, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load notification").Wrap(err)
	}
	if row == nil {
		return apperror.NotFound("notification")
	}
	if err := s.assertAccess(actorRole, actorUserID, row.UserID); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("notification")
		}
		return apperror.Internal("delete notification").Wrap(err)
	}
	return nil
}

func (s *Service) resolveListUser(actorRole domain.Role, actorUserID uuid.UUID, target *uuid.UUID) (uuid.UUID, error) {
	if rbac.IsStaff(actorRole) && target != nil {
		return *target, nil
	}
	return actorUserID, nil
}

func (s *Service) assertAccess(actorRole domain.Role, actorUserID, recipientID uuid.UUID) error {
	_ = s
	if rbac.IsStaff(actorRole) {
		return nil
	}
	if actorUserID != recipientID {
		return apperror.New(apperror.KindForbidden, "forbidden", "you may only access your own notifications")
	}
	return nil
}
