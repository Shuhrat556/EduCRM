package model

import (
	"time"

	"github.com/google/uuid"
)

// StudentGroupMembership enforces at most one group per student user (PK = user_id).
type StudentGroupMembership struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
}

func (StudentGroupMembership) TableName() string {
	return "student_group_memberships"
}
