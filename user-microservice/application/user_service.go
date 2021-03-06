package application

import (
	"context"
	"errors"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
	"user-microservice/model"
	"user-microservice/startup/config"
)

type UserService struct {
	store            model.UserStore
	config           *config.Config
	connectionClient connectionService.ConnectionServiceClient
}

func NewUserService(store model.UserStore, config *config.Config) *UserService {
	return &UserService{
		store:            store,
		config:           config,
		connectionClient: services.NewConnectionClient(fmt.Sprintf("%s:%s", config.ConnectionServiceHost, config.ConnectionServicePort)),
	}
}

func (service *UserService) Get(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	Log.Info("Getting user with id: " + id.Hex())
	return service.store.Get(ctx, id)
}

func (service *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	Log.Info("Getting all users")
	return service.store.GetAll(ctx)
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	Log.Info("Create new user with username: " + user.Username)
	_, err := service.store.GetByUsername(ctx, user.Username)
	if err == nil {
		Log.Warn("Username already exists")
		return nil, errors.New("username already exists")
	}

	if user.Email != "" {
		_, err = service.store.GetByEmail(ctx, user.Email)
		if err == nil {
			Log.Warn("Email already exists")
			return nil, errors.New("email already exists")
		}
	}

	err = service.IsPasswordOk(user.Password)

	if err != nil {
		Log.Warn("Password is to week")
		return nil, err
	}

	hashedPassword, err := security.BcryptGenerateFromPassword(user.Password)
	if err != nil {
		Log.Error("Error with bcrypy library with hashing password")
		return nil, err
	}
	user.Password = hashedPassword

	user.Confirmed = false
	user.ConfirmationId = uuid.New().String()
	err = SendConfirmationMail(ctx, user)
	if err != nil {
		Log.Error("Mail sending for creating new user failed")
		return nil, err
	}

	Log.Info("Created new user with username: " + user.Username)
	return service.store.Create(ctx, user)
}

func (service *UserService) IsPasswordOk(password string) error {
	if len(password) < 8 {
		Log.Warn("Password is too week")
		return errors.New("Password must be atleast 8 characters")
	}

	match, _ := regexp.MatchString("[0-9]", password)
	if !match {
		Log.Warn("Password is too week")
		return errors.New("Password must contain atleast 1 number")
	}

	match, _ = regexp.MatchString("[A-Z]", password)
	if !match {
		Log.Warn("Password is too week")
		return errors.New("Password must contain atleast 1 upper case")
	}

	match, _ = regexp.MatchString("[a-z]", password)
	if !match {
		Log.Warn("Password is too week")
		return errors.New("Password must contain atleast 1 lower case")
	}

	match, _ = regexp.MatchString("[.,<>/?|';:!@#$%^&*()_+=-]", password)
	if !match {
		Log.Warn("Password is too week")
		return errors.New("Password must contain atleast 1 special characher")
	}

	/*for _, commonPassword := range service.config.CommonPasswords {
		if strings.Contains(commonPassword, password) || strings.Contains(password, commonPassword) {
			return errors.New("Password must not be a common password or containts common. (" + commonPassword + ")")
		}
	}*/
	err := service.CheckIsPasswordInCommonPasswords(password)
	if err != nil {
		Log.Warn("Password is too common")
		return err
	}

	return nil
}

func (service *UserService) CheckIsPasswordInCommonPasswords(password string) error {
	numRoutines := 10
	c := make(chan string)

	step := len(service.config.CommonPasswords) / numRoutines

	for i := 0; i < numRoutines-1; i++ {
		go contain(password, service.config.CommonPasswords[step*i:step*(i+1)], c)
	}
	go contain(password, service.config.CommonPasswords[step*(numRoutines-1):len(service.config.CommonPasswords)], c)

	for i := 0; i < numRoutines; i++ {
		common := <-c
		if common != "" {
			return errors.New("Password must not be a common password or containts common. (" + common + ")")
		}
	}
	return nil
}

func contain(password string, subarray []string, c chan string) {
	for _, commonPassword := range subarray {
		if strings.Contains(commonPassword, password) || strings.Contains(password, commonPassword) {
			c <- commonPassword
			return
		}
	}
	c <- ""
}

