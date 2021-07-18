package status_bitmask

//type Code uint32
type StatusBitmask uint32

//type APINames Code

const (
	AllStatus = StatusBitmask(^uint32(0))

	NoStatus = StatusBitmask(uint32(0))

	// this should be used when you want a placeholder or you want to stop some process.
	Reject StatusBitmask = 1 << iota // token is malformed

	Login

	// Wait for VerificationCode status
	WaitVerificationCode
)

func BitMaskCheck(PermissionA, PermissionB StatusBitmask) bool {
	return PermissionA&PermissionB != 0
}
func BitMaskMarge(bb ...uint64) uint64 {
	if len(bb) == 1 {
		return bb[0]
	}
	for _, b := range bb[1:] {
		bb[0] |= b
	}
	return bb[0]
}
