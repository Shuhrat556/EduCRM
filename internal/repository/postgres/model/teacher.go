package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Teacher is the GORM model for teachers.
type Teacher struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	FullName          string    `gorm:"not null;size:255"`
	Phone             *string   `gorm:"size:32;uniqueIndex"`
	Email             *string   `gorm:"size:255;uniqueIndex"`
	Specialization    *string   `gorm:"size:255"`
	PhotoURL          *string   `gorm:"size:2048"`
	PhotoStorageKey   *string   `gorm:"size:512"`
	PhotoContentType  *string   `gorm:"size:128"`
	PhotoOriginalName *string   `gorm:"size:255"`
	Status            string    `gorm:"not null;size:16;index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Teacher) TableName() string {
	return "teachers"
}

// BeforeCreate assigns ID when missing.
func (t *Teacher) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