func (service *UserService) Update(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	Log.Info("Updating user with id:" + userId.Hex())
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return nil, err
	}
	user.Role = existUser.Role
	user.Username = existUser.Username
	user.Password = existUser.Password
	user.Confirmed = existUser.Confirmed
	user.ConfirmationId = existUser.ConfirmationId
	Log.Info("user with id: " + userId.Hex() + " updated")
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) UpdatePassword(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	Log.Warn("Updating password for user with id: " + userId.Hex())
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return nil, err
	}
	user.Role = existUser.Role
	Log.Info("Updated user with id: " + userId.Hex())
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) UpdateUsername(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	Log.Info("Updating username for user with id: " + userId.Hex())
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return nil, err
	}
	user.Role = existUser.Role
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	Log.Info("Deleting user with id: " + id.Hex())
	return service.store.Delete(ctx, id)
}

func (service *UserService) DeleteAll(ctx context.Context) {
	Log.Info("Deleting all users")
	service.store.DeleteAll(ctx)
}

func (service *UserService) Search(ctx context.Context, searchParam string, userId string) ([]*model.User, error) {
	Log.Info("Searching users by user with id:" + userId)
	users, err := service.store.GetAllWithoutAdmins(ctx)

	if err != nil {
		Log.Info("Error occuered in search of user with id: " + userId)
		return nil, err
	}
	searchParam = strings.ToLower(searchParam)
	var retVal []*model.User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Username), searchParam) || strings.Contains(strings.ToLower(user.Name), searchParam) || strings.Contains(strings.ToLower(user.Surname), searchParam) || strings.Contains(strings.ToLower(user.Email), searchParam) {
			isBlocked, err := service.connectionClient.IsBlockedAny(ctx, &connectionService.Block{UserId: userId, BlockUserId: user.Id.Hex()})

			if err != nil {
				Log.Info("Error occuered in search of user with id: " + userId)
				return nil, err
			}
			if !isBlocked.Blocked {
				retVal = append(retVal, user)
			}
		}

	}
	Log.Info("Search succeed of user with id: " + userId)
	return retVal, nil
}

func (service *UserService) IsUserPrivate(ctx context.Context, id primitive.ObjectID) (bool, error) {
	Log.Info("Checking is user with id: " + id.Hex() + " private")
	user, err := service.Get(ctx, id)
	if err != nil {
		Log.Warn("Unexciting user with id: " + id.Hex())
		return false, err
	}
	return user.Private, nil
}

func (service *UserService) RecoverPassword(ctx context.Context, in *userService.NewPasswordRecoveryRequest) error {
	Log.Info("Start password recovering with id" + in.RecoveryId)
	recoveryId, err := primitive.ObjectIDFromHex(in.RecoveryId)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}

	if in.PasswordRecovery.NewPassword != in.PasswordRecovery.ConfirmPassword {
		return errors.New("Passwords are not the same")
	}

	passwordRecoveryRequest, err := service.store.GetPasswordRecoveryRequest(ctx, recoveryId)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}
	now := time.Now()
	if passwordRecoveryRequest.ValidTo.Before(now) {
		Log.Error("Recovery link has expired or used already with id:" + recoveryId.Hex())
		return errors.New("Recovery link has expired or used already")
	}

	userId, err := primitive.ObjectIDFromHex(passwordRecoveryRequest.UserId)
	if err != nil {
		Log.Error("Invalid id:" + passwordRecoveryRequest.UserId)
		return err
	}
	err = service.recoverUserPassword(ctx, userId, in.PasswordRecovery.NewPassword)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}

	service.store.DeletePasswordRecoveryRequest(ctx, passwordRecoveryRequest.Id)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}

	return nil
}

func (service *UserService) recoverUserPassword(ctx context.Context, userId primitive.ObjectID, newPassword string) error {
	Log.Info("Start password recovering for user with id:" + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		return err
	}

	err = service.IsPasswordOk(newPassword)
	if err != nil {
		return err
	}

	hashedPassword, err := security.BcryptGenerateFromPassword(newPassword)
	if err != nil {
		Log.Error("Unexpected error with bcrypt library")
		return err
	}
	user.Password = hashedPassword

	_, err = service.store.Update(ctx, userId, user)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}

	return nil
}

func (service *UserService) ChangeProfilePrivacy(ctx context.Context, userId primitive.ObjectID) error {
	Log.Info("Change profile privacy for user with id:" + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}

	user.Private = !user.Private

	if !user.Private {
		_, err = service.connectionClient.ApproveAllConnection(ctx, &connectionService.UserIdRequest{UserId: userId.Hex()})
		if err != nil {
			Log.Error("Unexpected error with database occurred")
			return err
		}
	}

	_, err = service.Update(ctx, user.Id, user)
	if err != nil {
		Log.Error("Unexpected error with database occurred")
		return err
	}
	Log.Info("Profile with id:" + userId.Hex() + " changed privacy")
	return nil
}
