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

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mediaCollection *mongo.Collection = OpenCollection(Client, "media")

func ListBuckets(c *gin.Context) {

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

	c.JSON(http.StatusOK, result.Buckets)
}

func ListBucketContents(c *gin.Context){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := s3.New(sess)

	input := &s3.ListObjectsInput{
		Bucket:  aws.String("video-share-nlatham"),
		MaxKeys: aws.Int64(2),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
	c.JSON(http.StatusOK, result)

}

func AddMedia(c *gin.Context) {

	//TODO:
	//	Go to the users media list and add the media id
	//  Upload the media to S3

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var media models.Media

	if err := c.BindJSON(&media); err != nil {
		msg := fmt.Sprintf("Could not bind media")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(media)
	if validationErr != nil {
		msg := fmt.Sprintf("Could not validate media")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(validationErr)
		return
	}

	media.ID = primitive.NewObjectID()

	userEmail := media.Owner

	result, insertErr := mediaCollection.InsertOne(ctx, media)
	if insertErr != nil {
		msg := fmt.Sprintf("media object was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user); err != nil {
		msg := fmt.Sprintf("Could not get user to add media to")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(err)
		return
	}

	// updatedMedia := append(user.Media, media.ID)
	_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID},
		bson.D{
			{"$push", bson.D{{"media", media.ID}}},
		},
	)
	if updateErr != nil {
		msg := fmt.Sprintf("Could not assign media to user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(updateErr)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result)
}

//get all media
func GetAllMedia(c *gin.Context) {

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

//get list media
func GetListOfMedia(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var json = &struct {
		Media []primitive.ObjectID `form:"media" json:"media" binding:"required"`
	}{}
	if err := c.Bind(json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERR": "WRONG_INPUT"})
		fmt.Println("test", err.Error(), "TEST")
		return
	}
	mediaIDs := json.Media

	var mediaList []bson.M
	for _, id := range mediaIDs {
		var media bson.M
		if err := mediaCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&media); err != nil {
			fmt.Println(err)
		} else {
			mediaList = append(mediaList, media)
		}
	}

	defer cancel()

	fmt.Println(mediaList)

	c.JSON(http.StatusOK, mediaList)
}

//get single media
func GetSingleMedia(c *gin.Context) {

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
func DeleteAllMedia(c *gin.Context) {
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
func DeleteSingleMedia(c *gin.Context) {
	// TODO
	// delete the media file from S3
	// for each of the accessors, go to their accesible media list and delete the media entry
	// go to the users media list and delete the media entry
	mediaID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(mediaID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//get media object
	var media models.Media
	if err := mediaCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&media); err != nil {
		msg := fmt.Sprintf("Could not get media object")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(err)
		return
	}

	//delete media from accessors
	accessors := media.Viewers
	for _, element := range accessors {
		var user models.User

		//get the user
		if err := userCollection.FindOne(ctx, bson.M{"email": element}).Decode(&user); err != nil {
			msg := fmt.Sprintf("Could not get user accessor object")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(err)
			return
		}

		//delete the media from the access to list
		// updatedAccess := findAndDelete(user.Shared, media.ID)
		_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID},
			bson.D{
				{"$pull", bson.D{{"shared", media.ID}}},
			},
		)
		if updateErr != nil {
			msg := fmt.Sprintf("Could not remove accessor for user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(updateErr)
			return
		}

	}

	//delete the media from the owner
	var user models.User
	//get the owner
	if err := userCollection.FindOne(ctx, bson.M{"email": media.Owner}).Decode(&user); err != nil {
		msg := fmt.Sprintf("Could not get owner object")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(err)
		return
	}

	//delete the media from owners media list
	// updatedMedia := findAndDelete(user.Media, media.ID)
	_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID},
		bson.D{
			{"$pull", bson.D{{"media", media.ID}}},
		},
	)
	if updateErr != nil {
		msg := fmt.Sprintf("Could not remove accessor for user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(updateErr)
		return
	}

	//delete the media object
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
//  accessor: id of the accessor we are adding
//	action : delete or add
// }
//TODO
//pass in owner as header and make sure it matches - will need to pull media object to get owner?
func ChangeAccessor(c *gin.Context) {
	// Takes in an id and updates the media
	//	uses the media object to put the new media
	//  looks for the accessor among the users
	//		if add, then adds the media id to the list of access to for that user
	//		if delete, then deletes that accessible media for that user

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var json = &struct {
		Accessor string `form:"accessor" json:"accessor" binding:"required"`
		Action   string `form:"action" json:"action" binding:"required"`
	}{}
	if c.Bind(json) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERR": "WRONG_INPUT"})
		return
	}

	accessor := json.Accessor
	action := json.Action
	// viewers := json.Viewers

	mediaID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(mediaID)

	if action != "add" && action != "delete" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " action not provided"})
		return
	}

	//get the accessor user
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": accessor}).Decode(&user); err != nil {
		msg := fmt.Sprintf("Could not get user accessor object")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(err)
		return
	}

	var result *mongo.UpdateResult
	//update the media object
	if action == "add" {
		addResult, updateErr := mediaCollection.UpdateOne(ctx, bson.M{"_id": docID},
			bson.D{
				{"$push", bson.D{{"viewers", accessor}}},
			},
		)
		if updateErr != nil {
			msg := fmt.Sprintf("Could not add accessor to media")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(updateErr)
			return
		}
		result = addResult
	} else {
		deleteResult, updateErr := mediaCollection.UpdateOne(ctx, bson.M{"_id": docID},
			bson.D{
				{"$pull", bson.D{{"viewers", accessor}}},
			},
		)
		if updateErr != nil {
			msg := fmt.Sprintf("Could not update media object")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(updateErr)
			return
		}
		result = deleteResult
	}

	if action == "delete" {
		//delete the media from the access to list
		// updatedAccess := findAndDelete(user.Shared, docID)
		_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID},
			bson.D{
				{"$pull", bson.D{{"shared", docID}}},
			},
		)
		if updateErr != nil {
			msg := fmt.Sprintf("Could not remove accessor for user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(updateErr)
			return
		}
	}
	if action == "add" {
		//add the media to the access to list
		// updatedAccess := append(user.Shared, docID)
		_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID},
			bson.D{
				{"$push", bson.D{{"shared", docID}}},
			},
		)
		if updateErr != nil {
			msg := fmt.Sprintf("Could not assign media to user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(updateErr)
			return
		}
	}

	defer cancel()

	c.JSON(http.StatusOK, result.ModifiedCount)
}

// TODO
func GetAccessibleMedia(c *gin.Context) {
	//takes in the body the user and using the email, does a search in all media where the email is in the accessors
}

func findAndDelete(s []primitive.ObjectID, itemToDelete primitive.ObjectID) []primitive.ObjectID {
	var new = make([]primitive.ObjectID, 0)
	for _, i := range s {
		if i != itemToDelete {
			new = append(new, i)
		}
	}
	return new
}
