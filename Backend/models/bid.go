package models

import(
	"time"
)


type Bid struct {
	ID				uint			`gorm:"primaryKey" json:"id"`
	ProductID		uint			`json:"product_id"`
	Product			Product			`gorm:"foreignKey:productID" json:"product"`
	CategoryID		uint			`json:"category_id"`
	Category		Category		`gorm:"foreignKey:CategoryID" json:"category"`
	Amount			float64			`gorm:"not null" json:"amount"`
	IsActive		bool			`gorm:"default:false" json:"is_active"`
	CreatedAt       time.Time		`json:"created_at"`
}

