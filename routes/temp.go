package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"log"

	// "server/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
)

var validate = validator.New()

var orderCollection *mongo.Collection = OpenCollection(Client, "orders")

//get all orders
func GetOrders(c *gin.Context){

	log.Println("Test 1")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	var orders []bson.M

	cursor, err := orderCollection.Find(ctx, bson.M{})

	
	log.Println("Test 2")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	
	log.Println("Test 3")
	
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	
	log.Println("Test 4")

	defer cancel()

	fmt.Println(orders)

	c.JSON(http.StatusOK, orders)
}