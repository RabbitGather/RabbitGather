package claims

import (
	"crypto/subtle"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"rabbit_gather/src/auth/status_bitmask"
	//"rabbit_gather/src/auth"
	"rabbit_gather/src/logger"
	"time"
)

type ClaimsName string

var log = logger.NewLoggerWrapper("auth.Claims")

var DefaultExpiresAtTimeDuration = time.Hour * 12
var DefaultNotBeforeTimeDuration = time.Second * 2

const StandardClaimsName ClaimsName = "standard_claims"

type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

func (c StandardClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

// Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyAudience(cmp string, req bool) bool {
	return verifyAud(c.Audience, cmp, req)
}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

// Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(c.IssuedAt, cmp, req)
}

// Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(c.Issuer, cmp, req)
}

// Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(c.NotBefore, cmp, req)
}

// ----- helpers

func verifyAud(aud string, cmp string, required bool) bool {
	if aud == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(aud), []byte(cmp)) != 0 {
		return true
	} else {
		return false
	}
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
		return true
	} else {
		return false
	}
}

func verifyNbf(nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	return now >= nbf
}

const StatusClaimsName ClaimsName = "status_claims"

type StatusClaims struct {
	StatusBitmask status_bitmask.StatusBitmask `json:"status_bitmask"`
}

func (s StatusClaims) Valid() error {
	return nil
}

func (s StatusClaims) AppendBitMask(code status_bitmask.StatusBitmask) {

}

const UserClaimsName ClaimsName = "user_claims"

type UserClaims struct {
	UserName string `json:"user_name"`
	UserID   uint32 `json:"user_id"`
}

func (u UserClaims) Valid() error {
	return nil
}

type UtilityClaims map[ClaimsName]jwt.Claims

func (u UtilityClaims) GetSubClaims(name ClaimsName) (jwt.Claims, bool) {
	claims, exist := u[name]
	return claims, exist
}

func (u UtilityClaims) SetSubClaims(name ClaimsName, claim jwt.Claims) {
	u[name] = claim
	//switch claim.(type) {
	//case *StatusClaims:
	//case *UserClaims:
	//case *StandardClaims:
	//default:
	//	panic("not supported Claim type")
	//}
}

func (u UtilityClaims) Valid() error {
	var err error
	for name, claim := range u {
		if claim == nil {
			continue
		}
		e := claim.Valid()
		if e != nil {
			if err == nil {
				err = fmt.Errorf("%s: %w", name, e)
				continue
			}
			err = fmt.Errorf("%s ->%s: %w", err.Error(), name, e)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (u UtilityClaims) ParseMapClaims(mapClaims jwt.MapClaims) error {
	//u["sodas"] = StatusClaims{StatusBitmask: NoStatus}
	//log.TempLog().Println(u)
	rawMap := map[string]interface{}(mapClaims)
	for s, i := range rawMap {
		switch s {
		case string(StatusClaimsName):
			var sc StatusClaims
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &sc,
				TagName:  "json",
			})
			if err != nil {
				return err
			}
			err = decoder.Decode(i)
			if err != nil {
				return err
			}
			u[ClaimsName(s)] = sc
		case string(UserClaimsName):
			var sc UserClaims
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &sc,
				TagName:  "json",
			})
			if err != nil {
				return err
			}
			err = decoder.Decode(i)
			if err != nil {
				return err
			}
			u[ClaimsName(s)] = sc

		case string(StandardClaimsName):
			var sc StandardClaims
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &sc,
				TagName:  "json",
			})
			if err != nil {
				return err
			}
			err = decoder.Decode(i)
			if err != nil {
				return err
			}
			u[ClaimsName(s)] = sc
		}
	}
	return nil
}

//
//func (u UtilityClaims) MappingClaim(inputClaimPointer interface{}) error {
//	log.TempLog().Println(u["status_claims"])
//	switch inputClaimPointer.(type) {
//	case *StatusClaims:
//		claim, exist := u.GetSubClaims(StatusClaimsName)
//		if !exist {
//			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
//		}
//		a := claim.(StatusClaims)
//		//inputClaimPointer = claim.(*StatusClaims)
//		*(inputClaimPointer.(*StatusClaims)) = a
//		log.TempLog().Println(inputClaimPointer)
//	case *UserClaims:
//		claim, exist := u.GetSubClaims(UserClaimsName)
//		if !exist {
//			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
//		}
//		a := claim.(UserClaims)
//		*(inputClaimPointer.(*UserClaims)) = a
//	case *StandardClaims:
//		claim, exist := u.GetSubClaims(StandardClaimsName)
//		if !exist {
//			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
//		}
//		a := claim.(StandardClaims)
//		*(inputClaimPointer.(*StandardClaims)) = a
//	default:
//		panic("not supported Claim type")
//	}
//	return nil
//}
