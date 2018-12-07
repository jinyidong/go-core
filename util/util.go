package util

import (
	"github.com/json-iterator/go"
)

func StructToJson(model interface{}) string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(model)

	if err != nil {
		return ""
	}
	return string(data)
}

func ByteToStruct(data []byte, object interface{}) error {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	err := json.Unmarshal(data, object)

	return err
}
