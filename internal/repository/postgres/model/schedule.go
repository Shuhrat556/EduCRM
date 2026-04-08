package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule is the GORM model for weekly recurring slots.
type Schedule struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID       uuid.UUID `gorm:"type:uuid;not null;index"`
	TeacherID     uuid.UUID `gorm:"type:uuid;not null;index"`
	RoomID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Weekday       int16     `gorm:"not null;index"`
	StartMinutes  int       `gorm:"not null"`
	EndMinutes    int       `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Schedule) TableName() string {
	return "schedules"
}

// BeforeCreate assigns ID when missing.
func (s *Schedule) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
