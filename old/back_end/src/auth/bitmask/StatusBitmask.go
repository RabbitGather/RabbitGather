package bitmask

type StatusBitmask uint32

// status
const (

	// AllStatus should not be use in JWT but server permission filter only
	AllStatus = StatusBitmask(^uint32(0))

	// NoStatus should be used when you want to allow all kinds of Status through.
	// NoStatus should not be use in JWT but server permission filter only
	NoStatus = StatusBitmask(uint32(0))

	// Reject should be used when you want to close a port temporary.
	// Reject should not be use in JWT but server permission filter only
	Reject StatusBitmask = 1 << iota

	// Login means the user already login
	Login

	// WaitVerificationCode means the client is waiting for verification code
	WaitVerificationCode
)

// The MaskCheck is a shortcut of BitmaskA&BitmaskB != 0
func MaskCheck(BitmaskA StatusBitmask, BitmaskB ...StatusBitmask) bool {
	for _, bm := range BitmaskB {
		if BitmaskA&bm != 0 {
			return true
		}
	}
	return false
}

// The BitMaskMarge marge all the input bitmask with OR
func BitMaskMarge(bb ...uint64) uint64 {
	if len(bb) == 1 {
		return bb[0]
	}
	for _, b := range bb[1:] {
		bb[0] |= b
	}
	return bb[0]
}
