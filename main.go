package main

import (
	"os"
	"net/http"

	"video-share/routes"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

// simulate some private data
var secrets = gin.H{
    "foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
    "austin": gin.H{"email": "austin@example.com", "phone": "666"},
    "lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {

	port := os.Getenv("PORT")

	// if port == "" {
	// 	port = "8000"
	// }

	router := gin.New()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"X-Auth-Token"}
	config.AllowAllOrigins = true
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
  
	router.Use(cors.New(config))

	// authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
    //     "foo":    "bar",
    //     "austin": "1234",
    //     "lena":   "hello2",
    //     "manu":   "4321",
    // }))

	// authorized.GET("/secrets", func(c *gin.Context) {
    //     // get user, it was set by the BasicAuth middleware
    //     user := c.MustGet(gin.AuthUserKey).(string)
    //     if secret, ok := secrets[user]; ok {
    //         c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
    //     } else {
    //         c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
    //     }
    // })
	
	router.LoadHTMLGlob("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	needAPIKey := router.Group("/")
	needAPIKey.Use(routes.JWTAuthMiddleware())

	needAPIKey.GET("/users", routes.GetUsers)

	needAPIKey.POST("/user/add", routes.AddUser)

	needAPIKey.DELETE("/users/delete", routes.DeleteAllUsers)

	router.Run(":" + port)
}
