package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents an order placed by a buyer
type Order struct {
	ID           string       `gorm:"type:char(36);primaryKey"`
	BuyerID      string       `gorm:"type:char(36);not null;index;column:buyer_id"`
	FarmerID     string       `gorm:"type:char(36);not null;index;column:farmer_id"`
	Status       string       `gorm:"type:enum('pending','accepted','rejected','shipped','delivered');default:'pending'"`
	DeliveryMode string       `gorm:"type:enum('pickup','courier');not null;column:delivery_mode"`
	TotalAmount  float64      `gorm:"type:decimal(10,2);not null;column:total_amount"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`
	Buyer        User         `gorm:"foreignKey:BuyerID;references:ID;constraint:OnDelete:CASCADE"`
	Farmer       User         `gorm:"foreignKey:FarmerID;references:ID;constraint:OnDelete:CASCADE"`
	OrderItems   []OrderItem  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate generates UUID if not set
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = generateUUID()
	}
	return nil
}
