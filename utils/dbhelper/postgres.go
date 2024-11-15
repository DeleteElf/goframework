package dbhelper

import (
	"errors"
	"fmt"
	"github.com/deleteelf/goframework/utils/loghelper"
	"github.com/deleteelf/goframework/utils/stringhelper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type PostgresDB struct {
	DbBase
}

// 默认会自动打开和关闭连接，如果启用了事务，则在提交或回滚事务之前，不会关闭连接
func NewPostgresDB(config DbConfig) *PostgresDB {
	return &PostgresDB{
		DbBase{
			Config: config,
		},
	}
}

func (pg *PostgresDB) Open() bool {
	var err error
	//conf := MyNamingStrategy{
	//	ColumnPrefix: "f_",
	//}
	//conf.TablePrefix = "t_"
	//conf.IdentifierMaxLength = 64
	//conf.SingularTable = false
	pg.db, err = gorm.Open(postgres.Open(pg.Config.ConnectionString), &gorm.Config{
		SkipDefaultTransaction: pg.Config.SkipDefaultTransaction,
		Logger:                 logger.Default.LogMode(pg.Config.LogLevel),
		NamingStrategy:         schema.NamingStrategy{},
	})
	if err != nil {
		loghelper.GetLogManager().ErrorFormat("数据库连接失败！！%s", pg.Config.ConnectionString)
		return false
	}
	loghelper.GetLogManager().InfoFormat("数据库连接成功！！%s", pg.Config.ConnectionString)
	return true
}

func (pg *PostgresDB) Close() bool {
	sqlDb, err := pg.db.DB()
	if err != nil {
		loghelper.GetLogManager().Error(err)
		return false
	}
	err = sqlDb.Close()
	if err != nil {
		loghelper.GetLogManager().Error(err)
		return false
	}
	return true
}

func (pg *PostgresDB) BeginTransaction() bool {
	if !pg.isInTransaction {
		pg.db.Begin()
		pg.isInTransaction = true
		return true
	}
	return false
}

// 提交事务
func (pg *PostgresDB) CommitTransaction() bool {
	if pg.isInTransaction {
		pg.db.Commit()
		pg.isInTransaction = false
		return true
	}
	return false
}

// 回滚事务
func (pg *PostgresDB) RollbackTransaction() bool {
	if pg.isInTransaction {
		pg.db.Rollback()
		pg.isInTransaction = false
		return true
	}
	return false
}

// 是否处于事务中
func (pg *PostgresDB) IsInTransaction() bool {
	return pg.isInTransaction
}

func (pg *PostgresDB) AutoMigrate(model ModelInterface) error {
	if pg.Open() {
		defer pg.Close()
		return pg.db.AutoMigrate(model)
	}
	return nil
}
func (pg *PostgresDB) Save(model ModelInterface) {
	if pg.Open() {
		defer pg.Close()
		pg.db.Save(model)
	}
}

func (pg *PostgresDB) SelectById(model ModelInterface, id any) {
	//反射的案例，不过gorm已经做好反射了
	//t := reflect.TypeFor[T1]()
	//val := reflect.New(t).Elem()
	//result := val.Interface().(T1)
	if pg.Open() {
		defer pg.Close()
		err := pg.db.First(model, id).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("根据Id查询的数据不存在！！！")
			break
		default:
			break
		}

	}
}

// 根据条件查询数据，dest传入数组指针
func (pg *PostgresDB) SelectByCondition(dest interface{}, query string, conds ...any) {
	if pg.Open() {
		defer pg.Close()
		err := pg.db.Where(query, conds...).Find(dest).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		default:
			break
		}
	}
}

// 使用原始sql语句查询数据，支持通过配置SafeColumn进行数据字段保护，自动转化为驼峰命名法的字段
func (pg *PostgresDB) QueryData(sql string, conds ...any) *DataTable {
	result := new(DataTable)
	if pg.Open() {
		defer pg.Close()
		ctx := pg.db.Raw(sql, conds...)
		if pg.Config.SafeColumn {
			rows, err := pg.db.Raw(sql, conds...).Rows()
			defer rows.Close()
			if err != nil {
				loghelper.GetLogManager().Error("获取行数据出错！！")
			}
			columns, err1 := rows.Columns()
			if err1 != nil {
				loghelper.GetLogManager().Error("获取列数据出错！！")
			}
			ctx.Statement.ColumnMapping = map[string]string{}
			for _, column := range columns {
				ctx.Statement.ColumnMapping[column] = stringhelper.ConvertToCamel(column)
			}
		}
		err := ctx.Scan(&result.Rows).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		case gorm.ErrDryRunModeUnsupported:
			loghelper.GetLogManager().Error("ErrDryRunModeUnsupported！！！")
			break
		default:
			break
		}
	}
	return result
}

