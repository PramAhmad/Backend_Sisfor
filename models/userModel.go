package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(100)"`
	Password  string `gorm:"type:varchar(100)"`
	Role      int8   `gorm:"default:0"`
	Is_delete int8   `gorm:"default:0"`
}
