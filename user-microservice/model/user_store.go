package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	Get(id primitive.ObjectID) (*User, error)
	GetAll() ([]*User, error)
	Create(user *User) (*User, error)
	Update(userId primitive.ObjectID, user *User) (*User, error)
	Delete(id primitive.ObjectID) error
	DeleteAll()
}
