package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Database struct {
	client *mongo.Client
}

func New() (*Database, error) {
	database := &Database{}

	// Konekcija na mongo cluster!
	mongoUri := os.Getenv("MONGODB_URI")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + mongoUri + "?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	database.client = client

	return database, nil
}

func (database *Database) GetUserRepository() *mongo.Collection {
	userRepository := database.client.Database("usersDB").Collection("users")
	fmt.Println(userRepository.Name())

	return userRepository
}

func (database *Database) GetExperienceRepository() *mongo.Collection {
	experienceRepository := database.client.Database("userDatabase").Collection("experiences")
	fmt.Println(experienceRepository.Name())

	return experienceRepository
}

func (database *Database) Close() error {
	err := database.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

	return nil
}
