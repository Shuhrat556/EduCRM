package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupRepository implements repository.GroupRepository.
type GroupRepository struct {
	db *gorm.DB
}

var _ repository.GroupRepository = (*GroupRepository)(nil)

// NewGroupRepository constructs GroupRepository.
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// Create inserts a group.
func (r *GroupRepository) Create(ctx context.Context, g *domain.Group) error {
	m, err := groupToModel(g)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// Update updates scalar fields.
func (r *GroupRepository) Update(ctx context.Context, g *domain.Group) error {
	m, err := groupToModel(g)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Group{}).Where("id = ?", g.ID).Updates(map[string]any{
		"name":               m.Name,
		"subject_id":         m.SubjectID,
		"teacher_id":         m.TeacherID,
		"room_id":            m.RoomID,
		"start_date":         m.StartDate,
		"end_date":           m.EndDate,
		"monthly_fee_minor":  m.MonthlyFeeMinor,
		"status":             m.Status,
		"updated_at":         m.UpdatedAt,
	}).Error
}

// Delete removes a group by id.
func (r *GroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Group{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads a group by primary key.
func (r *GroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	var m model.Group
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return groupToDomain(&m)
}

// List returns paginated groups.
func (r *GroupRepository) List(ctx context.Context, p repository.GroupListParams) ([]domain.Group, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.Group{})
		if p.Search != "" {
			term := "%" + escapeLikePattern(p.Search) + "%"
			q = q.Where("groups.name ILIKE ? ESCAPE '\\'", term)
		}
		if p.Status != nil {
			q = q.Where("groups.status = ?", string(*p.Status))
		}
		if p.TeacherID != nil {
			q = q.Where("groups.teacher_id = ?", *p.TeacherID)
		}
		if p.SubjectID != nil {
			q = q.Where("groups.subject_id = ?", *p.SubjectID)
		}
		if p.RoomID != nil {
			q = q.Where("groups.room_id = ?", *p.RoomID)
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
	var rows []model.Group
	if err := build().Order("groups.start_date DESC, groups.name ASC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.Group, 0, len(rows))
	for i := range rows {
		g, err := groupToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *g)
	}
	return out, total, nil
}

func groupToDomain(m *model.Group) (*domain.Group, error) {
	st, err := domain.ParseGroupStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Group{
		ID:              m.ID,
		Name:            strings.TrimSpace(m.Name),
		SubjectID:       m.SubjectID,
		TeacherID:       m.TeacherID,
		RoomID:          m.RoomID,
		StartDate:       m.StartDate.UTC(),
		EndDate:         m.EndDate.UTC(),
		MonthlyFeeMinor: m.MonthlyFeeMinor,
		Status:          st,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}, nil
}

func groupToModel(g *domain.Group) (*model.Group, error) {
	return &model.Group{
		ID:              g.ID,
		Name:            g.Name,
		SubjectID:       g.SubjectID,
		TeacherID:       g.TeacherID,
		RoomID:          g.RoomID,
		StartDate:       g.StartDate,
		EndDate:         g.EndDate,
		MonthlyFeeMinor: g.MonthlyFeeMinor,
		Status:          string(g.Status),
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}, nil
}
