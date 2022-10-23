package claims

// The UserClaimsName is the string name represented of UserClaim.
// Will be use as the key in UtilityClaim.
const UserClaimsName = "user_claims"

type UserClaim struct {
	// The UserName is the name in the user table on DB
	UserName string `json:"user_name"`
	// The UserID is the serial ID in the user table on DB
	UserID uint32 `json:"user_id"`
}

func (u UserClaim) Valid() error {
	return nil
}
