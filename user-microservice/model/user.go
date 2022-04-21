package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Surname     string             `json:"surname"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phoneNumber"`
	Gender      Gender             `json:"gender"`
	BirthDate   time.Time          `json:"birthDate"`
	Username    string             `json:"username"`
	Password    string             `json:"password"`
	Bio         string             `json:"bio"`
	Skills      []string           `json:"skills"`
	Interests   []string           `json:"interests"`
	Private     bool               `json:"private"`
	Role        UserRole
}

type UserRole string

const (
	ADMIN UserRole = "ADMIN"
	USER           = "USER"
)

type Gender int64

const (
	MALE Gender = iota
	FEMALE
)
