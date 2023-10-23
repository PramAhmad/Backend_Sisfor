package controllers

import (
	"backend-berita/initializers"
	"backend-berita/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	var body struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room := models.Room{
		Title: body.Title,
		Desc:  body.Desc,
	}
	result := initializers.DB.Create(&room)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal Membuat Room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": room})
}

func GetRoom(c *gin.Context) {

	var room []models.Room
	initializers.DB.Find(&room, "is_delete = ?", 0)
	if len(room) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room Tidak Di Temukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": room})
}

func GetDetailRoom(c *gin.Context) {
	if c.GetHeader("Role") != "admin" || c.GetHeader("Role") != "bendahara" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var room models.Room
	id := c.Param("id")
	result := initializers.DB.First(&room, id)
	// cek jika room sudah terdelete
	if room.Is_delete == 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room Telah Di Hapus"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room Tidak Di Temukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": room})

}

func UpdateRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")

	result := initializers.DB.First(&room, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room Tidak Di Temukan"})
		return
	}
	var body struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room.Title = body.Title
	room.Desc = body.Desc
	result = initializers.DB.Save(&room)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal Mengupdate Room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": room})
}

func DeleteRoom(c *gin.Context) {
	var room models.Room

	result := initializers.DB.First(&room, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room Tidak Di Temukan"})
		return
	}
	// update is_delete to 1
	room.Is_delete = 1
	result = initializers.DB.Save(&room)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal Menghapus Room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "Room Berhasil Di Hapus"})
}
