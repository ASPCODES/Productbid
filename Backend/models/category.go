package models

import (
	"time"
)


type Category struct {
	ID 		  		uint				`gorm:"primaryKey" json:"id"`
	Name      		string				`gorm:"not null" json:"name"`
	Slug      		string				`gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt       time.Time			`gorm:"created_at"`		
}

