package entities

// IObject 对象基类
type IObject interface{}

// IModel 数据模型基类
type IModel struct {
	IObject
}

type IConfig interface {
	IObject
}

// BaseConfig 配置对象基类
type BaseConfig struct {
	IConfig `json:"-"`
	Name    string `json:",omitempty"`
	Enable  bool   `json:",default=true"`
}
