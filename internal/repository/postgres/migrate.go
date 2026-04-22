package postgres

import (
	"time"

	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AutoMigrate ensures application tables exist (use versioned SQL in production when preferred).
func AutoMigrate(db *gorm.DB) error {
	// Keep AutoMigrate safe on existing databases.
	// Postgres cannot add a NOT NULL column to a non-empty table unless we provide
	// a default/backfill path.
	if err := ensureUsersFullName(db); err != nil {
		return err
	}

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
		&model.TeacherGroupSubjectAssignment{},
	); err != nil {
		return err
	}
	return ensureDefaultSubject(db)
}

func ensureUsersFullName(db *gorm.DB) error {
	// If the `users` table doesn't exist yet, AutoMigrate will create it and apply defaults.
	if !db.Migrator().HasTable(&model.User{}) {
		return nil
	}

	// Ensure column exists and is compatible with NOT NULL requirement.
	// Use raw SQL to be explicit and idempotent.
	if err := db.Exec(`
		ALTER TABLE users
			ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		UPDATE users
		SET full_name = ''
		WHERE full_name IS NULL;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		ALTER TABLE users
			ALTER COLUMN full_name SET DEFAULT '';
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		ALTER TABLE users
			ALTER COLUMN full_name SET NOT NULL;
	`).Error; err != nil {
		return err
	}

	return nil
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
