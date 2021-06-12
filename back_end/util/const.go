package util

import "net/http"

const IDENTIFICATION_SYMBOL_KEY = "IDENTIFICATION_SYMBOL"
const IDENTIFICATION_SYMBOL = "fvqejfopj3/5<>?>9rm2ur#$TW 0924#$@T$#T$#^"

func CheckIDENTIFICATION_SYMBOL(req *http.Request) bool {
	return req.Header.Get(IDENTIFICATION_SYMBOL_KEY) == IDENTIFICATION_SYMBOL

}
