package web

type IdRequest struct {
	Id string `path:"id"` //设置请求的格式要求
}
type Response[T any] struct {
	Code string `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data T      `json:"data,omitempty"`
}

type ResponseString struct {
	Response[string]
}

type ResponseBytes struct {
	Response[[]byte]
}

type ResponseBool struct {
	Response[bool]
}

type ResponseFloat struct {
	Response[float64]
}

type ResponseInt struct {
	Response[int]
}

type ResponseInt8 struct {
	Response[int8]
}
type ResponseInt16 struct {
	Response[int16]
}
type ResponseInt64 struct {
	Response[int64]
}

type ResponseModel struct {
	Response[interface{}]
}
