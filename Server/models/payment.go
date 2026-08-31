package models

import(
	"time"
)


type Payment struct {
	ID					uint		`gorm:"primaryKey" json:"id"`
	BidID				uint		`json:"bid_id"`
	Bid					Bid			`gorm:"foreignKey:BidID" json:"bid"`
	Amount				float64		`json:"amount"`
	Status				string		`gorm:"default:pending" json:"status"` // pending, success, failed
	DodoSessionID		string		`json:"dodo_session_id"`
	DodoPaymentID 		string    	`json:"dodo_payment_id"`
	PayerEmail			string		`json:"payer_email"`
	CreatedAt			time.Time	`json:"created_at"`
}

