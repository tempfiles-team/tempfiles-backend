package response

type Response struct {
	// "요청이 정상적으로 처리되었습니다" | "서버에서 일시적인 오류가 발생했어요"
	Message string `json:"message"`
	// SUCCESS | FAIL
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func NewResponse(message string, status string, data interface{}) *Response {
	return &Response{
		Message: message,
		Status:  status,
		Data:    data,
	}
}

func NewSuccessMessageResponse(message string) *Response {
	return NewResponse(message, "SUCCESS", nil)
}

func NewFailMessageResponse(message string) *Response {
	return NewResponse(message, "FAIL", nil)
}

func NewSuccessDataResponse(data interface{}) *Response {
	return NewResponse("요청이 정상적으로 처리되었습니다", "SUCCESS", data)
}

func NewFailDataResponse(data interface{}) *Response {
	return NewResponse("서버에서 일시적인 오류가 발생했어요", "FAIL", data)
}
