package auth

//type Code uint32
type APIPermissionBitmask uint32

//type APINames Code

const (
	// Everyone can use
	Public = APIPermissionBitmask(^uint32(0))

	// Only Admin can use
	Admin = APIPermissionBitmask(uint32(0))

	// Only Login User can use
	Login APIPermissionBitmask = 1 << iota // token is malformed

	// Only VIP User can use
	VIP
)

func APIAuthorizationCheck(PermissionA, PermissionB APIPermissionBitmask) bool {
	return PermissionA&PermissionB != 0
}
