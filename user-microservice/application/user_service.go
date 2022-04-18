package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	return service.store.Create(ctx, user)
}

func (service *UserService) Update(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}

func (service *UserService) DeleteAll(ctx context.Context) {
	service.store.DeleteAll(ctx)
}
