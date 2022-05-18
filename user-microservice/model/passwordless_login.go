package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PasswordlessLogin struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId       string             `json:"userId"`
	CreationTime time.Time          `json:"creationTime"`
}
