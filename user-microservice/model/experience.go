package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Experience struct {
	Id             primitive.ObjectID `json:"experienceId"`
	UserId         string             `json:"userId"`
	Name           string             `json:"name"`
	Title          string             `json:"title"`
	StartDate      time.Time          `json:"startDate"`
	EndDate        time.Time          `json:"endDate"`
	ExperienceType bool               `json:"experienceType"`
}
