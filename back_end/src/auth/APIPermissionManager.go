package auth

type PermissionCode uint32

// API list
const (
	SearchArticle PermissionCode = 1 << iota // Token is malformed
	PostArticle
	Login
	GetPeerID
	PeerWebsockt
)

func APIAuthorizationCheck(PermissionA, PermissionB PermissionCode) bool {
	return PermissionA&PermissionB != 0
}
