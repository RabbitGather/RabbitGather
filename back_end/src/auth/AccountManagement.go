package auth

import (
	"fmt"
)

type AccountManagement struct {
}

func (m AccountManagement) GetUserByName(username string) *UserAccount {
	fmt.Println("Not implemented : GetUserByName")
	return &UserAccount{}
}

func (m AccountManagement) CheckUserAndPassword(username, password string) error {
	fmt.Println("Not implemented : checkUserAndPassword")
	return nil
}
