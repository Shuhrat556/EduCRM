package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room is the GORM model for rooms.
type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null;size:255;index"`
	Capacity    int       `gorm:"not null"`
	Description *string   `gorm:"type:text"`
	Status      string    `gorm:"not null;size:16;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Room) TableName() string {
	return "rooms"
}

// BeforeCreate assigns ID when missing.
func (r *Room) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
