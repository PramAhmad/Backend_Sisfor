package main

import (
	"backend-berita/controllers"
	"backend-berita/controllers/auth"
	"backend-berita/initializers"
	"backend-berita/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	// bendahara route
	bendaharaRouter := r.Group("/")
	bendaharaRouter.Use(auth.Auth(), middleware.IsBendahara())
	// admin route
	adminRouter := r.Group("/")
	adminRouter.Use(auth.Auth(), middleware.IsAdmin())
	// mahasiswa route
	mahasiswaRouter := r.Group("/")
	mahasiswaRouter.Use(auth.Auth(), middleware.IsMahasiswa())

	// activity
	r.POST("/activity", controllers.CreateData)
	r.GET("/activity", controllers.GetData)
	r.GET("/activity/:id", controllers.GetDetail)
	bendaharaRouter.PUT("/activity/:id", controllers.UpdateData)
	adminRouter.DELETE("/activity/:id")

	// room
	r.POST("/room", controllers.CreateRoom)
	bendaharaRouter.GET("/room", controllers.GetRoom)
	bendaharaRouter.GET("/room/:id", controllers.GetDetailRoom)
	bendaharaRouter.PUT("/room/:id", controllers.UpdateRoom)
	bendaharaRouter.DELETE("/room/:id", controllers.DeleteRoom)

	// mahasiswa
	adminRouter.POST("/mahasiswa", controllers.CreateMahasiswa)
	adminRouter.GET("/mahasiswa", controllers.GetMahasiswa)
	r.GET("/mahasiswa/:id", controllers.GetDetailMahasiswa)
	r.PUT("/mahasiswa/:id", controllers.UpdateMahasiswa)
	r.GET("/payment/mahasiswa/:id", controllers.GetPaymentByMahasiswa)
	adminRouter.DELETE("/mahasiswa/:id", controllers.DeleteMahasiswa)

	//payment
	r.POST("/payment", controllers.CreatePayment)
	r.GET("/payment/:id", controllers.GetDetailPayment)
	r.GET("/payment/room/:id", controllers.GetPaymentByRoom)
	r.PUT("/payment/:id", controllers.UpdatePayment)
	r.DELETE("/payment/:id", controllers.DeletePayment)

	// masukan
	mahasiswaRouter.GET("/masukan", controllers.GetMasukan)
	r.POST("/masukan", controllers.CreateMasukan)

	// Auth user
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.GenerateToken)
	r.POST("/logout", controllers.Logout)
	r.Run()

}
