package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"
	// "log"

	"video-share/models"
	// "server/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
)

var validate = validator.New()

var userCollection *mongo.Collection = OpenCollection(Client, "users")

//add a user
func AddUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	user.ID = primitive.NewObjectID()

	result, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("user was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

//get all users
func GetUsers(c *gin.Context){

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	var users []bson.M

	cursor, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(users)

	c.JSON(http.StatusOK, users)
}

//get single user
func GetUser(c *gin.Context){

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user bson.M

	if err := userCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
}

//deletes all users
func DeleteAllUsers(c * gin.Context){

	// TODO: delete all media - no users = no media

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := userCollection.DeleteMany(ctx, bson.M{})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)
}

//deletes one user
func DeleteUser(c * gin.Context){

	// TODO: go through the media list and delete all media for that user

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": docID})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)
}
