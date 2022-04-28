package application

import (
	"context"
	"errors"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"user-microservice/model"
)

type UserService struct {
	store model.UserStore
}

func NewUserService(store model.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Get(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return service.store.Get(ctx, id)
}

func (service *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return service.store.GetAll(ctx)
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := service.store.GetByUsername(ctx, user.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	if user.Email != "" {
		_, err = service.store.GetByEmail(ctx, user.Email)
		if err == nil {
			return nil, errors.New("email already exists")
		}
	}

	hashedPassword, err := security.BcryptGenerateFromPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return service.store.Create(ctx, user)
}

func (service *UserService) Update(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Role = existUser.Role
	user.Username = existUser.Username
	user.Password = existUser.Password
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) UpdatePassword(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Role = existUser.Role
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}

func (service *UserService) DeleteAll(ctx context.Context) {
	service.store.DeleteAll(ctx)
}

func (service *UserService) Search(ctx context.Context, searchParam string) ([]*model.User, error) {
	users, err := service.store.GetAllWithoutAdmins(ctx)

	if err != nil {
		return nil, err
	}
	searchParam = strings.ToLower(searchParam)
	var retVal []*model.User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Username), searchParam) || strings.Contains(strings.ToLower(user.Name), searchParam) || strings.Contains(strings.ToLower(user.Surname), searchParam) || strings.Contains(strings.ToLower(user.Email), searchParam) {
			retVal = append(retVal, user)
		}
	}
	return retVal, nil
}

func (service *UserService) IsUserPrivate(ctx context.Context, id primitive.ObjectID) (bool, error) {
	user, err := service.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return user.Private, nil
}
