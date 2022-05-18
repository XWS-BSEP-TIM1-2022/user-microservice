package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PasswordRecoveryRequest struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId  string             `json:"userId"`
	ValidTo time.Time          `json:"validTo"`
}
