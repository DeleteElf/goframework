package entities

import (
	"github.com/deleteelf/goframework/utils/stringhelper"
	"reflect"
	"time"
)

type IdData interface {
	string | int | uint | int32 | uint32 | int64 | uint64 //id支持的类型
}

type ModelInterface interface {
	TableName() string //如果强制要求每个对象都必须书写映射，则取消此注释
}

//
//type ParentChildrenInterface interface {
//}

type Model struct {
}

func (model Model) TableName() string {
	t := reflect.TypeOf(model)
	return "t_" + stringhelper.ConvertCamelToSnakeWithDefault(t.Name())
}

type Bean[T IdData] struct {
	Model
	Id         T         `gorm:"column:f_id;primaryKey" json:"id" xml:"id,attr"` //默认会使用Id作为主键
	Active     bool      `gorm:"column:f_active;default:true" json:"active" xml:"active,attr"`
	CreateTime time.Time `gorm:"column:f_create_time;default:now()" json:"createTime" xml:"createTime,attr"` //默认当前时间
	ModifyTime time.Time `gorm:"column:f_modify_time;default:now()" json:"modifyTime" xml:"modifyTime,attr"`
	Remark     *string   `gorm:"column:f_remark" json:"remark" xml:"remark,cdata"` //定义指针是为了支持空值
}

type Entity struct {
	Name string `gorm:"column:f_name;type:varchar(20);" json:"name" xml:"name,attr"` //定义有名称的实体
}

type Parent[T IdData] struct {
	//ParentChildrenInterface
	Parent T `gorm:"column:f_parent_id" json:"parentId"  xml:"parentId,attr"` //定义有父子关系的结构
}

// 系统用户
type UserInfo struct {
	Bean[int]         //匿名扩展
	Entity            //扁平式扩展，而非继承
	Account   string  `gorm:"column:f_account;type:varchar(20);not null" json:"account"`
	Password  string  `gorm:"column:f_password;type:varchar(35);not null" json:"-"`  //密码不允许对外输出
	Email     *string `gorm:"column:f_email;type:varchar(30);" json:"email"`         //定义指针是为了支持空值
	Telephone *string `gorm:"column:f_telephone;type:varchar(20);" json:"telephone"` //定义有名称的实体
}

// 系统用户的表名定义
func (UserInfo) TableName() string {
	return "sys.t_user_info"
}
