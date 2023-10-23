package models

import "gorm.io/gorm"

type Masukan struct {
	gorm.Model
	Nama  string `gorm:"type:varchar(100)"`
	Pesan string `gorm:"type:varchar(100)"`
}
