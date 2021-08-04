package claims

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type UtilityClaims map[string]jwt.Claims

var NoSuchClaimsError = errors.New("no such claim")
var UnknownClaimsError = errors.New("unknown claims type")

func (u UtilityClaims) GetSubClaims(subClaims interface{}) error {
	var exist = false
	var claims interface{}
	switch t := subClaims.(type) {
	case *StatusClaims:
		claims, exist = u[StatusClaimsName]
		if exist {
			*t = claims.(StatusClaims)
		}
	case *StandardClaims:
		claims, exist = u[StandardClaimsName]
		if exist {
			*t = claims.(StandardClaims)
		}
	case *UserClaims:
		claims, exist = u[UserClaimsName]
		if exist {
			*t = claims.(UserClaims)
		}
	default:
		return UnknownClaimsError
	}
	if !exist {
		return NoSuchClaimsError
	}
	return nil
}

func (u UtilityClaims) SetSubClaims(name string, claim jwt.Claims) {
	u[name] = claim
}
func (u UtilityClaims) RemoveSubClaims(name string) {
	delete(u, name)
}

func (u UtilityClaims) Valid() error {
	var err error
	for name, claim := range u {
		if claim == nil {
			return errors.New("nil sub claims found in UtilityClaims")
		}
		//cm,ok :=claim.(map[string]interface{})//jwt.MapClaims(claim.(map[string]interface{}))
		//if !ok {
		//	err := fmt.Errorf("found non jwt.Claims sub claim in UtilityClaims: %s",fmt.Sprint(claim))
		//	log.ERROR.Println(err.Error())
		//	//debug.PrintStack()
		//	return err
		//}
		//c := jwt.MapClaims(cm)
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

// turm map claims into UtilityClaims
func (u UtilityClaims) ParseMapClaims(claims jwt.MapClaims) error {
	for claimType, claimValue := range claims {
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
		case StandardClaimsName:
			var theClaim StandardClaims
			err := decode(&theClaim, claimValue)
			if err != nil {
				return err
			} else {
				u[claimType] = theClaim
			}
		case StatusClaimsName:
			var theClaim StatusClaims
			err := decode(&theClaim, claimValue)
			if err != nil {
				return err
			} else {
				u[claimType] = theClaim
			}
		case UserClaimsName:
			var theClaim UserClaims
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

//func (u UtilityClaims) ParseMapClaims(mapClaims jwt.MapClaims) error {
//	//u["sodas"] = StatusClaims{StatusBitmask: NoStatus}
//	//log.TempLog().Println(u)
//	rawMap := map[string]interface{}(mapClaims)
//	for s, i := range rawMap {
//		switch s {
//		case string(StatusClaimsName):
//			var sc StatusClaims
//			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
//				Metadata: nil,
//				Result:   &sc,
//				TagName:  "json",
//			})
//			if err != nil {
//				return err
//			}
//			err = decoder.Decode(i)
//			if err != nil {
//				return err
//			}
//			u[s] = sc
//		case string(UserClaimsName):
//			var sc UserClaims
//			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
//				Metadata: nil,
//				Result:   &sc,
//				TagName:  "json",
//			})
//			if err != nil {
//				return err
//			}
//			err = decoder.Decode(i)
//			if err != nil {
//				return err
//			}
//			u[s] = sc
//
//		case string(StandardClaimsName):
//			var sc StandardClaims
//			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
//				Metadata: nil,
//				Result:   &sc,
//				TagName:  "json",
//			})
//			if err != nil {
//				return err
//			}
//			err = decoder.Decode(i)
//			if err != nil {
//				return err
//			}
//			u[s] = sc
//		}
//	}
//	return nil
//}
