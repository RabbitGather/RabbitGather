package user_account

import (
	"database/sql"
	"errors"
	"rabbit_gather/src/auth/bitmask"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/auth/token/claims"
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
	err := util.ParseFileJsonConfig(&config, "config/auth_user_account.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.NewOperator(db_operator.Mysql, config.DatabaseConfig)
}

// CreateUserStruct is a struct required when creating  a new user
type CreateUserStruct struct {
	Username   string                `json:"username"`
	Password   string                `json:"password"`
	Email      string                `json:"email,omitempty"`
	Phone      string                `json:"phone,omitempty"`
	Permission bitmask.StatusBitmask `json:"permission"`
}

// ErrUserNameConflict will be thrown when the username is already exist in user table
var ErrUserNameConflict = errors.New("user name already exist")
var ErrEmailConflict = errors.New("email already exist")
var ErrPhoneConflict = errors.New("phone already exist")

// CreateNewUserAccount creat a new user record in DB tables, and return *UserAccount.
// throw ErrUserNameConflict when the username given is duplicate in user table.
// throw ErrEmailConflict when the email given is duplicate in user_information table.
// throw ErrEmailConflict when the phone given is duplicate in user_information table.
func CreateNewUserAccount(userinfo CreateUserStruct) (*UserAccount, error) {
	_, exist, err := checkUserExist(userinfo.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrUserNameConflict
	}
	exist, err = checkPhoneExist(userinfo.Phone)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrPhoneConflict
	}
	exist, err = checkEmailExist(userinfo.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrEmailConflict
	}
	tx, err := dbOperator.Begin()
	defer func(stmt *sql.Tx) {
		e := tx.Rollback()
		if e != sql.ErrTxDone && e != nil {
			log.ERROR.Println(e.Error())
		}
	}(tx)
	if err != nil {
		return nil, err
	}

	insertUserStatement := tx.Stmt(dbOperator.Statement("insert into user ( name, password,randomSalt, api_permission_bitmask) value (?,?,?,?);\n"))
	insertUserInformationStatement := tx.Stmt(dbOperator.Statement("insert into `user_information`(`user`,`email`,`phone`) value (?,?,?);"))

	password, randomSalt, err := util.HashPassword(userinfo.Password)
	if err != nil {
		return nil, err
	}
	res, err := insertUserStatement.Exec(userinfo.Username, password, randomSalt, uint32(userinfo.Permission))
	if err != nil {
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	_, err = insertUserInformationStatement.Exec(id, userinfo.Email, userinfo.Phone)
	if err != nil {
		log.DEBUG.Println(err.Error())
		return nil, err
	}
	if e := tx.Commit(); e != nil {
		log.DEBUG.Println(e.Error())
		return nil, e
	}
	return &UserAccount{
		UserName:             userinfo.Username,
		UserID:               uint32(id),
		APIPermissionBitmask: userinfo.Permission,
	}, nil
}

// checkEmailExist check if the user's email exist.
func checkEmailExist(email string) (bool, error) {
	statement := dbOperator.Statement("select `email` from `user_information`where email =?;")
	err := statement.QueryRow(email).Scan(new(string))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.ERROR.Println("error when scanning userid")
		return false, err
	}
	return true, nil
}

// checkPhoneExist check if the user's phone exist.
func checkPhoneExist(phone string) (bool, error) {
	statement := dbOperator.Statement("select `phone` from `user_information`where phone =?;")
	err := statement.QueryRow(phone).Scan(new(string))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.ERROR.Println("error when scanning userid")
		return false, err
	}
	return true, nil
}

// checkUserExist check if the user account exist.
func checkUserExist(username string) (uint64, bool, error) {
	statement := dbOperator.Statement("select id from user where name = ? limit 1;")
	var userid = uint64(0)
	err := statement.QueryRow(username).Scan(&userid)
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

// The UserAccount is represented of a user account provide some function to
// operate user account.
type UserAccount struct {
	UserName             string
	UserID               uint32
	APIPermissionBitmask bitmask.StatusBitmask
}

// GetUserAccountByName get the user's id and api_permission_bitmask from DB by username.
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
		APIPermissionBitmask: bitmask.StatusBitmask(api_permission_bitmask),
	}, nil
}

var ErrorPasswordWrong = errors.New("password wrong")

// CheckPassword check if the password given is the password of this user.
func (u *UserAccount) CheckPassword(password string) error {
	statement := dbOperator.Statement("select password,randomSalt from user where id = ? limit 1;")
	var passwordHash string
	var randomSalt string
	err := statement.QueryRow(u.UserID).Scan(&passwordHash, &randomSalt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("error when get password from user table on db")
		}
		return err
	}
	if util.CheckPasswordHash(password, passwordHash, randomSalt) {
		return nil
	} else {
		return ErrorPasswordWrong
	}
}

//func (u *UserAccount) GetUserClaim() claims.UserClaim {
//	return claims.UserClaim{
//		UserName: u.UserName,
//		UserID:   u.UserID,
//	}
//}

// NewToken Create a token with claims with this user situations
func (u *UserAccount) NewToken() (string, error) {
	theClaims := claims.UtilityClaim{
		claims.StatusClaimsName: claims.StatusClaim{
			StatusBitmask: u.APIPermissionBitmask,
		},
		claims.UserClaimsName: claims.UserClaim{
			UserName: u.UserName,
			UserID:   u.UserID,
		},
	}
	signedTokenString, err := token.SignToken(&theClaims)

	if err != nil {
		log.ERROR.Println(err.Error())
		return "", err
	}
	return signedTokenString, nil
}
