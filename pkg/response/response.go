package response

import "fmt"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(message string) Response {
	return Response{Code: 200, Message: message}
}

func (r Response) Status(code int) Response {
	r.Code = code
	return r
}

func (r Response) ToString() string {
	resp := fmt.Sprintf(`{"code":%d,"message":"%s"}`, r.Code, r.Message)

	return resp
}
