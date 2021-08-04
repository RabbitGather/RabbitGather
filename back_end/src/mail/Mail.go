package mail

import (
	"errors"
	"fmt"

	//"fmt"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

type GmailSender struct {
	userMail     string
	userPassword string
	auth         smtp.Auth
}

func (s *GmailSender) SendMail(msg, subject, to string) error {
	m := fmt.Sprintf("Subject: %s\n%s", subject, msg)
	err := smtp.SendMail("smtp.gmail.com:587", s.auth, s.userMail, []string{to}, []byte(m))
	if err != nil {
		return err
	}
	return nil
}
func NewGmailSender(userMail, username, password string) *GmailSender {
	return &GmailSender{
		userMail: userMail,
		auth:     LoginAuth(username, password),
	}
}
