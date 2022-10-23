package claims

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rabbit_gather/util"
)

// The UtilityClaim is a shells that carry all Claims
type UtilityClaim map[string]jwt.Claims

// The NoSuchClaimsError will be thrown when the specific SubClaims is not exist in UtilityClaim
var NoSuchClaimsError = errors.New("no such claim")

// The UnknownClaimsError will be thrown when the specific SubClaims type is not supported
var UnknownClaimsError = errors.New("unknown claims type")

// The GetSubClaims receiving a Claim pointer, it will be pointed to the
// corresponding Claim struct if it exists in UtilityClaim
// Throw NoSuchClaimsError when the specific Claim not exist.
// Throw UnknownClaimsError when the giving Claim pointer is not supported.
func (u UtilityClaim) GetSubClaims(subClaims interface{}) error {
	var exist = false
	var claims interface{}
	switch t := subClaims.(type) {
	case *StatusClaim:
		claims, exist = u[StatusClaimsName]
		if exist {
			*t = claims.(StatusClaim)
		}
	case *StandardClaim:
		claims, exist = u[StandardClaimsName]
		if exist {
			*t = claims.(StandardClaim)
		}
	case *UserClaim:
		claims, exist = u[UserClaimsName]
		if exist {
			*t = claims.(UserClaim)
		}
	default:
		return UnknownClaimsError
	}
	if !exist {
		return NoSuchClaimsError
	}
	return nil
}

// SetSubClaims put a new sub claim in UtilityClaim
func (u UtilityClaim) SetSubClaims(name string, claim jwt.Claims) {
	u[name] = claim
}

// RemoveSubClaims remove the sub claim in UtilityClaim
func (u UtilityClaim) RemoveSubClaims(name string) {
	delete(u, name)
}

// The Valid in UtilityClaim will call all the sub claim's Valid() and concatenated all errors if existed.
func (u UtilityClaim) Valid() error {
	wrapper := util.ErrorWrapper{Err: nil}
	for name, claim := range u {
		if claim == nil {
			panic("nil sub claims found in UtilityClaim")
		}
		e := claim.Valid()
		if e != nil {
			wrapper.Wrap(fmt.Errorf("%s: %s", name, e.Error()))
		}
	}
	return wrapper.GetError()
}
