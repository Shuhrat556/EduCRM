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

// RoomRepository implements repository.RoomRepository.
type RoomRepository struct {
	db *gorm.DB
}

var _ repository.RoomRepository = (*RoomRepository)(nil)

// NewRoomRepository constructs a RoomRepository.
func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create inserts a room.
func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	m, err := roomToModel(room)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// Update persists scalar fields.
func (r *RoomRepository) Update(ctx context.Context, room *domain.Room) error {
	m, err := roomToModel(room)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Room{}).Where("id = ?", room.ID).Updates(map[string]any{
		"name":        m.Name,
		"capacity":    m.Capacity,
		"description": m.Description,
		"status":      m.Status,
		"updated_at":  m.UpdatedAt,
	}).Error
}

// Delete removes a room by id.
func (r *RoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Room{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads a room by primary key.
func (r *RoomRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	var m model.Room
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return roomToDomain(&m)
}

// List returns a paginated slice and total count.
func (r *RoomRepository) List(ctx context.Context, p repository.RoomListParams) ([]domain.Room, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.Room{})
		if p.Search != "" {
			term := "%" + escapeLikePattern(p.Search) + "%"
			q = q.Where(
				"(rooms.name ILIKE ? ESCAPE '\\' OR COALESCE(rooms.description,'') ILIKE ? ESCAPE '\\')",
				term, term,
			)
		}
		if p.Status != nil {
			q = q.Where("rooms.status = ?", string(*p.Status))
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
	var rows []model.Room
	if err := build().Order("rooms.name ASC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.Room, 0, len(rows))
	for i := range rows {
		d, err := roomToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *d)
	}
	return out, total, nil
}

func roomToDomain(m *model.Room) (*domain.Room, error) {
	st, err := domain.ParseRoomStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Room{
		ID:          m.ID,
		Name:        m.Name,
		Capacity:    m.Capacity,
		Description: m.Description,
		Status:      st,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

func roomToModel(r *domain.Room) (*model.Room, error) {
	return &model.Room{
		ID:          r.ID,
		Name:        r.Name,
		Capacity:    r.Capacity,
		Description: r.Description,
		Status:      string(r.Status),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}
