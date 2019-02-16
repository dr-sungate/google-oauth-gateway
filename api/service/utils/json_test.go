package utils

import (
	"reflect"
	"testing"
)

func Test_JsonToByte_ByteToJson(t *testing.T) {
	inputdata1 := map[string]interface{}{
		"message1": "ok",
		"message2": "ng",
	}
	result, err := JsonToByte(inputdata1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log("success JsonToByte")
	t.Log(string(result))
	resultmap, err := ByteToJson(result)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if reflect.DeepEqual(inputdata1, resultmap) == false {
		t.Fatalf("failed convert ByteToJson ")
	}
	t.Log("success ByteToJson")
}
