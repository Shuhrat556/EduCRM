package user

import (
	"context"
	"errors"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	authpkg "github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service contains user management use cases.
type Service struct {
	repo repository.UserRepository
}

// NewService constructs a user service.
func NewService(repo repository.UserRepository) *Service {
	return &Service{repo: repo}
}

// UserPublic is a safe projection for APIs.
type UserPublic struct {
	ID        uuid.UUID   `json:"id"`
	Email     *string     `json:"email,omitempty"`
	Phone     *string     `json:"phone,omitempty"`
	Role      domain.Role `json:"role"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ListResult is a paginated user list.
type ListResult struct {
	Items    []UserPublic `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// CreateInput holds validated fields for user creation.
type CreateInput struct {
	Email    *string
	Phone    *string
	Password string
	Role     domain.Role
}

// UpdateInput holds optional field updates.
type UpdateInput struct {
	Email    *string
	Phone    *string
	Password *string
	Role     *string
	IsActive *bool
}

// Create persists a new user.
func (s *Service) Create(ctx context.Context, actor domain.Role, in CreateInput) (*UserPublic, error) {
	if err := AssertCanCreateTargetRole(actor, in.Role); err != nil {
		return nil, err
	}
	email := domain.NormalizeEmail(in.Email)
	phone := domain.NormalizePhone(in.Phone)
	if email == nil && phone == nil {
		return nil, apperror.Validation("contact", "at least one of email or phone is required")
	}
	if err := s.assertUniqueContact(ctx, email, phone, nil); err != nil {
		return nil, err
	}
	hash, err := authpkg.HashPassword(in.Password)
	if err != nil {
		return nil, apperror.Internal("hash password").Wrap(err)
	}
	now := time.Now().UTC()
	u := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		Phone:        phone,
		PasswordHash: hash,
		Role:         in.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, wrapUserRepoErr(err)
	}
	return publicPtr(u), nil
}

// GetByID returns a user if the actor may access them.
func (s *Service) GetByID(ctx context.Context, actor domain.Role, id uuid.UUID) (*UserPublic, error) {
	if err := AssertActorCanManageUsers(actor); err != nil {
		return nil, err
	}
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.NotFound("user")
	}
	if err := AssertActorCanAccessTargetUser(actor, u.Role); err != nil {
		return nil, apperror.NotFound("user")
	}
	return publicPtr(u), nil
}

// List returns a filtered page of users.
func (s *Service) List(ctx context.Context, actor domain.Role, params repository.UserListParams) (*ListResult, error) {
	if err := AssertActorCanManageUsers(actor); err != nil {
		return nil, err
	}
	if actor == domain.RoleAdmin {
		params.ExcludeRoles = []domain.Role{domain.RoleSuperAdmin, domain.RoleAdmin}
	}
	users, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, apperror.Internal("list users").Wrap(err)
	}
	items := make([]UserPublic, 0, len(users))
	for i := range users {
		items = append(items, *publicPtr(&users[i]))
	}
	page := params.Page
	if page < 1 {
		page = 1
	}
	size := params.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return &ListResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

// Update applies changes to a user.
func (s *Service) Update(ctx context.Context, actor domain.Role, actorUserID, id uuid.UUID, in UpdateInput) (*UserPublic, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.NotFound("user")
	}
	if err := AssertActorCanAccessTargetUser(actor, u.Role); err != nil {
		return nil, apperror.NotFound("user")
	}
	newRole := u.Role
	if in.Role != nil {
		parsed, perr := domain.ParseRole(*in.Role)
		if perr != nil {
			return nil, apperror.Validation("role", "invalid role")
		}
		if err := AssertCanAssignRoleOnUpdate(actor, parsed); err != nil {
			return nil, err
		}
		newRole = parsed
	}
	if err := AssertActorCanAccessTargetUser(actor, newRole); err != nil {
		return nil, apperror.Forbidden("insufficient permissions for the requested role")
	}
	if in.Email != nil {
		u.Email = domain.NormalizeEmail(in.Email)
	}
	if in.Phone != nil {
		u.Phone = domain.NormalizePhone(in.Phone)
	}
	if u.Email == nil && u.Phone == nil {
		return nil, apperror.Validation("contact", "at least one of email or phone is required")
	}
	if in.Password != nil && *in.Password != "" {
		hash, herr := authpkg.HashPassword(*in.Password)
		if herr != nil {
			return nil, apperror.Internal("hash password").Wrap(herr)
		}
		u.PasswordHash = hash
	}
	u.Role = newRole
	if in.IsActive != nil {
		if !*in.IsActive && id == actorUserID {
			return nil, apperror.Validation("is_active", "cannot deactivate your own account")
		}
		u.IsActive = *in.IsActive
	}
	u.UpdatedAt = time.Now().UTC()
	if err := s.assertUniqueContact(ctx, u.Email, u.Phone, &u.ID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, wrapUserRepoErr(err)
	}
	return publicPtr(u), nil
}

// SetActive updates only the activation flag.
func (s *Service) SetActive(ctx context.Context, actor domain.Role, actorUserID, id uuid.UUID, active bool) (*UserPublic, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.NotFound("user")
	}
	if err := AssertActorCanAccessTargetUser(actor, u.Role); err != nil {
		return nil, apperror.NotFound("user")
	}
	if !active && actorUserID == id {
		return nil, apperror.Validation("is_active", "cannot deactivate your own account")
	}
	u.IsActive = active
	u.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, wrapUserRepoErr(err)
	}
	return publicPtr(u), nil
}

// Delete removes a user.
func (s *Service) Delete(ctx context.Context, actor domain.Role, actorID uuid.UUID, id uuid.UUID) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return apperror.NotFound("user")
	}
	if err := AssertActorCanAccessTargetUser(actor, u.Role); err != nil {
		return apperror.NotFound("user")
	}
	if id == actorID {
		return apperror.Validation("user", "cannot delete your own account")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user")
		}
		return apperror.Internal("delete user").Wrap(err)
	}
	return nil
}

func (s *Service) assertUniqueContact(ctx context.Context, email, phone *string, excludeID *uuid.UUID) error {
	if email != nil {
		taken, err := s.repo.EmailTaken(ctx, *email, excludeID)
		if err != nil {
			return apperror.Internal("check email").Wrap(err)
		}
		if taken {
			return apperror.Conflict("email_taken", "email is already in use")
		}
	}
	if phone != nil {
		taken, err := s.repo.PhoneTaken(ctx, *phone, excludeID)
		if err != nil {
			return apperror.Internal("check phone").Wrap(err)
		}
		if taken {
			return apperror.Conflict("phone_taken", "phone is already in use")
		}
	}
	return nil
}

func wrapUserRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrDuplicate) {
		return apperror.Conflict("unique_violation", "email or phone already in use")
	}
	return apperror.Internal("persist user").Wrap(err)
}

func publicPtr(u *domain.User) *UserPublic {
	return &UserPublic{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
