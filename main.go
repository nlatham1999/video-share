package main

import (
	"net/http"
	"os"

	"video-share/routes"

	// "github.com/nlatham1999/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func CORSMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//         c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//         c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token, content-type")
//         c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

//         if c.Request.Method == "OPTIONS" {
//             c.AbortWithStatus(204)
//             return
//         }

//         c.Next()
//     }
// }

func main() {

	port := os.Getenv("PORT")

	// if port == "" {
	// 	port = "8000"
	// }

	router := gin.New()
	router.Use(gin.Logger())

	// router.Use(CORSMiddleware())

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"X-Auth-Token", "content-type"}
	config.ExposeHeaders = []string{"Content-Length"}
	// config.AllowAllOrigins = true
	config.AllowOrigins = []string{"https://www.videoshare.app"}

	router.Use(cors.New(config))

	router.LoadHTMLGlob("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	needAPIKey := router.Group("/")
	needAPIKey.Use(routes.JWTAuthMiddleware())

	needAPIKey.GET("/users", routes.GetUsers)
	needAPIKey.GET("/user/:id", routes.GetUser)

	needAPIKey.POST("/user/add", routes.AddUser)

	needAPIKey.DELETE("/users/delete", routes.DeleteAllUsers)
	needAPIKey.DELETE("/user/delete/:id", routes.DeleteUser)

	needAPIKey.GET("/media/buckets", routes.ListBuckets)
	needAPIKey.GET("/media/bucket-contents", routes.ListBucketContents)
	needAPIKey.GET("/media/empty-bucket", routes.EmptyBucket)
	needAPIKey.GET("media/all", routes.GetAllMedia)
	needAPIKey.GET("media/single/:id", routes.GetSingleMedia)
	needAPIKey.GET("/media/get-presigned-url/:location", routes.GetPreSignedUrl)

	needAPIKey.POST("media/list", routes.GetListOfMedia) //using post since we are sending data through body
	needAPIKey.POST("/media/add", routes.AddMedia)

	needAPIKey.POST("/media/post-media", routes.UploadMedia)

	needAPIKey.PUT("/media/change-accessor/:id", routes.ChangeAccessor)

	needAPIKey.DELETE("/media/delete-all", routes.DeleteAllMedia)
	needAPIKey.DELETE("media/delete/:id", routes.DeleteSingleMedia)

	router.Run(":" + port)
}
