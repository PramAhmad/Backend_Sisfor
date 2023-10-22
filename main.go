package main

import (
	"backend-berita/controllers"
	"backend-berita/controllers/auth"
	"backend-berita/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()
	// allow cross for everything

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	// activity
	// secured
	authGroup := r.Group("/")
	authGroup.Use(auth.Auth())

	r.POST("/activity", controllers.CreateData)
	r.GET("/activity", controllers.GetData)
	authGroup.GET("/activity/:id", controllers.GetDetail)
	authGroup.PUT("/activity/:id", controllers.UpdateData)
	authGroup.DELETE("/activity/:id")
	// room
	authGroup.POST("/room", controllers.CreateRoom)
	authGroup.GET("/room", controllers.GetRoom)
	authGroup.GET("/room/:id", controllers.GetDetailRoom)
	authGroup.PUT("/room/:id", controllers.UpdateRoom)
	authGroup.DELETE("/room/:id", controllers.DeleteRoom)

	// mahasiswa
	r.POST("/mahasiswa", controllers.CreateMahasiswa)
	r.GET("/mahasiswa", controllers.GetMahasiswa)
	r.GET("/mahasiswa/:id", controllers.GetDetailMahasiswa)
	r.PUT("/mahasiswa/:id", controllers.UpdateMahasiswa)
	r.GET("/payment/mahasiswa/:id", controllers.GetPaymentByMahasiswa)

	//payment
	r.POST("/payment", controllers.CreatePayment)
	r.GET("/payment/:id", controllers.GetDetailPayment)
	r.GET("/payment/room/:id", controllers.GetPaymentByRoom)
	r.PUT("/payment/:id", controllers.UpdatePayment)
	r.DELETE("/payment/:id", controllers.DeletePayment)

	// Auth user
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.GenerateToken)

	r.Run()

}
