package postgres

import (
	"time"

	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AutoMigrate ensures application tables exist (use versioned SQL in production when preferred).
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Subject{},
		&model.Room{},
		&model.Teacher{},
		&model.Group{},
		&model.StudentGroupMembership{},
		&model.UserTeacherLink{},
		&model.Schedule{},
		&model.Attendance{},
		&model.Grade{},
		&model.Payment{},
		&model.FileMetadata{},
		&model.Notification{},
	); err != nil {
		return err
	}
	return ensureDefaultSubject(db)
}

func ensureDefaultSubject(db *gorm.DB) error {
	var n int64
	if err := db.Model(&model.Subject{}).Where("code = ?", "GEN").Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	now := time.Now().UTC()
	s := &model.Subject{
		ID:        uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		Name:      "General",
		Code:      "GEN",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	return db.Create(s).Error
}
