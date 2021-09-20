package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"

	
	"video-share/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
)

var mediaCollection *mongo.Collection = OpenCollection(Client, "media")

func ListBuckets(c *gin.Context){
	
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func ListAllObjectsWithinABucket(c * gin.Context){

}

func AddMedia(c *gin.Context) {

	//TODO:
	//	Go to the users media list and add the media id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var media models.Media

	if err := c.BindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(media)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	media.ID = primitive.NewObjectID()

	result, insertErr := mediaCollection.InsertOne(ctx, media)
	if insertErr != nil {
		msg := fmt.Sprintf("media object was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

//get all media
func GetAllMedia(c *gin.Context){

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	var media []bson.M

	cursor, err := mediaCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	
	if err = cursor.All(ctx, &media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(media)

	c.JSON(http.StatusOK, media)
}

//get single media
func GetSingleMedia(c *gin.Context){

	mediaID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(mediaID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var media bson.M

	if err := mediaCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(media)

	c.JSON(http.StatusOK, media)
}

//deletes all media
func DeleteAllMedia(c * gin.Context){
	// Todo: go through all the users and make all media and accessible media lists empty

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := mediaCollection.DeleteMany(ctx, bson.M{})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)
}

//deletes one media object
func DeleteSingleMedia(c * gin.Context){
	// TODO
	// delete the media file from S3
	// for each of the accessors, go to their accesible media list and delete the media entry
	// go to the users media list and delete the media entry
	mediaID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(mediaID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := mediaCollection.DeleteOne(ctx, bson.M{"_id": docID})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)
}

// TODO
//body statement should take the form:
// {
//	media : media object
//  accessor: {email: email of the accessor we are adding, action: delete or add}
// }
func ChangeAccessor(c * gin.Context){
	// Takes in an id and updates the media
	//	uses the media object to put the new media
	//  looks for the accessor among the users
	//		if add, then adds the media id to the list of access to for that user
	//		if delete, then deletes that accessible media for that user
}

// TODO
func GetAccessibleMedia(c * gin.Context){
	//takes in the body the user and using the email, does a search in all media where the email is in the accessors
}

