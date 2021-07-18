package user_account

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"rabbit_gather/src/auth/claims"
	"rabbit_gather/src/auth/status_bitmask"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

var log = logger.NewLoggerWrapper("user_account")
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
	Username   string                       `json:"username"`
	Password   string                       `json:"password"`
	Email      string                       `json:"email,omitempty"`
	Phone      string                       `json:"phone,omitempty"`
	Permission status_bitmask.StatusBitmask `json:"permission"`
}

var UserNameExist = errors.New("user name already exist")

func CreateNewUserAccount(userinfo UserInformation) (*UserAccount, error) {
	_, exist, err := CheckUserExist(userinfo.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, UserNameExist
	}

	statment := dbOperator.Statement("insert into user ( name, password, api_permission_bitmask) value (?,?,?);")
	password, err := HashPassword(userinfo.Password)
	if err != nil {
		return nil, err
	}
	res, err := statment.Exec(userinfo.Username, password, uint32(userinfo.Permission))
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
		UserName:             userinfo.Username,
		UserID:               uint32(id),
		APIPermissionBitmask: userinfo.Permission,
	}, nil
}

func CheckUserExist(username string) (uint64, bool, error) {
	statment := dbOperator.Statement("select id from user where name = ? limit 1;")
	var userid = uint64(0)
	err := statment.QueryRow(username).Scan(&userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil

		}
		log.ERROR.Println("error when scanning userid")
		return 0, false, err
	}
	if userid != 0 {
		return userid, true, nil
	} else {
		return 0, false, nil
	}
}

type UserAccount struct {
	UserName             string
	UserID               uint32
	APIPermissionBitmask status_bitmask.StatusBitmask
}

const NormalUserPermission = status_bitmask.Login

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
		APIPermissionBitmask: status_bitmask.StatusBitmask(api_permission_bitmask),
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
func (u *UserAccount) NewToken(status status_bitmask.StatusBitmask) (string, error) {
	theClaims := claims.UtilityClaims{
		claims.StatusClaimsName: claims.StatusClaims{
			StatusBitmask: status,
		},
		claims.UserClaimsName: claims.UserClaims{
			UserName: u.UserName,
			UserID:   u.UserID,
		},
	}
	signedTokenString, err := token.SignToken(&theClaims)
	//token := JWTToken{}
	//token, err := NewSignedToken(&theConst)
	if err != nil {
		log.ERROR.Println(err.Error())
		return "", err
	}
	return signedTokenString, nil
}

//func GetUserByClaims(userClaims *UtilityClaims) user_account {
//	panic("Not implement")
//}
