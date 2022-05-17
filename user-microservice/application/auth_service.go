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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"user-microservice/model"

	_ "io/ioutil"

	qr "rsc.io/qr"
)

type AuthService struct {
	store      model.UserStore
	jwtManager *token.JwtManager
}

func NewAuthService(store model.UserStore, manager *token.JwtManager) *AuthService {
	return &AuthService{
		store:      store,
		jwtManager: manager,
	}
}

func (service *AuthService) Login(ctx context.Context, in *userService.CredentialsRequest) (*userService.LoginResponse, error) {
	user, err := service.getUser(ctx, in.Credentials.Username)
	if err == nil && security.BcryptCompareHashAndPassword(user.Password, in.Credentials.Password) == nil {
		jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))
		if err != nil {
			return nil, err
		}

		if user.TFAEnabled {
			return &userService.LoginResponse{UserId: user.Id.Hex()}, nil
		}
		return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken, IsPrivate: user.Private}, nil
	}
	return nil, errors.New("wrong username or password")
}

func (service *AuthService) getUser(ctx context.Context, username string) (*model.User, error) {
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

func (service *AuthService) IsAuthenticated(ctx context.Context, jwtToken string) (model.UserRole, error) {
	ok := service.jwtManager.IsUserAuthorized(jwtToken)
	if ok != nil {
		return "", ok
	}
	userRole, err := service.jwtManager.GetRoleFromToken(jwtToken)
	if err != nil {
		return "", err
	}
	return model.UserRole(userRole), nil
}

func (service *AuthService) CheckPassword(ctx context.Context, password string, userId primitive.ObjectID) (bool, error) {
	user, err := service.store.Get(ctx, userId)
	if err == nil && security.BcryptCompareHashAndPassword(user.Password, password) != nil {
		return true, nil
	}
	return false, errors.New("wrong password")
}

func (service *AuthService) GetQR2FA(ctx context.Context, userId primitive.ObjectID) ([]byte, error) {
	user, err := service.store.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	secret := make([]byte, 10)
	_, err = rand.Read(secret)
	if err != nil {
		panic(err)
	}

	user.TFASecret = base32.StdEncoding.EncodeToString(secret)
	user.TFAEnabled = false
	service.store.Update(ctx, userId, user)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
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
	return code.PNG(), nil
}

func (service *AuthService) Verify2fa(ctx context.Context, userId primitive.ObjectID, code string) (*userService.LoginResponse, error) {
	user, err := service.store.Get(ctx, userId)

	if err != nil {
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
		return nil, err
	}
	if !val {
		return nil, errors.New("Not recognize code")
	}

	jwtToken, err := service.jwtManager.GenerateJWT(user.Id.Hex(), user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &userService.LoginResponse{UserId: user.Id.Hex(), Email: user.Email, Role: string(user.Role), Token: jwtToken, IsPrivate: user.Private}, nil
}

func (service *AuthService) Disable2fa(ctx context.Context, userId primitive.ObjectID) error {
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		return err
	}
	user.TFAEnabled = false
	user.TFASecret = ""
	user, err = service.store.Update(ctx, userId, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) Enable2FA(ctx context.Context, userId primitive.ObjectID, code string) error {
	_, err := service.Verify2fa(ctx, userId, code)
	if err != nil {
		return err
	}
	user, err := service.store.Get(ctx, userId)
	user.TFAEnabled = true
	user, err = service.store.Update(ctx, userId, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) GetApiToken(ctx context.Context, userId primitive.ObjectID) (string, error) {
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		return "", err
	}
	return user.ApiToken, nil
}

func (service *AuthService) CreateApiToken(ctx context.Context, userId primitive.ObjectID) (string, error) {
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		return "", err
	}

	user.ApiToken = uuid.New().String()
	_, err = service.store.Update(ctx, userId, user)
	if err != nil {
		return "", err
	}
	return user.ApiToken, nil
}

func (service *AuthService) RemoveApiToken(ctx context.Context, userId primitive.ObjectID) error {
	user, err := service.store.Get(ctx, userId)
	if err != nil {
		return err
	}
	user.ApiToken = ""
	_, err = service.store.Update(ctx, userId, user)
	if err != nil {
		return err
	}
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
