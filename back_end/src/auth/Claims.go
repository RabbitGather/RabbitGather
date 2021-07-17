package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type ClaimsName string

const (
	StatusClaimsName   ClaimsName = "status_claims"
	UserClaimsName     ClaimsName = "user_claims"
	StandardClaimsName ClaimsName = "standard_claims"
)

type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

func (s StandardClaims) Valid() error {
	panic("implement me")
}

type StatusClaims struct {
	StatusBitmask StatusBitmask `json:"status_bitmask"`
}

func (s StatusClaims) Valid() error {
	panic("implement me")
}

func (s StatusClaims) AppendBitMask(code StatusBitmask) {

	//newToken, err := token.AppendBitMask(auth.WaitVerificationCode)
}

type UserClaims struct {
	UserName string `json:"user_name"`
	UserID   uint32 `json:"user_id"`
}

func (u UserClaims) Valid() error {
	panic("implement me")
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

func (u UtilityClaims) ParseMapClaims(mapClaims *jwt.MapClaims) error {
	//u["sodas"] = StatusClaims{StatusBitmask: NoStatus}
	//log.TempLog().Println(u)
	rawMap := map[string]interface{}(*mapClaims)
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
	log.TempLog().Println(u)

	return nil
}

func (u UtilityClaims) MappingClaim(inputClaimPointer interface{}) error {
	switch inputClaimPointer.(type) {
	case *StatusClaims:
		claim, exist := u.GetSubClaims(StatusClaimsName)
		if !exist {
			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
		}
		*(inputClaimPointer.(*StatusClaims)) = claim.(StatusClaims)
	case *UserClaims:
		claim, exist := u.GetSubClaims(UserClaimsName)
		if !exist {
			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
		}
		*(inputClaimPointer.(*UserClaims)) = claim.(UserClaims)
	case *StandardClaims:
		claim, exist := u.GetSubClaims(StandardClaimsName)
		if !exist {
			return fmt.Errorf("the claim:%s is not in this UtilityClaims", fmt.Sprint(inputClaimPointer))
		}
		*(inputClaimPointer.(*StandardClaims)) = claim.(StandardClaims)
	default:
		panic("not supported Claim type")
	}
	return nil
}