// 保存数据表，保存前，需要在数据表中设置表名和主键，仅支持单表数据更新，且数据必须包含主键数据
// 如需事务支持，请在调用此方法前开启事务，并在完成此方法后，提交或回归事务
func (pg *PostgresDB) SaveData(dataTale DataTable) (int, error) {
	if dataTale.Rows == nil || len(dataTale.Rows) == 0 {
		return 0, nil
	}
	if len(dataTale.TableName) == 0 {
		return 0, errors.New("未设置表名")
	}
	if len(dataTale.PkColumnName) == 0 {
		dataTale.PkColumnName = "f_id"
	}
	dataTale.PkColumnName = strings.ToLower(dataTale.PkColumnName)

	if len(dataTale.CreateTimeColumn) == 0 {
		dataTale.CreateTimeColumn = "f_create_time"
	}
	dataTale.CreateTimeColumn = strings.ToLower(dataTale.CreateTimeColumn)

	if len(dataTale.ModifyTimeColumn) == 0 {
		dataTale.ModifyTimeColumn = "f_modify_time"
	}
	dataTale.ModifyTimeColumn = strings.ToLower(dataTale.ModifyTimeColumn)

	if len(dataTale.ColumnPrefix) == 0 {
		dataTale.ColumnPrefix = "f_"
	}
	sqlFormat := "select * from %s where %s=?"
	updateFormat := "update %s set %s where %s=?"
	insertFormat := "insert into %s (%s) value (%s)"
	query := fmt.Sprintf(sqlFormat, dataTale.TableName, dataTale.PkColumnName)
	idKey := stringhelper.ConvertToCamel(dataTale.PkColumnName)
	for i, row := range dataTale.Rows {
		//从数据库检索出数据表对应的数据，以决定是新增还是修改
		isInsert := false
		if row[idKey] == nil {
			isInsert = true
		} else {
			dt := pg.QueryData(query, row[idKey])
			if len(dt.Rows) == 0 { //新增
				isInsert = true
			}
		}
		var columnStr string
		count := len(row)
		rowData := make([]any, count)
		var index int
		if isInsert { //新增
			var valueStr string
			for key, value := range row {
				columnName := strings.ToLower(dataTale.ColumnPrefix + stringhelper.ConvertCamelToSnakeWithDefault(key))
				if len(columnStr) != 0 {
					columnStr += ","
					valueStr += ","
				}
				columnStr += columnName
				valueStr += "?"
				if columnName == dataTale.CreateTimeColumn || columnName == dataTale.ModifyTimeColumn { //接管，不允许外部写入
					rowData[index] = time.Now()
				} else {
					rowData[index] = value
				}
				index++
			}
			sql := fmt.Sprintf(insertFormat, dataTale.TableName, columnStr, valueStr)
			ctx := pg.db.Exec(sql, rowData...)
			if ctx.Error != nil {
				return i, ctx.Error
			}
		} else {
			for key, value := range row {
				columnName := strings.ToLower(dataTale.ColumnPrefix + stringhelper.ConvertCamelToSnakeWithDefault(key))
				if columnName == dataTale.CreateTimeColumn {
					continue
				}
				if columnName == dataTale.PkColumnName {
					rowData[count-1] = value //主键条件直接写到最后一个参数
					continue
				}
				if len(columnStr) != 0 {
					columnStr += ","
				}
				columnStr += columnName + "=?"
				if columnName == dataTale.ModifyTimeColumn { //接管，不允许外部写入
					rowData[index] = time.Now()
				} else {
					rowData[index] = value
				}
				index++
			}
			sql := fmt.Sprintf(updateFormat, dataTale.TableName, columnStr, dataTale.PkColumnName)
			ctx := pg.db.Exec(sql, rowData...)
			if ctx.Error != nil {
				return i, ctx.Error
			}
		}
	}
	return len(dataTale.Rows), nil
}
