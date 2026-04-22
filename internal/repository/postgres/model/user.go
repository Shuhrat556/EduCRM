package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the GORM model for users.
type User struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FullName            string     `gorm:"not null;size:255;default:''"`
	Username            *string    `gorm:"size:64;uniqueIndex"`
	Email               *string    `gorm:"size:255;uniqueIndex"`
	Phone               *string    `gorm:"size:32;uniqueIndex"`
	PasswordHash        string     `gorm:"not null"`
	Role                string     `gorm:"not null;size:32"`
	IsActive            bool       `gorm:"not null;default:true;index"`
	ForcePasswordChange bool       `gorm:"not null;default:false"`
	CreatedByUserID     *uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate assigns a primary key when missing (e.g. future registration flows).
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
