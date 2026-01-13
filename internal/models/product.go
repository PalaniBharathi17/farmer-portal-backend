package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product listing by a farmer
type Product struct {
	ID           string    `gorm:"type:char(36);primaryKey"`
	FarmerID     string    `gorm:"type:char(36);not null;index;column:farmer_id"`
	CropName     string    `gorm:"type:varchar(255);not null;index;column:crop_name"`
	Quantity     float64   `gorm:"type:decimal(10,2);not null"`
	Unit         string    `gorm:"type:varchar(50);not null"`
	PricePerUnit float64   `gorm:"type:decimal(10,2);not null;column:price_per_unit"`
	State        string    `gorm:"type:varchar(100);not null"`
	City         string    `gorm:"type:varchar(100);not null"`
	Pincode      string    `gorm:"type:varchar(10);not null;index"`
	Status       string    `gorm:"type:enum('active','closed','sold');default:'active'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	Farmer       User      `gorm:"foreignKey:FarmerID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// BeforeCreate generates UUID if not set
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}
