package web

type ResponseResult[T any] struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,optional"`
	Data T      `json:"data,optional"`
}

type ResponseString struct {
	ResponseResult[string]
}

type ResponseBytes struct {
	ResponseResult[[]byte]
}

type ResponseBool struct {
	ResponseResult[bool]
}

type ResponseFloat struct {
	ResponseResult[float32]
}

type ResponseDouble struct {
	ResponseResult[float64]
}

type ResponseInt struct {
	ResponseResult[int]
}

type ResponseInt8 struct {
	ResponseResult[int8]
}
type ResponseInt16 struct {
	ResponseResult[int16]
}
type ResponseInt64 struct {
	ResponseResult[int64]
}

type ResponseModel struct {
	ResponseResult[interface{}]
}

type ResponseJsonArray struct {
	ResponseResult[interface{}]
}

type ResponseJsonObject struct {
	ResponseResult[interface{}]
}
type ResponseJsonBody struct {
	ResponseResult[interface{}]
}
