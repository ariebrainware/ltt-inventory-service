package model

import "gorm.io/gorm"

type InventoryMaster struct {
	gorm.Model          // contains ID as product_id (primary key), created_at, updated_at, deleted_at fields
	ProductName string  `gorm:"type:varchar(100);not null"`
	Category    string  `gorm:"type:varchar(100);not null"` // e.g., electronics, clothing, groceries
	Brand       string  `gorm:"type:varchar(100);not null"`
	StockInHand int     `gorm:"not null"`
	UnitPrice   float64 `gorm:"not null"`
	SupplierID  uint    `gorm:"not null"`  // Foreign key for Supplier
	Remarks     string  `gorm:"type:text"` // Optional remarks
}
