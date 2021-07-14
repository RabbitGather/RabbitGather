package util

import "net/http"

const IDENTIFICATION_SYMBOL_KEY = "IDENTIFICATION_SYMBOL"
const IDENTIFICATION_SYMBOL = "fvqejfopj3/5<>?>9rm2ur#$TW 0924#$@T$#T$#^"

const ClientIP_KEY = "fue8asodxn8fewj8snxfpei"

func CheckIDENTIFICATION_SYMBOL(req *http.Request) bool {
	return req.Header.Get(IDENTIFICATION_SYMBOL_KEY) == IDENTIFICATION_SYMBOL

}
