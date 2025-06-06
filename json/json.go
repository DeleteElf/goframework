package json

import "encoding/json"

// JsonObject 定义JsonObject类型
type JsonObject map[string]interface{}

// JsonArray 定义JsonArray类型
type JsonArray []interface{}

//  func GetJsonObjectString(jsonObject JsonObject) ([]byte, error) {
//	return json.Marshal(jsonObject)
//}
//
//func GetJsonArrayString(jsonArray JsonArray) ([]byte, error) {
//	return json.Marshal(jsonArray)
//}

// GetJsonObject 转成json对象
func GetJsonObject(str []byte) (JsonObject, error) {
	result := JsonObject{}
	err := GetObject(str, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetJsonArray 转成json数组
func GetJsonArray(str []byte) (JsonArray, error) {
	result := JsonArray{}
	err := GetObject(str, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetObject 数据转成对象
func GetObject(str []byte, v any) error {
	return json.Unmarshal(str, v)
}

// ToJsonString 对象转成json
func ToJsonString(v any) ([]byte, error) {
	return json.Marshal(v)
}
