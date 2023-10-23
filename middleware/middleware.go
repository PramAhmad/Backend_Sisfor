package middleware

import "github.com/gin-gonic/gin"

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("role") != "admin" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func IsBendahara() gin.HandlerFunc {
	return func(c *gin.Context) {
		// yang bisa akses bendahara dan admin
		if c.GetHeader("role") != "admin" && c.GetHeader("role") != "bendahara" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
func IsMahasiswa() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("role") == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
