package account_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/status_bitmask"
	"rabbit_gather/src/auth/user_account"
)

type SignupUserInput struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
	Email            string `json:"email,omitempty"`
	Phone            string `json:"phone,omitempty"`
}

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
	log.DEBUG.Println("Email : ", userinput.Email)
	log.DEBUG.Println("Phone : ", userinput.Phone)
	log.DEBUG.Println("VerificationCode : ", userinput.VerificationCode)

	bindingPackage, err := w.getVerificationCodeBindingPackage(userinput.VerificationCode)
	dropVerificationCode(userinput.VerificationCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "verification code wrong"})
		log.DEBUG.Printf("VerificationCode not exist: %s", err.Error())
		return
	}

	if !bindingPackage.Verify(userinput, VerificationCodePurposeSignup) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "verification code wrong"})
		log.DEBUG.Printf("VerificationCode Verify fail")
		return
	}

	userAccount, err := user_account.CreateNewUserAccount(user_account.UserInformation{
		Username:   userinput.Username,
		Password:   userinput.Password,
		Email:      userinput.Email,
		Phone:      userinput.Phone,
		Permission: user_account.NormalUserPermission,
	})
	if err != nil {
		if err == user_account.UserNameExist {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "user name exist."})
			log.ERROR.Printf("CreateNewUserAccount error : %s", err.Error())
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "Create Account error"})
			log.ERROR.Printf("CreateNewUserAccount error : %s", err.Error())
			return
		}

	}

	userToken, err := userAccount.NewToken(status_bitmask.Login)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "create new user token error"})
		log.ERROR.Printf("Signup - GetToken error : %s", err.Error())
		return
	}

	c.JSON(200, gin.H{
		"token": userToken,
	})
}
