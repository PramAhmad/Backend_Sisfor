package controllers

import (
	"backend-berita/initializers"
	"backend-berita/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMasukan(c *gin.Context) {
	var body struct {
		Nama  string `json:"isi"`
		Pesan string `json:"pesan"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	masukan := models.Masukan{
		Nama:  body.Nama,
		Pesan: body.Pesan,
	}
	result := initializers.DB.Create(&masukan)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat masukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": masukan})

}

func GetMasukan(c *gin.Context) {

	var masukan []models.Masukan
	initializers.DB.Find(&masukan)
	if len(masukan) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": masukan})

}
