package claims

import (
	"errors"
	"rabbit_gather/src/auth/bitmask"
)

// The StatusClaim hold the StatusBitmask which represents the user's status now
type StatusClaim struct {
	StatusBitmask bitmask.StatusBitmask `json:"status_bitmask"`
}

// The StatusClaimsName is the string name represented of StatusClaim.
// Will be use as the key in UtilityClaim.
const StatusClaimsName = "status_claims"

func (s StatusClaim) Valid() error {
	switch s.StatusBitmask {
	case bitmask.NoStatus:
		return errors.New("the status_bitmask should not be \"NoStatus\"")
	case bitmask.AllStatus:
		return errors.New("the status_bitmask should not be \"AllStatus\"")
	case bitmask.Reject:
		return errors.New("the status_bitmask should not be \"Reject\"")
	default:
		return nil
	}
}

// AppendBitMask append a new status in "StatusBitmask" field
func (s *StatusClaim) AppendBitMask(code bitmask.StatusBitmask) {
	if s.StatusBitmask&code == 0 {
		s.StatusBitmask |= code
	}
}
