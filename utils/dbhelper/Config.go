package dbhelper

import (
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DbType int

const (
	Postgres DbType = iota
	Oracle
	MySql
)

type DbConfig struct {
	ConnectionString string
	DbType           DbType
	gorm.Config      //扩展gorm的配置
}
type DbInterface interface {
	//打开数据库连接
	Open() bool
	//关闭数据库连接
	Close() bool
	SelectById(bean ModelInterface, id any) Model
	SelectByCondition(bean ModelInterface, conds ...any) Model
}

type DbBase struct {
	//DbInterface
	Config DbConfig
	db     *gorm.DB
}

// 参考： https://gorm.io/docs/gorm_config.html
func CreateDb(connectionString string, dbType DbType) DbInterface {
	config := DbConfig{ConnectionString: connectionString, DbType: dbType}
	config.SkipDefaultTransaction = true
	config.NamingStrategy = schema.NamingStrategy{
		TablePrefix: "t_",
	}
	switch dbType {
	case Postgres:
		pg := NewPostgresDB(config)
		return pg
		break
	default:
		break
	}
	loghelper.GetLogManager().Error("暂不支持的数据库类型！！")
	return nil
}
