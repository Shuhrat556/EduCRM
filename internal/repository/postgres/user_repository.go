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

// UserRepository implements repository.UserRepository with GORM.
type UserRepository struct {
	db *gorm.DB
}

var _ repository.UserRepository = (*UserRepository)(nil)

// NewUserRepository constructs a UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByLogin resolves login as email (case-insensitive) if it contains "@", otherwise phone (trimmed).
// Only active users can authenticate.
func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return nil, nil
	}
	var m model.User
	q := r.db.WithContext(ctx).Where("is_active = ?", true)
	if strings.Contains(login, "@") {
		q = q.Where("LOWER(email) = ?", strings.ToLower(login))
	} else {
		q = q.Where("phone = ?", login)
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userToDomain(&m)
}

// FindByID returns a user by primary key (any activation state).
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userToDomain(&m)
}

// Create inserts a new user.
func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	m, err := domainToModel(u)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicate
		}
		return err
	}
	return nil
}

// Update persists all scalar fields on an existing user.
func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	m, err := domainToModel(u)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", u.ID).Updates(map[string]any{
		"email":         m.Email,
		"phone":         m.Phone,
		"password_hash": m.PasswordHash,
		"role":          m.Role,
		"is_active":     m.IsActive,
		"updated_at":    m.UpdatedAt,
	}).Error
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicate
		}
		return err
	}
	return nil
}

// Delete removes a user (refresh tokens cascade when FK is set).
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// List returns a page of users and the total count for the same filter.
func (r *UserRepository) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.User{})
		if p.Search != "" {
			term := "%" + escapeLikePattern(p.Search) + "%"
			q = q.Where(
				"(LOWER(COALESCE(email, '')) LIKE LOWER(?) ESCAPE '\\' OR COALESCE(phone, '') LIKE ? ESCAPE '\\')",
				term, term,
			)
		}
		if p.Role != nil {
			q = q.Where("role = ?", string(*p.Role))
		}
		if len(p.ExcludeRoles) > 0 {
			ex := make([]string, 0, len(p.ExcludeRoles))
			for _, role := range p.ExcludeRoles {
				ex = append(ex, string(role))
			}
			q = q.Where("role NOT IN ?", ex)
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
	var rows []model.User
	if err := build().Order("created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.User, 0, len(rows))
	for i := range rows {
		u, err := userToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *u)
	}
	return out, total, nil
}

// EmailTaken reports whether another user already owns this email (case-insensitive).
func (r *UserRepository) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return false, nil
	}
	q := r.db.WithContext(ctx).Model(&model.User{}).Where("LOWER(email) = ?", email)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

// PhoneTaken reports whether another user already owns this phone.
func (r *UserRepository) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false, nil
	}
	q := r.db.WithContext(ctx).Model(&model.User{}).Where("phone = ?", phone)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func userToDomain(m *model.User) (*domain.User, error) {
	role, err := domain.ParseRole(m.Role)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		Phone:        m.Phone,
		PasswordHash: m.PasswordHash,
		Role:         role,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

func domainToModel(u *domain.User) (*model.User, error) {
	return &model.User{
		ID:           u.ID,
		Email:        u.Email,
		Phone:        u.Phone,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}, nil
}
