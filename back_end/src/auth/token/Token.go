package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

// The TokenKey is the key that use in the return and request header token field.
const TokenKey = "token"

var log = logger.NewLoggerWrapper("token")

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

//var Issuer string
var SignMethod jwt.SigningMethod

const RS256 = "RS256"

func init() {
	type Config struct {
		JwtPrivateKeyFile string `json:"jwt_private_key_file"`
		JwtPublicKeyFile  string `json:"jwt_public_key_file"`
		SignMethod        string `json:"sign_method"`
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

	privateKey = getPrivateKey(config.JwtPrivateKeyFile)
	publicKey = getPublicKey(config.JwtPublicKeyFile)
}

// ParseToken Parse the jwt token from string and fill in UtilityClaim.
// Return error when the input JWT token not Valid
func ParseToken(rawTokenString string) (*claims.UtilityClaim, error) {
	var mapClaims jwt.MapClaims
	utilityClaims := claims.UtilityClaim{}
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
	err = parseMapClaims(utilityClaims, mapClaims)
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

// parse the jwt.MapClaims and fill sub claims into UtilityClaim
func parseMapClaims(u claims.UtilityClaim, cl jwt.MapClaims) error {
	for claimType, claimValue := range cl {
		decode := func(res, target interface{}) error {
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   res,
				TagName:  "json",
			})
			if err != nil {
				return err
			}
			err = decoder.Decode(target)
			if err != nil {
				return err
			}
			return nil
		}
		switch claimType {
		case claims.StandardClaimsName:
			var theClaim claims.StandardClaim
			err := decode(&theClaim, claimValue)
			if err != nil {
				return err
			} else {
				u[claimType] = theClaim
			}
		case claims.StatusClaimsName:
			var theClaim claims.StatusClaim
			err := decode(&theClaim, claimValue)
			if err != nil {
				return err
			} else {
				u[claimType] = theClaim
			}
		case claims.UserClaimsName:
			var theClaim claims.UserClaim
			err := decode(&theClaim, claimValue)
			if err != nil {
				return err
			} else {
				u[claimType] = theClaim
			}
		}
	}
	return nil
}

// Check if the token format right.
func checkTokenWhenParse(token *jwt.Token) error {
	if token.Method != SignMethod {
		return errors.New("wrong SignMethod")
	}
	return nil
}

// SignToken Sign a new JWT token.
func SignToken(utilityClaims *claims.UtilityClaim) (string, error) {
	// put default StandardClaim in if not set
	if _, exist := (*utilityClaims)[claims.StandardClaimsName]; !exist {
		utilityClaims.SetSubClaims(claims.StandardClaimsName, claims.NewStandardClaims())
	}

	signedString, err := jwt.NewWithClaims(SignMethod, utilityClaims).SignedString(privateKey)
	if err != nil {
		log.ERROR.Println("NewSignedToken Error")
		return "", err
	}
	return signedString, nil
}

//
//func UpdateStatus(existToken string, code bitmask.StatusBitmask) (string, error) {
//	if existToken == "" {
//		return "", errors.New("input \"\" as token")
//	} else {
//		uc, err := ParseToken(existToken)
//		if err != nil {
//			log.DEBUG.Println("fail with ParseToken: ", err.Error())
//			return "", err
//		}
//
//		var statusClaims claims.StatusClaim
//		err = uc.GetSubClaims(&statusClaims)
//		if err != nil {
//			if err == claims.NoSuchClaimsError {
//				log.DEBUG.Println("status claims not exist: ", err.Error())
//				return "", err
//			} else if err == claims.UnknownClaimsError {
//				log.ERROR.Println("UnknownClaimsError: ", err.Error())
//				return "", err
//			} else {
//				panic(err.Error())
//			}
//
//		}
//		statusClaims.AppendBitMask(code)
//
//		token, err := SignToken(uc)
//		if err != nil {
//			log.DEBUG.Println("error when Sign token: ", err.Error())
//			return "", err
//		}
//		return token, nil
//	}
//}
