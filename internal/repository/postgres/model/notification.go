package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification is the GORM model for in-app notifications.
type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Type      string     `gorm:"not null;size:64;index"`
	Title     string     `gorm:"not null;size:512"`
	Body      string     `gorm:"not null;type:text"`
	ReadAt    *time.Time `gorm:"index"`
	Metadata  []byte     `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate assigns ID when missing.
func (n *Notification) BeforeCreate(_ *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
