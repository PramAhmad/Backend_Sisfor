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

func CreateData(c *gin.Context) {
	var body struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Foto  string `json:"foto"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedImage, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar tidak ditemukan"})
		return
	}
	// validate file type
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), uploadedImage.Filename)
	if uploadedImage.Header.Get("Content-Type") != "image/jpeg" && uploadedImage.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar tidak didukung"})
		return
	}

	// move image
	if err := c.SaveUploadedFile(uploadedImage, "images/"+filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar"})
		return
	}
	// vaidate required
	if body.Title == "" || body.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title dan Desc harus diisi"})
		return
	}
	// max image 3mb
	if uploadedImage.Size > 3<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar terlalu besar"})
		return
	}

	activity := models.Activity{
		Title: body.Title,
		Desc:  body.Desc,
		Foto:  os.Getenv("URL_LINK") + "images/" + filename,
	}

	result := initializers.DB.Create(&activity)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": activity})
}

func GetData(c *gin.Context) {
	// get token
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	var activity []models.Activity
	initializers.DB.Find(&activity)
	if len(activity) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": activity})

}

func GetDetail(c *gin.Context) {
	var activity models.Activity
	id := c.Param("id")
	result := initializers.DB.First(&activity, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": activity})
}

func UpdateData(c *gin.Context) {
	var activity models.Activity
	id := c.Param("id")
	result := initializers.DB.First(&activity, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada data"})
		return
	}
	var body struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Foto  string `json:"foto"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity.Title = body.Title
	activity.Desc = body.Desc
	activity.Foto = body.Foto
	result = initializers.DB.Save(&activity)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengupdate data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": activity})

}

func DeleteData(c *gin.Context) {
	var activity models.Activity
	id := c.Param("id")

	result := initializers.DB.Delete(&activity, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "gakda data"})
	}

}
