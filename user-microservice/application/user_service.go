package application

import (
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

func (service *UserService) Get(id primitive.ObjectID) (*model.User, error) {
	return service.store.Get(id)
}

func (service *UserService) GetAll() ([]*model.User, error) {
	return service.store.GetAll()
}

func (service *UserService) Create(user *model.User) (*model.User, error) {
	return service.store.Create(user)
}

func (service *UserService) Update(userId primitive.ObjectID, user *model.User) (*model.User, error) {
	return service.store.Update(userId, user)
}

func (service *UserService) Delete(id primitive.ObjectID) error {
	return service.store.Delete(id)
}

func (service *UserService) DeleteAll() {
	service.store.DeleteAll()
}
