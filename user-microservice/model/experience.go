package model

import "time"

type Experience struct {
	Name           string    `json:"name"`
	Title          string    `json:"title"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	ExperienceType bool      `json:"experienceType"`
}
