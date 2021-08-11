package account_management

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/bitmask"
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
		ServerMailAddr         string `json:"server_mail_addr"`
		Username               string `json:"username"`
		Password               string `json:"password"`
		VerificationCodeLength int    `json:"verification_code_length"`
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "config/mail_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	gmailsender = mail.NewGmailSender(config.ServerMailAddr, config.Username, config.Password)
	verificationCodeLength = config.VerificationCodeLength
}

var verificationCodeLength int
var gmailsender *mail.GmailSender
var redisClient = redis_db.GetClient(0)

const (
	VerificationCodePurposeSignup = "signup"
)
const (
	EMAIL string = "email"
	Phone string = "phone"
)

// The SentVerificationCodeHandler handle send VerificationCode request
func (w *AccountManagement) SentVerificationCodeHandler(c *gin.Context) {
	err := sentVerificationCodeSecurity(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"err": "wrong input"},
		)
		log.DEBUG.Println("sentVerificationCodeSecurity Error: ", err.Error())
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

	verificationCode := NewVerificationCode()
	switch userInput.Type {
	case EMAIL:
		err = sentVerificationCodeToMail(verificationCode, userInput.Email)
		if err != nil {
			err = fmt.Errorf("SentMail: %w", sentVerificationCodeToMail(verificationCode, userInput.Email))
		}
	case Phone:
		err = sentVerificationCodeToPhone(verificationCode, userInput.Phone)
		if err != nil {
			err = fmt.Errorf("SentPhone: %w", sentVerificationCodeToMail(verificationCode, userInput.Email))
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

	utilityClaims, err := server.ContextAnalyzer(c).GetUtilityClaim()
	if err != nil {
		if err == server.ErrTokenIsEmpty {
			utilityClaims = &claims.UtilityClaim{
				claims.StandardClaimsName: claims.NewStandardClaims(),
				claims.StatusClaimsName:   claims.StatusClaim{StatusBitmask: bitmask.WaitVerificationCode},
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
			log.DEBUG.Println("error when GetUtilityClaim: ", err.Error())
			return
		}
	} else {
		var statusClaims claims.StatusClaim
		err = utilityClaims.GetSubClaims(&statusClaims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "token wrong"})
			log.DEBUG.Println("GetSubClaims error : ", err.Error())
			return
		}
		statusClaims.AppendBitMask(bitmask.WaitVerificationCode)
	}

	var standardClaims claims.StandardClaim
	_ = utilityClaims.GetSubClaims(&standardClaims)
	userInput.TokenID = standardClaims.Id

	stringToken, err := token.SignToken(utilityClaims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
		log.DEBUG.Println("SignToken error : ", err.Error())
		return
	}

	cacheVerificationCode(verificationCode, &userInput)
	c.JSON(http.StatusOK, gin.H{token.TokenKey: stringToken})
}

type VerificationCodeBindingPackage struct {
	SentTime time.Time `json:"sent_time" `
	Purpose  string    `json:"purpose" binding:"required"`
	Type     string    `json:"type" binding:"required"`
	Email    string    `json:"email,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	TokenID  string    `json:"-"`
}

// Verify verified if the userInput and purpose same as the time VerificationCode sent.
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

func (w *AccountManagement) getVerificationCodeBindingPackage(code string) (*VerificationCodeBindingPackage, error) {
	var vpk VerificationCodeBindingPackage
	err := redisClient.Get(context.Background(), code, &vpk)
	if err != nil {
		return nil, err
	}
	return &vpk, nil
}

func sentVerificationCodeSecurity(c *gin.Context) error {
	log.DEBUG.Println("Not implemented - sentVerificationCodeSecurity")
	return nil
}

func NewVerificationCode() string {
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", verificationCodeLength), util.RandomInLength(verificationCodeLength)) //util.NewVerificationCodeWithLength(verificationCodeLength)
}

func sentVerificationCodeToPhone(code, phone string) error {
	return fmt.Errorf("NOT IMPLEMENT - sentVerificationCode to Phone: %s, Code:%s", phone, code)
}

func sentVerificationCodeToMail(code, email string) error {
	log.DEBUG.Println("Sent VerificationCode To Mail: ", email)
	err := gmailsender.SendMail("Verification Code For RabbitGather", fmt.Sprintf("-- %s --", code), email)
	return err
}

const VerificationCodeTimeout = time.Minute * 30

func cacheVerificationCode(code string, pkg *VerificationCodeBindingPackage) {
	err := redisClient.Set(context.Background(), code, *pkg, VerificationCodeTimeout)
	if err != nil {
		log.ERROR.Println("fail to catch VerificationCode")
	}
}

func dropVerificationCode(code string) {
	log.DEBUG.Println("dropVerificationCode: ", code)
	redisClient.Del(context.Background(), code)
}
