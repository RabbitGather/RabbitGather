package claims

import (
	"crypto/subtle"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rabbit_gather/util"
	"time"
)

// The StandardClaim implement the RFC7519 - Section 4 standard JWT Claims
// https://datatracker.ietf.org/doc/html/rfc7519#section-4
type StandardClaim struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

// The StandardClaimsName is the string name represented of StandardClaim.
// Will be use as the key in UtilityClaim.
const StandardClaimsName = "standard_claims"

// DefaultExpiresAtTimeDuration is the Default ExpiresAt(exp) value
var DefaultExpiresAtTimeDuration = time.Hour * 12

// DefaultNotBeforeTimeDuration is the Default NotBefore(nbf) value
var DefaultNotBeforeTimeDuration = time.Second * 2

var DefaultIssuer = "meowalien.com"

// Get a new Standard Claims with defult setting.
func NewStandardClaims() StandardClaim {
	nowTime := time.Now()
	return StandardClaim{
		Audience:  "",
		ExpiresAt: nowTime.Add(DefaultExpiresAtTimeDuration).Unix(),
		Id:        util.Snowflake().String(),
		IssuedAt:  nowTime.Unix(),
		Issuer:    DefaultIssuer,
		NotBefore: nowTime.Add(DefaultNotBeforeTimeDuration).Unix(),
		Subject:   "",
	}
}

func (c StandardClaim) Valid() error {
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

// VerifyAudience Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaim) VerifyAudience(cmp string, req bool) bool {
	return verifyAud(c.Audience, cmp, req)
}

// VerifyExpiresAt Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaim) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

// VerifyIssuedAt Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaim) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(c.IssuedAt, cmp, req)
}

// VerifyIssuer Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaim) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(c.Issuer, cmp, req)
}

// VerifyNotBefore Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaim) VerifyNotBefore(cmp int64, req bool) bool {
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
