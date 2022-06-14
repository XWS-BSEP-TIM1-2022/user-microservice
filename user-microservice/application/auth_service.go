package application

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"github.com/dgryski/dgoogauth"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"time"
	"user-microservice/model"

	_ "io/ioutil"

	qr "rsc.io/qr"
)

type AuthService struct {
	store      model.UserStore
	jwtManager *token.JwtManager
}

var Log = logrus.New()

func NewAuthService(store model.UserStore, manager *token.JwtManager) *AuthService {
	return &AuthService{
		store:      store,
		jwtManager: manager,
	}
}

func (service *AuthService) Login(ctx context.Context, in *userService.CredentialsRequest) (*userService.LoginResponse, error) {
	Log.Info("User with username: " + in.Credentials.Username + " try to login")
	user, err := service.getUser(ctx, in.Credentials.Username)
	if err == nil && security.BcryptCompareHashAndPassword(user.Password, in.Credentials.Password) == nil {
		if user.Confirmed == false {
			Log.Warn("User with username: " + in.Credentials.Username + " entered wrong password")
			return nil, errors.New("unconfirmed registration")
		}

		jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))
		if err != nil {
			Log.Error("User with username: " + in.Credentials.Username + " get error while generating JWT")
			return nil, err
		}

		if user.TFAEnabled {
			Log.Info("User with username: " + in.Credentials.Username + " started TFA")
			return &userService.LoginResponse{UserId: user.Id.Hex()}, nil
		}
		Log.Info("User with username: " + in.Credentials.Username + " logged in")
		return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken, IsPrivate: user.Private}, nil
	}
	return nil, errors.New("wrong username or password")
}

