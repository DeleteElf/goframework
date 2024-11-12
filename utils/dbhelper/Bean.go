package dbhelper

type IdData interface {
	string | ~int // int | uint | int32 | uint32 | int64 | uint64 //id支持的类型
}

type Bean interface {
	//TableName() string //如果强制要求每个对象都必须书写映射，则取消此注释
}

type BeanBase[T IdData] struct {
	Bean
	Id     T    `gorm:"column:f_id;primaryKey"` //默认会使用Id作为主键
	Active bool `gorm:"column:f_active;default:true"`
}

type Entry struct {
	Name string `gorm:"column:f_name"`
}

type Parent[T IdData] struct {
	Parent T `gorm:"column:f_parent_id"`
}

// 系统用户
type UserInfo struct {
	BeanBase[int]
	Entry
	Account  string  `gorm:"column:f_account"`
	Password string  `gorm:"column:f_password"`
	Email    *string `gorm:"column:f_email"` //定义指针是为了支持空值
}

// 系统用户的表名定义
func (UserInfo) TableName() string {
	return "t_user_info"
}
