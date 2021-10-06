package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID         	primitive.ObjectID 	`bson:"_id"`
	Name 		*string				`bson:"name"`
	Location	*string 			`bson:"location` //s3 location
	Owner       *string   			`bson:"owner"`
	Mediatype	*string				`bson:"mediatype"`
	Viewers		[]*string			`bson:"viewers"` //list of users that can access it
	//why email instead of id for viewers? 
	//	- when we display the media, we wont have to to any extra calls to get the email
	//	- emails are going to be a unique identifier as well 
}