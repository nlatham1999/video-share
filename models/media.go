package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID         	primitive.ObjectID 	`bson:"_id"`
	Name 		*string				`bson:"name"`
	Location	*string 			`bson:"location` //s3 location
	Owner       primitive.ObjectID   `bson:"owner"`
	Viewers		[]*string 			`bson:"viewers"` //list of emails that can access it
}