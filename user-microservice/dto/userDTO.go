package dto

import (
	"context"
	"encoding/json"
	"github.com/milossimic/rest/tracer"
	"io"
	"time"
)

type RequestUser struct {
	Id          string    `json:"id"`
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

func DecodeUserBody(ctx context.Context, r io.Reader) (*RequestUser, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt RequestUser
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}
