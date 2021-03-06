package application

import (
	"context"
	"net/smtp"
	"user-microservice/application/smtp_login"
	"user-microservice/model"
	cfg "user-microservice/startup/config"
)

func SendConfirmationMail(ctx context.Context, user *model.User) error {
	config := cfg.NewConfig()
	from := config.Email
	password := config.EmailPassword

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

func SendEmailForPasswordRecovery(ctx context.Context, user *model.User, passwordRecoveryId string) error {
	Log.Info("Starting send email for password recovery for user with id: " + user.Id.Hex())
	config := cfg.NewConfig()
	from := config.Email
	password := config.EmailPassword

	toEmailAddress := user.Email
	to := []string{toEmailAddress}

	host := "smtp-mail.outlook.com"
	port := "587"
	address := host + ":" + port
	url := "https://localhost:4200/create-new-password/" + passwordRecoveryId

	subject := "Subject: Reset your password\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "\nPozdrav " + user.Name + ",<br>" + "Da biste restartovali svoju lozinku, posetite sledeću stranicu:<br>" + "<h1><a href=" + url + " target=\"_self\">PROMENI LOZINKU</a></h1> " + "Hvala,<br>" + "Dislinkt."
	message := []byte(subject + mime + body)

	//auth := smtp.PlainAuth("", from, password, host)
	auth := smtp_login.LoginAuth(from, password)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}

	Log.Info("Email send successfully for password recovery for user with id: " + user.Id.Hex())
	return nil
}

func SendEmailForPasswordlessLogin(ctx context.Context, user *model.User, passwordlessId string) error {
	Log.Info("Starting send email for passwordless login for user with id: " + user.Id.Hex())
	config := cfg.NewConfig()
	from := config.Email
	password := config.EmailPassword

	toEmailAddress := user.Email
	to := []string{toEmailAddress}

	host := "smtp-mail.outlook.com"
	port := "587"
	address := host + ":" + port
	url := "https://localhost:4200/login/" + user.Id.Hex() + "/" + passwordlessId

	subject := "Passwordless LOGIN:\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "\nPozdrav " + user.Name + ",<br>" + "Kliknite na link da biste se logovali:<br>" + "<h1><a href=" + url + " target=\"_self\">ULOGUJ SE</a></h1> " + "Pozdrav,<br>" + "Dislinkt."
	message := []byte(subject + mime + body)

	auth := smtp_login.LoginAuth(from, password)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}

	Log.Info("Email send successfully for passwordless login for user with id: " + user.Id.Hex())
	return nil
}
