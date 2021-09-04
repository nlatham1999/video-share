package main

import (
	"os"

	"main/routes"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

    router.Use(cors.Default())

	
	router.GET("/", routes.TempRequest)
	router.Run(":" + port)
}
