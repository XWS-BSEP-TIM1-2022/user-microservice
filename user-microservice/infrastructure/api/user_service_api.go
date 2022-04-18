package api

import (
	"context"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
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
	span := tracer.StartSpanFromContextMetadata(ctx, "GetRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := in.UserId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(ctx, objectId)
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
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	users, err := handler.service.GetAll(ctx)
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
	span := tracer.StartSpanFromContextMetadata(ctx, "PostRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := handler.service.Create(ctx, mapUserPb(in.User))
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
	span := tracer.StartSpanFromContextMetadata(ctx, "UpdateRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id, _ := primitive.ObjectIDFromHex(in.UserId)
	user, err := handler.service.Update(ctx, id, mapUserPb(in.User))
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
	span := tracer.StartSpanFromContextMetadata(ctx, "DeleteRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id, _ := primitive.ObjectIDFromHex(in.UserId)
	handler.service.Delete(ctx, id)
	response := &userService.EmptyRequest{}
	return response, nil
}
