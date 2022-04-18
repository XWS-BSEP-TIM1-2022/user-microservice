package api

import (
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"user-microservice/model"
)

func mapUser(user *model.User) *userService.User {
	userPb := &userService.User{
		Id:          user.Id.Hex(),
		Name:        user.Name,
		Surname:     user.Surname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		BirthDate:   user.BirthDate.String(),
		Username:    user.Username,
		Password:    "",
		Bio:         user.Bio,
		Skills:      user.Skills,
		Interests:   user.Interests,
		Private:     user.Private,
	}
	return userPb
}
func mapUserPb(userPb *userService.User) *model.User {
	id, _ := primitive.ObjectIDFromHex(userPb.Id)
	t, _ := time.Parse(userPb.BirthDate, userPb.BirthDate)
	user := &model.User{
		Id:          id,
		Name:        userPb.Name,
		Surname:     userPb.Surname,
		Email:       userPb.Email,
		PhoneNumber: userPb.PhoneNumber,
		Gender:      userPb.Gender,
		BirthDate:   t,
		Username:    userPb.Username,
		Password:    "",
		Bio:         userPb.Bio,
		Skills:      userPb.Skills,
		Interests:   userPb.Interests,
		Private:     userPb.Private,
	}
	return user
}
