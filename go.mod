module video-share

go 1.16

require (
	github.com/aws/aws-sdk-go v1.40.42
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0
	go.mongodb.org/mongo-driver v1.7.2
)

replace github.com/gin-contrib/cors => ../cors
