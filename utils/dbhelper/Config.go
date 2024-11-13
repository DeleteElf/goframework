package dbhelper

import (
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	//开始事务，开始事务后，连接会始终打开，直到调用提交事务、回滚事务、关闭连接等操作
	BeginTransaction() bool
	//提交事务
	CommitTransaction() bool
	//回滚事务
	RollbackTransaction() bool
	//是否处于事务中
	IsInTransaction() bool
	//自动更新表结构，慎重使用此方法
	AutoMigrate(model ModelInterface) error
	//保存或更新数据
	Save(model ModelInterface)
	//根据id查询对象数据
	SelectById(model ModelInterface, id any)
	//根据条件查询对象数据
	SelectByCondition(dest interface{}, condition string, conds ...any)
	//查询数据
	QueryData(sql string, conds ...any) *DataTable
}

type DataTable struct {
	Rows []map[string]interface{}
}

type DbBase struct {
	//DbInterface
	Config DbConfig
	db     *gorm.DB
}

// 参考： https://gorm.io/docs/gorm_config.html
func CreateDb(connectionString string, dbType DbType, logLevel logger.LogLevel) DbInterface {
	config := DbConfig{ConnectionString: connectionString, DbType: dbType}
	config.SkipDefaultTransaction = true
	config.NamingStrategy = schema.NamingStrategy{
		TablePrefix: "t_",
	}
	config.Logger = logger.Default.LogMode(logLevel)
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
