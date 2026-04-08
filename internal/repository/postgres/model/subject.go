package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subject is the GORM model for subjects (courses).
type Subject struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null;size:255"`
	Description *string   `gorm:"type:text"`
	Code        string    `gorm:"not null;size:64;uniqueIndex"`
	Status      string    `gorm:"not null;size:16;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Subject) TableName() string {
	return "subjects"
}

// BeforeCreate assigns ID when missing.
func (s *Subject) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
