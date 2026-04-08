package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenRepository implements repository.RefreshTokenRepository.
type RefreshTokenRepository struct {
	db *gorm.DB
}

var _ repository.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

// NewRefreshTokenRepository constructs a RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a new refresh token hash.
func (r *RefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	row := model.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(&row).Error
}

// FindValidByHash returns the user ID for a non-expired token hash.
func (r *RefreshTokenRepository) FindValidByHash(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var row model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, repository.ErrNotFound
		}
		return uuid.Nil, err
	}
	return row.UserID, nil
}

// DeleteByHash removes a session by hash.
func (r *RefreshTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	res := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&model.RefreshToken{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// DeleteAllForUser revokes every refresh session for a user.
func (r *RefreshTokenRepository) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}

// Replace rotates a refresh token in one transaction.
func (r *RefreshTokenRepository) Replace(ctx context.Context, userID uuid.UUID, oldHash, newHash string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		var old model.RefreshToken
		if err := tx.Where("token_hash = ? AND user_id = ? AND expires_at > ?", oldHash, userID, now).
			First(&old).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return repository.ErrNotFound
			}
			return err
		}
		if err := tx.Delete(&model.RefreshToken{}, "id = ?", old.ID).Error; err != nil {
			return err
		}
		row := model.RefreshToken{
			ID:        uuid.New(),
			UserID:    userID,
			TokenHash: newHash,
			ExpiresAt: expiresAt,
		}
		return tx.Create(&row).Error
	})
}
