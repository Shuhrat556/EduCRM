package model

import (
	"github.com/google/uuid"
)

// UserTeacherLink binds a users row (teacher login) to a teachers profile.
type UserTeacherLink struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeacherID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
}

func (UserTeacherLink) TableName() string {
	return "user_teacher_links"
}
