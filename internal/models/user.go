package models

// User represents a user in the system
type User struct {
	Base
	Phone        string `gorm:"type:varchar(20);uniqueIndex;not null"`
	Name         string `gorm:"type:varchar(255);not null"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Role         string `gorm:"type:varchar(20);not null"` // farmer, buyer, admin
	IsVerified   bool   `gorm:"default:false"`
	IsActive     bool   `gorm:"default:true"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
