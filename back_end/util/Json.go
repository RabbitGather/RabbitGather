package util

import (
	"encoding/json"
	"os"
)

func ParseJsonConfic(sst interface{},filePath string)(err error){
	configFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(sst)
	if err != nil {
		return
	}
	err =configFile.Close()
	return
}