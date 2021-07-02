package auth

import "fmt"

type AccountManagement struct {

}
type UserAccount struct {}

type JWTToken struct {

}
func (*JWTToken)ToString()string{
	fmt.Println("Not implemented : ToString")

	return "THIS_WILL_BE_JWT_ACCESS_TOKEN"
}
func (*UserAccount)NewToken()*JWTToken{
	fmt.Println("Not implemented : NewToken")

	return &JWTToken{}
}

func (m AccountManagement) GetUserByName(username string) *UserAccount {
	fmt.Println("Not implemented : GetUserByName")
	return &UserAccount{}
}

func (m AccountManagement) CheckUserAndPassword(username ,password string) error {
	fmt.Println("Not implemented : checkUserAndPassword")
return nil
}