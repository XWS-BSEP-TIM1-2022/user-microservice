package api

import (
	"context"
	"errors"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-microservice/application"
	"user-microservice/model"
)

type UserHandler struct {
	userService.UnimplementedUserServiceServer
	service           *application.UserService
	authService       *application.AuthService
	experienceService *application.ExperienceService
}

func NewUserHandler(
	service *application.UserService,
	authService *application.AuthService,
	experienceService *application.ExperienceService) *UserHandler {
	return &UserHandler{
		service:           service,
		authService:       authService,
		experienceService: experienceService,
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

func (handler *UserHandler) GetAllRequest(ctx context.Context, in *userService.EmptyRequest) (*userService.UsersResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	users, err := handler.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := &userService.UsersResponse{
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

	if in.User.Password != in.User.ConfirmPassword {
		return nil, errors.New("passwords not match")
	}

	if in.User.Name == "" || in.User.Surname == "" || in.User.Email == "" || in.User.BirthDate == "" || in.User.Username == "" || in.User.Password == "" {
		return nil, errors.New("not entered required fields")
	}

	userFromRequest := mapUserPb(in.User)
	userFromRequest.Role = model.USER
	user, err := handler.service.Create(ctx, userFromRequest)
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &userService.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) PostAdminRequest(ctx context.Context, in *userService.UserRequest) (*userService.GetResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "PostAdminRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	userFromRequest := mapUserPb(in.User)
	userFromRequest.Role = model.ADMIN
	user, err := handler.service.Create(ctx, userFromRequest)
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

	if in.User.Name == "" || in.User.Surname == "" || in.User.Email == "" || in.User.BirthDate == "" || in.User.Username == "" {
		return nil, errors.New("not entered required fields")
	}

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

func (handler *UserHandler) LoginRequest(ctx context.Context, in *userService.CredentialsRequest) (*userService.LoginResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "LoginRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return handler.authService.Login(ctx, in)
}

func (handler *UserHandler) SearchUsersRequest(ctx context.Context, in *userService.SearchRequest) (*userService.UsersResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "SearchUsersRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	users, err := handler.service.Search(ctx, in.SearchParam)
	if err != nil {
		return nil, err
	}
	response := &userService.UsersResponse{
		Users: []*userService.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) IsUserAuthenticated(ctx context.Context, in *userService.AuthRequest) (*userService.AuthResponse, error) {
	userRole, err := handler.authService.IsAuthenticated(ctx, in.JwtToken)
	if err != nil {
		return nil, err
	}
	return &userService.AuthResponse{UserRole: string(userRole)}, nil
}

func (handler *UserHandler) UpdatePasswordRequest(ctx context.Context, in *userService.NewPasswordRequest) (*userService.GetResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "UpdatePasswordRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id := in.GetNewPassword().GetUserId()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(ctx, objectId)
	if err != nil {
		return nil, err
	}
	good, err := handler.authService.CheckPassword(ctx, in.NewPassword.Password, user.Id)
	if err != nil {
		return nil, err
	}

	user.Password, _ = security.BcryptGenerateFromPassword(in.NewPassword.Password)
	if good {
		_, err = handler.service.UpdatePassword(ctx, user.Id, user)
		if err != nil {
			return nil, err
		}
		response := &userService.GetResponse{
			User: mapUser(user),
		}
		return response, nil
	}
	return nil, err
}

func (handler *UserHandler) PostExperienceRequest(ctx context.Context, in *userService.NewExperienceRequest) (*userService.NewExperienceResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "PostExperienceRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	experienceFromRequest := mapExperiencePb(in.Experience)

	experience, err := handler.experienceService.Create(ctx, experienceFromRequest)

	if err != nil {
		return nil, err
	}
	experiencePb := mapExperience(experience)
	response := &userService.NewExperienceResponse{
		Experience: experiencePb,
	}
	return response, nil
}

func (handler *UserHandler) GetAllUsersExperienceRequest(ctx context.Context, in *userService.ExperienceRequest) (*userService.ExperienceResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllUsersExperienceRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	/*
		userId, err := primitive.ObjectIDFromHex(in.GetUserId())
		if err != nil {
			panic(err)
		}
	*/

	experiences, _ := handler.experienceService.GetByUserId(ctx, in.UserId)

	response := &userService.ExperienceResponse{
		Experiences: []*userService.Experience{},
	}
	for _, experience := range experiences {
		current := mapExperience(experience)
		response.Experiences = append(response.Experiences, current)
	}
	return response, nil
}

func (handler *UserHandler) DeleteExperienceRequest(ctx context.Context, in *userService.DeleteUsersExperienceRequest) (*userService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "DeleteExperienceRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id, _ := primitive.ObjectIDFromHex(in.ExperienceId)
	handler.experienceService.Delete(ctx, id)
	response := &userService.EmptyRequest{}
	return response, nil
}

func (handler *UserHandler) IsUserPrivateRequest(ctx context.Context, in *userService.UserIdRequest) (*userService.PrivateResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "IsUserPrivateRequest")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	id, _ := primitive.ObjectIDFromHex(in.UserId)

	isPrivate, err := handler.service.IsUserPrivate(ctx, id)

	if err != nil {
		return nil, err
	}

	return &userService.PrivateResponse{IsPrivate: isPrivate}, nil
}
