package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, userId primitive.ObjectID, user *User) (*User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteAll(ctx context.Context)
	GetAllWithoutAdmins(ctx context.Context) ([]*User, error)

	//experience
	GetExperiencesByUserId(ctx context.Context, id string) ([]*Experience, error)
	CreateExperience(ctx context.Context, experience *Experience) (*Experience, error)
	UpdateExperience(ctx context.Context, experienceId primitive.ObjectID, experience *Experience) (*Experience, error)
	DeleteExperience(ctx context.Context, id primitive.ObjectID) error
}
