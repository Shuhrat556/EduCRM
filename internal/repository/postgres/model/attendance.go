package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attendance is the GORM model for per-lesson attendance.
type Attendance struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey"`
	StudentID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:ux_attendance_slot"`
	GroupID             uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:ux_attendance_slot"`
	LessonDate          time.Time  `gorm:"type:date;not null;uniqueIndex:ux_attendance_slot"`
	Status              string     `gorm:"not null;size:16"`
	Comment             *string    `gorm:"type:text"`
	MarkedByTeacherID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (Attendance) TableName() string {
	return "attendances"
}

// BeforeCreate assigns ID when missing.
func (a *Attendance) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
