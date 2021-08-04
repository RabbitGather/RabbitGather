package claims

import "rabbit_gather/src/auth/status_bitmask"

const StatusClaimsName = "status_claims"

type StatusClaims struct {
	StatusBitmask status_bitmask.StatusBitmask `json:"status_bitmask"`
}

func (s StatusClaims) Valid() error {
	return nil
}

func (s *StatusClaims) AppendBitMask(code status_bitmask.StatusBitmask) {
	if s.StatusBitmask&code == 0 {
		s.StatusBitmask |= code
	}
}
