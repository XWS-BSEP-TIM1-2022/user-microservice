package api

import (
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
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
		Gender:      int64(user.Gender),
		BirthDate:   user.BirthDate.String(),
		Username:    user.Username,
		Password:    "",
		Bio:         user.Bio,
		Skills:      user.Skills,
		Interests:   user.Interests,
		Private:     user.Private,
		Role:        string(user.Role),
		TFAEnabled:  user.TFAEnabled,
	}
	return userPb
}
func mapUserPb(userPb *userService.User) *model.User {
	id, _ := primitive.ObjectIDFromHex(userPb.Id)
	t := time.Now()
	if userPb.BirthDate != "" {
		dateString := strings.Split(userPb.BirthDate, "T")
		date := strings.Split(dateString[0], "-")
		year, _ := strconv.Atoi(date[0])
		month, _ := strconv.Atoi(date[1])
		day, _ := strconv.Atoi(date[2])
		t = time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	}
	user := &model.User{
		Id:          id,
		Name:        userPb.Name,
		Surname:     userPb.Surname,
		Email:       userPb.Email,
		PhoneNumber: userPb.PhoneNumber,
		Gender:      model.Gender(userPb.Gender),
		BirthDate:   t,
		Username:    userPb.Username,
		Password:    userPb.Password,
		Bio:         userPb.Bio,
		Skills:      userPb.Skills,
		Interests:   userPb.Interests,
		Private:     userPb.Private,
		Role:        model.UserRole(userPb.Role),
	}
	return user
}

func mapExperience(experience *model.Experience) *userService.Experience {
	experiencePb := &userService.Experience{
		Id:             experience.Id.Hex(),
		UserId:         experience.UserId,
		Name:           experience.Name,
		Title:          experience.Title,
		ExperienceType: experience.ExperienceType,
		StartDate:      experience.StartDate.String(),
		EndDate:        experience.EndDate.String(),
	}
	return experiencePb
}

func mapExperiencePb(experiencePb *userService.Experience) *model.Experience {
	id, _ := primitive.ObjectIDFromHex(experiencePb.Id)
	start := time.Now()
	if experiencePb.StartDate != "" {
		dateString := strings.Split(experiencePb.StartDate, "T")
		date := strings.Split(dateString[0], "-")
		year, _ := strconv.Atoi(date[0])
		month, _ := strconv.Atoi(date[1])
		day, _ := strconv.Atoi(date[2])
		start = time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	}

	end := time.Now()
	if experiencePb.StartDate != "" {
		dateString := strings.Split(experiencePb.StartDate, "T")
		date := strings.Split(dateString[0], "-")
		year, _ := strconv.Atoi(date[0])
		month, _ := strconv.Atoi(date[1])
		day, _ := strconv.Atoi(date[2])
		end = time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	}
	experience := &model.Experience{
		Id:             id,
		Name:           experiencePb.Name,
		Title:          experiencePb.Title,
		UserId:         experiencePb.UserId,
		ExperienceType: experiencePb.ExperienceType,
		StartDate:      start,
		EndDate:        end,
	}
	return experience
}
