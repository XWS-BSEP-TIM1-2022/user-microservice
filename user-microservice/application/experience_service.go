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

func (service *ExperienceService) GetByUserId(ctx context.Context, id string) ([]*model.Experience, error) {
	Log.Info("Getting all experience for user with id: " + id)
	return service.store.GetExperiencesByUserId(ctx, id)
}

func (service *ExperienceService) Create(ctx context.Context, experience *model.Experience) (*model.Experience, error) {
	Log.Info("Creating new experience for user with id: " + experience.UserId)
	return service.store.CreateExperience(ctx, experience)
}

func (service *ExperienceService) Delete(ctx context.Context, expId primitive.ObjectID) error {
	Log.Info("Deleting experience with id: " + expId.Hex())
	return service.store.DeleteExperience(ctx, expId)

}
