package models

import(
	"time"
)


type Product struct {
	ID						uint			`gorm:"primaryKey" json:"id"`
	HandleOrURL				string			`gorm:"uniqueIndex;not null" json:"handle_or_url"`
	Name					string			`gorm:"not null" json:"name"`
	Tagline					string			`json:"tagline"`
	LogoURL					string			`json:"logo_url"`
	CategoryID     			uint			`json:"category_id"`
	Category				Category		`gorm:"foreignKey:CategoryID" json:"category"`
	ContactEmail			string			`json:"contact_email"`
	CreatedAt				time.Time		`json:"created_at"`
}

