package model

import "gorm.io/gorm"

type InventoryMaster struct {
	gorm.Model          // contains ID as product_id (primary key), created_at, updated_at, deleted_at fields
	ProductName string  `gorm:"type:varchar(100);not null" json:"product_name"`
	Category    string  `gorm:"type:varchar(100);not null" json:"category"` // e.g., electronics, clothing, groceries
	Brand       string  `gorm:"type:varchar(100);not null" json:"brand"`
	StockInHand int     `gorm:"not null" json:"stock_in_hand"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
	SupplierID  uint    `gorm:"not null" json:"supplier_id"` // Foreign key for Supplier
	Remarks     string  `gorm:"type:text" json:"remarks"`    // Optional remarks
}
