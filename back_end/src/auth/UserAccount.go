package auth

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"

	//"github.com/dgrijalva/jwt-go"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/util"
)

var dbOperator db_operator.DBOperator

func init() {
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/auth_user_account.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.GetOperator(db_operator.Mysql, config.DatabaseConfig)
}

type UserInformation struct {
	Username   string        `json:"username"`
	Password   string        `json:"password"`
	Email      string        `json:"email,omitempty"`
	Phone      string        `json:"phone,omitempty"`
	Permission StatusBitmask `json:"permission"`
}

func CreateNewUserAccount(userinfo UserInformation) (*UserAccount, error) {
	statment := dbOperator.Statement("insert into user ( name, password, api_permission_bitmask) value (?,?,?);")
	password, err := HashPassword(userinfo.Password)
	if err != nil {
		return nil, err
	}
	res, err := statment.Exec(userinfo.Username, password, uint32(permission))
	if err != nil {
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	return &UserAccount{
		UserName:             username,
		UserID:               uint32(id),
		APIPermissionBitmask: permission,
	}, nil
}

type UserAccount struct {
	UserName             string
	UserID               uint32
	APIPermissionBitmask StatusBitmask
}

const NormalUserPermission = Login

func GetUserAccountByName(username string) (*UserAccount, error) {
	statment := dbOperator.Statement("select id,api_permission_bitmask from user where name = ? limit 1;")
	var id uint32
	var api_permission_bitmask uint32
	err := statment.QueryRow(username).Scan(&id, &api_permission_bitmask)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username not exist")
		}
		return nil, err
	}

	return &UserAccount{
		UserID:               id,
		UserName:             username,
		APIPermissionBitmask: StatusBitmask(api_permission_bitmask),
	}, nil
}

var ErrorPasswordWrong = errors.New("password wrong")

func (u *UserAccount) CheckPassword(password string) error {
	statment := dbOperator.Statement("select password from user where id = ? limit 1;")
	var passwordHash string
	err := statment.QueryRow(u.UserID).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("error when get password from user table on db")
		}
		return err
	}
	if CheckPasswordHash(password, passwordHash) {
		return nil
	} else {
		return ErrorPasswordWrong
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Create a token with this user claims
func (u *UserAccount) NewToken() (string, error) {
	theClaims := UtilityClaims{
		StatusClaimsName: StatusClaims{},
		UserClaimsName: UserClaims{
			UserName: u.UserName,
			UserID:   u.UserID,
		},
	}
	signedTokenString, err := SignToken(&theClaims)
	//token := JWTToken{}
	//token, err := NewSignedToken(&theConst)
	if err != nil {
		log.ERROR.Println(err.Error())
		return "", err
	}
	return signedTokenString, nil
}

//func GetUserByClaims(userClaims *UtilityClaims) UserAccount {
//	panic("Not implement")
//}
