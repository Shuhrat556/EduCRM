package postgres

import (
	"context"
	"errors"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationRepository implements repository.NotificationRepository.
type NotificationRepository struct {
	db *gorm.DB
}

var _ repository.NotificationRepository = (*NotificationRepository)(nil)

// NewNotificationRepository constructs NotificationRepository.
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create inserts a notification.
func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	m, err := notificationToModel(n)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// Update persists read state (read_at, updated_at).
func (r *NotificationRepository) Update(ctx context.Context, n *domain.Notification) error {
	res := r.db.WithContext(ctx).Model(&model.Notification{}).Where("id = ?", n.ID).Updates(map[string]any{
		"read_at":    n.ReadAt,
		"updated_at": n.UpdatedAt,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads by id.
func (r *NotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	var m model.Notification
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return notificationToDomain(&m)
}

// List paginates notifications for a user.
func (r *NotificationRepository) List(ctx context.Context, p repository.NotificationListParams) ([]domain.Notification, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", p.UserID)
		if p.UnreadOnly {
			q = q.Where("read_at IS NULL")
		}
		return q
	}
	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := p.Page
	if page < 1 {
		page = 1
	}
	size := p.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	var rows []model.Notification
	if err := build().Order("created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.Notification, 0, len(rows))
	for i := range rows {
		d, err := notificationToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *d)
	}
	return out, total, nil
}

// Delete removes a notification row.
func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Notification{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func notificationToDomain(m *model.Notification) (*domain.Notification, error) {
	nt, err := domain.ParseNotificationType(m.Type)
	if err != nil {
		return nil, err
	}
	var meta []byte
	if len(m.Metadata) > 0 {
		meta = append([]byte(nil), m.Metadata...)
	}
	return &domain.Notification{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      nt,
		Title:     m.Title,
		Body:      m.Body,
		ReadAt:    m.ReadAt,
		Metadata:  meta,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func notificationToModel(n *domain.Notification) (*model.Notification, error) {
	var meta []byte
	if len(n.Metadata) > 0 {
		meta = append([]byte(nil), n.Metadata...)
	}
	return &model.Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      string(n.Type),
		Title:     n.Title,
		Body:      n.Body,
		ReadAt:    n.ReadAt,
		Metadata:  meta,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}, nil
}
