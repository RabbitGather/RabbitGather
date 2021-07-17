package util

import "io/ioutil"

var plainTextMap = map[string]string{}

func GetFileStoredPlainText(fileName string) string {
	if resStr, exist := plainTextMap[fileName]; exist {
		return resStr
	}
	btarr, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var resStr = string(btarr)
	plainTextMap[fileName] = resStr
	return resStr
}
