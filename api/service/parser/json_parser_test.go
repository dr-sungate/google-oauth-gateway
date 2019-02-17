package parser

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
	t.Log(resultmap)
	if reflect.DeepEqual(inputdata1, resultmap) == false {
		t.Fatalf("failed ByteToJson ")
	}
	t.Log("success ByteToJson")
}

func Test_JsonToByte_ByteToJson2(t *testing.T) {
	var inputdata1 map[string]interface{}
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
	t.Log(resultmap)
	if resultmap != nil {
		t.Fatalf("failed ByteToJson ")
	}
	t.Log("success ByteToJson")
}

func Test_GetJsonItems(t *testing.T) {
	inputdata1 := `{
	"message1": "ok",
	"message2": "ng",
	"message3": [
		{"message3_1": "data1"},
		{"message3_2": "data2"}
	]
}`

	getitem, err := GetJsonItems([]byte(inputdata1), "message1")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(getitem)
	if getitem != "ok" {
		t.Fatalf("failed GetJsonItems ")
	}
	t.Log("success GetJsonItems")
	getitem2, err := GetJsonItems([]byte(inputdata1), "message4")
	if err == nil {
		t.Fatalf("failed test invalid GetJsonItems")
	}
	t.Log(getitem2)
	t.Log(err)
	t.Log("success invalid GetJsonItems")

	resultmap := make([]interface{}, 0)
	data1 := map[string]interface{}{"message3_1": "data1"}
	data2 := map[string]interface{}{"message3_2": "data2"}
	resultmap = append(resultmap, data1)
	resultmap = append(resultmap, data2)
	getitem3, err := GetJsonItems([]byte(inputdata1), "message3")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if reflect.DeepEqual(getitem3, resultmap) == false {
		t.Fatalf("failed array GetJsonItems ")
	}
	t.Log("success array GetJsonItems")
}

func Test_RemoveJsonItems(t *testing.T) {
	inputdata := `{
	"message1": "ok",
	"message2": "ng",
	"message3": [
		{"message3_1": "data1"},
		{"message3_2": "data2"}
	]
}`

	resultdata1 := `{
	"message2": "ng",
	"message3": [
		{"message3_1": "data1"},
		{"message3_2": "data2"}
	]
}`

	resultdata2 := `{
	"message1": "ok",
	"message2": "ng"
}`

	inputdatamap, err := ByteToJson([]byte(inputdata))
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	resultdatamap1, err := ByteToJson([]byte(resultdata1))
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	resultdatamap2, err := ByteToJson([]byte(resultdata2))
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	removeditem1, err := RemoveJsonItems(inputdatamap, []string{"message1"})
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(removeditem1)
	if reflect.DeepEqual(removeditem1, resultdatamap1) == false {
		t.Fatalf("failed RemoveJsonItems 1")
	}
	removeditem2, err := RemoveJsonItems(inputdatamap, []string{"message3"})
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(removeditem2)
	if reflect.DeepEqual(removeditem2, resultdatamap2) == false {
		t.Fatalf("failed RemoveJsonItems 2")
	}

	t.Log("success RemoveJsonItems")
}

func Test_AddJsonItems(t *testing.T) {
	inputdata := `{

	"message2": "ng"
}`

	resultdata1 := `{
    "message1": "ok",
    "message2": "ng"
}`

	resultdata2 := `{
    "message2": "ng",
    "message3": [
        {
            "message3_1": "data1"
        },
        {
            "message3_2": "data2"
        }
    ]
}`

	addmap1 := make([]map[string]interface{}, 0)
	addmapone1 := map[string]interface{}{"message1": "ok"}
	addmap1 = append(addmap1, addmapone1)

	additem1, err := AddJsonItems([]byte(inputdata), addmap1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(string(additem1))
	if reflect.DeepEqual(additem1, []byte(resultdata1)) == false {
		t.Fatalf("failed AddJsonItems 1")
	}

	addmap2 := make([]map[string]interface{}, 0)
	addarraymap := make([]interface{}, 0)
	data1 := map[string]interface{}{"message3_1": "data1"}
	data2 := map[string]interface{}{"message3_2": "data2"}
	addarraymap = append(addarraymap, data1)
	addarraymap = append(addarraymap, data2)
	addmapone2 := map[string]interface{}{"message3": addarraymap}
	addmap2 = append(addmap2, addmapone2)

	additem2, err := AddJsonItems([]byte(inputdata), addmap2)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(string(additem2))
	if reflect.DeepEqual(additem2, []byte(resultdata2)) == false {
		t.Fatalf("failed AddJsonItems 2")
	}

	t.Log("success AddJsonItems")
}
