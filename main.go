package main

import (
	"backend-berita/controllers"
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
	r.POST("/activity", controllers.CreateData)
	r.GET("/activity", controllers.GetData)
	r.GET("/activity/:id", controllers.GetDetail)
	r.PUT("/activity/:id", controllers.UpdateData)
	r.DELETE("/activity/:id")
	// room
	r.POST("/room", controllers.CreateRoom)
	r.GET("/room", controllers.GetRoom)
	r.GET("/room/:id", controllers.GetDetailRoom)
	r.PUT("/room/:id", controllers.UpdateRoom)
	r.DELETE("/room/:id", controllers.DeleteRoom)

	// mahasiswa
	r.POST("/mahasiswa", controllers.CreateMahasiswa)
	r.GET("/mahasiswa", controllers.GetMahasiswa)
	r.GET("/mahasiswa/:id", controllers.GetDetailMahasiswa)

	//payment
	r.POST("/payment", controllers.CreatePayment)
	r.GET("/payment/:id", controllers.GetPaymentByRoom)
	r.GET("/payment/maha/:id", controllers.GetPaymentByMahasiswa)
	r.Run()

}
