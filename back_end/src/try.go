package main

import (
	"fmt"
)

type ABC string

func (a ABC) name() string {
	return string(a)
}

const (
	ValidationErrorMalformed        uint32 = 1 << iota // Token is malformed
	ValidationErrorUnverifiable                        // Token could not be verified because of signing problems
	ValidationErrorSignatureInvalid                    // Signature validation failed

	// Standard Claim validation errors
	ValidationErrorAudience      // AUD validation failed
	ValidationErrorExpired       // EXP validation failed
	ValidationErrorIssuedAt      // IAT validation failed
	ValidationErrorIssuer        // ISS validation failed
	ValidationErrorNotValidYet   // NBF validation failed
	ValidationErrorId            // JTI validation failed
	ValidationErrorClaimsInvalid // Generic claims validation error
)

func main() {
	//if 1&1{
	//}
	//var Errors =  uint32(1)
	fmt.Println(1 << 0)
	//fmt.Println(jwt.ValidationErrorMalformed )
	//fmt.Println(Errors)
	//fmt.Println(Errors&jwt.ValidationErrorMalformed )
	//fmt.Println(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet )
	//fmt.Println(1 &(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) )
	//
	////fmt.Println(jwt.ValidationErrorMalformed )
	////fmt.Println(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet)
	//if Errors&jwt.ValidationErrorMalformed != 0 {
	//	fmt.Println("That's not even a token")
	//} else if Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	//	fmt.Println("Timing is everything")
	//} else {
	//	fmt.Println("Couldn't handle this token:")
	//}

	//a := ABC("456")
	//fmt.Println(a.name())
	//fmt.Println("4444444")
	//theuuid := uuid.New()
	//a := [16]byte(uuid.New())
	//p ,_:=filepath.Abs("ssl/crt/meowalien_com.crt")
	//fmt.Println("ssl/crt/meowalien_com.crt  --  ",p) // F:\GoTest\GoTest\master.exe <nil>
	//fmt.Println(string(a[:]))
	//fmt.Println(string(uuid.New().NodeID()))
}
