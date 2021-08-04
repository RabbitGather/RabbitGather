package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"rabbit_gather/src/auth/status_bitmask"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
	"time"
)

const TokenHeaderKey = "token"
const UtilityClaimKey = "utility_claim"

var log = logger.NewLoggerWrapper("token")
var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

var Issuer string
var SignMethod jwt.SigningMethod

const RS256 = "RS256"

func init() {
	type Config struct {
		JwtPrivateKeyFile string `json:"jwt_private_key_file"`
		JwtPublicKeyFile  string `json:"jwt_public_key_file"`
		//ExpiresAtTimeDuration int    `json:"expires_at_time"`
		//NotBeforeTimeDuration int    `json:"not_before_time"`
		Issuer     string `json:"issuer"`
		SignMethod string `json:"sign_method"`
	}

	var config Config
	err := util.ParseJsonConfic(&config, "config/JWT.config.json")
	if err != nil {
		panic(err.Error())
	}

	switch config.SignMethod {
	case RS256:
		SignMethod = jwt.SigningMethodRS256
	default:
		panic("Not supported SignMethod: " + config.SignMethod)
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
	//ExpiresAtTimeDuration = time.Duration(config.ExpiresAtTimeDuration)
	//NotBeforeTimeDuration = time.Duration(config.NotBeforeTimeDuration)
	privateKey = getPrivateKey(config.JwtPrivateKeyFile)
	publicKey = getPublicKey(config.JwtPublicKeyFile)
	Issuer = config.Issuer
}

// Get a new Standard Claims with defult setting.
func NewStandardClaims() claims.StandardClaims {
	nowTime := time.Now()
	return claims.StandardClaims{
		Audience:  "",
		ExpiresAt: nowTime.Add(claims.DefaultExpiresAtTimeDuration).Unix(),
		Id:        util.Snowflake().String(),
		IssuedAt:  nowTime.Unix(),
		Issuer:    Issuer,
		NotBefore: nowTime.Add(claims.DefaultNotBeforeTimeDuration).Unix(),
		Subject:   "",
	}
}

// Parse the jwt token from string and fill in claims.
// Return error when the input JWT token not Vaild
func ParseToken(rawTokenString string) (*claims.UtilityClaims, error) {
	var mapClaims jwt.MapClaims
	utilityClaims := claims.UtilityClaims{}
	//a:=jwt.MapClaims(utilityClaims)
	token, err := jwt.ParseWithClaims(rawTokenString, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		e := checkTokenWhenParse(token)
		if e != nil {
			log.DEBUG.Println(e.Error())
		}
		return publicKey, e
	})
	if err != nil {
		return nil, err
	}
	err = utilityClaims.ParseMapClaims(mapClaims)
	if err != nil {
		log.ERROR.Println("Error when parse token: ", err.Error())
		return nil, err
	}
	if !token.Valid {
		err := fmt.Errorf("this token is not valid")
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	return &utilityClaims, nil
}

// Check if the token format right.
func checkTokenWhenParse(token *jwt.Token) error {
	if token.Method != SignMethod {
		return errors.New("wrong SignMethod")
	}
	return nil
}

// Sign a new token.
func SignToken(utilityClaims *claims.UtilityClaims) (string, error) {
	// put default StandardClaims in if not set
	if _, exist := (*utilityClaims)[claims.StandardClaimsName]; !exist {
		utilityClaims.SetSubClaims(claims.StandardClaimsName, NewStandardClaims())
	}

	signedString, err := jwt.NewWithClaims(SignMethod, utilityClaims).SignedString(privateKey)
	if err != nil {
		log.ERROR.Println("NewSignedToken Error")
		return "", err
	}
	return signedString, nil
}

func UpdateStatus(existToken string, code status_bitmask.StatusBitmask) (string, error) {
	if existToken == "" {
		return "", errors.New("input \"\" as token")
	} else {
		uc, err := ParseToken(existToken)
		if err != nil {
			log.DEBUG.Println("fail with ParseToken: ", err.Error())
			return "", err
		}

		var statusClaims claims.StatusClaims
		err = uc.GetSubClaims(&statusClaims)
		if err != nil {
			if err == claims.NoSuchClaimsError {
				log.DEBUG.Println("status claims not exist: ", err.Error())
				return "", err
			} else if err == claims.UnknownClaimsError {
				log.ERROR.Println("UnknownClaimsError: ", err.Error())
				return "", err
			} else {
				panic(err.Error())
			}

		}
		statusClaims.AppendBitMask(code)

		token, err := SignToken(uc)
		if err != nil {
			log.DEBUG.Println("error when Sign token: ", err.Error())
			return "", err
		}
		return token, nil
	}
}

//
//type JWTToken struct {
//	jwt.Token
//	signedString string
//}
//
//func (t *JWTToken) GetSignedString() string {
//	if t.signedString == "" {
//		panic("signedString is empty")
//	}
//	return t.signedString
//}
//
//func (t *JWTToken) AppendBitMask(code StatusBitmask) (*JWTToken, error) {
//	permissionClaims, ok := t.Claims.(*UtilityClaims)
//	if !ok {
//		return nil, errors.New("The Claims is not a UtilityClaims")
//	}
//	if BitMaskCheck(permissionClaims.PermissionBitmask, code) {
//		return t, nil
//	} else {
//		permissionClaims.PermissionBitmask = permissionClaims.PermissionBitmask | code
//	}
//	newToken, err := NewSignedToken(permissionClaims)
//	return newToken, err
//}

// ParseToken Parse the signed token string into claims
//func ParseToken(signedTokenString string, claims jwt.Claims) (*JWTToken, error) {
//if signedTokenString == "" {
//	return nil, errors.New("the token string is \"\"")
//}
//token, err := jwt.ParseWithClaims(signedTokenString, claims, func(token *jwt.Token) (interface{}, error) {
//	e := checkTokenWhenParse(token)
//	return publicKey, e
//})
//if err != nil {
//	return nil, err
//}
//jwtToken := &JWTToken{
//	Token:        *token,
//	signedString: signedTokenString,
//}
//return jwtToken, nil
//}

// NewSignedToken Create new Signed token
//func NewSignedToken(claims jwt.Claims) (*JWTToken, error) {
//	token := jwt.NewWithClaims(SignMethod, claims)
//	signedString, err := token.SignedString(privateKey)
//	if err != nil {
//		log.DEBUG.Println("NewSignedToken Error: ", err.Error())
//		return nil, err
//	}
//	token, err = jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
//		return publicKey, nil
//	})
//	if err != nil {
//		log.DEBUG.Println("ParseWithClaims Error: ", err.Error())
//		return nil, err
//	}
//	jwtToken := &JWTToken{
//		Token:        *token,
//		signedString: signedString,
//	}
//	return jwtToken, nil
//}
