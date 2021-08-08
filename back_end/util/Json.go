package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"os"
)

func ParseJsonConfic(sst interface{}, filePath string) (err error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(sst)
	if err != nil {
		return
	}
	err = configFile.Close()
	return
}

func ParseRequestJson(rawbody io.ReadCloser, st interface{}) error {
	body := json.NewDecoder(rawbody)
	body.DisallowUnknownFields()
	err := body.Decode(st)
	if err != nil {
		return err
	}
	return nil
}

// SQLJsonAble provide a simple default solution
// to make struct use as Json in SQL query
type SQLJsonAble struct {
}

func (m SQLJsonAble) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}

func (m *SQLJsonAble) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case []uint8:
		source = src.([]byte)
	case nil:
		return nil
	default:
		return errors.New("incompatible type")
	}
	err := json.Unmarshal(source, src)
	if err != nil {
		return err
	}
	return nil
}
