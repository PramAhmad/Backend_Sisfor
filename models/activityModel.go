package models

import (
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	Title     string
	Desc      string
	Foto      string
	Is_delete int8 `gorm:"default:0"`
}
