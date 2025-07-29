package response

type ErrorResponseSchema struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}

type IResponse[T any] interface {
	Status(code int) Response[T]
}

func New[T any](data T) Response[T] {
	return Response[T]{Code: 200, Data: data}
}

func (r Response[T]) Status(code int) Response[T] {
	r.Code = code
	return r
}
