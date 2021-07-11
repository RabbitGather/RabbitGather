package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserAccount struct{}

func (*UserAccount) NewToken() (*JWTToken, error) {
	fmt.Println("Not implemented : NewToken")
	//a :=
	theConst := struct {
		jwt.StandardClaims
		Text string `json:"text"`
	}{Text: "THIS_WILL_BE_JWT_ACCESS_TOKEN"}
	//token := JWTToken{}
	token, err := NewSignedToken(theConst)
	if err != nil {
		return nil, err
	}
	return token, nil
}
