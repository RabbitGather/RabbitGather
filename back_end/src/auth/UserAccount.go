package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserAccount struct{}

func (*UserAccount) NewToken() (*JWTToken, error) {
	fmt.Println("Not implemented : NewToken")
	//a :=
	theConst := UserClaims{
		StandardClaims: jwt.StandardClaims{},
		UserName:       "THIS_WILL_BE_USER_NAME",
		UserID:         1,
		APIPermission:  Login | PostArticle | SearchArticle,
	}
	//token := JWTToken{}
	token, err := NewSignedToken(theConst)
	if err != nil {
		return nil, err
	}
	return token, nil
}

type UserClaims struct {
	jwt.StandardClaims
	UserName      string
	UserID        int
	APIPermission PermissionCode
}

func GetUserByClaims(userClaims *UserClaims) UserAccount {
	panic("Not implement")
}
