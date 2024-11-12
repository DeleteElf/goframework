package dbhelper

import (
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm"
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
	SelectById(bean Bean, id interface {
		string | int | uint | int32 | uint32 | int64 | uint64 //id支持的类型
	}) Bean
	SelectByCondition(bean Bean, conds ...any) Bean
}

type DbBase struct {
	DbInterface
	Config DbConfig
	db     *gorm.DB
}

func CreateDb(connectionString string, dbType DbType) DbInterface {
	config := DbConfig{ConnectionString: connectionString, DbType: dbType}
	config.SkipDefaultTransaction = true
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