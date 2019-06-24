package tools

import (
	"encoding/json"
)

// 返回 JSON 字符串
func ReturnJson(o interface{}) string {
	data, _ := json.Marshal(o)
	return string(data)
}
