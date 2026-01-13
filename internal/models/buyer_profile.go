package models

import (
	"time"

	"gorm.io/gorm"
)

// BuyerProfile represents a buyer's profile information
type BuyerProfile struct {
	BuyerID     string    `gorm:"type:char(36);primaryKey;column:buyer_id"`
	BuyerType   string    `gorm:"type:enum('individual','restaurant','vendor');not null;column:buyer_type"`
	BusinessName string   `gorm:"type:varchar(255);column:business_name"`
	GSTNumber   string    `gorm:"type:varchar(50);column:gst_number"`
	State       string    `gorm:"type:varchar(100);not null"`
	City        string    `gorm:"type:varchar(100);not null"`
	Pincode     string    `gorm:"type:varchar(10);not null;index"`
	Address     string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	User        User      `gorm:"foreignKey:BuyerID;references:ID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate generates UUID if not set
func (b *BuyerProfile) BeforeCreate(tx *gorm.DB) error {
	if b.BuyerID == "" {
		b.BuyerID = generateUUID()
	}
	return nil
}
