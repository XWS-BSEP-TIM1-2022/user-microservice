package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-microservice/model"
)

type ExperienceService struct {
	store model.UserStore
}

func NewExperienceService(store model.UserStore) *ExperienceService {
	return &ExperienceService{
		store: store,
	}
}

func (service *ExperienceService) GetByUserId(ctx context.Context, id primitive.ObjectID) ([]*model.Experience, error) {
	return service.store.GetExperiencesByUserId(ctx, id)
}

func (service *ExperienceService) Create(ctx context.Context, experience *model.Experience) (*model.Experience, error) {
	return service.store.CreateExperience(ctx, experience)
}
