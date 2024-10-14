package model

import "time"

type Card struct {
	ID             uint64 `gorm:"primaryKey"`
	CardHolderName string
	CardNumber     string `gorm:"index"`
	CVV            string
	ExpiryDate     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
