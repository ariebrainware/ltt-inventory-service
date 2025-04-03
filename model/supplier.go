package model

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Address string `gorm:"type:varchar(255);not null" json:"address"`
	Phone   string `gorm:"type:varchar(15);not null" json:"phone"`
	Email   string `gorm:"type:varchar(100);not null" json:"email"`
	Remarks string `gorm:"type:text" json:"remarks"` // Optional remarks
}
