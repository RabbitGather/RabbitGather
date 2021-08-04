package claims

const UserClaimsName = "user_claims"

type UserClaims struct {
	UserName string `json:"user_name"`
	UserID   uint32 `json:"user_id"`
}

func (u UserClaims) Valid() error {
	return nil
}
