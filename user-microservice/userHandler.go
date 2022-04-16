package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/milossimic/rest/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mime"
	"net/http"
	"user-microservice/dto"
	"user-microservice/model"
)

func (ts *userServer) createUserHandler(writer http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createUserHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling user create at %s\n", req.URL.Path)),
	)

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(writer, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)
	reqUser, err := dto.DecodeUserBody(ctx, req.Body)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:        reqUser.Name,
		Surname:     reqUser.Surname,
		Email:       reqUser.Email,
		PhoneNumber: reqUser.PhoneNumber,
		Gender:      reqUser.Gender,
		BirthDate:   reqUser.BirthDate,
		Username:    reqUser.Username,
		Password:    reqUser.Password,
		Bio:         reqUser.Bio,
		Skills:      reqUser.Skills,
		Interests:   reqUser.Interests,
		Private:     reqUser.Private}

	//insertResult, err := ts.userCollection.InsertOne(context.TODO(), user)
	insertResult, err := ts.userCollection.InsertOne(ctx, user)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(ctx, writer, insertResult)
	writer.WriteHeader(http.StatusCreated)
}

func (ts *userServer) updateUserHandler(writer http.ResponseWriter, request *http.Request) {
	span := tracer.StartSpanFromRequest("createUserHandler", ts.tracer, request)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling user create at %s\n", request.URL.Path)),
	)

	// Enforce a JSON Content-Type.
	contentType := request.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(writer, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)
	requestUser, err := dto.DecodeUserBody(ctx, request.Body)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:        requestUser.Name,
		Surname:     requestUser.Surname,
		Email:       requestUser.Email,
		PhoneNumber: requestUser.PhoneNumber,
		Gender:      requestUser.Gender,
		BirthDate:   requestUser.BirthDate,
		Username:    requestUser.Username,
		Password:    requestUser.Password,
		Bio:         requestUser.Bio,
		Skills:      requestUser.Skills,
		Interests:   requestUser.Interests,
		Private:     requestUser.Private}

	updatedUser := bson.M{
		"$set": user,
	}

	id, _ := mux.Vars(request)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	filter := bson.D{{"_id", objectId}}
	updateResult, err := ts.userCollection.UpdateOne(ctx, filter, updatedUser)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(ctx, writer, updateResult)
	writer.WriteHeader(http.StatusCreated)
}

func (ts *userServer) getAllUsersHandler(writer http.ResponseWriter, request *http.Request) {
	span := tracer.StartSpanFromRequest("getAllUsersHandler", ts.tracer, request)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all users at %s\n", request.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	findOptions := options.Find()
	// Finding multiple documents returns a cursor
	cursor, err := ts.userCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []*model.User
	// Iterate through the cursor
	//for cursor.Next(context.TODO()) {
	for cursor.Next(ctx) {
		var elem model.User
		err := cursor.Decode(&elem)
		if err != nil {
			tracer.LogError(span, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, &elem)
	}

	if err := cursor.Err(); err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close the cursor once finished
	//cursor.Close(context.TODO())
	cursor.Close(ctx)

	renderJSON(ctx, writer, users)
}

func (ts *userServer) getUserHandler(writer http.ResponseWriter, request *http.Request) {
	span := tracer.StartSpanFromRequest("getUserHandler", ts.tracer, request)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get user at %s\n", request.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	var user model.User

	id, _ := mux.Vars(request)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.D{{"_id", objectId}}
	err = ts.userCollection.FindOne(ctx, filter).Decode(&user)
	//err := ts.userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(ctx, writer, user)
}

func (ts *userServer) deleteUserHandler(writer http.ResponseWriter, request *http.Request) {
	span := tracer.StartSpanFromRequest("deleteUserHandler", ts.tracer, request)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete user at %s\n", request.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id, _ := mux.Vars(request)["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.D{{"_id", objectId}}
	_, err = ts.userCollection.DeleteOne(ctx, filter)
	//_, err := ts.userCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		tracer.LogError(span, err)
		http.Error(writer, err.Error(), http.StatusNotFound)
	}

	writer.WriteHeader(http.StatusNoContent)
}
