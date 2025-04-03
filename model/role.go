package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID          uint32       `gorm:"primary_key;auto_increment" json:"id"`
	Name        string       `gorm:"type:varchar(100);not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID   uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}
