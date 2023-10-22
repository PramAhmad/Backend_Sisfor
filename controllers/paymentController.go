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

func GetDetailPayment(c *gin.Context) {

	var payment models.Payment
	id := c.Param("id")
	result := initializers.DB.Preload("Room").Preload("Mahasiswa").First(&payment, id).Where("is_delete", 0)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	if payment.Is_delete == 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment Telah Di Hapus"})
		return
	}
	mahasiswa := gin.H{
		"id":    payment.Mahasiswa.ID,
		"npm":   payment.Mahasiswa.Npm,
		"nama":  payment.Mahasiswa.Nama,
		"kelas": payment.Mahasiswa.Kelas,
	}
	room := gin.H{
		"id":    payment.Room.ID,
		"title": payment.Room.Title,
		"desc":  payment.Room.Desc,
	}
	json := gin.H{
		"id":          payment.ID,
		"total_bayar": payment.Total,
		"mahasiswa":   mahasiswa,
		"room":        room,
		"addby":       payment.Addby,
		"is_delete":   payment.Is_delete,
	}
	// marsal json

	c.JSON(http.StatusOK, gin.H{"result": json})
}

func GetPaymentByRoom(c *gin.Context) {
	var room models.Room

	id := c.Param("id")

	result := initializers.DB.Preload("Payment").First(&room, "id = ?", id).Order("created_at desc")
	queryroom := gin.H{
		"id":    room.ID,
		"title": room.Title,
		"desc":  room.Desc,
	}
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
		"room":   queryroom,
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
	var mahasiswa models.Mahasiswa
	result := initializers.DB.Preload("Payment").First(&mahasiswa, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	if len(mahasiswa.Payment) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}

	var paymentsWithRoomInfo []PaymentResult
	for _, payment := range mahasiswa.Payment {
		roomID := payment.RoomID
		var room models.Room
		initializers.DB.First(&room, roomID)
		paymentInfo := PaymentResult{
			Total: payment.Total,
			Room:  room,
			Addby: payment.Addby,
		}

		paymentsWithRoomInfo = append(paymentsWithRoomInfo, paymentInfo)

	}
	jsonmahasiswa := gin.H{
		"id":    mahasiswa.ID,
		"npm":   mahasiswa.Npm,
		"nama":  mahasiswa.Nama,
		"kelas": mahasiswa.Kelas,
	}

	c.JSON(http.StatusOK, gin.H{
		"mahasiswa": jsonmahasiswa,
		"count":     len(paymentsWithRoomInfo),
		"result":    paymentsWithRoomInfo,
	})

}

func UpdatePayment(c *gin.Context) {
	var payment models.Payment
	id := c.Param("id")
	result := initializers.DB.Preload("Mahasiswa").Preload("Room").First(&payment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment tidak di temukan"})
		return
	}
	var body struct {
		Total int64 `json:"total"`
		Addby int64 `json:"addby"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment.Total = body.Total
	payment.Addby = body.Addby
	result = initializers.DB.Save(&payment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengupdate data"})
		return
	}
	mahasiswa := gin.H{
		"id":    payment.Mahasiswa.ID,
		"npm":   payment.Mahasiswa.Npm,
		"nama":  payment.Mahasiswa.Nama,
		"kelas": payment.Mahasiswa.Kelas,
	}
	room := gin.H{
		"id":    payment.Room.ID,
		"title": payment.Room.Title,
		"desc":  payment.Room.Desc,
	}

	json := gin.H{
		"id":          payment.ID,
		"total_bayar": payment.Total,
		"mahasiswa":   mahasiswa,
		"room":        room,
		"addby":       payment.Addby,
		"is_delete":   payment.Is_delete,
	}
	c.JSON(http.StatusOK, gin.H{"result": json})
}

func DeletePayment(c *gin.Context) {
	var payment models.Payment
	id := c.Param("id")
	result := initializers.DB.First(&payment, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment tidak di temukan"})
		return
	}
	payment.Is_delete = 1
	result = initializers.DB.Save(&payment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "Payment Berhasil Di Hapus"})
}
