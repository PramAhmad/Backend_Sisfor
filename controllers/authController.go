package controllers

import (
	"Backend_SI/controllers/auth"
	"Backend_SI/initializers"
	"Backend_SI/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, Password, and Role must be provided"})
		return
	}

	// Hash the password before saving it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}
	user := models.User{
		Username: body.Username,
		Password: string(hashedPassword),
		Role:     body.Role,
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": user})
}

type TokenRequest struct {
	Username string `json:"Username"`
	Role     int8   `json:"role"`
	Password string `json:"password"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := initializers.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// validasi login
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// generate token

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

func Logout(context *gin.Context) {
	context.Header("Authorization", "")
	context.Header("role", "")
	context.JSON(http.StatusOK, gin.H{"message": "Logout Success"})
}
