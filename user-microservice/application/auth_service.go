package application

import (
	"context"
	"errors"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"user-microservice/model"
)

type AuthService struct {
	store      model.UserStore
	jwtManager *token.JwtManager
}

func NewAuthService(store model.UserStore, manager *token.JwtManager) *AuthService {
	return &AuthService{
		store:      store,
		jwtManager: manager,
	}
}

func (service *AuthService) Login(ctx context.Context, in *userService.CredentialsRequest) (*userService.LoginResponse, error) {
	user, err := service.getUser(ctx, in.Credentials.Username)
	if err == nil && security.BcryptCompareHashAndPassword(user.Password, in.Credentials.Password) == nil {
		jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))
		if err != nil {
			return nil, err
		}
		return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken}, nil
	}
	return nil, errors.New("wrong username or password")
}

func (service *AuthService) getUser(ctx context.Context, username string) (*model.User, error) {
	user, err := service.store.GetByEmail(ctx, username)
	if err == nil {
		return user, nil
	}
	user, err = service.store.GetByUsername(ctx, username)
	if err == nil {
		return user, nil
	}
	return nil, err
}
