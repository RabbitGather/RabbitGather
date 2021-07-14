package auth

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
	"time"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	type Config struct {
		JwtPrivateKeyFile     string `json:"jwt_private_key_file"`
		JwtPublicKeyFile      string `json:"jwt_public_key_file"`
		ExpiresAtTimeDuration int    `json:"expires_at_time_sec"`
		NotBeforeTimeDuration int    `json:"not_before_time_sec"`
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
	ExpiresAtTimeDuration = time.Duration(config.ExpiresAtTimeDuration)
	NotBeforeTimeDuration = time.Duration(config.NotBeforeTimeDuration)
	privateKey = getPrivateKey(config.JwtPrivateKeyFile)
	publicKey = getPublicKey(config.JwtPublicKeyFile)
}

var ExpiresAtTimeDuration time.Duration
var NotBeforeTimeDuration time.Duration

func NewStandardClaims() *jwt.StandardClaims {
	nowTime := time.Now()
	return &jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: nowTime.Add(ExpiresAtTimeDuration).Unix(),
		Id:        util.Snowflake().String(),
		IssuedAt:  nowTime.Unix(),
		Issuer:    "meowalien.com",
		NotBefore: nowTime.Add(NotBeforeTimeDuration).Unix(),
		Subject:   "",
	}
}

const TokenHeaderKey = "token"

type JWTToken struct {
	jwt.Token
	signedString string
}

func (t *JWTToken) GetSignedString() string {
	if t.signedString == "" {
		panic("signedString is empty")
	}
	return t.signedString
}

func (t *JWTToken) AppendBitMask(code APIPermissionBitmask) (*JWTToken, error) {
	permissionClaims, ok := t.Claims.(*PermissionClaims)
	if !ok {
		return nil, errors.New("The Claims is not a PermissionClaims")
	}
	if BitMaskCheck(permissionClaims.PermissionBitmask, code) {
		return t, nil
	} else {
		permissionClaims.PermissionBitmask = permissionClaims.PermissionBitmask | code
	}
	newToken, err := NewSignedToken(permissionClaims)
	return newToken, err
}

var JWTTokenSigningMethod = jwt.SigningMethodRS256

// ParseToken Parse the signed token string into claims
func ParseToken(signedTokenString string, claims jwt.Claims) (*JWTToken, error) {
	if signedTokenString == "" {
		return nil, errors.New("the token string is \"\"")
	}
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

var log = logger.NewLoggerWrapper("auth.Token")

// NewSignedToken Create new Signed token
func NewSignedToken(claims jwt.Claims) (*JWTToken, error) {
	token := jwt.NewWithClaims(JWTTokenSigningMethod, claims)
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		log.DEBUG.Println("NewSignedToken Error: ", err.Error())
		return nil, err
	}
	token, err = jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		log.DEBUG.Println("ParseWithClaims Error: ", err.Error())
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
