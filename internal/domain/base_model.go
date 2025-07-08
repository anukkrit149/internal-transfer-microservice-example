package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (b *Base) BeforeUpdate(db *gorm.DB) error {
	now := time.Now()
	b.UpdatedAt = &now
	return nil
}

func (b *Base) BeforeCreate(db *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (Base) TableName() string {
	panic("implement me")
}
