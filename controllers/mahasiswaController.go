package controllers

import (
	"Backend_SI/initializers"
	"Backend_SI/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateMahasiswa(c *gin.Context) {
	var body struct {
		Npm       int64  `json:"npm"`
		Nama      string `json:"nama"`
		Kelas     string `json:"kelas"`
		Foto      string `json:"foto"`
		Instagram string `json:"instagram"`
		Twitter   string `json:"twitter"`
		Facebook  string `json:"facebook"`
		Linkedin  string `json:"linkedin"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uploadedImage, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), uploadedImage.Filename)
	// upload image
	if uploadedImage.Header.Get("Content-Type") != "image/jpeg" && uploadedImage.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar tidak didukung"})
		return
	}

	// move image
	if err := c.SaveUploadedFile(uploadedImage, "images/"+filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar"})
		return
	}
	//  required
	if body.Nama == "" || body.Kelas == "" || body.Npm == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NPM,Nama Dan Kelas Harus Di Isi"})
		return
	}
	// max foto 3mb
	if uploadedImage.Size > 3<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar terlalu besar"})
		return
	}
	// save  db
	mahasiswa := models.Mahasiswa{
		Npm:       body.Npm,
		Nama:      body.Nama,
		Kelas:     body.Kelas,
		Foto:      os.Getenv("URL_LINK") + "images/" + filename,
		Instagram: body.Instagram,
		Twitter:   body.Twitter,
		Facebook:  body.Facebook,
		Linkedin:  body.Linkedin,
	}
	result := initializers.DB.Create(&mahasiswa)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat mahasiswa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": mahasiswa})
}

func GetMahasiswa(c *gin.Context) {
	{
		var mahasiswa []models.Mahasiswa
		initializers.DB.Find(&mahasiswa, "is_delete = ?", 0)

		if len(mahasiswa) <= 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Tidak Di Temukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": mahasiswa})
	}
}

func GetDetailMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa
	id := c.Param("id")
	result := initializers.DB.Find(&mahasiswa, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Tidak Di Temukan"})
		return
	}
	if mahasiswa.Is_delete == 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Telah Di Hapus"})
		return
	}
	var payments []models.Payment
	initializers.DB.Find(&payments, "mahasiswa_id = ?", id)

	// kalkulasi total payment
	roomTotals := make(map[string]int64)

	for _, payment := range payments {
		roomID := payment.RoomID
		var room models.Room
		initializers.DB.First(&room, roomID)
		roomName := room.Title
		roomTotals[roomName] += payment.Total
	}

	c.JSON(http.StatusOK, gin.H{
		"result":      mahasiswa,
		"Total Bayar": roomTotals,
	})

}

func UpdateMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa
	id := c.Param("id")

	result := initializers.DB.First(&mahasiswa, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Tidak Di Temukan"})
		return
	}
	var body struct {
		Npm       int64  `json:"npm"`
		Nama      string `json:"nama"`
		Kelas     string `json:"kelas"`
		Foto      string `json:"foto"`
		Instagram string `json:"instagram"`
		Twitter   string `json:"twitter"`
		Facebook  string `json:"facebook"`
		Linkedin  string `json:"linkedin"`
	}
	// jika tidak ada request image pakeyang dulu
	if body.Foto == "" {
		body.Foto = mahasiswa.Foto
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uploadedImage, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), uploadedImage.Filename)
	// upload image
	if uploadedImage.Header.Get("Content-Type") != "image/jpeg" && uploadedImage.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar tidak didukung"})
		return
	}

	// move image
	if err := c.SaveUploadedFile(uploadedImage, "images/"+filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar"})
		return
	}

	// max foto 3mb
	if uploadedImage.Size > 3<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar terlalu besar"})
		return
	}
	// save  db
	mahasiswa.Npm = body.Npm
	mahasiswa.Nama = body.Nama
	mahasiswa.Kelas = body.Kelas
	mahasiswa.Foto = os.Getenv("URL_LINK") + "images/" + filename
	mahasiswa.Instagram = body.Instagram
	mahasiswa.Twitter = body.Twitter
	mahasiswa.Facebook = body.Facebook
	mahasiswa.Linkedin = body.Linkedin
	result = initializers.DB.Save(&mahasiswa)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat mahasiswa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": mahasiswa})

}

func DeleteMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa

	result := initializers.DB.First(&mahasiswa, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Tidak Di Temukan"})
		return
	}
	// update is_delete to 1
	mahasiswa.Is_delete = 1
	result = initializers.DB.Save(&mahasiswa)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menghapus mahasiswa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": mahasiswa})
}
