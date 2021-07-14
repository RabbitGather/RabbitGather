package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/util"
)

func CreateNewUserAccount(username, password string, permission APIPermissionBitmask) (*UserAccount, error) {
	statment := dbOperator.Statement("insert into user ( name, password, api_permission_bitmask) value (?,?,?);")
	res, err := statment.Exec(username, password, uint32(permission))
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
	APIPermissionBitmask APIPermissionBitmask
	token                *JWTToken
}

func (u *UserAccount) GetToken() (*JWTToken, error) {
	fmt.Println("Not implemented : GetToken")
	if u.token != nil {
		return u.token, nil
	}
	token, err := u.NewToken()
	if err != nil {
		return nil, err
	}
	u.token = token
	return token, nil
}

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

var ErrorPasswordWrong = errors.New("password wrong")

func (u *UserAccount) CheckPassword(password string) error {
	statment := dbOperator.Statement("select password from user where id = ? limit 1;\n")
	var theUserPassword string
	err := statment.QueryRow(u.UserID).Scan(&theUserPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("error when get password from user table on db")
		}
		return err
	}
	if theUserPassword == password {
		return nil
	} else {
		return ErrorPasswordWrong
	}
}

func (u *UserAccount) NewToken() (*JWTToken, error) {
	theConst := NormalUserClaims{
		PermissionClaims: PermissionClaims{
			StandardClaims:    *NewStandardClaims(),
			PermissionBitmask: u.APIPermissionBitmask,
		},
		UserName: u.UserName,
		UserID:   u.UserID,
	}
	//token := JWTToken{}
	token, err := NewSignedToken(&theConst)
	if err != nil {
		log.ERROR.Println(err.Error())
		return nil, err
	}
	return token, nil
}

type PermissionClaims struct {
	jwt.StandardClaims `json:"standard_claims"`
	PermissionBitmask  APIPermissionBitmask `json:"api_permission_bitmask"`
}

type NormalUserClaims struct {
	PermissionClaims
	UserName string `json:"user_name"`
	UserID   uint32 `json:"user_id"`
}

func GetUserByClaims(userClaims *PermissionClaims) UserAccount {
	panic("Not implement")
}
