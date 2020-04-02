package tools

import "encoding/json"

func ToJsonString(value interface{}) string {
	str, _ := json.Marshal(value)
	return string(str)
}
