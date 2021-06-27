package util

import (
	"encoding/json"
	"io"
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

func ParseRequestJson(rawbody io.ReadCloser,st interface{})error{
	body := json.NewDecoder(rawbody)
	body.DisallowUnknownFields()
	err := body.Decode(st)
	if err != nil {
		return err
	}
	return  nil
}
