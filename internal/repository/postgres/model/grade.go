package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Grade is the GORM model for weekly ratings.
type Grade struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	StudentID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_grade_weekly"`
	TeacherID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_grade_weekly"`
	GroupID        uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_grade_weekly"`
	WeekStartDate  time.Time `gorm:"type:date;not null;uniqueIndex:ux_grade_weekly"`
	GradeType      string    `gorm:"not null;size:32;uniqueIndex:ux_grade_weekly"`
	GradeValue     float64   `gorm:"not null"`
	Comment        *string   `gorm:"type:text"`
	GradedAt       time.Time `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Grade) TableName() string {
	return "grades"
}

// BeforeCreate assigns ID when missing.
func (g *Grade) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
