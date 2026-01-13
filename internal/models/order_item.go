package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item within an order
type OrderItem struct {
	ID           string    `gorm:"type:char(36);primaryKey"`
	OrderID      string    `gorm:"type:char(36);not null;index;column:order_id"`
	ProductID    string    `gorm:"type:char(36);not null;index;column:product_id"`
	Quantity     float64   `gorm:"type:decimal(10,2);not null"`
	PricePerUnit float64   `gorm:"type:decimal(10,2);not null;column:price_per_unit"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	Order        Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	Product      Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate generates UUID if not set
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = generateUUID()
	}
	return nil
}
