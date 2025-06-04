package entities

// IObject 对象基类
type IObject struct{}

// IModel 数据模型基类
type IModel struct {
	IObject
}

// IConfig 配置对象基类
type IConfig struct {
	IObject
}
