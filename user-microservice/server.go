package main

import (
	"github.com/milossimic/rest/tracer"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"user-microservice/database"
)

type userServer struct {
	databaseClient       *database.Database
	userCollection       *mongo.Collection
	experienceCollection *mongo.Collection
	tracer               opentracing.Tracer
	closer               io.Closer
}

const name = "user_service"

func NewUserServer1() (*userServer, error) {
	databaseClient, err := database.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(tracer)
	return &userServer{
		databaseClient:       databaseClient,
		userCollection:       databaseClient.GetUserRepository(),
		experienceCollection: databaseClient.GetExperienceRepository(),
		tracer:               tracer,
		closer:               closer,
	}, nil
}

func (s *userServer) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *userServer) GetCloser() io.Closer {
	return s.closer
}

func (s *userServer) CloseTracer() error {
	return s.closer.Close()
}

func (s *userServer) CloseDB() error {
	return s.databaseClient.Close()
}
