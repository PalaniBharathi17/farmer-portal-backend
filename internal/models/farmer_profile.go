package models

import (
	"time"

	"gorm.io/gorm"
)

// FarmerProfile represents a farmer's profile information
type FarmerProfile struct {
	FarmerID     string    `gorm:"type:char(36);primaryKey;column:farmer_id"`
	FarmName     string    `gorm:"type:varchar(255);not null;column:farm_name"`
	State        string    `gorm:"type:varchar(100);not null"`
	City         string    `gorm:"type:varchar(100);not null"`
	Pincode      string    `gorm:"type:varchar(10);not null;index"`
	Address      string    `gorm:"type:text"`
	FarmSizeAcres float64  `gorm:"type:decimal(10,2);column:farm_size_acres"`
	Rating       float64   `gorm:"type:decimal(3,2);default:0.00"`
	TotalOrders  int       `gorm:"default:0;column:total_orders"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	User         User      `gorm:"foreignKey:FarmerID;references:ID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate generates UUID if not set
func (f *FarmerProfile) BeforeCreate(tx *gorm.DB) error {
	if f.FarmerID == "" {
		f.FarmerID = generateUUID()
	}
	return nil
}
