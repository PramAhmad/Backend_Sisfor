package models

import (
	"gorm.io/gorm"
)

type Mahasiswa struct {
	gorm.Model
	Npm       int64
	Nama      string `gorm:"type:varchar(100)"`
	Kelas     string `gorm:"type:varchar(100)"`
	Foto      string `gorm:"type:varchar(100)"`
	Instagram string
	Twitter   string
	Facebook  string
	Linkedin  string
	Is_delete int8 `gorm:"default:0"`
	Payment   []Payment
}
