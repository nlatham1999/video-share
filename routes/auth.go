package routes

import (
	"os"
	
    "github.com/gin-gonic/gin"
)

//JWTAuthMiddleware middleware
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        validateToken(c)
        c.Next()
    }
}

func validateToken(c *gin.Context) {
    token := c.Request.Header.Get("X-Auth-Token")

	ApiKey := os.Getenv("API_KEY")

    if token == "" {
        c.AbortWithStatus(401)
    } else if token == ApiKey {
        c.Next()
    } else {
        c.AbortWithStatus(401)
    }
}