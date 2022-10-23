package claims

import (
	"rabbit_gather/src/logger"
)

var log = logger.NewLoggerWrapper("token.claims")

// The UtilityClaimKey is the key use when set parsed UtilityClaimKey in the gin.Context
const UtilityClaimKey = "utility_claim"
