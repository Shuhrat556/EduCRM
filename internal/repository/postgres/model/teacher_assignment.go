package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeacherGroupSubjectAssignment links a teacher to a group+subject they may teach.
type TeacherGroupSubjectAssignment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeacherID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_tgsa_tg_sub"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:ux_tgsa_tg_sub"`
	SubjectID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:ux_tgsa_tg_sub"`
	CreatedAt time.Time
}

func (TeacherGroupSubjectAssignment) TableName() string {
	return "teacher_group_subject_assignments"
}

func (m *TeacherGroupSubjectAssignment) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
