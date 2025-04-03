package model

import (
	"database/sql"
	"time"
)

type DeletedAt sql.NullTime

type InventoryDetails struct {
	TransactionID          uint   `gorm:"primaryKey;not null" json:"transaction_id"`
	InventoryMasterID      uint   `gorm:"not null" json:"inventory_master_id"` // as product_id, the naming just for clarity
	TransactionType        string `gorm:"not null" json:"transaction_type"`    // e.g., purchase, sale, return
	Quantity               int    `gorm:"not null" json:"quantity"`
	TransactionDate        string `gorm:"type:varchar(100);not null" json:"transaction_date"`        // e.g., "2023-10-01"
	TransactionDescription string `gorm:"type:varchar(255);not null" json:"transaction_description"` // e.g., order number, customer name
	TransactionAmount      int    `gorm:"not null" json:"transaction_amount"`                        // Total amount for the transaction
	TransactionStatus      string `gorm:"type:varchar(50);not null" json:"transaction_status"`       // e.g., completed, pending, cancelled
	TransactionNotes       string `gorm:"type:text" json:"transaction_notes"`                        // Optional notes for the transaction
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              DeletedAt `gorm:"index"`
}