func (service *AuthService) getUser(ctx context.Context, username string) (*model.User, error) {
	Log.Info("Getting user by id or email: " + username)
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

func (service *AuthService) ConfirmRegistration(ctx context.Context, in *userService.ConfirmationRequest) (*userService.ConfirmationResponse, error) {
	Log.Info("Confirmation registration with id : " + in.ConfirmationId)
	user, err := service.store.GetByConfirmationId(ctx, in.ConfirmationId)
	if err != nil {
		Log.Error("User with given confirmationId: " + in.ConfirmationId + " does not exist")
		return &userService.ConfirmationResponse{ResponseMessage: "user with given confirmationId does not exist"}, err
	}
	user.Confirmed = true
	_, err = service.store.Update(ctx, user.Id, user)
	if err != nil {
		return nil, err
	}
	Log.Info("Successfully confirmed registration with id : " + in.ConfirmationId)
	return &userService.ConfirmationResponse{ResponseMessage: "successfully confirmed registration"}, nil
}

func (service *AuthService) IsAuthenticated(ctx context.Context, jwtToken string) (model.UserRole, error) {
	ok := service.jwtManager.IsUserAuthorized(jwtToken)
	if ok != nil {
		Log.Warn("Unauthorized user")
		return "", ok
	}
	userRole, err := service.jwtManager.GetRoleFromToken(jwtToken)
	if err != nil {
		Log.Warn("Jwt is not valid")
		return "", err
	}
	return model.UserRole(userRole), nil
}

func (service *AuthService) CheckPassword(ctx context.Context, password string, userId primitive.ObjectID) (bool, error) {
	Log.Info("Checking password of user with id:" + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err == nil && security.BcryptCompareHashAndPassword(user.Password, password) != nil {
		Log.Info("Valid password of user with id:" + userId.Hex())
		return true, nil
	}
	Log.Warn("Invalid password of user with id:" + userId.Hex())
	return false, errors.New("wrong password")
}

func (service *AuthService) CheckUsername(ctx context.Context, username string) (bool, error) {
	Log.Info("Checking username: " + username)
	_, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		Log.Info("Valid username: " + username)
		return true, nil
	}
	Log.Info("Invalid username: " + username)
	return false, nil
}

func (service *AuthService) GetQR2FA(ctx context.Context, userId primitive.ObjectID) ([]byte, error) {
	Log.Info("Getting QR2FA for user with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)

	if err != nil {
		Log.Warn("User with id: " + userId.Hex() + " doesn't exits")
		return nil, err
	}

	secret := make([]byte, 10)
	_, err = rand.Read(secret)
	if err != nil {
		Log.Error("Secret making stop working")
		panic(err)
	}

	user.TFASecret = base32.StdEncoding.EncodeToString(secret)
	user.TFAEnabled = false
	service.store.Update(ctx, userId, user)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		Log.Error("Server stopped working due to bad url creating service")
		panic(err)
	}

	URL.Path += "/" + url.PathEscape("Dislinkt") + ":" + url.PathEscape(user.Username)

	params := url.Values{}
	params.Add("secret", user.TFASecret)
	params.Add("issuer", "Dislinkt")

	URL.RawQuery = params.Encode()
	fmt.Printf("URL is %s\n", URL.String())

	code, err := qr.Encode(URL.String(), qr.Q)

	if err != nil {
		return nil, err
	}
	Log.Info("Returning QR2FA for user with id: " + userId.Hex())
	return code.PNG(), nil
}

func (service *AuthService) Verify2fa(ctx context.Context, userId primitive.ObjectID, code string) (*userService.LoginResponse, error) {
	Log.Info("Verifying 2FA for user with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)

	if err != nil {
		Log.Warn("Invalid 2FA for user with id: " + userId.Hex())
		return nil, err
	}

	otpc := &dgoogauth.OTPConfig{
		Secret:      user.TFASecret,
		WindowSize:  3,
		HotpCounter: 0,
		// UTC:         true,
	}
	val, err := otpc.Authenticate(code)
	if err != nil {
		Log.Warn("Invalid 2FA for user with id: " + userId.Hex())
		return nil, err
	}
	if !val {
		Log.Warn("Invalid 2FA for user with id: " + userId.Hex())
		return nil, errors.New("Not recognize code")
	}

	jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))
	if err != nil {
		Log.Warn("Invalid 2FA for user with id: " + userId.Hex())
		return nil, err
	}

	Log.Warn("Successful 2FA for user with id: " + userId.Hex())
	return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken, IsPrivate: user.Private}, nil
}

func (service *AuthService) Disable2fa(ctx context.Context, userId primitive.ObjectID) error {
	Log.Info("Disabling 2FA for use with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return err
	}
	user.TFAEnabled = false
	user.TFASecret = ""
	user, err = service.store.Update(ctx, userId, user)
	if err != nil {
		Log.Error("2FA was not disabled due to error, for user with id: " + userId.Hex())
		return err
	}
	Log.Info("2FA disabled for use with id: " + userId.Hex())
	return nil
}

func (service *AuthService) Enable2FA(ctx context.Context, userId primitive.ObjectID, code string) error {
	Log.Info("Enabling 2FA for use with id: " + userId.Hex())
	_, err := service.Verify2fa(ctx, userId, code)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return err
	}
	user, err := service.store.Get(ctx, userId)
	user.TFAEnabled = true
	user, err = service.store.Update(ctx, userId, user)
	if err != nil {
		Log.Error("2FA was not enabled due to error, for user with id: " + userId.Hex())
		return err
	}
	Log.Info("2FA enabled for use with id: " + userId.Hex())
	return nil
}

func (service *AuthService) GetApiToken(ctx context.Context, userId primitive.ObjectID) (string, error) {
	Log.Info("Getting API token for user with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Error("Can not get api token due to error, for user with id: " + userId.Hex())
		return "", err
	}
	Log.Info("API token got successful for user with id: " + userId.Hex())
	return user.ApiToken, nil
}

