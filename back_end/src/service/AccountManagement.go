package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"net/http"
	"rabbit_gather/src/auth"
	"rabbit_gather/src/db_operator"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

func (w *AccountManagement) LoginHandler(c *gin.Context) {
	log.DEBUG.Println("login")
	userinput := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := c.ShouldBindJSON(&userinput)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"err": "wrong input",
			})
		log.DEBUG.Printf("ShouldBindJSON error : %s", err.Error())
		return
	}
	log.DEBUG.Println("Username : ", userinput.Username)
	log.DEBUG.Println("Password : ", userinput.Password)

	userInst, err := w.GetUserAccountByName(userinput.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "wrong input"})
		log.DEBUG.Printf("GetUserAccountByName error: %s", err.Error())
		return
	}
	err = userInst.CheckPassword(userinput.Password)
	if err != nil {
		if err == auth.ErrorPasswordWrong {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "wrong input"})
			log.DEBUG.Println("password wrong: ", err.Error())
			return
		} else {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"err": "CheckPassword error"})
			log.DEBUG.Printf("CheckPassword Error: %s", err.Error())
			return
		}

	}
	userToken, err := userInst.GetToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "GetToken error"})
		log.DEBUG.Printf("login - GetToken error : %s", err.Error())
		return
	}
	c.JSON(200, gin.H{
		"ok":    true,
		"err":   "",
		"token": userToken.GetSignedString(),
	})
}

const NormalUserPermission = auth.Login

func (w *AccountManagement) SignupHandler(c *gin.Context) {
	var userinput SignupUserInput
	if err := c.ShouldBindJSON(&userinput); err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"err": "wrong input"},
		)
		return
	}

	log.DEBUG.Println("Username : ", userinput.Username)
	log.DEBUG.Println("Password : ", userinput.Password)
	log.DEBUG.Println("VerificationCode : ", userinput.VerificationCode)

	bindingPackage, err := w.getVerificationCodeBindingPackage(userinput.VerificationCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "verification_code wrong"})
		log.DEBUG.Printf("VerificationCode wrong : %s", err.Error())
		return
	}
	if !bindingPackage.Verify(userinput) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "verification_code wrong"})
		log.DEBUG.Println("VerificationCode Verify fail : %s")
		return
	}

	userAccount, err := auth.CreateNewUserAccount(userinput.Username, userinput.Password, NormalUserPermission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "create new user account error"})
		log.ERROR.Printf("Signup - GetToken error : %s", err.Error())
		return
	}
	userToken, err := userAccount.GetToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "create new user token error"})
		log.ERROR.Printf("Signup - GetToken error : %s", err.Error())
		return
	}
	w.dropVerificationCode(userinput.VerificationCode)
	c.JSON(200, gin.H{
		"token": userToken.GetSignedString(),
	})
}

func (w *AccountManagement) SentVerificationCodeHandler(c *gin.Context) {
	// check if repeat request in too short time.
	log.TempLog().Println("AccountManagement GetClientIP(c): ", util.GetClientIP(c))
	log.TempLog().Println("AccountManagement Header: ", pretty.Sprint(c.Request.Header))

	var userinput VerificationCodeBindingPackage
	if err := c.ShouldBindJSON(&userinput); err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"err": "wrong input"},
		)
		log.DEBUG.Println("wrong input when ShouldBindJSON: ", err.Error())
		return
	}
	verificationCode := w.NewVerificationCode()
	switch userinput.Type {
	case EMAIL:
		w.sentVerificationCodeToMail(verificationCode, userinput.Email)
	case Phone:
		w.sentVerificationCodeToPhone(verificationCode, userinput.Phone)
	default:
		log.ERROR.Println("unknown type")
		return
	}
	existToken := c.GetHeader(auth.TokenHeaderKey)
	if existToken == "" {
		//	create token
		claims := auth.PermissionClaims{
			StandardClaims:    *auth.NewStandardClaims(),
			PermissionBitmask: auth.WaitVerificationCode,
		}
		token, err := auth.NewSignedToken(&claims)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"err": "fail to create WaitVerificationCode token"},
			)
			log.DEBUG.Println("fail to create WaitVerificationCode token: ", err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{auth.TokenHeaderKey: token.GetSignedString()})
	} else {
		//	append permission bit mask
		var claims auth.PermissionClaims
		token, err := auth.ParseToken(existToken, &claims)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"err": "fail with ParseToken"},
			)
			log.DEBUG.Println("fail with ParseToken")

			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"err": "token not Valid"},
			)
			log.DEBUG.Println("token not Valid")
			return
		}
		newToken, err := token.AppendBitMask(auth.WaitVerificationCode)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"err": "error when AppendBitMask"},
			)
			log.DEBUG.Println("error when AppendBitMask: ", err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{auth.TokenHeaderKey: newToken.GetSignedString()})
	}

	w.catchVerificationCode(verificationCode, &userinput)
}

type AccountManagement struct {
}

var log *logger.LoggerWrapper

func init() {
	log = logger.NewLoggerWrapper("AccountManagement")
}

var dbOperator db_operator.DBOperator

func init() {
	type Config struct {
		DatabaseConfig db_operator.DatabaseConnectionConfiguration `json:"database_config"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/account_management.config.json")
	if err != nil {
		panic(err.Error())
	}
	dbOperator = db_operator.GetOperator(db_operator.Mysql, config.DatabaseConfig)
}

func (m *AccountManagement) GetUserAccountByName(username string) (*auth.UserAccount, error) {
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

	return &auth.UserAccount{
		UserID:               id,
		UserName:             username,
		APIPermissionBitmask: auth.APIPermissionBitmask(api_permission_bitmask),
	}, nil
}

type SignupUserInput struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
	Email            string `json:"email,omitempty"`
	Phone            string `json:"phone,omitempty"`
}

func (w *AccountManagement) getVerificationCodeBindingPackage(code string) (*VerificationCodeBindingPackage, error) {
	// --- here must have format check
	pkg, exist := verificationCodeMap[code]
	if !exist {
		return nil, errors.New("verification code not found")
	}
	return pkg, nil
}

type VerificationCodeBindingPackage struct {
	Type  string `json:"type"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (p *VerificationCodeBindingPackage) Verify(userinput SignupUserInput) bool {
	switch p.Type {
	case EMAIL:
		return userinput.Email == p.Email
	case Phone:
		return userinput.Phone == p.Phone
	default:
		log.ERROR.Println("unknown type")
		return false
	}
}

var verificationCodeMap = map[string]*VerificationCodeBindingPackage{}

const (
	EMAIL string = "email"
	Phone string = "phone"
)

func (w *AccountManagement) catchVerificationCode(code string, pkg *VerificationCodeBindingPackage) {
	log.DEBUG.Println("catchVerificationCode: ", code)
	verificationCodeMap[code] = pkg
}
func (w *AccountManagement) sentVerificationCodeToMail(code, email string) error {
	log.DEBUG.Printf("NOT IMPLEMENT - sentVerificationCode to Mail: %s, Code:%s", email, code)
	return nil
}
func (w *AccountManagement) sentVerificationCodeToPhone(code, phone string) error {
	log.DEBUG.Printf("NOT IMPLEMENT - sentVerificationCode to Phone: %s, Code:%s", phone, code)
	return nil
}

const VerificationCodeLength = 4

func (w *AccountManagement) NewVerificationCode() string {
	//a := fmt.Sprintf("\\%")//"%4d"
	return util.NewVerificationCodeWithLength(VerificationCodeLength)
}

func (w *AccountManagement) dropVerificationCode(code string) {
	delete(verificationCodeMap, code)
}
