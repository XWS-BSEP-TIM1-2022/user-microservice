package application

import (
	"context"
	"errors"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/security"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/smtp"
	"regexp"
	"strings"
	"user-microservice/application/smtp_login"
	"user-microservice/model"
	"user-microservice/startup/config"
)

type UserService struct {
	store  model.UserStore
	config *config.Config
}

func NewUserService(store model.UserStore, config *config.Config) *UserService {
	return &UserService{
		store:  store,
		config: config,
	}
}

func (service *UserService) Get(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return service.store.Get(ctx, id)
}

func (service *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return service.store.GetAll(ctx)
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := service.store.GetByUsername(ctx, user.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	if user.Email != "" {
		_, err = service.store.GetByEmail(ctx, user.Email)
		if err == nil {
			return nil, errors.New("email already exists")
		}
	}

	err = service.IsPasswordOk(user.Password)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := security.BcryptGenerateFromPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	user.Confirmed = false
	user.ConfirmationId = uuid.New().String()
	err = sendConfirmationMail(ctx, user)
	if err != nil {
		return nil, err
	}

	return service.store.Create(ctx, user)
}

func (service *UserService) IsPasswordOk(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be atleast 8 characters")
	}

	match, _ := regexp.MatchString("[0-9]", password)
	if !match {
		return errors.New("Password must contain atleast 1 number")
	}

	match, _ = regexp.MatchString("[A-Z]", password)
	if !match {
		return errors.New("Password must contain atleast 1 upper case")
	}

	match, _ = regexp.MatchString("[a-z]", password)
	if !match {
		return errors.New("Password must contain atleast 1 lower case")
	}

	match, _ = regexp.MatchString("[.,<>/?|';:!@#$%^&*()_+=-]", password)
	if !match {
		return errors.New("Password must contain atleast 1 special characher")
	}

	/*for _, commonPassword := range service.config.CommonPasswords {
		if strings.Contains(commonPassword, password) || strings.Contains(password, commonPassword) {
			return errors.New("Password must not be a common password or containts common. (" + commonPassword + ")")
		}
	}*/
	err := service.CheckIsPasswordInCommonPasswords(password)
	if err != nil {
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
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Role = existUser.Role
	user.Username = existUser.Username
	user.Password = existUser.Password
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) UpdatePassword(ctx context.Context, userId primitive.ObjectID, user *model.User) (*model.User, error) {
	existUser, err := service.store.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Role = existUser.Role
	return service.store.Update(ctx, userId, user)
}

func (service *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return service.store.Delete(ctx, id)
}

func (service *UserService) DeleteAll(ctx context.Context) {
	service.store.DeleteAll(ctx)
}

func (service *UserService) Search(ctx context.Context, searchParam string) ([]*model.User, error) {
	users, err := service.store.GetAllWithoutAdmins(ctx)

	if err != nil {
		return nil, err
	}
	searchParam = strings.ToLower(searchParam)
	var retVal []*model.User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Username), searchParam) || strings.Contains(strings.ToLower(user.Name), searchParam) || strings.Contains(strings.ToLower(user.Surname), searchParam) || strings.Contains(strings.ToLower(user.Email), searchParam) {
			retVal = append(retVal, user)
		}
	}
	return retVal, nil
}

func (service *UserService) IsUserPrivate(ctx context.Context, id primitive.ObjectID) (bool, error) {
	user, err := service.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return user.Private, nil
}

func sendConfirmationMail(ctx context.Context, user *model.User) error {
	from := "xwstim1@outlook.com"
	password := "XWS.tim1"

	toEmailAddress := user.Email
	to := []string{toEmailAddress}

	host := "smtp-mail.outlook.com"
	port := "587"
	address := host + ":" + port
	url := "https://localhost:8090/auth/verify/" + user.ConfirmationId

	subject := "Subject: Verify your account on dislinkt\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "\nPozdrav " + user.Name + ",<br>" + "Da biste verifikovali svoj nalog, posetite sledeću stranicu:<br>" + "<h1><a href=" + url + " target=\"_self\">VERIFIKUJ</a></h1> " + "Hvala,<br>" + "Dislinkt."
	message := []byte(subject + mime + body)

	//auth := smtp.PlainAuth("", from, password, host)
	auth := smtp_login.LoginAuth(from, password)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}
