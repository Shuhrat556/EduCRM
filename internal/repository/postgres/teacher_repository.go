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

// TeacherRepository implements repository.TeacherRepository.
type TeacherRepository struct {
	db *gorm.DB
}

var _ repository.TeacherRepository = (*TeacherRepository)(nil)

// NewTeacherRepository constructs a TeacherRepository.
func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// Create inserts a teacher.
func (r *TeacherRepository) Create(ctx context.Context, t *domain.Teacher) error {
	m, err := teacherToModel(t)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicate
		}
		return err
	}
	return nil
}

// Update updates scalar teacher fields.
func (r *TeacherRepository) Update(ctx context.Context, t *domain.Teacher) error {
	m, err := teacherToModel(t)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(&model.Teacher{}).Where("id = ?", t.ID).Updates(map[string]any{
		"full_name":           m.FullName,
		"phone":               m.Phone,
		"email":               m.Email,
		"specialization":      m.Specialization,
		"photo_url":           m.PhotoURL,
		"photo_storage_key":   m.PhotoStorageKey,
		"photo_content_type":  m.PhotoContentType,
		"photo_original_name": m.PhotoOriginalName,
		"status":              m.Status,
		"updated_at":          m.UpdatedAt,
	}).Error
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicate
		}
		return err
	}
	return nil
}

// Delete removes a teacher. Fails with ErrReferenced if groups still reference this teacher.
func (r *TeacherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Teacher{}, "id = ?", id)
	if res.Error != nil {
		if isForeignKeyViolation(res.Error) {
			return repository.ErrReferenced
		}
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads a teacher and groups where groups.teacher_id = id.
func (r *TeacherRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Teacher, []domain.GroupBrief, error) {
	var m model.Teacher
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	t, err := teacherToDomain(&m)
	if err != nil {
		return nil, nil, err
	}
	var rows []model.Group
	if err := r.db.WithContext(ctx).Model(&model.Group{}).
		Where("teacher_id = ?", id).
		Order("name ASC").
		Find(&rows).Error; err != nil {
		return nil, nil, err
	}
	groups := make([]domain.GroupBrief, 0, len(rows))
	for i := range rows {
		groups = append(groups, domain.GroupBrief{ID: rows[i].ID, Name: strings.TrimSpace(rows[i].Name)})
	}
	return t, groups, nil
}

// List returns paginated teachers with group counts.
func (r *TeacherRepository) List(ctx context.Context, p repository.TeacherListParams) ([]repository.TeacherListEntry, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.Teacher{})
		if p.Search != "" {
			term := "%" + escapeLikePattern(p.Search) + "%"
			q = q.Where(
				`(teachers.full_name ILIKE ? ESCAPE '\' OR LOWER(COALESCE(teachers.email,'')) LIKE LOWER(?) ESCAPE '\' OR COALESCE(teachers.phone,'') LIKE ? ESCAPE '\' OR COALESCE(teachers.specialization,'') ILIKE ? ESCAPE '\')`,
				term, term, term, term,
			)
		}
		if p.Status != nil {
			q = q.Where("teachers.status = ?", string(*p.Status))
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

	var teachers []model.Teacher
	if err := build().Order("teachers.created_at DESC").Limit(size).Offset(offset).Find(&teachers).Error; err != nil {
		return nil, 0, err
	}
	if len(teachers) == 0 {
		return nil, total, nil
	}
	ids := make([]uuid.UUID, len(teachers))
	for i := range teachers {
		ids[i] = teachers[i].ID
	}
	type cnt struct {
		TeacherID uuid.UUID `gorm:"column:teacher_id"`
		N         int64     `gorm:"column:n"`
	}
	var counts []cnt
	if err := r.db.WithContext(ctx).Model(&model.Group{}).
		Select("teacher_id, COUNT(*) as n").
		Where("teacher_id IN ?", ids).
		Group("teacher_id").
		Scan(&counts).Error; err != nil {
		return nil, 0, err
	}
	countBy := make(map[uuid.UUID]int64, len(counts))
	for _, c := range counts {
		countBy[c.TeacherID] = c.N
	}
	out := make([]repository.TeacherListEntry, 0, len(teachers))
	for i := range teachers {
		t, err := teacherToDomain(&teachers[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, repository.TeacherListEntry{
			Teacher:    *t,
			GroupCount: int(countBy[teachers[i].ID]),
		})
	}
	return out, total, nil
}

// EmailTaken checks email uniqueness (case-insensitive).
func (r *TeacherRepository) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return false, nil
	}
	q := r.db.WithContext(ctx).Model(&model.Teacher{}).Where("LOWER(email) = ?", email)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

// PhoneTaken checks phone uniqueness.
func (r *TeacherRepository) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false, nil
	}
	q := r.db.WithContext(ctx).Model(&model.Teacher{}).Where("phone = ?", phone)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func teacherToDomain(m *model.Teacher) (*domain.Teacher, error) {
	st, err := domain.ParseTeacherStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Teacher{
		ID:                m.ID,
		FullName:          m.FullName,
		Phone:             m.Phone,
		Email:             m.Email,
		Specialization:    m.Specialization,
		PhotoURL:          m.PhotoURL,
		PhotoStorageKey:   m.PhotoStorageKey,
		PhotoContentType:  m.PhotoContentType,
		PhotoOriginalName: m.PhotoOriginalName,
		Status:            st,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}, nil
}

func teacherToModel(t *domain.Teacher) (*model.Teacher, error) {
	return &model.Teacher{
		ID:                t.ID,
		FullName:          t.FullName,
		Phone:             t.Phone,
		Email:             t.Email,
		Specialization:    t.Specialization,
		PhotoURL:          t.PhotoURL,
		PhotoStorageKey:   t.PhotoStorageKey,
		PhotoContentType:  t.PhotoContentType,
		PhotoOriginalName: t.PhotoOriginalName,
		Status:            string(t.Status),
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}, nil
}
