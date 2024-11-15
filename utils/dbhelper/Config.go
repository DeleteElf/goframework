package dbhelper

import (
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	//gorm.Config      //扩展gorm的配置
	SkipDefaultTransaction bool
	SafeColumn             bool
	LogLevel               logger.LogLevel
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

// 数据表
type DataTable struct {
	//行数据集合，key在配置SafeColumn=true时，仅支持驼峰命名法，写入和读取都会按ColumnPrefix自动转换，配置SafeColumn=false时，按key直接映射数据库字段
	Rows []map[string]interface{}
	//数据的表名，仅当写入数据时才需要,填写完成的表名
	TableName string
	//主键字段名称，如果没有设置，则默认为f_id字段
	PkColumnName string
	//字段的前缀，如果没有配置，则默认为为f_
	ColumnPrefix string
	//创建时间的字段名，默认为f_create_time
	CreateTimeColumn string
	//修改时间的字段名，默认为f_modify_time
	ModifyTimeColumn string
}
type DbBase struct {

	//DbInterface
	Config          DbConfig
	db              *gorm.DB
	isInTransaction bool
}

//
//type MyNamingStrategy struct {
//	schema.NamingStrategy
//
//	ColumnPrefix string
//}
//
//// 重写字段名规则，使其支持f_
//func (ns MyNamingStrategy) ColumnName(table, column string) string {
//	return ns.ColumnPrefix + ns.toDBName(column) //不用进行复数转换
//}
//
//var (
//	// https://github.com/golang/lint/blob/master/lint.go#L770
//	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
//	commonInitialismsReplacer *strings.Replacer
//)
//
//func init() {
//	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
//	for _, initialism := range commonInitialisms {
//		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, cases.Title(language.Und).String(initialism))
//	}
//	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
//}
//
//func (ns MyNamingStrategy) toDBName(name string) string {
//	if name == "" {
//		return ""
//	}
//
//	if ns.NameReplacer != nil {
//		tmpName := ns.NameReplacer.Replace(name)
//
//		if tmpName == "" {
//			return name
//		}
//
//		name = tmpName
//	}
//
//	if ns.NoLowerCase {
//		return name
//	}
//
//	var (
//		value                          = commonInitialismsReplacer.Replace(name)
//		buf                            strings.Builder
//		lastCase, nextCase, nextNumber bool // upper case == true
//		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
//	)
//
//	for i, v := range value[:len(value)-1] {
//		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
//		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'
//
//		if curCase {
//			if lastCase && (nextCase || nextNumber) {
//				buf.WriteRune(v + 32)
//			} else {
//				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
//					buf.WriteByte('_')
//				}
//				buf.WriteRune(v + 32)
//			}
//		} else {
//			buf.WriteRune(v)
//		}
//
//		lastCase = curCase
//		curCase = nextCase
//	}
//
//	if curCase {
//		if !lastCase && len(value) > 1 {
//			buf.WriteByte('_')
//		}
//		buf.WriteByte(value[len(value)-1] + 32)
//	} else {
//		buf.WriteByte(value[len(value)-1])
//	}
//	ret := buf.String()
//	return ret
//}

// 参考： https://gorm.io/docs/gorm_config.html
func CreateDb(connectionString string, dbType DbType, logLevel logger.LogLevel) DbInterface {
	config := DbConfig{ConnectionString: connectionString, DbType: dbType}
	config.LogLevel = logLevel
	return CreateDbByConfig(config)
}

func CreateDbByConfig(config DbConfig) DbInterface {
	config.SkipDefaultTransaction = true
	switch config.DbType {
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
