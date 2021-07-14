package auth

//type Code uint32
type APIPermissionBitmask uint32

//type APINames Code

const (
	// public open api
	Public = APIPermissionBitmask(^uint32(0))

	// Only Admin can use
	Admin = APIPermissionBitmask(uint32(0))
	//Admin =NoAccess-1

	// login status
	Login APIPermissionBitmask = 1 << iota // token is malformed

	//// after the page loaded
	//PageLoad

	// Wait for VerificationCode status
	WaitVerificationCode
)

func BitMaskCheck(PermissionA, PermissionB APIPermissionBitmask) bool {
	return PermissionA&PermissionB != 0
}
