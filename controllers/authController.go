package controllers

import (
	"backend-berita/controllers/auth"
	"backend-berita/initializers"
	"backend-berita/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     int8   `json:"role"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// required
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, Password Dan Role Harus Di Isi"})
		return
	}
	// save  db
	user := models.User{
		Username: body.Username,
		Password: body.Password,
		Role:     body.Role,
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": user})

}

type TokenRequest struct {
	Username string `json:"Username"`
	Role     int8   `json:"role"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := initializers.DB.Where("username = ?", request.Username).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Role, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// set header role

	// switchcase the role 1=admin 2=bendahara 3= mahasiswa
	var role string
	switch user.Role {
	case 1:
		role = "admin"
	case 2:
		role = "bendahara"
	case 3:
		role = "mahasiswa"
	}

	context.Header("role", role)
	context.Header("Authorization", tokenString)
	context.JSON(http.StatusOK, gin.H{"token": tokenString,
		"role": role})
}

func DestroyToken(context *gin.Context) {

	context.JSON(http.StatusOK, gin.H{"message": "Token destroyed"})
}
