package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/auth/token/claims"
)

type ContextAnalyzer struct {
	*gin.Context
}

func (c ContextAnalyzer) GetUtilityClaims() (*claims.UtilityClaims, bool, error) {
	if ut, exist := c.Get(token.TokenHeaderKey); exist {
		utilityClaims, ok := ut.(*claims.UtilityClaims)
		if !ok {
			return nil, false, errors.New("the token stored in context is not type of *claims.UtilityClaims")
		}
		return utilityClaims, true, nil
	}
	tokenRawString := c.GetHeader(token.TokenHeaderKey)
	if tokenRawString == "" {
		return nil, false, nil
	}
	utilityClaims, err := token.ParseToken(tokenRawString)
	if err != nil {
		return nil, false, err
	}
	c.Set(token.TokenHeaderKey, utilityClaims)
	return utilityClaims, true, nil
}
