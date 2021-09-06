package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         	primitive.ObjectID 	`bson:"_id"`
	Email       *string            	`json:"email"`
	Media		[]primitive.ObjectID `json:"media"`
}