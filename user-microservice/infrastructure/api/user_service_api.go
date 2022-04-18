package api

import (
	"context"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-microservice/application"
)

type UserHandler struct {
	userService.UnimplementedUserServiceServer
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) GetRequest(ctx context.Context, in *userService.UserIdRequest) (*userService.GetResponse, error) {
	id := in.UserId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &userService.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) GetAllRequest(ctx context.Context, in *userService.EmptyRequest) (*userService.GetAllUsers, error) {
	users, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &userService.GetAllUsers{
		Users: []*userService.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) PostRequest(ctx context.Context, in *userService.UserRequest) (*userService.GetResponse, error) {
	user, err := handler.service.Create(mapUserPb(in.User))
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &userService.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) UpdateRequest(ctx context.Context, in *userService.UserRequest) (*userService.GetResponse, error) {
	id, _ := primitive.ObjectIDFromHex(in.UserId)
	user, err := handler.service.Update(id, mapUserPb(in.User))
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &userService.GetResponse{
		User: userPb,
	}
	return response, nil
}
func (handler *UserHandler) DeleteRequest(ctx context.Context, in *userService.UserIdRequest) (*userService.EmptyRequest, error) {
	id, _ := primitive.ObjectIDFromHex(in.UserId)
	handler.service.Delete(id)
	response := &userService.EmptyRequest{}
	return response, nil
}
