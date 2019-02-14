package parser

import (
	"bytes"
	"encoding/json"
	"github.com/mattn/go-jsonpointer"
	"strings"
)

const MARSHAL_PREFIX = ""
const MARSHAL_INDENT = "    "

func JsonToByte(jsondata interface{}) ([]byte, error) {
	return json.MarshalIndent(jsondata, MARSHAL_PREFIX, MARSHAL_INDENT)
}

func ByteToJson(jsonbyte []byte) (interface{}, error) {
	var jsondata interface{}
	err := json.Unmarshal(jsonbyte, &jsondata)
	return jsondata, err
}

func RemoveJsonItems(jsondata interface{}, removelist []string) (interface{}, error) {
	if jsondata == nil || removelist == nil || len(removelist) == 0 {
		return jsondata, nil
	}
	for _, remove := range removelist {
		if removedjsondata, err := jsonpointer.Remove(jsondata, MakeJsonpointerKey(remove)); err != nil {
			return jsondata, err
		} else {
			jsondata = removedjsondata
		}
	}
	return jsondata, nil

}
func GetJsonItems(jsondata []byte, searchpointer string) (interface{}, error) {
	var jsonobj interface{}
	if bytes.Equal(jsondata, make([]byte, 0)) == false {
		if err := json.Unmarshal(jsondata, &jsonobj); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte("{}"), &jsonobj); err != nil {
			return nil, err
		}
	}
	return jsonpointer.Get(jsonobj, MakeJsonpointerKey(searchpointer))
}

func AddJsonItems(jsondata []byte, addlist []map[string]interface{}) ([]byte, error) {
	var jsonobj interface{}
	if bytes.Equal(jsondata, make([]byte, 0)) == false {
		if err := json.Unmarshal(jsondata, &jsonobj); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte("{}"), &jsonobj); err != nil {
			return nil, err
		}
	}
	for _, add := range addlist {
		for key, val := range add {
			if err := jsonpointer.Set(jsonobj, MakeJsonpointerKey(key), val); err != nil {
				return nil, err
			}
		}
	}
	return json.MarshalIndent(jsonobj, MARSHAL_PREFIX, MARSHAL_INDENT)

}

func MakeJsonpointerKey(key string) string {
	if strings.HasPrefix(key, "/") == false {
		return "/" + key
	}
	return key
}
