package main

import (
	"backend-berita/initializers"
	"backend-berita/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Activity{})
	initializers.DB.AutoMigrate(&models.Room{})
	initializers.DB.AutoMigrate(&models.Mahasiswa{})
	initializers.DB.AutoMigrate(&models.Payment{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Masukan{})
}
