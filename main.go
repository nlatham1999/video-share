package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"video-share/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/codegangsta/negroni"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	// "github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "https://videoshare/api"
		fmt.Println(token)
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		// Verify 'iss' claim
		iss := os.Getenv("AUTH0_DOMAIN")
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMid := *jwtMiddleware
		if err := jwtMid.CheckJWT(c.Writer, c.Request); err != nil {
			fmt.Println("ERROR: ", err)
			c.AbortWithStatus(401)
		}
	}
}

func main() {

	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"X-Auth-Token", "content-type", "authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	router.LoadHTMLGlob("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	needAPIKey := router.Group("/")
	// needAPIKey.Use(routes.JWTAuthMiddleware())
	fmt.Println("TESTING")
	needAPIKey.Use(checkJWT())

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

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")

	if err != nil {
		fmt.Println("error1")
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		fmt.Println("error2")
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		fmt.Println("Unable to find appropriate key")
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
