package account_management

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/status_bitmask"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/mail"
	"rabbit_gather/src/redis_db"
	"rabbit_gather/src/server"
	"rabbit_gather/util"
	//"sync"
	"time"
)

func init() {
	type Config struct {
		ServerMailAddr string `json:"server_mail_addr"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/mail_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	gmailsender = mail.NewGmailSender(config.ServerMailAddr, config.Username, config.Password)
}

var gmailsender *mail.GmailSender

const (
	EMAIL string = "email"
	Phone string = "phone"
)
const VerificationCodeExpirationDuration = time.Minute * 3

func (w *AccountManagement) SentVerificationCodeHandler(c *gin.Context) {
	err := w.SentVerificationCodeSecurity(c)
	if err != nil {
		c.Abort()
		log.DEBUG.Println("SentVerificationCodeSecurity Error: ", err.Error())
		return
	}

	var userInput VerificationCodeBindingPackage
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"err": "wrong input"},
		)
		log.DEBUG.Println("ShouldBindJSON - wrong input: ", err.Error())
		return
	}

	log.TempLog().Println("userInput: ", userInput)

	verificationCode := w.NewVerificationCode()
	switch userInput.Type {
	case EMAIL:
		err = w.sentVerificationCodeToMail(verificationCode, userInput.Email)
		if err != nil {
			err = fmt.Errorf("SentMail: %w", w.sentVerificationCodeToMail(verificationCode, userInput.Email))
		}
	case Phone:
		err = w.sentVerificationCodeToPhone(verificationCode, userInput.Phone)
		if err != nil {
			err = fmt.Errorf("SentPhone: %w", w.sentVerificationCodeToMail(verificationCode, userInput.Email))
		}
	default:
		log.ERROR.Println("unknown type:", userInput.Type)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"err": "fail to sent VerificationCode"},
		)
		log.ERROR.Println("fail to sent VerificationCode: ", err.Error())
		return
	}
	utilityClaims, exist, err := server.ContextAnalyzer{Context: c}.GetUtilityClaims()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
		log.DEBUG.Println("error when GetUtilityClaims: ", err.Error())
		return
	} else {
		if exist {
			var statusClaims claims.StatusClaims
			err = utilityClaims.GetSubClaims(&statusClaims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "token wrong"})
				log.DEBUG.Println("GetSubClaims error : ", err.Error())
				return
			}
			statusClaims.AppendBitMask(status_bitmask.WaitVerificationCode)
		} else {
			utilityClaims = &claims.UtilityClaims{
				claims.StandardClaimsName: token.NewStandardClaims(),
				claims.StatusClaimsName:   claims.StatusClaims{StatusBitmask: status_bitmask.WaitVerificationCode},
			}
		}
	}
	var standardClaims claims.StandardClaims
	_ = utilityClaims.GetSubClaims(&standardClaims)
	userInput.TokenID = standardClaims.Id

	stringToken, err := token.SignToken(utilityClaims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
		log.DEBUG.Println("SignToken error : ", err.Error())
		return
	}

	w.catchVerificationCode(verificationCode, &userInput)
	c.JSON(http.StatusOK, gin.H{token.TokenHeaderKey: stringToken})
}

func (w *AccountManagement) SentVerificationCodeSecurity(c *gin.Context) error {
	log.DEBUG.Println("Not implemented - SentVerificationCodeSecurity")
	return nil
}

var client = redis_db.GetClient(0)

func (w *AccountManagement) catchVerificationCode(code string, pkg *VerificationCodeBindingPackage) {
	err := client.Set(context.Background(), code, *pkg, cleanVerificationCodeMapTimeout)
	if err != nil {
		log.ERROR.Println("fail to catch VerificationCode")
	}
	//pkg.SentTime = time.Now()
	//log.DEBUG.Printf("catchVerificationCode:%s ,Pkg: %s", code, fmt.Sprint(*pkg))
	//verificationCodeMap.Store(code, pkg)
}

func (w *AccountManagement) sentVerificationCodeToPhone(code, phone string) error {
	return fmt.Errorf("NOT IMPLEMENT - sentVerificationCode to Phone: %s, Code:%s", phone, code)
}

func (w *AccountManagement) sentVerificationCodeToMail(code, email string) error {
	err := gmailsender.SendMail(fmt.Sprintf("-- %s --", code), "Verification Code For RabbitGather", email)
	return err
}

const VerificationCodeLength = 4

func (w *AccountManagement) NewVerificationCode() string {
	return util.NewVerificationCodeWithLength(VerificationCodeLength)
}

const (
	VerificationCodePurposeSignup = "signup"
)

type VerificationCodeBindingPackage struct {
	SentTime time.Time `json:"sent_time" `
	Purpose  string    `json:"purpose" binding:"required"`
	Type     string    `json:"type" binding:"required"`
	Email    string    `json:"email,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	TokenID  string    `json:"-"`
}

func (p *VerificationCodeBindingPackage) Verify(userInput SignupUserInput, purpose string) bool {
	if p.Purpose != purpose {
		return false
	}
	switch p.Type {
	case EMAIL:
		return userInput.Email == p.Email
	case Phone:
		return userInput.Phone == p.Phone
	default:
		log.ERROR.Println("unknown type")
		return false
	}
}

//var verificationCodeMap = sync.Map{}

func (w *AccountManagement) getVerificationCodeBindingPackage(code string) (*VerificationCodeBindingPackage, error) {
	var vpk VerificationCodeBindingPackage
	err := client.Get(context.Background(), code, &vpk)
	if err != nil {
		return nil, err
	}
	return &vpk, nil
	// --- here must have format check
	//pkg, exist := verificationCodeMap.Load(code)
	//if !exist {
	//	return nil, errors.New("verification code not found")
	//} else if time.Since(pkg.(*VerificationCodeBindingPackage).SentTime) > VerificationCodeExpirationDuration {
	//	dropVerificationCode(code)
	//	return nil, errors.New("timeout")
	//}
	//return pkg.(*VerificationCodeBindingPackage), nil
}

var cleanVerificationCodeMapTimeout time.Duration = time.Minute * 30

//var cleanVerificationCodeMapCancel = make(chan struct{})

//func init() {
//	util.RunAfterFuncWithCancel(cleanVerificationCodeMapTimeout, cleanVerificationCodeMap, cleanVerificationCodeMapCancel)
//}
//
//func cleanVerificationCodeMap() {
//	log.DEBUG.Println("cleanVerificationCodeMap running...")
//	verificationCodeMap.Range(func(key, value interface{}) bool {
//		pkg, ok := value.(*VerificationCodeBindingPackage)
//		if !ok {
//			log.ERROR.Println("not VerificationCodeBindingPackage object in verificationCodeMap")
//			return false
//		}
//		if time.Since(pkg.SentTime) > VerificationCodeExpirationDuration {
//			dropVerificationCode(key.(string))
//		}
//		return true
//	})
//}

func dropVerificationCode(code string) {
	log.DEBUG.Println("dropVerificationCode: ", code)
	client.Del(context.Background(), code)

	//verificationCodeMap.Delete(code)
}
