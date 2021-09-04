package main

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {

	port := os.Getenv("PORT")

	// if port == "" {
	// 	port = "8000"
	// }

	router := gin.New()
	router.Use(gin.Logger())

    router.Use(cors.Default())
	
	router.LoadHTMLGlob("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(":" + port)
}
