package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

const MARSHAL_PREFIX = ""
const MARSHAL_INDENT = "    "

func JsonDecode(i io.Reader, s interface{}) error {
	bytes, err := ioutil.ReadAll(i)
	if err != nil {
		return nil
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func JsonToByte(jsondata map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(jsondata, MARSHAL_PREFIX, MARSHAL_INDENT)
}

func ByteToJson(jsonbyte []byte) (map[string]interface{}, error) {
	var jsondata map[string]interface{}
	err := json.Unmarshal(jsonbyte, &jsondata)
	return jsondata, err
}
