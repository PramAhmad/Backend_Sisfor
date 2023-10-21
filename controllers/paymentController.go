package controllers

import (
	"backend-berita/initializers"
	"backend-berita/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var body struct {
		Total       int64 `json:"total"`
		MahasiswaID int64 `json:"mahasiswa_id"`
		RoomID      int64 `json:"room_id"`
		Addby       int64 `json:"addby"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Total == 0 || body.MahasiswaID == 0 || body.RoomID == 0 || body.Addby == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua Field Harus Di Isi"})
		return
	}

	// Fetch the Room data id
	var room models.Room
	if err := initializers.DB.First(&room, body.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room Tiadk Di Temukan"})
		return
	}

	var mahasiswa models.Mahasiswa
	if err := initializers.DB.First(&mahasiswa, body.MahasiswaID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa Tidak Di temukan"})
		return
	}

	payment := models.Payment{
		Total:       body.Total,
		MahasiswaID: body.MahasiswaID,
		RoomID:      body.RoomID,
		Addby:       body.Addby,
	}
	result := initializers.DB.Create(&payment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": payment})
}

func GetPaymentByRoom(c *gin.Context) {
	var room models.Room

	id := c.Param("id")

	result := initializers.DB.Preload("Payment").First(&room, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}

	if len(room.Payment) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	// slice mahasiswa
	var paymentsWithMahasiswaInfo []gin.H

	// looping Mahasiswanya
	for _, payment := range room.Payment {
		mahasiswaID := payment.MahasiswaID
		// query untuuk ke id mahasiswa
		var mahasiswa models.Mahasiswa
		initializers.DB.First(&mahasiswa, mahasiswaID)
		// map mahasiswa
		mahasiswaInfo := gin.H{
			"Npm":   mahasiswa.Npm,
			"Nama":  mahasiswa.Nama,
			"Kelas": mahasiswa.Kelas,
		}
		// Create map
		paymentInfo := gin.H{
			"id_payment": payment.ID,
			"total":      payment.Total,
			"addby":      payment.Addby,
			"mahasiswa":  mahasiswaInfo,
		}
		paymentsWithMahasiswaInfo = append(paymentsWithMahasiswaInfo, paymentInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"room":   room.Title,
		"count":  len(paymentsWithMahasiswaInfo),
		"result": paymentsWithMahasiswaInfo,
	})
}

type PaymentResult struct {
	Total int64       `json:"Total"`
	Room  models.Room `json:"Room"`
	Addby int64       `json:"Addby"`
}

func GetPaymentByMahasiswa(c *gin.Context) {
	var payments []models.Payment

	id := c.Param("id")

	result := initializers.DB.Preload("Room").Where("mahasiswa_id = ?", id).Find(&payments)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}

	if len(payments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}

	var paymentResults []PaymentResult

	for _, payment := range payments {
		paymentResult := PaymentResult{
			Total: payment.Total,
			Addby: payment.Addby,
			Room:  payment.Room,
		}
		paymentResults = append(paymentResults, paymentResult)
	}

	var mahasiswa models.Mahasiswa
	initializers.DB.First(&mahasiswa, id)
	if mahasiswa.Is_delete == 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa Telah Di Hapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":     len(paymentResults),
		"mahasiswa": mahasiswa,
		"result":    paymentResults,
	})
}
