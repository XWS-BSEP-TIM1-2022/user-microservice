package model

import "time"

type User struct {
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Gender      bool      `json:"gender"`
	BirthDate   time.Time `json:"birthDate"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Bio         string    `json:"bio"`
	Skills      []string  `json:"skills"`
	Interests   []string  `json:"interests"`
	Private     bool      `json:"private"`
}