func (service *AuthService) CreateApiToken(ctx context.Context, userId primitive.ObjectID) (string, error) {
	Log.Info("Creating API token for user with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return "", err
	}

	user.ApiToken = uuid.New().String()
	_, err = service.store.Update(ctx, userId, user)
	if err != nil {
		Log.Error("Can not create api token due to error, for user with id: " + userId.Hex())
		return "", err
	}
	Log.Info("API token created successful for user with id: " + userId.Hex())
	return user.ApiToken, nil
}

func (service *AuthService) RemoveApiToken(ctx context.Context, userId primitive.ObjectID) error {
	Log.Info("Removing API token for user with id: " + userId.Hex())
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		Log.Warn("Unexciting user with id: " + userId.Hex())
		return err
	}
	user.ApiToken = ""
	_, err = service.store.Update(ctx, userId, user)
	if err != nil {
		Log.Error("Can not remove api token due to error, for user with id: " + userId.Hex())
		return err
	}
	Log.Info("API token removed successful for user with id: " + userId.Hex())
	return nil
}

func (service *AuthService) IsApiTokenValid(ctx context.Context, token string) (string, error) {
	users, err := service.store.GetAll(ctx)
	if err != nil {
		return "", err
	}

	for _, user := range users {
		if user.ApiToken == token {
			return user.Id.Hex(), nil
		}
	}
	return "", errors.New("unauthorized")
}

func (service *AuthService) CreatePasswordRecoveryRequest(ctx context.Context, username string) error {
	Log.Info("Starting password recovery for user with username: " + username)
	user, err := service.getUser(ctx, username)
	if err != nil {
		Log.Warn("Unexciting user with username: " + username)
		return err
	}

	passwordRecoveryRequest := &model.PasswordRecoveryRequest{
		UserId:  user.Id.Hex(),
		ValidTo: time.Now().Local().Add(time.Minute * time.Duration(30)),
	}

	createdRequest, err := service.store.CreatePasswordRecoveryRequest(ctx, passwordRecoveryRequest)
	if err != nil {
		Log.Error("Cannot create password recovery request for user with username: " + username)
		return err
	}

	err = SendEmailForPasswordRecovery(ctx, user, createdRequest.Id.Hex())
	if err != nil {
		Log.Error("Cannot send password recovery mail for user with username: " + username)
		return err
	}

	Log.Info("Successful send email for password recovery for user with username: " + username)
	return nil
}

func (service *AuthService) PasswordlessLoginCreate(ctx context.Context, username string) error {
	Log.Info("Create passwordless login for user with username: " + username)
	user, err := service.getUser(ctx, username)
	if err != nil {
		Log.Warn("Unexciting user with username: " + username)
		return err
	}

	id, err := service.store.CreatePasswordlessRequest(ctx, user.Id)
	if err != nil {
		Log.Error("Cannot create passwordless login for user with id: " + user.Id.Hex())
		return err
	}

	err = SendEmailForPasswordlessLogin(ctx, user, id)
	if err != nil {
		Log.Error("Cannot send email for passwordless login for user with id: " + user.Id.Hex())
		return err
	}

	Log.Info("Passwordless login created successful for user with username: " + username)
	return nil
}

func (service *AuthService) PasswordlessLogin(ctx context.Context, userId primitive.ObjectID, loginId primitive.ObjectID) (*userService.LoginResponse, error) {
	Log.Info("Starting passwordless login for user with id: " + userId.Hex())
	found, err := service.store.GetPasswordlessRequest(ctx, userId, loginId)
	if err != nil {
		Log.Error("Cannot find passwordless login for user with id: " + userId.Hex())
		return nil, err
	}
	if !found {
		Log.Error("Cannot find passwordless login for user with id: " + userId.Hex())
		return nil, errors.New("not found")
	}

	user, err := service.store.Get(ctx, userId)
	jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))

	Log.Info("Successful passwordless login for user with id: " + userId.Hex())
	return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken, IsPrivate: user.Private}, nil
}
