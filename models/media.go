package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         	primitive.ObjectID 	`bson:"_id"`
	Location	*string 			`bson:"location` //s3 location
	Owner       primitive.ObjectID   `json:"owner"`
	Viewers		[]primitive.ObjectID `json:"viewers"`
}