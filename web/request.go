package web

// IdRequestPath 路径id参数的支持
type IdRequestPath struct {
	Id string `path:"id"` //设置请求的格式要求
}

// TypeIdRequestPath 路径type/id参数的支持
type TypeIdRequestPath struct {
	IdRequestPath
	Type string `path:"type"`
}

// RequestFormId 支持验证的http request请求id
type RequestFormId struct {
	Id string `form:"id" validate:“required”`
}

// RequestJsonId 支持验证的http json request请求id
type RequestJsonId struct {
	Id string `json:"id" validate:“required”`
}
