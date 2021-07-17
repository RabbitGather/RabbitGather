package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth"
	"rabbit_gather/util"
	"sync"
	"time"
)

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

	existToken := c.GetHeader(util.TokenHeaderKey)
	token, err := auth.UpdateStatus(existToken, auth.WaitVerificationCode)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"err": "fail to update Status"},
		)
		return
	}
	w.catchVerificationCode(verificationCode, &userInput)
	c.JSON(http.StatusOK, gin.H{util.TokenHeaderKey: token})
}

func (w *AccountManagement) SentVerificationCodeSecurity(c *gin.Context) error {
	log.DEBUG.Println("Not implemented - SentVerificationCodeSecurity")
	return nil
}

func (w *AccountManagement) catchVerificationCode(code string, pkg *VerificationCodeBindingPackage) {
	log.DEBUG.Println("catchVerificationCode: ", code)
	pkg.SentTime = time.Now()
	verificationCodeMap.Store(code, pkg)

}

func (w *AccountManagement) sentVerificationCodeToPhone(code, phone string) error {
	log.DEBUG.Printf("NOT IMPLEMENT - sentVerificationCode to Phone: %s, Code:%s", phone, code)
	return nil
}

func (w *AccountManagement) sentVerificationCodeToMail(code, email string) error {
	log.DEBUG.Printf("NOT IMPLEMENT - sentVerificationCode to Mail: %s, Code:%s", email, code)
	return nil
}

const VerificationCodeLength = 4

func (w *AccountManagement) NewVerificationCode() string {
	return util.NewVerificationCodeWithLength(VerificationCodeLength)
}

type VerificationCodeBindingPackage struct {
	SentTime time.Time
	Type     string `json:"type"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

func (p *VerificationCodeBindingPackage) Verify(userInput SignupUserInput) bool {
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

var verificationCodeMap = sync.Map{}

func (w *AccountManagement) getVerificationCodeBindingPackage(code string) (*VerificationCodeBindingPackage, error) {
	// --- here must have format check
	pkg, exist := verificationCodeMap.Load(code)
	if !exist {
		return nil, errors.New("verification code not found")
	} else if time.Since(pkg.(*VerificationCodeBindingPackage).SentTime) > VerificationCodeExpirationDuration {
		dropVerificationCode(code)
		return nil, errors.New("timeout")
	}

	return pkg.(*VerificationCodeBindingPackage), nil
}

var cleanVerificationCodeMapTimeout = time.Minute * 30

func init() {
	if util.DebugMode {
		cleanVerificationCodeMapTimeout = time.Minute * 2
	}
}

var cleanVerificationCodeMapCancel = make(chan struct{})

func init() {
	util.RunAfterFuncWithCancel(cleanVerificationCodeMapTimeout, cleanVerificationCodeMap, cleanVerificationCodeMapCancel)
}

func cleanVerificationCodeMap() {
	log.DEBUG.Println("cleanVerificationCodeMap running...")
	verificationCodeMap.Range(func(key, value interface{}) bool {
		pkg, ok := value.(*VerificationCodeBindingPackage)
		if !ok {
			log.ERROR.Println("not VerificationCodeBindingPackage object in verificationCodeMap")
			return false
		}
		if time.Since(pkg.SentTime) > VerificationCodeExpirationDuration {
			dropVerificationCode(key.(string))
		}
		return true
	})
}

func dropVerificationCode(code string) {
	log.DEBUG.Println("dropVerificationCode: ", code)

	verificationCodeMap.Delete(code)
}
