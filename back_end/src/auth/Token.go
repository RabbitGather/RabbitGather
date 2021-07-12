package auth

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"rabbit_gather/util"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	type Config struct {
		JwtPrivateKeyFile string `json:"jwt_private_key_file"`
		JwtPublicKeyFile  string `json:"jwt_public_key_file"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/JWT.config.json")
	if err != nil {
		panic(err.Error())
	}
	getPrivateKey := func(theJwtPrivatekeyfile string) (pk *rsa.PrivateKey) {
		privateKeyBytes, err := ioutil.ReadFile(theJwtPrivatekeyfile)
		if err != nil {
			panic(err.Error())
		}
		pk, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	getPublicKey := func(theTokenPublicKeyFile string) (pk *rsa.PublicKey) {
		publicKeyBytes, err := ioutil.ReadFile(theTokenPublicKeyFile)
		if err != nil {
			panic(err.Error())
		}
		pk, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	privateKey = getPrivateKey(config.JwtPrivateKeyFile)
	publicKey = getPublicKey(config.JwtPublicKeyFile)
}

type JWTToken struct {
	jwt.Token
	signedString string
}

func (t *JWTToken) GetSignedString() string {
	if t.signedString != "" {
		panic("signedString is empty")
	}
	return t.signedString
}

var JWTTokenSigningMethod = jwt.SigningMethodRS256

// ParseToken Parse the signed token string into claims
func ParseToken(signedTokenString string, claims jwt.Claims) (*JWTToken, error) {
	token, err := jwt.ParseWithClaims(signedTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		e := checkTokenWhenParse(token)
		return publicKey, e
	})
	if err != nil {
		return nil, err
	}
	jwtToken := &JWTToken{
		Token:        *token,
		signedString: signedTokenString,
	}
	return jwtToken, nil
}

// NewSignedToken Create new Signed token
func NewSignedToken(claims jwt.Claims) (*JWTToken, error) {
	token := jwt.NewWithClaims(JWTTokenSigningMethod, claims)
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	token, err = jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		//e := checkTokenWhenParse(token)
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	jwtToken := &JWTToken{
		Token:        *token,
		signedString: signedString,
	}
	return jwtToken, nil
}

func checkTokenWhenParse(token *jwt.Token) error {
	if token.Method != JWTTokenSigningMethod {
		return errors.New("token signed method wrong")
	}
	return nil
}
