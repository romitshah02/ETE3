package handlers

import (
	"log"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	log.Println("Router Setup Started.")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "authorization", "Content-Type", "content-type"}
	r.Use(cors.New(config))
	//tokenmiddleware := r.Group("/").Use(middleware.AuthMiddleware())

	r.POST("/user/register", Register)
	r.POST("/user/login", Login)
	r.POST("/movie/add", AddMovie)
	r.POST("/show/add", AddShowHandler)
	r.GET("/movie/get", GetAllMovies)
	r.GET("/show/get/:movie_id", GetShowsByMovie)
	r.POST("/show/book", BookSeats)
	r.GET("/show/seats/get/:show_id")

	return r
}
