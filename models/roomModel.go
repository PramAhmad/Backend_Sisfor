package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Title     string
	Desc      string
	Payment   []Payment
	Is_delete int8 `gorm:"default:0"`
}
