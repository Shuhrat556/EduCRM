package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group is the GORM model for class groups.
type Group struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name            string     `gorm:"not null;size:255"`
	SubjectID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	TeacherID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	RoomID          *uuid.UUID `gorm:"type:uuid;index"`
	StartDate       time.Time  `gorm:"type:date;not null"`
	EndDate         time.Time  `gorm:"type:date;not null"`
	MonthlyFeeMinor int64      `gorm:"not null"`
	Status          string     `gorm:"not null;size:16;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Group) TableName() string {
	return "groups"
}

// BeforeCreate assigns ID when missing.
func (g *Group) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
