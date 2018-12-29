package models

type ReturnMessage struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func Default() *ReturnMessage {
	return &ReturnMessage{
		Result:  true,
		Message: "It's OK!",
	}
}