package structs

// 发送消息结构体
type SendMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
