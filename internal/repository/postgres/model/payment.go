package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment is the GORM model for billing rows.
type Payment struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey"`
	StudentID           uuid.UUID      `gorm:"type:uuid;not null;index"`
	GroupID             uuid.UUID      `gorm:"type:uuid;not null;index"`
	AmountMinor         int64          `gorm:"not null"`
	Status              string         `gorm:"not null;size:32;index"`
	PaymentDate         *time.Time     `gorm:"type:date"`
	MonthFor            time.Time      `gorm:"type:date;not null;index"`
	PaymentType         string         `gorm:"not null;size:32;index"`
	Comment             *string        `gorm:"type:text"`
	IsFree              bool           `gorm:"not null;default:false;index"`
	DiscountAmountMinor int64          `gorm:"not null;default:0"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate assigns ID when missing.
func (p *Payment) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
